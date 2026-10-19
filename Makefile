all: proto install

compile:
	go build -C cmd/gown
	go build -C cmd/queue
	go build -C cmd/gown-srv
	go build -C cmd/gown-client

install: compile
	install cmd/gown/gown.exe cmd/queue/queue.exe cmd/gown-srv/gown-srv.exe cmd/gown-client/gown-client.exe

gnu:
	go build -C cmd/gown -o ../../gown
	go build -C cmd/queue -o ../../queue
	go build -C cmd/gown-srv -o ../../gown-srv
	go build -C cmd/gown-client/gown-client -o ../../gown-client
	install gown queue gown-srv

go-install:
	go install -C cmd/gown
	go install -C cmd/queue
	go install -C cmd/gown-srv
	go install -C cmd/gown-client

protoc: proto/listen.proto
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/listen.proto
	protoc --dart_out="grpc:."  proto/listen.proto
	rm proto/listen.pbserver.dart
	mv proto/*.dart ../tiara/lib/proto
