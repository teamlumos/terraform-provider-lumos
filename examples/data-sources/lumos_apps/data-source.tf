data "lumos_apps" "my_apps" {
  exact_match = false
  name_search = "...my_name_search..."
  page        = 3
  size        = 19
}