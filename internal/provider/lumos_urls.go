package provider

const (
	PERMISSIONS_URL      = "appstore/requestable_permissions"
	USERS_URL            = "users"
	PERMISSION_BY_ID_URL = PERMISSIONS_URL + "/%s"
	USER_BY_ID_URL       = USERS_URL + "/%s"
	USER_BY_EMAIL_URL    = USERS_URL + "?search_term=%s&exact_match=true"
)
