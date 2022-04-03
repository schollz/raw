package sampswap

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEstimateTempo(t *testing.T) {
	tempo, err := estimateBPM("this_is_4th_file_bpm_bpm123_45.wav")
	assert.Nil(t, err)
	assert.Equal(t, 123.0, tempo)
	tempo, err = estimateBPM("175-16-104.wav")
	assert.Nil(t, err)
	assert.Equal(t, 175.0, tempo)
}

func TestRunApp(t *testing.T) {
	ss := Init()
	ss.Seed = 18
	ss.DebugLevel = "debug"
	ss.FileIn = "175-16-104.wav"
	ss.FileOut = "test.wav"
	ss.BeatsOut = 32
	ss.ProbJump = 0.0
	ss.ProbReverse = 0.0
	ss.ProbStutter = 0.1
	assert.Nil(t, ss.Run())
	fmt.Printf("ss: %+v\n", ss)
}
