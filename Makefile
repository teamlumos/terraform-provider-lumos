.PHONY: provider
provider: overlay.yaml
	speakeasy run

overlay.yaml: openapi-original.json openapi.json
	speakeasy overlay compare --schemas openapi-original.json --schemas openapi.json > overlay.yaml

.PHONY: testacc
testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m