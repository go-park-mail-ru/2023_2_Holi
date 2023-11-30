#go test -coverprofile=coverage.out.tmp -coverpkg=./...  ./...
#cat coverage.out.tmp | grep -v _mock.go | grep -v _easyjson.go | grep -v .pb.go | grep -v _grpc.go | grep -v domain/mocks/ > coverage.out
#go tool cover -func=coverage.out
#go tool cover --html=coverage.out
#user

go test -coverprofile=coverage.out.tmp ./...
grep -v "domain/mocks" coverage.out.tmp > coverage.out
go tool cover -func=coverage.out
go tool cover --html=coverage.out
