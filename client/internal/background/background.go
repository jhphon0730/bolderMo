package background

import (
	"image"
	_ "image/png"
	"os"
)

func LoadBackground() (image.Image, error) {
	bg_file, err := os.Open("assets/bg.png")
	if err != nil {
		return nil, err
	}
	defer bg_file.Close()

	img, _, err := image.Decode(bg_file)
	if err != nil {
		return nil, err
	}

	return img, nil
}

func LoadCharImage() (image.Image, error) {
	bg_file, err := os.Open("assets/char.png")
	if err != nil {
		return nil, err
	}
	defer bg_file.Close()

	img, _, err := image.Decode(bg_file)
	if err != nil {
		return nil, err
	}

	return img, nil
}
