HELMFILE_DIRS := $(shell find . -type f -name 'helmfile.yaml' ! -path './helmfile.yaml' ! -path './common/*' -exec dirname {} \; | sed 's|^\./||')

.PHONY: help
help:
	@echo "Usage: make <target> [directory...]"
	@echo ""
	@echo "Targets:"
	@echo "  lint    : Lint helmfile.yaml files"
	@echo "  diff    : Show the difference between the last run and the current state"
	@echo "  apply   : Apply the changes to the cluster"
	@echo "  destroy : Destroy the resources in the cluster"
	@echo ""
	@echo "Examples:"
	@echo "  make lint                  : Lint all directories"
	@echo "  make diff redis            : Show diff for redis directory"
	@echo "  make apply redis mysql     : Apply changes to redis and mysql directories"
	@echo "  make sync nginx            : Sync resources in nginx directory"
	@echo "  make destroy nginx         : Destroy resources in nginx directory"
	@echo ""
	@echo "Note: If no directory is specified, the command will be executed for all directories."

# Function to retry a command up to 3 times
define retry_command
	echo "Running command: $1"; \
	for i in 1 2 3; do \
		if $1; then \
			exit 0; \
		fi; \
		echo "Command failed. Retry $$i/3"; \
		if [ $$i -eq 3 ]; then \
			exit 1; \
		fi; \
	done
endef

# Function to run helmfile commands with retry logic
define run_helmfile_command
	@if [ -n "$(filter-out $1,$(MAKECMDGOALS))" ]; then \
		for dir in $(filter-out $1,$(MAKECMDGOALS)); do \
			if [ -d "$$dir" ]; then \
				echo "$2 $$dir"; \
				cd "$$dir" && \
				if ! $(call retry_command,helmfile $1); then \
					echo "Command failed after retries in $$dir"; \
					return 1; \
				fi; \
				cd - > /dev/null; \
			else \
				echo "Directory $$dir not found"; \
				return 1; \
			fi; \
		done; \
	else \
		echo "$2 all directories"; \
		for dir in $(HELMFILE_DIRS); do \
			echo "Processing: $$dir"; \
			cd "$$dir"; \
			if ! @bash -c '$(call retry_command,helmfile $1)'; then \
				echo "Command failed after retries in $$dir"; \
			fi; \
			cd - > /dev/null; \
		done; \
	fi
endef

.PHONY: list
list:
	@for dir in $(HELMFILE_DIRS); do \
		echo "  $$dir"; \
	done

.PHONY: lint
lint:
	$(call run_helmfile_command,$@,Linting)

.PHONY: diff
diff:
	$(call run_helmfile_command,$@,Diffing)

.PHONY: apply
apply:
	$(call run_helmfile_command,$@,Applying)


.PHONY: sync
sync:
	$(call run_helmfile_command,$@,Syncing)

.PHONY: destroy
destroy:
	$(call run_helmfile_command,$@,Destroying)

.PHONY: fetch
fetch:
	@if [ -n "$(filter-out $@,$(MAKECMDGOALS))" ]; then \
		for dir in $(filter-out $@,$(MAKECMDGOALS)); do \
			if [ -d "$$dir" ]; then \
				echo "Fetching Helm charts for $$dir"; \
				cd "$$dir" && helmfile fetch --output-dir $(CURDIR)/tmp || exit 1; \
				cd - > /dev/null; \
			else \
				echo "Directory $$dir not found"; \
				exit 1; \
			fi; \
		done; \
	else \
		echo "Fetching Helm charts for all directories..."; \
		for dir in $(HELMFILE_DIRS); do \
			echo "Processing: $$dir"; \
			cd "$$dir" && helmfile fetch --output-dir $(CURDIR)/tmp || exit 1; \
			cd - > /dev/null; \
		done; \
	fi
	@echo "All Helm charts have been fetched to $(CURDIR)/tmp directory"

%:
	@:
