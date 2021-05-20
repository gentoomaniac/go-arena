# go-arena

_***This is a wip project!***_

## What?
go-arena is a [Robocode](https://robocode.sourceforge.io/) clone written in go and [ebiten](https://ebiten.org/)

You create code that controls your tank.
Last tank to survive is the winner.

![screenshot The arena](img/screenshot.png)

## Let them fight

### compile the robots

    # in the main directory
    make build

### start the battle

Specify the compiled bots with `-b` parameter:

    go run . -b bots/testbot/testbot.so -b bots/gentoobot/gentoobot.so

## write your own bot

Check out the code for [TestBot](bots/testbot.go).

To compile run:

    go build -buildmode=plugin -o newbot.so newbot.go

## State

So for this project is very rudimental but making progress.

Bots can't do much more than driving around and crash into walls (and "die" by that.)
