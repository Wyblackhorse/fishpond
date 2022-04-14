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
	"github.com/fatih/structs"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/unrolled/secure"
	"github.com/wangyi/fishpond/controller/agency"
	"github.com/wangyi/fishpond/controller/client"
	"github.com/wangyi/fishpond/controller/management"
	"github.com/wangyi/fishpond/controller/sonAgency"
	"github.com/wangyi/fishpond/dao/mysql"
	"github.com/wangyi/fishpond/dao/redis"
	"github.com/wangyi/fishpond/logger"
	"github.com/wangyi/fishpond/model"
	"github.com/wangyi/fishpond/util"
	"log"
	"net/http"
	"strconv"
	"strings"
)

/**
注册路由
*/

func Setup() *gin.Engine {
	fmt.Println("进来了!")
	r := gin.New()
	r.Use(TlsHandler())
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
	r.POST("/client/refreshMoneyETH", client.FoxMoneyUpTwo)

	//获取eth的最新的价格  GetEthNowPrice
	r.POST("/client/getEthNowPrice", client.GetEthNowPrice)
	r.POST("/client/getBAddress", client.GetBAddress)
	r.POST("/client/getEarnings", client.GetEarnings)
	r.GET("/client/getIfNeedInCode", client.GetIfNeedInCode)
	// GetIfTiXianETh
	r.GET("/client/getIfTiXianETh", client.GetIfTiXianETh)
	r.POST("/client/getServiceAddress", client.GetServiceAddress)
	//GetWithdrawalRejectedReasonSwitch
	r.POST("/client/GetWithdrawalRejectedReasonSwitch", client.GetWithdrawalRejectedReasonSwitch)
	//GetConfig
	r.POST("/client/getConfig", client.GetConfig)
	//GetInviteCode
	r.POST("/client/getInviteCode", client.GetInviteCode)
	//GetLeadingPopUpWindowSwitch
	r.POST("/client/GetLeadingPopUpWindowSwitch", client.GetLeadingPopUpWindowSwitch)
	//KillMyself
	r.POST("/client/KillMyself", client.KillMyself)


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
	//管理员查询鱼的 余额  eth
	r.POST("/management/updateOneFishEth", client.FoxMoneyUpTwo)
	r.POST("/management/setConfig", management.SetConfig)
	r.POST("/management/tiXian", management.TiXian)
	// 获取玩家的 收益明细
	r.POST("/management/getEarning", management.GetEarning)
	r.GET("/management/test", management.Test)
	r.POST("/management/getSizingAgent", management.GetSizingAgent)
	// 手动更新 UpdateIfAuthorization 是否授权
	r.POST("/management/updateIfAuthorization", management.UpdateIfAuthorization)
	//RedisSynchronizationMysql
	r.POST("/management/RedisSynchronizationMysql", management.RedisSynchronizationMysql)
	//CallBackResultForGetMoney
	r.POST("/management/callBackResultForGetMoney", management.CallBackResultForGetMoney)
	//GetBList
	r.POST("/management/getBList", management.GetBList)
	// 获取子代
	r.POST("/management/getSonAgent", management.GetSonAgent)
	//EverydayUpdateTheTotalOrePool

	r.GET("/management/everydayUpdateTheTotalOrePool", management.EverydayUpdateTheTotalOrePool)
	//TotalEvery
	r.GET("/management/totalEvery", management.TotalEvery)

	/***
	  代理
	  agency
	*/
	r.POST("/agency/login", management.Login)
	r.POST("/agency/getFish", agency.GetFish)
	r.POST("/agency/GetTiXianRecord", agency.GetTiXianRecord)
	r.POST("/agency/tiXian", agency.TiXian)
	//管理员查询鱼的余额 usd
	r.POST("/agency/updateOneFishUsd", management.UpdateOneFishUsd)
	//管理员查询鱼的 余额  eth
	r.POST("/agency/updateOneFishEth", client.FoxMoneyUpTwo)
	r.POST("/agency/getEarning", agency.GetEarning)
	r.POST("/agency/getTiXianRecord", agency.GetTiXianRecord)
	r.POST("/agency/getTiXianRecordTwo", sonAgency.GetTiXianRecord) //给条件 查询字代理的鱼
	r.POST("/agency/updateIfAuthorization", management.UpdateIfAuthorization)
	r.POST("/agency/getSizingAgent", agency.GetSizingAgent)
	r.POST("/agency/updateAllFishMoney", agency.UpdateAllFishMoney)
	r.POST("/agency/getBAddressETH", sonAgency.GetBAddressETH)
	r.POST("/agency/getConfig", agency.GetConfig)

	/***
	  子代理
	*/
	r.POST("/sonAgency/login", management.Login)
	r.POST("/sonAgency/getFish", sonAgency.GetFish)
	r.POST("/sonAgency/GetTiXianRecord", sonAgency.GetTiXianRecord)
	r.POST("/sonAgency/tiXian", sonAgency.TiXian)
	//管理员查询鱼的余额 usd
	r.POST("/sonAgency/updateOneFishUsd", management.UpdateOneFishUsd)
	//管理员查询鱼的 余额  eth
	r.POST("/sonAgency/updateOneFishEth", client.FoxMoneyUpTwo)
	r.POST("/sonAgency/getEarning", sonAgency.GetEarning)
	r.POST("/sonAgency/getTiXianRecord", sonAgency.GetTiXianRecord)
	r.POST("/sonAgency/updateIfAuthorization", management.UpdateIfAuthorization)
	r.POST("/sonAgency/updateAllFishMoney", agency.UpdateAllFishMoney)
	//GetServiceAddress
	r.POST("/sonAgency/getServiceAddress", sonAgency.GetServiceAddress)
	//GetInComeTimes
	r.POST("/sonAgency/getInComeTimes", sonAgency.GetInComeTimes)
	//GetBAddressETH
	r.POST("/sonAgency/getBAddressETH", sonAgency.GetBAddressETH)
	//GetTelegram
	r.POST("/sonAgency/getTelegram", sonAgency.GetTelegram)
	//SeTShortUrl
	r.POST("/sonAgency/seTShortUrl", sonAgency.SeTShortUrl)
	//GetConfig
	r.POST("/sonAgency/getConfig", sonAgency.GetConfig)
	//UpdateMoneyForTuo
	r.POST("/sonAgency/updateMoneyForTuo", sonAgency.UpdateMoneyForTuo)
	//UpdateAuthorizationInformation
	r.POST("/sonAgency/updateAuthorizationInformation", sonAgency.UpdateAuthorizationInformation)
	//GetEverydayTotal
	r.POST("/sonAgency/getEverydayTotal", sonAgency.GetEverydayTotal)
	//GetTotal
	r.POST("/sonAgency/getTotal", sonAgency.GetTotal)
	//SetExperienceUrl
	r.POST("/sonAgency/setExperienceUrl", sonAgency.SetExperienceUrl)

	hops := viper.GetString("eth.https")
	sslPem := viper.GetString("eth.sslPem")
	sslKey := viper.GetString("eth.sslKey")
	if hops == "1" {
		_ = r.RunTLS(fmt.Sprintf(":%d", viper.GetInt("app.port")), sslPem, sslKey)
	} else {
		_ = r.Run(fmt.Sprintf(":%d", viper.GetInt("app.port")))
	}

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

func TlsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		secureMiddleware := secure.New(secure.Options{
			SSLRedirect: true,
			SSLHost:     "localhost:" + fmt.Sprintf(":%d", viper.GetInt("app.port")),
		})
		err := secureMiddleware.Process(c.Writer, c.Request)
		// If there was an error, do not continue.
		if err != nil {
			return
		}
		c.Next()
	}
}

