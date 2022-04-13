package song

import (
	"fmt"
	"testing"

	"github.com/pelletier/go-toml"
	"github.com/stretchr/testify/assert"
)

func TestSong(t *testing.T) {
	var s Song
	assert.Nil(t, toml.Unmarshal([]byte(`Tempo = 175.0
Bars = 96.0
Seed = 0

[[track]]
  Name = "chords"
  Structure = "AABCDE"

  [[track.part]]
    Name = "A"
    [track.part.ss]
      FileIn = "/home/zns/go/src/github.com/schollz/raw/samples/chords/*.wav"
      FilterIn = 8.0
      FilterOut = 0.0
      ProbJump = 0.0
      ProbPitch = 0.0
      ProbRereverb = 0.0
      ProbReverb = 0.0
      ProbReverse = 0.0
      ProbSlow = 0.0
      ProbStutter = 0.2
      Sidechain = 0.0
      SilenceAppend = 0.0
      SilencePrepend = 0.0
      Tapedeck = false

[[track]]
  Name = "drums"
  Structure = "AABBCCDGEHFF"

  [[track.part]]
    Name = "A"
    [track.part.ss]
      FileIn =  "/home/zns/go/src/github.com/schollz/raw/samples/drums/*.wav"
      FilterIn = 8.0
      FilterOut = 0.0
      ProbJump = 0.1
      ProbPitch = 0.0
      ProbRereverb = 0.01
      ProbReverb = 0.0
      ProbReverse = 0.01
      ProbSlow = 0.0
      ProbStutter = 0.15
      Sidechain = 0.0
      SilenceAppend = 0.0
      SilencePrepend = 0.0
      Tapedeck = true

[[track]]
  Name = "vocals"
  Structure = "ABCDEF"

  [[track.part]]
    Name = "A"
    [track.part.ss]
      FileIn =  "/home/zns/go/src/github.com/schollz/raw/samples/vocals/*.wav"
      FilterIn = 8.0
      FilterOut = 0.0
      ProbJump = 0.05
      ProbPitch = 0.0
      ProbRereverb = 0.05
      ProbReverb = 0.0
      ProbReverse = 0.07
      ProbSlow = 0.0
      ProbStutter = 0.1
      Sidechain = 0.0
      SilenceAppend = 0.0
      SilencePrepend = 0.0
      Tapedeck = false



`), &s))
	fmt.Printf("%+v\n", s.Tracks[0].Parts[0].SampSwap)
	// s.Tracks = append(s.Tracks, Track{Name: "chords", Structure: "ABABC"})
	// s.Tracks = append(s.Tracks, Track{Name: "vocals", Structure: "DEEBFGH", NameSync: "chords"})
	assert.Nil(t, s.Generate())
}
