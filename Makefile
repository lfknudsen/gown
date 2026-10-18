all: install

compile:
	go build -C cmd/gown
	go build -C cmd/queue

install: compile
	install cmd/gown/gown.exe cmd/queue/queue.exe

gnu:
	go build -C cmd/gown -o ../../gown
	go build -C cmd/queue -o ../../queue
	install gown queue

go-install:
	go install -C cmd/gown
	go install -C cmd/queue