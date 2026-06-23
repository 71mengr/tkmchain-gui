package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

// loadImage loads an image resource from the given path and returns a *canvas.Image.
// If width or height are >0 they are used to set a minimum size for layout hints.
// The path is relative to the current working directory (the executable's working dir).
func loadImage(path string, width, height int) *canvas.Image {
	res, err := fyne.LoadResourceFromPath(path)
	if err != nil {
		return nil
	}
	img := canvas.NewImageFromResource(res)
	img.FillMode = canvas.ImageFillContain
	if width > 0 && height > 0 {
		img.SetMinSize(fyne.NewSize(float32(width), float32(height)))
	}
	return img
}
