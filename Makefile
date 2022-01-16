all: clean createdirs build

gpython: src/gpython/main.go
	cd src/gpython; \
		go build -buildmode=plugin -o gpython.so
	cp src/gpython/gpython.so bin/.
	rm -rf src/gpython/gpython.so

starlightd: src/starlightd/main.go
	cd src/starlightd; \
		go build
	cp src/starlightd/starlightd bin/.
	rm -rf src/starlightd/starlightd

taylor: src/taylor/main.go
	cd src/taylor; \
		go build
	cp src/taylor/taylor bin/.
	rm -rf src/taylor/taylor

clean:
	rm -rf bin

createdirs:
	mkdir bin

build: gpython starlightd taylor
	ls bin -la
