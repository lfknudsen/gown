all:
	go build -C cmd/gown/ -o ../../gown.exe
	go build -C cmd/queue/ -o ../../queue.exe
	install gown.exe queue.exe