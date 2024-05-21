package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	policyapi "github.com/aviva-verde/policyapiclient"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func main() {

	fmt.Println("Hello World")

	fetchS3 := true
	fetchNCDYears := false

	inputFile := "../ncd-policyid-submissionid.csv"

	// read csv file
	readCsvFile(inputFile, fetchS3, fetchNCDYears)

	// fmt.Println("S3 Key: ", s3Key)
	fmt.Println("Done")
}

func readCsvFile(inputFile string, fetchS3 bool, fetchNCDYears bool) {

	// Open the file
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

	fmt.Println("started reading file with heading: ", heading)

	for {
		// Read each record from csv
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		if fetchS3 {
			key := prepfetchS3Key(record[0], record[1])
			s3Keys := fetchS3Key(key)
			file, err := os.OpenFile("output-s3keys.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

			if err != nil {
				log.Fatal(err)
			}

			if len(s3Keys) > 0 {
				for _, s3Key := range s3Keys {
					writeRecord := record
					writeRecord = append(record, s3Key)
					file.WriteString(writeRecord[0] + "," + writeRecord[1] + "," + writeRecord[2] + "\n")
				}
			}
			defer file.Close()
		}

		if fetchNCDYears {
			file, err := os.OpenFile("output-ncd-years.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

			if err != nil {
				log.Fatal(err)
			}

			ncdYears := fetchNCDYEars(record[0])
			record = append(record, strconv.Itoa(int(ncdYears)))
			file.WriteString(record[0] + "," + record[1] + "," + record[2] + "\n")
			defer file.Close()
		}

	}

	csvfile.Close()
}

func prepfetchS3Key(policyId string, submissionid string) string {
	return "policy/" + policyId + "/submission/" + submissionid
}

func fetchS3Key(pk string) (keys []string) {
	tableName := "ncd-api-ncdTable"
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		panic(err)
	}
	svc := dynamodb.NewFromConfig(cfg)
	out, err := svc.Query(context.TODO(), &dynamodb.QueryInput{
		TableName: aws.String(tableName),
		// IndexName:              aws.String("_pk"),
		KeyConditionExpression: aws.String("#_pk = :_pk"),
		ExpressionAttributeNames: map[string]string{
			"#_pk": "_pk",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":_pk": &types.AttributeValueMemberS{Value: pk},
		},
	})
	if err != nil {
		panic(err)
	}

	if len(out.Items) > 0 {
		for _, item := range out.Items {
			if item["currentStatus"].(*types.AttributeValueMemberS).Value == "upload_complete" {
				keys = append(keys, item["s3Key"].(*types.AttributeValueMemberS).Value)
			}
		}
	}

	return keys
}

func fetchNCDYEars(pk string) int64 {
	pc, err := policyapi.NewClient("https://yz4by9nuy0.execute-api.eu-west-1.amazonaws.com/prod/")
	if err != nil {
		panic(err)
	}
	p, ok, err := pc.GetLatestPolicyVersion(context.Background(), pk)
	if err != nil {
		panic(err)
	}
	if !ok {
		panic("No policy found")
	}
	return p.CoverDetails.NoClaimsBonus
}
