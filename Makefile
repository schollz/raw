
install:
	sudo -H python3 -m pip install -r requirements.txt
	git clone https://github.com/schollz/sendosc
	cd sendosc && go install -v
	rm -rf sendosc
