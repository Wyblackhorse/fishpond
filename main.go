package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/wangyi/fishpond/dao/mysql"
	"github.com/wangyi/fishpond/dao/redis"
	"github.com/wangyi/fishpond/logger"
	"github.com/wangyi/fishpond/setting"
	"go.uber.org/zap"
	"net/http"
)

func main() {





	//加载配置
	if err := setting.Init(); err != nil {
		fmt.Println("配置文件初始化事变", err)
		return
	}

	//初始化日志
	if err := logger.Init(); err != nil {
		fmt.Println("日志初始化失败", err)
		return
	}

	defer zap.L().Sync() //缓存日志追加到日志文件中
	zap.L().Debug("LaLa")

	//链接数据库
	if err := mysql.Init(); err != nil {
		fmt.Println("mysql 链接失败,", err)
		return
	}

	defer mysql.Close()

	//redis 初始化
	//4.初始化redis连接
	if err := redis.Init(); err != nil {
		fmt.Println("redis文件初始化失败：", err)
		return
	}
	defer redis.Close()

	// 1.创建路
	r := gin.New()








	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	// 2.绑定路由规则，执行的函数
	// gin.Context，封装了request和response
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "hello World!")
	})
	// 3.监听端口，默认在8080
	// Run("里面不指定端口号默认为8080")
	r.Run(":2345")

}
