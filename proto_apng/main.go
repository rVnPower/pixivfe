package main

import (
	"archive/zip"
	"fmt"
	"image"
	"io"
	"os"

	"github.com/kettek/apng"
)

func main() {
	filename := "119477424_ugoira600x600.zip"
	f, err := os.Create("out.apng")
	if err != nil {
		panic(err)
	}
	err = zip2apng(filename, f, 1000, 1000)
	if err != nil {
		panic(err)
	}
}

// APNG spec: https://wiki.mozilla.org/APNG_Specification

// each frame will be delayed (delayNumerator/delayDenominator) seconds
func zip2apng(filename string, w io.Writer, delayNumerator, delayDenominator uint16) error {
	// Open a zip archive for reading.
	r, err := zip.OpenReader(filename)
	if err != nil {
		return err
	}
	defer r.Close()

	img_apng := apng.APNG{}

	// Iterate through the files in the archive,
	// printing some of their contents.
	for _, f := range r.File {
		fmt.Printf("Contents of %s:\n", f.Name)
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()
		img, img_format, err := image.Decode(rc)
		_ = img_format
		if err != nil {
			return err
		}
		frame := apng.Frame{
			Image:            img,
			DelayNumerator:   delayNumerator,
			DelayDenominator: delayDenominator,
		}
		img_apng.Frames = append(img_apng.Frames, frame)
	}

	return apng.Encode(w, img_apng)
}
