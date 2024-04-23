package gfx

import (
	_ "embed"
)

type AnimationType int

const (
	Fire AnimationType = iota
)

//go:embed fire_transparent.gif
var FireGif []byte
