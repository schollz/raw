package supercollider

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSuperCollider(t *testing.T) {
	assert.Nil(t, Start())
	fname2, err := Effect(scPath("../sox/sample.wav"), "tapedeck")
	assert.Nil(t, err)
	fmt.Println(fname2)
	assert.Nil(t, Stop())
	os.Rename(fname2, "test.wav")
}
