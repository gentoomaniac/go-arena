package gfx

import (
	"fmt"
	"image/gif"
	"io"

	"github.com/hajimehoshi/ebiten"
)

type Animation struct {
	currentFrameIndex int
	Frames            []*ebiten.Image
	AnimationSpeed    int
	FrameClock        int
	Width             int
	Height            int
}

func (a Animation) GetCurrentFrame() *ebiten.Image {
	return a.Frames[a.currentFrameIndex]
}
func (a *Animation) GetNextFrame() *ebiten.Image {
	if a.currentFrameIndex < len(a.Frames)-1 {
		a.currentFrameIndex++
	} else {
		a.Reset()
	}
	return a.Frames[a.currentFrameIndex]
}
func (a *Animation) Reset() {
	a.currentFrameIndex = 0
}

func (a *Animation) GetFrame() *ebiten.Image {
	var frame *ebiten.Image

	if a.FrameClock%a.AnimationSpeed == 0 {
		frame = a.GetNextFrame()
	} else {
		frame = a.GetCurrentFrame()
	}

	if a.FrameClock > 60 {
		a.FrameClock = 0
	} else {
		a.FrameClock++
	}

	return frame
}

func (a Animation) GetDuration() int {
	return len(a.Frames) * a.AnimationSpeed
}

func AnimationFromGIF(r io.Reader) (*Animation, error) {
	a := &Animation{}

	gif, err := gif.DecodeAll(r)
	if err != nil {
		return nil, err
	}
	for _, image := range gif.Image {
		frame, err := ebiten.NewImageFromImage(image, ebiten.FilterDefault)
		if err != nil {
			return nil, fmt.Errorf("couldn't load image from gif frame: %s", err)
		}
		a.Frames = append(a.Frames, frame)
	}

	a.Width, a.Height = a.Frames[0].Size()

	return a, nil
}
