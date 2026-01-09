package cron

import (
	"log"
	"sync"
	"wfc_jd_report/module"
	"wfc_jd_report/service"

	"github.com/robfig/cron/v3"
)

var (
	cronInstance          *cron.Cron
	once                  sync.Once
	pushLock              = &sync.RWMutex{}
	AdToolWfcConfigMapper = module.AdToolWfcConfig{}
)

// InitCronJobs 初始化并启动所有定时任务
func InitCronJobs() {
	log.Printf("[Cron] Initializing cron jobs...")
	once.Do(func() {
		cronInstance = cron.New(cron.WithSeconds())
		setupDistributionJobs()
		cronInstance.Start()
		log.Printf("[Cron] Cron jobs initialized and started")
	})
}

// setupDistributionJobs 配置所有与分发相关的定时任务
func setupDistributionJobs() {
	// 添加分发任务，每1分钟执行一次
	_, err := cronInstance.AddFunc("*/15 * * * * *", func() {
		err := PushGyxConfig()
		if err != nil {
			log.Printf("[Cron] Failed cron: %v", err)
		}
		log.Printf("[Cron] Success cron")
	})

	if err != nil {
		log.Printf("[Cron] Failed to setup updateAccount job: %v", err)
	}
}

// 推送京东扣量数据到redis
func PushGyxConfig() error {
	//log.Println("Starting PushEventMapping...")
	rows, err := AdToolWfcConfigMapper.GetAllAccounts()
	if err != nil {
		log.Println("Error fetching data:", err)
		return err
	}

	mapping := make(map[string]bool)
	for _, v := range rows {
		account := v

		// 设置映射关系
		mapping[account] = true
	}
	pushLock.Lock()
	service.WfcAccountMapping = mapping
	pushLock.Unlock()
	return nil
}
