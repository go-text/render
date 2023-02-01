package render

import (
	"image/color"
	"image/draw"
	"math"

	"github.com/benoitkugler/textlayout/fonts"
	"github.com/go-text/typesetting/shaping"
	"github.com/srwiley/rasterx"
	"golang.org/x/image/math/fixed"
)

type Renderer struct {
	Face               fonts.Face // TODO []Face?
	FontSize, PixScale float32
	Color              color.Color
}

func (r *Renderer) DrawString(str string, img draw.Image) int {
	sh := &shaping.HarfbuzzShaper{}
	in := shaping.Input{
		Text:     []rune(str),
		RunStart: 0,
		RunEnd:   len(str),
		Face:     r.Face,
		Size:     fixed.I(int(r.FontSize * r.PixScale)),
	}
	out := sh.Shape(in)
	return r.DrawShapedRunAt(out, img, 0, out.LineBounds.Ascent.Ceil())
}

func (r *Renderer) DrawStringAt(str string, img draw.Image, x, y int) int {
	sh := &shaping.HarfbuzzShaper{}
	in := shaping.Input{
		Text:     []rune(str),
		RunStart: 0,
		RunEnd:   len(str),
		Face:     r.Face,
		Size:     fixed.I(int(r.FontSize * r.PixScale)),
	}
	return r.DrawShapedRunAt(sh.Shape(in), img, x, y)
}

func (r *Renderer) DrawShapedRunAt(run shaping.Output, img draw.Image, startX, startY int) int {
	scale := r.FontSize * r.PixScale / float32(run.Face.Upem())

	b := img.Bounds()
	scanner := rasterx.NewScannerGV(b.Dx(), b.Dy(), img, b)
	f := rasterx.NewFiller(b.Dx(), b.Dy(), scanner)
	f.SetColor(r.Color)
	point := uint16(float32(run.Face.Upem()) * r.PixScale)
	x := float32(startX)
	y := float32(startY)
	for _, g := range run.Glyphs {
		x -= fixed266ToFloat(g.XOffset) * r.PixScale
		outline, _ := run.Face.GlyphData(g.GlyphID, point, point).(fonts.GlyphOutline)

		for _, s := range outline.Segments {
			switch s.Op {
			case fonts.SegmentOpMoveTo:
				f.Start(fixed.Point26_6{X: floatToFixed266(s.Args[0].X*scale + x), Y: floatToFixed266(-s.Args[0].Y*scale + y)})
			case fonts.SegmentOpLineTo:
				f.Line(fixed.Point26_6{X: floatToFixed266(s.Args[0].X*scale + x), Y: floatToFixed266(-s.Args[0].Y*scale + y)})
			case fonts.SegmentOpQuadTo:
				f.QuadBezier(fixed.Point26_6{X: floatToFixed266(s.Args[0].X*scale + x), Y: floatToFixed266(-s.Args[0].Y*scale + y)},
					fixed.Point26_6{X: floatToFixed266(s.Args[1].X*scale + x), Y: floatToFixed266(-s.Args[1].Y*scale + y)})
			case fonts.SegmentOpCubeTo:
				f.CubeBezier(fixed.Point26_6{X: floatToFixed266(s.Args[0].X*scale + x), Y: floatToFixed266(-s.Args[0].Y*scale + y)},
					fixed.Point26_6{X: floatToFixed266(s.Args[1].X*scale + x), Y: floatToFixed266(-s.Args[1].Y*scale + y)},
					fixed.Point26_6{X: floatToFixed266(s.Args[2].X*scale + x), Y: floatToFixed266(-s.Args[2].Y*scale + y)})
			}
		}
		f.Stop(true)

		x += fixed266ToFloat(g.XAdvance)
	}
	f.Draw()
	return int(math.Ceil(float64(x)))
}

func fixed266ToFloat(i fixed.Int26_6) float32 {
	return float32(float64(i) / 64)
}

func floatToFixed266(f float32) fixed.Int26_6 {
	return fixed.Int26_6(int(float64(f) * 64))
}
