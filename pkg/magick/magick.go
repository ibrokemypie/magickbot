package magick

import (
	"fmt"
	"path"

	"gopkg.in/gographics/imagick.v3/imagick"
)

func Implode(files []string, iterations int) {
	imagick.Initialize()
	defer imagick.Terminate()

	mw := imagick.NewMagickWand()
	defer mw.Destroy()

	for i := 0; i < iterations; i++ {
		for k, file := range files {
			err := mw.ReadImage(file)
			if err != nil {
				panic(err)
			}

			err = mw.ImplodeImage(.5, imagick.INTERPOLATE_PIXEL_UNDEFINED)
			if err != nil {
				panic(err)
			}

			outputFile := file
			if i == 0 {
				outputFile = "/tmp/out." + path.Base(file)
				files[k] = outputFile
			}

			err = mw.WriteImage(outputFile)
			if err != nil {
				panic(err)
			}

			fmt.Println("imploded image " + outputFile)
		}
	}
}

func Explode(files []string, iterations int) {
	imagick.Initialize()
	defer imagick.Terminate()

	mw := imagick.NewMagickWand()
	defer mw.Destroy()

	for i := 0; i < iterations; i++ {
		for k, file := range files {
			err := mw.ReadImage(file)
			if err != nil {
				panic(err)
			}

			err = mw.ImplodeImage(-.5, imagick.INTERPOLATE_PIXEL_UNDEFINED)
			if err != nil {
				panic(err)
			}

			outputFile := file
			if i == 0 {
				outputFile = "/tmp/out." + path.Base(file)
				files[k] = outputFile
			}

			err = mw.WriteImage(outputFile)
			if err != nil {
				panic(err)
			}

			fmt.Println("exploded image " + outputFile)
		}
	}
}
