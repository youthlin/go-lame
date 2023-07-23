# go-lame

Yet another simple wrapper of libmp3lame for golang. 
It focuses on converting __raw PCM__, and __WAV__ files into mp3. 

> Forked from https://github.com/sunicy/go-lame
> 
> fixed [relative import paths are not supported in module mode](https://github.com/sunicy/go-lame/issues/4)
> (sunicy/go-lame#issues-4)

## Install
```
go get -u github.com/youthlin/go-lame
```

### Cmdline tool
- [mp3-encoder](./cmd/mp3-encoder/)

```
go install github.com/youthlin/go-lame/cmd/mp3-encoder@latest
# execute the binary to show usage
mp3-incoder
```

## Examples

### PCM to MP3
```go
func PcmToMp3(pcmFileName, mp3FileName string) {
	pcmFile, _ := os.OpenFile(pcmFileName, os.O_RDONLY, 0555)
	mp3File, _ := os.OpenFile(mp3FileName, os.O_CREATE | os.O_WRONLY | os.O_TRUNC, 0755)
	defer mp3File.Close()
	wr, err := lame.NewWriter(mp3File)
	if err != nil {
		panic("cannot create lame writer, err: " + err.Error())
	}
	wr.InSampleRate = 24000  // input sample rate, default: 24000
	wr.InNumChannels = 1     // number of channels: 1, default: 1
	wr.OutSampleRate = 24000 // output sample rate, default 0: same to input
	wr.OutNumChannels = 1    // default 0: same to input
	wr.OutQuality = 0        // 0: highest(default); 9: lowest 
	
	io.Copy(wr, pcmFile)
	wr.Close()
}
```

### WAV to MP3

```go
func WavToMp3(wavFileName, mp3FileName string) {
	// open files
	wavFile, _ := os.OpenFile(wavFileName, os.O_RDONLY, 0555)
	mp3File, _ := os.OpenFile(mp3FileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	defer mp3File.Close()

	// parsing wav info
	// NOTE: reader position moves even if it is not a wav file
	wavHdr, err := lame.ReadWavHeader(wavFile)
	if err != nil {
		panic("not a wav file, err=" + err.Error())
	}

	wr, _ := lame.NewWriter(mp3File)
	wr.EncodeOptions = wavHdr.ToEncodeOptions()
	io.Copy(wr, wavFile) // wavFile's pos has been changed!
	wr.Close()
}
```

## Roadmap

- [x] Wrapping functions from libmp3lame
- [x] WavFile parsing support
- [x] Shortcut to using wrapped functions
- [x] Supporting parsing both little-endian and big-endian PCM files
- [ ] Thorough tests 
- [ ] Supporting bit depth other than 16