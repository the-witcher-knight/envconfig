test:
	go test -v -coverpkg=./... -coverprofile=profile.cov ./...

benchmark:
	go test -bench=.

build:
	go build .

