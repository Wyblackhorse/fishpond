/**
 * @Author $
 * @Description //TODO $
 * @Date $ $
 * @Param $
 * @return $
 **/
package model

/**
  每日统计数据
*/
type everydayData struct {
	ID            uint    `gorm:"primaryKey;comment:'主键'"`
	Created       int64   //创建时间
	AdminId       int     `gorm:"int(11);comment:'属于那个代理';index"`
	Date          string  //日期
	RegisterCount int     //注册个数
	TiXianCount   int     //提现个数
	Authorization int     //授权个数
	ChouQuMoney   float64 `gorm:"type:decimal(10,2)"` //抽取金额
	TiXianMoney   float64 `gorm:"type:decimal(10,2)"` //提现金额

}
