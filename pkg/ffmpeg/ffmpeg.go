package ffmpeg

import (
	"io/ioutil"
	"os"

	"github.com/xfrr/goffmpeg/transcoder"
)

func Amplify(input string) (*os.File, error) {
	trans := new(transcoder.Transcoder)

	outputPath, err := ioutil.TempFile("", "vm.*.oga")
	if err != nil {
		return nil, err
	}

	err = trans.Initialize(input, outputPath.Name())

	if err != nil {
		return nil, err
	}

	trans.MediaFile().SetAudioFilter("loudnorm=I=-24:LRA=5:TP=-3")
	trans.MediaFile().SetAudioCodec("libopus")

	done := trans.Run(false)

	err = <-done

	return outputPath, nil
}
