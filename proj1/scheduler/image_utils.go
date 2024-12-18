package scheduler

import (
	"fmt"
	"image"
	"image/color"
	"proj1/png"
)

func PadImage(img *png.Image, padding int) *png.Image {
    bounds := img.Bounds
    in := png.GetInImage(img)
    paddedBounds := image.Rect(0, 0, bounds.Dx()+2*padding, bounds.Dy()+2*padding)
    paddedImage := image.NewRGBA64(paddedBounds)

    for y := paddedBounds.Min.Y; y < paddedBounds.Max.Y; y++ {
        for x := paddedBounds.Min.X; x < paddedBounds.Max.X; x++ {
            paddedImage.Set(x, y, color.RGBA64{0, 0, 0, 0}) 
        }
    }

    for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
        for x := bounds.Min.X; x < bounds.Max.X; x++ {
            paddedImage.Set(x+padding, y+padding, in.At(x, y))
        }
    }

    newImg := &png.Image{
        Bounds: paddedBounds,
    }
    png.SetInImage(newImg, paddedImage)
    png.SetOutImage(newImg, image.NewRGBA64(paddedBounds))
    return newImg
}

func CopyFromPaddedImage(img *png.Image, paddedImage *png.Image, padding int) {
    inPadded := png.GetInImage(paddedImage)
    inImg := png.GetInImage(img)

    imgBounds := img.Bounds

    for y := imgBounds.Min.Y; y < imgBounds.Max.Y; y++ {
        for x := imgBounds.Min.X; x < imgBounds.Max.X; x++ {
            color := inPadded.At(x+padding, y+padding)
            inImg.Set(x, y, color)
        }
    }
}


func ApplyEffects(inputPath string, outputPath string, effects []string) error { 
	img, err := png.Load(inputPath)
	if err != nil {
		return err
	}
	for _, effect := range effects {
		switch effect {
			case "B", "E", "S":
				kernel := getKernel(effect)
				ApplyConvolution(img, kernel)
			case "G":
            	ApplyGrayscale(img)
			default:
				fmt.Printf("Unknown effect: %s\n", effect)
			}
	}
	return img.Save(outputPath)
	
}


func ApplyGrayscale(img *png.Image) {
    bounds := img.Bounds
    out := png.GetOutImage(img)
    for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
        for x := bounds.Min.X; x < bounds.Max.X; x++ {
            r, g, b, _ := png.GetInImage(img).At(x, y).RGBA()
            greyC := uint16((r + g + b) / 3)
            out.Set(x, y, color.RGBA64{greyC, greyC, greyC, 65535})
        }
    }
    
    png.SetInImage(img, out)
    png.SetOutImage(img, out)    
}


