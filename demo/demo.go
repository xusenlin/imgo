package main

import (
	"github.com/xusenlin/imgo"
	"image/png"
	"os"
)

func main() {

	pic, err := imgo.NewPng("./cat.png")

	if err != nil {
		panic(err)
	}

	searchPic, err := imgo.NewPng("./x.png")

	if err != nil {
		panic(err)
	}

	replace, err := imgo.NewPng("./text.png")

	if err != nil {
		panic(err)
	}


	p, err := pic.ReplaceAll(searchPic, replace)
	if err != nil {
		panic(err)
	}
	file, err := os.Create("./dst.jpg")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	err = png.Encode(file, p)
	if err != nil {
		panic(err)
	}
}
