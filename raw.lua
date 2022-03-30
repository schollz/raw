math.randomseed(os.time())

SENDOSC="sendosc"
WORKDIR="/tmp/raw/"
PROGRESSFILE="/tmp/raw/progress"
NRTREADY="/tmp/nrt-scready"
debugging=true
function os.file_exists(name)
  local f=io.open(name,"r")
  if f~=nil then
    io.close(f)
    return true
  else
    return false
  end
end

function os.capture(cmd,raw)
  local f=assert(io.popen(cmd,'r'))
  local s=assert(f:read('*a'))
  f:close()
  if raw then return s end
  s=string.gsub(s,'^%s+','')
  s=string.gsub(s,'%s+$','')
  s=string.gsub(s,'[\n\r]+',' ')
  return s
end

function os.cmd(cmd)
  if debugging then
    print(cmd)
  end
  os.execute(cmd.." 2>&1")
end

local snare_file="snare.wav"
local kick_file="kick.wav"
local debugging=false
local save_onset=false
local charset={}

-- qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM1234567890
for i=48,57 do table.insert(charset,string.char(i)) end
for i=65,90 do table.insert(charset,string.char(i)) end
for i=97,122 do table.insert(charset,string.char(i)) end

function string.random(length)
  if length>0 then
    return string.random(length-1)..charset[math.random(1,#charset)]
  else
    return ""
  end
end

function string.random_filename(suffix,prefix)
  suffix=suffix or ".wav"
  prefix=prefix or WORKDIR.."breaktemp-"
  return prefix..string.random(8)..suffix
end

function math.round(number,quant)
  if quant==0 then
    return number
  else
    return math.floor(number/(quant or 1)+0.5)*(quant or 1)
  end
end

function math.std(numbers)
  local mu=math.average(numbers)
  local sum=0
  for _,v in ipairs(numbers) do
    sum=sum+(v-mu)^2
  end
  return math.sqrt(sum/#numbers)
end

function math.average(numbers)
  if next(numbers)==nil then
    do return end
  end
  local total=0
  for _,v in ipairs(numbers) do
    total=total+v
  end
  return total/#numbers
end

function math.trim(numbers,std_num)
  local mu=math.average(numbers)
  local std=math.std(numbers)
  local new_numbers={}
  for _,v in ipairs(numbers) do
    if v>mu-(std*std_num) and v<mu+(std*std_num) then
      table.insert(new_numbers,v)
    end
  end
  return math.average(new_numbers)
end

function table.clone(org)
  return {table.unpack(org)}
end

function table.merge(t1,t2)
  n=#t1
  for i=1,#t2 do
    t1[n+i]=t2[i]
  end
end

function table.reverse(t)
  local len=#t
  for i=len-1,1,-1 do
    t[len]=table.remove(t,i)
  end
end

function table.permute(t,n,count)
  n=n or #t
  for i=1,count or n do
    local j=math.random(i,n)
    t[i],t[j]=t[j],t[i]
  end
end

function table.shuffle(tbl)
  for i=#tbl,2,-1 do
    local j=math.random(i)
    tbl[i],tbl[j]=tbl[j],tbl[i]
  end
end

function table.add(t,scalar)
  for i,_ in ipairs(t) do
    t[i]=t[i]+scalar
  end
end

function table.is_empty(t)
  return next(t)==nil
end

function table.get_rotation(t)
  local t2={}
  local v1=0
  for i,v in ipairs(t) do
    if i>1 then
      table.insert(t2,v)
    else
      v1=v
    end
  end
  table.insert(t2,v1)
  return t2
end

function table.average(t)
  local sum=0
  for _,v in pairs(t) do
    sum=sum+v
  end
  return sum/#t
end

function table.rotate(t)
  for i,v in ipairs(table.get_rotation(t)) do
    t[i]=v
  end
end

function table.rotatex(t,d)
  if d<0 then
    table.reverse(t)
  end
  local d_abs=math.abs(d)
  if d_abs>0 then
    for i=1,d_abs do
      table.rotate(t)
    end
  end
  if d<0 then
    table.reverse(t)
  end
end

function table.clone(org)
  return {table.unpack(org)}
end

function table.copy(t)
  local t2={}
  for i,v in ipairs(t) do
    table.insert(t2,v)
  end
  return t2
end

function table.print(t)
  for i,v in ipairs(t) do
    print(i,v)
  end
end

function table.print_matrix(m)
  for _,v in ipairs(m) do
    local s=""
    for _,v2 in ipairs(v) do
      s=s..v2.." "
    end
    print(s)
  end
end

function table.get_change(m)
  local total_change=0
  for col=1,#m[1] do
    local last_val=0
    for row=1,#m do
      local val=m[row][col]
      if row>1 then
        total_change=total_change+math.abs(val-last_val)
      end
      last_val=val
    end
  end
  return total_change
end

function table.minimize_row_changes(m)
  local m_=table.clone(m)
  -- generate random rotations
  local best_change=100000
  local best_m={}
  for i=1,10000 do
    -- rotate a row randomly
    local random_row=math.random(1,#m)
    m_[random_row]=table.get_rotation(m_[random_row])
    local change=table.get_change(m_)
    if change<best_change then
      best_change=change
      best_m=table.clone(m_)
    end
  end
  return best_m
  -- table.print_matrix(best_m)
end

function table.contains(t,x)
  for _,v in ipairs(t) do
    if v==x then
      do return true end
    end
  end
  return false
end

-- lfo goes from 0 to 1
function math.lfo(t,period,phase)
  return (math.sin(2*3.14159265/period+phase)+1)/2
end

local audio={}

function audio.tempo(fname)
  local audio_length=audio.length(fname)
  local closet_bpm={0,100000}
  for bpm=100,200 do 
    local measures=audio_length/((60/bpm)*4)
    if math.round(measures)%2==0 then 
      local dif=math.abs(math.round(measures)-measures)
      dif=dif-math.round(measures)/60
      -- print(bpm,math.round(measures),measures,dif)
      if dif<closet_bpm[2] then 
        closet_bpm[2]=dif 
        closet_bpm[1]=bpm
      end
    end
  end
  return closet_bpm[1]
end

function audio.length(fname)
  local s=os.capture("sox "..fname.." -n stat 2>&1  | grep Length | awk '{print $3}'")
  return tonumber(s)
end

function audio.mean_norm(fname)
  local s=os.capture("sox "..fname.." -n stat 2>&1  | grep Mean | grep norm | awk '{print $3}'")
  return tonumber(s)
end

function audio.silence_add(fname,silence_length)
  local sample_rate,channels=audio.get_info(fname)
  local silence_file=string.random_filename()
  local fname2=string.random_filename()
  -- first create the silence
  os.cmd("sox -n -r "..sample_rate.." -c "..channels.." "..silence_file.." trim 0.0 "..silence_length)
  -- combine with original file
  os.cmd("sox "..fname.." "..silence_file.." "..fname2)
  os.cmd("rm -f "..silence_file)
  return fname2
end

function audio.silence_trim(fname)
  local fname2=string.random_filename()
  os.cmd("sox "..fname.." "..fname2.." silence 1 0.1 0.025% reverse silence 1 0.1 0.025% reverse")
  return fname2
end

function audio.trim(fname,start,length)
  local fname2=string.random_filename()
  if length==nil then
    os.cmd("sox "..fname.." "..fname2.." trim "..start)
  else
    os.cmd("sox "..fname.." "..fname2.." trim "..start.." "..length)
  end
  return fname2
end

function audio.reverse(fname)
  local fname2=string.random_filename()
  os.cmd(string.format("sox %s %s reverse",fname,fname2))
  return fname2
end

function audio.pitch(fname,notes)
  local fname2=string.random_filename()
  os.cmd(string.format("sox %s %s pitch %d",fname,fname2,notes*100))
  return fname2
end

function audio.join(fnames)
  local fname2=string.random_filename()
  os.cmd(string.format("sox %s %s",table.concat(fnames," "),fname2))
  return fname2
end

function audio.repeat_n(fname,repeats)
  local fname2=string.random_filename()
  os.cmd(string.format("sox %s %s repeat %d",fname,fname2,repeats))
  return fname2
end

function audio.retempo_speed(fname,old_tempo,new_tempo)
  local fname2=string.random_filename()
  os.cmd(string.format("sox %s %s speed %f rate -v 48k",fname,fname2,new_tempo/old_tempo))
  return fname2
end

function audio.retempo_stretch(fname,old_tempo,new_tempo)
  local fname2=string.random_filename()
  os.cmd(string.format("sox %s %s tempo -m %f",fname,fname2,new_tempo/old_tempo))
  return fname2
end

function audio.slowdown(fname,slowdown)
  local fname2=string.random_filename()
  os.cmd(string.format("sox %s %s tempo -m %f",fname,fname2,slowdown))
  return fname2
end

function audio.get_info(fname)
  local sample_rate=tonumber(os.capture("sox --i "..fname.." | grep 'Sample Rate' | awk '{print $4}'"))
  local channels=tonumber(os.capture("sox --i "..fname.." | grep 'Channels' | awk '{print $3}'"))
  return sample_rate,channels
end

-- copy_and_paste2 finds best positionn, but does not keep timing
function audio.copy_and_paste2(fname,copy_start,copy_stop,paste_start)
	local copy_length=copy_stop-copy_start
  if copy_length==nil or copy_length<0.05 then 
    do return fname end 
  end
	local piece=string.random_filename()
	local part1=string.random_filename()
	local part2=string.random_filename()
	local fname2=string.random_filename()
	local e=5/1000 
	local l=5/1000
	os.cmd(string.format("sox %s %s trim %f %f",fname,piece,copy_start-e-l,copy_stop-copy_start+e+l+e))
	os.cmd(string.format("sox %s %s trim 0 %f",fname,part1,paste_start+e))
	os.cmd(string.format("sox %s %s trim %f",fname,part2,paste_start+copy_stop-copy_start-e-l))
	os.cmd(string.format("sox %s %s %s %s splice %f %f",part1,piece,part2,fname2,paste_start+e,paste_start+e+copy_stop-copy_start+e+l+e))
  os.cmd(string.format("rm -f %s %s %s",part1,part2,piece))
	return fname2
end

function audio.copy_and_paste(fname,copy_start,copy_stop,paste_start,crossfade)
	local copy_length=copy_stop-copy_start
  if copy_length==nil or copy_length<0.05 then 
    do return fname end 
  end
	local piece=string.random_filename()
	local part1=string.random_filename()
	local part2=string.random_filename()
	local fname2=string.random_filename()
	local splice1=string.random_filename()
	local e=crossfade or 0.1 
	local l=0 -- no leeway
	os.cmd(string.format("sox %s %s trim %f %f",fname,piece,copy_start-e,copy_length+2*e))
	os.cmd(string.format("sox %s %s trim 0 %f",fname,part1,paste_start+e))
	os.cmd(string.format("sox %s %s trim %f",fname,part2,paste_start+copy_length-e))
	os.cmd(string.format("sox %s %s %s splice %f,%f,%f",part1,piece,splice1,paste_start+e,e,l))
	os.cmd(string.format("sox %s %s %s splice %f,%f,%f",splice1,part2,fname2,paste_start+copy_length+e,e,l))
  os.cmd(string.format("rm -f %s %s %s %s",piece,part1,part2,splice1))
	return fname2
end

-- pastes any piece into a place in the audio
-- assumes that the piece has "crossfade" length on both sides
-- in addition to its current length
function audio.paste(fname,piece,paste_start,crossfade)
	local copy_length=audio.length(piece)
  if copy_length==nil then 
    do return fname end 
  end
	local part1=string.random_filename()
	local part2=string.random_filename()
	local fname2=string.random_filename()
	local splice1=string.random_filename()
	local e=crossfade or 0.1 
	local l=0 -- no leeway
	os.cmd(string.format("sox %s %s trim 0 %f",fname,part1,paste_start+e))
	os.cmd(string.format("sox %s %s trim %f",fname,part2,paste_start+copy_length-e*3))
	os.cmd(string.format("sox %s %s %s splice %f,%f,%f",part1,piece,splice1,paste_start+e,e,l))
	os.cmd(string.format("sox %s %s %s splice %f,%f,%f",splice1,part2,fname2,paste_start+copy_length+e,e,l))
  os.cmd(string.format("rm -f %s %s %s",part1,part2,splice1))
	return fname2
end

function audio.sample_rate(fname,sr,ch)
  sr=sr or 48000
  ch=ch or 2
  local fname2=string.random_filename()
  os.cmd(string.format("sox -r %d %s %s",sr,fname,fname2))
  return fname2
end

function audio.gain(fname,gain)
  local fname2=string.random_filename()
  os.cmd(string.format("sox %s %s gain %f",fname,fname2,gain))
  return fname2
end

function audio.stretch(fname,stretch)
  local fname2=string.random_filename()
  os.cmd(string.format("sox %s %s stretch %f",fname,fname2,stretch))
  return fname2
end

function audio.stutter(fname,stutter_length,pos_start,count,crossfade_piece,crossfade_stutter,gain_amt)
	crossfade_piece=0.1 or crossfade_piece
	crossfade_stutter=0.005 or crossfade_stutter
	local partFirst=string.random_filename()
	local partMiddle=string.random_filename()
	local partLast=string.random_filename()
	os.cmd(string.format("sox %s %s trim %f %f",fname,partFirst,pos_start-crossfade_piece,stutter_length+crossfade_piece+crossfade_stutter))
	os.cmd(string.format("sox %s %s trim %f %f",fname,partMiddle,pos_start-crossfade_stutter,stutter_length+crossfade_stutter+crossfade_stutter))
	os.cmd(string.format("sox %s %s trim %f %f",fname,partLast,pos_start-crossfade_stutter,stutter_length+crossfade_piece+crossfade_stutter))
  gain_amt=gain_amt or (count>8 and -1.5 or -2)
	for i=1,count do 
		local fnameNext=""
		if i==1 then 
			fnameNext=audio.gain(partFirst,gain_amt*(count-i))
		else
			fnameNext=string.random_filename()
      local fnameMid=i<count and partMiddle or partLast 
      if gain_amt~=0 then 
        fnameMid=audio.gain(fnameMid,gain_amt*(count-i))
      end
			os.cmd(string.format("sox %s %s %s splice %f,%f,0",fname2,fnameMid,fnameNext,audio.length(fname2),crossfade_stutter))
		end
		fname2=fnameNext
	end
	return fname2
end

function audio.supercollider_effect(fname,effect,f1,f2,f3,f4)
  local fname2=string.random_filename()
  local durationScaling=1 
  if effect=="reverberate" then 
    durationScaling=4
  end
  f1=f1 or 0
  f2=f2 or 0
  f3=f3 or 0
  f4=f4 or 0
  os.cmd(string.format(SENDOSC..' --host 127.0.0.1 --addr "/score" --port 47113 --recv-port 47888 -s %s -s %s -s %s -s %s -s 47888 -s %f -s %f -s %f -s %f',fname,fname2,effect,durationScaling,f1,f2,f3,f4))
  return fname2
end


function run()
  local progress=0
  local fname="sample.aiff"
  local fname_out="result.wav"
  local target_beats=32
  local new_tempo=nil
  local input_tempo=nil
  local p_reverse=10
  local p_stutter=5
  local p_slow=0
  local p_jump=80
  local p_pitch=10
  local p_reverb=2
  local p_revreverb=2
  local print_help=false
  local tapedeck=false
  local server_started=false
  local retempo_switch=1 
  local RETEMPO_SPEED=1 
  local RETEMPO_STRETCH=2
  local RETEMPO_NONE=3
  local filter_in_beats=4
  local filter_out_beats=4
  local add_kick=false
  local add_sidechain=false
  for i,v in ipairs(arg) do
    if string.find(v,"input") and string.find(v,"tempo") then
      input_tempo=tonumber(arg[i+1]) or input_tempo
    elseif string.find(v,"server") and string.find(v,"started") then
      server_started=true
    elseif string.find(v,"filter") and string.find(v,"in") then
      filter_in_beats=tonumber(arg[i+1])
    elseif string.find(v,"filter") and string.find(v,"out") then
      filter_out_beats=tonumber(arg[i+1])
    elseif string.find(v,"retempo") and string.find(v,"repitch") then
      retempo_switch=RETEMPO_SPEED
    elseif string.find(v,"retempo") and string.find(v,"stretch") then
      retempo_switch=RETEMPO_STRETCH
    elseif string.find(v,"retempo") and string.find(v,"none") then
      retempo_switch=RETEMPO_NONE
    elseif string.find(v,"save") and string.find(v,"onset") then
      save_onset=true
    elseif string.find(v,"strip") and string.find(v,"silence") then
      strip_silence=true
    elseif string.find(v,"make") and string.find(v,"movie") then
      make_movie=true
    elseif string.find(v,"input") and string.find(v,"file") then
      fname=arg[i+1]
      for j=2,100 do
        local a=arg[i+j]
        if a==nil or string.find(a,"-")  then
          break
        end
        fname=fname.." "..a
      end
      if string.find(fname," ") then
        local fname2=string.random_filename()
        print("cp '"..fname.."' "..fname2)
        os.cmd("cp '"..fname.."' "..fname2)
        fname=fname2
      end
    elseif string.find(v,"-output") then
      fname_out=arg[i+1]
    elseif string.find(v,"reverse") then
      p_reverse=tonumber(arg[i+1]) or p_reverse
    elseif string.find(v,"slowdown") then
      p_slow=tonumber(arg[i+1]) or p_slow
    elseif string.find(v,"stutter") then
      p_stutter=tonumber(arg[i+1]) or p_stutter
    elseif string.find(v,"tapedeck") then
      tapedeck=true
    elseif string.find(v,"addkick") then
      add_kick=true
    elseif string.find(v,"addsidechain") then
      add_sidechain=true
    elseif string.find(v,"pitch") then
      p_pitch=tonumber(arg[i+1]) or p_pitch
    elseif string.find(v,"jump") then
      p_jump=tonumber(arg[i+1]) or p_jump
    elseif string.find(v,"revreverb") then
      p_revreverb=tonumber(arg[i+1]) or p_revreverb
    elseif string.find(v,"reverb") then
      p_reverb=tonumber(arg[i+1]) or p_reverb
    elseif string.find(v,"help") then
      print_help=true
    elseif string.find(v,"target") and string.find(v,"beats") then
      target_beats=tonumber(arg[i+1])
    elseif string.find(v,"target") and string.find(v,"tempo") then
      new_tempo=tonumber(arg[i+1]) or new_tempo
      if new_tempo then 
        new_tempo=math.floor(new_tempo)
      end
    elseif string.find(v,"-sendosc") then
      SENDOSC=arg[i+1]
    elseif string.find(v,"-debug") then
      debugging=true
    end
  end

  if #arg<2 or print_help then
    print([[NAME
   
      dnb.lua - generative drum & bass
   
  DESCRIPTION
   
    -i, --input string
        input filename
   
    --input-tempo value
        tempo of input file (defaults to determine automatically)
   
    -o, --output string
        output filename
   
    -b, --beats value
        number of beats
   
    -t, --tempo value
        tempo of generated beat
   
    --debug
        debug mode (lots of output)
   
    --reverse value
        probability of reversing (0-100%, default 10%)
   
    --stutter value
        probability of stutter (0-100%, default 5%)
   
    --pitch value
        probability of pitch up (0-100%, default 10%)
   
    --reverb value
        probability of adding reverb tail to kick/snare (0-100%, default 2%)
  ]])
  else
    if os.capture("command -v sox")=="" then
      print("need to install sox first")
    else
      os.cmd("rm -rf "..WORKDIR)
      os.cmd("mkdir -p "..WORKDIR)
      os.cmd("echo 0 > "..PROGRESSFILE)
      if not server_started then 
        os.cmd(SENDOSC..' --host 127.0.0.1 --addr "/quit" --port 47113')
        os.cmd("sclang raw.sc &")
      end
      while not os.file_exists(NRTREADY) do 
        os.execute("sleep 0.1")
      end
      -- convert to 48000
      original_fname=fname
      fname=audio.sample_rate(fname,48000,2)
      if not os.file_exists(fname) then
        print("!!!! ERROR !!!!")
        do return end 
      end
      fname=audio.silence_trim(fname)
      if not os.file_exists(fname) then
        print("!!!! ERROR !!!!")
        do return end 
      end
      local bpm=nil
      if input_tempo==nil then 
        for word in string.gmatch(original_fname, '([^_]+)') do
          print(word)
          if string.find(word,"bpm") then 
            bpm=word:match("bpm%d+")
          end
        end
        if bpm~=nil then
          bpm=bpm:match("%d+")
        end
        if bpm==nil then 
          bpm=audio.tempo(fname)
        end
        local bpms={bpm/4,bpm/2,bpm,bpm*2,bpm*4}
        local closest=100000
        for _, b in ipairs(bpms) do 
          local diff=math.abs(new_tempo-b)
          if diff<closest then 
            closest=diff
            bpm=b
          end          
        end
      else
        bpm=input_tempo 
      end
      print("bpm: "..bpm)
      local beats=math.floor(math.round(audio.length(fname)/(60/bpm)))

      fname=audio.silence_add(fname,2)
      if not os.file_exists(fname) then
        print("!!!! ERROR !!!!")
        do return end 
      end

      fname=audio.trim(fname,0,beats*60/bpm)
      if not os.file_exists(fname) then
        print("!!!! ERROR !!!!")
        do return end 
      end

      beats=audio.length(fname)/(60/bpm)
      while beats<target_beats do 
        fname=audio.repeat_n(fname,2)
        beats=audio.length(fname)/(60/bpm)
        print(beats)
      end
      fname=audio.trim(fname,0,target_beats*60/bpm)
      if not os.file_exists(fname) then
        print("!!!! ERROR !!!!")
        do return end 
      end

      local total_beats=math.floor(audio.length(fname)/(60/bpm))
      print(total_beats)
      os.cmd("cp "..fname.." original.wav")
      local total_things=math.floor(total_beats*p_pitch/100)
        +math.floor(total_beats*p_jump/100)
        +math.floor(total_beats*p_reverse/100)
        +math.floor(total_beats*p_revreverb/100)
        +math.floor(total_beats*p_reverb/100)
        +math.floor(total_beats*p_slow/100)
        +math.floor(total_beats*p_stutter/100/2)*2+3

      local fname_original=fname
      -- copy and pitch and paste
      for i=1,math.floor(total_beats*p_pitch/100) do 
        progress=progress+1
        os.cmd('echo '..(math.round(progress/total_things*1000)/10).." >> "..PROGRESSFILE)
        local length_beat=math.random(1,4)/8
        local start_beat=math.random(1,total_beats-length_beat*8)
        local paste_beat=start_beat
        local crossfade=0.095
        local piece=audio.pitch(audio.trim(fname_original,60/bpm*start_beat-crossfade,60/bpm*length_beat+crossfade*2),math.random(2,3))
        fname=audio.paste(fname,piece,60/bpm*paste_beat,crossfade)
        if not os.file_exists(fname) then
          print("!!!! ERROR !!!!")
          do return end 
        end
      end
      -- basic copy and paste
      for i=1,math.floor(total_beats*p_jump/100) do 
        progress=progress+1
        os.cmd('echo '..(math.round(progress/total_things*1000)/10).." >> "..PROGRESSFILE)
        local length_beat=math.random(1,3) -- eigth notes
        local start_beat=math.random(1,math.floor(total_beats-length_beat/2))*2
        local paste_beat=math.random(1,math.floor(total_beats-length_beat/2))*2
        local crossfade=0.05
        fname=audio.copy_and_paste(fname,60/bpm/2*start_beat,60/bpm/2*(start_beat+length_beat),60/bpm/2*paste_beat,crossfade)
        if not os.file_exists(fname) then
          print("!!!! ERROR !!!!")
          do return end 
        end
      end
      -- copy and reverse and paste
      for i=1,math.floor(total_beats*p_reverse/100) do 
        progress=progress+1
        os.cmd('echo '..(math.round(progress/total_things*1000)/10).." >> "..PROGRESSFILE)
        local length_beat=math.random(1,3) --eight notes
        local start_beat=math.random(1,math.floor(total_beats-length_beat/2))*2
        local paste_beat=math.random(1,math.floor(total_beats-length_beat/2))*2
        local crossfade=0.05
        local piece=audio.reverse(audio.trim(fname_original,60/bpm/2*start_beat-crossfade,60/bpm/2*length_beat+crossfade*2))
        fname=audio.paste(fname,piece,60/bpm/2*paste_beat,crossfade)
        if not os.file_exists(fname) then
          print("!!!! ERROR !!!!")
          do return end 
        end
      end      
      -- copy and reverberate and paste
      for i=1,math.floor(total_beats*p_reverb/100) do 
        progress=progress+1
        os.cmd('echo '..(math.round(progress/total_things*1000)/10).." >> "..PROGRESSFILE)
        local length_beat=math.random(1,2)
        local start_beat=math.random(3,total_beats-length_beat-1)
        local paste_beat=math.random(2,math.floor(total_beats-length_beat*2))*2
        local crossfade=0.1
        local piece=audio.trim(fname_original,60/bpm*start_beat-crossfade,60/bpm/4*length_beat+crossfade*2)
        piece=audio.supercollider_effect(piece,"reverberate")
        fname=audio.paste(fname,piece,60/bpm/2*paste_beat,crossfade)
        if not os.file_exists(fname) then
          print("!!!! ERROR !!!!")
          do return end 
        end
      end
      -- copy and reverberate and reverse and paste
      for i=1,math.floor(total_beats*p_revreverb/100) do 
        progress=progress+1
        os.cmd('echo '..(math.round(progress/total_things*1000)/10).." >> "..PROGRESSFILE)
        local length_beat=math.random(1,4)
        local start_beat=math.random(1,total_beats-length_beat)
        local paste_beat=math.random(1,2*(total_beats-length_beat))
        local crossfade=0.01
        local piece=audio.trim(fname_original,60/bpm*start_beat-crossfade,60/bpm/4*length_beat+crossfade*2)
        piece=audio.supercollider_effect(piece,"reverberate")
        piece=audio.reverse(piece)
        fname=audio.paste(fname,piece,60/bpm/2*paste_beat,crossfade)
        if not os.file_exists(fname) then
          print("!!!! ERROR !!!!")
          do return end 
        end
      end
      -- copy and stutter and paste
      for i=1,math.floor(total_beats*p_stutter/100/2) do 
        progress=progress+1
        os.cmd('echo '..(math.round(progress/total_things*1000)/10).." >> "..PROGRESSFILE)
        local crossfade=0.05
        local beat_start=math.random(4,total_beats-4)
        local stutters=math.random(1,3)*4
        local piece=audio.stutter(fname_original,60/bpm/4,60/bpm*beat_start,stutters,crossfade,0.001,nil)
        piece=audio.supercollider_effect(piece,"lpf_rampup")
        fname=audio.paste(fname,piece,60/bpm/4*math.random(12,total_beats*4-16),crossfade)
        if not os.file_exists(fname) then
          print("!!!! ERROR !!!!")
          do return end 
        end
      end
      -- copy and slow and paste
      for i=1,math.floor(total_beats*p_slow/100) do 
        progress=progress+1
        os.cmd('echo '..(math.round(progress/total_things*1000)/10).." >> "..PROGRESSFILE)
        local length_beat=math.random(1,3)
        local start_beat=math.random(1,total_beats-length_beat-1)
        local paste_beat=start_beat
        local crossfade=0.05
        local piece=audio.trim(fname,60/bpm*start_beat-crossfade*0.75,60/bpm*length_beat+crossfade*0.75*2)
        piece=audio.slowdown(piece,0.75)
        fname=audio.paste(fname,piece,60/bpm*paste_beat,crossfade)
        if not os.file_exists(fname) then
          print("!!!! ERROR !!!!")
          do return end 
        end
      end
      if retempo_switch~=RETEMPO_NONE and new_tempo~=nil and new_tempo~=bpm  then 
        print("retempo_switch",retempo_switch)
        if retempo_switch==RETEMPO_SPEED then 
          fname=audio.retempo_speed(fname,bpm,new_tempo)
        elseif retempo_switch==RETEMPO_STRETCH then 
          print("stretching temp",bpm,new_tempo)
          fname=audio.retempo_stretch(fname,bpm,new_tempo)
        end
        bpm=new_tempo
      end

      fname=audio.supercollider_effect(fname,"filter_in_out",filter_in_beats*60/bpm,filter_out_beats*60/bpm)
      for tries=1,30 do 
        if os.file_exists(fname) then
          break
        end
        os.execute("sleep 0.1")
      end
      if not os.file_exists(fname) then
        print("!!!! ERROR !!!!")
        do return end 
      end

      -- bitcrush
      -- local bitcrush_occurences=4
      -- local bitcrush_length=4*60/bpm
      -- fname=audio.supercollider_effect(fname,"bitcrush",bitcrush_occurences,bitcrush_length)
      if add_kick or add_sidechain then 
        local sidechain_occurences=4
        local sidechain_length=8*60/bpm
        local sidechain_bpm=bpm
        if sidechain_bpm>100 then 
          sidechain_bpm=sidechain_bpm/2
        end
        fname=audio.supercollider_effect(fname,add_kick and "kick" or "sidechain",sidechain_occurences,sidechain_length,sidechain_bpm)
      end
      if tapedeck then 
        fname=audio.supercollider_effect(fname,"tapedeck",filter_in_beats*60/bpm,filter_out_beats*60/bpm)
        for tries=1,30 do 
          if os.file_exists(fname) then
            break
          end
          os.execute("sleep 0.1")
        end
        if not os.file_exists(fname) then
          print("!!!! ERROR !!!!")
          do return end 
        end
      end
      os.cmd("mv "..fname.." "..fname_out)
      if not server_started then 
        os.cmd(SENDOSC..' --host 127.0.0.1 --addr "/quit" --port 47113')
      end
      os.cmd("rm -rf "..WORKDIR)
    end
  end
end

run()
