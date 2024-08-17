package serve

import (
	"bytes"
	"io"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/net/html"
)

// not working as middleware for some reason

func Rewrite_Link_preload(c *fiber.Ctx) error {
	err := c.Next()
	if err != nil {
		return err
	}

	contentType := string(c.Response().Header.ContentType())
	// log.Print("ct ", contentType)
	if contentType == "text/html; charset=utf-8" {
		images_to_preload := []string{}

		// here, the response body is not valid
		// no idea how to get c.Render to work here.
		// reading gofiber's code, c.Render uses bytebufferpool
		// i fail to understand how that makes sense
		z := html.NewTokenizer(bytes.NewReader(c.Response().Body()))
	scan_tokens:
		for {
			token_type := z.Next()
			switch token_type {
			case html.ErrorToken:
				err := z.Err()
				if err == io.EOF {
					break scan_tokens
				}
				return err
			case html.StartTagToken, html.SelfClosingTagToken, html.EndTagToken:
				tag_name, hasAttr := z.TagName()
				log.Print("tag name ", string(tag_name))
				if string(tag_name) == "img" && hasAttr {
					for hasAttr {
						var key, val []byte
						key, val, hasAttr = z.TagAttr()
						log.Print("key val ", string(key), string(val))
						if string(key) == "src" && strings.HasPrefix(string(val), "/proxy/i.pximg.net/img-master/") {
							images_to_preload = append(images_to_preload, string(val))
						}
					}
				}
			}
		}
		log.Print("preload ", images_to_preload)
	}

	return nil
}
