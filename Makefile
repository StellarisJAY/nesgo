export GAME = PacMan.nes
export FUNC = TestPPU_RenderTile
export TRACE = false
build:
	@go build
run:build
	@./nesgo -game games/$(GAME) -trace $(TRACE)
test:build
	@./nesgo -game games/nestest.nes -trace true
gotest:build
	@go test -v -run $(FUNC) ./
