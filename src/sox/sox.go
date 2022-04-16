package sox

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	log "github.com/schollz/logger"
	"github.com/schollz/progressbar/v3"
)

// TempDir is where the temporary intermediate files are held
var TempDir = os.TempDir()

// TempPrefix is a unique indicator of the temporary files
var TempPrefix = "sox"

// TempType is the type of file to be generated (should be "wav")
var TempType = "wav"

var GarbageCollection = false

func Tmpfile() string {
	randBytes := make([]byte, 16)
	rand.Read(randBytes)
	return filepath.Join(TempDir, TempPrefix+hex.EncodeToString(randBytes)+"."+TempType)
}

func init() {
	log.SetLevel("info")
	stdout, _, _ := run("sox", "--help")
	if !strings.Contains(stdout, "SoX") {
		panic("need to install sox")
	}
}

func run(args ...string) (string, string, error) {
	log.Trace(strings.Join(args, " "))
	baseCmd := args[0]
	cmdArgs := args[1:]
	cmd := exec.Command(baseCmd, cmdArgs...)
	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb
	err := cmd.Run()
	if err != nil {
		log.Errorf("%s -> '%s'", strings.Join(args, " "), err.Error())
		log.Error(outb.String())
		log.Error(errb.String())
	}
	return outb.String(), errb.String(), err
}

// MustString returns only the first argument of any function, as a string
func MustString(t ...interface{}) string {
	if len(t) > 0 {
		return t[0].(string)
	}
	return ""
}

// MustFloat returns only the first argument of any function, as a float
func MustFloat(t ...interface{}) float64 {
	if len(t) > 0 {
		return t[0].(float64)
	}
	return 0.0
}

// Clean will remove files created after each function
func Clean() (err error) {
	files, err := filepath.Glob(path.Join(TempDir, TempPrefix+"*."+TempType))
	if err != nil {
		return err
	}
	bar := progressbar.NewOptions(len(files),
		progressbar.OptionSetDescription("cleaning"),
		progressbar.OptionSetPredictTime(true))
	for _, fname := range files {
		bar.Add(1)
		log.Tracef("removing %s", fname)
		err = os.Remove(fname)
		if err != nil {
			return
		}
	}
	return
}

// Info returns the sample rate and number of channels for file
func Info(fname string) (samplerate int, channels int, err error) {
	stdout, stderr, err := run("sox", "--i", fname)
	if err != nil {
		return
	}
	stdout += stderr
	for _, line := range strings.Split(stdout, "\n") {
		if strings.Contains(line, "Channels") && channels == 0 {
			parts := strings.Fields(line)
			channels, err = strconv.Atoi(parts[len(parts)-1])
			if err != nil {
				return
			}
		} else if strings.Contains(line, "Sample Rate") && samplerate == 0 {
			parts := strings.Fields(line)
			samplerate, err = strconv.Atoi(parts[len(parts)-1])
			if err != nil {
				return
			}
		}
	}
	return
}

// Length returns the length of the file in seconds
func Length(fname string) (length float64, err error) {
	stdout, stderr, err := run("sox", fname, "-n", "stat")
	if err != nil {
		return
	}
	stdout += stderr
	for _, line := range strings.Split(stdout, "\n") {
		if strings.Contains(line, "Length") {
			parts := strings.Fields(line)
			length, err = strconv.ParseFloat(parts[len(parts)-1], 64)
			return
		}
	}
	return
}

// SilenceAppend appends silence to a file
func SilenceAppend(fname string, length float64) (fname2 string, err error) {
	samplerate, channels, err := Info(fname)
	if err != nil {
		return
	}
	silencefile := Tmpfile()
	defer os.Remove(silencefile)
	fname2 = Tmpfile()
	// generate silence
	_, _, err = run("sox", "-n", "-r", fmt.Sprint(samplerate), "-c", fmt.Sprint(channels), silencefile, "trim", "0.0", fmt.Sprint(length))
	if err != nil {
		return
	}
	// combine with original file
	_, _, err = run("sox", fname, silencefile, fname2)
	if err != nil {
		return
	}
	if GarbageCollection {
		os.Remove(fname)
	}
	return
}

// SilencePrepend prepends silence to a file
func SilencePrepend(fname string, length float64) (fname2 string, err error) {
	samplerate, channels, err := Info(fname)
	if err != nil {
		return
	}
	silencefile := Tmpfile()
	defer os.Remove(silencefile)
	fname2 = Tmpfile()
	// generate silence
	_, _, err = run("sox", "-n", "-r", fmt.Sprint(samplerate), "-c", fmt.Sprint(channels), silencefile, "trim", "0.0", fmt.Sprint(length))
	if err != nil {
		return
	}
	// combine with original file
	_, _, err = run("sox", silencefile, fname, fname2)
	if err != nil {
		return
	}
	return
}

