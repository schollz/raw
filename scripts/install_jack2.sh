sudo apt-get install libsamplerate0-dev libsndfile1-dev libasound2-dev libavahi-client-dev libreadline-dev libfftw3-dev libudev-dev libncurses5-dev cmake git

cd ~ # or cd into the directory where you'd like to build jack2
git clone git://github.com/jackaudio/jack2 --depth 1 /tmp/jack2
cd /tmp/jack2
./waf configure --alsa --libdir=/usr/lib/arm-linux-gnueabihf/
./waf build
sudo ./waf install
sudo ldconfig
cd ..
rm -rf /tmp/jack2
sudo sh -c "echo @audio - memlock 256000 >> /etc/security/limits.conf"
sudo sh -c "echo @audio - rtprio 75 >> /etc/security/limits.conf"
