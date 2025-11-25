# Increment patch version and create tag/release
# Usage:
# make release TYPE=patch   (default)
# make release TYPE=minor
# make release TYPE=major

REPO := c0mm4nd/tronetl
TYPE ?= patch
PREFIX ?= v

current_tag := $(shell git tag --list 'v*' --sort=-v:refname | head -n1)
ifdef current_tag
base_version := $(subst $(PREFIX),,$(current_tag))
else
base_version := 0.0.0
endif
major := $(word 1,$(subst ., ,$(base_version)))
minor := $(word 2,$(subst ., ,$(base_version)))
patch := $(word 3,$(subst ., ,$(base_version)))

ifeq ($(TYPE),major)
  major := $(shell echo $$((major + 1)))
  minor := 0
  patch := 0
else ifeq ($(TYPE),minor)
  minor := $(shell echo $$((minor + 1)))
  patch := 0
else
  patch := $(shell echo $$((patch + 1)))
endif

new_version := $(major).$(minor).$(patch)
new_tag := $(PREFIX)$(new_version)

.PHONY: release show-version

show-version:
	@echo Current tag: $(current_tag)
	@echo New tag: $(new_tag)

release: show-version
	@if git rev-parse $(new_tag) >/dev/null 2>&1; then echo "Tag $(new_tag) already exists" && exit 1; fi
	git tag -a $(new_tag) -m "Release $(new_tag)"
	git push origin $(new_tag)
	@echo "Pushed tag $(new_tag). GitHub Action will create the release."
