package render

import (
	"bytes"
	"image"
	"image/draw"
	_ "image/jpeg" // load image formats for users of the API
	_ "image/png"
	"log"

	"github.com/go-text/typesetting/shaping"
	_ "golang.org/x/image/tiff" // load image formats for users of the API

	"github.com/go-text/typesetting/opentype/api"
	"github.com/nfnt/resize"
)

func (r *Renderer) drawBitmap(g shaping.Glyph, bitmap api.GlyphBitmap, img draw.Image, x, y float32) (advance float32, err error) {
	adv := float32(0)

	switch bitmap.Format {
	case api.BlackAndWhite:
		log.Println("black and white - TODO")
	case api.JPG, api.PNG, api.TIFF:
		pix, _, err := image.Decode(bytes.NewReader(bitmap.Data))
		if err != nil {
			return 0, err
		}

		h := r.FontSize * r.PixScale
		scaled := resize.Resize(uint(h), uint(h), pix, resize.Bicubic)
		draw.Draw(img, scaled.Bounds().Add(image.Point{X: int(x), Y: int(y)}), scaled, image.Point{}, draw.Over)

		adv = h
	}

	if bitmap.Outline != nil {
		adv += r.drawOutline(g, *bitmap.Outline, r.filler, r.fillerScale, x, y)
	}
	return adv, nil
}
