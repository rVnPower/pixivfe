package routes

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"

	"codeberg.org/vnpower/pixivfe/v2/session"
	"github.com/tidwall/gjson"
)

func pixivPostRequest(r CompatRequest, url, payload, token, csrf string, isJSON bool) error {
	requestBody := []byte(payload)

	req, err := http.NewRequestWithContext(r.Context(), "POST", url, bytes.NewBuffer(requestBody))
 if err != nil {
   return err
 }
	req.Header.Add("User-Agent", "Mozilla/5.0")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Cookie", "PHPSESSID="+token)
	req.Header.Add("x-csrf-token", csrf)

	if isJSON {
		req.Header.Add("Content-Type", "application/json; charset=utf-8")
	} else {
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=utf-8")
	}
	// req.AddCookie(&http.Cookie{
	// 	Name:  "PHPSESSID",
	// 	Value: token,
	// })

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.New("Failed to do this action.")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return errors.New("Cannot parse the response from Pixiv. Please report this issue.")
	}
	body_s := string(body)
	if !gjson.Valid(body_s) {
		return fmt.Errorf("Invalid JSON: %v", body_s)
	}
	errr := gjson.Get(body_s, "error")

	if !errr.Exists() {
		return errors.New("Incompatible request body.")
	}

	if errr.Bool() {
		return errors.New("Pixiv: Invalid request.")
	}
	return nil
}

func AddBookmarkRoute(w http.ResponseWriter, r CompatRequest) error {
	token := session.GetPixivToken(r.Request)
	csrf := session.GetCookie(r.Request, session.Cookie_CSRF)

	if token == "" || csrf == "" {
		return PromptUserToLoginPage(w, r)
	}

	id := r.Params("id")
	if id == "" {
		return errors.New("No ID provided.")
	}

	URL := "https://www.pixiv.net/ajax/illusts/bookmarks/add"
	payload := fmt.Sprintf(`{
"illust_id": "%s",
"restrict": 0,
"comment": "",
"tags": []
}`, id)
	if err := pixivPostRequest(r, URL, payload, token, csrf, true); err != nil {
		return err
	}

	return SendString(w, "Success")
}

func DeleteBookmarkRoute(w http.ResponseWriter, r CompatRequest) error {
	token := session.GetPixivToken(r.Request)
	csrf := session.GetCookie(r.Request, session.Cookie_CSRF)

	if token == "" || csrf == "" {
		return PromptUserToLoginPage(w, r)
	}

	id := r.Params("id")
	if id == "" {
		return errors.New("No ID provided.")
	}

	// You can't unlike
	URL := "https://www.pixiv.net/ajax/illusts/bookmarks/delete"
	payload := fmt.Sprintf(`bookmark_id=%s`, id)
	if err := pixivPostRequest(r, URL, payload, token, csrf, false); err != nil {
		return err
	}

	return SendString(w, "Success")
}

func LikeRoute(w http.ResponseWriter, r CompatRequest) error {
	token := session.GetPixivToken(r.Request)
	csrf := session.GetCookie(r.Request, session.Cookie_CSRF)

	if token == "" || csrf == "" {
		return PromptUserToLoginPage(w, r)
	}

	id := r.Params("id")
	if id == "" {
		return errors.New("No ID provided.")
	}

	URL := "https://www.pixiv.net/ajax/illusts/like"
	payload := fmt.Sprintf(`{"illust_id": "%s"}`, id)
	if err := pixivPostRequest(r, URL, payload, token, csrf, true); err != nil {
		return err
	}

	return SendString(w, "Success")
}

func SendString(w http.ResponseWriter, text string) error {
	w.Header().Set("content-type", "text/plain")
	_, err :=  w.Write([]byte(text))
	return err
}