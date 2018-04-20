package imgwebserver

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/astaxie/beego"
	"github.com/disintegration/imaging"
	"imgwebserver/g"
	"io/ioutil"
	"strconv"
	"strings"
)

type HomeController struct {
	beego.Controller
}

func (this *HomeController) Get() {
	this.Data["host"] = this.Ctx.Request.Host
	this.Data["version"] = g.Version()
	this.TplName = "index.html"
}

type OriginalController struct {
	beego.Controller
	ImgPath string
}

func (this *OriginalController) Prepare() {
	this.ImgPath = beego.AppConfig.String("imgpath")
}

func (this *OriginalController) Get() {
	filename := this.Ctx.Input.Param(":filename")
	fileext := this.Ctx.Input.Param(":ext")

	if filename == "" || fileext == "" {
		this.Redirect("/notfound", 404)
		return
	}
	filename = strings.ToLower(filename)
	fileext = strings.ToLower(fileext)

	filename = fmt.Sprintf("%s%s.%s", this.ImgPath, filename, fileext)

	//log.Infof("filename : %s", filename)

	filebytes, err := ioutil.ReadFile(filename)
	if err != nil {
		this.Redirect("/notfound", 404)
		return
	}
	this.Ctx.Output.Body(filebytes)
}

type ResizeController struct {
	beego.Controller
	ImgPath string
}

func (this *ResizeController) Prepare() {
	this.ImgPath = beego.AppConfig.String("imgpath")
}

func (this *ResizeController) Get() {
	filename := this.Ctx.Input.Param(":filename")
	fileext := this.Ctx.Input.Param(":ext")
	if filename == "" || fileext == "" {
		this.Redirect("/notfound", 404)
		return
	}

	filename = strings.ToLower(filename)
	fileext = strings.ToLower(fileext)

	if !(fileext == "jpg" || fileext == "jpeg" || fileext == "png") {
		//except jpg/png of other image can not resize
		this.Redirect("/notfound", 404)
		return
	}

	width, err_w := strconv.Atoi(this.Ctx.Input.Param(":width"))
	height, err_h := strconv.Atoi(this.Ctx.Input.Param(":height"))
	if err_w != nil || err_h != nil {
		log.Warnf("width/height error : %v %v", err_w, err_h)
		this.Redirect("/notfound", 404)
		return
	}

	original_filename := fmt.Sprintf("%s%s.%s", this.ImgPath, filename, fileext)
	resize_filename := fmt.Sprintf("%s%s_%dx%d.%s", this.ImgPath, filename, width, height, fileext)

	log.Debugf("original filename : %s", original_filename)
	log.Debugf("resize filename : %s", resize_filename)
	log.Debugf("filesize : %dx%d", width, height)

	//图片文件存在直接输出
	filebytes, err := ioutil.ReadFile(resize_filename)
	if err == nil {
		this.Ctx.Output.Body(filebytes)
		return
	}

	//图片文件不存在，读取原图生成压缩图片
	img, err1 := imaging.Open(original_filename)
	if err1 != nil {
		log.Warnf("failed to open image: %v", err1)
		this.Redirect("/notfound", 404)
		return
	}

	original_width := img.Bounds().Dx()
	original_hegith := img.Bounds().Dy()

	//原图宽大于高，以宽为准，反之亦然
	if original_width > original_hegith {
		height = 0
	} else if original_width < original_hegith {
		width = 0
	}
	img = imaging.Resize(img, width, height, imaging.Lanczos)

	err = imaging.Save(img, resize_filename)
	if err != nil {
		log.Warnf("failed to save new image: %v", err)
		return
	}

	var format imaging.Format

	if fileext == "jpg" || fileext == "jpeg" {
		format = imaging.JPEG
	} else {
		format = imaging.PNG
	}

	imaging.Encode(this.Ctx.ResponseWriter, img, format)
}
