package sampswap

import (
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
	DebugLevel     string
	Seed           int64
	FileIn         string
	FileOut        string
	FileOriginal   string
	TempoIn        float64
	TempoOut       float64
	BeatsIn        float64
	BeatsOut       float64
	ProbStutter    float64
	ProbReverse    float64
	ProbSlow       float64
	ProbJump       float64
	ProbPitch      float64
	ProbReverb     float64
	ProbRereverb   float64
	Tapedeck       bool
	FilterIn       float64
	FilterOut      float64
	ReTempoNone    bool
	ReTempoStretch bool
	ReTempoSpeed   bool
}

func Init() (ss *SampSwap) {
	return &SampSwap{
		ReTempoStretch: true,
	}
}

func (ss *SampSwap) Run() (err error) {
	if ss.DebugLevel != "" {
		log.SetLevel(ss.DebugLevel)
	}
	if ss.Seed == 0 {
		ss.Seed = time.Now().UnixNano()
	}
	rand.Seed(ss.Seed)
	var fname string
	defer supercollider.Stop()

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
	// TOOD: make optional
	// find closest multiple of tempoout to tempoin
	foodiff := 1000000.0
	for _, bpm := range []float64{ss.TempoOut / 8, ss.TempoOut / 4, ss.TempoOut / 2, ss.TempoOut, ss.TempoOut * 2, ss.TempoOut * 4, ss.TempoOut * 8} {
		if math.Abs(bpm-ss.TempoIn) < foodiff {
			foodiff = math.Abs(bpm - ss.TempoIn)
			ss.TempoOut = bpm
		}
	}
	log.Infof("tempo in: %f bpm", ss.TempoIn)
	log.Infof("tempo out: %f bpm", ss.TempoOut)

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
	beats := math.Floor(math.Round(sox.MustFloat(sox.Length(fname)) / (60 / ss.TempoIn)))
	for {
		if beats >= ss.BeatsOut {
			break
		}
		fname, err = sox.Repeat(fname, 1)
		if err != nil {
			return
		}
		beats = math.Floor(math.Round(sox.MustFloat(sox.Length(fname)) / (60 / ss.TempoIn)))
	}

	// trim off excess beats
	log.Tracef("trimming to %fs", ss.BeatsOut*60/ss.TempoIn)
	fname, err = sox.Trim(fname, 0, ss.BeatsOut*60/ss.TempoIn)
	if err != nil {
		return
	}
	log.Infof("beats: %f", ss.BeatsOut)
	ss.FileOriginal = fname

	total := 2 + ss.BeatsOut*(ss.ProbPitch+ss.ProbJump+
		ss.ProbReverse+ss.ProbStutter+ss.ProbRereverb)
	bar := progressbar.Default(int64(total))
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
	}

	// retempo
	if ss.ReTempoSpeed {
		fname, err = sox.RetempoSpeed(fname, ss.TempoIn, ss.TempoOut)
	} else if ss.ReTempoStretch {
		fname, err = sox.RetempoStretch(fname, ss.TempoIn, ss.TempoOut)
	}
	if err != nil {
		return
	}

	// fname, err = supercollider.Effect(fname, "kick", ss.BeatsOut/4, 1, ss.TempoIn)
	// if err != nil {
	// 	return
	// }

	if ss.FilterIn > 0 || ss.FilterOut > 0 {
		fname, err = supercollider.Effect(fname, "filter_in_out",
			ss.FilterIn*60/ss.TempoIn, ss.FilterOut*60/ss.TempoOut)
		if err != nil {
			return
		}
	}
	bar.Add(1)

	if ss.Tapedeck {
		fname, err = supercollider.Effect(fname, "tapedeck")
		if err != nil {
			return
		}
	}
	bar.Add(1)

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
		return
	}
	// add reverberate to it
	piece, err = supercollider.Effect(piece, "reverberate")
	if err != nil {
		return
	}
	// reverse it
	piece, err = sox.Reverse(piece)
	if err != nil {
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
	start_beat := randF(1, ss.BeatsOut-4)
	stutters := randF(1, 3) * 4
	paste_beat := randF(12, ss.BeatsOut*4-16)
	crossfade := 0.05
	// do the stuter
	piece, err := sox.Stutter(ss.FileOriginal, 60/ss.TempoIn/4,
		60/ss.TempoIn*start_beat, stutters, crossfade, 0.001)
	if err != nil {
		return
	}
	// add lpf ramp to it
	piece, err = supercollider.Effect(piece, "lpf_rampup")
	if err != nil {
		return
	}
	// paste it
	fname, err = sox.Paste(fname, piece, 60/ss.TempoIn/4*paste_beat, crossfade)
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
		return
	}

	// reverse it
	piece, err = sox.Reverse(piece)
	if err != nil {
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
