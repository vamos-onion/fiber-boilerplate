GO_CMDS := main
GO_MODULE := $(shell grep '^module ' $(DIR_ROOT)/go.mod | awk '{print $$2}')
GO_SQLC := $(subst $(DIR_SQLC_CONFIG)/,,$(wildcard $(DIR_SQLC_CONFIG)/*.go))
GO_SRCS := $(shell find $(DIR_ROOT) -name '*.go')
GO_TESTS := $(shell find $(DIR_ROOT) -name '*_test.go' -exec dirname {} \; | sort -u)
GO_STAMP_TIDY := $(DIR_OUT)/stamp-go-mod-tidy

GO_RUN_MODELS := $(foreach sqlc,$(GO_SQLC),run-$(sqlc))
GO_BUILD_TARGETS := $(foreach cmd,$(GO_CMDS),build-$(cmd))
GO_RUN_TARGETS := $(foreach cmd,$(GO_CMDS),run-$(cmd))

.PHONY: all build clean test

help-body::
	$(call HELP_HEADING,-,Go targets)
	$(HELP) "* all" "test & build"
	$(HELP) "* build" "build all local codes"
	$(foreach cmd,$(GO_CMDS),$(HELP) "  * build-$(cmd)" "build $(cmd)";)
	$(HELP) "* run" "use below targets instead"
	$(foreach cmd,$(GO_CMDS),$(HELP) "  * run-$(cmd)" "run $(cmd)";)
	$(HELP) "* test" "test packages"
	$(HELP) "* tidy" "add missing and remove unused modules"
	$(HELP) "* clean" "remove object files and cached files"
	echo

all: test build

sqlc: $(GO_RUN_MODELS)

build: $(GO_RUN_MODELS) $(GO_BUILD_TARGETS)

run:
	echo "* Use below targets instead"
	$(foreach cmd,$(GO_CMDS),$(HELP) "  * run-$(cmd)" "run $(cmd)")

test: tidy
	for dir in $(GO_TESTS); do \
		cd $$dir; \
		go test -v; \
	done

tidy: $(GO_STAMP_TIDY)

clean:
	go clean
	rm -f $(GO_STAMP_TIDY)

$(GO_RUN_MODELS):
	echo "Building sqlc file"; \
	go run $(patsubst run-%.go,%,$@)/$(subst run-,,$@); \
	sqlc generate -f $(DIR_OUT)/merged_sqlc.yaml

$(GO_BUILD_TARGETS): tidy
	export CMD=$(subst build-,,$@); \
	echo Building $(GO_MODULE); \
	go build -o $(DIR_OUT)/$(GO_MODULE) $(DIR_CMD)/$$CMD.go

$(GO_RUN_TARGETS): tidy
	export CMD=$(subst run-,,$@); \
	echo Running $$CMD; \
	go run $(DIR_CMD)/$$CMD.go

$(GO_STAMP_TIDY): $(GO_SRCS) openapi-generate
	go mod tidy
	touch $@