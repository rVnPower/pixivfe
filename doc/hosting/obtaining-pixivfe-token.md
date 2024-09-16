# Obtaining the `PIXIVFE_TOKEN` cookie

This guide covers how to obtain the `PIXIVFE_TOKEN` cookie from your Pixiv account, which is necessary for authenticating with the Pixiv API.

!!! warning
    You should create an entirely new account for this to avoid account theft. And also, PixivFE will get contents **from your account.** You might not want people to know what kind of illustrations you like :P
    
    For now, the only page that may contain contents that is relevant to you is the discovery page. Be careful if you are using your main account.

## 1. Log in to Pixiv

Log in to the Pixiv account you want to use. Upon logging in, you should see the Pixiv landing page. If you are already logged in, simply navigate to the landing page.

![The URL of the landing page](https://files.catbox.moe/7dbv3e.png)

## 2. Open developer tools

### For Firefox

Press `F12` to open the Firefox Developer Tools. Switch to the `Storage` tab.

![Storage tab on Firefox](https://files.catbox.moe/mra6rs.png)

### For Chrome

Press `F12` to open the Chrome Developer Tools. Switch to the `Application` tab.

![Application tab on Chrome](https://files.catbox.moe/jqpcw2.png)

## 3. Locate the Cookie

### For Firefox

In the left sidebar, expand the `Cookies` section and select `www.pixiv.net`. This is where you will find your authentication cookie.

Locate the cookie with the key `PHPSESSID`. The value next to this key is your account's token.

![Cookie on Firefox](https://files.catbox.moe/zb16o8.png)

### For Chrome

In the left sidebar, find the `Storage` section. Expand the `Cookies` subsection and select `www.pixiv.net`. This is where you will find your authentication cookie.

Locate the cookie with the key `PHPSESSID`. The value next to this key is your account's token.

![PHPSESSID on Chrome-based browsers](https://files.catbox.moe/8wu9f0.png)

## 4. Set the environment variable

Copy the token value obtained in the previous step. If deploying with Docker, set it as the `PIXIVFE_TOKEN` environment variable in your configuration.

## 5. Enabling R-18G Artworks (Optional)

For PixivFE to show R-18G artworks, the account used by PixivFE has to enable the "Show ero-guro content (R-18G)" option on Pixiv. Here's how to do it:

1. Go to Pixiv's [display settings page](https://www.pixiv.net/settings/viewing).
2. Enable the "Show ero-guro content (R-18G)" option.

To test if R-18G content is now visible:

1. Go to this [search endpoint](https://www.pixiv.net/ajax/search/artworks/gore).
2. Search for any appearances of "R-18G" in the results.
3. If you disable the R-18G option and search again, you shouldn't see any R-18G artworks in the results.

## Additional notes

- The token format resembles: `123456_AaBbccDDeeFFggHHIiJjkkllmMnnooPP`
    - The underscore separates your **member ID (left side)** from a **random string (right side)**
- Logging out of Pixiv will reset the token. Always verify your token is current before reporting issues.
- Some Chrome-related content was sourced from [Nandaka's guide](https://github.com/Nandaka/PixivUtil2/wiki#pixiv-login-using-cookie).
