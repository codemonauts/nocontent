package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
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

func addLabel(img *image.RGBA, width int, height int, label string) {
	col := color.RGBA{255, 255, 255, 255}

	// Single character is 13px tall and 7px wide
	// Get middle of image and apply offset depending on label length
	x := (width / 2) - (len(label)*7)/2
	y := (height / 2) + 6

	point := fixed.Point26_6{fixed.Int26_6(x * 64), fixed.Int26_6(y * 64)}

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(col),
		Face: basicfont.Face7x13,
		Dot:  point,
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

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	headers := map[string]string{
		"Content-Type":                "image/png",
		"Access-Control-Allow-Origin": "*",
	}

	fmt.Println("%+v", request)
	params := request.QueryStringParameters

	x, _ := strconv.Atoi(params["x"])
	y, _ := strconv.Atoi(params["y"])
	label := fmt.Sprintf("%d x %d", x, y)
	bgColor, _ := parseHexColor(params["bg"])

	image := createImage(x, y, bgColor)
	addLabel(image, x, y, label)

	buf := new(bytes.Buffer)
	png.Encode(buf, image)

	b64Image := base64.StdEncoding.EncodeToString(buf.Bytes())

	return events.APIGatewayProxyResponse{StatusCode: 200, Headers: headers, Body: b64Image, IsBase64Encoded: true}, nil
}

func main() {
	lambda.Start(Handler)
}
