# SPDX-FileCopyrightText: 2024 Dominik Wombacher <dominik@wombacher.cc>
# SPDX-FileCopyrightText: 2019 The SourceHut API Contributors
#
# SPDX-License-Identifier: BSD-2-Clause

.POSIX:
.SUFFIXES:

# Standardize on NetBSD style builtins.
.CURDIR ?= $(CURDIR)
.CURDIR ?= $(PWD)

GOFILES!=find . -name '*.go'

GO=go
TAGS=
VERSION!=git describe --tags --dirty 2>/dev/null
COMMIT!=git rev-parse --short HEAD 2>/dev/null

GOLDFLAGS =-s -w
GOLDFLAGS+=-X main.commit=$(COMMIT)
GOLDFLAGS+=-X main.version=$(VERSION)
GOLDFLAGS+=-extldflags $(LDFLAGS)
GCFLAGS  = all=-trimpath=$(.CURDIR)
ASMFLAGS = all=-trimpath=$(.CURDIR)

srht: go.mod
	$(GO) build \
		-trimpath \
		-gcflags="$(GCFLAGS)" \
		-asmflags="$(ASMFLAGS)" \
		-tags "$(TAGS)" \
		-o $@ \
		-ldflags "$(GOLDFLAGS)" \
		./cmd/srht/

clean:
	rm srht

test-cmd:
	@echo Test: cmd/srht
	cd cmd/srht; go test -cover -fullpath -v

test-testlog:
	@echo Test: internal/testlog
	cd internal/testlog; go test -cover -fullpath -v

test-sourcehut:
	@echo Test: sourcehut
	go test -cover -fullpath -v

test: test-sourcehut test-cmd test-testlog

bump:
	@echo bump go dependencies and module versions
	go get -u ./...
	go mod tidy
