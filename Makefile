all:
	go build -ldflags "-s -w"

grpc:
	protoc --go_out=. --go_opt=paths=source_relative \
	--go-grpc_out=. --go-grpc_opt=paths=source_relative \
	internal/types/fr24/feed.proto
