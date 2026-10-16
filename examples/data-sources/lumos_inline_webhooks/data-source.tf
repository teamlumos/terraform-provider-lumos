# Finds an inline webhook by name, for example to validate access requests.
data "lumos_inline_webhooks" "all" {}

output "request_validation_webhook_id" {
  value = one([for w in data.lumos_inline_webhooks.all.inline_webhooks : w.id if w.name == "Check manager approval"])
}
