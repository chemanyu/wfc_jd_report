package module

import (
	mysqldb "wfc_jd_report/common/mysql"
)

// AdToolWfcConfig WFC配置表结构
type AdToolWfcConfig struct {
	ID         uint   `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Account    string `gorm:"column:account;type:varchar(64);default:'';comment:微方程需要的账户" json:"account"`
	DelFlag    int8   `gorm:"column:del_flag;type:tinyint(1);default:0;comment:删除(0:正常;1:删除)" json:"del_flag"`
	CreateAt   int    `gorm:"column:create_at;default:0;comment:创建人" json:"create_at"`
	CreateTime int    `gorm:"column:create_time;default:0;comment:创建时间" json:"create_time"`
	UpdateAt   int    `gorm:"column:update_at;default:0;comment:修改人" json:"update_at"`
	UpdateTime int    `gorm:"column:update_time;default:0;comment:修改时间" json:"update_time"`
}

// TableName 设置表名
func (AdToolWfcConfig) TableName() string {
	return "ad_tool_wfc_config"
}

// AdToolWfcConfigMapper 数据映射器
var AdToolWfcConfigMapper = AdToolWfcConfig{}

// FindAll 查找所有未删除的记录
func (a *AdToolWfcConfig) FindAll() ([]AdToolWfcConfig, error) {
	db := mysqldb.GetConnected()
	var configs []AdToolWfcConfig
	err := db.Where("del_flag = 0").Find(&configs).Error
	if err != nil {
		return nil, err
	}
	return configs, nil
}

// GetAllAccounts 获取所有有效的账户列表
func (a *AdToolWfcConfig) GetAllAccounts() ([]string, error) {
	configs, err := a.FindAll()
	if err != nil {
		return nil, err
	}

	accounts := make([]string, 0, len(configs))
	for _, config := range configs {
		if config.Account != "" {
			accounts = append(accounts, config.Account)
		}
	}
	return accounts, nil
}
