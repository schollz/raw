# raw

random audio workstation. basically a command-line program to generate entire songs from short sample loops.

for me, making an album of music means devising a single process of making music and using that process to generate tracks. my past processes were focused on melody and harmony. for [*generations*](https://infinitedigits.bandcamp.com/album/generations) I found a process around generative chord changes. for [*18,000*](https://infinitedigits.bandcamp.com/album/18000) I used a coding framework to write small snippets for generative audio. for [*at the place*](https://infinitedigits.bandcamp.com/album/at-the-place) I tried to devise a paper-and-pen method of developing melodies. 

in this current process I'm focusing on effects and structure. I plan on only using samples and trying to elevate the samples by restructuring and juxtaposed effects chains - and have its creation based entirely in a terminal. this current project is based off my other projcets [dnb.lua](https://github.com/schollz/dnb.lua), [makebreakbeat](https://github.com/schollz/makebreakbeat/), and [sampswap](https://github.com/schollz/sampswap/). this one should be better capable of cross-platformability.

## requirements

- python3 numpy+scipy
- sox
- supercollider

you can install python3 anyway you want, and then install the dependencies with:

```
python3 -m pip install numpy scipy
```

`sox` [can be downloaded here](https://sourceforge.net/projects/sox/) (or Linux: `apt get sox`).

SuperCollider is spread across several places (because you'll need to install all the plugins):

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

## usage

```
> git clone https://github.com/schollz/raw
> cd raw
> go build -v
> ./raw stemstitch --config example/config.toml
```

## confg file

see `config-example.toml` and its directory for how to structure samples for a song and add configuration. basically, the config file lets you setup the song structure for each track, and setup effects based on probabilities (0-1).


## license


MIT