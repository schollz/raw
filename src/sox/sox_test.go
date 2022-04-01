package sox

import (
	"math"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {
	stdout, stderr, err := run("sox", "--help")
	assert.Nil(t, err)
	assert.True(t, strings.Contains(stdout, "SoX"))
	assert.Empty(t, stderr)
}

func TestLength(t *testing.T) {
	length, err := Length("sample.wav")
	assert.Nil(t, err)
	assert.Equal(t, 1.133354, length)
}

func TestSampleRate(t *testing.T) {
	samplerate, channnels, err := SampleRate("sample.wav")
	assert.Nil(t, err)
	assert.Equal(t, 48000, samplerate)
	assert.Equal(t, 2, channnels)
}

func TestSilence(t *testing.T) {
	fname2, err := SilenceAppend("sample.wav", 1)
	assert.Nil(t, err)
	length1, _ := Length("sample.wav")
	length2, _ := Length(fname2)
	assert.Less(t, math.Abs(length2-length1-1), 0.00001)

	fname2, err = SilencePrepend("sample.wav", 1)
	assert.Nil(t, err)
	length1, _ = Length("sample.wav")
	length2, _ = Length(fname2)
	assert.Less(t, math.Abs(length2-length1-1), 0.00001)

	fname3 := MustString(SilenceTrim(fname2))
	length3 := MustFloat(Length(fname3))
	assert.Greater(t, length2-length3, 1.0)

	os.Rename(fname3, "test.wav")
}

func TestClean(t *testing.T) {
	assert.Nil(t, Clean())
}
