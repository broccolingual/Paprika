.PHONY: clean install run build

clean:
	rm -rf ./bin/*

install:
	# install 3rd party library
	# go install ~
	# go mod download
	go mod tidy

run:
	go run main.go window.go keyboard.go editor.go term.go syntax.go

build: install clean
	GOOS=linux go build -ldflags="-s -w -buildid=" -trimpath -o bin/main main.go window.go keyboard.go editor.go term.go syntax.go
