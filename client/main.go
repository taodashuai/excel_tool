package main

import (
	"awesome2/controller"
	"fmt"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"github.com/kataras/iris/view"
	"os"
	"strconv"
)

func main() {
	e, _ := strconv.ParseInt(strconv.Itoa(0x00040), 16, 10)
	fmt.Println(e)

	app := iris.New()
	app.RegisterView(ViewHelper())
	app.StaticWeb("/web", LocalPath()+"web")
	mvc.Configure(app.Party("/").Layout("layout/layout.html"), Route)
	err := app.Run(iris.Addr(":8080"), iris.WithoutServerError(iris.ErrServerClosed))
	if err != nil {
		fmt.Println(err)
	}

}

// 解析html
func ViewHelper() *view.HTMLEngine {
	tpl := iris.HTML(LocalPath()+"web/html", ".html").Reload(false)
	return tpl
}

// 路由管理
func Route(app *mvc.Application) {
	app.Party("/").Handle(new(controller.IndexController))
}

// 获取项目路径
func LocalPath() string {
	path, _ := os.Getwd()
	return path+"/"
}
