package sampswap

import (
	"fmt"
	"io"
	"math"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	log "github.com/schollz/logger"
	"github.com/schollz/progressbar/v3"
	"github.com/schollz/raw/src/sox"
	"github.com/schollz/raw/src/supercollider"
)

type SampSwap struct {
	DebugLevel          string
	Seed                int64
	FileIn              string
	FileOut             string
	FileOriginal        string
	TempoIn             float64
	TempoOut            float64
	BeatsIn             float64
	BeatsOut            float64
	ProbStutter         float64
	ProbReverse         float64
	ProbSlow            float64
	ProbJump            float64
	ProbPitch           float64
	ProbReverb          float64
	ProbRereverb        float64
	Sidechain           float64
	Tapedeck            bool
	FilterIn            float64
	FilterOut           float64
	SilencePrepend      float64 // number of beats
	SilenceAppend       float64 // number of beats
	ReTempoNone         bool    // ignores retempoing
	ReTempoSpeed        bool    // ignores pitch
	showProgress        bool
	doStopSuperCollider bool
}

func Init() (ss *SampSwap) {
	return &SampSwap{}
}

func (ss *SampSwap) Run() (err error) {
	if ss.FileIn == "" {
		err = fmt.Errorf("no input file")
		fmt.Println("HI")
	}
	if err != nil {
		log.Error(err)
		return
	}
	if ss.DebugLevel != "" {
		log.SetLevel(ss.DebugLevel)
	}
	if ss.Seed == 0 {
		ss.Seed = time.Now().UnixNano()
	}
	rand.Seed(ss.Seed)
	var fname string
	if ss.doStopSuperCollider {
		defer supercollider.Stop()
	}

	// convert to 48000
	fname, err = sox.SampleRate(ss.FileIn, 48000)
	if err != nil {
		log.Error(err)
		return
	}

	// trim audio
	fname, err = sox.SilenceTrim(fname)
	if err != nil {
		log.Error(err)
		return
	}

	// estimate bpm
	if ss.TempoIn == 0 {
		ss.TempoIn, err = estimateBPM(fname, ss.FileIn)
		if err != nil {
			log.Error(err)
			return
		}
	}
	if ss.TempoOut == 0 {
		ss.TempoOut = ss.TempoIn
	}

	// determine average number of beats
	ss.BeatsIn = math.Floor(math.Round(sox.MustFloat(sox.Length(fname)) / (60 / ss.TempoIn)))
	if ss.BeatsOut == 0 {
		ss.BeatsOut = ss.BeatsIn
	}

	// add silence and truncate to the average number of beats
	fname, err = sox.SilenceAppend(fname, 2)
	if err != nil {
		log.Error(err)
		return
	}
	fname, err = sox.Trim(fname, 0, ss.BeatsIn*60/ss.TempoIn)
	if err != nil {
		log.Error(err)
		return
	}

	// TODO: make optional
	// find closest multiple of tempoout to tempoin
	foodiff := 1000000.0
	bestBPMDivision := 1.0
	for _, bpmDivision := range []float64{1 / 8, 1 / 4, 1 / 2, 1, 2, 4, 8} {
		if math.Abs(bpmDivision*ss.TempoOut-ss.TempoIn) < foodiff {
			foodiff = math.Abs(bpmDivision*ss.TempoOut - ss.TempoIn)
			bestBPMDivision = bpmDivision
		}
	}
	log.Debugf("tempo in: %f bpm", ss.TempoIn)
	log.Debugf("beats out (requested): %f", ss.BeatsOut)
	ss.TempoOut = ss.TempoOut * bestBPMDivision
	ss.BeatsOut = ss.BeatsOut * bestBPMDivision
	log.Debugf("tempo out: %f bpm", ss.TempoOut)
	log.Debugf("beats out (determined): %f", ss.BeatsOut)

	// add repeats until we reach the number of wanted beats
	// subtract off the beats of silence
	ss.BeatsOut = ss.BeatsOut - ss.SilencePrepend - ss.SilenceAppend
	if ss.BeatsOut < 4 {
		err = fmt.Errorf("too much silence!")
		log.Error(err)
		return
	}
	beats := math.Floor(math.Round(sox.MustFloat(sox.Length(fname)) / (60 / ss.TempoIn)))
	for {
		if beats >= ss.BeatsOut {
			break
		}
		fname, err = sox.Repeat(fname, 1)
		if err != nil {
			log.Error(err)
			return
		}
		beats = math.Floor(math.Round(sox.MustFloat(sox.Length(fname)) / (60 / ss.TempoIn)))
	}

	// trim off excess beats
	log.Tracef("trimming to %fs", ss.BeatsOut*60/ss.TempoIn)
	fname, err = sox.Trim(fname, 0, ss.BeatsOut*60/ss.TempoIn)
	if err != nil {
		log.Error(err)
		return
	}
	log.Debugf("beats: %f", ss.BeatsOut)
	ss.FileOriginal = fname

	total := 2 + ss.BeatsOut*(ss.ProbPitch+ss.ProbJump+
		ss.ProbReverse+ss.ProbStutter+ss.ProbRereverb)
	bar := progressbar.NewOptions(int(total), progressbar.OptionSetVisibility(ss.showProgress))
	for i := 0.0; i < ss.BeatsOut*ss.ProbPitch; i++ {
		fname = ss.pitch(fname)
		bar.Add(1)
	}
	for i := 0.0; i < ss.BeatsOut*ss.ProbJump; i++ {
		fname = ss.jump(fname)
		bar.Add(1)
	}
	for i := 0.0; i < ss.BeatsOut*ss.ProbReverse; i++ {
		fname = ss.reverse(fname)
		bar.Add(1)
	}
	for i := 0.0; i < ss.BeatsOut*ss.ProbRereverb; i++ {
		fname = ss.rereverb(fname)
		bar.Add(1)
	}
	for i := 0.0; i < ss.BeatsOut*ss.ProbStutter; i++ {
		fname = ss.stutter(fname)
		bar.Add(1)
		break
	}

	// add silence to beginning and end if asked for
	if ss.SilencePrepend > 0 {
		fname, err = sox.SilencePrepend(fname, 60/ss.TempoIn*ss.SilencePrepend)
	}
	if ss.SilenceAppend > 0 {
		fname, err = sox.SilenceAppend(fname, 60/ss.TempoIn*ss.SilenceAppend)
	}

	// retempo
	if ss.TempoIn != ss.TempoOut {
		if ss.ReTempoSpeed {
			fname, err = sox.RetempoSpeed(fname, ss.TempoIn, ss.TempoOut)
		} else if ss.ReTempoNone {
		} else {
			fname, err = sox.RetempoStretch(fname, ss.TempoIn, ss.TempoOut)
		}
		if err != nil {
			log.Error(err)
			return
		}
	}

	if ss.Sidechain > 0 {
		fname, err = supercollider.Effect(fname, "sidechain", ss.BeatsOut/ss.Sidechain, 1, ss.TempoIn)
		if err != nil {
			log.Error(err)
			return
		}
	}

	if ss.FilterIn > 0 || ss.FilterOut > 0 {
		fname, err = supercollider.Effect(fname, "filter_in_out",
			ss.FilterIn*60/ss.TempoIn, ss.FilterOut*60/ss.TempoOut)
		if err != nil {
			log.Error(err)
			return
		}
	}
	bar.Add(1)

	if ss.Tapedeck {
		fname, err = supercollider.Effect(fname, "tapedeck")
		if err != nil {
			log.Error(err)
			return
		}
	}
	bar.Add(1)

	if ss.FileOut == "" {
		ss.FileOut = sox.Tmpfile()
	}

	err = os.Rename(fname, ss.FileOut)
	return
}

