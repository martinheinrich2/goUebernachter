# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'


# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

## tidy: format code and tidy modfile
.PHONY: tidy
tidy:
	go fmt ./...
	go mod tidy -v

## audit: run quality control checks
.PHONY: audit
audit:
	go mod verify
	go vet ./...
	go run honnef.co/go/tools/cmd/staticcheck@latest -checks=all,-ST1000,-U1000 ./...
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...
	go test -race -buildvcs -vet=off ./...


# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

## test: run all tests
.PHONY: test
test:
	go test -v -race -buildvcs ./...

## test/cover: run all tests and display coverage
.PHONY: test/cover
test/cover:
	go test -v -race -buildvcs -coverprofile=/tmp/coverage.out ./...
	go tool cover -html=/tmp/coverage.out

## build: build the cmd/web application
.PHONY: build
build:
	## go build -o=/tmp/bin/web ./cmd/web
	go build -o=/tmp/bin/web ./cmd/web

## run: run the cmd/web application
.PHONY: run
run: build
	/tmp/bin/web

## run/live: run the application with reloading on file changes
.PHONY: run/live
run/live:
	go run github.com/cosmtrek/air@v1.52.1 \
		--build.cmd "make build" --build.bin "/tmp/bin/web" --build.delay "100" \
		--build.exclude_dir "" \
		--build.include_ext "go, tpl, tmpl, html, css, scss, js, ts, sql, jpeg, jpg, gif, png, bmp, svg, webp, ico" \
		--misc.clean_on_exit "true"


# ==================================================================================== #
# SQL MIGRATIONS
# ==================================================================================== #

## migrations/new name=$1: create a new database migration
#.PHONY: migrations/new
#migrations/new:
#	GOOSE_DRIVER=sqlite3 GOOSE_DBSTRING=./reception.db goose -dir internal/migrations create articles sql

## migrations/up: apply all up database migrations
.PHONY: migrations/up
migrations/up:
	GOOSE_DRIVER=sqlite3 GOOSE_DBSTRING=./sqlite.db goose -dir internal/migrations up

## migrations/down: apply all down database migrations
.PHONY: migrations/down
migrations/down:
	GOOSE_DRIVER=sqlite3 GOOSE_DBSTRING=./sqlite.db goose -dir internal/migrations down

## migrations/version: print the current in-use migration version
.PHONY: migrations/version
migrations/version:
	GOOSE_DRIVER=sqlite3 GOOSE_DBSTRING=./sqlite.db goose -dir internal/migrations version

## migrations/status: print the current migrations status
.PHONY: migrations/status
migrations/status:
	GOOSE_DRIVER=sqlite3 GOOSE_DBSTRING=./sqlite.db goose -dir internal/migrations status

## migrations/create: create a new database migration
.PHONY: migrations/create
migrations/create:
	GOOSE_DRIVER=sqlite3 GOOSE_DBSTRING=./sqlite.db goose -dir internal/migrations create ${name} sql