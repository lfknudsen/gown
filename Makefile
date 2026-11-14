all: proto install

compile:
	go build -C cmd/gown -o ../..
	go build -C cmd/queue -o ../..
	go build -C cmd/gown-srv -o ../..
	go build -C cmd/gown-client -o ../..
	go build -C cmd/gown-worker -o ../..

install: compile
	install gown.exe queue.exe gown-srv.exe gown-client.exe gown-worker.exe

gnu: compile
	install gown queue gown-srv gown-client gown-worker

go-install:
	go install -C cmd/gown
	go install -C cmd/queue
	go install -C cmd/gown-srv
	go install -C cmd/gown-client

protoc: proto/listen/listen.proto
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/listen/listen.proto --descriptor_set_out=docs/descriptor.pb --include_source_info

doc: docs/descriptor.pb
	cd docs; sabledocs

tiara: proto/listen/listen.proto
	protoc --dart_out="grpc:."  proto/listen.proto --descriptor_set_out=docs/descriptor.pb --include_source_info
	rm proto/listen/listen.pbserver.dart
	mv proto/listen/*.dart ../tiara/lib/proto
	cp docs/descriptor.pb ../tiara/docs/descriptor.pb
	cd ../tiara/docs; sabledocs
