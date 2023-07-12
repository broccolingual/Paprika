.PHONY: clean install run build

clean:
	rm -rf ./bin/*

install:
	# install 3rd party library
	# go install ~
	# go mod download
	go mod tidy

run:
	go run main.go view.go event.go keyboard.go editor.go syntax.go

build: install clean
	GOOS=linux go build -ldflags="-s -w -buildid=" -trimpath -o bin/xanadu main.go view.go event.go keyboard.go editor.go syntax.go
