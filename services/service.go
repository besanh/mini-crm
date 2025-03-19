package services

import "github.com/besanh/mini-crm/pkgs/oauth"

var (
	OAUTH2CONFIG               *oauth.OAuth2Config
	ENABLE_LOGIN_MULTI_SESSION bool = false

	// Google url get user info
	GOOGLE_URL_USER_INFO string = ""
)

const (
	OAUTH2_TOKEN string = "oauth2_token"

	// State in callback url
	OAUTH2_STATE string = "mini_crm_state"
)
