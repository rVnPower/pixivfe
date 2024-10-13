package routes

import (
	"fmt"
	"net/http"
	"net/url"

	"codeberg.org/vnpower/pixivfe/v2/core"
	"codeberg.org/vnpower/pixivfe/v2/i18n"
	"codeberg.org/vnpower/pixivfe/v2/server/audit"
	"codeberg.org/vnpower/pixivfe/v2/server/session"
	"codeberg.org/vnpower/pixivfe/v2/server/utils"
	"go.uber.org/zap"
)

// getLogger initializes the audit logger lazily
func getLogger() *zap.Logger {
	return audit.GetLogger()
}

// NOTE: is the csrf protection by the upstream Pixiv API itself good enough, or do we need to implement our own?

func AddBookmarkRoute(w http.ResponseWriter, r *http.Request) error {
	logger := getLogger()

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
	_, err := core.API_POST(r.Context(), URL, payload, token, csrf, true)
	if err != nil {
		logger.Error("API call failed", zap.Error(err))
		return err
	}

	utils.RedirectToWhenceYouCame(w, r)
	return nil
}

func DeleteBookmarkRoute(w http.ResponseWriter, r *http.Request) error {
	logger := getLogger()

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

	// You can't unlike
	URL := "https://www.pixiv.net/ajax/illusts/bookmarks/delete"
	payload := fmt.Sprintf(`bookmark_id=%s`, id)
	_, err := core.API_POST(r.Context(), URL, payload, token, csrf, false)
	if err != nil {
		logger.Error("API call failed", zap.Error(err))
		return err
	}

	utils.RedirectToWhenceYouCame(w, r)
	return nil
}

func LikeRoute(w http.ResponseWriter, r *http.Request) error {
	logger := getLogger()

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
	_, err := core.API_POST(r.Context(), URL, payload, token, csrf, true)
	if err != nil {
		logger.Error("API call failed", zap.Error(err))
		return err
	}

	utils.RedirectToWhenceYouCame(w, r)
	return nil
}

func FollowUserRoute(w http.ResponseWriter, r *http.Request) error {
	logger := getLogger()

	logger.Debug("FollowUserRoute called")

	if r.Method != http.MethodPost {
		logger.Debug("Method not allowed", zap.String("method", r.Method))
		return i18n.Error("Method not allowed")
	}

	token := session.GetUserToken(r)
	csrf := session.GetCookie(r, session.Cookie_CSRF)

	if token == "" || csrf == "" {
		logger.Debug("User not logged in or missing CSRF")
		return PromptUserToLoginPage(w, r)
	}

	id := GetPathVar(r, "id")
	if id == "" {
		logger.Debug("No user ID provided")
		return i18n.Error("No user ID provided.")
	}
	logger.Debug("Following user", zap.String("user_id", id))

	isPrivate := r.FormValue("private") == "true"
	restrict := "0"
	if isPrivate {
		restrict = "1"
	}
	logger.Debug("Follow privacy setting", zap.Bool("isPrivate", isPrivate))

	URL := "https://www.pixiv.net/bookmark_add.php"
	payload := url.Values{
		"mode":     {"add"},
		"type":     {"user"},
		"user_id":  {id},
		"tag":      {""},
		"restrict": {restrict},
		"format":   {"json"},
	}.Encode()

	logger.Debug("Making API call to follow user", zap.String("URL", URL))
	resp, err := core.API_POST(r.Context(), URL, payload, token, csrf, false)
	if err != nil {
		logger.Error("API call failed", zap.Error(err))
		return err
	}
	logger.Debug("API call successful", zap.Int("StatusCode", resp.StatusCode), zap.String("Body", resp.Body))

	logger.Debug("Redirecting user")
	utils.RedirectToWhenceYouCame(w, r)
	return nil
}

func UnfollowUserRoute(w http.ResponseWriter, r *http.Request) error {
	logger := getLogger()

	logger.Debug("UnfollowUserRoute called")

	if r.Method != http.MethodPost {
		logger.Debug("Method not allowed", zap.String("method", r.Method))
		return i18n.Error("Method not allowed")
	}

	token := session.GetUserToken(r)
	csrf := session.GetCookie(r, session.Cookie_CSRF)

	if token == "" || csrf == "" {
		logger.Debug("User not logged in or missing CSRF")
		return PromptUserToLoginPage(w, r)
	}

	id := GetPathVar(r, "id")
	if id == "" {
		logger.Debug("No user ID provided")
		return i18n.Error("No user ID provided.")
	}
	logger.Debug("Unfollowing user", zap.String("user_id", id))

	URL := "https://www.pixiv.net/rpc_group_setting.php"
	payload := url.Values{
		"mode": {"del"},
		"type": {"bookuser"},
		"id":   {id},
	}.Encode()

	logger.Debug("Making API call to unfollow user", zap.String("URL", URL))
	_, err := core.API_POST(r.Context(), URL, payload, token, csrf, false)
	if err != nil {
		logger.Error("API call failed", zap.Error(err))
		return err
	}

	logger.Debug("API call successful")

	logger.Debug("Redirecting user")
	utils.RedirectToWhenceYouCame(w, r)
	return nil
}
