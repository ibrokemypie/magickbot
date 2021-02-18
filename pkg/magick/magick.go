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
	"moarjpeg",
	"deepfry",
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
					argument = 5
				} else if argument > maxIterations {
					argument = maxIterations
				}

				magick(mw, float64(width), float64(height), float64(argument))
			}
		case "moarjpeg":
			{
				if argument < 1 {
					argument = 1
				} else if argument > maxIterations {
					argument = maxIterations
				}

				jpegify(mw, argument, file)

			}
		case "deepfry":
			{
				err = deepfry(mw, width, height, file)
				if err != nil {
					return -1, err
				}

				argument = -1
			}
		default:
			{
				return -1, errors.New("Unsupported command, try \"help\"")
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

func jpegify(mw *imagick.MagickWand, iterations int, file string) error {
	file = strings.TrimSuffix(file, filepath.Ext(file)) + ".jpg"

	for i := 0; i < iterations; i++ {
		if i > 1 {
			err := mw.WriteImage(file)
			if err != nil {
				return err
			}

			err = mw.ReadImage(file)
			if err != nil {
				return err
			}
		}

		err := mw.SetImageAlphaChannel(imagick.ALPHA_CHANNEL_OPAQUE)
		if err != nil {
			return err
		}

		err = mw.SetImageInterlaceScheme(imagick.INTERLACE_JPEG)
		if err != nil {
			return err
		}

		err = mw.SetImageCompression(imagick.COMPRESSION_JPEG)
		if err != nil {
			return err
		}

		err = mw.SetImageCompressionQuality(15)
		if err != nil {
			return err
		}

		err = mw.SharpenImage(0, 4)
		if err != nil {
			return err
		}
	}

	return nil
}

func magick(mw *imagick.MagickWand, width, height, scale float64) error {
	scaleOne := float64(1)
	scaleTwo := float64(2)
	if scale != 0 {
		scaleOne = float64(scale) / 2
		scaleTwo = float64(scale)
	}

	err := mw.LiquidRescaleImage(uint(width/2), uint(height/2), scaleOne, 0)
	if err != nil {
		return err
	}

	err = mw.LiquidRescaleImage(uint(width*1.5), uint(height*1.5), scaleTwo, 0)
	if err != nil {
		return err
	}

	return nil
}

func deepfry(mw *imagick.MagickWand, width, height uint, file string) error {
	orangeOverlay := imagick.NewPixelWand()
	orangeOverlay.SetColor("#992604")
	orangeOverlay.SetAlpha(65)

	orangeOverlayImage := imagick.NewMagickWand()
	err := orangeOverlayImage.NewImage(width, height, orangeOverlay)
	if err != nil {
		return err
	}

	err = mw.CompositeImage(orangeOverlayImage, imagick.COMPOSITE_OP_OVERLAY, true, 0, 0)
	if err != nil {
		return err
	}

	err = mw.ContrastImage(true)
	if err != nil {
		return err
	}

	err = mw.ModulateImage(100, 150, 100)
	if err != nil {
		return err
	}

	err = jpegify(mw, 3, file)
	if err != nil {
		return err
	}

	return nil
}
