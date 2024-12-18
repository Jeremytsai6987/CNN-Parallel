package scheduler

import (
	"image"
)

type Slice struct {
	ReadStartY, ReadEndY, WriteStartY, WriteEndY int
}

func SliceImage(bounds image.Rectangle, numSlices int, overlap int) []Slice {
    sliceHeight := (bounds.Dy() + numSlices - 1) / numSlices // 向上取整
    slices := make([]Slice, numSlices)
    for i := 0; i < numSlices; i++ {
        startY := bounds.Min.Y + i*sliceHeight
        endY := startY + sliceHeight
        if endY > bounds.Max.Y {
            endY = bounds.Max.Y
        }
        readStartY := startY
        readEndY := endY
        if i > 0 {
            readStartY -= overlap
        }
        if i < numSlices-1 {
            readEndY += overlap
        }
        slices[i] = Slice{
            ReadStartY:  readStartY,
            ReadEndY:    readEndY,
            WriteStartY: startY,
            WriteEndY:   endY,
        }
    }
    return slices
}






