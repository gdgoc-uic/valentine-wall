package main

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"io"
	"log"
	"os"
	"time"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"github.com/patrickmn/go-cache"
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

var latoLight *truetype.Font
var nanumPenScript *truetype.Font

func loadFont(path string) *truetype.Font {
	fontData, err := os.ReadFile(path)
	if err != nil {
		log.Panicln(err)
	}

	fontT, err := truetype.Parse(fontData)
	if err != nil {
		log.Panicln(err)
	}

	return fontT
}

func init() {
	latoLight = loadFont("./renderer_assets/fonts/lato/lato-v22-latin-regular.ttf")
	nanumPenScript = loadFont("./renderer_assets/fonts/nanum-pen-script/nanum-pen-script-v15-latin-regular.ttf")
}

type ImageRenderer struct {
	// for chrome-based image renderer
	ChromeCtx  context.Context
	Template   *template.Template
	CacheStore *cache.Cache
}

func (ctx *ImageRenderer) Render(itype ImageType, message Message) ([]byte, error) {
	// use cached image if available
	imageCacheKey := fmt.Sprintf("image/%s", message.ID)
	if cachedImage, isImageCached := ctx.CacheStore.Get(imageCacheKey); isImageCached && cachedImage != nil {
		log.Println("using cached image...")
		return cachedImage.([]byte), nil
	}

	imgBuf := &bytes.Buffer{}
	var err error
	if ctx.ChromeCtx != nil {
		// use alternative gg-based mode if not connected to chrome
		err = generateImagePNGChrome(imgBuf, ctx.ChromeCtx, ctx.Template, RendererContext{
			Message:    message,
			BackendURL: baseUrl,
		})
	}

	if err != nil || imgBuf.Len() == 0 {
		passivePrintError(err)
		if err2 := generateImagePNG(imgBuf, imageTypeTwitter, message); err2 != nil {
			return nil, err2
		}
	}

	if imgBuf.Len() != 0 {
		ctx.CacheStore.Set(imageCacheKey, imgBuf.Bytes(), cache.DefaultExpiration)
	}
	return imgBuf.Bytes(), nil
}

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

type RendererContext struct {
	Message
	BackendURL string
}

func generateImagePNGChrome(wr io.Writer, parentChromeCtx context.Context, tmpl *template.Template, rctx RendererContext) error {
	// compile template first
	output := &bytes.Buffer{}
	if err := tmpl.Execute(output, rctx); err != nil {
		return err
	}

	ctx, cancel := chromedp.NewContext(parentChromeCtx)
	defer cancel()

	actx, acancel := context.WithTimeout(ctx, 15*time.Second)
	defer acancel()

	// render to browser
	var buf []byte
	actions := chromedp.Tasks{
		chromedp.Navigate("about:blank"),
		chromedp.ActionFunc(func(c context.Context) error {
			frameTree, err := page.GetFrameTree().Do(c)
			if err != nil {
				return err
			}
			return page.SetDocumentContent(frameTree.Frame.ID, output.String()).Do(c)
		}),
		chromedp.Screenshot("#image-preview", &buf, chromedp.NodeVisible),
	}

	if err := chromedp.Run(actx, actions...); err != nil {
		return err
	} else if _, err := wr.Write(buf); err != nil {
		return err
	}

	return nil
}
