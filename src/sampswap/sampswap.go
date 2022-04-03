package sampswap

import (
	"math"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"strings"

	log "github.com/schollz/logger"
	"github.com/schollz/raw/src/sox"
)

type SampSwap struct {
	FileIn        string
	FileOut       string
	FileOriginal  string
	TempoIn       float64
	TempoOut      float64
	TempoEstimate float64
	BeatsIn       float64
	BeatsOut      float64
	ProbStutter   float64
	ProbReverse   float64
	ProbSlow      float64
	ProbJump      float64
	ProbPitch     float64
	ProbReverb    float64
	ProbRereverb  float64
	Tapedeck      bool
	FilterIn      float64
	FilterOut     float64
	RetempoSwitch int // 0-none,1=speed,2=stretch
}

func Init() (ss *SampSwap) {
	return &SampSwap{
		RetempoSwitch: 1,
	}
}

func (ss *SampSwap) Run() (err error) {
	var fname string

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
		ss.TempoIn, err = estimateBPM(fname)
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

	// add repeats until we reach the number of wanted beats
	ss.BeatsIn = math.Floor(math.Round(sox.MustFloat(sox.Length(fname)) / (60 / ss.TempoIn)))
	for {
		if ss.BeatsIn >= ss.BeatsOut {
			break
		}
		fname, err = sox.Repeat(fname, 1)
		if err != nil {
			return
		}
		ss.BeatsIn = math.Floor(math.Round(sox.MustFloat(sox.Length(fname)) / (60 / ss.TempoIn)))
	}

	// trim off excess beats
	log.Tracef("trimming to %fs", ss.BeatsOut*60/ss.TempoIn)
	fname, err = sox.Trim(fname, 0, ss.BeatsOut*60/ss.TempoIn)
	if err != nil {
		return
	}

	for i := 0.0; i < ss.BeatsIn*ss.ProbPitch; i++ {
		fname = ss.pitch(fname)
	}
	for i := 0.0; i < ss.BeatsIn*ss.ProbJump; i++ {
		fname = ss.jump(fname)
	}

	err = os.Rename(fname, ss.FileOut)
	return
}

func (ss *SampSwap) jump(fname string) (fname2 string) {
	var err error
	fname2 = fname
	length_beat := randF(1, 5) / 2
	start_beat := randF(1, ss.BeatsIn-length_beat)
	paste_beat := randF(1, ss.BeatsIn-length_beat)
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
	start_beat := randF(1, ss.BeatsIn-length_beat*8)
	paste_beat := start_beat
	crossfade := 0.095
	var piece string
	piece, err = sox.Trim(ss.FileOriginal,
		60/ss.TempoIn*start_beat-crossfade,
		60/ss.TempoIn*length_beat+crossfade)
	if err != nil {
		return
	}
	piece, err = sox.Pitch(piece, rand.Intn(3)+1)
	if err != nil {
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

func estimateBPM(fname string) (bpm float64, err error) {
	// first see if the bpm appears in the filename
	r, _ := regexp.Compile(`bpm(\d+)`)
	bpm, err = strconv.ParseFloat(strings.TrimPrefix(r.FindString(fname), "bpm"), 64)
	if err == nil {
		return
	}

	// assume the file trimmed and guess it based on the length
	audioLength, err := sox.Length(fname)
	if err != nil {
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

	return
}
