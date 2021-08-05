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

func parseHexColor(s string, fallback color.RGBA) color.RGBA {
	c := color.RGBA{}
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
		return fallback
	}
	return c
}

func addLabel(img *image.RGBA, width int, height int, label string, size int, fg color.RGBA) {
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
	dpi := 72.0
	d := &font.Drawer{
		Dst: img,
		Src: image.NewUniform(fg),
		Face: truetype.NewFace(f, &truetype.Options{
			Size:    float64(size),
			DPI:     dpi,
			Hinting: h,
		}),
	}
	d.Dot = fixed.Point26_6{
		X: (fixed.I(width) - d.MeasureString(label)) / 2,
		Y: (fixed.I(height) + fixed.I(int(float64(size)*0.6))) / 2,
	}
	d.DrawString(label)
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

func validatePixelDimension(value int) int {
	if value <= 0 || 4000 < value {
		value = 200
	}

	return value
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
	x = validatePixelDimension(x)

	y, _ := strconv.Atoi(params["y"])
	y = validatePixelDimension(y)

	fontSize, _ := strconv.Atoi(params["fontSize"])

	label, _ := params["label"]
	if label == "" || len(label) > 20 {
		label = fmt.Sprintf("%d x %d", x, y)
	}

	bgColor := parseHexColor(params["bg"], color.RGBA{255, 255, 255, 255})

	fgColor := parseHexColor(params["fg"], color.RGBA{51, 51, 51, 255})

	image := createImage(x, y, bgColor)
	addLabel(image, x, y, label, fontSize, fgColor)

	buf := new(bytes.Buffer)
	png.Encode(buf, image)

	b64Image := base64.StdEncoding.EncodeToString(buf.Bytes())
	fmt.Println(b64Image)

	return events.APIGatewayProxyResponse{StatusCode: 200, Headers: headers, Body: b64Image, IsBase64Encoded: true}, nil
}

func main() {
	lambda.Start(Handler)
}
