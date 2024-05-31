package provider

const (
	PERMISSIONS_URL                = "appstore/requestable_permissions"
	PERMISSION_BY_ID_URL           = PERMISSIONS_URL + "/%s"
	USERS_URL                      = "users"
	USER_BY_ID_URL                 = USERS_URL + "/%s"
	USER_BY_EMAIL_URL              = USERS_URL + "?search_term=%s&exact_match=true"
	GROUPS_URL                     = "groups"
	GROUP_BY_ID_URL                = GROUPS_URL + "/%s"
	GROUP_BY_NAME_URL              = GROUPS_URL + "?name=%s&exact_match=true"
	APPS_URL                       = "apps"
	APP_BY_ID_URL                  = APPS_URL + "/%s"
	APPSTORE_APPS_URL              = "appstore/apps"
	APPSTORE_APP_BY_ID_URL         = APPSTORE_APPS_URL + "/%s"
	APPSTORE_APP_SETTING_BY_ID_URL = APPSTORE_APP_BY_ID_URL + "/settings"
	APPSTORE_APP_BY_NAME_URL       = APPSTORE_APPS_URL + "?name_search=%s&exact_match=true"
)
