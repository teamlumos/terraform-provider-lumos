.PHONY: provider
provider: overlay-v1.yaml
	@echo "Syncing gen.yaml version..."
	@VERSION="$(VERSION)"; \
	if [ -z "$$VERSION" ]; then \
	  VERSION=$$(git describe --tags --abbrev=0 --match 'v1.*' 2>/dev/null | sed 's/^v//'); \
	fi; \
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
overlay: openapi-v1.json original-openapi.json
	@echo "Generating overlay-v1.yaml from openapi-v1.json and original-openapi.json..."
	speakeasy overlay compare --before original-openapi.json --after openapi-v1.json > overlay-v1.yaml
	@echo "✓ overlay-v1.yaml generated"


.PHONY: testacc
testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m
