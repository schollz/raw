package supercollider

import (
	_ "embed"
	"os/exec"
	"runtime"
	"time"
)

//go:embed raw.sc
var sclangCode string

var sclang = "sclang"
var cmdSuperCollider *exec.Cmd
var ready = false

func init() {
	if runtime.GOOS == "windows" {
		// TODO: search the program files for the "SuperCollider" sclang binary
	}
}

func blockUntilReady() {
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
	go func() {
		cmdSuperCollider = exec.Command("sclang") // TODO: run sclang on the embeded file
		err = cmdSuperCollider.Start()
	}()
	go func() {
		for {
			// check whether the "ready" file exists from the sclang server
			time.Sleep(100 * time.Millisecond)
		}
	}()
	return
}

// Stop will close supercollider
func Stop() (err error) {
	if cmdSuperCollider != nil {
		err = cmdSuperCollider.Process.Kill()
	}
	return
}
