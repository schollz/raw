# raw

random audio workstation.

like a [daw](https://en.wikipedia.org/wiki/Digital_audio_workstation) but most everything is stochastic drawn from a probability distribution you can specify.

## requirements

requires [`sox`](https://sourceforge.net/projects/sox/) which you can download or (on Linux) install with `apt get sox`.

you'll also need SuperCollider+plugins:

- [SuperCollider](https://supercollider.github.io/downloads)
- [SuperCollider sc3-plugins](http://supercollider.github.io/sc3-plugins/)
- [SuperCollider mi-ugens](https://github.com/v7b1/mi-UGens/releases)
- [SuperCollider ported plugins](https://github.com/madskjeldgaard/portedplugins/releases)

you can pretty easily install them by source (i.e. for Raspberry Pi where binaries aren't available):

```
git clone https://github.com/schollz/raw && cd raw
sudo apt install cmake
make supercollider
make sc3plugins
make miugens
make portedplugins
```