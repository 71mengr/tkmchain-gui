package ui

import (
	_ "embed"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

// This file embeds image assets from the repository's images/ folder and exposes them
// as fyne.Resources and convenience helpers. Embedding makes the assets available
// inside the compiled binary so the UI doesn't depend on external files at runtime.

//go:embed assets/images/icon-16.png
var icon16Bytes []byte

//go:embed assets/images/icon-32.png
var icon32Bytes []byte

//go:embed assets/images/icon-48.png
var icon48Bytes []byte

//go:embed assets/images/icon-64.png
var icon64Bytes []byte

//go:embed assets/images/icon-128.png
var icon128Bytes []byte

//go:embed assets/images/icon-256.png
var icon256Bytes []byte

//go:embed assets/images/icon.png
var iconPNGBytes []byte

//go:embed assets/images/icon.ico
var iconIcoBytes []byte

var (
	// Public resources (use these directly in widget.NewButtonWithIcon, app.SetIcon, etc.)
	Icon16  = fyne.NewStaticResource("icon-16.png", icon16Bytes)
	Icon32  = fyne.NewStaticResource("icon-32.png", icon32Bytes)
	Icon48  = fyne.NewStaticResource("icon-48.png", icon48Bytes)
	Icon64  = fyne.NewStaticResource("icon-64.png", icon64Bytes)
	Icon128 = fyne.NewStaticResource("icon-128.png", icon128Bytes)
	Icon256 = fyne.NewStaticResource("icon-256.png", icon256Bytes)
	IconPNG  = fyne.NewStaticResource("icon.png", iconPNGBytes)
	IconICO  = fyne.NewStaticResource("icon.ico", iconIcoBytes)
)

// ImageFromResource returns a *canvas.Image configured for common use in the UI.
// Pass desired minimum width/height (0 to leave unconstrained).
func ImageFromResource(res fyne.Resource, w, h int) *canvas.Image {
	if res == nil {
		return nil
	}
	img := canvas.NewImageFromResource(res)
	img.FillMode = canvas.ImageFillContain
	if w > 0 && h > 0 {
		img.SetMinSize(fyne.NewSize(float32(w), float32(h)))
	}
	return img
}

// Fallback helper: returns a short description of the embedded resources for debugging.
func EmbeddedResourcesInfo() string {
	return fmt.Sprintf("embedded: %d resources (16/32/48/64/128/256/png/ico)", 8)
}
