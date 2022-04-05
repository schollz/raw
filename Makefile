
test:
	LOGGER=debug go test -v -cover ./...

docker-run:
	docker run --rm -ti raw:latest

docker-build:
	docker build -t raw . 

samples:
	wget https://github.com/schollz/raw/releases/download/samples/samples.tar.gz
	tar -xvzf samples.tar.gz

install:
	sudo -H python3 -m pip install -r requirements.txt
	git clone https://github.com/schollz/sendosc
	cd sendosc && go install -v
	rm -rf sendosc

supercollider:
	sudo apt-get install libjack-jackd2-dev libsamplerate0-dev libsndfile1-dev libasound2-dev libavahi-client-dev libreadline-dev libfftw3-dev libudev-dev libncurses5-dev cmake git
	git clone --recurse-submodules https://github.com/supercollider/supercollider.git /tmp/scbuild
	mkdir /tmp/scbuild/build && cd build
	# For a GUI-less build:
	cd /tmp/scbuild/build && cmake -DCMAKE_BUILD_TYPE=Release -DSUPERNOVA=OFF -DSC_ED=OFF -DSC_EL=OFF -DSC_VIM=ON -DNATIVE=ON -DSC_IDE=OFF -DNO_X11=ON -DSC_QT=OFF ..
	cd /tmp/scbuild/build && cmake --build . --config Release --target all -- -j3
	cd /tmp/scbuild/build && sudo cmake --build . --config Release --target install
	sudo ldconfig

miugens:
	rm -rf /tmp/supercollider-Version-3.12.2 /tmp/mi-UGens
	cd /tmp/ && wget https://github.com/supercollider/supercollider/archive/refs/tags/Version-3.12.2.zip
	cd /tmp/ && unzip Version-3.12.2.zip
	cd /tmp/ && git clone --recurse-submodules https://github.com/v7b1/mi-UGens.git
	cd /tmp/mi-UGens && bash build.sh /tmp/supercollider-Version-3.12.2
	cp -r /tmp/mi-UGens/build/mi-UGens ~/.local/share/SuperCollider/Extensions/

portedplugins:
	rm -rf /tmp/portedplugins
	rm -rf /tmp/supercollider
	git clone --recurse-submodules https://github.com/supercollider/supercollider.git /tmp/supercollider
	git clone --recurse-submodules https://github.com/madskjeldgaard/portedplugins.git /tmp/portedplugins
	mkdir /tmp/portedplugins/build
	cd /tmp/portedplugins/build && cmake .. -DCMAKE_BUILD_TYPE='Release' -DSC_PATH="/tmp/supercollider" -DCMAKE_INSTALL_PREFIX="$(HOME)/.local/share/SuperCollider/Extensions"
	cd /tmp/portedplugins/build && cmake --build . --config Release
	cd /tmp/portedplugins/build && cmake --build . --config Release --target install
	rm -rf /tmp/portedplugins
	rm -rf /tmp/supercollider

sc3plugins:
	rm -rf /tmp/sc3p /tmp/sc3
	git clone --recursive https://github.com/supercollider/sc3-plugins.git   /tmp/sc3p
	git clone --recursive https://github.com/supercollider/supercollider.git /tmp/sc3
	rm -rf /tmp/sc3p /tmp/sc3
	mkdir /tmp/sc3p/build
	cd /tmp/sc3p/build && cmake -DSC_PATH=/tmp/sc3 -DCMAKE_BUILD_TYPE=Release -DSUPERNOVA=OFF ..
	cd /tmp/sc3p/build && cmake --build . --config Release
	cd /tmp/sc3p/build && sudo cmake --build . --config Release --target install
	rm -rf /tmp/sc3*
