/**
 * @Author $
 * @Description //TODO $
 * @Date $ $
 * @Param $   子代理 管理
 * @return $
 **/
package agency

import (
	"crypto/md5"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/wangyi/fishpond/dao/mysql"
	"github.com/wangyi/fishpond/model"
	"github.com/wangyi/fishpond/util"
	"strconv"
	"time"
)

/***

  获取  分销代理
*/

func GetSizingAgent(c *gin.Context) {
	_, err2 := c.Get("who")
	if !err2 {
		return
	}
	action := c.PostForm("action")
	if action == "GET" {
		data := make([]model.Admin, 0)
		belong := c.PostForm("belong")
		err := mysql.DB.Where("level > 0").Where("belong=?", belong).Find(&data).Error
		if err != nil {
			util.JsonWrite(c, -101, nil, "获取失败")
			return
		}
		util.JsonWrite(c, 200, data, "获取成功")
		return
	}

	if action == "ADD" {
		username := c.PostForm("username")
		password := c.PostForm("password")
		password = fmt.Sprintf("%x", md5.Sum([]byte(password)))
		level := c.PostForm("level")
		belongOne := c.PostForm("belong")
		belong, _ := strconv.Atoi(belongOne)
		//判断这个用户是否存在
		admin := model.Admin{}
		err := mysql.DB.Where("username=?", username).First(&admin).Error
		if err == nil {
			util.JsonWrite(c, -101, nil, "不要重复添加!")
			return
		}
		var TheOnlyInvited string
		for i := 0; i < 5; i++ {
			TheOnlyInvited = util.RandStr(6)
			err := mysql.DB.Where("the_only_invited=?", TheOnlyInvited).First(&model.Admin{}).Error
			if err != nil {
				break
			}
		}

		a, _ := strconv.Atoi(level)
		add := model.Admin{
			Username:       username,
			Password:       password,
			Level:          a,
			Updated:        time.Now().Unix(),
			Created:        time.Now().Unix(),
			Token:          util.RandStr(36),
			Belong:         belong,
			TheOnlyInvited: TheOnlyInvited,
		}
		err = mysql.DB.Save(&add).Error
		if err != nil {
			util.JsonWrite(c, -101, nil, "保存成功")
			return
		}
		util.JsonWrite(c, 200, nil, "添加成功!")

		return
	}
	if action == "UPDATE" {
		id := c.PostForm("id")
		//判断  这个管理员是否存在
		admin := model.Admin{}
		err := mysql.DB.Where("id=?", id).First(&admin).Error
		if err != nil {
			util.JsonWrite(c, -101, nil, "这个管理员不存在!")
			return
		}
		username := c.PostForm("username")

		password := c.PostForm("password")
		password = fmt.Sprintf("%x", md5.Sum([]byte(password)))
		up := model.Admin{
			Username: username,
			Password: password,
		}
		err = mysql.DB.Model(&model.Admin{}).Where("id=?", id).Update(&up).Error
		if err != nil {
			util.JsonWrite(c, -101, nil, "更新失败")
			return
		}
		util.JsonWrite(c, 200, nil, "修改成功")
		return
	}

	if action == "DEL" {
		id := c.PostForm("id")
		admin := model.Admin{}
		err := mysql.DB.Where("id=?", id).First(&admin).Error
		if err != nil {
			util.JsonWrite(c, -101, nil, "这个管理员不存在!")
			return
		}

		err = mysql.DB.Delete(&model.Admin{}, id).Error
		if err != nil {
			util.JsonWrite(c, 101, nil, "删除失败")

			return
		}
		util.JsonWrite(c, 200, nil, "删除成功")

		return

	}
	return
}
