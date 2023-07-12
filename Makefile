export GAME = snake.nes
build:
	@go build
run:build
	@./nesgo games/$(GAME)
