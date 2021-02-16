package magick

import (
	"fmt"
	"path"

	"gopkg.in/gographics/imagick.v3/imagick"
)

func Implode(files []string) {
	imagick.Initialize()
	defer imagick.Terminate()

	mw := imagick.NewMagickWand()
	defer mw.Destroy()

	for _, file := range files {
		err := mw.ReadImage(file)
		if err != nil {
			panic(err)
		}

		err = mw.ImplodeImage(.5, imagick.INTERPOLATE_PIXEL_UNDEFINED)
		if err != nil {
			panic(err)
		}

		outputFile := "/tmp/out." + path.Base(file)

		err = mw.WriteImage(outputFile)
		if err != nil {
			panic(err)
		}

		fmt.Println("imploded image " + outputFile)
	}
}

func Explode(files []string) {
	imagick.Initialize()
	defer imagick.Terminate()

	mw := imagick.NewMagickWand()
	defer mw.Destroy()

	for _, file := range files {
		err := mw.ReadImage(file)
		if err != nil {
			panic(err)
		}

		err = mw.ImplodeImage(-.5, imagick.INTERPOLATE_PIXEL_UNDEFINED)
		if err != nil {
			panic(err)
		}

		outputFile := "/tmp/out." + path.Base(file)

		err = mw.WriteImage(outputFile)
		if err != nil {
			panic(err)
		}

		fmt.Println("exploded image " + outputFile)
	}
}
