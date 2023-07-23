# mp3-encoder
将 pcm/wav 文件编码为 mp3.

配合 [silk](https://github.com/youthlin/silk) 库或 [silk-decoder](https://github.com/youthlin/silk/tree/main/cmd/silk-decoder) 命令行可以将国内通信软件的 amr 文件转为 mp3 格式。

## Install
```
go install github.com/youthlin/go-lame/cmd/mp3-encoder@latest
# execute to see usage
mp3-encoder
```

Usage

```
Usage: mp3-encoder -i <input file> -o <output file> [settings]
  -i <input file>       the input file name, wav or pcm
  -o <output file>      the output mp3 file name
  [settings]
    -inBigEndian[=false]        if the input file is in big-endian (default false)
    -inSampleRate <Hz>          sample rate of input file (default 24000)
    -inChannels <num>           channels of input file, 1 or 2 (default 1)
    -inBits <num>               the bit count of each sample (default 16)
    -outSampleRate <Hz>         sample rate of output file (default 0, means same of input)
    -outChannels <num>          channels of output file, 1 or 2 (default 0, means same of input)
    -quality <num>              quality, 0-9, 0-highest, 9-lowest (default 0)
    -lang <path>                path to po/mo file or dir

用法: mp3-encoder -i <输入文件> -o <输出文件> [选项]
  -i <输入文件>         wav 或 pcm 格式的输入文件名
  -o <输出文件>         输出的 mp3 文件名
  [选项]
    -inBigEndian[=false]        输入文件是否是大端序(默认值：false)
    -inSampleRate <赫兹>        输入文件的采样率(默认值：24000 赫兹)
    -inChannels <声道数>        输入文件的声道数，1 或 2（默认值：1）
    -inBits <位数>              输入文件的位数（默认值：16）
    -outSampleRate <赫兹>       输出文件采样率（默认值是 0, 表示使用和输入文件相同的值）
    -outChannels <声道数>       输出文件声道数，1 或 2（默认值是 0, 表示使用和输入文件相同的值）
    -quality <质量>             输出质量，0-9, 0 表示最好，9 表示最差（默认值：0）
    -lang <语言路径>            指向 po/mo 文件或其文件夹

```
## 翻译
提取翻译字符串
```
xgettext -C -c=TRANSLATORS: --from-code=UTF-8 -o messages.pot -kT:1  *.go
```

生成翻译文件
```
msginit -i messages.pot -l zh_CN
```

重新提取后, 更新翻译文件
```
msgmerge -U zh_CN.po messages.pot
```
