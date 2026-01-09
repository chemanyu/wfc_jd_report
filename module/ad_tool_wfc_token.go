package module

import (
	"errors"
	mysqldb "wfc_jd_report/common/mysql"
)

// AdToolWfcToken WFC Token表结构
type AdToolWfcToken struct {
	ID          uint   `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UnionId     string `gorm:"column:union_id;type:varchar(64);default:'';comment:京东联盟ID" json:"union_id"`
	MsAppkey    string `gorm:"column:ms_appkey;type:varchar(64);default:'';comment:美数给到wfc的appKey" json:"ms_appkey"`
	JdMediaId   string `gorm:"column:jd_media_id;type:varchar(64);default:'';comment:京东流量媒体ID" json:"jd_media_id"`
	JdAppkey    string `gorm:"column:jd_appkey;type:varchar(64);default:'';comment:流量媒体对应appKey" json:"jd_appkey"`
	JdSecretkey string `gorm:"column:jd_secretkey;type:varchar(64);default:'';comment:流量媒体对应secretKey" json:"jd_secretkey"`
	DelFlag     int8   `gorm:"column:del_flag;type:tinyint(1);default:0;comment:删除(0:正常;1:删除)" json:"del_flag"`
	CreateAt    int    `gorm:"column:create_at;default:0;comment:创建人" json:"create_at"`
	CreateTime  int    `gorm:"column:create_time;default:0;comment:创建时间" json:"create_time"`
	UpdateAt    int    `gorm:"column:update_at;default:0;comment:修改人" json:"update_at"`
	UpdateTime  int    `gorm:"column:update_time;default:0;comment:修改时间" json:"update_time"`
}

// TableName 设置表名
func (AdToolWfcToken) TableName() string {
	return "ad_tool_wfc_token"
}

// AdToolWfcTokenMapper 数据映射器
var AdToolWfcTokenMapper = AdToolWfcToken{}

// FindByMsAppkey 根据美数appkey查找记录（核心方法）
func (a *AdToolWfcToken) FindByMsAppkey(msAppkey string) (*AdToolWfcToken, error) {
	if msAppkey == "" {
		return nil, errors.New("ms_appkey不能为空")
	}

	db := mysqldb.GetConnected()
	var token AdToolWfcToken
	err := db.Where("ms_appkey = ? AND del_flag = 0", msAppkey).First(&token).Error
	if err != nil {
		return nil, err
	}
	return &token, nil
}

// GetJdCredentialsByMsAppkey 通过美数appkey查询京东appkey和secretkey
func (a *AdToolWfcToken) GetJdCredentialsByMsAppkey(msAppkey string) (jdAppkey, jdSecretkey string, err error) {
	token, err := a.FindByMsAppkey(msAppkey)
	if err != nil {
		return "", "", err
	}

	return token.JdAppkey, token.JdSecretkey, nil
}