// FFT
func FFT(fname string) (data string, err error) {
	_, data, err = run("sox", fname, "-n", "stat", "-freq")
	return
}

// Norm normalizes the audio
func Norm(fname string) (fname2 string, err error) {
	fname2 = Tmpfile()
	_, _, err = run("sox", fname, fname2, "norm")
	return
}

// SilenceTrim trims silence around a file
func SilenceTrim(fname string) (fname2 string, err error) {
	fname2 = Tmpfile()
	_, _, err = run("sox", fname, fname2, "silence", "1", "0.1", `0.025%`, "reverse", "silence", "1", "0.1", `0.25%`, "reverse")
	return
}

// SilenceTrimEnd trims silence from end of file
func SilenceTrimEnd(fname string) (fname2 string, err error) {
	fname2 = Tmpfile()
	_, _, err = run("sox", fname, fname2, "reverse", "silence", "1", "0.1", `0.25%`, "reverse")
	return
}

// Trim will trim the audio from the start point (with optional length)
func Trim(fname string, start float64, length ...float64) (fname2 string, err error) {
	fname2 = Tmpfile()
	if len(length) > 0 {
		_, _, err = run("sox", fname, fname2, "trim", fmt.Sprint(start), fmt.Sprint(length[0]))
	} else {
		_, _, err = run("sox", fname, fname2, "trim", fmt.Sprint(start))
	}
	return
}

// Reverse will reverse the audio
func Reverse(fname string) (fname2 string, err error) {
	fname2 = Tmpfile()
	_, _, err = run("sox", fname, fname2, "reverse")
	return
}

// Pitch repitched the audio
func Pitch(fname string, notes int) (fname2 string, err error) {
	fname2 = Tmpfile()
	_, _, err = run("sox", fname, fname2, "pitch", fmt.Sprintf("%d", notes*100))
	return
}

// Join will concatonate the files
func Join(fnames ...string) (fname2 string, err error) {
	fname2 = Tmpfile()
	fnames = append(fnames, fname2)
	_, _, err = run(append([]string{"sox"}, fnames...)...)
	return
}

// Mix will mix the files
func Mix(fnames ...string) (fname2 string, err error) {
	fname2 = Tmpfile()
	fnames = append(fnames, fname2)
	fnames = append(fnames, "norm")
	_, _, err = run(append([]string{"sox", "-m"}, fnames...)...)
	return
}

// Repeat will add n repeats to the audio
func Repeat(fname string, n int) (fname2 string, err error) {
	fname2 = Tmpfile()
	_, _, err = run("sox", fname, fname2, "repeat", fmt.Sprintf("%d", n))
	return
}

// RetempoSpeed will change the tempo of the audio and pitch
func RetempoSpeed(fname string, old_tempo float64, new_tempo float64) (fname2 string, err error) {
	fname2 = Tmpfile()
	_, _, err = run("sox", fname, fname2, "speed", fmt.Sprint(new_tempo/old_tempo), "rate", "-v", "48k")
	return
}

// RetempoStretch will change the tempo of the audio trying to keep pitch similar
func RetempoStretch(fname string, old_tempo float64, new_tempo float64) (fname2 string, err error) {
	fname2 = Tmpfile()
	_, _, err = run("sox", fname, fname2, "tempo", "-m", fmt.Sprint(new_tempo/old_tempo))
	return
}

// RetempoStretch will change the tempo of the audio trying to keep pitch similar
func Slowdown(fname string, slowdown float64) (fname2 string, err error) {
	fname2 = Tmpfile()
	_, _, err = run("sox", fname, fname2, "tempo", "-m", fmt.Sprint(slowdown))
	return
}

