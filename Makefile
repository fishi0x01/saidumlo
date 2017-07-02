RELEASE_PLATFORMS := linux darwin

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
	for platform in $(RELEASE_PLATFORMS); do\
		env GOOS=linux GOARCH=amd64 go build -o sdl ./src/*; \
		zip sdl_${version}_$${platform}_amd64.zip sdl; \
		sha256sum sdl_${version}_$${platform}_amd64.zip >> SHA256SUMS; \
		rm sdl; \
	done

clean:
	rm -f test/prod-bar || true
	rm -f test/prod-foo || true
	rm -f test/vault_pid || true
