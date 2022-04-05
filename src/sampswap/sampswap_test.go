package sampswap

import (
	"fmt"
	"testing"

	"github.com/schollz/raw/src/supercollider"
	"github.com/stretchr/testify/assert"
)

func TestEstimateTempo(t *testing.T) {
	tempo, err := estimateBPM("this_is_4th_file_bpm_bpm123_45.wav")
	assert.Nil(t, err)
	assert.Equal(t, 123.0, tempo)
	tempo, err = estimateBPM("175-16-104.wav")
	assert.Nil(t, err)
	assert.Equal(t, 175.0, tempo)
	tempo, err = estimateBPM("../../samples/chords/MiniKit1_160_Chords_Amin_keyAmin_bpm160.wav")
	assert.Nil(t, err)
	assert.Equal(t, 160.0, tempo)
}

func TestRunApp(t *testing.T) {
	go supercollider.Start()
	ss := Init()
	// ss.Seed = 18
	// ss.DebugLevel = ""
	ss.FileIn = "../sox/sample.wav"
	ss.FileIn = "175-16-104.wav"
	ss.FileOut = "test.wav"
	ss.BeatsOut = 16
	ss.ProbJump = 0.0
	ss.ProbReverse = 0.0
	ss.ProbStutter = 0.1
	ss.ProbRereverb = 0.0
	ss.Tapedeck = false
	ss.FilterIn = 0
	ss.FilterOut = 0
	// ss.TempoOut = 160
	ss.SilencePrepend = 0
	ss.SilenceAppend = 0
	ss.doStopSuperCollider = true
	assert.Nil(t, ss.Run())
	fmt.Printf("ss: %+v\n", ss)
}
