/**
 * @Author $
 * @Description //TODO $
 * @Date $ $
 * @Param $
 * @return $
 **/
package util

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/shopspring/decimal"
	"math/big"

	"math/rand"
	"net/http"
	"time"
)

func RandStr(length int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	rand.Seed(time.Now().UnixNano() + int64(rand.Intn(100)))
	for i := 0; i < length; i++ {
		result = append(result, bytes[rand.Intn(len(bytes))])
	}
	return string(result)
}

func CreateToken(Rdb *redis.Client) string {

	for i := 0; i < 5; i++ {
		str := RandStr(36)
		_, err := Rdb.HGet("TOKEN_USER", str).Result()
		if err != nil {
			return str
		}
	}
	return ""
}

/**
判断 字符串是否在这个数组中
*/
func InArray(target string, strArray []string) bool {
	for _, element := range strArray {
		if target == element {
			return true
		}
	}
	return false
}

/**
  返回方法
*/
func JsonWrite(context *gin.Context, status int, result interface{}, msg string) {
	context.JSON(http.StatusOK, gin.H{
		"code":   status,
		"result": result,
		"msg":    msg,
	})
}

//func GetAccountMoneyUsdT(id int, foxAddress string, db *gorm.DB) {
//	resp, err := http.Get("https://etherscan.io/address/" + foxAddress)
//	if err != nil {
//		return
//	}
//	defer resp.Body.Close()
//	body, err1 := ioutil.ReadAll(resp.Body)
//	if err1 != nil {
//		return
//	}
//	//fmt.Println(string(body))
//	//解析正则表达式，如果成功返回解释器
//	reg1 := regexp.MustCompile(`<div class="col-md-8">\$(\d+)`)
//	if reg1 == nil { //解释失败，返回nil
//		return
//	}
//	//根据规则提取关键信息
//	result1 := reg1.FindAllStringSubmatch(string(body), -1)
//	maxMoney, err := strconv.ParseFloat(result1[0][1], 64)
//	up := model.Fish{
//		Money: maxMoney,
//	}
//	db.Model(&model.Fish{}).Where("id=?", id).Update(up)
//	return
//
//}

func ToDecimal(ivalue interface{}, decimals int) decimal.Decimal {
	value := new(big.Int)
	switch v := ivalue.(type) {
	case string:
		value.SetString(v, 10)
	case *big.Int:
		value = v
	}

	mul := decimal.NewFromFloat(float64(10)).Pow(decimal.NewFromFloat(float64(decimals)))
	num, _ := decimal.NewFromString(value.String())
	result := num.Div(mul)

	return result
}

//生成邀请码

func SetInCode(token string, Rdb *redis.Client) {

	//生成邀请码

}
