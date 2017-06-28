build:
	go build -o bin/sdl ./src/*

deps:
	go get gopkg.in/alecthomas/kingpin.v2
	go get github.com/fatih/color
	go get gopkg.in/yaml.v2

verify:
	./test/test.pre.sh
	./test/test.sh
	./test/test.post.sh
