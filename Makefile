t="coverage.txt"

test:
	go test ./... -cover

coverage:
	go test -coverprofile=$t ./... && go tool cover -html=$t && unlink $t
