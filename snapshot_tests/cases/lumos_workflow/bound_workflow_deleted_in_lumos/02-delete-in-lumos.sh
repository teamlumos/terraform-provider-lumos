# The bound workflow is deleted in Lumos. If the API refuses to delete a bound
# workflow, this step and the case fail, with the refusal in terraform.log.
api DELETE "/v1/workflows/$(state_id lumos_workflow.this)"
