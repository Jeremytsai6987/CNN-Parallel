package scheduler

import (
	"image/color"
	"math"
	"proj1/png"
)


func getKernel(effect string) [9]float64 {
    switch effect {
    case "B": // Blur effect
        return [9]float64{
    1.0 / 9.0, 1.0 / 9.0, 1.0 / 9.0,
    1.0 / 9.0, 1.0 / 9.0, 1.0 / 9.0,
    1.0 / 9.0, 1.0 / 9.0, 1.0 / 9.0,
}
    case "E": // Edge detection effect
        return [9]float64{
   -1, -1, -1,
   -1,  8, -1,
   -1, -1, -1,
}
    case "S": // Sharpen effect
        return [9]float64{
    0, -1,  0,
   -1,  5, -1,
    0, -1,  0,
}
    default:
        return [9]float64{}
    }
}

func FlipKernel(kernel [9]float64) [9]float64 {
    return [9]float64{
        kernel[8], kernel[7], kernel[6], 
        kernel[5], kernel[4], kernel[3], 
        kernel[2], kernel[1], kernel[0], 
    }
}

func clamp(value float64) uint16 {
    return uint16(math.Max(0, math.Min(65535, value)))
}

func ApplyConvolution(img *png.Image, kernel [9]float64) {
    flippedKernel := FlipKernel(kernel) 

    bounds := img.Bounds
    paddedImg := PadImage(img, 1)  
    in := png.GetInImage(paddedImg)  
	out := png.GetOutImage(img)

    padding := 1 

    for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
        for x := bounds.Min.X; x < bounds.Max.X; x++ {
            var rSum, gSum, bSum float64

            k := 0
            for ky := -1; ky <= 1; ky++ {
                for kx := -1; kx <= 1; kx++ {
                    r, g, b, _ := in.At(x+kx+padding, y+ky+padding).RGBA()

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

    png.SetInImage(img, out)  
    png.SetOutImage(img, out)
}

