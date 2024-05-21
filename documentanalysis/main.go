package main

import (
	"context"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/textract"

	"github.com/aws/aws-sdk-go-v2/service/textract/types"
	"github.com/aws/aws-sdk-go/aws"
	fitz "github.com/gen2brain/go-fitz"
	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/llgcode/draw2d/draw2dkit"
)

var textractClient *textract.Client
var s3Client *s3.Client
var page int32

func main() {
	fmt.Println("Started Analysing")

	s3BucketName := ""
	objectKey := ""
	adapterId := ""
	textQueryString := "What is the No Claims Discount (NCD) Years"

	cfg, err := config.LoadDefaultConfig(context.Background())

	if err != nil {
		log.Fatal("failed to load config ", err)
	}

	textractClient = textract.NewFromConfig(cfg)
	s3Client = s3.NewFromConfig(cfg)

	extension := filepath.Ext(objectKey)

	// Remove the dot from the extension
	extension = extension[1:]

	if extension != "pdf" {
		resp, err := textractClient.AnalyzeDocument(context.Background(), &textract.AnalyzeDocumentInput{
			Document: &types.Document{
				S3Object: &types.S3Object{
					Bucket: &s3BucketName,
					Name:   &objectKey,
				},
			},
			AdaptersConfig: &types.AdaptersConfig{
				Adapters: []types.Adapter{
					{
						AdapterId: &adapterId,
						Version:   aws.String("2"),
						Pages:     []string{"*"},
					},
				},
			},
			FeatureTypes: []types.FeatureType{
				types.FeatureTypeQueries,
			},
			QueriesConfig: &types.QueriesConfig{
				Queries: []types.Query{
					{
						Text: &textQueryString,
					},
				},
			},
		})

		if err != nil {
			panic(err)
		}

		for _, block := range resp.Blocks {
			if block.BlockType == types.BlockTypeQueryResult {
				fmt.Println("Years found: ", *block.Text)
				fmt.Println("Confidence: ", *block.Confidence)
				fmt.Println("Coordinates", block.Geometry.BoundingBox.Width, block.Geometry.BoundingBox.Height, block.Geometry.BoundingBox.Left, block.Geometry.BoundingBox.Top)
				DrawBoundingBoxV2("", float64(block.Geometry.BoundingBox.Top), float64(block.Geometry.BoundingBox.Left), float64(block.Geometry.BoundingBox.Width), float64(block.Geometry.BoundingBox.Height))
			}
		}

		if err != nil {
			panic(err)
		}
	} else {
		startResp, err := textractClient.StartDocumentAnalysis(context.Background(), &textract.StartDocumentAnalysisInput{
			DocumentLocation: &types.DocumentLocation{
				S3Object: &types.S3Object{
					Bucket: &s3BucketName,
					Name:   &objectKey,
				},
			},
			AdaptersConfig: &types.AdaptersConfig{
				Adapters: []types.Adapter{
					{
						AdapterId: &adapterId,
						Version:   aws.String("2"),
						Pages:     []string{"*"},
					},
				},
			},
			FeatureTypes: []types.FeatureType{
				types.FeatureTypeQueries,
			},
			QueriesConfig: &types.QueriesConfig{
				Queries: []types.Query{
					{
						Text: &textQueryString,
					},
				},
			},
		})

		if err != nil {
			panic(err)
		}

		fmt.Println("Job Id: ", *startResp.JobId)
		time.Sleep(60 * time.Second)

		resp, err := textractClient.GetDocumentAnalysis(context.Background(), &textract.GetDocumentAnalysisInput{
			JobId: startResp.JobId,
		})

		if err != nil {
			panic(err)
		}

		// var result string
		var geometry types.Geometry
		fmt.Println("Job Status: ", resp.JobStatus)
		for _, block := range resp.Blocks {
			if block.BlockType == types.BlockTypeQueryResult {
				fmt.Println("Years found: ", *block.Text)
				fmt.Println("Page found: ", *block.Page)
				page = *block.Page
				geometry = *block.Geometry
				fmt.Println("Confidence: ", *block.Confidence)
				fmt.Println("Coordinates", block.Geometry.BoundingBox.Width, block.Geometry.BoundingBox.Height, block.Geometry.BoundingBox.Left, block.Geometry.BoundingBox.Top)
			}
		}

		if err != nil {
			panic(err)
		}

		// Download PDF file from S3
		respS3, err := s3Client.GetObject(context.Background(), &s3.GetObjectInput{
			Bucket: &s3BucketName,
			Key:    &objectKey,
		})

		if err != nil {
			panic(err)
		}

		// Create a file to write the downloaded content to
		file, err := os.Create("drawnoutput/output.pdf")
		if err != nil {
			panic(err)
		}
		defer file.Close()

		_, err = io.Copy(file, respS3.Body)
		if err != nil {
			panic(err)
		}

		DrawBoundingBoxV2("drawnoutput/output.pdf", float64(geometry.BoundingBox.Top), float64(geometry.BoundingBox.Left), float64(geometry.BoundingBox.Width), float64(geometry.BoundingBox.Height))
	}

	fmt.Println("Done")
}

