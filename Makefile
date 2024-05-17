COLOR_OK=\\x1b[0;32m
COLOR_NONE=\x1b[0m
COLOR_WARNING=\x1b[33;05m
COLOR_LUMOS=\x1B[34;01m

help:
	@echo "$(COLOR_OK)Lumos Terraform Provider(COLOR_NONE)"
	@echo ""
	@echo "$(COLOR_WARNING)Usage:$(COLOR_NONE)"
	@echo "$(COLOR_OK)  make [command]$(COLOR_NONE)"
	@echo ""
	@echo "$(COLOR_WARNING)Available commands:$(COLOR_NONE)"
	@echo "$(COLOR_OK)  build               					  Generate the provider and documentation, and build it$(COLOR_NONE)"
	@echo "$(COLOR_OK)  provider-code-spec               Generates provider-code-spec.json from generator_config.yml and spec.json$(COLOR_NONE)"
	@echo "$(COLOR_OK)  provider              			  Generates Terraform Provider$(COLOR_NONE)"
	@echo "$(COLOR_OK)  sdocs               			  Generates documentation$(COLOR_NONE)"
	@echo "$(COLOR_OK)  testacc               		  			  Run acceptence tests$(COLOR_NONE)"
	@echo "$(COLOR_OK)  pull-spec               				  Pull down the most recent released version of the spec$(COLOR_NONE)"

build:
	@echo "$(COLOR_LUMOS)Building...$(COLOR_NONE)"
	make pull-spec
	make provider-code-spec
	make provider
	make docs

pull-spec:
	@echo "$(COLOR_LUMOS)Pulling in latest spec...$(COLOR_NONE)"
	rm -f spec.json
	rm -rf spec-raw
	git clone https://github.com/teamlumos/lumos-openapi-spec spec-raw
	cp spec-raw/openapi.json spec.json
	rm -fr spec-raw

provider-code-spec:
	@echo "$(COLOR_LUMOS)Generating Provider Code Spec...$(COLOR_NONE)"
	tfplugingen-openapi generate --config generator_config.yml --output tmp-provider-code-spec.json spec.json
	python add_plan_modifiers.py > provider-code-spec.json
	rm -f tmp-provider-code-spec.json
	@echo "$(COLOR_OK)Done!$(COLOR_NONE)"

provider:
	@echo "$(COLOR_LUMOS)Generating Provider...$(COLOR_NONE)"
	rm -rf internal/generated
	tfplugingen-framework generate all --input provider-code-spec.json --output internal/generated
	@echo "$(COLOR_OK)Done!$(COLOR_NONE)"

sdocs:
	@echo "$(COLOR_LUMOS)Generating Documentation...$(COLOR_NONE)"
	rm -rf docs
	tfplugindocs generate
	@echo "$(COLOR_OK)Done!$(COLOR_NONE)"


testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m
