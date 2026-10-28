variable "run_id" {
  description = "Random suffix for this run. Put it in the name of everything a case creates."
  type        = string
}

variable "server_url" {
  description = "From LUMOS_API_URL."
  type        = string
}

variable "group_id" {
  description = "From GROUP_ID. Cases that use it say `# requires: GROUP_ID`."
  type        = string
  default     = null
}

variable "tba_label" {
  description = "From LUMOS_TEST_TBA_LABEL. Cases that use it say `# requires: LUMOS_TEST_TBA_LABEL`."
  type        = string
  default     = null
}
