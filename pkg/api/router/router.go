package router

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func SetUp(e *gin.Engine)  {
	routeSet()

	//合并发布模式
	//config模式需要重构
	if viper.GetBool("project.merge") {
		e.LoadHTMLGlob("./pkg/ui/dist/*.html") // 添加入口index.html
		e.Static("/static", "./pkg/ui/dist/static")   // 添加资源路径
		e.StaticFile("/", "./pkg/ui/dist/index.html") //前端接口
	}
}

//路由区,前期先写在一堆
func routeSet()  {

}

