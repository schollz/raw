package song

import (
	"fmt"
	"testing"

	"github.com/pelletier/go-toml"
	"github.com/stretchr/testify/assert"
)

func TestSong(t *testing.T) {
	var s Song
	assert.Nil(t, toml.Unmarshal([]byte(`Tempo = 160.0
Bars = 160.0
Seed = 12

[[track]]
  Name = "chords"
  Structure = "ABABC"

  [[track.part]]
    Name = "A"
    [track.part.ss]
      FileIn = "../../samples/chords/*.wav"
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

  [[track.part]]
    Name = "B"
    [track.part.ss]
      FileIn = "../../samples/chords/*.wav"
      FilterIn = 0.0
      FilterOut = 0.0
      ProbJump = 0.0
      ProbPitch = 0.0
      ProbRereverb = 0.0
      ProbReverb = 0.0
      ProbReverse = 0.0
      ProbSlow = 0.0
      ProbStutter = 0.1
      Sidechain = 0.0
      SilenceAppend = 0.0
      SilencePrepend = 0.0
      Tapedeck = false

  [[track.part]]
    Name = "C"
    [track.part.ss]
      FileIn = "../../samples/chords/*.wav"
      FilterIn = 0.0
      FilterOut = 8.0
      ProbJump = 0.0
      ProbPitch = 0.0
      ProbRereverb = 0.0
      ProbReverb = 0.0
      ProbReverse = 0.0
      ProbSlow = 0.0
      ProbStutter = 0.0
      Sidechain = 0.0
      SilenceAppend = 0.0
      SilencePrepend = 0.0
      Tapedeck = false

[[track]]
  Name = "drums"
  Structure = "CCABCABC"

  [[track.part]]
    Name = "A"
    [track.part.ss]
      FileIn =  "../../samples/drums/*.wav"
      FilterIn = 8.0
      FilterOut = 0.0
      ProbJump = 0.1
      ProbPitch = 0.0
      ProbRereverb = 0.01
      ProbReverb = 0.0
      ProbReverse = 0.01
      ProbSlow = 0.0
      ProbStutter = 0.1
      Sidechain = 0.0
      SilenceAppend = 0.0
      SilencePrepend = 0.0
      Tapedeck = false

  [[track.part]]
    Name = "B"
    [track.part.ss]
      FileIn =  "../../samples/drums/*.wav"
      FilterIn = 8.0
      FilterOut = 0.0
      ProbJump = 0.1
      ProbPitch = 0.0
      ProbRereverb = 0.1
      ProbReverb = 0.0
      ProbReverse = 0.1
      ProbSlow = 0.0
      ProbStutter = 0.1
      Sidechain = 0.0
      SilenceAppend = 0.0
      SilencePrepend = 0.0
      Tapedeck = false

  [[track.part]]
    Name = "C"
    [track.part.ss]
      FileIn =  "../../samples/drums/*.wav"
      FilterIn = 1.0
      FilterOut = 0.0
      ProbJump = 0.1
      ProbPitch = 0.0
      ProbRereverb = 0.1
      ProbReverb = 0.0
      ProbReverse = 0.1
      ProbSlow = 0.0
      ProbStutter = 0.1
      Sidechain = 0.0
      SilenceAppend = 0.0
      SilencePrepend = 0.0
      Tapedeck = false


[[track]]
  Name = "vocals"
  Structure = "BAABCABC"

  [[track.part]]
    Name = "A"
    [track.part.ss]
      FileIn =  "../../samples/vocals/*.wav"
      FilterIn = 8.0
      FilterOut = 0.0
      ProbJump = 0.05
      ProbPitch = 0.0
      ProbRereverb = 0.05
      ProbReverb = 0.0
      ProbReverse = 0.05
      ProbSlow = 0.0
      ProbStutter = 0.05
      Sidechain = 0.0
      SilenceAppend = 0.0
      SilencePrepend = 4.0
      Tapedeck = false

  [[track.part]]
    Name = "B"
    [track.part.ss]
      FileIn =  "../../samples/vocals/*.wav"
      FilterIn = 8.0
      FilterOut = 0.0
      ProbJump = 0.1
      ProbPitch = 0.0
      ProbRereverb = 0.07
      ProbReverb = 0.0
      ProbReverse = 0.05
      ProbSlow = 0.0
      ProbStutter = 0.05
      Sidechain = 0.0
      SilenceAppend = 2.0
      SilencePrepend = 0.0
      Tapedeck = false

  [[track.part]]
    Name = "C"
    [track.part.ss]
      FileIn =  "../../samples/vocals/*.wav"
      FilterIn = 1.0
      FilterOut = 0.0
      ProbJump = 0.1
      ProbPitch = 0.0
      ProbRereverb = 0.05
      ProbReverb = 0.0
      ProbReverse = 0.05
      ProbSlow = 0.0
      ProbStutter = 0.05
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
