package sampswap

import (
	"fmt"
	"testing"

	log "github.com/schollz/logger"
	"github.com/stretchr/testify/assert"
)

func TestEstiamteTempo(t *testing.T) {
	tempo, err := estimateBPM("this_is_4th_file_bpm_bpm123_45.wav")
	assert.Nil(t, err)
	assert.Equal(t, 123.0, tempo)
	tempo, err = estimateBPM("175-16-104.wav")
	assert.Nil(t, err)
	assert.Equal(t, 175.0, tempo)
}

func TestRunApp(t *testing.T) {
	log.SetLevel("trace")
	ss := Init()
	ss.FileIn = "175-16-104.wav"
	ss.FileOut = "test.wav"
	ss.ProbJump = 0.5
	ss.BeatsOut = 16
	assert.Nil(t, ss.Run())
	fmt.Printf("ss: %+v\n", ss)
}
