package song

import (
	"fmt"
	"math"
	"os"
	"runtime"
	"sort"
	"strings"

	"github.com/pelletier/go-toml"
	log "github.com/schollz/logger"
	"github.com/schollz/raw/src/sampswap"
	"github.com/schollz/raw/src/sox"
	"github.com/schollz/raw/src/supercollider"
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
	FileOut        string
}

type Part struct {
	Name     string
	Start    float64
	Length   float64            `toml:"length"` // length in bars (each bar is 4 beats)
	SampSwap *sampswap.SampSwap `toml:"ss"`     // original sample for this part
}

func (s *Song) Generate() (err error) {
	// if there exist sampswap values, save them
	sampswapParts := make(map[int]map[string]*sampswap.SampSwap)
	for i, track := range s.Tracks {
		for _, part := range track.Parts {
			if part.SampSwap != nil && part.Name != "" {
				if _, ok := sampswapParts[i]; !ok {
					sampswapParts[i] = make(map[string]*sampswap.SampSwap)
				}
				sampswapParts[i][part.Name] = part.SampSwap
			}
		}
	}

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
		// add the known sampswap to each part
		for j, part := range s.Tracks[i].Parts {
			if _, ok := sampswapParts[i]; !ok {
				continue
			}
			if _, ok := sampswapParts[i][part.Name]; !ok {
				continue
			}
			s.Tracks[i].Parts[j].SampSwap = sampswapParts[i][part.Name]
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
	b, _ := toml.Marshal(s)
	fmt.Println(string(b))

	// run all the song components
	if err = s.RunAll(); err != nil {
		return
	}

	// combine all the song components
	if err = s.CombineAll(); err != nil {
		return
	}

	// rename the final file for each track
	for _, track := range s.Tracks {
		if track.FileOut != "" {
			newName := "temp.wav" // TODO change this
			log.Debugf("%s -> %s", track.FileOut, newName)
			os.Rename(track.FileOut, newName)
		}
	}

	// clean up everything
	sox.Clean()
	supercollider.Stop()
	return
}

func (s *Song) CombineAll() (err error) {
	// start worker group to generate the parts for each track
	numJobs := len(s.Tracks)
	type job struct {
		tracki int
	}
	type result struct {
		tracki int
		err    error
	}
	jobs := make(chan job, numJobs)
	results := make(chan result, numJobs)
	runtime.GOMAXPROCS(runtime.NumCPU())
	for i := 0; i < runtime.NumCPU(); i++ {
		go func(jobs <-chan job, results chan<- result) {
			for j := range jobs {
				// step 3: specify the work for the worker
				var r result
				r.tracki = j.tracki
				fileList := []string{}
				for _, part := range s.Tracks[j.tracki].Parts {
					if part.SampSwap != nil {
						if part.SampSwap.FileOut != "" {
							fileList = append(fileList, part.SampSwap.FileOut)
						}
					}
				}
				if len(fileList) > 1 {
					s.Tracks[j.tracki].FileOut, r.err = sox.Join(fileList...)
				} else if len(fileList) == 1 {
					s.Tracks[j.tracki].FileOut = fileList[0]
				}
				results <- r
			}
		}(jobs, results)
	}
	for tracki := range s.Tracks {
		jobs <- job{tracki: tracki}
	}
	close(jobs)
	for i := 0; i < numJobs; i++ {
		r := <-results
		if r.err != nil {
			// do something with error
			log.Errorf("%+v: %s", r, r.err)
			err = r.err
		}
	}
	return
}

func (s *Song) RunAll() (err error) {
	// start worker group to generate the parts for each track
	numJobs := 0
	for _, track := range s.Tracks {
		for _, part := range track.Parts {
			if part.SampSwap != nil {
				numJobs++
			}
		}
	}
	type job struct {
		tracki int
		parti  int
	}
	type result struct {
		tracki int
		parti  int
		err    error
	}
	jobs := make(chan job, numJobs)
	results := make(chan result, numJobs)
	runtime.GOMAXPROCS(runtime.NumCPU())
	for i := 0; i < runtime.NumCPU(); i++ {
		go func(jobs <-chan job, results chan<- result) {
			for j := range jobs {
				// step 3: specify the work for the worker
				var r result
				r.tracki = j.tracki
				r.parti = j.parti
				s.Tracks[j.tracki].Parts[j.parti].SampSwap.BeatsOut = s.Tracks[j.tracki].Parts[j.parti].Length * 4
				// TODO: add silence to tracks if they aren't the first track?
				r.err = s.Tracks[j.tracki].Parts[j.parti].SampSwap.Run()
				results <- r
			}
		}(jobs, results)
	}
	for tracki, track := range s.Tracks {
		for parti, part := range track.Parts {
			if part.SampSwap != nil {
				jobs <- job{tracki: tracki, parti: parti}
			}
		}
	}
	close(jobs)
	for i := 0; i < numJobs; i++ {
		r := <-results
		if r.err != nil {
			// do something with error
			log.Errorf("%+v: %s", r, r.err)
			err = r.err
		}
	}
	return
}
