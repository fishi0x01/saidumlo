RELEASE_PLATFORMS := linux darwin
VERSION := `git describe --tags`

build:
	go build -ldflags "-X main.sdlVersion=`git describe --tags`" -o bin/sdl ./src/main/*

deps:
	go get gopkg.in/alecthomas/kingpin.v2
	go get github.com/fatih/color
	go get gopkg.in/yaml.v2
	go get github.com/mitchellh/gox

verify:
	./test/test.pre.sh
	./test/test.sh
	./test/test.post.sh

release: deps build verify
	echo -n "" > SHA256SUMS
	${GOPATH}/bin/gox -ldflags="-X main.sdlVersion=`git describe --tags`" -osarch="linux/amd64" -osarch="darwin/amd64" -output="sdl_${VERSION}_{{.OS}}_{{.Arch}}" ./src/main/
	for platform in $(RELEASE_PLATFORMS); do\
		mv sdl_${VERSION}_$${platform}_amd64 sdl; \
		zip sdl_${VERSION}_$${platform}_amd64.zip sdl; \
		sha256sum sdl_${VERSION}_$${platform}_amd64.zip >> SHA256SUMS; \
		rm sdl; \
	done

clean:
	rm -f test/prod-bar || true
	rm -f test/prod-foo || true
	rm -f test/vault_pid || true
