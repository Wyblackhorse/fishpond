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
	"github.com/wangyi/fishpond/controller/client"
	"github.com/wangyi/fishpond/controller/management"
	"github.com/wangyi/fishpond/dao/mysql"
	"github.com/wangyi/fishpond/dao/redis"
	"github.com/wangyi/fishpond/logger"
	"github.com/wangyi/fishpond/model"
	"github.com/wangyi/fishpond/util"
	"log"
	"net/http"
	"strings"
)

/**
注册路由
*/

func Setup() *gin.Engine {

	r := gin.New()
	//添加记录日志的中间件
	r.Use(logger.GinLogger(), logger.GinRecovery(true), Cors())
	r.Static("/static", "./static")
	r.NoRoute(func(c *gin.Context) {
		if c.Request.URL.Path == "/" {
			fmt.Println("===")
			c.Redirect(http.StatusMovedPermanently, "/static/ethdefi/#/")
			return
		}
		c.String(http.StatusOK, "404 not found2222")
		return
	})
	//r.Use(checkUrl())
	r.Use(tokenCheck())

	/**
	  @user  client
	*/
	r.POST("/client/register", client.FishRegister)
	r.POST("/client/getInformation", client.GetInformation)
	r.POST("/client/getVipEarnings", client.GetVipEarnings)
	r.POST("/client/tiXian", client.TiXian)
	//
	r.POST("/client/foxMoneyUp", client.FoxMoneyUpTwo)
	// 用户检查注册吗
	r.POST("/client/checkInCode", client.CheckInCode)
	r.POST("/client/checkAuthorization", client.CheckAuthorization)
	r.POST("/client/refreshMoney", client.UpdateOneFishUsd)
	//获取eth的最新的价格  GetEthNowPrice
	r.POST("/client/getEthNowPrice", client.GetEthNowPrice)
	r.POST("/client/getBAddress",client.GetBAddress)

	/***
	  管理员
	  management
	*/
	r.POST("/management/login", management.Login)
	r.POST("/management/getVipEarnings", management.GetVipEarnings)
	r.POST("/management/getFish", management.GetFish)
	r.GET("/management/everydayToAddMoney", management.EverydayToAddMoney)
	r.POST("/management/getTiXianRecord", management.GetTiXianRecord)
	r.POST("/management/SetInvitationCode", management.SetInvitationCode)
	//管理员查询鱼的余额 usd
	r.POST("/management/updateOneFishUsd", management.UpdateOneFishUsd)

	r.POST("/management/setConfig", management.SetConfig)

	r.POST("/management/tiXian", management.TiXian)

	r.Run(fmt.Sprintf(":%d", viper.GetInt("app.port")))
	return r
}

func checkUrl() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("检查:" + c.Request.URL.Path)
		if c.Request.URL.Path == "/" {
			c.Redirect(http.StatusMovedPermanently, "/")
			return
		}

	}

}

type Token struct {
	//不能为空并且大于10
	Token string `form:"token" binding:"required,len=36"`
}

/**
  检查权限 token
*/
func tokenCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		//先判断白名单
		whiteList := []string{
			"/client/register", "/management/login", "/client/checkInCode",
		}

		// 不需要判断 token
		if !util.InArray(c.Request.URL.Path, whiteList) { //摆在白名单里面需要校验token
			//获取token 参数
			//token := c.PostForm("token")
			var token Token

			if err := c.ShouldBind(&token); err != nil {
				util.JsonWrite(c, -2, nil, "非法参数!")
				c.Abort()
				return
			}
			//判断是谁访问进来了 然后在判断token 是否非法
			who := strings.Split(c.Request.URL.Path, "/")
			if len(who) < 2 {
				fmt.Println("这个后期再说!")
				return
			}

			if who[1] == "client" {
				//获取token
				token := c.PostForm("token")
				foxAddress, err := redis.Rdb.HGet("TOKEN_USER", token).Result()
				if err != nil {
					util.JsonWrite(c, -2, nil, "token非法!")
					return
				}
				information, err := redis.Rdb.HGetAll("USER_" + foxAddress).Result()
				if err != nil {
					util.JsonWrite(c, -2, nil, "token非法")
					return
				}
				c.Set("who", information)
				//token 检查结束
			} else if who[1] == "management" {
				//fmt.Println("管理员进来了!")
				token := c.PostForm("token")
				//foxAddress, err := redis.Rdb.HGet("TOKEN_ADMIN", token).Result()
				//if err != nil {
				//	util.JsonWrite(c, -2, nil, "token非法")
				//	//panic(err)
				//	return
				//}
				//information, err := redis.Rdb.HGetAll("ADMIN_" + foxAddress).Result()
				//if err != nil {
				//	util.JsonWrite(c, -2, nil, "token非法")
				//	return
				//}
				admin := model.Admin{}
				err := mysql.DB.Where("token=?", token).Find(&admin).Error

				if err != nil {
					util.JsonWrite(c, -2, nil, "token非法,该管理员不存在")
					return
				}
				c.Set("who", admin)
			}

		}

		c.Next()
	}

}

/**
跨域设置
*/
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin") //请求头部
		if origin != "" {
			//接收客户端发送的origin （重要！）
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			//服务器支持的所有跨域请求的方法
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE")
			//允许跨域设置可以返回其他子段，可以自定义字段
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session")
			// 允许浏览器（客户端）可以解析的头部 （重要）
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers")
			//设置缓存时间
			c.Header("Access-Control-Max-Age", "172800")
			//允许客户端传递校验信息比如 cookie (重要)
			c.Header("Access-Control-Allow-Credentials", "true")
		}

		//允许类型校验
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "ok!")
		}

		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic info is: %v", err)
			}
		}()

		c.Next()
	}
}
