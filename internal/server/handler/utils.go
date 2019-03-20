package handler

import (
	"bytes"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io"

	"github.com/pkg/errors"
)

func EncodePNG(b []byte, writer io.Writer) error {
	img, err := png.Decode(bytes.NewReader(b))
	if err != nil {
		l.WithError(errors.WithStack(err)).Error("failed to decode png file")
		return errors.WithStack(err)
	}

	target := image.NewRGBA(img.Bounds())
	draw.Draw(target, img.Bounds(), img, image.Point{0, 0}, draw.Src)

	return png.Encode(writer, target)
}

func EncodeJPEG(b []byte, writer io.Writer) error {
	img, err := jpeg.Decode(bytes.NewReader(b))
	if err != nil {
		l.WithError(errors.WithStack(err)).Error("failed to decode png file")
		return errors.WithStack(err)
	}

	target := image.NewRGBA(img.Bounds())
	draw.Draw(target, img.Bounds(), img, image.Point{0, 0}, draw.Src)

	return jpeg.Encode(writer, target, nil)
}
