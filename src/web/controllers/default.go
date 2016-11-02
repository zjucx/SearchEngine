package controllers

import (
	"github.com/astaxie/beego"
	"segment"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	segment.Segment()
	c.Data["Website"] = "beego.me"
	c.Data["Email"] = "astaxie@gmail.com"
	c.TplName = "index.tpl"
}
