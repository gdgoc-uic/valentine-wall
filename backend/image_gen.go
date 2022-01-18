package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
)

type ImageType int

func (i ImageType) Size() (width int, height int) {
	switch i {
	case imageTypeFacebook:
		return 1200, 630
	case imageTypeTwitter:
		return 1200, 675
	default:
		return 0, 0
	}
}

const (
	imageTypeFacebook ImageType = 0
	imageTypeTwitter  ImageType = 1
)

var latoBold *truetype.Font
var latoRegular *truetype.Font
var latoLight *truetype.Font
var nanumPenScript *truetype.Font

func loadFont(path string) *truetype.Font {
	fontData, err := os.ReadFile(path)
	if err != nil {
		log.Fatalln(err)
	}

	fontT, err := truetype.Parse(fontData)
	if err != nil {
		log.Fatalln(err)
	}

	return fontT
}

func init() {
	latoLight = loadFont("./fonts/lato/Lato-Light.ttf")
	nanumPenScript = loadFont("./fonts/nanum-pen-script/NanumPenScript-Regular.ttf")
}

// TODO: add support for emojis
// TODO: render images using chrome with chromedp/puppeteer/playwright
func generateImagePNG(wr io.Writer, itype ImageType, message Message) error {
	margin := float64(50)
	doubleMargin := 2.0 * margin
	innerContainerMargin := 4.0 * margin
	width, height := itype.Size()
	dc := gg.NewContext(width, height)

	containerStartX := float64(width) - doubleMargin
	containerEndY := float64(height) - doubleMargin
	innerContainerStartX := float64(width) - innerContainerMargin
	centerX := float64(width) / 2
	centerY := float64(height) / 2

	dc.DrawRectangle(0, 0, float64(width), float64(height))
	dc.SetRGB255(251, 207, 232)
	dc.Fill()

	dc.DrawRoundedRectangle(margin, margin, containerStartX, containerEndY, 30.0*float64(width/height))
	dc.SetRGB255(250, 242, 186)
	dc.Fill()

	// content := `Molestie id vulputate condimentum tempus nisi interdum scelerisque odio habitasse consectetur suspendisse hac id eu adipiscing.Arcu ridiculus felis nec adipiscing gravida platea neque parturient posuere faucibus laoreet vestibulum feugiat.`
	content := message.Content
	fontReductionFactor := 1.0
	threshold := 110.0
	lineHeight := 1.5
	offsetY := 20.0

	for size := float64(len(content)); size > threshold; size -= threshold {
		fontReductionFactor -= 0.09
		lineHeight -= 0.23
		offsetY -= 7
	}

	fontSize := (float64(height) * 0.075) * fontReductionFactor

	// paper lines
	dc.SetRGB255(24, 74, 153)
	dc.SetLineWidth(2)
	for w := 20.0; doubleMargin+w < containerEndY; w += fontSize + 3 {
		y := doubleMargin + w
		dc.DrawLine(margin, y, float64(width)-margin, y)
		dc.Stroke()
	}

	dc.SetRGB255(0, 0, 0)
	dc.SetFontFace(truetype.NewFace(nanumPenScript, &truetype.Options{
		Size: fontSize,
	}))
	dc.DrawStringWrapped(content, centerX, centerY+offsetY, 0.5, 0.5, innerContainerStartX, lineHeight, gg.AlignCenter)

	dc.SetRGB255(10, 10, 10)
	dc.SetFontFace(truetype.NewFace(latoLight, &truetype.Options{
		Size: float64(width) * 0.02,
	}))
	dc.DrawStringWrapped(fmt.Sprintf("Posted on %s", message.CreatedAt), centerX, containerEndY+10, 0.5, 0.5, innerContainerStartX, 1, gg.AlignCenter)

	return dc.EncodePNG(wr)
}
