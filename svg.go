package render

import (
	"bytes"
	"image"
	"image/draw"
	"io"

	"github.com/go-text/typesetting/opentype/api"
	"github.com/go-text/typesetting/shaping"
	"github.com/srwiley/oksvg"
	"github.com/srwiley/rasterx"
)

func (r *Renderer) drawSVG(g shaping.Glyph, svg api.GlyphSVG, img draw.Image, x, y float32) (advance float32, err error) {
	h := r.FontSize * r.PixScale
	pix, err := renderSVGStream(bytes.NewReader(svg.Source), int(h), int(h))
	if err != nil {
		return 0, err
	}
	draw.Draw(img, pix.Bounds().Add(image.Point{X: int(x), Y: int(y)}), pix, image.Point{}, draw.Over)

	if len(svg.Outline.Segments) > 0 {
		h += r.drawOutline(g, svg.Outline, r.filler, r.fillerScale, x, y)
	}
	return h, nil
}

func renderSVGStream(stream io.Reader, width, height int) (*image.NRGBA, error) {
	icon, err := oksvg.ReadIconStream(stream)
	if err != nil {
		return nil, err
	}

	iconAspect := float32(icon.ViewBox.W / icon.ViewBox.H)
	viewAspect := float32(width) / float32(height)
	imgW, imgH := width, height
	if viewAspect > iconAspect {
		imgW = int(float32(height) * iconAspect)
	} else if viewAspect < iconAspect {
		imgH = int(float32(width) / iconAspect)
	}

	icon.SetTarget(icon.ViewBox.X, icon.ViewBox.Y, float64(imgW), float64(imgH))

	out := image.NewNRGBA(image.Rect(0, 0, imgW, imgH))
	scanner := rasterx.NewScannerGV(int(icon.ViewBox.W), int(icon.ViewBox.H), out, out.Bounds())
	raster := rasterx.NewDasher(width, height, scanner)

	icon.Draw(raster, 1)
	return out, nil
}
