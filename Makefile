NAME = gitb
PKG = github.com/vvatanabe/gitb
VERSION = $(shell gobump show -r .)
COMMIT = $$(git describe --tags --always)
DATE = $$(date '+%Y-%m-%d_%H:%M:%S')
BUILD_LDFLAGS = -X main.commit=$(COMMIT) -X main.date=$(DATE)
RELEASE_BUILD_LDFLAGS = -s -w $(BUILD_LDFLAGS)
u := $(if $(update),-u)

export GO111MODULE=on

.PHONY: devel-deps
devel-deps:
	sh -c '\
	tmpdir=$$(mktemp -d); \
	cd $$tmpdir; \
	go get ${u} \
	  github.com/mattn/goveralls           \
	  golang.org/x/lint/golint             \
	  github.com/x-motemen/gobump/cmd/gobump \
	  github.com/Songmu/ghch/cmd/ghch      \
	  github.com/Songmu/goxz/cmd/goxz      \
	  github.com/tcnksm/ghr'

.PHONY: cover
cover: devel-deps
	goveralls -coverprofile=coverage.out -service=travis-ci

.PHONY: lint
lint: devel-deps
	go vet ./... && \
	golint -set_exit_status ./...

.PHONY: test
test:
	go test -v -race -covermode=atomic -coverprofile=coverage.out ./...

.PHONY: bump
bump: devel-deps
	./bump

.PHONY: build
build:
	CGO_ENABLED=0 go build -ldflags="$(BUILD_LDFLAGS)" -o ./dist/current/$(NAME) .

.PHONY: install
install:
	CGO_ENABLED=0 go install -ldflags="$(BUILD_LDFLAGS)" .

.PHONY: crossbuild
crossbuild: devel-deps
	goxz -pv=$(VERSION) -arch=386,amd64 -build-ldflags="$(RELEASE_BUILD_LDFLAGS)" \
	  -d=./dist/$(VERSION) .

.PHONY: upload
upload: devel-deps
	ghr -username vvatanabe -replace $(VERSION) ./dist/$(VERSION)

.PHONY: release
release: bump crossbuild upload