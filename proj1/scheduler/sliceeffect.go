package scheduler

import (
	"fmt"
	"image/color"
	"proj1/png"
	"sync"
)

func ApplyEffectsToSlice(inputPath string, outputPath string, effects []string, numSlices int) error {    
    img, err := png.Load(inputPath) 
    if err != nil {
        return fmt.Errorf("failed to load image: %v", err)
    }

    // Pad the image to handle border effects during convolution
    padding := 1
    paddedImage := PadImage(img, padding)
    slices := SliceImage(paddedImage.Bounds, numSlices, padding)
    
    for _, effect := range effects {
        
        var wg sync.WaitGroup
        wg.Add(len(slices))

        for _, slice := range slices {
            go func(slice Slice) {
                defer wg.Done()
                switch effect {
                case "B", "E", "S":
                    kernel := getKernel(effect)
                    ApplyConvolutionToSlice(paddedImage, kernel, slice)
                case "G":
                    ApplyGrayscaleToSlice(paddedImage, slice, padding)
                default:
                    fmt.Printf("Unknown effect: %s\n", effect)
                }
            }(slice)
        }

        wg.Wait()
        png.SwapInOutImage(paddedImage)
    }

    CopyFromPaddedImage(img, paddedImage, padding)
	png.SwapInOutImage(img)
    err = img.Save(outputPath) // save based on out
    if err != nil {
        return fmt.Errorf("failed to save image: %v", err)
    }
    return nil
}

func ApplyGrayscaleToSlice(img *png.Image, slice Slice, padding int) {
    bounds := img.Bounds
    in := png.GetInImage(img)
    out := png.GetOutImage(img)

    if in == nil || out == nil {
        fmt.Printf("Error: nil image buffer in grayscale operation\n")
        return
    }

    for y := slice.WriteStartY; y < slice.WriteEndY; y++ {
        for x := bounds.Min.X; x < bounds.Max.X; x++ {
            pixel := in.At(x+padding, y+padding)
            r, g, b, a := pixel.RGBA()
            greyC := uint16((r + g + b) / 3)
            
            out.Set(x+padding, y+padding, color.RGBA64{greyC, greyC, greyC, uint16(a)})
        }
    }

}

func ApplyConvolutionToSlice(img *png.Image, kernel [9]float64, slice Slice) {
    bounds := img.Bounds
    in := png.GetInImage(img)
    out := png.GetOutImage(img)
    flippedKernel := FlipKernel(kernel)

    for y := slice.WriteStartY; y < slice.WriteEndY; y++ {
        for x := bounds.Min.X + 1; x < bounds.Max.X - 1; x++ {
            var rSum, gSum, bSum float64
            k := 0

            for ky := -1; ky <= 1; ky++ {
                for kx := -1; kx <= 1; kx++ {
                    r, g, b, _ := in.At(x+kx, y+ky).RGBA()
                    rSum += float64(r) * flippedKernel[k]
                    gSum += float64(g) * flippedKernel[k]
                    bSum += float64(b) * flippedKernel[k]
                    k++
                }
            }

            out.Set(x, y, color.RGBA64{
                R: clamp(rSum),
                G: clamp(gSum),
                B: clamp(bSum),
                A: 65535,
            })
        }
    }
}


