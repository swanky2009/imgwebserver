package imgwebserver

import (
	log "github.com/Sirupsen/logrus"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"imgwebserver/g"
)

type imageServer struct {
	ctx *imgContext
}

func (p *imageServer) Init() {
	beego.Router("/", &HomeController{}, "get:Get")
	beego.Get("/health", func(ctx *context.Context) {
		ctx.Output.Body([]byte("ok"))
	})
	beego.Get("/version", func(ctx *context.Context) {
		ctx.Output.Body([]byte(g.Version()))
	})
	beego.Get("/notfound", func(ctx *context.Context) {
		ctx.Output.Body([]byte("the image not found"))
	})

	beego.Router("/img/:filename([\\w]+)_:width([0-9]+)x:height([0-9]+).:ext([\\w]+)", &ResizeController{}, "get:Get")

	beego.Router("/img/:filename([\\w]+).:ext([\\w]+)", &OriginalController{}, "get:Get")
}

func (p *imageServer) Start() {

	var HTTPAddress = p.ctx.imgwebserver.getOpts().HTTPAddress

	log.Infof("HTTP: listening on %s", HTTPAddress)

	beego.AppConfig.Set("imgpath", p.ctx.imgwebserver.getOpts().UploadPath)

	beego.Run(HTTPAddress)
}
