package basic

import (
	"fmt"
	"log"

	"github.com/mangofeet/netrunner-alt-gen/art"
	"github.com/mangofeet/nrdb-go"
	"github.com/tdewolff/canvas"
)

type FrameProgram struct{}

func (FrameProgram) Draw(ctx *canvas.Context, card *nrdb.Printing) error {

	canvasWidth, canvasHeight := ctx.Size()

	strokeWidth := getStrokeWidth(ctx)

	log.Printf("strokeWidth: %f", strokeWidth)

	factionBaseColor := art.GetFactionBaseColor(card.Attributes.FactionID)
	factionColor := art.Darken(factionBaseColor, 0.811)

	ctx.Push()
	ctx.SetFillColor(bgColor)
	ctx.SetStrokeColor(textColor)
	ctx.SetStrokeWidth(strokeWidth)

	titleBoxHeight := getTitleBoxHeight(ctx)

	titleBoxTop := getTitleBoxTop(ctx)
	titleBoxBottom := titleBoxTop - titleBoxHeight
	titleBoxArcStart := canvasWidth - (canvasWidth / 3)
	titleBoxRight := canvasWidth - (canvasWidth / 16)
	titleBoxArcCP1 := titleBoxRight - (canvasWidth / 48)

	costContainerR := getCostContainerRadius(ctx)
	costContainerStart := getCostContainerStart(ctx)

	titlePath := &canvas.Path{}
	titlePath.MoveTo(0, titleBoxTop)

	// background for cost, top
	titlePath.LineTo(costContainerStart, titleBoxTop)
	titlePath.QuadTo(costContainerStart+costContainerR, titleBoxTop+(costContainerR), costContainerStart+(costContainerR*2), titleBoxTop)

	// arc down on right side
	titlePath.LineTo(titleBoxArcStart, titleBoxTop)
	titlePath.QuadTo(titleBoxArcCP1, titleBoxTop, titleBoxRight, titleBoxBottom)

	// background for cost, bottom
	titlePath.LineTo(costContainerStart+(costContainerR*2), titleBoxBottom)
	titlePath.QuadTo(costContainerStart+costContainerR, titleBoxBottom-(costContainerR), costContainerStart, titleBoxBottom)

	// finish title box
	titlePath.LineTo(0, titleBoxBottom)
	titlePath.Close()

	ctx.DrawPath(0, 0, titlePath)
	ctx.Pop()

	drawCostCircle(ctx, transparent)

	boxText, boxType := drawTextBox(ctx, canvasHeight/48, cornerRounded)

	drawInfluence(ctx, card, boxText.right, factionColor)

	// program strength
	ctx.Push()

	ctx.SetFillColor(factionColor)
	ctx.SetStrokeColor(textColor)
	ctx.SetStrokeWidth(strokeWidth)

	ctx.DrawPath(canvasWidth*-0.04, 0, strength(canvasWidth, canvasHeight))

	ctx.Pop()

	// mu icon
	muImage, err := loadGameAsset("Mu")
	if err != nil {
		return err
	}
	muImage = muImage.Transform(canvas.Identity.ReflectY()).Scale(0.05, 0.05)

	ctx.Push()

	ctx.SetFillColor(bgColor)
	ctx.SetStrokeColor(textColor)
	ctx.SetStrokeWidth(strokeWidth)

	muBoxX := costContainerStart + costContainerR*0.25
	muBoxY := titleBoxBottom - (muImage.Bounds().H * 0.8)
	muBoxW := muImage.Bounds().W + muImage.Bounds().W*0.35
	muBoxH := muImage.Bounds().H + muImage.Bounds().H*0.45

	boxPath := &canvas.Path{}

	boxPath.MoveTo(0, 0)
	boxPath.LineTo(muBoxW, 0)
	boxPath.LineTo(muBoxW, -1*muBoxH)
	boxPath.LineTo(0, -1*muBoxH)
	boxPath.Close()

	ctx.DrawPath(muBoxX, muBoxY, boxPath)

	ctx.Pop()

	ctx.Push()
	ctx.SetFillColor(textColor)

	muIconX := muBoxX
	muIconY := muBoxY + (muImage.Bounds().H * 0.2)

	ctx.DrawPath(muIconX, muIconY, muImage)

	ctx.Pop()

	// render card text

	// not sure how these sizes actually correlate to the weird
	// pixel/mm setup I'm using, but these work
	fontSizeTitle := titleBoxHeight * 2
	fontSizeCost := titleBoxHeight * 3
	fontSizeStr := titleBoxHeight * 4
	fontSizeCard := titleBoxHeight * 1.2

	titleTextX := costContainerStart + (costContainerR * 2) + (costContainerR / 3)
	titleTextY := titleBoxTop - titleBoxHeight*0.1
	ctx.DrawText(titleTextX, titleTextY, getCardText(getTitle(card), fontSizeTitle, titleBoxRight, titleBoxHeight, canvas.Left))
	// ctx.DrawText(titleTextX, titleTextY, canvas.NewTextLine(getFont(fontSizeTitle, canvas.FontRegular), getTitleText(card), canvas.Left))

	if card.Attributes.Cost != nil {
		costTextX := costContainerStart
		costTextY := titleBoxBottom + titleBoxHeight/2
		ctx.DrawText(costTextX, costTextY, canvas.NewTextBox(
			getFont(fontSizeCost, canvas.FontBlack), fmt.Sprint(*card.Attributes.Cost),
			costContainerR*2, 0,
			canvas.Center, canvas.Center, 0, 0))
	}

	strengthText := "-"
	if card.Attributes.Strength != nil {
		strengthText = fmt.Sprint(*card.Attributes.Strength)
	}

	strTextX := canvasWidth * 0.01
	strTextY := canvasHeight / 10
	ctx.DrawText(strTextX, strTextY, canvas.NewTextBox(
		getFont(fontSizeStr, canvas.FontBlack), strengthText,
		canvasWidth/5, 0,
		canvas.Center, canvas.Center, 0, 0))

	muText := ""
	if card.Attributes.MemoryCost != nil {
		muText = fmt.Sprint(*card.Attributes.MemoryCost)
	}

	muTextX := muBoxX - muBoxW*0.08
	muTextY := muBoxY
	ctx.DrawText(muTextX, muTextY, canvas.NewTextBox(
		getFont(fontSizeCard, canvas.FontBlack), muText,
		muBoxW, muBoxH,
		canvas.Center, canvas.Center, 0, 0))

	drawCardText(ctx, card, fontSizeCard, boxText.height*0.45, canvasWidth*0.06, boxText)
	drawTypeText(ctx, card, fontSizeCard, boxType)

	return nil
}