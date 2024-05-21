package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/textract"
	"github.com/aws/aws-sdk-go-v2/service/textract/types"
	"github.com/aws/aws-sdk-go/aws"
)

var textractClient *textract.Client

func main() {
	useAdapter := false
	fmt.Println("Started Analysing")

	s3BucketName := ""
	proofFile := ""
	adapterId := ""
	textQueryString := "What is the No Claims Discount (NCD) Years"

	cfg, err := config.LoadDefaultConfig(context.Background())

	if err != nil {
		log.Fatal("failed to load config ", err)
	}

	textractClient = textract.NewFromConfig(cfg)

	if useAdapter {
		// using adapter
		resp, err := textractClient.AnalyzeDocument(context.Background(), &textract.AnalyzeDocumentInput{
			Document: &types.Document{
				S3Object: &types.S3Object{
					Bucket: &s3BucketName,
					Name:   &proofFile,
				},
			},
			AdaptersConfig: &types.AdaptersConfig{
				Adapters: []types.Adapter{
					{
						AdapterId: &adapterId,
						Version:   aws.String("1.0"),
					},
				},
			},
		})

		if err != nil {
			panic(err)
		}

		fmt.Println("Response: ", resp)
	}

	// using keyword
	resp, err := textractClient.AnalyzeDocument(context.Background(), &textract.AnalyzeDocumentInput{
		Document: &types.Document{
			S3Object: &types.S3Object{
				Bucket: &s3BucketName,
				Name:   &proofFile,
			},
		},
		FeatureTypes: []types.FeatureType{
			types.FeatureTypeQueries,
			// types.FeatureTypeLayout,
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
		fmt.Println("Block Type: ", block.BlockType)
		if block.BlockType == types.BlockTypeQueryResult {
			fmt.Println("Line: ", block.Text)
		}
	}
	fmt.Println("Done")
}
