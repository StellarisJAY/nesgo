export GAME = PacMan.nes
export FUNC = TestPPU_RenderTile
build:
	@go build
run:build
	@./nesgo -game games/$(GAME)
test:build
	@./nesgo -game games/nestest.nes
trace:build
	@./nesgo -game games/$(GAME) -trace
gotest:build
	@go test -v -run $(FUNC) ./
