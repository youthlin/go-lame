package lame

import (
	"errors"
	"io"
)

// A helper for liblame, which is able to,
// 1. set quality
// 2. support mono (single channel) and stereo (2 channels ) mode
// 3. change sample rate
// 4. Big-endian and little-endian (input file)

type (
	// EncodeOptions is the options for encoder.
	// 编码的选项，包括输入输出文件采样率，声道数等设置项。
	EncodeOptions struct {
		InBigEndian     bool // true if it is in big-endian (default: false)
		InSampleRate    int  // Hz, e.g., 8000, 16000, 12800, 44100, etc. (default: 24000)
		InNumChannels   int  // count of channels, for mono ones, please remain 1, and 2 if stereo
		InBitsPerSample int  // typically 16 should be working fine. the bit count of each sample, e.g., 2Bytes/sample->16bits

		OutSampleRate  int // Hz (default: same to input)
		OutNumChannels int // channels of output, only support 1 or 2 (default: same to input)
		OutQuality     int // quality: 0-highest, 9-lowest (default: 0)
	}
	// Writer wraps the Output, Lame and EncodeOptions.
	// 编码器，将输入文件写入这个结构，会编码为 mp3 格式，并将结果写入输出文件。
	Writer struct {
		output io.Writer
		lame   *Lame
		EncodeOptions
	}
)

// ErrUnsupportedChannelNum is a error, throws when channels are unsupported.
// 输入文件声道数配置错误时返回。输入文件只能有 1 或 2 个声道。
var ErrUnsupportedChannelNum = errors.New("only 1 and 2 channels are supported")

// ErrNotInited is a error, throws then Write to a nil Writer.
// 当往一个未实例化的 Writer 中写入数据时抛出。
var ErrNotInited = errors.New("not inited, please use NewWriter first")

// NewWriter create a new writer, without initializing the Lame.
// 实例化一个带有默认选项的编码器，在写入数据前都可以对返回的实例进行配置修改。
func NewWriter(output io.Writer) (*Writer, error) {
	lame, err := NewLame()
	if err != nil {
		return nil, err
	}
	return &Writer{
		output: output,
		lame:   lame,
		EncodeOptions: EncodeOptions{
			InBigEndian:     false,
			InSampleRate:    24000,
			InNumChannels:   1,
			InBitsPerSample: 16,
			OutSampleRate:   0,
			OutNumChannels:  0,
			OutQuality:      0,
		},
	}, nil
}

// ForceUpdateParams forced to init the params inside.
// NOT NECESSARY.
// 强制更新内部设置，写入数据时会自动更新，外部没有必要关注改方法。
func (w *Writer) ForceUpdateParams() (err error) {
	if w.OutSampleRate == 0 {
		w.OutSampleRate = w.InSampleRate
	}
	if w.OutNumChannels == 0 {
		w.OutNumChannels = w.InNumChannels
	}
	var mode = MODE_STEREO // two channel
	if w.OutSampleRate == 1 {
		mode = MODE_MONO // one channel
	}

	if err = w.lame.SetInSampleRate(w.InSampleRate); err != nil {
		return
	}
	if err = w.lame.SetOutSampleRate(w.OutSampleRate); err != nil {
		return
	}
	if err = w.lame.SetNumChannels(w.InNumChannels); err != nil {
		return
	}
	if err = w.lame.SetMode(mode); err != nil {
		return
	}
	if err = w.lame.SetQuality(w.OutQuality); err != nil {
		return
	}
	if err = w.lame.InitParams(); err != nil {
		return
	}
	return nil
}

// Write NOT thread-safe!
// will check if we have lame object inside first!
// TODO: support 8bit/32bit!
// currently, we support 16bit depth only
// 写入数据并进行编码。非协程安全！目前只支持 16 位。
func (w *Writer) Write(p []byte) (n int, err error) {
	if w == nil || w.output == nil || w.lame == nil {
		return 0, ErrNotInited
	}
	if w.InNumChannels != 1 && w.InNumChannels != 2 {
		return 0, ErrUnsupportedChannelNum
	}
	if !w.lame.paramUpdated {
		if err = w.ForceUpdateParams(); err != nil {
			return 0, err
		}
	}

	var samples = make([]int16, len(p)/2)
	var lo, hi int16 = 0x0001, 0x0100
	if w.InBigEndian {
		lo, hi = 0x0100, 0x0001
	}
	for i := 0; i < len(samples); i++ {
		samples[i] = int16(p[i*2])*lo + int16(p[i*2+1])*hi
	}
	// inSample * (inRate / outRate) / (inNumChan / outNumChan)
	var outSampleCount = int(int64(len(samples)) * int64(w.InSampleRate) / int64(w.OutSampleRate) * int64(w.OutNumChannels) / int64(w.InNumChannels))
	var mp3BufSize = int(1.25*float32(outSampleCount) + 7200) // follow the instruction from LAME
	var mp3Buf = make([]byte, mp3BufSize)

	if w.InNumChannels == 1 {
		n, err = w.lame.EncodeInt16(samples, samples, mp3Buf)
	} else if w.InNumChannels == 2 {
		n, err = w.lame.EncodeInt16Interleaved(samples, mp3Buf)
	}
	if err != nil {
		return 0, err
	} else {
		_, err = w.output.Write(mp3Buf[:n])
		return 2 * len(samples), err
	}
}

// Close close the Writer, and write the residual data to output(if has).
// 如果有未写入 output 的数据，调用该方法将 flush 给 output.
func (w *Writer) Close() error {
	// try to get some residual data
	if residual, err := w.lame.EncodeFlush(); err != nil {
		return err
	} else {
		if len(residual) > 0 {
			_, err = w.output.Write(residual)
		}
		return err
	}
}
