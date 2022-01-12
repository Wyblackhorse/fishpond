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
	err := mysql.DB.Where("status=1").Where("username=?", username).Where("password=?", password).Where("ip=?", c.ClientIP()).First(&admin).Error
	if err != nil {
		util.JsonWrite(c, -101, nil, "登录失败")
		return
	}

	type Admin map[string]interface{}
	data := model.GetMenus(mysql.DB, strconv.Itoa(admin.Level))

	one := Admin{
		"data":  data,
		"admin": admin,
	}

	util.JsonWrite(c, 200, one, "登录成功")
	return
}
