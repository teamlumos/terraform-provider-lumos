.PHONY: provider
provider: overlay.yaml
	@echo "Running Speakeasy to generate the provider..."
	speakeasy run --skip-compile
	@echo "Running go mod tidy..."
	@go mod tidy
	@echo "Compiling provider..."
	@go build ./internal/provider/...
	@echo "Running tests..."
	@go test ./internal/provider/...
	@echo "Generating documentation..."
	@if [ -d docs/guides ]; then cp -R docs/guides .docs-guides-backup; fi
	@go generate ./...
	@if [ -d .docs-guides-backup ]; then rm -rf docs/guides && mv .docs-guides-backup docs/guides; fi
	@echo "✓ Successfully generated and compiled the provider and ran tests"


.PHONY: overlay
overlay: openapi.json original-openapi.json
	@echo "Generating overlay.yaml from openapi.json and original-openapi.json..."
	speakeasy overlay compare --schemas original-openapi.json --schemas openapi.json > overlay.yaml
	@echo "✓ overlay.yaml generated"


.PHONY: testacc
testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m
