package basic

import (
	"fmt"
	"image/color"
	"math"

	"github.com/mangofeet/netrunner-alt-gen/art"
	"github.com/mangofeet/nrdb-go"
	"github.com/tdewolff/canvas"
)

type FrameAgenda struct{}

func (FrameAgenda) Draw(ctx *canvas.Context, card *nrdb.Printing) error {

	canvasWidth, canvasHeight := ctx.Size()

	strokeWidth := getStrokeWidth(ctx)

	factionBaseColor := art.GetFactionBaseColor(card.Attributes.FactionID)
	factionColor := color.RGBA{
		R: uint8(math.Max(0, math.Min(float64(int64(factionBaseColor.R)-48), 255))),
		G: uint8(math.Max(0, math.Min(float64(int64(factionBaseColor.G)-48), 255))),
		B: uint8(math.Max(0, math.Min(float64(int64(factionBaseColor.B)-48), 255))),
		A: 0xff,
	}

	ctx.Push()
	ctx.SetFillColor(bgColor)
	ctx.SetStrokeColor(textColor)
	ctx.SetStrokeWidth(strokeWidth)

	costContainerR := getCostContainerRadius(ctx)

	titleBoxHeight := getTitleBoxHeight(ctx)

	titleBoxTop := getTitleBoxTop(ctx)
	titleBoxBottom := titleBoxTop - titleBoxHeight
	titleBoxRight := canvasWidth - costContainerR*3.25
	costContainerStart := titleBoxRight

	titlePath := &canvas.Path{}
	titlePath.MoveTo(0, titleBoxTop)
	titlePath.LineTo(titleBoxRight, titleBoxTop)
	titlePath.QuadTo(titleBoxRight-(costContainerR*0.5), titleBoxBottom+(titleBoxHeight*0.5), titleBoxRight, titleBoxBottom)
	titlePath.LineTo(0, titleBoxBottom)
	titlePath.Close()

	ctx.DrawPath(0, 0, titlePath)
	ctx.Pop()

	// outline for cost circle
	ctx.Push()
	ctx.SetFillColor(bgColor)
	ctx.SetStrokeColor(textColor)
	ctx.SetStrokeWidth(strokeWidth)

	path := canvas.Circle(costContainerR)
	ctx.DrawPath(costContainerStart+(costContainerR), titleBoxTop-(titleBoxHeight*0.5), path)

	ctx.Pop()

	var boxText, boxType textBoxDimensions
	if card.Attributes.TrashCost != nil {
		boxText, boxType = drawTextBoxTrashable(ctx, canvasHeight/192, cornerRounded)
	} else {
		boxText, boxType = drawTextBox(ctx, canvasHeight/192, cornerRounded)
	}

	drawInfluence(ctx, card, boxText.left, factionColor)

	// render card text

	// not sure how these sizes actually correlate to the weird
	// pixel/mm setup I'm using, but these work
	fontSizeTitle := titleBoxHeight * 2
	fontSizeCost := titleBoxHeight * 3
	fontSizeCard := titleBoxHeight * 1.2

	drawAgendaPoints(ctx, card, fontSizeCost)

	titleTextX := 0.0
	titleTextY := titleBoxTop - titleBoxHeight*0.1
	ctx.DrawText(titleTextX, titleTextY, getCardText(getTitle(card), fontSizeTitle, titleBoxRight-(costContainerR*0.5), titleBoxHeight, canvas.Right))

	if card.Attributes.AdvancementRequirement != nil {
		costTextX := costContainerStart
		costTextY := titleBoxBottom + titleBoxHeight/2
		ctx.DrawText(costTextX, costTextY, canvas.NewTextBox(
			getFont(fontSizeCost, canvas.FontBlack), fmt.Sprint(*card.Attributes.AdvancementRequirement),
			costContainerR*2, 0,
			canvas.Center, canvas.Center, 0, 0))
	}

	drawCardText(ctx, card, fontSizeCard, canvasHeight, 0, boxText)
	drawTypeText(ctx, card, fontSizeCard, boxType)

	return nil
}