func (ss *SampSwap) rereverb(fname string) (fname2 string) {
	var err error
	fname2 = fname
	length_beat := randF(1, 4)
	start_beat := randF(3, ss.BeatsOut-length_beat-1)
	paste_beat := randF(2, 2*(ss.BeatsOut-length_beat))
	crossfade := 0.05

	piece, err := sox.Trim(ss.FileOriginal, 60/ss.TempoIn*start_beat-crossfade,
		60/ss.TempoIn/4*length_beat+crossfade*2)
	if err != nil {
		log.Error(err)
		return
	}
	// add reverberate to it
	piece, err = supercollider.Effect(piece, "reverberate")
	if err != nil {
		log.Error(err)
		return
	}
	// reverse it
	piece, err = sox.Reverse(piece)
	if err != nil {
		log.Error(err)
		return
	}
	// paste it
	fname, err = sox.Paste(fname, piece, 60/ss.TempoIn/2*paste_beat, crossfade)
	if err == nil {
		fname2 = fname
	}
	return
}

func (ss *SampSwap) stutter(fname string) (fname2 string) {
	var err error
	fname2 = fname
	start_pos := randF(4, ss.BeatsOut-4) * 60 / ss.TempoIn
	stutters := randF(1, 3) * 4
	stutter_length := 60 / ss.TempoIn / 4
	paste_pos := start_pos - (stutters-1)*stutter_length
	crossfade := 0.01
	// do the stuter
	piece, err := sox.Stutter(ss.FileOriginal, stutter_length,
		start_pos, stutters, crossfade, 0.001)
	if err != nil {
		log.Error(err)
		return
	}
	// add lpf ramp to it
	piece, err = supercollider.Effect(piece, "lpf_rampup")
	if err != nil {
		log.Error(err)
		return
	}
	// paste it
	fname, err = sox.Paste(fname, piece, paste_pos, crossfade)
	if err == nil {
		fname2 = fname
	}
	return
}

