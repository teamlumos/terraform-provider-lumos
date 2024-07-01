resource "lumos_pre_approval_rule" "eng_dev" {
  app_id        = local.aws_app_id
  justification = "All developers should be pre-approved for dev access."

  preapproved_permissions = [
    { id = local.dev_permission_id }
  ]
  preapproved_groups = [
    { id = local.dev_group_id }
  ]
}
