package main

import (
	_ "web/routers"
	"github.com/astaxie/beego"
)

func main() {
	beego.SetStaticPath("/public", "public")
	beego.Run()
}
