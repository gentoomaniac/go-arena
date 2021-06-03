package ui

import (
	"bufio"
	"fmt"
	"image"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	FilenameFormatString = "ui/textures/alphabet-%s.png"
	ValidationRegex      = regexp.MustCompile("[^a-z0-9 !?$]+")
	SpaceWidth           = 40
	charset              = map[rune]*ebiten.Image{}
)

func loadRune(r rune, name string) {
	f, err := os.Open(fmt.Sprintf(FilenameFormatString, name))
	if err != nil {
		log.Panicf("could not find file for '%s'", r)
	}

	img, _, err := image.Decode(bufio.NewReader(f))
	if err != nil {
		log.Panicf("could not load rune image for '%s'", r)
	}

	charset[r] = ebiten.NewImageFromImage(img)
}

func init() {
	for c := 'a'; c <= 'z'; c++ {
		loadRune(c, string(c))
	}
	for c := '0'; c <= '9'; c++ {
		loadRune(c, string(c))
	}
	loadRune('!', "symbol-exclamation")
	loadRune('?', "symbol-question")
	loadRune('$', "symbol-dollar")
}

func NewText(s string) *Text {
	text := strings.ToLower(s)
	text = string(ValidationRegex.ReplaceAll([]byte(text), []byte("")))

	t := &Text{text: text}
	t.Image(true)

	return t
}

type Text struct {
	text  string
	image *ebiten.Image
}

func (t *Text) Image(refresh bool) *ebiten.Image {
	if t.image == nil || refresh {
		width := 0
		for _, c := range t.text {
			if c == ' ' {
				width += SpaceWidth
			} else {
				width += charset[c].Bounds().Dx()
			}
		}
		t.image = ebiten.NewImage(width, charset['a'].Bounds().Dy())

		opt := ebiten.DrawImageOptions{}
		for _, c := range t.text {
			if c == ' ' {
				opt.GeoM.Translate(float64(SpaceWidth), 0)
			} else {
				t.image.DrawImage(charset[c], &opt)
				opt.GeoM.Translate(float64(charset[c].Bounds().Dx()), 0)
			}
		}
	}

	return t.image
}
