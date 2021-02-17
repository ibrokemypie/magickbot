package magick

import (
	"errors"
	"path"

	"gopkg.in/gographics/imagick.v3/imagick"
)

type MagickCommand int

const (
	EXPLODE = iota + 1
	IMPLODE
)

func RunMagick(command MagickCommand, files []string, iterations int) error {
	imagick.Initialize()
	defer imagick.Terminate()

	mw := imagick.NewMagickWand()
	defer mw.Destroy()

	for i := 0; i < iterations; i++ {
		for k, file := range files {
			err := mw.ReadImage(file)
			if err != nil {
				return (err)
			}

			switch command {
			case EXPLODE:
				{
					err = mw.ImplodeImage(-.5, imagick.INTERPOLATE_PIXEL_UNDEFINED)
					if err != nil {
						return (err)
					}
				}
			case IMPLODE:
				{
					err = mw.ImplodeImage(.5, imagick.INTERPOLATE_PIXEL_UNDEFINED)
					if err != nil {
						return (err)
					}
				}
			default:
				{
					return errors.New("Unsupported command")
				}
			}

			outputFile := file
			if i == 0 {
				outputFile = "/tmp/out." + path.Base(file)
				files[k] = outputFile
			}

			err = mw.WriteImage(outputFile)
			if err != nil {
				return (err)
			}
		}
	}
	return nil
}
