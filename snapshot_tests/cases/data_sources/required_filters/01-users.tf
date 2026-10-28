# expect: error
# lumos_users needs search; it refuses at plan without calling the API.
data "lumos_users" "all" {}
