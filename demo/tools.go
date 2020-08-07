package main

import (
	"fmt"
	"github.com/xusenlin/imgo"
	"image/png"
	"io/ioutil"
	"os"
	"strings"
)



func main()  {
	const readDirName  = "./imgs/"
	const distDirName  = "./dist/"
	const searchPicName  =  "./x.png"
	const replacePicName  =  "./text.png"


	searchPic, err := imgo.NewPng(searchPicName)
	if err != nil {
		panic(err)
	}
	replace, err := imgo.NewPng(replacePicName)
	if err != nil {
		panic(err)
	}
	imgDir, err := ioutil.ReadDir(readDirName)
	if err != nil {
		panic(err)
	}
	for _, fileInfo := range imgDir {
		name := strings.ToLower(fileInfo.Name())
		if strings.HasSuffix(name, "png"){
			pic, err := imgo.NewPng(readDirName + fileInfo.Name())
			if  err != nil {
				fmt.Println(err)
				continue
			}
			err = replaceFunc(distDirName + name,pic,searchPic,replace)
			if  err != nil {
				fmt.Println(err)
				continue
			}
		}
		if strings.HasSuffix(name, "jpeg") ||
			strings.HasSuffix(name, "jpg") {
			pic, err := imgo.NewJpeg(readDirName + fileInfo.Name())
			if  err != nil {
				fmt.Println(err)
				continue
			}
			err = replaceFunc(distDirName + name,pic,searchPic,replace)
			if  err != nil {
				fmt.Println(err)
				continue
			}
		}
	}
}

func replaceFunc(path string,pic *imgo.Picture,searchPic *imgo.Picture,replace *imgo.Picture) error {
	p, err := pic.Replace(searchPic, replace)
	if err != nil {
		return err
	}
	file, err := os.Create(path)
	defer file.Close()
	if err != nil {
		return err
	}

	err = png.Encode(file, p)
	if  err != nil {
		return err
	}
	return nil
}