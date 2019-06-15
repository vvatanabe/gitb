NAME = gitb
PKG = github.com/vvatanabe/gitb
VERSION = $(shell gobump show -r .)
COMMIT = $$(git describe --tags --always)
DATE = $$(date '+%Y-%m-%d_%H:%M:%S')
BUILD_LDFLAGS = -X main.commit=$(COMMIT) -X main.date=$(DATE)
RELEASE_BUILD_LDFLAGS = -s -w $(BUILD_LDFLAGS)

ifeq ($(update),yes)
  u=-u
endif

.PHONY: devel-deps
devel-deps:
	export GO111MODULE=off && \
	go get ${u} github.com/mattn/goveralls && \
	go get ${u} golang.org/x/lint/golint && \
	go get ${u} github.com/motemen/gobump/cmd/gobump && \
	go get ${u} github.com/Songmu/ghch/cmd/ghch && \
	go get ${u} github.com/Songmu/goxz/cmd/goxz && \
	go get ${u} github.com/tcnksm/ghr

.PHONY: cover
cover: devel-deps
	export GO111MODULE=on && \
	goveralls -coverprofile=coverage.out -service=travis-ci

.PHONY: lint
lint: devel-deps
	export GO111MODULE=on && \
	go vet ./... && \
	golint -set_exit_status ./...

.PHONY: test
test:
	export GO111MODULE=on && \
	go test -v -race -covermode=atomic -coverprofile=coverage.out ./...

.PHONY: bump
bump: devel-deps
	export GO111MODULE=on && \
	./bump

.PHONY: build
build:
	export GO111MODULE=on && \
	go build -ldflags="$(BUILD_LDFLAGS)" -o ./dist/current/$(NAME) .

.PHONY: install
install:
	export GO111MODULE=on && \
	go install -ldflags="$(BUILD_LDFLAGS)" .

.PHONY: crossbuild
crossbuild: devel-deps
	export GO111MODULE=on && \
	goxz -pv=$(VERSION) -arch=386,amd64 -build-ldflags="$(RELEASE_BUILD_LDFLAGS)" \
	  -o=$(NAME) -d=./dist/$(VERSION) .

.PHONY: upload
upload: devel-deps
	ghr -username vvatanabe -replace $(VERSION) ./dist/$(VERSION)

.PHONY: release
release: bump crossbuild upload