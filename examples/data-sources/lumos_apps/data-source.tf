data "lumos_apps" "my_apps" {
  connection_source = "API"
  page              = 1
  size              = 50
}