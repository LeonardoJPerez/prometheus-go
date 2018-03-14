BINARY=prometheusTelemetry
GO_META_LINTER := $(GOPATH)/bin/gometalinter.v2

$(GO_META_LINTER):
	go get -u gopkg.in/alecthomas/gometalinter.v2	 
	gometalinter.v2 --install &> /dev/null

# Installs the dependencies, runs the test suite, and builds the project for a Docker container, outputs the result to 'deploy/foghorn'
build: $(GO_META_LINTER)
	make install-dep
	make lint
	make just-build

# Installs the dependencies, runs the test suite, and then installs the project
install: $(GO_META_LINTER)
	make install-dep	 
	make just-install

# Installs all the dependencies for the project
install-dep:
	go get -u github.com/golang/dep/cmd/dep
	dep ensure

# Runs the gometalinter on the project, and then runs the test suite
test: lint
	make just-test

# Runs the gometalinter on the project, checking for issues
lint: $(GO_META_LINTER)
	test -z $(gofmt -s -l $GO_FILES)
	gometalinter.v2 ./... --deadline=60s --vendor

# Runs the test suite for the project
just-test:
	go list ./... | grep -v /vendor | xargs go test -p 1 -timeout=20s

# Builds the project for a Docker container, outputs the result to 'deploy/foghorn'
just-build:
	go build

# Installs the project
just-install:	
	go install