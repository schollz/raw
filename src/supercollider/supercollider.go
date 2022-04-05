package supercollider

import (
	"bufio"
	"crypto/rand"
	_ "embed"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	mathr "math/rand"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/hypebeast/go-osc/osc"
	ps "github.com/mitchellh/go-ps"
	log "github.com/schollz/logger"
	psprocess "github.com/shirou/gopsutil/process"
)

//go:embed raw.sc
var sclangCode string

var sclang = "sclang"
var sclangFolder = "."
var sccodeFile = ""
var cmdSuperCollider *exec.Cmd
var ready = false
var starting = false
var READYFILE = ""
var mu sync.Mutex

// TempDir is where the temporary intermediate files are held
var TempDir = os.TempDir()

// TempPrefix is a unique indicator of the temporary files
var TempPrefix = "scstuff"

// TempType is the type of file to be generated (should be "wav")
var TempType = "wav"

func tmpfile() string {
	randBytes := make([]byte, 16)
	rand.Read(randBytes)
	return filepath.Join(TempDir, TempPrefix+hex.EncodeToString(randBytes)+"."+TempType)
}

func init() {
	randBytes := make([]byte, 16)
	rand.Read(randBytes)
	READYFILE = filepath.FromSlash(path.Join(os.TempDir(), hex.EncodeToString(randBytes)+".sc"))
	if runtime.GOOS == "windows" {
		folders, err := filepath.Glob("C:\\Program Files\\SuperCollider*")
		if err != nil {
			panic(err)
		}
		for _, folder := range folders {
			files, err := filepath.Glob(path.Join(folder, "sclang*"))
			if err != nil {
				panic(err)
			}
			for _, fname := range files {
				_, sclang = filepath.Split(fname)
				sclangFolder = folder
			}
		}
	}
	sclangCode = strings.Replace(sclangCode, "/tmp/nrt-scready", strings.Replace(READYFILE, `\`, `\\`, -1), -1)
	log.Debugf("using sclang: '%s'", sclang)

	processList, err := ps.Processes()
	if err != nil {
		log.Debug("ps.Processes() Failed, are you using windows?")
		return
	}
	// map ages
	for x := range processList {
		var process ps.Process
		process = processList[x]
		if strings.Contains(process.Executable(), "sclang") {
			panic("sclang already running, exit first")
		}
		if strings.Contains(process.Executable(), "scsynth") {
			panic("scsynth already running, exit first")
		}
	}
}

func scPath(f string) string {
	fabs, _ := filepath.Abs(f)
	fabs = filepath.FromSlash(fabs)
	fabs = strings.Replace(fabs, `\`, `\\`, -1)
	return fabs
}

func Effect(fname string, effect string, fs ...float64) (fname2 string, err error) {
	blockUntilReady()
	log.Tracef("getting lock for %s", fname)
	mu.Lock()
	defer mu.Unlock()
	log.Tracef("got lock for %s", fname)
	effectF := []float64{0, 0, 0, 0}
	for i, f := range fs {
		if i < 4 {
			effectF[i] = f
		}
	}
	fname2 = tmpfile()
	durationScaling := "1"
	if effect == "reverberate" {
		if mathr.Float64() < 0.5 {
			durationScaling = "2"
		} else {
			durationScaling = "2.5"
		}
	}
	scDoneFile := filepath.FromSlash(tmpfile())
	scDoneFile = strings.Replace(scDoneFile, `\`, `\\`, -1)
	defer os.Remove(scDoneFile)

	client := osc.NewClient("localhost", 47113)
	msg := osc.NewMessage("/score")
	msg.Append(fname)
	msg.Append(fname2)
	msg.Append(effect)
	msg.Append(durationScaling)
	for _, f := range effectF {
		msg.Append(float32(f))
	}
	msg.Append(scDoneFile)
	err = client.Send(msg)
	if err != nil {
		log.Error(err)
		return
	}

	// wait for the file to be written
	for {
		time.Sleep(25 * time.Millisecond)
		if _, err := os.Stat(scDoneFile); err == nil {
			log.Debug("supercollider is done")
			break
		}
	}

	return
}

func blockUntilReady() {
	Start()
	if ready {
		return
	}
	for {
		if ready {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func Start() (err error) {
	if starting {
		return
	}
	starting = true
	go func() {
		f, _ := ioutil.TempFile(os.TempDir(), "sccode")
		f.WriteString(sclangCode)
		f.Close()
		sccodeFile, _ = filepath.Abs(f.Name())
		sccodeFile = filepath.FromSlash(sccodeFile)
		cwd, _ := os.Getwd()
		if runtime.GOOS == "windows" {
			os.Chdir(sclangFolder)
		}
		cmdSuperCollider = exec.Command(sclang, sccodeFile) // TODO: run sclang on the embeded file
		log.Debug("starting supercollider")
		stdout, _ := cmdSuperCollider.StdoutPipe()
		stderr, _ := cmdSuperCollider.StderrPipe()
		err = cmdSuperCollider.Start()
		os.Chdir(cwd)
		go func() {
			scanner := bufio.NewScanner(stderr)
			for scanner.Scan() {
				m := scanner.Text()
				log.Tracef("%s", m)
			}
		}()
		go func() {
			scanner := bufio.NewScanner(stdout)
			for scanner.Scan() {
				m := scanner.Text()
				log.Tracef("%s", m)
			}
		}()
	}()
	go func() {
		log.Trace("watching for ready signal")
		for {
			// check whether the "ready" file exists from the sclang server
			if _, err := os.Stat(READYFILE); err == nil {
				log.Debug("supercollider is ready")
				ready = true
				break
			}
			time.Sleep(100 * time.Millisecond)
		}
	}()
	return
}

// Stop will close supercollider
func Stop() (err error) {
	if cmdSuperCollider != nil {
		err = cmdSuperCollider.Process.Kill()
		if err != nil {
			log.Error(err)
		}
		err = killProcess("scsynth")
		if err != nil {
			log.Error(err)
		}
		err = killProcess("sclang")
		if err != nil {
			log.Error(err)
		}
	}
	if sccodeFile != "" {
		os.Remove(sccodeFile)
	}
	os.Remove(READYFILE)
	err = clean()
	return
}

// Clean will remove files created after each function
func clean() (err error) {
	files, err := filepath.Glob(path.Join(TempDir, TempPrefix+"*."+TempType))
	if err != nil {
		return err
	}
	for _, fname := range files {
		log.Tracef("removing %s", fname)
		err = os.Remove(fname)
		if err != nil {
			return
		}
	}
	return
}

func killProcess(name string) error {
	processes, err := psprocess.Processes()
	if err != nil {
		return err
	}
	for _, p := range processes {
		n, err := p.Name()
		if err != nil {
			return err
		}
		if strings.Contains(n, name) {
			return p.Kill()
		}
	}
	return fmt.Errorf("process not found")
}
