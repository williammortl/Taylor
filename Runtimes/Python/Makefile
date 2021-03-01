all: clean debug

debug: src/main.c
	mkdir -p bin
	cd src && gcc main.c -g -o ../bin/main.exe

release: src/main.c
	mkdir -p bin
	cd src && gcc main.c -o ../bin/main.exe

clean:
	rm -rf bin
