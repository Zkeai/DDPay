package cron

import (
	"context"
	"github.com/Zkeai/go_template/common/logger"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
)

type CronService struct {
	c           *cron.Cron
	ctx         context.Context
	cancel      context.CancelFunc
	once        sync.Once
	taskManager *TaskManager
}

func NewCronService() *CronService {
	ctx, cancel := context.WithCancel(context.Background())
	return &CronService{
		ctx:         ctx,
		cancel:      cancel,
		taskManager: GetManager(), // 使用单例的 TaskManager
	}
}

// Start 启动 Cron 服务
func (s *CronService) Start() {
	logger.Info("CronService starting...")
	s.taskManager.Ctx = s.ctx
	s.c = cron.New(cron.WithSeconds())
	s.c.Start()
}

// Stop 停止 Cron 服务
func (s *CronService) Stop() {
	s.once.Do(func() {
		logger.Info("CronService stopping...")
		s.cancel()
		s.taskManager.Stop() // 停止所有任务
		s.c.Stop()           // 停止 cron 调度器
	})
}

// AddTask 动态添加任务
func (s *CronService) AddTask(spec string, cmd func(context.Context)) (cron.EntryID, error) {
	return s.taskManager.AddTask(spec, cmd)
}

// AddTaskOnce 动态添加只执行一次的任务
func (s *CronService) AddTaskOnce(d time.Duration, cmd func(context.Context)) int {
	return s.taskManager.AddTaskOnce(d, cmd)
}
