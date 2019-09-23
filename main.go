package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

func parseHexColor(s string) (c color.RGBA, err error) {
	c.A = 0xff

	hexToByte := func(b byte) byte {
		switch {
		case b >= '0' && b <= '9':
			return b - '0'
		case b >= 'a' && b <= 'f':
			return b - 'a' + 10
		case b >= 'A' && b <= 'F':
			return b - 'A' + 10
		}
		return 0
	}

	switch len(s) {
	case 6:
		c.R = hexToByte(s[0])<<4 + hexToByte(s[1])
		c.G = hexToByte(s[2])<<4 + hexToByte(s[3])
		c.B = hexToByte(s[4])<<4 + hexToByte(s[5])
	case 3:
		c.R = hexToByte(s[0]) * 17
		c.G = hexToByte(s[1]) * 17
		c.B = hexToByte(s[2]) * 17
	default:
		return c, fmt.Errorf("can't parse a color")
	}
	return
}

func addLable(img *image.RGBA, width int, height int, lable string, fg color.RGBA) {
	// Read the font data.
	fontBytes, err := ioutil.ReadFile("Inconsolata.ttf")
	if err != nil {
		log.Println(err)
		return
	}
	f, err := truetype.Parse(fontBytes)
	if err != nil {
		log.Println(err)
		return
	}
	// Draw the text.
	h := font.HintingFull
	size := 24.0
	dpi := 72.0
	d := &font.Drawer{
		Dst: img,
		Src: image.NewUniform(fg),
		Face: truetype.NewFace(f, &truetype.Options{
			Size:    size,
			DPI:     dpi,
			Hinting: h,
		}),
	}
	d.Dot = fixed.Point26_6{
		X: (fixed.I(width) - d.MeasureString(lable)) / 2,
		Y: fixed.I(height),
	}
	d.DrawString(lable)
}

func createImage(width int, height int, bgColor color.RGBA) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, bgColor)
		}
	}
	return img
}

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	headers := map[string]string{
		"Content-Type":                "image/png",
		"Expires":                     time.Now().AddDate(1, 0, 0).Format(http.TimeFormat),
		"Last-Modified":               time.Now().Format(http.TimeFormat),
		"Access-Control-Allow-Origin": "*",
	}

	params := request.QueryStringParameters

	x, _ := strconv.Atoi(params["x"])
	if x <= 0 || 2000 < x {
		x = 200
	}

	y, _ := strconv.Atoi(params["y"])
	if y <= 0 || 2000 < y {
		y = 200
	}

	lable, _ := params["lable"]
	if lable == "" || len(lable) > 20 {
		lable = fmt.Sprintf("%d x %d", x, y)
	}

	bgColor, err := parseHexColor(params["bg"])
	if err != nil {
		bgColor = color.RGBA{255, 255, 255, 255}
	}

	fgColor, err := parseHexColor(params["fg"])
	if err != nil {
		fgColor = color.RGBA{51, 51, 51, 255}
	}

	image := createImage(x, y, bgColor)
	addLable(image, x, y, lable, fgColor)

	buf := new(bytes.Buffer)
	png.Encode(buf, image)

	b64Image := base64.StdEncoding.EncodeToString(buf.Bytes())
	fmt.Println(b64Image)

	return events.APIGatewayProxyResponse{StatusCode: 200, Headers: headers, Body: b64Image, IsBase64Encoded: true}, nil
}

func main() {
	lambda.Start(Handler)
}
