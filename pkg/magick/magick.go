package magick

import (
	"errors"
	"math"
	"path"

	"github.com/spf13/viper"
	"gopkg.in/gographics/imagick.v3/imagick"
)

type MagickCommand string

const (
	EXPLODE = "explode"
	IMPLODE = "implode"
	MAGIK   = "magick"
)

func RunMagick(command MagickCommand, files []string, argument int) error {
	imagick.Initialize()
	defer imagick.Terminate()

	mw := imagick.NewMagickWand()
	defer mw.Destroy()

	maxPixels := uint(viper.GetInt("max_pixels"))
	maxIterations := viper.GetInt("max_iterations")

	// run specified iterations of operations  for each file
	for k, file := range files {
		err := mw.ReadImage(file)
		if err != nil {
			return (err)
		}

		height := mw.GetImageHeight()
		width := mw.GetImageWidth()

		// If the image has more than maxPixels pixels, resize it down to fit. This is to reduce the maximum utilisation from a single operation.
		if width*height > maxPixels {

			ratio := float64(width) / float64(height)
			newHeight := math.Sqrt(float64(maxPixels) / ratio)
			newWidth := float64(maxPixels) / newHeight

			mw.ResizeImage(uint(newWidth), uint(newWidth), imagick.FILTER_UNDEFINED)
			if err != nil {
				return (err)
			}

			height = mw.GetImageHeight()
			width = mw.GetImageWidth()
		}

		switch command {
		case EXPLODE:
			{
				if argument < 1 {
					argument = 1
				} else if argument > maxIterations {
					argument = maxIterations
				}

				for i := 0; i < argument; i++ {
					err = mw.ImplodeImage(-.5, imagick.INTERPOLATE_PIXEL_UNDEFINED)
					if err != nil {
						return (err)
					}
				}
			}
		case IMPLODE:
			{
				if argument < 1 {
					argument = 1
				} else if argument > maxIterations {
					argument = maxIterations
				}

				for i := 0; i < argument; i++ {
					err = mw.ImplodeImage(.5, imagick.INTERPOLATE_PIXEL_UNDEFINED)
					if err != nil {
						return (err)
					}
				}
			}
		case MAGIK:
			{
				if argument < 1 {
					argument = 1
				} else if argument > maxIterations {
					argument = maxIterations
				}

				err = mw.LiquidRescaleImage(uint(width/2), uint(height/2), 1*float64(argument), 0)
				if err != nil {
					return (err)
				}
				err = mw.LiquidRescaleImage(uint(float32(width)*1.5), uint(float32(height)*1.5), 2*float64(argument), 0)
				if err != nil {
					return (err)
				}
			}
		default:
			{
				return errors.New("Unsupported command")
			}
		}

		outputFile := "/tmp/out." + path.Base(file)
		files[k] = outputFile

		err = mw.WriteImage(outputFile)
		if err != nil {
			return (err)
		}
	}

	return nil
}
