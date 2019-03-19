default: build
build:
	go build -o terraform-provider-bouncr


test_assignments: testdeps
	BOUNCR_ACCOUNT=admin BOUNCR_PASSWORD=password TF_ACC=1 \
	go test -v -run TestBouncrAssignments_Basic

test_application: testdeps
	BOUNCR_ACCOUNT=admin BOUNCR_PASSWORD=password TF_ACC=1 \
	go test -v -run TestBouncrApplication_Basic

test: testdeps
	go test -v ./...

testdeps:
	go get -d -v -t ./...
	go get golang.org/x/lint/golint \
		golang.org/x/tools/cmd/cover \
		github.com/axw/gocov/gocov \
		github.com/mattn/goveralls