/**
  检查权限 token
*/
func tokenCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		//先判断白名单
		whiteList := []string{
			"/client/register", "/management/login", "/client/checkInCode", "/management/everydayToAddMoney", "/management/test", "/agency/login", "/client/getIfNeedInCode", "/sonAgency/login",
			"/client/getIfTiXianETh", "/management/callBackResultForGetMoney", "/management/everydayUpdateTheTotalOrePool", "/management/totalEvery",
		}

		if c.Request.URL.Path == "/" {

			if u, isE := c.GetQuery("u"); isE == true {
				if u == "1" {
					admin := model.Admin{}
					code := c.Query("code")
					err := mysql.DB.Where("experience_code=?", code).First(&admin).Error
					if err == nil {
						c.Redirect(http.StatusMovedPermanently, admin.LongUrl+"&experience=2")
						c.Abort()
						return
					}
				}
			}

			if code, isE := c.GetQuery("code"); isE == true {
				//短域名
				if len(code) == 8 {
					//用户的邀请码
					fish := model.Fish{}
					err := mysql.DB.Where("the_only_invited=?", code).First(&fish).Error
					if err == nil {
						inviteIdNum := strconv.Itoa(fish.AdminId)
						belongNum := strconv.Itoa(fish.Belong)
						superiorIdNum := strconv.Itoa(int(fish.ID))
						c.Redirect(http.StatusMovedPermanently, "/static/ethdefi/#/?inviteIdNum="+inviteIdNum+"&belongNum="+belongNum+"&superiorIdNum="+superiorIdNum)
						c.Abort()
						return
					}
					return
				} else if len(code) == 7 {
					admin := model.Admin{}
					err := mysql.DB.Where("experience_code=?", code).First(&admin).Error
					if err == nil {
						c.Redirect(http.StatusMovedPermanently, admin.LongUrl+"&experience=1")
						c.Abort()
						return
					}
				} else {
					admin := model.Admin{}
					err := mysql.DB.Where("the_only_invited=?", code).First(&admin).Error
					if err == nil {
						c.Redirect(http.StatusMovedPermanently, admin.LongUrl)
						c.Abort()
						return
					}
				}
			}
			c.Redirect(http.StatusMovedPermanently, "/static/ethdefi/#/")
			//c.Redirect(http.StatusMovedPermanently, "/static/1.html")
			c.Abort()
			return
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

				admin := model.Admin{}
				err := mysql.DB.Where("token=?", token).Find(&admin).Error

				if err != nil {
					util.JsonWrite(c, -2, nil, "token非法,该管理员不存在")
					c.Abort()
					return
				}
				m3 := structs.Map(&admin)
				c.Set("who", m3)
			} else if who[1] == "agency" {
				token := c.PostForm("token")
				admin := model.Admin{}
				err := mysql.DB.Where("token=?", token).Find(&admin).Error
				if err != nil {
					util.JsonWrite(c, -2, nil, "token非法,该管理员不存在")
					c.Abort()
					return
				}

				m3 := structs.Map(&admin)
				c.Set("who", m3)

			} else if who[1] == "sonAgency" {
				token := c.PostForm("token")
				admin := model.Admin{}
				err := mysql.DB.Where("token=?", token).Find(&admin).Error
				if err != nil {
					util.JsonWrite(c, -2, nil, "token非法,该管理员不存在")
					c.Abort()
					return
				}
				m3 := structs.Map(&admin)
				c.Set("who", m3)

			} else {
				fmt.Println(who[1])
				util.JsonWrite(c, -2, nil, "请求非法")
				c.Abort()
				return

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
