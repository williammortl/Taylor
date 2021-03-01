all: clean build

build: src/main.c
	mkdir -p bin
	cd src
	gcc main.c -o ../bin/main.exe
	cd ..

clean:
	rm -rf bin
