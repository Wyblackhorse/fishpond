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
	"strconv"
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

	action := c.PostForm("action")
	who, _ := c.Get("who")
	MapWho := who.(map[string]interface{})
	if action == "GET" {
		admin := model.Admin{}
		err := mysql.DB.Where("id=?", MapWho["ID"]).First(&admin).Error
		if err != nil {
			util.JsonWrite(c, -101, nil, "没有找到用户名")
			return
		}
		util.JsonWrite(c, 200, admin.LongUrl, "获取成功")

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
	util.JsonWrite(c, 200, admin.TheOnlyInvited, "设置成功")
	return
}

/**

设置体验金 链接
*/

func SetExperienceUrl(c *gin.Context) {
	action := c.PostForm("action")
	who, _ := c.Get("who")
	MapWho := who.(map[string]interface{})
	if action == "GET" {
		admin := model.Admin{}
		err := mysql.DB.Where("id=?", MapWho["ID"]).First(&admin).Error
		if err != nil {
			util.JsonWrite(c, -101, nil, "获取失败")
			return
		}

		util.JsonWrite(c, 200, admin, "获取成功")
		return
	}

	if action == "UPDATE" {

		if generate, isE := c.GetPostForm("generate"); isE == true {
			if generate == "generate" {
				admin := model.Admin{}
				err := mysql.DB.Where("id=?", MapWho["ID"]).First(&admin).Error

				if err != nil {
					util.JsonWrite(c, -101, nil, "非法参数")
					return
				}

				if admin.ExperienceCode == "" {
					for i := 0; i < 10; i++ {
						code := util.RandStr(7)
						err := mysql.DB.Where("experience_code=?", code).First(&model.Admin{}).Error
						if err != nil { //有错误  说明不存在这个code
							//更新 数据
							mysql.DB.Model(&model.Admin{}).Where("id=?", MapWho["ID"]).Update(&model.Admin{ExperienceCode: code})
							util.JsonWrite(c, 200, code, "设置成功")
							return
						}
					}
				}
				util.JsonWrite(c, -101, nil, "已经设置过了")
				return
			}
		}

		ups := model.Admin{}
		if generate, isE := c.GetPostForm("ExperienceTime"); isE == true {
			times, _ := strconv.ParseInt(generate, 10, 64)
			ups.ExperienceTime = times
		}

		//DefaultEarningsMoney
		if generate, isE := c.GetPostForm("DefaultEarningsMoney"); isE == true {
			times, _ := strconv.ParseFloat(generate, 64)
			ups.DefaultEarningsMoney = times
			err := mysql.DB.Model(&model.Admin{}).Where("id=?", MapWho["ID"]).Update(map[string]interface{}{"DefaultEarningsMoney":times}).Error
			if err != nil {
				util.JsonWrite(c, -101, nil, "修改失败")

				return
			}
			util.JsonWrite(c, 200, nil, "修改成功")
			return
		}

		//ExperienceMoney
		if generate, isE := c.GetPostForm("ExperienceMoney"); isE == true {
			times, _ := strconv.ParseFloat(generate, 64)
			ups.ExperienceMoney = times
		}

		err := mysql.DB.Model(&model.Admin{}).Where("id=?", MapWho["ID"]).Update(&ups).Error
		if err != nil {
			util.JsonWrite(c, -101, nil, "修改失败")
			return
		}

		util.JsonWrite(c, 200, nil, "修改成功")
		return

	}

}
