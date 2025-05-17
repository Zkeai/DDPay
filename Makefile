run-cron:
	go run ./cmd/cron

run-http:
	go run ./cmd/http
 # 启动 HTTP 服务并开启调试器（Delve）
debug-http:
	dlv debug ./cmd/http/main.go --headless --listen=:2345 --api-version=2 --accept-multiclient

run-grpc:
	go run ./cmd/grpc


protoc:
	protoc \
		--go_out=./internal/pb \
		--go_opt=paths=source_relative \
		--go-grpc_out=./internal/pb \
		--go-grpc_opt=paths=source_relative \
		--proto_path ../AIDog/internal/pb/ ../AIDog/internal/pb/*.proto