export GAME = snake
build:
	@go build
run:build
	@./nesgo games/$(GAME)
