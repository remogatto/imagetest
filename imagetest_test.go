package imagetest

import (
	"image"
	"image/png"
	"log"
	"os"
	"testing"

	pt "github.com/remogatto/prettytest"
)

type testSuite struct{ pt.Suite }

var cmp []float64

func loadImage(filename string) image.Image {
	// Open the file.
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Decode the image.
	img, err := png.Decode(file)
	if err != nil {
		log.Fatal(err)
	}
	return img
}

func (t *testSuite) TestCompare() {
	img1 := loadImage("testdata/img1.png")
	img2 := loadImage("testdata/img2.png")
	img4 := loadImage("testdata/img4.png")
	img5 := loadImage("testdata/img5.png")

	cmp = make([]float64, 0)
	tolerance := 0.08

	// compare two identical images
	compare := CompareDistance(img1, img1, Scale)
	t.True(compare == 0)
	cmp = append(cmp, compare)

	// compare two slightly different images (between the
	// tolerance interval)
	compare = CompareDistance(img1, img2, Scale)
	t.True(compare < tolerance)
	cmp = append(cmp, compare)

	// compare two different images (outside the tolerance
	// interval)
	compare = CompareDistance(img1, img4, Scale)
	t.True(compare > tolerance)
	cmp = append(cmp, compare)

	// compare two very different images (outside the tolerance
	// interval)
	compare = CompareDistance(img1, img5, Scale)
	t.True(compare > tolerance)
	cmp = append(cmp, compare)

}

func (t *testSuite) TestAdjustScale() {
	img1 := loadImage("testdata/img1.png")
	img3 := loadImage("testdata/img3.png")
	img6 := loadImage("testdata/img6.png")

	tolerance := 0.08

	// compare two identical images
	compare := CompareDistance(img1, img1, Scale)
	t.True(compare == 0)

	// compare two slightly different images (between the
	// tolerance interval but with different sizes)
	compare = CompareDistance(img1, img3, Scale)
	t.True(compare < tolerance)
	cmp = append(cmp, compare)

	// compare two identical images with different sizes
	compare = CompareDistance(img1, img6, Scale)
	t.True(compare < tolerance)
	cmp = append(cmp, compare)
}

func (t *testSuite) TestAdjustCenter() {
	img1 := loadImage("testdata/box_320_480.png")
	img2 := loadImage("testdata/box_160_240.png")
	img3 := loadImage("testdata/box_160_240_lime.png")

	// compare two identical images
	compare := CompareDistance(img1, img1, Center)
	t.True(compare == 0)

	// compare two images with different size containing the same
	// picture
	compare = CompareDistance(img1, img2, Center)
	t.True(compare == 0)

	// compare two images with same size containing the same
	// picture with different colors
	compare = CompareDistance(img2, img3, Center)
	t.True(compare > 0)

}

func TestCompare(t *testing.T) {
	pt.Run(t, new(testSuite))
}
