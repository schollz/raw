# raw

[![Go Reference](https://pkg.go.dev/badge/github.com/schollz/raw.svg)](https://pkg.go.dev/github.com/schollz/raw)

random audio workstation. wip. projected finish: june?

## idea

`<ramble>`

to me, making an album of music means devising a single process of making music and using that process to generate tracks.

my past processes were focused on melody and harmony. for my album [generations](https://infinitedigits.bandcamp.com/album/generations) I made a process around generative chord changes. for my album [18,000](https://infinitedigits.bandcamp.com/album/18000) I used a coding framework to write small snippets for generative audio. for my album [at the place](https://infinitedigits.bandcamp.com/album/at-the-place) I tried to devise a paper-and-pen method of developing melodies.

in this current process I'm focusing on effects and structure. I plan on only using samples and trying to elevate the samples by restructuring and juxtaposed effects chains. 

this current project is based off my other projcets [dnb.lua](https://github.com/schollz/dnb.lua), [makebreakbeat](https://github.com/schollz/makebreakbeat/), and [sampswap](https://github.com/schollz/sampswap/). this one should be better capable of cross-platformability.

`</ramble>`

## requirements

requires `sox` which [you can download](https://sourceforge.net/projects/sox/) (or Linux: `apt get sox`).

you'll also need SuperCollider+plugins:

- [SuperCollider](https://supercollider.github.io/downloads)
- [SuperCollider sc3-plugins](http://supercollider.github.io/sc3-plugins/)
- [SuperCollider mi-ugens](https://github.com/v7b1/mi-UGens/releases)
- [SuperCollider ported plugins](https://github.com/madskjeldgaard/portedplugins/releases)

you can click those links to download the files for your system. its also possible to install them by source (i.e. for Raspberry Pi where binaries aren't available):

```
git clone https://github.com/schollz/raw && cd raw
sudo apt install cmake
make supercollider
make sc3plugins
make miugens
make portedplugins
```


### sampswap units

`chords.toml`:

```toml
# map[string]*sampswap.SampSwap

[A]
FileIn = ""
FilterIn = 0.0
FilterOut = 0.0
ProbJump = 0.0
ProbPitch = 0.0
ProbRereverb = 0.0
ProbReverb = 0.0
ProbReverse = 0.0
ProbSlow = 0.0
ProbStutter = 0.1
ReTempoSpeed = false
Sidechain = 0.0
Tapedeck = false

[B]
FileIn = ""
FilterIn = 0.0
FilterOut = 0.0
ProbJump = 0.0
ProbPitch = 0.0
ProbRereverb = 0.0
ProbReverb = 0.0
ProbReverse = 0.0
ProbSlow = 0.0
ProbStutter = 0.1
ReTempoSpeed = false
Sidechain = 0.0
Tapedeck = false
```