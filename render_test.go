package render_test

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"testing"

	"github.com/go-text/render"
	"github.com/go-text/typesetting/font"
)

func Test_Render(t *testing.T) {
	img := image.NewNRGBA(image.Rect(0, 0, 350, 180))
	draw.Draw(img, img.Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)

	data, _ := os.Open("testdata/NotoSans-Regular.ttf")
	fn, _ := font.ParseTTF(data)

	r := &render.Renderer{
		Face:     fn,
		FontSize: 24,
		PixScale: float32(2),
		Color:    color.Black,
	}
	str := "Hello! ± सभमन"
	r.DrawString(str, img)
	r.DrawStringAt(str, img, 0, 100)

	r.Color = color.Gray{Y: 0xcc}
	r.DrawStringAt("baseline", img, 0, 180)

	data, _ = os.Open("testdata/NotoSans-Bold.ttf")
	r.Face, _ = font.ParseTTF(data)
	r.FontSize = 36
	r.Color = color.NRGBA{R: 0xcc, G: 0, B: 0x33, A: 0xbb}
	r.DrawStringAt("RedBold", img, 60, 140)

	w, _ := os.Create("testdata/out.png")
	png.Encode(w, img)
	w.Close()
}
