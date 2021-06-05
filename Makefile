build:
	make -C bots/gentoobot build
	make -C bots/testbot build

clean:
	make -C bots/gentoobot clean
	make -C bots/testbot clean
