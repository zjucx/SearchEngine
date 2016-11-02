package web

import (
	_ "web/routers"
	"github.com/astaxie/beego"
)

func Main() {
	beego.Run()
}
