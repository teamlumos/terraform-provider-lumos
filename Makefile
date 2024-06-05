.PHONY: provider
provider: overlay.yaml
	speakeasy run

overlay.yaml: openapi.json
	curl https://raw.githubusercontent.com/teamlumos/lumos-openapi-spec/main/openapi.json > original-openapi.json
	speakeasy overlay compare --schemas original-openapi.json --schemas openapi.json > overlay.yaml

.PHONY: testacc
testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m
