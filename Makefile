.PHONY: gen

gen:
	cd proto
	protoc -I /usr/local/include \
	-I vendor/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
	-I ./proto \
	--go_out=plugins=grpc:. \
	--grpc-gateway_out=logtostderr=true:. \
	--swagger_out=logtostderr=true:. \
	profile.proto

run:
	GRPC_GO_LOG_VERBOSITY_LEVEL=99 GRPC_GO_LOG_SEVERITY_LEVEL=info go run \
	-ldflags "-X github.com/mkorolyov/core/discovery/consul.env=dev" \
	-ldflags "-X github.com/mkorolyov/core/discovery/consul.ip=0.0.0.0" \
	-ldflags "-X github.com/mkorolyov/core/discovery/consul.port=9090" \
	cmd/profiles/main.go