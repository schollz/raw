
build:
	docker build -t raw . 

run:
	docker run --rm -ti raw:latest

install:
	sudo -H python3 -m pip install -r requirements.txt
	git clone https://github.com/schollz/sendosc
	cd sendosc && go install -v
	rm -rf sendosc

miugens:
	rm -rf /tmp/supercollider-Version-3.12.2 /tmp/mi-UGens
	cd /tmp/ && wget https://github.com/supercollider/supercollider/archive/refs/tags/Version-3.12.2.zip
	cd /tmp/ && unzip Version-3.12.2.zip
	cd /tmp/ && git clone --recurse-submodules https://github.com/v7b1/mi-UGens.git
	cd /tmp/mi-UGens && bash build.sh /tmp/supercollider-Version-3.12.2
	cp -r /tmp/mi-UGens/build/mi-UGens ~/.local/share/SuperCollider/Extensions/

portedplugins:

