import shutil
import sys
import tempfile
import os.path
import os
import random
import string
import subprocess

from scipy import signal
from scipy.io import wavfile
import numpy as np 

debug_depop = False
if debug_depop:
    import matplotlib.pyplot as plt


def rm(fname):
    try:
        os.remove(fname)
    except:
        pass

def tmpfile() -> str:
    random_string = ''.join(random.choices(string.ascii_uppercase + string.digits, k=10))
    return "sox"+random_string+".wav"

def copypaste(fname,fname2,startPos,endPos,pastePos,crossfade):
    copyLength=endPos-startPos
    piece=tmpfile()
    part1=tmpfile()
    part2=tmpfile()
    splice1=tmpfile()
    leeway=0.001
    cmd=["sox", fname,piece,"trim",str(startPos-crossfade),str(copyLength+2*crossfade)]
    print(cmd)
    os.system(" ".join(cmd))
    cmd=["sox", fname,part1,"trim","0",str(pastePos+crossfade)]
    print(cmd)
    os.system(" ".join(cmd))
    cmd=["sox", fname, part2, "trim", str(pastePos+copyLength-crossfade)]
    print(cmd)
    os.system(" ".join(cmd))
    cmd=["sox", part1, piece, splice1, "splice", "{},{},{}".format(pastePos+crossfade, crossfade, leeway)]
    print(cmd)
    os.system(" ".join(cmd))
    cmd=["sox", splice1, part2, fname2, "splice", "{},{},{}".format(pastePos+copyLength+crossfade, crossfade, leeway)]
    print(cmd)
    os.system(" ".join(cmd))
    rm(piece)
    rm(part1)
    rm(part2)
    rm(splice1)
    return fname2

def findpops(filename):
    # read in file
    samplerate, data = wavfile.read(filename)
    # collect data-type information and normalize
    original_type = data.dtype
    data_max = (np.iinfo(original_type).max)
    data = data.astype(float)
    data = data / data_max

    # compute basic properties
    length = data.shape[0] / samplerate
    num_channels = data.shape[1]
    time = np.linspace(0., length, data.shape[0])

    # stft of the audio
    alldata=data[:,0]
    if num_channels==2:
        alldata=np.add(data[:, 0],data[:, 1])
    f, t, Zxx = signal.stft(alldata, samplerate, window=signal.get_window("hann",256),nperseg=256)
    zmean = np.zeros(Zxx[1,:].shape)
    znum = 0
    for i in range(Zxx.shape[0]):
        z = 20*np.log10(np.abs(Zxx[i,:]))
        z = z - np.mean(z[-500:])
        zmean = np.add(zmean,z)
        znum += 1

    # average all the soft bands
    zmean = zmean / znum 
    zmean = zmean - np.mean(zmean)

    peaks, _ = signal.find_peaks(zmean,prominence=15,threshold=10,height=7,wlen=7)
    if debug_depop:
        print(peaks)
        plt.plot(t,zmean)
        plt.plot(t[peaks],zmean[peaks],'o')
        plt.show()
        # if len(peaks)<30:
        #     break


    pos=[]
    for i,peak in enumerate(peaks):
        pos.append(float(t[peak]))
    return pos

pops=findpops(sys.argv[1])
fname=sys.argv[1]
fname2=sys.argv[2]
if len(pops)==0:
    shutil.copy(fname,fname2)
for i,pop in enumerate(pops):
    print(pop)
    startPos=pop-0.4
    endPos=pop-0.1
    pastePos=pop-0.15
    crossfade=0.05
    copypaste(fname,fname2,startPos,endPos,pastePos,crossfade)
    fname=fname2