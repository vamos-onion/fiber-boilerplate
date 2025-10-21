OPENAPI_INDEX := $(DIR_API)/index.yaml
OPENAPI_SPEC := $(DIR_OUT)/openapi.yaml
OPENAPI_SRCS := $(shell find $(DIR_API) -name '*.yaml')
OPENAPI_DIR := $(DIR_SOURCE)/generated/serviceapi
OPENAPI_STAMP_VALIDATE := $(DIR_OUT)/stamp-openapi-validate
OPENAPI_STAMP_BUNDLE := $(DIR_OUT)/stamp-openapi-bundle
OPENAPI_STAMP_GENERATE := $(DIR_OUT)/stamp-openapi-generate

.PHONY: openapi-validate openapi-bundle openapi-generate openapi-upgrade

help-body::
	$(call HELP_HEADING,-,OpenAPI targets)
	$(HELP) "* openapi" "openapi-validate & openapi-bundle & openapi-generate"
	$(HELP) "* openapi-validate" "validates index $(OPENAPI_INDEX)"
	$(HELP) "* openapi-bundle" "bundles all specs into $(OPENAPI_SPEC)"
	$(HELP) "* openapi-generate" "generate codes from spec $(OPENAPI_SPEC)"
	echo

openapi: openapi-generate

openapi-validate: $(OPENAPI_STAMP_VALIDATE)

openapi-bundle: openapi-validate $(OPENAPI_STAMP_BUNDLE)

openapi-generate: openapi-bundle $(OPENAPI_STAMP_GENERATE)

openapi-upgrade:
	$(DIR_ROOT)/scripts/openapi.sh upgrade

$(OPENAPI_STAMP_VALIDATE): $(OPENAPI_SRCS)
	$(DIR_ROOT)/scripts/openapi.sh validate $(OPENAPI_INDEX)
	touch $@

$(OPENAPI_STAMP_BUNDLE): $(OPENAPI_SRCS)
	$(DIR_ROOT)/scripts/openapi.sh bundle $(OPENAPI_INDEX) $(OPENAPI_SPEC)
	touch $@

$(OPENAPI_STAMP_GENERATE): $(OPENAPI_SRCS)
	$(DIR_ROOT)/scripts/openapi.sh generate $(OPENAPI_SPEC) $(OPENAPI_DIR)
	touch $@
