#!/bin/bash

# https://github.com/supercollider/supercollider/blob/develop/README_RASPBERRY_PI.md
# building a GUI-less supercollider

# first run install_jack2.sh!

sudo apt-get install libsamplerate0-dev libsndfile1-dev libasound2-dev libavahi-client-dev libreadline-dev libfftw3-dev libudev-dev libncurses5-dev cmake git


git clone --recurse-submodules https://github.com/supercollider/supercollider.git /tmp/scbuild
cd /tmp/scbuild
mkdir build && cd build

# For a GUI-less build:
cmake -DCMAKE_BUILD_TYPE=Release -DSUPERNOVA=OFF -DSC_ED=OFF -DSC_EL=OFF -DSC_VIM=ON -DNATIVE=ON -DSC_IDE=OFF -DNO_X11=ON -DSC_QT=OFF ..

cmake --build . --config Release --target all -- -j3 
sudo cmake --build . --config Release --target install
sudo ldconfig

echo /usr/local/bin/jackd -P75 -p16 -dalsa -dhw:0 -r44100 -p1024 -n3 > ~/.jackdrc
