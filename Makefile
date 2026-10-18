all: install

compile:
	go build -C cmd/gown
	go build -C cmd/queue
	go build -C cmd/gown-srv

install: compile
	install cmd/gown/gown.exe cmd/queue/queue.exe cmd/gown-srv/gown-srv.exe

gnu:
	go build -C cmd/gown -o ../../gown
	go build -C cmd/queue -o ../../queue
	go build -C cmd/gown-srv -o ../../gown-srv
	install gown queue gown-srv

go-install:
	go install -C cmd/gown
	go install -C cmd/queue
	go install -C cmd/gown-srv