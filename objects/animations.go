package objects

import (
	"fmt"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

type Orientation int

const (
	Up    Orientation = 0
	Down              = 1
	Left              = 2
	Right             = 3
)

type Animation struct {
	currentFrameIndex int
	Frames            [][]*ebiten.Image
	AnimationSpeed    int
	FrameClock        int
	Width             int
	Height            int
}

func (a Animation) GetCurrentFrame(orientation Orientation) *ebiten.Image {
	return a.Frames[orientation][a.currentFrameIndex]
}
func (a *Animation) GetNextFrame(orientation Orientation) *ebiten.Image {
	if a.currentFrameIndex < len(a.Frames[orientation])-1 {
		a.currentFrameIndex++
	} else {
		a.Reset()
	}
	return a.Frames[orientation][a.currentFrameIndex]
}
func (a *Animation) Reset() {
	a.currentFrameIndex = 0
}

func (a *Animation) GetFrame(orientation Orientation) *ebiten.Image {
	var frame *ebiten.Image

	if a.FrameClock%a.AnimationSpeed == 0 {
		frame = a.GetNextFrame(orientation)
	} else {
		frame = a.GetCurrentFrame(orientation)
	}

	if a.FrameClock > 60 {
		a.FrameClock = 0
	} else {
		a.FrameClock++
	}

	return frame
}

func (a Animation) GetDuration() int {
	return len(a.Frames[Up]) * a.AnimationSpeed
}

// This should be solved in a better way. Maybe a Spritemap.
func (a *Animation) LoadFramesFromFiles(FramePathsUp []string, FramePathsDown []string, FramePathsLeft []string, FramePathsRight []string) error {
	a.Frames = make([][]*ebiten.Image, 4)

	for _, file := range FramePathsUp {
		frame, _, error := ebitenutil.NewImageFromFile(file, ebiten.FilterDefault)
		if error != nil {
			return fmt.Errorf("Couldn't load image from file: %s", file)
		}
		a.Frames[Up] = append(a.Frames[Up], frame)
	}

	for _, file := range FramePathsDown {
		frame, _, error := ebitenutil.NewImageFromFile(file, ebiten.FilterDefault)
		if error != nil {
			return fmt.Errorf("Couldn't load image from file: %s", file)
		}
		a.Frames[Down] = append(a.Frames[Down], frame)
	}

	for _, file := range FramePathsLeft {
		frame, _, error := ebitenutil.NewImageFromFile(file, ebiten.FilterDefault)
		if error != nil {
			return fmt.Errorf("Couldn't load image from file: %s", file)
		}
		a.Frames[Left] = append(a.Frames[Left], frame)
	}

	for _, file := range FramePathsRight {
		frame, _, error := ebitenutil.NewImageFromFile(file, ebiten.FilterDefault)
		if error != nil {
			return fmt.Errorf("Couldn't load image from file: %s", file)
		}

		a.Frames[Right] = append(a.Frames[Right], frame)
	}

	// assuming all images have the same dimensions
	a.Width, a.Height = a.Frames[Down][0].Size()

	return nil
}
