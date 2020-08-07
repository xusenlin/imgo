# imgo
Golang 图片工具箱

## 1.图片内容的查找替换

- `func (p *Picture) SetCompareAccuracy(compareAccuracy int)` 设置图片在查找过程找到图片后对比的精确度，1代表100%完全吻合。
- `func (p *Picture) SearchPic(searchPic *Picture) (bool, image.Rectangle)` 在大图中查找小图出现的一个区域
- `func (p *Picture) SearchAllPic(searchPic *Picture) (bool, []image.Rectangle)` 在大图中查找小图出现的多个区域
- `func (p *Picture) Replace(searchPic *Picture, replacer *Picture) (image.Image, error)` 在大图中查找并替换小图的一个区域
- `func (p *Picture) ReplaceAll(searchPic *Picture, replacer *Picture) (image.Image, error)`在大图中查找并替换小图的多个区域

### demo
```go
pic, err := imgo.NewPng("./cat.png")
searchPic, err := imgo.NewPng("./x.png")
replace, err := imgo.NewPng("./text.png")
p, err := pic.ReplaceAll(searchPic, replace)
```
> 如果你有几万张设计稿需要替换文字或者一些图案，此工具库可能很好的帮助你编写代码
