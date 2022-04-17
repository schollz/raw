import shutil
import sys
import tempfile
import os.path
import os
import random
import string
import subprocess

import wave
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


def _wav2array(nchannels, sampwidth, data):
    """data must be the string containing the bytes from the wav file."""
    num_samples, remainder = divmod(len(data), sampwidth * nchannels)
    if remainder > 0:
        raise ValueError(
            "The length of data is not a multiple of " "sampwidth * num_channels."
        )
    if sampwidth > 4:
        raise ValueError("sampwidth must not be greater than 4.")

    if sampwidth == 3:
        a = np.empty((num_samples, nchannels, 4), dtype=np.uint8)
        raw_bytes = np.frombuffer(data, dtype=np.uint8)
        a[:, :, :sampwidth] = raw_bytes.reshape(-1, nchannels, sampwidth)
        a[:, :, sampwidth:] = (a[:, :, sampwidth - 1 : sampwidth] >> 7) * 255
        result = a.view("<i4").reshape(a.shape[:-1])
    else:
        # 8 bit samples are stored as unsigned ints; others as signed ints.
        dt_char = "u" if sampwidth == 1 else "i"
        a = np.frombuffer(data, dtype="<%s%d" % (dt_char, sampwidth))
        result = a.reshape(-1, nchannels)
    return result


def readwav(file):
    """
    Read a wav file.
    Returns the frame rate, sample width (in bytes) and a numpy array
    containing the data.
    This function does not read compressed wav files.
    """
    wav = wave.open(file)
    rate = wav.getframerate()
    nchannels = wav.getnchannels()
    sampwidth = wav.getsampwidth()
    nframes = wav.getnframes()
    data = wav.readframes(nframes)
    wav.close()
    array = _wav2array(nchannels, sampwidth, data)
    return rate, sampwidth, array


def tmpfile() -> str:
    random_string = "".join(
        random.choices(string.ascii_uppercase + string.digits, k=10)
    )
    return os.path.join(tempfile.gettempdir(), "sox" + random_string + ".wav")


def copypaste(fname, fname2, startPos, endPos, pastePos, crossfade):
    copyLength = endPos - startPos
    piece = tmpfile()
    part1 = tmpfile()
    part2 = tmpfile()
    splice1 = tmpfile()
    leeway = 0.001
    cmd = [
        "sox",
        fname,
        piece,
        "trim",
        str(startPos - crossfade),
        str(copyLength + 2 * crossfade),
    ]
    print(cmd)
    os.system(" ".join(cmd))
    cmd = ["sox", fname, part1, "trim", "0", str(pastePos + crossfade)]
    print(cmd)
    os.system(" ".join(cmd))
    cmd = ["sox", fname, part2, "trim", str(pastePos + copyLength - crossfade)]
    print(cmd)
    os.system(" ".join(cmd))
    cmd = [
        "sox",
        part1,
        piece,
        splice1,
        "splice",
        "{},{},{}".format(pastePos + crossfade, crossfade, leeway),
    ]
    print(cmd)
    os.system(" ".join(cmd))
    cmd = [
        "sox",
        splice1,
        part2,
        fname2,
        "splice",
        "{},{},{}".format(pastePos + copyLength + crossfade, crossfade, leeway),
    ]
    print(cmd)
    os.system(" ".join(cmd))
    rm(piece)
    rm(part1)
    rm(part2)
    rm(splice1)
    return fname2


def findpops(filename):
    # read in file
    # samplerate, data = wavfile.read(filename)
    samplerate, _, data = readwav(filename)
    # print(data)
    # collect data-type information and normalize
    original_type = data.dtype
    data_max = np.iinfo(original_type).max
    data = data.astype(float)
    data = data / data_max

    # compute basic properties
    length = data.shape[0] / samplerate
    num_channels = data.shape[1]
    time = np.linspace(0.0, length, data.shape[0])

    # stft of the audio
    alldata = data[:, 0]
    if num_channels == 2:
        alldata = np.add(data[:, 0], data[:, 1])
    f, t, Zxx = signal.stft(
        alldata, samplerate, window=signal.get_window("hann", 256), nperseg=256
    )
    zmean = np.zeros(Zxx[1, :].shape)
    znum = 0
    for i in range(Zxx.shape[0]):
        try:
            z = 20 * np.log10(np.abs(Zxx[i, :]))
            z = z - np.mean(z[-500:])
            zmean = np.add(zmean, z)
            znum += 1
        except:
            pass

    # average all the soft bands
    zmean = zmean / znum
    zmean = zmean - np.mean(zmean)

    peaks, _ = signal.find_peaks(zmean, prominence=15, threshold=10, height=7, wlen=7)
    if debug_depop:
        print(peaks)
        plt.plot(t, zmean)
        plt.plot(t[peaks], zmean[peaks], "o")
        plt.show()
        # if len(peaks)<30:
        #     break

    pos = []
    for i, peak in enumerate(peaks):
        pos.append(float(t[peak]))
    return pos

    # tmp := Tmpfile()
    # defer os.Remove(tmp)
    # // make it 16-bit sined
    # _, _, err = run("sox", fname, "-b", "16", "-e", "signed-integer", tmp)
    # if err != nil {
    #     return
    # }


fname = sys.argv[1]
fname2 = sys.argv[2]
fnametmp = tmpfile()
os.system("sox {} -b 16 -e signed-integer {}".format(fname, fnametmp))
pops = findpops(fnametmp)
if len(pops) == 0:
    shutil.copy(fname, fname2)
for i, pop in enumerate(pops):
    print(pop)
    startPos = pop - 0.4
    endPos = pop - 0.1
    pastePos = pop - 0.15
    crossfade = 0.05
    copypaste(fname, fname2, startPos, endPos, pastePos, crossfade)
    fname = fname2
os.remove(fnametmp)
