package main

import (
	"bytes"
	"fmt"
	"image"
	_ "image/jpeg" // register JPEG decoder
	// xxx: if "panic: image: unknown format", register more formats may help

	"io"
	"os"

	"codeberg.org/vnpower/pixivfe/v2/proto_apng/zip"

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
	r, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer r.Close()

	img_apng := apng.APNG{}

	for {
		file, err := zip.ReadFile(r)
		if err == zip.ErrFormat {
			break
		}
		if err != nil {
			return err
		}
		fmt.Println("len", len(file.Content))
		img, img_format, err := image.Decode(bytes.NewReader(file.Content))
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
