package basic

import (
	"fmt"
	"os"
	"strings"

	"github.com/mangofeet/netrunner-alt-gen/art"
	"github.com/mangofeet/nrdb-go"
	"github.com/tdewolff/canvas"
)

func (fb FrameBasic) Back() art.Drawer {

	return art.DrawerFunc(func(ctx *canvas.Context, card *nrdb.Printing) error {

		canvasWidth, canvasHeight := ctx.Size()

		attributionBoxTop := canvasHeight * 0.9409
		// attributionBoxBottom := canvasHeight * 0.032
		attributionBoxBottom := canvasHeight * 0.87
		attributionBoxHeight := attributionBoxTop - attributionBoxBottom
		attributionBoxLeft := canvasWidth * 0.25
		attributionBoxRight := canvasWidth * 0.75
		attributionBoxRadius := canvasWidth * 0.01

		cliBoxTop := canvasHeight * 0.13
		cliBoxBottom := canvasHeight * 0.0591
		cliBoxHeight := cliBoxTop - cliBoxBottom
		cliBoxLeft := canvasWidth * 0.1
		cliBoxRight := canvasWidth * 0.9
		cliBoxRadius := canvasWidth * 0.01

		fb.drawRoundedBox(ctx, attributionBoxTop, attributionBoxBottom, attributionBoxLeft, attributionBoxRight, attributionBoxRadius)

		fb.drawRoundedBox(ctx, cliBoxTop, cliBoxBottom, cliBoxLeft, cliBoxRight, cliBoxRadius)

		attributionFontSize := attributionBoxHeight * 0.6
		attributionTextMaxWidth := (attributionBoxRight - attributionBoxLeft) * 0.9
		attributionTextMaxHeight := (attributionBoxTop - attributionBoxBottom) * 0.9

		attributionString := fmt.Sprintf("%s<BR>Generated using \"%s\" by %s<BR>netrunner-alt-gen %s", card.Attributes.Title, fb.Algorithm, fb.Designer, fb.Version)

		if fb.Algorithm == "" && fb.Designer == "" {
			attributionString = fmt.Sprintf("%s<BR>Layout by netrunner-alt-gen %s", card.Attributes.Title, fb.Version)
		}
		if fb.Algorithm == "" && fb.Designer != "" {
			attributionString = fmt.Sprintf("%s<BR>Design by %s<BR>Layout by netrunner-alt-gen %s", card.Attributes.Title, fb.Designer, fb.Version)
		}

		attributionText := fb.getVerticalFittedText(ctx, attributionString, attributionFontSize, attributionTextMaxWidth, attributionTextMaxHeight, canvas.Center)

		attributionTextX := attributionBoxLeft + ((attributionBoxRight - attributionBoxLeft) * 0.05)
		attributionTextY := (attributionBoxTop - (attributionBoxHeight-attributionText.Bounds().H)*0.5)
		ctx.DrawText(attributionTextX, attributionTextY, attributionText)

		cliFontSize := cliBoxHeight * 0.6
		cliTextMaxWidth := (cliBoxRight - cliBoxLeft) * 0.9
		cliTextMaxHeight := (cliBoxTop - cliBoxBottom) * 0.9
		cliString := getCLIText()

		cliText := fb.getVerticalFittedText(ctx, cliString, cliFontSize, cliTextMaxWidth, cliTextMaxHeight, canvas.Center)

		cliTextX := cliBoxLeft + ((cliBoxRight - cliBoxLeft) * 0.05)
		cliTextY := (cliBoxTop - (cliBoxHeight-cliText.Bounds().H)*0.5)
		ctx.DrawText(cliTextX, cliTextY, cliText)

		return nil
	})
}

func (fb FrameBasic) drawRoundedBox(ctx *canvas.Context, top, bottom, left, right, radius float64) {

	strokeWidth := getStrokeWidth(ctx)

	ctx.Push()
	ctx.SetFillColor(fb.getColorBG())
	ctx.SetStrokeColor(fb.getColorBorder())
	ctx.SetStrokeWidth(strokeWidth)

	path := &canvas.Path{}
	path.MoveTo(left, top-radius)
	path.QuadTo(left, top, left+radius, top)
	path.LineTo(right-radius, top)
	path.QuadTo(right, top, right, top-radius)
	path.LineTo(right, bottom+radius)
	path.QuadTo(right, bottom, right-radius, bottom)
	path.LineTo(left+radius, bottom)
	path.QuadTo(left, bottom, left, bottom+radius)
	path.Close()

	ctx.DrawPath(0, 0, path)
	ctx.Pop()
}

func getCLIText() string {
	var args []string

	var isFlag bool
	for _, arg := range os.Args {

		if len(arg) < 1 {
			continue
		}

		if arg == "-o" {
			continue
		}

		if arg == "--make-back" {
			continue
		}

		if arg[0] == '-' {
			isFlag = true
		}

		if !isFlag {
			continue
		}

		args = append(args, arg)

		if arg[0] != '-' {
			isFlag = false
		}

	}

	return strings.Join(args, " ")
}
