# Every pre-approval rule of the app.
data "lumos_pre_approval_rules" "deploy_dashboard" {
  app_id = "<app id>"
}
