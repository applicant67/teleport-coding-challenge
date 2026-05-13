SHELL := /bin/sh -c
MAKEFLAGS += --no-builtin-rules
.SUFFIXES:
.DELETE_ON_ERROR:

.PHONY: default
default: \
	.keys/server-key.pem \
	.keys/server-cert.pem \
	node_modules/.dirstamp \
	.env

.keys/server-key.pem .keys/server-cert.pem &:
	scripts/generate-keys .keys

node_modules/.dirstamp: package.json
	yarn
	touch "$@"

.env:
	cp template.env "$@"
