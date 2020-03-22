package fake

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/color"
	"image/png"
	"math/rand"
	"time"
	"vobook/utils"

	"github.com/nfnt/resize"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

func PictureBase64() (pic string, err error) {
	picData, err := Picture()
	if err != nil {
		return
	}
	pic = base64.StdEncoding.EncodeToString(picData)
	return
}

func Picture(optSize ...int) (data []byte, err error) {
	const baseSize = 15
	size := baseSize

	if len(optSize) == 1 && optSize[0] > size {
		size = optSize[0]
	}

	rand.Seed(time.Now().UTC().UnixNano())

	pic := image.NewRGBA(image.Rect(0, 0, baseSize, baseSize))
	clr := color.RGBA{
		A: 255,
		R: uint8(rand.Intn(230)),
		G: uint8(rand.Intn(230)),
		B: uint8(rand.Intn(230)),
	}

	for y := 0; y < baseSize; y++ {
		for x := 0; x < baseSize; x++ {
			pic.SetRGBA(x, y, clr)
		}
	}

	d := &font.Drawer{
		Dst: pic,
		Src: image.NewUniform(color.RGBA{
			A: 255,
			R: uint8(rand.Intn(255)),
			G: uint8(rand.Intn(255)),
			B: uint8(rand.Intn(255)),
		}),
		Face: basicfont.Face7x13,
		Dot: fixed.Point26_6{
			Y: fixed.Int26_6((10 + baseSize/2 - 5) * 64),
			X: fixed.Int26_6(1 * 64),
		},
	}
	d.DrawString(utils.RandomString(2))

	buff := new(bytes.Buffer)
	if size != baseSize {
		newPic := resize.Resize(uint(size), 0, pic, resize.NearestNeighbor)
		err = png.Encode(buff, newPic)
		if err != nil {
			return
		}
		return buff.Bytes(), nil
	}

	err = png.Encode(buff, pic)
	if err != nil {
		return
	}
	return buff.Bytes(), nil
}
