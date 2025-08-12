package resize

import (
	"image"
	"math/bits"
	"sync"
)

const bufferSize = 4 * 2048 * 2048 * 4 // 64 MiB

var pool = sync.Pool{
	New: func() any {
		return make([]uint8, 0, bufferSize)
	},
}

func GetPixelBuffer(size int) ([]uint8, func()) {
	if size <= bufferSize {
		pix := pool.Get().([]uint8)[:bufferSize]
		return pix, func() {
			pool.Put(pix[:0])
		}
	}
	pix := make([]uint8, size)
	return pix, func() {}
}

func NewRGBA(r image.Rectangle) (*image.RGBA, func()) {
	size := pixelBufferLength(4, r, "RGBA")
	pix, canc := GetPixelBuffer(size)
	return &image.RGBA{
		Pix:    pix,
		Stride: 4 * r.Dx(),
		Rect:   r,
	}, canc
}

// NewRGBA64 returns a new [RGBA64] image with the given bounds.
func NewRGBA64(r image.Rectangle) (*image.RGBA64, func()) {
	size := pixelBufferLength(8, r, "RGBA64")
	pix, canc := GetPixelBuffer(size)
	return &image.RGBA64{
		Pix:    pix,
		Stride: 8 * r.Dx(),
		Rect:   r,
	}, canc
}

// NewGray returns a new [Gray] image with the given bounds.
func NewGray(r image.Rectangle) (*image.Gray, func()) {
	size := pixelBufferLength(1, r, "Gray")
	pix, canc := GetPixelBuffer(size)
	return &image.Gray{
		Pix:    pix,
		Stride: 1 * r.Dx(),
		Rect:   r,
	}, canc
}

// NewGray16 returns a new [Gray16] image with the given bounds.
func NewGray16(r image.Rectangle) (*image.Gray16, func()) {
	size := pixelBufferLength(2, r, "Gray16")
	pix, canc := GetPixelBuffer(size)
	return &image.Gray16{
		Pix:    pix,
		Stride: 2 * r.Dx(),
		Rect:   r,
	}, canc
}

// NewNRGBA returns a new [NRGBA] image with the given bounds.
func NewNRGBA(r image.Rectangle) (*image.NRGBA, func()) {
	size := pixelBufferLength(4, r, "NRGBA")
	pix, canc := GetPixelBuffer(size)
	return &image.NRGBA{
		Pix:    pix,
		Stride: 4 * r.Dx(),
		Rect:   r,
	}, canc
}

// NewNRGBA64 returns a new [NRGBA64] image with the given bounds.
func NewNRGBA64(r image.Rectangle) (*image.NRGBA64, func()) {
	size := pixelBufferLength(8, r, "NRGBA64")
	pix, canc := GetPixelBuffer(size)
	return &image.NRGBA64{
		Pix:    pix,
		Stride: 8 * r.Dx(),
		Rect:   r,
	}, canc
}

// mul3NonNeg returns (x * y * z), unless at least one argument is negative or
// if the computation overflows the int type, in which case it returns -1.
func mul3NonNeg(x int, y int, z int) int {
	if (x < 0) || (y < 0) || (z < 0) {
		return -1
	}
	hi, lo := bits.Mul64(uint64(x), uint64(y))
	if hi != 0 {
		return -1
	}
	hi, lo = bits.Mul64(lo, uint64(z))
	if hi != 0 {
		return -1
	}
	a := int(lo)
	if (a < 0) || (uint64(a) != lo) {
		return -1
	}
	return a
}

// pixelBufferLength returns the length of the []uint8 typed Pix slice field
// for the NewXxx functions. Conceptually, this is just (bpp * width * height),
// but this function panics if at least one of those is negative or if the
// computation would overflow the int type.
//
// This panics instead of returning an error because of backwards
// compatibility. The NewXxx functions do not return an error.
func pixelBufferLength(bytesPerPixel int, r image.Rectangle, imageTypeName string) int {
	totalLength := mul3NonNeg(bytesPerPixel, r.Dx(), r.Dy())
	if totalLength < 0 {
		panic("image: New" + imageTypeName + " Rectangle has huge or negative dimensions")
	}
	return totalLength
}
