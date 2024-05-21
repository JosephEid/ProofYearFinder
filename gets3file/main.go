package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func main() {
	// Load AWS credentials and configuration
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		fmt.Println("Error loading AWS configuration:", err)
		return
	}

	// Create an S3 client
	client := s3.NewFromConfig(cfg)

	inputFile := "../fetchS3Key/output-s3keys.csv"

	csvfile, err := os.Open(inputFile)
	if err != nil {
		log.Fatalln("Couldn't open the csv file", err)
	}

	// Parse the file
	r := csv.NewReader(csvfile)
	heading, err := r.Read()
	if err == io.EOF {
		log.Fatal("No records found in the csv file")
	}
	fmt.Println(heading)
	records, err := r.ReadAll()
	for i, record := range records {
		// Define the bucket and object key you want to download
		bucketName := "ncd-api-ncduploadclean02f2e29c-cjtk1he5hrxp"
		objectKey := record[2]
		extension := filepath.Ext(objectKey)

		// Remove the dot from the extension
		extension = extension[1:]

		// Create a file to write the downloaded object to
		outputFile, err := os.Create(fmt.Sprintf("output/%s.%s", record[1], extension))
		if err != nil {
			fmt.Println("Error creating file:", err)
			return
		}
		defer outputFile.Close()

		// Call GetObject to retrieve the object from S3
		resp, err := client.GetObject(context.Background(), &s3.GetObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(objectKey),
		})
		if err != nil {
			fmt.Println("Error retrieving object:", err, objectKey)
			continue
		}

		// Copy the object data to the file
		_, err = io.Copy(outputFile, resp.Body)
		if err != nil {
			fmt.Println("Error writing to file:", err)
			return
		}

		fmt.Println("File downloaded successfully.")

		if i%10 == 0 {
			time.Sleep(time.Millisecond * 500)
		}

	}

}
