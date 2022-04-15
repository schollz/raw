package supercollider

import (
	"fmt"
	"math/rand"
	"os"
	"testing"

	log "github.com/schollz/logger"
	"github.com/stretchr/testify/assert"
)

func TestSuperCollider(t *testing.T) {
	log.SetLevel("trace")
	var fname2 string
	var err error
	fname2 = scPath("../sox/sample.wav")
	fname2, err = Effect(fname2, "tapedeck")
	assert.Nil(t, err)
	fname2, err = Effect(fname2, "lpf_rampup", 2, 0)
	assert.Nil(t, err)
	assert.Nil(t, os.Rename(fname2, "test.wav"))
	assert.Nil(t, Stop())
}

func TestSuperCollider2(t *testing.T) {
	log.SetLevel("trace")
	var fname2 string
	var err error
	fname2 = scPath("../sampswap/175-16-104.wav")
	seed := rand.Float64() * 1000
	fmt.Println(seed)
	fname2, err = Effect(fname2, "toggle", 175, 2, seed)
	assert.Nil(t, err)
	assert.Nil(t, os.Rename(fname2, "test.wav"))
	assert.Nil(t, Stop())
}
