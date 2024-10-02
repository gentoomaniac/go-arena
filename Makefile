phony: build-bots clean-bots build

build-bots:
	make -C pkg/bots/gentoobot build
	make -C pkg/bots/testbot build

clean-bots:
	make -C pkg/bots/gentoobot clean
	make -C pkg/bots/testbot clean

pkg/bots/gentoobot/gentoobot.so:
	make -C pkg/bots/gentoobot build

pkg/bots/testbot/testbot.so:
	make -C pkg/bots/testbot build

run: pkg/bots/gentoobot/gentoobot.so pkg/bots/testbot/testbot.so
	go run cmd/go-arena/main.go -vvv --map-path="maps/test.tmx" -b pkg/bots/testbot/testbot.so -b pkg/bots/gentoobot/gentoobot.so -b pkg/bots/testbot/testbot.so -b pkg/bots/gentoobot/gentoobot.so