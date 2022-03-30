#!/bin/bash

git clone --recursive https://github.com/supercollider/sc3-plugins.git   /tmp/sc3p
git clone --recursive https://github.com/supercollider/supercollider.git /tmp/sc3

cd /tmp/sc3p
mkdir build && cd build
cmake -DSC_PATH=/tmp/sc3 -DCMAKE_BUILD_TYPE=Release -DSUPERNOVA=ON ..
cmake --build . --config Release
cmake --build . --config Release --target install

rm -rf /tmp/sc3*
