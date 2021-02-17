package magick

import (
	"path"

	"gopkg.in/gographics/imagick.v3/imagick"
)

func Implode(files []string, iterations int) error {
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

			err = mw.ImplodeImage(.5, imagick.INTERPOLATE_PIXEL_UNDEFINED)
			if err != nil {
				return (err)
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

func Explode(files []string, iterations int) error {
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

			err = mw.ImplodeImage(-.5, imagick.INTERPOLATE_PIXEL_UNDEFINED)
			if err != nil {
				return (err)
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
