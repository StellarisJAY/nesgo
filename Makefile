all: win linux
win: web-win sdl-win
linux:web-linux sdl-linux
web-win:
	@ GOOS=windows go build -x -tags="web" -o nesgo_web.exe
web-linux:
	@ go build -x -tags="web" -o nesgo_web
sdl-win:
	@ CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc \
		GOOS=windows GOARCH=amd64 \
		CGO_LDFLAGS="-lmingw32 -lSDL2" \
		CGO_CFLAGS="-D_REENTRANT -O2" go build -tags="sdl" -x -o nesgo_sdl.exe
sdl-linux:
	@ go build -tags="sdl" -x -o nesgo_sdl
