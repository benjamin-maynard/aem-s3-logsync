package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var sourceBucket string
var targetBucket string
var bucketRegion string
var printOnly string

func main() {

	sourceBucket = os.Getenv("SOURCE_BUCKET_NAME")
	targetBucket = os.Getenv("TARGET_BUCKET_NAME")
	bucketRegion = os.Getenv("BUCKET_REGION")
	printOnly = os.Getenv("PRINT_ONLY")

	cmd := exec.Command("tail", "-f", "/locallog.log")

	// Create a pipe for the output of the script
	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error creating StdoutPipe for Cmd", err)
		return
	}

	scanner := bufio.NewScanner(cmdReader)
	go func() {
		for scanner.Scan() {

			// Pass the log line to the copytoS3() function to process
			go copytoS3(scanner.Text())

		}
	}()

	err = cmd.Start()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error starting Cmd", err)
		return
	}

	err = cmd.Wait()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error waiting for Cmd", err)
		return
	}

}

func copytoS3(file string) {

	// Split the File by last comma
	fileName := strings.Split(file, ",")[1]

	// Remove the Noise
	fileName = strings.ReplaceAll(fileName, "[", "")
	fileName = strings.ReplaceAll(fileName, " [/mnt2/s3-cache/upload/da/82/5d/", "")

	// Insert Missing "-" into filename from the log
	fileName = fileName[:4] + "-" + fileName[4:]

	if strings.ToLower(printOnly) != "true" {
		// Copy the Object from Source to Target
		sess := session.Must(session.NewSession())

		svc := s3.New(sess, &aws.Config{
			Region: aws.String(bucketRegion),
		})

		input := &s3.CopyObjectInput{
			Bucket:     aws.String(targetBucket),
			CopySource: aws.String((sourceBucket + "/" + fileName)),
			Key:        aws.String(fileName),
		}

		_, err := svc.CopyObject(input)
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				case s3.ErrCodeObjectNotInActiveTierError:
					fmt.Println(s3.ErrCodeObjectNotInActiveTierError, aerr.Error())
				default:
					fmt.Println(aerr.Error())
				}
			} else {
				// Print the error, cast err to awserr.Error to get the Code and
				// Message from an error.
				fmt.Println(err.Error())
			}
			return
		}

		fmt.Println("Successfully copied object " + fileName + ".")

	} else {
		fmt.Println("Print Only: Would have copied " + fileName + " from " + sourceBucket + "to" + targetBucket + ".")
	}

}
