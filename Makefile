export GAME = super.nes
export TEST = nestest.nes
export FILE = tests/cpu/nestest.nes

build:
	@GOOS=linux go build -x
run:build
	@./nesgo -game games/$(GAME) -scale 3 -interval 3ms
test:build
	@./nesgo -game tests/$(TEST)
trace:build
	@./nesgo -trace -game $(FILE)
disassemble:build
	@./nesgo -disassemble -game $(FILE)
color:build
	@./nesgo -game games/color.nes
win:
	@ CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc \
		GOOS=windows GOARCH=amd64 \
		CGO_LDFLAGS="-lmingw32 -lSDL2" \
		CGO_CFLAGS="-D_REENTRANT -O2" go build -x
