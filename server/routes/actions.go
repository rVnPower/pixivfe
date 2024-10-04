package routes

import (
	"fmt"
	"net/http"

	"codeberg.org/vnpower/pixivfe/v2/core"
	"codeberg.org/vnpower/pixivfe/v2/i18n"
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
		return i18n.Error("No ID provided.")
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

	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, `<button type="button" class="btn custom-btn-secondary bg-charcoal-surface2 text-nowrap me-2" hx-post="/self/deleteBookmark/%s" hx-target="#bookmark-button" hx-swap="outerHTML">
			<i class="bi bi-heart-fill me-2"></i>Bookmarked
		</button>`, id)
		return nil
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
		return i18n.Error("No ID provided.")
	}

	// You can't unlike
	URL := "https://www.pixiv.net/ajax/illusts/bookmarks/delete"
	payload := fmt.Sprintf(`bookmark_id=%s`, id)
	if err := core.API_POST(r.Context(), URL, payload, token, csrf, false); err != nil {
		return err
	}

	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, `<button type="button" class="btn custom-btn-secondary bg-charcoal-surface2 text-nowrap me-2" hx-post="/self/addBookmark/%s" hx-target="#bookmark-button" hx-swap="outerHTML">
			<i class="bi bi-heart me-2"></i>Bookmark
		</button>`, id)
		return nil
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
		return i18n.Error("No ID provided.")
	}

	URL := "https://www.pixiv.net/ajax/illusts/like"
	payload := fmt.Sprintf(`{"illust_id": "%s"}`, id)
	if err := core.API_POST(r.Context(), URL, payload, token, csrf, true); err != nil {
		return err
	}

	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, `<button type="button" class="btn custom-btn-secondary bg-charcoal-surface2 text-nowrap" disabled>
			<i class="bi bi-hand-thumbs-up-fill me-2"></i>Liked
		</button>`)
		return nil
	}

	utils.RedirectToWhenceYouCame(w, r)
	return nil
}
