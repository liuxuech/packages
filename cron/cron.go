package cron

import (
	"errors"
	"github.com/robfig/cron/v3"
)

var cronJob *LXC

// Cron
type LXC struct {
	entryIdMap map[string]cron.EntryID
	cron       *cron.Cron
}

func AddFunc(uid string, spec string, cmd func()) error {
	if entryId, err := cronJob.cron.AddFunc(spec, cmd); err == nil {
		return err
	} else {
		cronJob.entryIdMap[uid] = entryId
	}
	return nil
}

func AddJob(uid string, spec string, cmd cron.Job) error {
	if entryId, err := cronJob.cron.AddJob(spec, cmd); err == nil {
		return err
	} else {
		cronJob.entryIdMap[uid] = entryId
	}
	return nil
}

func Add(uid string, spec string, cmd interface{}) error {
	switch cmd.(type) {
	case cron.Job:
		return AddJob(uid, spec, cmd.(cron.Job))
	case func():
		return AddFunc(uid, spec, cmd.(func()))
	default:
		return errors.New("错误的类型")
	}
}

// 删除某个 Job
func Remove(uid string) {
	cronJob.cron.Remove(cronJob.entryIdMap[uid])
}

func init() {
	var lxc LXC
	lxc.entryIdMap = make(map[string]cron.EntryID)
	lxc.cron = cron.New(cron.WithSeconds())

	cronJob = &lxc

	// 启动一个goroutine执行任务
	lxc.cron.Start()
}
