package main

import (
	"fmt"
	"image"
	"log"

	"github.com/gentoomaniac/ebitmx"
	"github.com/hajimehoshi/ebiten"
)

const mapPath = "maps/test.tmx"
const (
	screenWidth  = 1080
	screenHeight = 720
)

func run() {
	scalingFactor := .4

	tmxMap, error := ebitmx.LoadFromFile(mapPath)
	if error != nil {
		log.Fatal(error)
	}

	startPosition := image.Point{tmxMap.PixelWidth / 2, 2400}

	tmxMap.CameraBounds = image.Rect(0, 0, screenWidth, screenHeight)
	tmxMap.CameraPosition = startPosition
	fmt.Printf("Map dimensions: %d/%d\n", tmxMap.PixelWidth, tmxMap.PixelHeight)

	game := NewGame().WithMap(tmxMap).WithScalingFactor(scalingFactor)

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Go Arena")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
