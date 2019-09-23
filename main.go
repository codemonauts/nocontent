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
	"math"
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
		c.R = 255
		c.G = 255
		c.B = 255
	}
	return
}

func addLabel(img *image.RGBA, width int, height int, label string, fg color.RGBA) {
	// Read the font data.
	fontBytes, err := ioutil.ReadFile("Inconsolata.otf")
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
	size := 24
	dpi := 72
	d := &font.Drawer{
		Dst: img,
		Src: image.NewUniform(fg),
		Face: truetype.NewFace(f, &truetype.Options{
			Size:    size,
			DPI:     dpi,
			Hinting: h,
		}),
	}
	y := 10 + int(math.Ceil(size**dpi/72))
	dy := int(math.Ceil(size * *spacing * dpi / 72))
	d.Dot = fixed.Point26_6{
		X: (fixed.I(imgW) - d.MeasureString(title)) / 2,
		Y: fixed.I(y),
	}
	d.DrawString(title)
	y += dy
	for _, s := range text {
		d.Dot = fixed.P(10, y)
		d.DrawString(s)
		y += dy
	}
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
	y, _ := strconv.Atoi(params["y"])
	label := fmt.Sprintf("%d x %d", x, y)
	bgColor, _ := parseHexColor(params["bg"])
	fgColor, _ := parseHexColor(params["fg"])

	image := createImage(x, y, bgColor)
	addLabel(image, x, y, label, fgColor)

	buf := new(bytes.Buffer)
	png.Encode(buf, image)

	b64Image := base64.StdEncoding.EncodeToString(buf.Bytes())

	return events.APIGatewayProxyResponse{StatusCode: 200, Headers: headers, Body: b64Image, IsBase64Encoded: true}, nil
}

func main() {
	lambda.Start(Handler)
}
