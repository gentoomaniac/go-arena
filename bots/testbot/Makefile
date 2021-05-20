name=testbot

${name}.so:
	go build -buildmode=plugin -o ${name}.so ${name}.go

build: ${name}.so

clean:
	rm ${name}.so