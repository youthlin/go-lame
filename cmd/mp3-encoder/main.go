package main

import (
	"embed"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/youthlin/go-lame"
	"github.com/youthlin/t"
)

//go:embed *.po
var poFiles embed.FS
var args appArgs

func main() {
	parseArgs()
	if args.lang == "" {
		t.LoadFS(poFiles)
	} else {
		t.Load(args.lang)
	}

	if args.input == "" || args.output == "" {
		printUsage()
		fmt.Println(t.T("[Error] Both input file and output file are required.\n"))
		os.Exit(1)
	}

	input, err := os.Open(args.input)
	if err != nil {
		fmt.Println(t.T("failed to open input file %q: %+v", args.input, err))
		os.Exit(1)
	}
	defer input.Close()
	output, err := os.OpenFile(args.output, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		fmt.Println(t.T("failed to open output file %q: %+v", args.output, err))
		os.Exit(1)
	}
	defer output.Close()

	wr, err := lame.NewWriter(output)
	if err != nil {
		fmt.Println(t.T("failed to create lame: %+v", err))
		os.Exit(1)
	}
	defer wr.Close()
	hdr, err := lame.ReadWavHeader(input)
	if err == nil { // is wav
		wr.EncodeOptions = hdr.ToEncodeOptions()
		// read header, no need to seek to start
	} else { // may pcm
		// need seek to start
		input.Seek(0, io.SeekStart)
		wr.EncodeOptions = args.EncodeOptions
	}

	_, err = io.Copy(wr, input)
	if err != nil {
		fmt.Println(t.T("failed to encode: %+v", err))
		os.Exit(1)
	}
}

type appArgs struct {
	lame.EncodeOptions
	input  string
	output string
	lang   string
}

func parseArgs() {
	flag.StringVar(&args.input, "i", "", "")
	flag.StringVar(&args.output, "o", "", "")
	flag.StringVar(&args.lang, "l", "", "")
	flag.BoolVar(&args.InBigEndian, "inBigEndian", false, "")
	flag.IntVar(&args.InSampleRate, "inSampleRate", 24000, "")
	flag.IntVar(&args.InNumChannels, "inChannels", 1, "")
	flag.IntVar(&args.InBitsPerSample, "inBits", 16, "")
	flag.IntVar(&args.OutSampleRate, "outSampleRate", 0, "")
	flag.IntVar(&args.OutNumChannels, "outChannels", 0, "")
	flag.IntVar(&args.OutQuality, "quality", 0, "")
	flag.Usage = printUsage
	flag.Parse()
}

func printUsage() {
	fmt.Println()
	fmt.Println(t.T("Usage: %s -i <input file> -o <output file> [settings]", os.Args[0]))
	fmt.Println(t.T("  -i <input file>\tthe input file name, wav or pcm"))
	fmt.Println(t.T("  -o <output file>\tthe output mp3 file name"))
	fmt.Println(t.T("  [settings]"))
	fmt.Println(t.T("    -inBigEndian[=false]\tif the input file is in big-endian (default false)"))
	fmt.Println(t.T("    -inSampleRate <Hz>\t\tsample rate of input file (default 24000)"))
	fmt.Println(t.T("    -inChannels <num>\t\tchannels of input file, 1 or 2 (default 1)"))
	fmt.Println(t.T("    -inBits <num>\t\tthe bit count of each sample (default 16)"))
	fmt.Println(t.T("    -outSampleRate <Hz>\t\tsample rate of output file (default 0, means same of input)"))
	fmt.Println(t.T("    -outChannels <num>\t\tchannels of output file, 1 or 2 (default 0, means same of input)"))
	fmt.Println(t.T("    -quality <num>\t\tquality, 0-9, 0-highest, 9-lowest (default 0)"))
	fmt.Println(t.T("    -lang <path>\t\tpath to po/mo file or dir"))
	fmt.Println()
}
