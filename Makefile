generate:
	mockgen -package=fus -destination=./fus/echo_context_mock.go "github.com/labstack/echo/v4" Context 
	go generate

test: generate
	go test -v ./...

test-only: 
	go test -v ./...

