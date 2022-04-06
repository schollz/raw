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
Seed = 0

[[track]]
  Name = "chords"
  Structure = "ABABC"

  [[track.part]]
    Name = "A"
    [track.part.ss]
      FileIn = "../../samples/chords/MiniKit1_160_Chords_Amin_keyAmin_bpm160.wav"
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
      FileIn = "../../samples/chords/KIT6_160_Chords_Future_Modulation_6_Amin_keyAmin_bpm160.wav"
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
      FileIn = "../../samples/chords/OS_NC_160_Am_Showdown_Synths_keyAmin_bpm160.wav"
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
  Structure = "DEFGCH"

  [[track.part]]
    Name = "D"
    [track.part.ss]
      FileIn =  "../../samples/drums/*.wav"
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
      FileIn = "../../samples/chords/KIT6_160_Chords_Future_Modulation_6_Amin_keyAmin_bpm160.wav"
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
      FileIn = "../../samples/chords/OS_NC_160_Am_Showdown_Synths_keyAmin_bpm160.wav"
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
`), &s))
	fmt.Printf("%+v\n", s.Tracks[0].Parts[0].SampSwap)
	// s.Tracks = append(s.Tracks, Track{Name: "chords", Structure: "ABABC"})
	// s.Tracks = append(s.Tracks, Track{Name: "vocals", Structure: "DEEBFGH", NameSync: "chords"})
	assert.Nil(t, s.Generate())
}
