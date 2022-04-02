package sampswap

import (
	"regexp"
	"strconv"
	"strings"

	log "github.com/schollz/logger"
	"github.com/schollz/raw/src/sox"
)

type SampSwap struct {
	FileIn        string
	FileOut       string
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

	return
}

func estimateBPM(fname string) (bpm float64, err error) {
	r, _ := regexp.Compile(`bpm(\d+)`)
	bpm, err = strconv.ParseFloat(strings.TrimPrefix(r.FindString(fname), "bpm"), 64)
	if err == nil {
		return
	}
	return
}
