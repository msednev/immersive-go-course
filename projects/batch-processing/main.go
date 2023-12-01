package main

import (
	"encoding/csv"
	"flag"
	"log"
	"os"
	"io"
	"net/http"

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

func ParseCsv(r io.Reader) ([]string) {

	reader := csv.NewReader(r)
	reader.FieldsPerRecord = 1
	result := make([]string, 0)
	for i := 0;;i++ {
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
	if contentType != "image/jpg" {
		log.Fatalf("Bad content type %v", contentType)
	}

	written, err = io.Copy(dst, resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return written
} 

func main() {
	// Accept --input and --output arguments for the images
	inputFilepath := flag.String("input", "", "A path to an csv with a list of files to be processed")
	outputFilepath := flag.String("output", "", "A path to where the csv with results should be written")
	flag.Parse()

	// Ensure that both flags were set
	if *inputFilepath == "" || *outputFilepath == "" {
		flag.Usage()
		os.Exit(1)
	}

	// inputfiles := ParseCsv(csv)

	// Set up imagemagick
	imagick.Initialize()
	defer imagick.Terminate()

	// Log what we're going to do
	log.Printf("processing: %q to %q\n", *inputFilepath, *outputFilepath)

	// Build a Converter struct that will use imagick
	c := &Converter{
		cmd: imagick.ConvertImageCommand,
	}

	// Do the conversion!
	err := c.Grayscale(*inputFilepath, *outputFilepath)
	if err != nil {
		log.Printf("error: %v\n", err)
	}

	// Log what we did
	log.Printf("processed: %q to %q\n", *inputFilepath, *outputFilepath)
}
