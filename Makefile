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

compile:
	go build -o bin/profiles \
	-ldflags "-X github.com/mkorolyov/core/discovery/consul.env=dev -X github.com/mkorolyov/core/discovery/consul.ip=127.0.0.1 -X github.com/mkorolyov/core/discovery/consul.port=9090" \
	cmd/profiles/main.go

run: compile
	GRPC_GO_LOG_VERBOSITY_LEVEL=99 GRPC_GO_LOG_SEVERITY_LEVEL=info bin/profiles

docker:
	GOOS=linux GOARCH=arm go build -v -o bin/linux_profiles \
	-ldflags "-X github.com/mkorolyov/core/discovery/consul.env=dev -X github.com/mkorolyov/core/discovery/consul.ip=0.0.0.0 -X github.com/mkorolyov/core/discovery/consul.port=9090"
	docker build -t profiles .

docker_run: docker
	docker run -m 1024m --memory-swap=1024m --cpus="1" -p 9090:9090 --rm profiles