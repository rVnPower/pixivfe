package routes

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net/http"

	"codeberg.org/vnpower/pixivfe/v2/core"
	"codeberg.org/vnpower/pixivfe/v2/i18n"
	"codeberg.org/vnpower/pixivfe/v2/server/session"
	"codeberg.org/vnpower/pixivfe/v2/server/utils"
)

// NOTE: is the csrf protection by the upstream Pixiv API itself good enough, or do we need to implement our own?

func AddBookmarkRoute(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodPost {
		return i18n.Error("Method not allowed")
	}

	token := session.GetUserToken(r)
	csrf := session.GetCookie(r, session.Cookie_CSRF)

	if token == "" || csrf == "" {
		return PromptUserToLoginPage(w, r)
	}

	id := GetPathVar(r, "id")
	if id == "" {
		return i18n.Error("No ID provided.")
	}

	URL := "https://www.pixiv.net/ajax/illusts/bookmarks/add"
	payload := fmt.Sprintf(`{
        "illust_id": "%s",
        "restrict": 0,
        "comment": "",
        "tags": []
    }`, id)

	contentType := "application/json; charset=utf-8"
	_, err := core.API_POST(r.Context(), URL, payload, token, csrf, contentType)
	if err != nil {
		return err
	}

	utils.RedirectToWhenceYouCame(w, r)
	return nil
}

func DeleteBookmarkRoute(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodPost {
		return i18n.Error("Method not allowed")
	}

	token := session.GetUserToken(r)
	csrf := session.GetCookie(r, session.Cookie_CSRF)

	if token == "" || csrf == "" {
		return PromptUserToLoginPage(w, r)
	}

	id := GetPathVar(r, "id")
	if id == "" {
		return i18n.Error("No ID provided.")
	}

	URL := "https://www.pixiv.net/ajax/illusts/bookmarks/delete"
	payload := fmt.Sprintf("bookmark_id=%s", id)

	contentType := "application/x-www-form-urlencoded; charset=utf-8"
	_, err := core.API_POST(r.Context(), URL, payload, token, csrf, contentType)
	if err != nil {
		return err
	}

	utils.RedirectToWhenceYouCame(w, r)
	return nil
}

func LikeRoute(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodPost {
		return i18n.Error("Method not allowed")
	}

	token := session.GetUserToken(r)
	csrf := session.GetCookie(r, session.Cookie_CSRF)

	if token == "" || csrf == "" {
		return PromptUserToLoginPage(w, r)
	}

	id := GetPathVar(r, "id")
	if id == "" {
		return i18n.Error("No ID provided.")
	}

	URL := "https://www.pixiv.net/ajax/illusts/like"
	payload := fmt.Sprintf(`{"illust_id": "%s"}`, id)

	contentType := "application/json; charset=utf-8"
	_, err := core.API_POST(r.Context(), URL, payload, token, csrf, contentType)
	if err != nil {
		return err
	}

	utils.RedirectToWhenceYouCame(w, r)
	return nil
}

/*
NOTE: we're using the mobile API for FollowUserRoute and UnfollowUserRoute since it's an actual AJAX API
			instead of some weird php thing for the usual desktop routes (/bookmark_add.php and /rpc_group_setting.php)

			the desktop routes return HTML for the pixiv SPA when they feel like it and don't return helpful responses
			when you send a request that doesn't perfectly meet their specifications, making troubleshooting a nightmare

			for comparison, the mobile API worked first try without any issues

			interestingly enough, replicating the requests for the desktop routes via cURL worked fine but a Go implementation
			just refused to work
*/

func FollowUserRoute(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodPost {
		return i18n.Error("Method not allowed")
	}

	token := session.GetUserToken(r)
	csrf := session.GetCookie(r, session.Cookie_CSRF)

	if token == "" || csrf == "" {
		return PromptUserToLoginPage(w, r)
	}

	id := GetPathVar(r, "id")
	if id == "" {
		return i18n.Error("No user ID provided.")
	}

	isPrivate := r.FormValue("private") == "true"
	restrict := "0"
	if isPrivate {
		restrict = "1"
	}

	URL := "https://www.pixiv.net/touch/ajax_api/ajax_api.php"

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.WriteField("mode", "add_bookmark_user")
	writer.WriteField("restrict", restrict)
	writer.WriteField("user_id", id)
	writer.Close()

	fields := map[string]string{
		"mode":     "add_bookmark_user",
		"restrict": restrict,
		"user_id":  id,
	}
	_, err := core.API_POST(r.Context(), URL, fields, token, csrf, "")
	if err != nil {
		return err
	}

	utils.RedirectToWhenceYouCame(w, r)
	return nil
}

func UnfollowUserRoute(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodPost {
		return i18n.Error("Method not allowed")
	}

	token := session.GetUserToken(r)
	csrf := session.GetCookie(r, session.Cookie_CSRF)

	if token == "" || csrf == "" {
		return PromptUserToLoginPage(w, r)
	}

	id := GetPathVar(r, "id")
	if id == "" {
		return i18n.Error("No user ID provided.")
	}

	URL := "https://www.pixiv.net/touch/ajax_api/ajax_api.php"

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.WriteField("mode", "delete_bookmark_user")
	writer.WriteField("user_id", id)
	writer.Close()

	fields := map[string]string{
		"mode":    "delete_bookmark_user",
		"user_id": id,
	}
	_, err := core.API_POST(r.Context(), URL, fields, token, csrf, "")
	if err != nil {
		return err
	}

	utils.RedirectToWhenceYouCame(w, r)
	return nil
}
