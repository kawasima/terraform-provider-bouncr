VERSION=v0.1.0
ARCH=amd64
OS=linux

TARGET_BINARY=terraform-provider-bouncr_$(VERSION)

TERRAFORM_PLUGIN_DIR=$(HOME)/.terraform.d/plugins/$(OS)_$(ARCH)

.PHONY: $(TARGET_BINARY)

build: $(TARGET_BINARY)

$(TARGET_BINARY):
	go build -o $(TARGET_BINARY)


test_oidc_provider: testdeps
	BOUNCR_ACCOUNT=admin BOUNCR_PASSWORD=password TF_ACC=1 \
	go test -v ./... -run TestBouncrOidcProvider_

test: testdeps
	BOUNCR_ACCOUNT=admin BOUNCR_PASSWORD=password TF_ACC=1 \
	go test -v ./...

testdeps:
	go get -d -v -t ./...
	go get golang.org/x/lint/golint \
		golang.org/x/tools/cmd/cover \
		github.com/axw/gocov/gocov \

install_plugin_locally: $(TARGET_BINARY)
	mkdir -p $(TERRAFORM_PLUGIN_DIR)
	cp ./$(TARGET_BINARY) $(TERRAFORM_PLUGIN_DIR)/
