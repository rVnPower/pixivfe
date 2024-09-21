package routes

import (
	"errors"
	"fmt"
	"net/http"

	"codeberg.org/vnpower/pixivfe/v2/core"
	"codeberg.org/vnpower/pixivfe/v2/server/session"
	"codeberg.org/vnpower/pixivfe/v2/server/utils"
)

func AddBookmarkRoute(w http.ResponseWriter, r *http.Request) error {
	token := session.GetUserToken(r)
	csrf := session.GetCookie(r, session.Cookie_CSRF)

	if token == "" || csrf == "" {
		return PromptUserToLoginPage(w, r)
	}

	id := GetPathVar(r, "id")
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
	if err := core.API_POST(r.Context(), URL, payload, token, csrf, true); err != nil {
		return err
	}

	utils.RedirectToWhenceYouCame(w, r)
	return nil
}

func DeleteBookmarkRoute(w http.ResponseWriter, r *http.Request) error {
	token := session.GetUserToken(r)
	csrf := session.GetCookie(r, session.Cookie_CSRF)

	if token == "" || csrf == "" {
		return PromptUserToLoginPage(w, r)
	}

	id := GetPathVar(r, "id")
	if id == "" {
		return errors.New("No ID provided.")
	}

	// You can't unlike
	URL := "https://www.pixiv.net/ajax/illusts/bookmarks/delete"
	payload := fmt.Sprintf(`bookmark_id=%s`, id)
	if err := core.API_POST(r.Context(), URL, payload, token, csrf, false); err != nil {
		return err
	}

	utils.RedirectToWhenceYouCame(w, r)
	return nil
}

func LikeRoute(w http.ResponseWriter, r *http.Request) error {
	token := session.GetUserToken(r)
	csrf := session.GetCookie(r, session.Cookie_CSRF)

	if token == "" || csrf == "" {
		return PromptUserToLoginPage(w, r)
	}

	id := GetPathVar(r, "id")
	if id == "" {
		return errors.New("No ID provided.")
	}

	URL := "https://www.pixiv.net/ajax/illusts/like"
	payload := fmt.Sprintf(`{"illust_id": "%s"}`, id)
	if err := core.API_POST(r.Context(), URL, payload, token, csrf, true); err != nil {
		return err
	}

	utils.RedirectToWhenceYouCame(w, r)
	return nil
}
