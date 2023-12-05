package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"gopkg.in/gographics/imagick.v2/imagick"
)

type ConvertImageCommand func(args []string) (*imagick.ImageCommandResult, error)

type Converter struct {
	cmd ConvertImageCommand
}

func (c *Converter) Grayscale(inputFilepath string, outputFilepath string) error {
	// Convert the image to grayscale using imagemagick
	// We are directly calling the convert command
	_, err := c.cmd([]string{
		"convert", inputFilepath, "-set", "colorspace", "Gray", outputFilepath,
	})
	return err
}

// Parse provided csv and return an array of urls
func ParseCsv(r io.Reader) []string {

	reader := csv.NewReader(r)
	reader.FieldsPerRecord = 1
	result := make([]string, 0)
	for i := 0; ; i++ {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println(err)
		}
		// checking if first column has header called 'url'
		if i == 0 && record[0] != "url" {
			log.Fatalf("Expected column `url`, found %v", record)
		}
		if i > 0 {
			result = append(result, record[0])
		}
	}
	return result
}

func DownloadImage(url string, dst io.Writer) (written int64) {
	resp, err := http.Get(url)
	defer resp.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Bad response status: %v", resp.Status)
	}

	contentType := resp.Header.Get("content-type")
	if contentType != "image/jpeg" {
		log.Fatalf("Bad content type %v", contentType)
	}

	written, err = io.Copy(dst, resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return written
}

func UploadImage(bucket *string, filename *string) error {
	file, err := os.Open(*filename)
	defer file.Close()
	if err != nil {
		return err
	}
	sess, err := session.NewSession()
	if err != nil {
		return err
	}
	uploader := s3manager.NewUploader(sess)
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: bucket,
		Key:    filename,
		Body:   file,
	})
	return nil
}

func main() {
	// Accept --input and --output arguments for the csv files
	inputCsv := flag.String("input", "", "A path to an csv with a list of files to be processed")
	outputCsv := flag.String("output", "", "A path to where the csv with results should be written")
	flag.Parse()

	// Ensure that both flags were set
	if *inputCsv == "" || *outputCsv == "" {
		flag.Usage()
		os.Exit(1)
	}
	file, err := os.Open(*inputCsv)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
	inputUrls := ParseCsv(file)
	inputFiles := make([]string, len(inputUrls))
	for i, url := range inputUrls {
		filename := filepath.Join("inputs", strconv.Itoa(i) + ".jpg")
		file, err := os.Create(filename)
		if err != nil {
			log.Println(err)
		}
		DownloadImage(url, file)
		inputFiles[i] = filename
		file.Close()
	}

	// Set up imagemagick
	imagick.Initialize()
	defer imagick.Terminate()

	// Build a Converter struct that will use imagick
	c := &Converter{
		cmd: imagick.ConvertImageCommand,
	}

	// Do the conversion!
	fmt.Println(inputFiles)
	for _, inputFile := range inputFiles {
		// Log what we're going to do
		inputFileNameExt := filepath.Ext(inputFile)
		inputFileNameBase := filepath.Base(inputFile)
		inputFileWithoutExt := inputFileNameBase[:len(inputFileNameBase) - len(inputFileNameExt)]
		outputFile := filepath.Join("outputs", inputFileWithoutExt + "_bw" + inputFileNameExt)
		log.Printf("processing: %q to %q\n", inputFile, outputFile)
		err = c.Grayscale(inputFile, outputFile)
		if err != nil {
			log.Printf("error: %v\n", err)
		}
		// Log what we did
		log.Printf("processed: %q to %q\n", inputFile, outputFile)
	}
}
