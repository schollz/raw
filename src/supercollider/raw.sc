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
    Out.ar(out, LeakDC.ar(snd));
}).load(nrtServer);
SynthDef("lpf_rampdown", {
    arg out=0,  dur=30, f1,f2,f3,f4;
    var duration=BufDur.ir(0);
    var snd = PlayBuf.ar(2,0,BufRateScale.kr(0));
    snd=LPF.ar(snd,XLine.kr(20000,200,duration));
    snd = snd * EnvGen.ar(Env.new([0, 1, 1, 0], [0.005,dur-0.01,0.005]), doneAction:2);
    Out.ar(out, LeakDC.ar(snd));
}).load(nrtServer);
SynthDef("dec_ramp", {
    arg out=0,  dur=30, f1,f2,f3,f4;
    var duration=BufDur.ir(0);
    var snd = PlayBuf.ar(2,0,BufRateScale.kr(0));
    snd=SelectX.ar(Line.kr(0,1,duration/4),[snd,Decimator.ar(snd,8000,8)]);
    snd = snd * EnvGen.ar(Env.new([0, 1, 1, 0], [0.005,dur-0.01,0.005]), doneAction:2);
    Out.ar(out, LeakDC.ar(snd));
}).load(nrtServer);
SynthDef("dec", {
    arg out=0,  dur=30, f1,f2,f3,f4;
    var duration=BufDur.ir(0);
    var snd = PlayBuf.ar(2,0,BufRateScale.kr(0));
    snd=Decimator.ar(snd,8000,8);
    snd = snd * EnvGen.ar(Env.new([0, 1, 1, 0], [0.005,dur-0.01,0.005]), doneAction:2);
    Out.ar(out, LeakDC.ar(snd));
}).load(nrtServer);
SynthDef("reverberate", {
    arg out=0,  dur=30, f1,f2,f3,f4;
    var duration=BufDur.ir(0);
    var snd = PlayBuf.ar(2,0,BufRateScale.kr(0));
    snd=SelectX.ar(XLine.kr(0,1,duration/4),[snd,Greyhole.ar(snd* EnvGen.ar(Env.new([0, 1, 1, 0], [0.1,dur-0.2,0.1]), doneAction:2))]);
    snd = snd * EnvGen.ar(Env.new([0, 1, 1, 0], [0.1,dur-0.2,0.1]), doneAction:2);
    Out.ar(out, LeakDC.ar(snd));
}).load(nrtServer);
SynthDef("filter_in_out", {
    arg out=0,  dur=30, f1,f2,f3,f4;
    var duration=BufDur.ir(0);
    var snd = PlayBuf.ar(2,0,BufRateScale.kr(0));
    snd = RLPF.ar(snd,
        LinExp.kr(EnvGen.kr(Env.new([0.1, 1, 1, 0.1], [f1,dur-f1-f2,f2])),0.1,1,100,20000),
        0.6);
    snd = snd * EnvGen.ar(Env.new([0, 1, 1, 0], [0.005,dur-0.01,0.005]), doneAction:2);
    Out.ar(out, LeakDC.ar(snd));
}).load(nrtServer);
SynthDef("tapedeck", {
    arg out=0,  dur=30,f1,f2,f3,f4,
    amp=0.7,tape_wet=0.90,tape_bias=0.9,saturation=0.9,drive=0.9,
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
    Out.ar(out, LeakDC.ar(snd));
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
    Out.ar(out, LeakDC.ar(snd));
}).load(nrtServer);
SynthDef("tapestop", {
	arg out=0, dur=30,f1=120,f2=1,f3=0,f4;
	var duration=BufDur.ir(0);
	var snd;
	var bpm=f1;
	var times=f2;
	var seed=f3;
	var rate,imp,tries,prob,length_time,env;
	RandSeed.kr(Impulse.kr(0),seed);
	rate=bpm/60;
	imp=Impulse.kr(rate);
	tries = duration/(60/bpm);
	prob=times/tries;
	length_time=Demand.kr(imp,0,Dseq([1,2,4,8,16]*60/bpm,inf));
	env=EnvGen.kr(Env.new([0,1-(0.125/2),1-(0.125/2),0],[length_time/9*4,length_time/9,length_time/9*4],\sine),imp*(LFNoise0.kr(rate).range(0,1)<prob));
	rate=BufRateScale.ir(0)*(1-env);
	snd = PlayBuf.ar(2,0,rate);
	Out.ar(out, LeakDC.ar(snd));
}).load(nrtServer);
SynthDef("togglefilter", {
	arg out=0, dur=30,f1=120,f2=1,f3=0,f4;
	var duration=BufDur.ir(0);
	var snd = PlayBuf.ar(2,0,BufRateScale.kr(0));
	var bpm=f1;
	var times=f2;
    var seed=f3;
    var rate,imp,tries,prob,env,length_time;
    RandSeed.ir(Impulse.kr(0),seed);
	rate=bpm/60/16;
	imp=Impulse.kr(rate);
	tries = duration/rate;
	prob=times/tries;
	length_time=Demand.kr(imp,0,Dseq([16,24,32]*60/bpm,inf));
	env=EnvGen.kr(Env.new([0,1,1,0],[length_time/3,length_time/3,length_time/3]),imp*(LFNoise0.kr(rate).range(0,1)<prob));
	snd = RLPF.ar(snd,LinExp.kr(1-env,0,1,100,20000),0.707);
	Out.ar(out, LeakDC.ar(snd));
}).load(nrtServer);
SynthDef("toggle", {
	arg out=0, dur=30,f1=120,f2=1,f3=0,f4;
	var duration=BufDur.ir(0);
	var snd = PlayBuf.ar(2,0,BufRateScale.kr(0));
	var bpm=f1;
	var times=f2;
    var seed=f3;
    var rate,imp,tries,prob,env,length_time;
    RandSeed.ir(Impulse.kr(0),seed);
	rate=bpm/60/16;
	imp=Impulse.kr(rate);
	tries = duration/rate;
	prob=times/tries;
	length_time=Demand.kr(imp,0,Dseq([2,4,8,12,16,24,32]*60/bpm,inf));
	env=EnvGen.kr(Env.new([0,1,1,0],[length_time/12,length_time*10/12,length_time/12]),imp*(LFNoise0.kr(rate).range(0,1)<prob));
	Out.ar(out, LeakDC.ar(snd*(1-env)));
}).load(nrtServer);
SynthDef("oneworddelay", {
    arg out=0, dur=30,f1=120,f2=1,f3,f4;
    var duration=BufDur.ir(0);
    var snd = PlayBuf.ar(2,0,BufRateScale.kr(0));
    var gate,sndDelay;
    var bpm=f1;
    gate=Lag.ar(DetectSilence.ar(snd,amp:0.01,time:60/bpm/4,doneAction:0),60/bpm/4);
    sndDelay=CombC.ar(snd,60/bpm*2,60/bpm*2,5)*0.5;
    sndDelay=Pan2.ar(sndDelay,SinOsc.kr(bpm/60));
    snd=SelectX.ar(gate,[snd,sndDelay]);
    snd = snd * EnvGen.ar(Env.new([0, 1, 1, 0], [0.0005,dur-0.001,0.0005]), doneAction:2);
    Out.ar(out, LeakDC.ar(snd));
}).load(nrtServer);
SynthDef("sidechain", {
    arg out=0,  dur=30,f1=1,f2=1,f3,f4;
    var duration=BufDur.ir(0);
    var bcimp=Impulse.kr(f1/duration); // times per song
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
    Out.ar(out, LeakDC.ar(snd));
}).load(nrtServer);
SynthDef("kick", {
    arg out=0,  dur=30,f1=1,f2=1,f3,f4;
    var duration=BufDur.ir(0);
    var bcimp=Dust.kr(f1/duration); // times per song
    var timescale=f2*TChoose.kr(bcimp,[0.5,0.75,1,1.25,1.5,1.75]);
    var bctrig=Trig.kr(bcimp,timescale/2);
    var bc=EnvGen.kr(Env.asr(timescale/2,1,timescale/2,[TChoose.kr(bcimp,[-1,-2,-4,-8]),TChoose.kr(bcimp,[2,4,8,16])]),bctrig);
    var snd = PlayBuf.ar(2,0,BufRateScale.kr(0));
    var qntrig=Impulse.kr(f1/duration); // f3 = bpm
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
    Out.ar(out, LeakDC.ar(snd));
}).load(nrtServer);

scoreFn={
    arg inFile,outFile,synthDefinition,durationScaling,scDoneFile,f1,f2,f3,f4;
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
                        ["writing",scDoneFile].postln;
                        0.15.wait;
                        File.new(scDoneFile, "w");
                        ["finished: wrote ",scDoneFile].postln;
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
            var f1=msg[5].asFloat;
            var f2=msg[6].asFloat;
            var f3=msg[7].asFloat;
            var f4=msg[8].asFloat;
            var scDoneFile=msg[9].asString;
            [msg, time, addr, recvPort].postln;
            scoreFn.value(inFile,outFile,synthDefinition,durationScaling,scDoneFile,f1,f2,f3,f4);
        }, '/score',recvPort:47113);
        0.15.wait;
        "writing ready file".postln;
        File.new("/tmp/nrt-scready", "w");
        ["ready, wrote ","/tmp/nrt-scready"].postln;
    }.play;
});
)
