package song

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSong(t *testing.T) {
	s := &Song{
		Bars: 160,
	}
	s.Tracks = append(s.Tracks, Track{Name: "chords", Structure: "ABABC"})
	s.Tracks = append(s.Tracks, Track{Name: "vocals", Structure: "DEEBFGH", NameSync: "chords"})
	assert.Nil(t, s.Generate())
}
