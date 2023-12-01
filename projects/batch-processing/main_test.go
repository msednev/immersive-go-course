package main

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"gopkg.in/gographics/imagick.v2/imagick"
)

func TestParseCsv(t *testing.T) {
	r := strings.NewReader(
		"url\nhttps://url1.com\nhttps://url2.com\nhttps://url3.com\n",
	)
	expected := []string{
		"https://url1.com",
		"https://url2.com",
		"https://url3.com",
	}
	actual := ParseCsv(r)
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("incorrect result, expected %v, got %v", expected, actual)
	}

}

func TestDownloadImage(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "image/jpg")
		w.Write([]byte("OK"))
	}))
	defer ts.Close()
	buf := new(bytes.Buffer)
	DownloadImage(ts.URL, buf)
	if buf.String() != "OK" {
		t.Fatal("")
	}
}

func TestGrayscaleMockError(t *testing.T) {
	c := &Converter{
		cmd: func(args []string) (*imagick.ImageCommandResult, error) {
			return nil, errors.New("not implemented")
		},
	}

	err := c.Grayscale("input.jpg", "output.jpg")
	if err == nil {
		t.Fatal("expected error")
	}

}

func TestGrayscaleMockCall(t *testing.T) {
	var args []string
	expected := []string{"convert", "input.jpg", "-set", "colorspace", "Gray", "output.jpg"}
	c := &Converter{
		cmd: func(a []string) (*imagick.ImageCommandResult, error) {
			args = a
			return &imagick.ImageCommandResult{
				Info: nil,
				Meta: "",
			}, nil
		},
	}

	err := c.Grayscale("input.jpg", "output.jpg")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(expected, args) {
		t.Fatalf("incorrect arguments: expected %v, got %v", expected, args)
	}
}
