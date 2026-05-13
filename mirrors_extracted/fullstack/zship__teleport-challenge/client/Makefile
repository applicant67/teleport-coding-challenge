SHELL := /bin/sh -c
MAKEFLAGS += --no-builtin-rules
.SUFFIXES:
.DELETE_ON_ERROR:

export PATH := node_modules/.bin:$(PATH)

SRC_FILES = $(shell find src -type f)

.PHONY: default
default: build

build: build/index.html build/index.js build/reset.css build/style.css

build/index.html: src/index.html
	mkdir -p "$(@D)"
	cp "$^" "$@"

build/index.js: node_modules/.dirstamp $(SRC_FILES)
	mkdir -p "$(@D)"
	./node_modules/.bin/tsc
	webpack --config webpack.dev.config.ts
	touch "$@"

build/style.css: src/style.css
	mkdir -p "$(@D)"
	cp "$^" "$@"

build/reset.css:
	mkdir -p "$(@D)"
	curl -L https://meyerweb.com/eric/tools/css/reset/reset.css > "$@"

node_modules/.dirstamp: package.json
	yarn
	touch "$@"
