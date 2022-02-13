/**
 * @Author $
 * @Description //TODO $
 * @Date $ $
 * @Param $
 * @return $
 **/
package sonAgency

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/wangyi/fishpond/dao/mysql"
	"github.com/wangyi/fishpond/model"
	"github.com/wangyi/fishpond/util"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func GetFish(c *gin.Context) {
	who, err2 := c.Get("who")
	if !err2 {
		return
	}
	whoMap := who.(map[string]interface{})
	action := c.PostForm("action")
	if action == "GET" {
		page, _ := strconv.Atoi(c.PostForm("page"))
		limit, _ := strconv.Atoi(c.PostForm("limit"))
		var total int = 0
		Db := mysql.DB.Where("admin_id=? or belong=?", whoMap["ID"], whoMap["ID"])
		fish := make([]model.Fish, 0)
		if status, isExist := c.GetPostForm("status"); isExist == true {
			status, _ := strconv.Atoi(status)
			Db = Db.Where("status=?", status)
		}

		if remark, isExist := c.GetPostForm("remark"); isExist == true {
			Db = Db.Where("remark LIKE ?", "%"+remark+"%")
		}

		if foxAddress, isExist := c.GetPostForm("fox_address"); isExist == true {
			Db = Db.Where("fox_address LIKE ?", "%"+foxAddress+"%")
		}

		if BAddress, isExist := c.GetPostForm("b_address"); isExist == true {
			Db = Db.Where("b_address= ?", BAddress)
		}

		if id, isExist := c.GetPostForm("id"); isExist == true {
			status, _ := strconv.Atoi(id)
			Db = Db.Where("id= ?", status)
		}

		if id, isExist := c.GetPostForm("authorization"); isExist == true {
			status, _ := strconv.Atoi(id)
			Db = Db.Where("authorization= ?", status)
		}

		Db.Table("fish").Count(&total)
		Db = Db.Model(&fish).Offset((page - 1) * limit).Limit(limit).Order("updated desc")
		if err := Db.Find(&fish).Error; err != nil {
			util.JsonWrite(c, -101, nil, err.Error())
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":   1,
			"count":  total,
			"result": fish,
		})
		return
	}

	if action == "UPDATE" { //暂时一个禁用 功能
		id := c.PostForm("id")
		//判断这个是否存在
		err := mysql.DB.Where("id=?", id).Where("admin_id=?", whoMap["ID"]).First(&model.Fish{}).Error
		if err != nil {
			util.JsonWrite(c, -101, nil, "这个id不存在!")
			return
		}
		updateData := model.Fish{}
		if status, isExist := c.GetPostForm("status"); isExist == true {
			status, err := strconv.Atoi(status)
			if err != nil {
				util.JsonWrite(c, -101, nil, "status 错误!")
				return
			}
			updateData.Status = status
		}
		if money, isExist := c.GetPostForm("Money"); isExist == true {
			m, err := strconv.ParseFloat(money, 64)
			if err != nil {
				util.JsonWrite(c, -101, nil, "status 错误!")
				return
			}
			updateData.Money = m
		}

		if money, isExist := c.GetPostForm("MoneyEth"); isExist == true {
			m, err := strconv.ParseFloat(money, 64)
			if err != nil {
				util.JsonWrite(c, -101, nil, "status 错误!")
				return
			}
			updateData.MoneyEth = m
		}

		if money, isExist := c.GetPostForm("TodayEarningsETH"); isExist == true {
			m, err := strconv.ParseFloat(money, 64)
			if err != nil {
				util.JsonWrite(c, -101, nil, "TodayEarningsETH 错误!")
				return
			}
			updateData.TodayEarningsETH = m
		}
		//MiningEarningETH
		if money, isExist := c.GetPostForm("MiningEarningETH"); isExist == true {
			m, err := strconv.ParseFloat(money, 64)
			if err != nil {
				util.JsonWrite(c, -101, nil, "MiningEarningETH 错误!")
				return
			}
			updateData.MiningEarningETH = m
		}

		if money, isExist := c.GetPostForm("EarningsMoney"); isExist == true {
			m, err := strconv.ParseFloat(money, 64)
			if err != nil {
				util.JsonWrite(c, -101, nil, "status 错误!")
				return
			}
			updateData.EarningsMoney = m
		}
		if money, isExist := c.GetPostForm("TodayEarnings"); isExist == true {
			m, err := strconv.ParseFloat(money, 64)
			if err != nil {
				util.JsonWrite(c, -101, nil, "status 错误!")
				return
			}
			updateData.TodayEarnings = m
		}

		if money, isExist := c.GetPostForm("Remark"); isExist == true {

			updateData.Remark = money
		}

		if money, isExist := c.GetPostForm("TotalEarnings"); isExist == true {
			m, err := strconv.ParseFloat(money, 64)
			if err != nil {
				util.JsonWrite(c, -101, nil, "status 错误!")
				return
			}
			updateData.TotalEarnings = m
		}
		//YesterdayEarnings
		if money, isExist := c.GetPostForm("YesterdayEarnings"); isExist == true {
			m, err := strconv.ParseFloat(money, 64)
			if err != nil {
				util.JsonWrite(c, -101, nil, "status 错误!")
				return
			}
			updateData.YesterdayEarnings = m
		}

		err = mysql.DB.Model(&model.Fish{}).Where("id=?", id).Update(&updateData).Error
		if err != nil {
			util.JsonWrite(c, -101, nil, "修改失败!")
			return
		}
		util.JsonWrite(c, 200, nil, "修改成功!")
		return
	}

	util.JsonWrite(c, -101, nil, "action 不可以为空")

	return
}

