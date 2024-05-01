package basic

import (
	"github.com/mangofeet/netrunner-alt-gen/art"
	"github.com/mangofeet/nrdb-go"
	"github.com/tdewolff/canvas"
)

func (fb FrameBasic) Upgrade() art.Drawer {
	return art.DrawerFunc(func(ctx *canvas.Context, card *nrdb.Printing) error {

		canvasWidth, canvasHeight := ctx.Size()

		strokeWidth := getStrokeWidth(ctx)

		factionColor := fb.getColor(card)

		ctx.Push()
		ctx.SetFillColor(bgColor)
		ctx.SetStrokeColor(textColor)
		ctx.SetStrokeWidth(strokeWidth)

		titleBoxHeight := getTitleBoxHeight(ctx)
		fontSizeCost := titleBoxHeight * 2.3
		boxResIcon, err := drawRezCost(ctx, card, fontSizeCost)
		if err != nil {
			return err
		}

		titleBoxTop := getTitleBoxTop(ctx)
		titleBoxBottom := titleBoxTop - titleBoxHeight
		titleBoxLeftOut := boxResIcon.left + (boxResIcon.width * 1.2)
		titleBoxLeftIn := titleBoxLeftOut + titleBoxHeight*0.8

		titlePath := &canvas.Path{}
		titlePath.MoveTo(titleBoxLeftIn, titleBoxTop)
		titlePath.LineTo(canvasWidth, titleBoxTop)
		titlePath.LineTo(canvasWidth, titleBoxBottom)
		titlePath.LineTo(titleBoxLeftOut, titleBoxBottom)
		titlePath.Close()

		ctx.DrawPath(0, 0, titlePath)
		ctx.Pop()

		var boxText, boxType textBoxDimensions
		if card.Attributes.TrashCost != nil {
			boxText, boxType = drawTextBoxTrashable(ctx, canvasHeight/192, cornerRounded)
		} else {
			boxText, boxType = drawTextBox(ctx, canvasHeight/192, cornerRounded)
		}

		drawInfluenceAndOrFactionSymbol(ctx, card, boxText.left, factionColor)

		if _, err := drawTrashCost(ctx, card); err != nil {
			return err
		}
		// render card text

		// not sure how these sizes actually correlate to the weird
		// pixel/mm setup I'm using, but these work
		fontSizeTitle := titleBoxHeight * 2
		fontSizeCard := titleBoxHeight * 1.2

		titleTextX := titleBoxLeftIn
		if card.Attributes.IsUnique { // unique diamon fits better in the angled end here
			titleTextX = titleBoxLeftIn - titleBoxHeight*0.2
		}

		titleTextMaxWidth := (canvasWidth * 0.9) - titleBoxLeftIn

		titleText := getTitleText(ctx, card, fontSizeTitle, titleTextMaxWidth, titleBoxHeight, canvas.Left)
		titleTextY := (titleBoxTop - (titleBoxHeight-titleText.Bounds().H)*0.5)
		ctx.DrawText(titleTextX, titleTextY, titleText)
		// canvas.NewTextLine(getFont(fontSizeTitle, canvas.FontRegular), getTitleText(card), canvas.Left))

		drawCardText(ctx, card, fontSizeCard, boxText.height*0.6, canvasWidth*0.02, boxText, fb.getAdditionalText()...)
		drawTypeText(ctx, card, fontSizeCard, boxType)

		return nil

	})
}
