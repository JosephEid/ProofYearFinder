package main

import (
	"image"
	"image/color"
	"image/draw"
	"log"
	"os"

	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/llgcode/draw2d/draw2dkit"
)

// func DrawBoundingBoxes(top, left, width, height float64) {
func main() {
	// width := 200
	// height := 200

	existingImageFile, err := os.Open("eaafd6c8-a586-43c7-9d1f-e4d4c177905b.png")
	if err != nil {
		log.Fatal("Error opening file: ", err)
	}
	defer existingImageFile.Close()

	myimage, _, err := image.Decode(existingImageFile)
	if err != nil {
		log.Fatal("Error decoding image: ", err)
	}

	img := imageToRGBA(myimage)

	gc := draw2dimg.NewGraphicContext(img)

	gc.SetStrokeColor(color.NRGBA{255, 255, 255, 255})
	gc.SetFillColor(color.NRGBA{255, 255, 255, 255})

	gc.SetStrokeColor(color.NRGBA{255, 0, 0, 255})
	gc.SetLineWidth(1)

	// Draw a circle
	draw2dkit.Circle(gc, 100, 100, 50)
	gc.Stroke()

	draw2dimg.SaveToPngFile("TestCircle.png", img)

}

func DrawSomething() {
	existingImageFile, err := os.Open("eaafd6c8-a586-43c7-9d1f-e4d4c177905b.png")
	if err != nil {
		log.Fatal("Error opening file: ", err)
	}
	defer existingImageFile.Close()

	// Calling the generic image.Decode() will tell give us the data
	// and type of image it is as a string. We expect "png"
	// imConfig, _, err := image.DecodeConfig(existingImageFile)
	// if err != nil {
	// 	log.Fatal()
	// }

	// newOffset, err := existingImageFile.Seek(0, 0)
	// if err != nil {
	// 	log.Fatal("Error resetting offset: ", err)
	// }
	// fmt.Println("file reset offset: ", newOffset)

	myimage, _, err := image.Decode(existingImageFile)
	if err != nil {
		log.Fatal("Error decoding image: ", err)
	}

	// fmt.Println("image height: ", imConfig.Height)

	// Initialize the graphic context on an RGBA image
	// dest := image.NewRGBA(myimage)
	gc := draw2dimg.NewGraphicContext(imageToRGBA(myimage))

	// // Set some properties
	gc.SetFillColor(color.RGBA{0x44, 0xff, 0x44, 0xff})
	gc.SetStrokeColor(color.RGBA{0x44, 0x44, 0x44, 0xff})
	gc.SetLineWidth(5)

	// Draw a closed shape
	gc.BeginPath() // Initialize a new path
	draw2dkit.Rectangle(gc, 10, 10, 30, 30)
	// gc.MoveTo(200, 200) // Move to a position to start the new path
	// CurveRectangle(gc, 200, 200, 100, 50, color.RGBA{0x44, 0xff, 0x44, 0xff}, color.RGBA{0x44, 0x44, 0x44, 0xff})
	defer gc.Close()
	gc.Stroke()

	// Save to file
	draw2dimg.SaveToPngFile("hello.png", myimage)

	// fmt.Println("Top: ", top)
	// fmt.Println("Left: ", left)
	// fmt.Println("Width: ", width)
	// fmt.Println("Height: ", height)

	// new_png_file := "./outputImage.png" // output image will live here

	// existingImageFile, err := os.Open("1711961501985403200.png")
	// if err != nil {
	// 	// Handle error
	// }
	// defer existingImageFile.Close()

	// // Calling the generic image.Decode() will tell give us the data
	// // and type of image it is as a string. We expect "png"
	// myimage, _, err := image.Decode(existingImageFile)
	// if err != nil {
	// 	// Handle error
	// }

	// b := myimage.Bounds()
	// m := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	// // draw.Draw(m, m.Bounds(), myimage, b.Min, draw.Src)

	// // myimage := image.NewRGBA(image.Rect(0, 0, 220, 220)) // x1,y1,  x2,y2 of background rectangle
	// // mygreen := color.RGBA{0, 100, 0, 255} //  R, G, B, Alpha

	// // // backfill entire background surface with color mygreen
	// // draw.Draw(m, myimage.Bounds(), &image.Uniform{mygreen}, image.Point{0, 0}, draw.Src)

	// red_rect := image.Rect(60, 80, 120, 160) //  geometry of 2nd rectangle which we draw atop above rectangle
	// myred := color.RGBA{200, 0, 0, 255}

	// // create a red rectangle atop the green surface
	// draw.Draw(m, red_rect, &image.Uniform{myred}, image.Point{50, 50}, draw.Src)

	// myfile, err := os.Create(new_png_file) // ... now lets save output image
	// if err != nil {
	// 	log.Error("Error creating file: ", err)
	// 	panic(err)
	// }
	// defer myfile.Close()
	// png.Encode(myfile, myimage) // output file /tmp/two_rectangles.png
}

