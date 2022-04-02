package sampswap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEstiamteTempo(t *testing.T) {
	tempo, err := estimateBPM("this_is_4th_file_bpm_bpm123_45.wav")
	assert.Nil(t, err)
	assert.Equal(t, 123.0, tempo)
}
