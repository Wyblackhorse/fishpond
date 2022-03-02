/**
 * @Author $
 * @Description //TODO $
 * @Date $ $
 * @Param $
 * @return $
 **/
package management

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

/**
  管理员登录
*/
func Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	password = fmt.Sprintf("%x", md5.Sum([]byte(password)))

	//
	admin := model.Admin{}
	err := mysql.DB.Where("status=1").Where("username=?", username).Where("password=?", password).First(&admin).Error
	if err != nil {
		util.JsonWrite(c, -101, nil, "登录失败")
		return
	}

	ArrayUrl := []string{"/management/login", "/agency/login", "/sonAgency/login"}
	for k, v := range ArrayUrl {
		if v == c.Request.URL.Path {
			if k != admin.Level {
				util.JsonWrite(c, -101, nil, "登录失败")
				return
			}
		}
	}

	//fmt.Println(c.ClientIP())
	//ss := strings.Split(admin.Ip, ",")
	//
	//if !util.InArray(c.ClientIP(), ss) {
	//	util.JsonWrite(c, -101, nil, "登录失败")
	//	return
	//}

	type Admin map[string]interface{}
	data := model.GetMenus(mysql.DB, strconv.Itoa(admin.Level))

	one := Admin{
		"data":  data,
		"admin": admin,
	}

	util.JsonWrite(c, 200, one, "登录成功")
	return
}

/***

  获取  分销代理  总代
*/

func GetSizingAgent(c *gin.Context) {
	_, err2 := c.Get("who")
	if !err2 {
		return
	}
	action := c.PostForm("action")
	if action == "GET" {
		data := make([]model.Admin, 0)
		err := mysql.DB.Where("level =1").Find(&data).Error
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
		//判断这个用户是否存在
		admin := model.Admin{}
		err := mysql.DB.Where("username=?", username).First(&admin).Error
		if err == nil {
			util.JsonWrite(c, -101, nil, "不要重复添加!")
			return
		}
		a, _ := strconv.Atoi(level)
		add := model.Admin{
			Username:       username,
			Password:       password,
			Level:          a,
			Updated:        time.Now().Unix(),
			Created:        time.Now().Unix(),
			Token:          util.RandStr(36),
			TheOnlyInvited: util.RandStr(40),
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

		if QRCodeSwitch, isE := c.GetPostForm("QRCodeSwitch"); isE == true {
			QR, _ := strconv.Atoi(QRCodeSwitch)
			up.QRCodeSwitch = QR //并且修改这个 总代下的所有子代
			admins := make([]model.Admin, 0)
			err := mysql.DB.Where("belong=?", id).Find(&admins).Error
			if err == nil {
				for _, k := range admins {
					mysql.DB.Model(&model.Admin{}).Where("id=?", k.ID).Update(&model.Admin{QRCodeSwitch: QR})
				}
			}

		}

		err = mysql.DB.Model(&model.Admin{}).Where("id=?", id).Update(&up).Error
		if err != nil {
			util.JsonWrite(c, -101, nil, "更新失败")
			return
		}

		util.JsonWrite(c, 200, nil, "修改成功")
		return

	}
	return
}

/**
  获取子代
*/
func GetSonAgent(c *gin.Context) {
	action := c.PostForm("action")
	if action == "GET" {
		adminId := c.PostForm("adminID")
		admin := make([]model.Admin, 0)
		err := mysql.DB.Where("belong=?", adminId).Find(&admin).Error
		if err != nil {
			util.JsonWrite(c, -101, nil, "获取失败")
			return
		}
		util.JsonWrite(c, 200, admin, "获取成功")
		return
	}
}
