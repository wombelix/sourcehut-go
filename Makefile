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

srht: cmd/srht/go.mod go.mod $(GOFILES)
	cd cmd/srht; \
	$(GO) build \
		-gcflags="$(GCFLAGS)" \
		-asmflags="$(ASMFLAGS)" \
		-tags "$(TAGS)" \
		-o ../../$@ \
		-ldflags "$(GOLDFLAGS)"
