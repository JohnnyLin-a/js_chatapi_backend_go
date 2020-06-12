go env -w GO111MODULE=on
go mod download
go build ./cmd/chatapisrv/main.go
go build ./cmd/dbsetup/dbsetup.go