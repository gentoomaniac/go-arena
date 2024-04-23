phony: build-bots clean-bots build

build-bots:
	make -C pkg/bots/gentoobot build
	make -C pkg/bots/testbot build

clean-bots:
	make -C pkg/bots/gentoobot clean
	make -C pkg/bots/testbot clean