func CopyPaste(fname string, startPos float64, endPos float64, pastePos float64, crossfade float64, leeway0 ...float64) (fname2 string, err error) {
	copyLength := endPos - startPos
	if copyLength < 0.05 {
		fname2 = fname
		return
	}
	piece := Tmpfile()
	part1 := Tmpfile()
	part2 := Tmpfile()
	splice1 := Tmpfile()
	defer os.Remove(piece)
	defer os.Remove(part1)
	defer os.Remove(part2)
	defer os.Remove(splice1)
	fname2 = Tmpfile()
	leeway := 0.0
	if len(leeway0) > 0 {
		leeway = leeway0[0]
	}
	// 	os.cmd(string.format("sox %s %s trim %f %f",fname,piece,copy_start-e,copy_length+2*e))
	_, _, err = run("sox", fname, piece, "trim", fmt.Sprint(startPos-crossfade), fmt.Sprint(copyLength+2*crossfade))
	if err != nil {
		log.Error(err)
		fname2 = fname
		return
	}

	// 	os.cmd(string.format("sox %s %s trim 0 %f",fname,part1,paste_start+e))
	_, _, err = run("sox", fname, part1, "trim", "0", fmt.Sprint(pastePos+crossfade))
	if err != nil {
		log.Error(err)
		fname2 = fname
		return
	}

	// 	os.cmd(string.format("sox %s %s trim %f",fname,part2,paste_start+copy_length-e))
	_, _, err = run("sox", fname, part2, "trim", fmt.Sprint(pastePos+copyLength-crossfade))
	if err != nil {
		log.Error(err)
		fname2 = fname
		return
	}

	// 	os.cmd(string.format("sox %s %s %s splice %f,%f,%f",part1,piece,splice1,paste_start+e,e,l))
	_, _, err = run("sox", part1, piece, splice1, "splice", fmt.Sprintf("%f,%f,%f", pastePos+crossfade, crossfade, leeway))
	if err != nil {
		log.Error(err)
		fname2 = fname
		return
	}

	// 	os.cmd(string.format("sox %s %s %s splice %f,%f,%f",splice1,part2,fname2,paste_start+copy_length+e,e,l))
	_, _, err = run("sox", splice1, part2, fname2, "splice", fmt.Sprintf("%f,%f,%f", pastePos+copyLength+crossfade, crossfade, leeway))
	if err != nil {
		log.Error(err)
		fname2 = fname
		return
	}

	return
}

// Paste pastes any piece into a place in the audio, assumes that the piece has "crossfade" length on both sides
// in addition to its current length.
func Paste(fname string, piece string, pasteStart float64, crossfade float64) (fname2 string, err error) {
	copyLength, err := Length(piece)
	if err != nil {
		log.Error(err)
		fname2 = fname
		return
	}
	part1 := Tmpfile()
	part2 := Tmpfile()
	splice1 := Tmpfile()
	defer os.Remove(part1)
	defer os.Remove(part2)
	defer os.Remove(splice1)
	fname2 = Tmpfile()
	leeway := 0.0

	// 	os.cmd(string.format("sox %s %s trim 0 %f",fname,part1,paste_start+e))
	_, _, err = run("sox", fname, part1, "trim", "0", fmt.Sprint(pasteStart+crossfade))
	if err != nil {
		log.Error(err)
		fname2 = fname
		return
	}
	// copy(part1, "1.wav")

	// 	os.cmd(string.format("sox %s %s trim %f",fname,part2,paste_start+copy_length-e*3))
	_, _, err = run("sox", fname, part2, "trim", fmt.Sprint(pasteStart+copyLength-crossfade*3))
	if err != nil {
		log.Error(err)
		fname2 = fname
		return
	}
	// copy(part2, "2.wav")

	// 	os.cmd(string.format("sox %s %s %s splice %f,%f,%f",part1,piece,splice1,paste_start+e,e,l))
	_, _, err = run("sox", part1, piece, splice1, "splice", fmt.Sprintf("%f,%f,%f", pasteStart+crossfade, crossfade, leeway))
	if err != nil {
		log.Error(err)
		fname2 = fname
		return
	}
	// copy(splice1, "3.wav")

	// 	os.cmd(string.format("sox %s %s %s splice %f,%f,%f",splice1,part2,fname2,paste_start+copy_length+e,e,l))
	_, _, err = run("sox", splice1, part2, fname2, "splice", fmt.Sprintf("%f,%f,%f", pasteStart+copyLength+crossfade, crossfade, leeway))
	if err != nil {
		log.Error(err)
		fname2 = fname
		return
	}
	// copy(fname2, "4.wav")

	return
}

// SampleRate changes the sample rate
func SampleRate(fname string, srCh ...int) (fname2 string, err error) {
	sampleRate := int(48000)
	if len(srCh) > 0 {
		sampleRate = srCh[0]
	}
	fname2 = Tmpfile()
	_, _, err = run("sox", fname, fname2, "rate", fmt.Sprint(sampleRate))
	return
}

// PostProcess
func PostProcess(fname string, gain float64) (fname2 string, err error) {
	fname2 = Tmpfile()
	_, _, err = run("sox", fname, fname2, "reverse", "silence", "1", "0.1", `0.25%`, "reverse", "gain", fmt.Sprint(gain))
	return
}

// Gain applies gain
func Gain(fname string, gain float64) (fname2 string, err error) {
	fname2 = Tmpfile()
	_, _, err = run("sox", fname, fname2, "gain", fmt.Sprint(gain))
	return
}

