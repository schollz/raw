package sox

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	log "github.com/schollz/logger"
)

// TempDir is where the temporary intermediate files are held
var TempDir = os.TempDir()

// TempPrefix is a unique indicator of the temporary files
var TempPrefix = "sox"

// TempType is the type of file to be generated (should be "wav")
var TempType = "wav"

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
		log.Errorf("%s: '%s'", strings.Join(args, " "), err.Error())
	}
	return outb.String(), errb.String(), err
}

func tmpfile() string {
	randBytes := make([]byte, 16)
	rand.Read(randBytes)
	return filepath.Join(TempDir, TempPrefix+hex.EncodeToString(randBytes)+"."+TempType)
}

// Clean will remove files created after each function
func Clean() (err error) {
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

// SampleRate returns the sample rate and number of channels for file
func SampleRate(fname string) (samplerate int, channels int, err error) {
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
	samplerate, channels, err := SampleRate(fname)
	if err != nil {
		return
	}
	silencefile := tmpfile()
	defer os.Remove(silencefile)
	fname2 = tmpfile()
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

	return
}

// SilencePrepend prepends silence to a file
func SilencePrepend(fname string, length float64) (fname2 string, err error) {
	samplerate, channels, err := SampleRate(fname)
	if err != nil {
		return
	}
	silencefile := tmpfile()
	defer os.Remove(silencefile)
	fname2 = tmpfile()
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

// SilenceTrim trims silence around a file
func SilenceTrim(fname string) (fname2 string, err error) {
	fname2 = tmpfile()
	_, _, err = run("sox", fname, fname2, "silence", "1", "0.1", `0.025%`, "reverse", "silence", "1", "0.1", `0.25%`, "reverse")
	if err != nil {
		return
	}
	return
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

// TODO

// function audio.trim(fname,start,length)
//   local fname2=string.random_filename()
//   if length==nil then
//     os.cmd("sox "..fname.." "..fname2.." trim "..start)
//   else
//     os.cmd("sox "..fname.." "..fname2.." trim "..start.." "..length)
//   end
//   return fname2
// end

// function audio.reverse(fname)
//   local fname2=string.random_filename()
//   os.cmd(string.format("sox %s %s reverse",fname,fname2))
//   return fname2
// end

// function audio.pitch(fname,notes)
//   local fname2=string.random_filename()
//   os.cmd(string.format("sox %s %s pitch %d",fname,fname2,notes*100))
//   return fname2
// end

// function audio.join(fnames)
//   local fname2=string.random_filename()
//   os.cmd(string.format("sox %s %s",table.concat(fnames," "),fname2))
//   return fname2
// end

// function audio.repeat_n(fname,repeats)
//   local fname2=string.random_filename()
//   os.cmd(string.format("sox %s %s repeat %d",fname,fname2,repeats))
//   return fname2
// end

// function audio.retempo_speed(fname,old_tempo,new_tempo)
//   local fname2=string.random_filename()
//   os.cmd(string.format("sox %s %s speed %f rate -v 48k",fname,fname2,new_tempo/old_tempo))
//   return fname2
// end

// function audio.retempo_stretch(fname,old_tempo,new_tempo)
//   local fname2=string.random_filename()
//   os.cmd(string.format("sox %s %s tempo -m %f",fname,fname2,new_tempo/old_tempo))
//   return fname2
// end

// function audio.slowdown(fname,slowdown)
//   local fname2=string.random_filename()
//   os.cmd(string.format("sox %s %s tempo -m %f",fname,fname2,slowdown))
//   return fname2
// end

// function audio.get_info(fname)
//   local sample_rate=tonumber(os.capture("sox --i "..fname.." | grep 'Sample Rate' | awk '{print $4}'"))
//   local channels=tonumber(os.capture("sox --i "..fname.." | grep 'Channels' | awk '{print $3}'"))
//   return sample_rate,channels
// end

// -- copy_and_paste2 finds best positionn, but does not keep timing
// function audio.copy_and_paste2(fname,copy_start,copy_stop,paste_start)
// 	local copy_length=copy_stop-copy_start
//   if copy_length==nil or copy_length<0.05 then
//     do return fname end
//   end
// 	local piece=string.random_filename()
// 	local part1=string.random_filename()
// 	local part2=string.random_filename()
// 	local fname2=string.random_filename()
// 	local e=5/1000
// 	local l=5/1000
// 	os.cmd(string.format("sox %s %s trim %f %f",fname,piece,copy_start-e-l,copy_stop-copy_start+e+l+e))
// 	os.cmd(string.format("sox %s %s trim 0 %f",fname,part1,paste_start+e))
// 	os.cmd(string.format("sox %s %s trim %f",fname,part2,paste_start+copy_stop-copy_start-e-l))
// 	os.cmd(string.format("sox %s %s %s %s splice %f %f",part1,piece,part2,fname2,paste_start+e,paste_start+e+copy_stop-copy_start+e+l+e))
//   os.cmd(string.format("rm -f %s %s %s",part1,part2,piece))
// 	return fname2
// end

// function audio.copy_and_paste(fname,copy_start,copy_stop,paste_start,crossfade)
// 	local copy_length=copy_stop-copy_start
//   if copy_length==nil or copy_length<0.05 then
//     do return fname end
//   end
// 	local piece=string.random_filename()
// 	local part1=string.random_filename()
// 	local part2=string.random_filename()
// 	local fname2=string.random_filename()
// 	local splice1=string.random_filename()
// 	local e=crossfade or 0.1
// 	local l=0 -- no leeway
// 	os.cmd(string.format("sox %s %s trim %f %f",fname,piece,copy_start-e,copy_length+2*e))
// 	os.cmd(string.format("sox %s %s trim 0 %f",fname,part1,paste_start+e))
// 	os.cmd(string.format("sox %s %s trim %f",fname,part2,paste_start+copy_length-e))
// 	os.cmd(string.format("sox %s %s %s splice %f,%f,%f",part1,piece,splice1,paste_start+e,e,l))
// 	os.cmd(string.format("sox %s %s %s splice %f,%f,%f",splice1,part2,fname2,paste_start+copy_length+e,e,l))
//   os.cmd(string.format("rm -f %s %s %s %s",piece,part1,part2,splice1))
// 	return fname2
// end

// -- pastes any piece into a place in the audio
// -- assumes that the piece has "crossfade" length on both sides
// -- in addition to its current length
// function audio.paste(fname,piece,paste_start,crossfade)
// 	local copy_length=audio.length(piece)
//   if copy_length==nil then
//     do return fname end
//   end
// 	local part1=string.random_filename()
// 	local part2=string.random_filename()
// 	local fname2=string.random_filename()
// 	local splice1=string.random_filename()
// 	local e=crossfade or 0.1
// 	local l=0 -- no leeway
// 	os.cmd(string.format("sox %s %s trim 0 %f",fname,part1,paste_start+e))
// 	os.cmd(string.format("sox %s %s trim %f",fname,part2,paste_start+copy_length-e*3))
// 	os.cmd(string.format("sox %s %s %s splice %f,%f,%f",part1,piece,splice1,paste_start+e,e,l))
// 	os.cmd(string.format("sox %s %s %s splice %f,%f,%f",splice1,part2,fname2,paste_start+copy_length+e,e,l))
//   os.cmd(string.format("rm -f %s %s %s",part1,part2,splice1))
// 	return fname2
// end

// function audio.sample_rate(fname,sr,ch)
//   sr=sr or 48000
//   ch=ch or 2
//   local fname2=string.random_filename()
//   os.cmd(string.format("sox -r %d %s %s",sr,fname,fname2))
//   return fname2
// end

// function audio.gain(fname,gain)
//   local fname2=string.random_filename()
//   os.cmd(string.format("sox %s %s gain %f",fname,fname2,gain))
//   return fname2
// end

// function audio.stretch(fname,stretch)
//   local fname2=string.random_filename()
//   os.cmd(string.format("sox %s %s stretch %f",fname,fname2,stretch))
//   return fname2
// end

// function audio.stutter(fname,stutter_length,pos_start,count,crossfade_piece,crossfade_stutter,gain_amt)
// 	crossfade_piece=0.1 or crossfade_piece
// 	crossfade_stutter=0.005 or crossfade_stutter
// 	local partFirst=string.random_filename()
// 	local partMiddle=string.random_filename()
// 	local partLast=string.random_filename()
// 	os.cmd(string.format("sox %s %s trim %f %f",fname,partFirst,pos_start-crossfade_piece,stutter_length+crossfade_piece+crossfade_stutter))
// 	os.cmd(string.format("sox %s %s trim %f %f",fname,partMiddle,pos_start-crossfade_stutter,stutter_length+crossfade_stutter+crossfade_stutter))
// 	os.cmd(string.format("sox %s %s trim %f %f",fname,partLast,pos_start-crossfade_stutter,stutter_length+crossfade_piece+crossfade_stutter))
//   gain_amt=gain_amt or (count>8 and -1.5 or -2)
// 	for i=1,count do
// 		local fnameNext=""
// 		if i==1 then
// 			fnameNext=audio.gain(partFirst,gain_amt*(count-i))
// 		else
// 			fnameNext=string.random_filename()
//       local fnameMid=i<count and partMiddle or partLast
//       if gain_amt~=0 then
//         fnameMid=audio.gain(fnameMid,gain_amt*(count-i))
//       end
// 			os.cmd(string.format("sox %s %s %s splice %f,%f,0",fname2,fnameMid,fnameNext,audio.length(fname2),crossfade_stutter))
// 		end
// 		fname2=fnameNext
// 	end
// 	return fname2
// end
