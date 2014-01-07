package imagetest

import (
	"image"

	"github.com/nfnt/resize"
)

// Compare compares two images with a naive "distance" algorithm. It
// returns a distance percentage value from 0 (identical images) to
// 100 (totally different images). Images of different sizes are
// automatically scaled before compare.
func CompareDistance(img1, img2 image.Image) (distance float64) {
	bounds := img1.(image.Image).Bounds().Intersect(img2.(image.Image).Bounds())
	width, height := uint(bounds.Size().X), uint(bounds.Size().Y)
	imgCmp1 := resize.Resize(width, height, img1.(image.Image), resize.Bilinear)
	imgCmp2 := resize.Resize(width, height, img2.(image.Image), resize.Bilinear)
	avg := 0.0
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r1, g1, b1, a1 := imgCmp1.At(x, y).RGBA()
			r2, g2, b2, a2 := imgCmp2.At(x, y).RGBA()

			// normalize
			dr := float64(byte(r1))/255 - float64(byte(r2))/255
			dg := float64(byte(g1))/255 - float64(byte(g2))/255
			db := float64(byte(b1))/255 - float64(byte(b2))/255
			da := float64(byte(a1))/255 - float64(byte(a2))/255
			avg += dr*dr + dg*dg + db*db + da*da
		}
	}
	return avg / float64(4*width*height) * 100
}
