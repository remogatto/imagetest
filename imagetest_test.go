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
	img3 := loadImage("testdata/img3.png")
	img4 := loadImage("testdata/img4.png")
	img5 := loadImage("testdata/img5.png")
	img6 := loadImage("testdata/img6.png")

	cmp = make([]float64, 0)
	tolerance := 8.0

	// compare two identical images
	compare := CompareDistance(img1, img1)
	t.True(compare == 0)
	cmp = append(cmp, compare)

	// compare two slightly different images (between the
	// tolerance interval)
	compare = CompareDistance(img1, img2)
	t.True(compare < tolerance)
	cmp = append(cmp, compare)

	// compare two slightly different images (between the
	// tolerance interval but with different sizes)
	compare = CompareDistance(img1, img3)
	t.True(compare < tolerance)
	cmp = append(cmp, compare)

	// compare two different images (outside the tolerance
	// interval)
	compare = CompareDistance(img1, img4)
	t.True(compare > tolerance)
	cmp = append(cmp, compare)

	// compare two very different images (outside the tolerance
	// interval)
	compare = CompareDistance(img1, img5)
	t.True(compare > tolerance)
	cmp = append(cmp, compare)

	// compare two identical images with different sizes
	compare = CompareDistance(img1, img6)
	t.True(compare < tolerance)
	cmp = append(cmp, compare)
}

func TestCompare(t *testing.T) {
	pt.Run(t, new(testSuite))
}
