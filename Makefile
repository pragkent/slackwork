PACKAGES?=$$(go list ./... | grep -v vendor)
GOFMT_FILES?=$$(find . -name '*.go' | grep -v vendor)
GIT_COMMIT?=$$(git rev-parse --short HEAD)
GIT_DIRTY?=$$(test -n "`git status --porcelain`" && echo "+DIRTY" || true)

bin:
	go build -ldflags "-X main.GitCommit=${GIT_COMMIT}${GIT_DIRTY}" -o bin/slackwork

static:
	CGO_ENABLED=0 go build -ldflags "-X main.GitCommit=${GIT_COMMIT}${GIT_DIRTY}" -o bin/slackwork

test:
	go test $(PACKAGES)

testrace:
	go test -race $(PACKAGES)

testcover:
	go test -cover $(PACKAGES)

vet:
	go vet $(PACKAGES)

fmt:
	gofmt -w $(GOFMT_FILES)

.PHONY: bin package test testrace testcover vet fmt