func (ss *SampSwap) reverse(fname string) (fname2 string) {
	var err error
	fname2 = fname
	length_beat := randF(1, 5) / 2
	start_beat := randF(1, ss.BeatsOut-length_beat)
	paste_beat := randF(1, ss.BeatsOut-length_beat)
	crossfade := 0.05

	// grab a piece
	piece, err := sox.Trim(fname, 60/ss.TempoIn*start_beat-crossfade,
		60/ss.TempoIn*(length_beat)+2*crossfade)
	if err != nil {
		log.Error(err)
		return
	}

	// reverse it
	piece, err = sox.Reverse(piece)
	if err != nil {
		log.Error(err)
		return
	}

	// paste it
	fname, err = sox.Paste(fname, piece, 60/ss.TempoIn*paste_beat, crossfade)
	if err == nil {
		fname2 = fname
	}
	return
}

func (ss *SampSwap) jump(fname string) (fname2 string) {
	var err error
	fname2 = fname
	length_beat := randF(1, 5) / 2
	start_beat := randF(1, ss.BeatsOut-length_beat)
	paste_beat := randF(1, ss.BeatsOut-length_beat)
	crossfade := 0.05
	log.Tracef("jump - length %f start %f paste %f", length_beat, start_beat, paste_beat)
	fname, err = sox.CopyPaste(fname, 60/ss.TempoIn*start_beat,
		60/ss.TempoIn*(start_beat+length_beat),
		60/ss.TempoIn*paste_beat,
		crossfade)
	if err == nil {
		fname2 = fname
	}
	return
}

func (ss *SampSwap) pitch(fname string) (fname2 string) {
	var err error
	fname2 = fname
	length_beat := randF(1, 4) / 8
	start_beat := randF(1, ss.BeatsOut-length_beat*8)
	paste_beat := start_beat
	crossfade := 0.095
	var piece string
	piece, err = sox.Trim(ss.FileOriginal,
		60/ss.TempoIn*start_beat-crossfade,
		60/ss.TempoIn*length_beat+crossfade)
	if err != nil {
		log.Error(err)
		return
	}
	piece, err = sox.Pitch(piece, rand.Intn(3)+1)
	if err != nil {
		log.Error(err)
		return
	}
	fname, err = sox.Paste(fname, piece, 60/ss.TempoIn*paste_beat, crossfade)
	if err == nil {
		fname2 = fname
	}
	return
}

func randF(min, max float64) float64 {
	return math.Round(min + rand.Float64()*(max-min))
}

func estimateBPM(fname string, originalName ...string) (bpm float64, err error) {
	// first see if the bpm appears in the filename
	r, _ := regexp.Compile(`bpm(\d+)`)
	tryName := fname
	if len(originalName) > 0 {
		tryName = originalName[0]
	}
	bpm, err = strconv.ParseFloat(strings.TrimPrefix(r.FindString(tryName), "bpm"), 64)
	if err == nil {
		log.Debugf("estimated bpm from file name: %f", bpm)
		return
	}

	// assume the file trimmed and guess it based on the length
	audioLength, err := sox.Length(fname)
	if err != nil {
		log.Error(err)
		return
	}

	// find the closest number of even measures
	closestDiff := 1000000.0
	for b := 100.0; b < 200; b++ {
		measures := audioLength / ((60.0 / b) * 4)
		measuresRounded := math.Round(measures)
		if int(measuresRounded)%2 == 0 {
			dif := math.Abs(measuresRounded - measures)
			if int(measuresRounded)%8 == 0 {
				dif = dif / 2
			}
			if dif < closestDiff {
				closestDiff = dif
				bpm = b
			}
		}
	}
	log.Debugf("estimated bpm from file length: %f", bpm)

	return
}

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
