/**
 * @Author $
 * @Description //TODO $
 * @Date $ $
 * @Param $
 * @return $
 **/
package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/wangyi/fishpond/logger"
)

/**
注册路由
*/

func Setup() *gin.Engine {

	r := gin.New()
	//添加记录日志的中间件
	r.Use(logger.GinLogger(), logger.GinRecovery(true))
	r.Run(fmt.Sprintf(":%d", viper.GetInt("app.port")))
	return r
}
