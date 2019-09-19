package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"image/png"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

func addLabel(img *image.RGBA, width int, height int, label string) {
	col := color.RGBA{255, 255, 255, 255}

	// Single character is 13px tall and 7px wide
	x := (width / 2) - (len(label)*7)/2
	y := (height / 2) - 6

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

	x := 400
	y := 200
	label := fmt.Sprintf("%d x %d", x, y)

	image := createImage(x, y, color.RGBA{50, 50, 50, 25})
	addLabel(image, x, y, label)

	buf := new(bytes.Buffer)
	png.Encode(buf, image)

	b64Image := base64.StdEncoding.EncodeToString(buf.Bytes())

	return events.APIGatewayProxyResponse{StatusCode: 200, Headers: headers, Body: b64Image, IsBase64Encoded: true}, nil
}

func main() {
	lambda.Start(Handler)
}