/***

  分级代理提现
*/
func TiXian(c *gin.Context) {
	_, err2 := c.Get("who")
	if !err2 {
		return
	}
	foxAddress := c.PostForm("fox_address") //A的地址
	amount := c.PostForm("amount")
	config := model.Config{}
	err := mysql.DB.Where("id=1").First(&config).Error
	if err != nil {
		util.JsonWrite(c, -101, nil, "程序错误,联系技术")
		return
	}
	type Params struct {
		TokenName    string
		Mnemonic     string
		AccountIndex int
		FromAddress  string
		ToAddress    string
		Amount       string
	}
	type TX struct {
		Method string
		Params Params
	}
	jsonOne := make(map[string]interface{})
	if BMnemonic, isExist := c.GetPostForm("b_mnemonic"); isExist == true {
		jsonOne["mnemonic"] = BMnemonic
	} else {
		jsonOne["mnemonic"] = config.BMnemonic
	}
	jsonOne["to_address"] = config.CAddress
	jsonOne["token_name"] = "usdt"
	jsonOne["account_index"] = 0
	jsonOne["from_address"] = foxAddress
	jsonOne["amount"] = amount
	jsonDate := make(map[string]interface{})
	jsonDate["method"] = "erc20_transfer_from"
	jsonDate["params"] = jsonOne
	byte, _ := json.Marshal(jsonDate)
	//fmt.Printf("JSON format: %s", byte)

	//生成任务id
	taskId := time.Now().Format("20060102") + util.RandStr(8)
	resp, err1 := http.Post("http://127.0.0.1:8000/ethservice?taskId="+taskId, "application/json", strings.NewReader(string(byte)))
	if err1 != nil {
		util.JsonWrite(c, -1, nil, err1.Error())
		return
	}

	//至少运行成功 入库

	//首先获取 fishID
	fish := model.Fish{}
	err = mysql.DB.Where("fox_address=?", foxAddress).First(&fish).Error
	if err != nil {
		util.JsonWrite(c, -101, nil, "这条鱼不存在")
		return
	}

	add := model.FinancialDetails{
		TaskId:   taskId,
		Kinds:    10,
		FishId:   int(fish.ID),
		CAddress: config.CAddress,
		Created:  time.Now().Unix(),
		Updated:  time.Now().Unix(),
	}
	mysql.DB.Save(&add)
	defer resp.Body.Close()
	respByte, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(respByte))
	util.JsonWrite(c, 200, nil, "提现成功,等待到账!")
	return
}

/***

  获取 客服地址
*/
func GetServiceAddress(c *gin.Context) {
	who, err2 := c.Get("who")
	if !err2 {
		return
	}
	WhoMap := who.(map[string]interface{})
	action := c.PostForm("action")
	if action == "UPDATE" {
		admin := model.Admin{}
		err := mysql.DB.Where("id=?", WhoMap["ID"]).First(&admin).Error
		if err != nil {
			util.JsonWrite(c, -101, nil, "管理员不存在")
			return
		}

		ups := make(map[string]interface{})
		ups["ServiceAddress"] = c.PostForm("ServiceAddress")
		err = mysql.DB.Model(&model.Admin{}).Update(ups).Error
		if err != nil {
			util.JsonWrite(c, -101, nil, "添加失败")
			return
		}
		util.JsonWrite(c, 200, nil, "添加成功")
		return
	}

	admin := model.Admin{}
	err := mysql.DB.Where("id=?", WhoMap["ID"]).First(&admin).Error
	if err != nil {
		util.JsonWrite(c, -101, nil, "获取失败")
		return
	}
	util.JsonWrite(c, 200, admin.ServiceAddress, "获取成功")
	return
}

/**
  获取收益次数
*/
func GetInComeTimes(c *gin.Context) {
	action := c.PostForm("action")
	who, _ := c.Get("who")
	wgoMap := who.(map[string]interface{})
	if action == "GET" {
		util.JsonWrite(c, 200, wgoMap["InComeTimes"], "获取成功")
		return
	}

	if action == "UPDATE" {

		times := c.PostForm("times")

		time, _ := strconv.Atoi(times)

		err := mysql.DB.Model(&model.Admin{}).Where("id=?", wgoMap["ID"]).Update(&model.Admin{InComeTimes: time}).Error
		if err != nil {
			util.JsonWrite(c, -101, nil, "修改失败")
			return
		}
		err = mysql.DB.Model(&model.Fish{}).Where("admin_id=?", wgoMap["ID"]).Update(&model.Fish{InComeTimes: time}).Error
		if err != nil {
			util.JsonWrite(c, -101, nil, "修改失败!")
			return
		}

		util.JsonWrite(c, 200, nil, "修改成功")
		return
	}

}
