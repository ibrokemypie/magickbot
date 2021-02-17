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
	MAGIK   = "magik"
)

func RunMagick(command MagickCommand, files []string, iterations int) error {
	imagick.Initialize()
	defer imagick.Terminate()

	mw := imagick.NewMagickWand()
	defer mw.Destroy()

	maxPixels := uint(viper.GetInt("max_pixels"))
	maxSide := uint(math.Sqrt(float64(maxPixels)))

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
			mw.ResizeImage(maxSide, maxSide, imagick.FILTER_UNDEFINED)
			if err != nil {
				return (err)
			}

			height = mw.GetImageHeight()
			width = mw.GetImageWidth()
		}

		//  Run the magick operation iterations number of times
		for i := 0; i < iterations; i++ {
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
			case MAGIK:
				{
					err = mw.LiquidRescaleImage(uint(width/2), uint(height/2), 1, 0)
					if err != nil {
						return (err)
					}
					err = mw.LiquidRescaleImage(uint(float32(width)*1.5), uint(float32(height)*1.5), 2, 0)
					if err != nil {
						return (err)
					}
				}
			default:
				{
					return errors.New("Unsupported command")
				}
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
