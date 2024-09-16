data "lumos_apps" "my_apps" {
  exact_match = true
  name_search = "...my_name_search..."
  page        = 3
  size        = 19
}