// Stretch does a time stretch
func Stretch(fname string, stretch float64) (fname2 string, err error) {
	fname2 = Tmpfile()
	_, _, err = run("sox", fname, fname2, "stretch", fmt.Sprint(stretch))
	return
}

// Stutter does a stutter effect
func Stutter(fname string, stutter_length float64, pos_start float64, count float64, xfadePieceStutterGain ...float64) (fname2 string, err error) {
	crossfade_piece := 0.1
	crossfade_stutter := 0.005
	gain_amt := -2.0
	if count > 8 {
		gain_amt = -1.5
	}
	if len(xfadePieceStutterGain) > 0 {
		crossfade_piece = xfadePieceStutterGain[0]
	}
	if len(xfadePieceStutterGain) > 1 {
		crossfade_stutter = xfadePieceStutterGain[1]
	}
	if len(xfadePieceStutterGain) > 2 {
		gain_amt = xfadePieceStutterGain[2]
	}

	partFirst := Tmpfile()
	partMiddle := Tmpfile()
	partLast := Tmpfile()
	defer os.Remove(partFirst)
	defer os.Remove(partMiddle)
	defer os.Remove(partLast)

	// 	os.cmd(string.format("sox %s %s trim %f %f",fname,partFirst,pos_start-crossfade_piece,stutter_length+crossfade_piece+crossfade_stutter))
	_, _, err = run("sox", fname, partFirst, "trim",
		fmt.Sprint(pos_start-crossfade_piece), fmt.Sprint(stutter_length+crossfade_piece+crossfade_stutter))
	if err != nil {
		log.Error(err)
		fname2 = fname
		return
	}
	// 	os.cmd(string.format("sox %s %s trim %f %f",fname,partMiddle,pos_start-crossfade_stutter,stutter_length+crossfade_stutter+crossfade_stutter))
	_, _, err = run("sox", fname, partMiddle, "trim", fmt.Sprint(pos_start-crossfade_stutter),
		fmt.Sprint(stutter_length+crossfade_stutter+crossfade_stutter))
	if err != nil {
		log.Error(err)
		fname2 = fname
		return
	}
	// 	os.cmd(string.format("sox %s %s trim %f %f",fname,partLast,pos_start-crossfade_stutter,stutter_length+crossfade_piece+crossfade_stutter))
	_, _, err = run("sox", fname, partLast, "trim", fmt.Sprint(pos_start-crossfade_stutter),
		fmt.Sprint(stutter_length+crossfade_piece+crossfade_stutter))
	if err != nil {
		log.Error(err)
		fname2 = fname
		return
	}
	for i := 1.0; i <= count; i++ {
		fnameNext := ""
		if i == 1 {
			fnameNext, err = Gain(partFirst, gain_amt*(count-i))
			if err != nil {
				log.Errorf("stutter %f: %s", i, err.Error())
				fname2 = fname
				return
			}
		} else {
			fnameNext = Tmpfile()
			fnameMid := partLast
			if i < count {
				fnameMid = partMiddle
			}
			if gain_amt != 0 {
				var foo string
				foo, err = Gain(fnameMid, gain_amt*(count-i))
				if err != nil {
					log.Errorf("stutter %f: %s", i, err.Error())
					fname2 = fname
					return
				}
				fnameMid = foo
			}
			var fname2Length float64
			fname2Length, err = Length(fname2)
			if err != nil {
				log.Errorf("no length %f: %s", i, err.Error())
				fname2 = fname
				return
			}

			// os.cmd(string.format("sox %s %s %s splice %f,%f,0",fname2,fnameMid,fnameNext,audio.length(fname2),crossfade_stutter))
			_, _, err = run("sox", fname2, fnameMid, fnameNext, "splice", fmt.Sprintf("%f,%f,0",
				fname2Length, crossfade_stutter))
			if err != nil {
				log.Errorf("stutter %f: %s", i, err.Error())
				fname2 = fname
				return
			}
		}
		fname2 = fnameNext
	}
	return
}

// 	for i=1,count do
// 		local fnameNext=""
// 		if i==1 then
// 			fnameNext=audio.gain(partFirst,gain_amt*(count-i))
// 		else
// 			fnameNext=string.random_filename()
//          local fnameMid=i<count and partMiddle or partLast
//          if gain_amt~=0 then
//            fnameMid=audio.gain(fnameMid,gain_amt*(count-i))
//          end
// 			os.cmd(string.format("sox %s %s %s splice %f,%f,0",fname2,fnameMid,fnameNext,audio.length(fname2),crossfade_stutter))
// 		end
// 		fname2=fnameNext
// 	end
// 	return fname2
// end

func copy(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}
