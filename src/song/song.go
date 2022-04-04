package song

import (
	"fmt"
	"math"
	"sort"
	"strings"

	"github.com/pelletier/go-toml"
	log "github.com/schollz/logger"
	"github.com/schollz/raw/src/sampswap"
)

type Song struct {
	Tempo  float64
	Bars   float64 // each bar is 4 beats
	Seed   int64
	Tracks []Track `toml:"track"`
}

type Track struct {
	Name           string
	Structure      string   // "ABABC"
	StructureArray []string // []string{"A","B"}
	NameSync       string
	Parts          []Part `toml:"part"`
}

type Part struct {
	Name     string
	Start    float64
	Length   float64            `toml:"length"` // length in bars (each bar is 4 beats)
	SampSwap *sampswap.SampSwap `toml:"ss"`     // original sample for this part
}

func (s *Song) Generate() (err error) {
	// first run through, generate the part positions for each track
	namesToIndex := make(map[string]int)
	for i, track := range s.Tracks {
		namesToIndex[track.Name] = i
	}
	for i, track := range s.Tracks {
		s.Tracks[i].StructureArray = strings.Split(track.Structure, "")
		log.Tracef("track%d: %v", i, s.Tracks[i].StructureArray)
		s.Tracks[i].Parts = []Part{}
		for j, name := range s.Tracks[i].StructureArray {
			p := Part{Name: name, Start: math.Round(s.Bars * float64(j) / float64(len(s.Tracks[i].StructureArray)))}
			log.Debugf("part: %v", p)
			s.Tracks[i].Parts = append(s.Tracks[i].Parts, p)
		}
		// find the tracks in the synced track
		if _, ok := namesToIndex[track.NameSync]; ok {
			for _, partOther := range s.Tracks[namesToIndex[track.NameSync]].Parts {
				if !strings.Contains(track.Structure, partOther.Name) {
					continue
				}
				s.Tracks[i].Parts = append(s.Tracks[i].Parts, partOther)
			}
		}
		// sort the parts
		sort.Slice(s.Tracks[i].Parts, func(m, n int) bool {
			return s.Tracks[i].Parts[m].Start < s.Tracks[i].Parts[n].Start
		})
		// determine the lengths of each part
		for j, part := range s.Tracks[i].Parts {
			nextStart := s.Bars
			if j < len(s.Tracks[i].Parts)-1 {
				nextStart = s.Tracks[i].Parts[j+1].Start
			}
			s.Tracks[i].Parts[j].Length = nextStart - part.Start
		}
	}
	// b, _ := json.MarshalIndent(s, " ", " ")
	s.Tracks[0].Parts[0].SampSwap = &sampswap.SampSwap{}
	s.Tracks[1].Parts[0].SampSwap = &sampswap.SampSwap{}
	b, _ := toml.Marshal(s)
	fmt.Println(string(b))
	return
}