func DrawBoundingBoxV2(existingImage string, top, left, width, height float64) {

	extension := filepath.Ext(existingImage)

	// Remove the dot from the extension
	extension = extension[1:]
	var myimage image.Image
	var imConfig image.Config
	if extension != "pdf" {
		existingImageFile, err := os.Open(existingImage)
		if err != nil {
			log.Fatal("Error opening image: ", err)
		}
		defer existingImageFile.Close()

		// Calling the generic image.Decode() will tell give us the data
		// and type of image it is as a string. We expect "png"
		imConfig, _, err = image.DecodeConfig(existingImageFile)
		if err != nil {
			log.Fatal("Error decoding image config: ", err)
		}

		_, err = existingImageFile.Seek(0, 0)
		if err != nil {
			log.Fatal("Error resetting offset: ", err)
		}

		myimage, _, err = image.Decode(existingImageFile)
		if err != nil {
			log.Fatal("Error decoding image: ", err)
		}
	} else {
		doc, err := fitz.New(existingImage)
		if err != nil {
			panic(err)
		}

		defer doc.Close()

		// Extract pages as images
		for n := 0; n < doc.NumPage(); n++ {
			if n != int(page-1) {
				continue
			}
			img, err := doc.Image(n)
			if err != nil {
				panic(err)
			}

			// Save the modified image to a new file.
			outputFile, err := os.Create("drawnoutput/output.png")
			if err != nil {
				panic(err)
			}
			defer outputFile.Close()

			// Encode the image as PNG and write to file.
			err = png.Encode(outputFile, img)
			myimage = img
			existingImageFile, err := os.Open("drawnoutput/output.png")
			if err != nil {
				log.Fatal("Error opening image: ", err)
			}
			defer existingImageFile.Close()

			// Calling the generic image.Decode() will tell give us the data
			// and type of image it is as a string. We expect "png"
			imConfig, _, err = image.DecodeConfig(existingImageFile)
			if err != nil {
				log.Fatal("Error decoding image config: ", err)
			}

			if err != nil {
				panic(err)
			}
		}

	}

	x := left * float64(imConfig.Width)
	y := top * float64(imConfig.Height)

	img := imageToRGBA(myimage)

	gc := draw2dimg.NewGraphicContext(img)

	gc.SetFillColor(color.NRGBA{255, 255, 255, 255})

	gc.SetStrokeColor(color.NRGBA{255, 0, 0, 100})
	gc.SetLineWidth(10)

	draw2dkit.Circle(gc, x, y, 50)
	gc.Stroke()

	draw2dimg.SaveToPngFile("output.png", img)
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