func imageToRGBA(src image.Image) *image.RGBA {

	// No conversion needed if image is an *image.RGBA.
	if dst, ok := src.(*image.RGBA); ok {
		return dst
	}

	// Use the image/draw package to convert to *image.RGBA.
	b := src.Bounds()
	dst := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(dst, dst.Bounds(), src, b.Min, draw.Src)
	return dst
}

func CurveRectangle(gc draw2d.GraphicContext, x0, y0,
	rectWidth, rectHeight float64, stroke, fill color.Color) {
	radius := (rectWidth + rectHeight) / 4

	x1 := x0 + rectWidth
	y1 := y0 + rectHeight
	if rectWidth/2 < radius {
		if rectHeight/2 < radius {
			gc.MoveTo(x0, (y0+y1)/2)
			gc.CubicCurveTo(x0, y0, x0, y0, (x0+x1)/2, y0)
			gc.CubicCurveTo(x1, y0, x1, y0, x1, (y0+y1)/2)
			gc.CubicCurveTo(x1, y1, x1, y1, (x1+x0)/2, y1)
			gc.CubicCurveTo(x0, y1, x0, y1, x0, (y0+y1)/2)
		} else {
			gc.MoveTo(x0, y0+radius)
			gc.CubicCurveTo(x0, y0, x0, y0, (x0+x1)/2, y0)
			gc.CubicCurveTo(x1, y0, x1, y0, x1, y0+radius)
			gc.LineTo(x1, y1-radius)
			gc.CubicCurveTo(x1, y1, x1, y1, (x1+x0)/2, y1)
			gc.CubicCurveTo(x0, y1, x0, y1, x0, y1-radius)
		}
	} else {
		if rectHeight/2 < radius {
			gc.MoveTo(x0, (y0+y1)/2)
			gc.CubicCurveTo(x0, y0, x0, y0, x0+radius, y0)
			gc.LineTo(x1-radius, y0)
			gc.CubicCurveTo(x1, y0, x1, y0, x1, (y0+y1)/2)
			gc.CubicCurveTo(x1, y1, x1, y1, x1-radius, y1)
			gc.LineTo(x0+radius, y1)
			gc.CubicCurveTo(x0, y1, x0, y1, x0, (y0+y1)/2)
		} else {
			gc.MoveTo(x0, y0+radius)
			gc.CubicCurveTo(x0, y0, x0, y0, x0+radius, y0)
			gc.LineTo(x1-radius, y0)
			gc.CubicCurveTo(x1, y0, x1, y0, x1, y0+radius)
			gc.LineTo(x1, y1-radius)
			gc.CubicCurveTo(x1, y1, x1, y1, x1-radius, y1)
			gc.LineTo(x0+radius, y1)
			gc.CubicCurveTo(x0, y1, x0, y1, x0, y1-radius)
		}
	}
	gc.Close()
	gc.SetStrokeColor(stroke)
	gc.SetFillColor(fill)
	gc.SetLineWidth(10.0)
	gc.FillStroke()
}
