/*
@Time : 29/1/2021 公元 16:10
@Author : philiphu
@File : schedule_task
@Software: GoLand
*/
package common

import (
	"errors"
	"github.com/robfig/cron"
)

type QyTask struct {
	//任务名称
	Name string
	//需要执行的任务
	Run func()
	//任务执行的周期配置，如，每天凌晨一点执行一次 0 0 1 * * ？
	Spec string
}

func StartSchedule(tasks []*QyTask, async bool) error {
	//ctx := context.Background()
	if len(tasks) == 0 {
		//log.Errorf(ctx, "start time task,but no task")
		return errors.New("no time task")
	}
	c := cron.New()
	for _, t := range tasks {
		//log.Errorf(ctx, "start time task")
		if err := c.AddFunc(t.Spec, t.Run); err != nil {
			//log.Errorf(ctx, "start time task,but add task error error is %v",err)
		}
	}
	if async {
		c.Start()
	} else {
		c.Run()
	}
	return nil
}
