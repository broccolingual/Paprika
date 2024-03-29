.PHONY: clean install run build

clean:
	rm -rf ./bin/*

install:
	# install 3rd party library
	go get golang.org/x/sys/unix
	go get github.com/pkg/term/termios
	go mod tidy

run:
	go run main.go view.go event.go keyboard.go editor.go syntax.go

build: install clean
	GOOS=linux go build -ldflags="-s -w -buildid=" -trimpath -o bin/paprika main.go view.go event.go keyboard.go editor.go syntax.go
