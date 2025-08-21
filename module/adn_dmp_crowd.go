package module

import (
	mysqldb "dmp_distribution/common/mysql"
	"time"
)

// AdnDmpCrowd 人群表结构
type AdnDmpCrowd struct {
	ID              int       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	CrowdID         int       `gorm:"column:crowd_id;type:varchar(50);not null;comment:人群ID" json:"crowd_id"`
	CrowdName       string    `gorm:"column:crowd_name;type:varchar(100);not null;comment:人群名称" json:"crowd_name"`
	Desc            string    `gorm:"column:desc;type:text;comment:人群描述" json:"desc"`
	InvolveMember   int       `gorm:"column:involve_member;type:int;default:0;comment:包含人数" json:"involve_member"`
	CreateID        int       `gorm:"column:create_id;type:int;comment:创建者ID" json:"create_id"`
	CreateName      string    `gorm:"column:create_name;type:varchar(100);comment:创建者名称" json:"create_name"`
	CrowdCreateTime int       `gorm:"column:crowd_create_time;type:int;default:0;comment:人群包创建时间" json:"crowd_create_time"`
	CrowdUpdateTime int       `gorm:"column:crowd_update_time;type:int;default:0;comment:人群包更新时间" json:"crowd_update_time"`
	CreateTime      time.Time `gorm:"column:create_time;type:datetime;default:CURRENT_TIMESTAMP;comment:创建时间" json:"create_time"`
	UpdateTime      time.Time `gorm:"column:update_time;type:datetime;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;comment:更新时间" json:"update_time"`
}

// TableName 设置表名
func (AdnDmpCrowd) TableName() string {
	return "adn_dmp_crowd"
}

// AdnDmpCrowdMapper 数据映射器
var AdnDmpCrowdMapper = AdnDmpCrowd{}

// Insert 插入单条记录
func (a *AdnDmpCrowd) Insert(crowd AdnDmpCrowd) error {
	db := mysqldb.GetAdnConnected()
	return db.Create(&crowd).Error
}

// BatchInsert 批量插入记录
func (a *AdnDmpCrowd) BatchInsert(crowds []AdnDmpCrowd) error {
	if len(crowds) == 0 {
		return nil
	}

	db := mysqldb.GetAdnConnected()
	return db.CreateInBatches(crowds, 1000).Error // 每批1000条记录
}

// Update 更新记录
func (a *AdnDmpCrowd) Update() error {
	db := mysqldb.GetAdnConnected()
	return db.Save(a).Error
}

// UpdateByID 根据ID更新记录
func (a *AdnDmpCrowd) UpdateByID(id int, updates map[string]interface{}) error {
	db := mysqldb.GetAdnConnected()
	return db.Model(&AdnDmpCrowd{}).Where("id = ?", id).Updates(updates).Error
}

// UpdateByCrowdID 根据人群ID更新记录
func (a *AdnDmpCrowd) UpdateByCrowdID(crowdID string, updates map[string]interface{}) error {
	db := mysqldb.GetAdnConnected()
	return db.Model(&AdnDmpCrowd{}).Where("crowd_id = ?", crowdID).Updates(updates).Error
}

// FindByID 根据ID查找记录
func (a *AdnDmpCrowd) FindByID(id int) (*AdnDmpCrowd, error) {
	db := mysqldb.GetAdnConnected()
	var crowd AdnDmpCrowd
	err := db.Where("id = ? AND is_deleted = 0", id).First(&crowd).Error
	if err != nil {
		return nil, err
	}
	return &crowd, nil
}

// FindByCrowdID 根据人群ID查找记录
func (a *AdnDmpCrowd) FindByCrowdID(crowdID string) (*AdnDmpCrowd, error) {
	db := mysqldb.GetAdnConnected()
	var crowd AdnDmpCrowd
	err := db.Where("crowd_id = ? AND is_deleted = 0", crowdID).First(&crowd).Error
	if err != nil {
		return nil, err
	}
	return &crowd, nil
}
