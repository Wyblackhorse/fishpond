/**
 * @Author $
 * @Description //TODO $
 * @Date $ $
 * @Param $
 * @return $
 **/
package sonAgency

import (
	"github.com/gin-gonic/gin"
	"github.com/wangyi/fishpond/dao/mysql"
	"github.com/wangyi/fishpond/model"
	"github.com/wangyi/fishpond/util"
)

/**
  飞机通知
*/
func GetTelegram(c *gin.Context) {
	action := c.PostForm("action")

	who, _ := c.Get("who")
	mapWho := who.(map[string]interface{})

	if action == "GET" {

		admin := model.Admin{}
		err := mysql.DB.Where("id=?", mapWho["ID"]).First(&admin).Error
		if err != nil {
			util.JsonWrite(c, -101, nil, "获取失败")
			return
		}

		util.JsonWrite(c, 200, admin, "获取成功")
		return
	}

	if action == "UPDATE" {

		if _, isExist := c.GetPostForm("TelegramToken"); isExist != true {
			util.JsonWrite(c, -101, nil, "TelegramToken 不可以为空")
			return
		}
		if _, isExist := c.GetPostForm("TelegramChatId"); isExist != true {
			util.JsonWrite(c, -101, nil, "TelegramChatId 不可以为空")
			return
		}

		TelegramToken := c.PostForm("TelegramToken")
		TelegramChatId := c.PostForm("TelegramChatId")

		err := mysql.DB.Model(&model.Admin{}).Where("id=?", mapWho["ID"]).Update(&model.Admin{TelegramChatId: TelegramChatId, TelegramToken: TelegramToken}).Error
		if err != nil {
			util.JsonWrite(c, -101, nil, "更新失败")

			return
		}
		util.JsonWrite(c, 200, nil, "更新成功")
		return
	}

}

/**
设置短域名
*/

func SeTShortUrl(c *gin.Context) {

	action:=c.PostForm("action")
	who, _ := c.Get("who")
	MapWho := who.(map[string]string)
	if  action=="GET"{
		admin := model.Admin{}
		err := mysql.DB.Where("id=?", MapWho["ID"]).First(&admin).Error
		if err != nil {
			util.JsonWrite(c, -101, nil, "没有找到用户名")
			return
		}
		util.JsonWrite(c, -101, admin.LongUrl, "没有找到用户名")

		return
	}


	long := c.PostForm("long")

	admin := model.Admin{}
	err := mysql.DB.Where("id=?", MapWho["ID"]).First(&admin).Error
	if err != nil {
		util.JsonWrite(c, -101, nil, "没有找到用户名")
		return
	}
	err = mysql.DB.Model(&model.Admin{}).Where("id=?", admin.ID).Update(&model.Admin{LongUrl: long}).Error
	if err != nil {
		util.JsonWrite(c, -101, nil, "设置失败")
		return
	}
	util.JsonWrite(c, 200, nil, "设置成功")
	return
}
