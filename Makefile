.PHONY: provider
provider: overlay.yaml
	@echo "Syncing gen.yaml version from latest git tag..."
	@VERSION=$$(git describe --tags --abbrev=0 2>/dev/null | sed 's/^v//'); \
	if [ -n "$$VERSION" ]; then \
	  sed "s/^  version: .*/  version: $$VERSION/" .speakeasy/gen.yaml > .speakeasy/gen.yaml.tmp && mv .speakeasy/gen.yaml.tmp .speakeasy/gen.yaml; \
	  echo "  Set terraform.version to $$VERSION"; \
	fi
	@echo "Running Speakeasy to generate the provider..."
	speakeasy run --skip-compile --skip-versioning
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
