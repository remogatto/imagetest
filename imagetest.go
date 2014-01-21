package imagetest

import (
	"image"
	"image/color"
	"image/draw"

	"github.com/nfnt/resize"
)

var (
	// A default scaler.
	Scale = Scaler{}

	// A default centerer with black fill color.
	Center = Centerer{color.Black}
)

type Adjuster interface {
	// Adjust adjusts the images before the comparaison.
	Adjust(img1, img2 image.Image) (image.Image, image.Image)
}

// Scaler is an Adjuster that resize the two images to the same size
// of the smaller one.
type Scaler struct{}

// Adjust adjusts the size of the two images to fit the size of the
// smaller one.
func (s Scaler) Adjust(img1, img2 image.Image) (image.Image, image.Image) {
	bounds := img1.Bounds().Intersect(img2.Bounds())
	width, height := uint(bounds.Size().X), uint(bounds.Size().Y)
	return resize.Resize(width, height, img1, resize.Bilinear), resize.Resize(width, height, img2, resize.Bilinear)
}

// Centerer is an Adjuster that aligns the centers's images and fill
// the background of smaller image with a color. Default fill color is
// black.
type Centerer struct {
	FillColor color.Color
}

// Adjust adjusts the two images aligning ther centers. The background
// of the smaller image is filled with FillColor. The returned images
// are of the same size.
func (c Centerer) Adjust(img1, img2 image.Image) (image.Image, image.Image) {
	img1Rect := img1.Bounds()
	img2Rect := img2.Bounds()
	backgroundRect := img1Rect.Union(img2Rect)

	dstImg1 := image.NewRGBA(backgroundRect)
	dstImg2 := image.NewRGBA(backgroundRect)

	// Fill destination images with FillColor
	draw.Draw(dstImg1, dstImg1.Bounds(), &image.Uniform{c.FillColor}, image.ZP, draw.Src)
	draw.Draw(dstImg2, dstImg2.Bounds(), &image.Uniform{c.FillColor}, image.ZP, draw.Src)

	// Copy img1 to the center of dstImg1
	dp := image.Point{
		(backgroundRect.Max.X-backgroundRect.Min.X)/2 - (img1Rect.Max.X-img1Rect.Min.X)/2,
		(backgroundRect.Max.Y-backgroundRect.Min.Y)/2 - (img1Rect.Max.Y-img1Rect.Min.Y)/2,
	}
	r := image.Rectangle{dp, dp.Add(img1Rect.Size())}
	draw.Draw(dstImg1, r, img1, image.ZP, draw.Src)

	// Copy img2 to the center of dstImg2
	dp = image.Point{
		(backgroundRect.Max.X-backgroundRect.Min.X)/2 - (img2Rect.Max.X-img2Rect.Min.X)/2,
		(backgroundRect.Max.Y-backgroundRect.Min.Y)/2 - (img2Rect.Max.Y-img2Rect.Min.Y)/2,
	}
	r = image.Rectangle{dp, dp.Add(img2Rect.Size())}
	draw.Draw(dstImg2, r, img2, image.ZP, draw.Src)

	return dstImg1, dstImg2
}

// Compare compares two images with a naive "distance" algorithm. It
// returns a distance percentage value from 0 (identical images) to 1
// (totally different images). Images are adjusted by Adjuster before
// comparaison.
func CompareDistance(img1 image.Image, img2 image.Image, adjust Adjuster) (distance float64) {
	imgCmp1, imgCmp2 := adjust.Adjust(img1, img2)
	bounds := imgCmp1.Bounds()
	sum := 0.0
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r1, g1, b1, a1 := imgCmp1.At(x, y).RGBA()
			r2, g2, b2, a2 := imgCmp2.At(x, y).RGBA()

			// normalize
			dr := float64(byte(r1))/255 - float64(byte(r2))/255
			dg := float64(byte(g1))/255 - float64(byte(g2))/255
			db := float64(byte(b1))/255 - float64(byte(b2))/255
			da := float64(byte(a1))/255 - float64(byte(a2))/255

			sum += dr*dr + dg*dg + db*db + da*da
		}
	}
	return sum / float64(4*bounds.Dx()*bounds.Dy())
}
