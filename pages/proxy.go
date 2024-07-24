package pages

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"

	"image"
	_ "image/jpeg" // register JPEG decoder

	// xxx: if "panic: image: unknown format", register more formats may help

	core_http "codeberg.org/vnpower/pixivfe/v2/core/http"
	"codeberg.org/vnpower/pixivfe/v2/core/zip"
	"github.com/gofiber/fiber/v2"
	"github.com/kettek/apng"
	"github.com/tidwall/gjson"
)

func SPximgProxy(c *fiber.Ctx) error {
	URL := fmt.Sprintf("https://s.pximg.net/%s", c.Params("*"))
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return err
	}
	req = req.WithContext(c.Context())

	// Make the request
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	c.Set("Content-Type", resp.Header.Get("Content-Type"))

	return c.Send([]byte(body))
}

func IPximgProxy(c *fiber.Ctx) error {
	URL := fmt.Sprintf("https://i.pximg.net/%s", c.Params("*"))
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return err
	}
	req = req.WithContext(c.Context())
	req.Header.Add("Referer", "https://www.pixiv.net/")

	// Make the request
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	c.Set("Content-Type", resp.Header.Get("Content-Type"))

	return c.Send([]byte(body))
}

// relies on ugoira.com
func UgoiraProxy_mp4(c *fiber.Ctx) error {
	URL := fmt.Sprintf("https://ugoira.com/api/mp4/%s", c.Params("*"))
	req, err := http.NewRequestWithContext(c.Context(), "GET", URL, nil)
	if err != nil {
		return err
	}
	req = req.WithContext(c.Context())

	// Make the request
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	c.Set("Content-Type", resp.Header.Get("Content-Type"))

	return c.Send([]byte(body))

	// todo:
	// delay... where do you get the delay?

}

func UgoiraProxy_apng(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return errors.New("No ID provided.")
	}
	
	URL_meta := core_http.GetUgoiraMetaURL(id)
	
	ugoira_metadata, err := core_http.UnwrapWebAPIRequest(c.Context(), URL_meta, "")
	if err != nil {
		return err
	}
	delays := []uint16{}
	for _, o := range gjson.Get(ugoira_metadata, "frames").Array() {
		delays = append(delays, uint16(o.Get("delay").Int()))
	}

	URL_zip := gjson.Get(ugoira_metadata, "originalSrc").Str

	req, err := http.NewRequestWithContext(c.Context(), "GET", URL_zip, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Referer", "https://www.pixiv.net/")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("while reading respose body: %w", err)
	}
	if resp.StatusCode != 200 {

		return errors.New(string(body))
	}
	r := bytes.NewReader(body)
	c.Set("Content-Type", "image/apng")
	return zip2apng(r, c, delays)
}

// each frame will be delayed $delayNumerator/delayDenominator$ seconds
func zip2apng(reader_zip io.ReadSeeker, writer_apng io.Writer, delays []uint16) error {
	img_apng := apng.APNG{}

	i := 0
	for {
		file, err := zip.ReadFile(reader_zip)
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
			DelayNumerator:   delays[i],
			DelayDenominator: 1000,
		}
		i += 1
		img_apng.Frames = append(img_apng.Frames, frame)
	}

	return apng.Encode(writer_apng, img_apng)
}
