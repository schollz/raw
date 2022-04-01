// sclang render.scd
(
var oscScore;
var mainServer;
var nrtServer;
var serverOptions;
var scoreFn;

mainServer = Server(\sampswap_nrt, NetAddr("127.0.0.1", 47112));
serverOptions=ServerOptions.new.numOutputBusChannels_(2);
serverOptions.sampleRate=48000;
nrtServer = Server(\nrt, NetAddr("127.0.0.1", 47114), options:serverOptions);
SynthDef("lpf_rampup", {
    arg out=0,  dur=30, f1,f2,f3,f4;
    var duration=BufDur.ir(0);
    var snd = PlayBuf.ar(2,0,BufRateScale.kr(0));
    snd=LPF.ar(snd,XLine.kr(200,20000,duration));
    snd = snd * EnvGen.ar(Env.new([0, 1, 1, 0], [0.005,dur-0.01,0.005]), doneAction:2);
    Out.ar(out, snd);
}).load(nrtServer);
SynthDef("lpf_rampdown", {
    arg out=0,  dur=30, f1,f2,f3,f4;
    var duration=BufDur.ir(0);
    var snd = PlayBuf.ar(2,0,BufRateScale.kr(0));
    snd=LPF.ar(snd,XLine.kr(20000,200,duration));
    snd = snd * EnvGen.ar(Env.new([0, 1, 1, 0], [0.005,dur-0.01,0.005]), doneAction:2);
    Out.ar(out, snd);
}).load(nrtServer);
SynthDef("dec_ramp", {
    arg out=0,  dur=30, f1,f2,f3,f4;
    var duration=BufDur.ir(0);
    var snd = PlayBuf.ar(2,0,BufRateScale.kr(0));
    snd=SelectX.ar(Line.kr(0,1,duration/4),[snd,Decimator.ar(snd,8000,8)]);
    snd = snd * EnvGen.ar(Env.new([0, 1, 1, 0], [0.005,dur-0.01,0.005]), doneAction:2);
    Out.ar(out, snd);
}).load(nrtServer);
SynthDef("dec", {
    arg out=0,  dur=30, f1,f2,f3,f4;
    var duration=BufDur.ir(0);
    var snd = PlayBuf.ar(2,0,BufRateScale.kr(0));
    snd=Decimator.ar(snd,8000,8);
    snd = snd * EnvGen.ar(Env.new([0, 1, 1, 0], [0.005,dur-0.01,0.005]), doneAction:2);
    Out.ar(out, snd);
}).load(nrtServer);
SynthDef("reverberate", {
    arg out=0,  dur=30, f1,f2,f3,f4;
    var duration=BufDur.ir(0);
    var snd = PlayBuf.ar(2,0,BufRateScale.kr(0));
    snd=SelectX.ar(XLine.kr(0,1,duration/4),[snd,Greyhole.ar(snd* EnvGen.ar(Env.new([0, 1, 1, 0], [0.1,dur-0.2,0.1]), doneAction:2))]);
    snd=LeakDC.ar(snd);
    snd = snd * EnvGen.ar(Env.new([0, 1, 1, 0], [0.1,dur-0.2,0.1]), doneAction:2);
    Out.ar(out, snd);
}).load(nrtServer);
SynthDef("filter_in_out", {
    arg out=0,  dur=30, f1,f2,f3,f4;
    var duration=BufDur.ir(0);
    var snd = PlayBuf.ar(2,0,BufRateScale.kr(0));
    snd = RLPF.ar(snd,
        LinExp.kr(EnvGen.kr(Env.new([0.1, 1, 1, 0.1], [f1,dur-f1-f2,f2])),0.1,1,100,20000),
        0.6);
    snd = snd * EnvGen.ar(Env.new([0, 1, 1, 0], [0.005,dur-0.01,0.005]), doneAction:2);
    Out.ar(out, snd);
}).load(nrtServer);
SynthDef("tapedeck", {
    arg out=0,  dur=30,f1,f2,f3,f4,
    amp=0.9,tape_wet=0.90,tape_bias=0.9,saturation=0.9,drive=0.9,
    tape_oversample=1,mode=0,
    dist_wet=0.05,drivegain=0.5,dist_bias=0.5,lowgain=0.1,highgain=0.1,
    shelvingfreq=600,dist_oversample=1,
    hpf=60,hpfqr=0.6,
    lpf=18000,lpfqr=0.6;
    var duration=BufDur.ir(0);
    var snd = PlayBuf.ar(2,0,BufRateScale.kr(0));
    snd=snd*amp;
    snd=SelectX.ar(Lag.kr(tape_wet,1),[snd,AnalogTape.ar(snd,tape_bias,saturation,drive,tape_oversample,mode)]);    
    snd=SelectX.ar(Lag.kr(dist_wet/10,1),[snd,AnalogVintageDistortion.ar(snd,drivegain,dist_bias,lowgain,highgain,shelvingfreq,dist_oversample)]);              
    snd=RHPF.ar(snd,hpf,hpfqr);
    snd=RLPF.ar(snd,lpf,lpfqr);
    snd = snd * EnvGen.ar(Env.new([0, 1, 1, 0], [0.005,dur-0.01,0.005]), doneAction:2);
    Out.ar(out, snd);
}).load(nrtServer);
SynthDef("bitcrush", {
    arg out=0,  dur=30,f1=1,f2=1,f3,f4;
    var duration=BufDur.ir(0);
    var bcimp=Dust.kr(f1/duration); // times per song
    var timescale=f2*TChoose.kr(bcimp,[0.5,0.75,1,1.25,1.5,1.75]);
    var bctrig=Trig.kr(bcimp,timescale/2);
    var bc=EnvGen.kr(Env.asr(timescale/2,1,timescale/2,[TChoose.kr(bcimp,[-1,-2,-4,-8]),TChoose.kr(bcimp,[2,4,8,16])]),bctrig);
    var snd = PlayBuf.ar(2,0,BufRateScale.kr(0));
    var freq, hasFreq,bitcrushed;
    # freq, hasFreq = Pitch.kr(snd, ampThreshold: 0.02, median: 7);
    bitcrushed=HPF.ar(LPF.ar(MiMu.ar(Decimator.ar(snd,Clip.kr(freq[0],20,2000)*8,8)),2000),100);

    snd=SelectX.ar(bc,[snd,bitcrushed]);
    snd = snd * EnvGen.ar(Env.new([0, 1, 1, 0], [0.005,dur-0.01,0.005]), doneAction:2);
    Out.ar(out, snd);
}).load(nrtServer);
SynthDef("sidechain", {
    arg out=0,  dur=30,f1=1,f2=1,f3,f4;
    var duration=BufDur.ir(0);
    var bcimp=Dust.kr(f1/duration); // times per song
    var timescale=f2*TChoose.kr(bcimp,[0.5,0.75,1,1.25,1.5,1.75]);
    var bctrig=Trig.kr(bcimp,timescale/2);
    var bc=EnvGen.kr(Env.asr(timescale/2,1,timescale/2,[TChoose.kr(bcimp,[-1,-2,-4,-8]),TChoose.kr(bcimp,[2,4,8,16])]),bctrig);
    var snd = PlayBuf.ar(2,0,BufRateScale.kr(0));
    var qntrig=Impulse.kr(f3/60); // f3 = bpm
    var kick=MiPlaits.ar(
        pitch:28,
        harm:SinOsc.kr(Rand(0,1),Rand(0,3)).range(0.1,0.6),
        morph:SinOsc.kr(Rand(0,1),Rand(0,3)).range(0.4,0.6),
        timbre:SinOsc.kr(Rand(0,1),Rand(0,3)).range(0.4,0.6),
        decay:SinOsc.kr(Rand(0,1),Rand(0,3)).range(0.4,0.6),
        engine:13,
        trigger:qntrig, 
    )*EnvGen.ar(Env.perc(0.005,1),qntrig);
    var sidechained=Compander.ar(snd,kick,0.005,1,0.05,0.01,0.1);

    snd=sidechained;//SelectX.ar(bc,[snd,sidechained]);
    snd = snd * EnvGen.ar(Env.new([0, 1, 1, 0], [0.005,dur-0.01,0.005]), doneAction:2);
    Out.ar(out, snd);
}).load(nrtServer);
SynthDef("kick", {
    arg out=0,  dur=30,f1=1,f2=1,f3,f4;
    var duration=BufDur.ir(0);
    var bcimp=Dust.kr(f1/duration); // times per song
    var timescale=f2*TChoose.kr(bcimp,[0.5,0.75,1,1.25,1.5,1.75]);
    var bctrig=Trig.kr(bcimp,timescale/2);
    var bc=EnvGen.kr(Env.asr(timescale/2,1,timescale/2,[TChoose.kr(bcimp,[-1,-2,-4,-8]),TChoose.kr(bcimp,[2,4,8,16])]),bctrig);
    var snd = PlayBuf.ar(2,0,BufRateScale.kr(0));
    var qntrig=Impulse.kr(f3/60); // f3 = bpm
    var kick=MiPlaits.ar(
        pitch:28,
        harm:SinOsc.kr(Rand(0,1),Rand(0,3)).range(0.4,0.6),
        morph:SinOsc.kr(Rand(0,1),Rand(0,3)).range(0.4,0.6),
        timbre:SinOsc.kr(Rand(0,1),Rand(0,3)).range(0.4,0.6),
        decay:SinOsc.kr(Rand(0,1),Rand(0,3)).range(0.4,0.6),
        engine:13,
        trigger:qntrig, 
    )*EnvGen.ar(Env.perc(0.005,1),qntrig);
    var sidechained=Compander.ar(snd,kick,0.005,1,0.05,0.01,0.1);

    snd=sidechained+kick;//SelectX.ar(bc,[snd,sidechained]);
    snd = snd * EnvGen.ar(Env.new([0, 1, 1, 0], [0.005,dur-0.01,0.005]), doneAction:2);
    Out.ar(out, snd);
}).load(nrtServer);

scoreFn={
    arg inFile,outFile,synthDefinition,durationScaling,oscCallbackPort,f1,f2,f3,f4;
    Buffer.read(mainServer,inFile,action:{
        arg buf;
        Routine {
            var buffer;
            var score;
            var duration=buf.duration*durationScaling;

            "defining score".postln;
            score = [
                [0.0, ['/s_new', synthDefinition, 1000, 0, 0, \dur,duration,\f1,f1,\f2,f2,\f3,f3,\f4,f4]],
                [0.0, ['/b_allocRead', 0, inFile]],
                [duration, [\c_set, 0, 0]] // dummy to end
            ];

            "recording score".postln;
            Score(score).recordNRT(
                outputFilePath: outFile,
                sampleRate: 48000,
                headerFormat: "wav",
                sampleFormat: "int24",
                options: nrtServer.options,
                duration: duration,
                action: {
                    Routine {
                        postln("done rendering: " ++ outFile);
                        0.15.wait;
                        NetAddr.new("localhost",oscCallbackPort).sendMsg("/quit");
                    }.play;
                }
            );
        }.play;
    });
};
mainServer.waitForBoot({
    Routine {
        "registring osc for score".postln;
        oscScore = OSCFunc({ arg msg, time, addr, recvPort;
            var inFile=msg[1].asString;
            var outFile=msg[2].asString;
            var synthDefinition=msg[3].asSymbol;
            var durationScaling=msg[4].asFloat;
            var oscCallbackPort=msg[5].asInteger;
            var f1=msg[6].asFloat;
            var f2=msg[7].asFloat;
            var f3=msg[8].asFloat;
            var f4=msg[9].asFloat;
            [msg, time, addr, recvPort].postln;
            scoreFn.value(inFile,outFile,synthDefinition,durationScaling,oscCallbackPort,f1,f2,f3,f4);
            "finished".postln;
        }, '/score',recvPort:47113);
        0.15.wait;
        "writing ready file".postln;
        File.new("/tmp/nrt-scready", "w");
        "ready".postln;
    }.play;
});
)
