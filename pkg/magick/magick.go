package magick

import (
	"errors"
	"math"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/spf13/viper"
	"gopkg.in/gographics/imagick.v3/imagick"
)

var MagickCommands = []string{
	"explode",
	"implode",
	"magick",
	"compress",
}

func RunMagick(command string, files []string, argument int) (int, error) {
	imagick.Initialize()
	defer imagick.Terminate()

	mw := imagick.NewMagickWand()
	defer mw.Destroy()

	maxInputPixels := viper.GetUint("max_pixels_in")
	maxOutputPixels := viper.GetUint("max_pixels_out")
	maxIterations := viper.GetInt("max_iterations")

	// run specified iterations of operations  for each file
	for k, file := range files {
		err := mw.ReadImage(file)
		if err != nil {
			return -1, err
		}

		height := mw.GetImageHeight()
		width := mw.GetImageWidth()

		if height*width > maxInputPixels {
			return -1, errors.New("Input too large! Maximum pixels: " + strconv.Itoa(int(maxInputPixels)))
		}

		// If the image has more than maxPixels pixels, resize it down to fit. This is to reduce the maximum utilisation from a single operation.
		if width*height > maxOutputPixels {

			ratio := float64(width) / float64(height)
			newHeight := math.Sqrt(float64(maxOutputPixels) / ratio)
			newWidth := float64(maxOutputPixels) / newHeight

			mw.ResizeImage(uint(newWidth), uint(newWidth), imagick.FILTER_UNDEFINED)
			if err != nil {
				return -1, err
			}

			height = mw.GetImageHeight()
			width = mw.GetImageWidth()
		}

		switch command {
		case "explode":
			{
				if argument < 1 {
					argument = 1
				} else if argument > maxIterations {
					argument = maxIterations
				}

				for i := 0; i < argument; i++ {
					err = mw.ImplodeImage(-.5, imagick.INTERPOLATE_PIXEL_UNDEFINED)
					if err != nil {
						return -1, err
					}
				}
			}
		case "implode":
			{
				if argument < 1 {
					argument = 1
				} else if argument > maxIterations {
					argument = maxIterations
				}

				for i := 0; i < argument; i++ {
					err = mw.ImplodeImage(.5, imagick.INTERPOLATE_PIXEL_UNDEFINED)
					if err != nil {
						return -1, err
					}
				}
			}
		case "magick":
			{
				if argument < 1 {
					argument = 1
				} else if argument > maxIterations {
					argument = maxIterations
				}

				err = mw.LiquidRescaleImage(uint(width/2), uint(height/2), 1*float64(argument), 0)
				if err != nil {
					return -1, err
				}
				err = mw.LiquidRescaleImage(uint(float32(width)*1.5), uint(float32(height)*1.5), 2*float64(argument), 0)
				if err != nil {
					return -1, err
				}
			}
		case "compress":
			{
				if argument < 1 {
					argument = 1
				} else if argument > maxIterations {
					argument = maxIterations
				}

				file = strings.TrimSuffix(file, filepath.Ext(file)) + ".jpg"

				for i := 0; i < argument; i++ {
					if i > 1 {
						err = mw.WriteImage(file)
						if err != nil {
							return -1, err
						}

						err := mw.ReadImage(file)
						if err != nil {
							return -1, err
						}
					}

					err = mw.SetImageAlphaChannel(imagick.ALPHA_CHANNEL_OPAQUE)
					if err != nil {
						return -1, err
					}

					err = mw.SetImageInterlaceScheme(imagick.INTERLACE_JPEG)
					if err != nil {
						return -1, err
					}

					err = mw.SetImageCompression(imagick.COMPRESSION_JPEG)
					if err != nil {
						return -1, err
					}

					err = mw.SetImageCompressionQuality(15)
					if err != nil {
						return -1, err
					}

					err = mw.SharpenImage(0, 4)
					if err != nil {
						return -1, err
					}

				}
			}
		default:
			{
				return -1, errors.New("Unsupported command")
			}
		}

		outputFile := "/tmp/out." + path.Base(file)
		files[k] = outputFile

		err = mw.WriteImage(outputFile)
		if err != nil {
			return -1, err
		}
	}

	return argument, nil
}
