go test -coverprofile=coverage.out.tmp ./...
grep -v "domain/mocks" coverage.out.tmp > coverage.out
go tool cover -func=coverage.out
go tool cover --html=coverage.out
