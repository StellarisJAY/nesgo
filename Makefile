sdl-win: # 编译windows版sdl本地客户端
	@ CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc \
		GOOS=windows GOARCH=amd64 \
		CGO_LDFLAGS="-lmingw32 -lSDL2" \
		CGO_CFLAGS="-D_REENTRANT -O2" go build -tags="sdl" -x -o nesgo_sdl.exe
sdl-linux: # 编译linux版sdl本地客户端
	@ go build -tags="sdl" -x -o nesgo_sdl
.PHONY: sdl-linux sdl-win
