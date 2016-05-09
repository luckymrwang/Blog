package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	. "myblog/models"
	// "myblog/utils"
)

type LoginController struct {
	beego.Controller
}

func (this *LoginController) Get() {
	username := this.GetSession("username")
	if username == nil {
		this.TplName = "admin/login.html"
	} else {
		this.Redirect("/admin", 302)
	}
}

func (this *LoginController) Post() {
	username := this.GetString("login_username")
	password := this.GetString("login_password")
	fmt.Println("username", username)
	fmt.Println("password", password)
	if username == "" || password == "" {
		this.Redirect("/login", 302)
		return
	}

	user, err := GetUserByName(username)
	fmt.Println("user", user)
	if err != nil {
		this.Data["msg"] = "此用户不存在"
	} else {
		// password = utils.Md5(password)
		if password == user.Password {
			this.SetSession("username", username)
			this.Redirect("/admin", 302)
		} else {
			this.Data["msg"] = "登录失败"
		}
	}
	this.TplName = "admin/login.html"
}

type LogoutController struct {
	beego.Controller
}

func (this *LogoutController) Get() {
	this.DelSession("username")
	this.Redirect("/admin", 302)
}
