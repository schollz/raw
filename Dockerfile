FROM ubuntu:20.04 as builder
RUN apt-get update
RUN DEBIAN_FRONTEND=noninteractive TZ=Etc/UTC apt-get install -y libjack-jackd2-dev libsamplerate0-dev libsndfile1-dev libasound2-dev libavahi-client-dev libreadline-dev libfftw3-dev libudev-dev libncurses5-dev cmake git unzip wget build-essential g++

RUN mkdir compile
WORKDIR /compile

# build supercollider
RUN git clone --recurse-submodules https://github.com/SuperCollider/SuperCollider.git /compile/sc
WORKDIR /compile/sc
RUN git checkout Version-3.12.2
RUN git submodule update --init --recursive
RUN mkdir build
WORKDIR /compile/sc/build
RUN cmake -DCMAKE_BUILD_TYPE=Release -DSUPERNOVA=OFF -DSC_ED=OFF -DSC_EL=OFF -DSC_VIM=ON -DNATIVE=ON -DSC_IDE=OFF -DNO_X11=ON -DSC_QT=OFF ..
RUN cmake --build . --config Release --target all -- -j3
RUN cmake --build . --config Release --target install
RUN ldconfig
RUN echo /usr/local/bin/jackd -P75 -p16 -dalsa -dhw:0 -r48000 -p1024 -n3 > /root/.jackdrc


# sc3 extensions
WORKDIR /compile
RUN git clone --recursive https://github.com/supercollider/sc3-plugins.git
WORKDIR /compile/sc3-plugins
RUN git checkout Version-3.11.1
RUN mkdir build
WORKDIR /compile/sc3-plugins/build
RUN cmake -DSC_PATH=/compile/sc -DCMAKE_BUILD_TYPE=Release -DSUPERNOVA=OFF ..
RUN cmake --build . --config Release
RUN cmake --build . --config Release --target install


# ported-plugins
RUN git clone --recurse-submodules https://github.com/madskjeldgaard/portedplugins.git /compile/pp
RUN mkdir /compile/pp/build
WORKDIR /compile/pp/build
RUN cmake .. -DCMAKE_BUILD_TYPE='Release' -DSC_PATH=/compile/sc -DCMAKE_INSTALL_PREFIX=/root/.local/share/SuperCollider/Extensions
RUN cmake --build . --config Release
RUN cmake --build . --config Release --target install


# mi-ugens
RUN git clone --recurse-submodules https://github.com/v7b1/mi-UGens.git /compile/miugens
WORKDIR /compile/miugens
RUN bash build.sh /compile/sc
RUN cp -r /compile/miugens/build/mi-UGens /root/.local/share/SuperCollider/Extensions/

# build sendosc
WORKDIR /compile
RUN wget https://go.dev/dl/go1.18.linux-amd64.tar.gz
RUN rm -rf /usr/local/go && tar -C /usr/local -xzf go1.18.linux-amd64.tar.gz
RUN /usr/local/go/bin/go install -v -x github.com/schollz/sendosc@latest

RUN pwd
RUN pwd
# build the slimmed down image
FROM ubuntu:20.04
RUN apt-get update
RUN DEBIAN_FRONTEND=noninteractive TZ=Etc/UTC apt-get install -y libjack-jackd2-dev libsamplerate0-dev libsndfile1-dev libasound2-dev libavahi-client-dev libreadline-dev libfftw3-dev libudev-dev libncurses5-dev lua5.3 python3
COPY --from=builder /usr/local/include/SuperCollider /usr/local/include/SuperCollider
COPY --from=builder /usr/local/lib/SuperCollider /usr/local/lib/SuperCollider
COPY --from=builder /usr/local/share/SuperCollider /usr/local/share/SuperCollider
COPY --from=builder /usr/local/bin/sclang /usr/local/bin/sclang
COPY --from=builder /usr/local/bin/scsynth /usr/local/bin/scsynth
COPY --from=builder /root/go/bin/sendosc /usr/local/bin/sendosc
COPY --from=builder /root/.local /root/.local
COPY scripts/raw.lua /root/raw.lua
COPY scripts/raw.sc /root/raw.sc
WORKDIR /root
CMD ["/bin/bash"]
