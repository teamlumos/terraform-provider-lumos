package provider

const (
	PERMISSIONS_URL          = "appstore/requestable_permissions"
	PERMISSION_BY_ID_URL     = PERMISSIONS_URL + "/%s"
	USERS_URL                = "users"
	USER_BY_ID_URL           = USERS_URL + "/%s"
	USER_BY_EMAIL_URL        = USERS_URL + "?search_term=%s&exact_match=true"
	GROUPS_URL               = "groups"
	GROUP_BY_ID_URL          = GROUPS_URL + "/%s"
	GROUP_BY_NAME_URL        = GROUPS_URL + "?name=%s&exact_match=true"
	APPSTORE_APPS_URL        = "appstore/apps"
	APPSTORE_APP_BY_ID_URL   = APPSTORE_APPS_URL + "/%s"
	APPSTORE_APP_BY_NAME_URL = APPSTORE_APPS_URL + "?search_term=%s&exact_match=true"
)
