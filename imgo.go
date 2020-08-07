package imgo

import (
	"errors"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"os"
)

const defaultCompareAccuracy = 10 //查找图片的精确值，默认查找图片平均有10分之一的像素对应即认为两部分图片一样。

type Picture struct {
	Img             image.Image
	Width           int
	Height          int
	Path            string
	CompareAccuracy int
}

func NewJpeg(path string) (*Picture, error) {

	read, err := os.Open(path)
	if err != nil {
		return &Picture{}, err
	}
	defer read.Close()

	img, err := jpeg.Decode(read)
	if err != nil {
		return &Picture{}, err
	}

	return newPic(img, path), nil
}

func NewPng(path string) (*Picture, error) {

	read, err := os.Open(path)
	if err != nil {
		return &Picture{}, err
	}
	defer read.Close()

	img, err := png.Decode(read)
	if err != nil {
		return &Picture{}, err
	}

	return newPic(img, path), nil
}

func (p *Picture) SetCompareAccuracy(compareAccuracy int) {
	p.CompareAccuracy = compareAccuracy
}

func (p *Picture) SearchPic(searchPic *Picture) (bool, image.Rectangle) {
	rectangles := seekPos(p, searchPic, true)
	if len(rectangles) == 0 {
		return false, image.Rectangle{}
	}
	return true, rectangles[0]
}

func (p *Picture) SearchAllPic(searchPic *Picture) (bool, []image.Rectangle) {
	rectangles := seekPos(p, searchPic, false)
	if len(rectangles) == 0 {
		return false, rectangles
	}
	return true, rectangles
}

func (p *Picture) Replace(searchPic *Picture, replacer *Picture) (image.Image, error) {

	if searchPic.Width != replacer.Width || searchPic.Height != replacer.Height {
		return p.Img, errors.New("查找和替换的图片大小不一致")
	}

	isExist, rectangle := p.SearchPic(searchPic)
	if !isExist {
		return p.Img, errors.New("在" + p.Path + "并未发现" + searchPic.Path)
	}

	dst := p.Img
	if dst, ok := dst.(draw.Image); ok {
		draw.Draw(dst, rectangle, replacer.Img, image.Point{}, draw.Src)
	}
	return dst, nil
}

func (p *Picture) ReplaceAll(searchPic *Picture, replacer *Picture) (image.Image, error) {

	if searchPic.Width != replacer.Width || searchPic.Height != replacer.Height {
		return p.Img, errors.New("查找和替换的图片大小不一致")
	}

	isExist, rectangles := p.SearchAllPic(searchPic)
	if !isExist {
		return p.Img, errors.New("在" + p.Path + "并未发现" + searchPic.Path)
	}

	dst := p.Img
	if dst, ok := dst.(draw.Image); ok {
		for _,rectangle := range rectangles{
			draw.Draw(dst, rectangle, replacer.Img, image.Point{}, draw.Src)
		}
	}
	return dst, nil
}

func newPic(img image.Image, path string) *Picture {

	rectangle := img.Bounds()
	w := rectangle.Max.X
	h := rectangle.Max.Y
	return &Picture{
		Img:             img,
		Width:           w,
		Height:          h,
		Path:            path,
		CompareAccuracy: defaultCompareAccuracy,
	}
}

func scanAreaOk(intX, intY int, p, searchPic *Picture) bool {
	if p.CompareAccuracy < 1 {
		p.SetCompareAccuracy(1)
	}
	for y := 0; y <= searchPic.Height-1; y += p.CompareAccuracy {
		for x := 0; x <= searchPic.Width-1; x += p.CompareAccuracy {
			if searchPic.Img.At(x, y) != p.Img.At(intX+x, intY+y) {
				return false
			}
		}
	}
	return true
}

func seekPos(p *Picture, searchPic *Picture, searchOnce bool) []image.Rectangle {
	var rectangles []image.Rectangle
	if searchPic.Width > p.Width || searchPic.Height > p.Height {
		return rectangles
	}
	for y := 0; y <= (p.Height - searchPic.Height); y++ {
		for x := 0; x <= (p.Width - searchPic.Width); x++ {
			if searchPic.Img.At(0, 0) != p.Img.At(x, y) ||
				searchPic.Img.At(searchPic.Width - 1, 0) != p.Img.At(x + searchPic.Width - 1, y) ||
				searchPic.Img.At(searchPic.Width - 1, searchPic.Height - 1) != p.Img.At(x + searchPic.Width - 1, y + searchPic.Height - 1) ||
				searchPic.Img.At(0, searchPic.Height - 1) != p.Img.At(x, y + searchPic.Height - 1) { //四个角只要有一个颜色对应不上直接跳到下一次
				continue
			}

			if !scanAreaOk(x, y, p, searchPic) { //四个角对上了在扫描区域，不成功直接下一次，
				continue
			}

			min := image.Point{X: x, Y: y}
			max := image.Point{X: x + searchPic.Width, Y: y + searchPic.Height}
			rectangles = append(rectangles, image.Rectangle{Min: min, Max: max})
			if searchOnce {
				return rectangles
			}
		}
	}
	return rectangles
}
