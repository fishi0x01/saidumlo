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

release: deps build verify
	echo -n "" > SHA256SUMS
	env GOOS=linux GOARCH=amd64 go build -o sdl ./src/*
	zip sdl_${version}_linux_amd64.zip sdl
	sha256sum sdl_${version}_linux_amd64.zip >> SHA256SUMS
	rm sdl
	env GOOS=darwin GOARCH=amd64 go build -o sdl ./src/*
	zip sdl_${version}_darwin_amd64.zip sdl
	sha256sum sdl_${version}_darwin_amd64.zip >> SHA256SUMS
	rm sdl

clean:
	rm -f test/prod-bar || true
	rm -f test/prod-foo || true
	rm -f test/vault_pid || true
