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
	img := image.NewNRGBA(image.Rect(0, 0, 420, 180))
	draw.Draw(img, img.Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)

	data, _ := os.Open("testdata/NotoSans-Regular.ttf")
	f, _ := font.ParseTTF(data)

	r := &render.Renderer{
		FontSize: 48,
		Color:    color.Black,
	}
	str := "Hello! ± ज्या"
	r.DrawString(str, img, f)
	r.DrawStringAt(str, img, 0, 100, f)

	r.PixScale = 2
	r.Color = color.Gray{Y: 0xcc}
	r.DrawStringAt("baseline", img, 0, 180, f)

	data, _ = os.Open("testdata/NotoSans-Bold.ttf")
	f, _ = font.ParseTTF(data)
	r.FontSize = 36
	r.Color = color.NRGBA{R: 0xcc, G: 0, B: 0x33, A: 0x99}
	x := r.DrawStringAt("Red", img, 60, 140, f)
	r.DrawStringAt("Bold", img, x, 140, f)

	// from https://github.com/adobe-fonts/emojione-color, MIT license
	data, _ = os.Open("testdata/EmojiOneColor.otf")
	f, _ = font.ParseTTF(data)
	r.FontSize = 36
	r.DrawStringAt("😀🎁", img, 270, 80, f)

	data, _ = os.Open("testdata/Greybeard-22px.ttf")
	f, _ = font.ParseTTF(data)
	r.FontSize = 22
	r.Color = color.NRGBA{R: 0xcc, G: 0x66, B: 0x33, A: 0xcc}
	r.DrawStringAt("\uE0A2░", img, 362, 149, f)

	w, _ := os.Create("testdata/out.png")
	png.Encode(w, img)
	w.Close()
}

func TestRender_PixScaleAdvance(t *testing.T) {
	img := image.NewNRGBA(image.Rect(0, 0, 350, 180))

	data, _ := os.Open("testdata/NotoSans-Regular.ttf")
	f, _ := font.ParseTTF(data)

	r := &render.Renderer{
		FontSize: 48,
		Color:    color.Black,
	}
	str := "Testing"
	adv0 := r.DrawString(str, img, f)

	r.PixScale = 1 // instead of the zero value
	adv1 := r.DrawString(str, img, f)
	if adv0 != adv1 {
		t.Error("unscaled font did not advance as default")
	}

	r.PixScale = 2
	adv2 := r.DrawString(str, img, f)
	if adv2 <= int(float32(adv1)*1.9) || adv2 >= int(float32(adv1)*2.1) {
		t.Error("scaled font did not advance proportionately")
	}
}
