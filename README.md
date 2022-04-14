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


## config file


```
Tempo = 175.0
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
    Name = "B"
    [track.part.ss]
      FileIn =  "/home/zns/go/src/github.com/schollz/raw/samples/vocals/*.wav"
      FilterIn = 8.0
      FilterOut = 0.0
      ProbJump = 0.05
      ProbPitch = 0.0
      ProbRereverb = 0.05
      ProbReverb = 0.0
      ProbReverse = 0.05
      ProbSlow = 0.0
      ProbStutter = 0.1
      Sidechain = 0.0
      SilenceAppend = 0.0
      SilencePrepend = 0.0
      Tapedeck = false

  [[track.part]]
    Name = "A"
    [track.part.ss]
      FileIn =  "/home/zns/go/src/github.com/schollz/raw/samples/vocals/*.wav"
      FilterIn = 16.0
      FilterOut = 4.0
      ProbJump = 0.05
      ProbPitch = 0.0
      ProbRereverb = 0.05
      ProbReverb = 0.0
      ProbReverse = 0.05
      ProbSlow = 0.0
      ProbStutter = 0.1
      Sidechain = 0.0
      SilenceAppend = 0.0
      SilencePrepend = 0.0
      Tapedeck = false
```

## todo

- add one-word delay as effect:

```supercollider
(
Routine {
  SynthDef("oneWordDelay",{
    arg in,out;
    var snd,gate,sndDelay;
    snd=In.ar(in,2);
    // one word ping pong delay
    gate=Lag.ar(DetectSilence.ar(snd,amp:0.01,time:60/175/4,doneAction:0),60/175/4).poll;
    sndDelay=CombC.ar(snd,60/175*2,60/175*2,5)*0.25;
    sndDelay=Pan2.ar(sndDelay,SinOsc.kr(175/60));
    snd=SelectX.ar(gate,[snd,sndDelay]);
    Out.ar(out,snd);
  }).add;
  s.sync;

  SynthDef("simplePlayer",{
    arg out,buf;
    var snd;
    snd=PlayBuf.ar(2,buf);
    Out.ar(out,snd);
  }).add;
  s.sync;


  ~delayBus = Bus.audio(s,2);
  s.sync;
  b = Buffer.read(s, "/home/zns/go/src/github.com/schollz/raw/close-to-me/track2.wav");
  s.sync;

  ~player=Synth.new("simplePlayer",[\out,~delayBus,\buf,b]);
  s.sync;
  ~oneWordDelay=Synth.after(~player,"oneWordDelay",[\out,0,\in,~delayBus]);
  s.sync;

}.play;
)




b=Buffer.read(s,"/home/zns/go/src/github.com/schollz/raw/close-to-me/track2.wav");
b.play;
(
{
  var snd=WhiteNoise.ar(0.1)!2;
  snd=BPF.ar(snd,MouseX.kr(100,10000,1));
  snd=Greyhole.ar(snd,2,1,2);
  Out.ar(0,snd);
}.play;
)

~buf = Buffer.alloc(s, s.sampleRate * 5);

(
{
    var trig, freq;
    var dry, snd, phase, iois;
    var cumulDelay;
    var decayFactor, lpfCutoff, numEchoes, ioiScale, ioiDecayFactor;

    // play with these!
    decayFactor = 0.65;
    lpfCutoff = 500;
    numEchoes = 16;
    ioiScale = 0.4;
    ioiDecayFactor = 0.75;

    // example input -- you can substitute other stuff here
    trig = Dust.ar(1);
    dry = Pulse.ar(100 * TIRand.ar(1, 12, trig));
    dry = dry * Decay2.ar(trig, 0.01, 0.2);

    phase = DelTapWr.ar(~buf, dry);
    iois = Array.geom(numEchoes, ioiScale, ioiDecayFactor);
    cumulDelay = iois.sum;
    (cumulDelay > ~buf.duration).if {
        Error("use a larger buffer").throw;
    };
    snd = Silent.ar;
    iois.reverse.do { |ioi|
        snd = (snd + DelTapRd.ar(~buf, phase, cumulDelay, interp: 4)) * decayFactor;
        // one-pole lowpass -- LPF was too aggressive
        snd = OnePole.ar(snd, lpfCutoff / SampleRate.ir);
        cumulDelay = cumulDelay - ioi;
    };
    snd = snd + dry;
    
    snd * 0.5!2;
}.play;

// contributors so far: Nathan Ho
)
```




```
arg bpm=160,times=4; // times = number of times you'd like it to go
var rate = bpm/60/16; // every 16 beats
var imp = Impulse.kr(rate); // impulse every 16 beats
var tries = BufDuration.ir(b)/(16*60/bpm); // tries are the number of 16-beat bars
var probability = times/tries; // tries = probability * times
var length_time=Demand.kr(imp,Dseq([1,2,4,8,16,32]*60/bpm,inf)); // length of each
```


## todo

- track effects
- adding in filter sweeps
- adding in track toggles
- adding in tape slowdowns
- adding in low pass filter for a second or tow
