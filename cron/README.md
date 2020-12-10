# cron

* 安装
> go get github.com/robfig/cron/v3

* Run 和 Start 的区别
    * Run 阻塞执行
    * 开启一个goroutine执行任务

* 固定时间间隔
    * @every 1s 每秒触发一下
    * @every 1m 每分钟触发一下
    * @every 1h 每小时触发一下
    * @every 1m2s 每1分2秒触发一下
> time.ParseDuration 支持的格式都可以用在这里

* 预定义时间规则
    * @hourly   每小时的0分0秒触发
    * @daily    每天的0时0分0秒触发
    * @weekly   每周的第一天的0时0分0秒触发，注意：每周的第一天是周日
    * @monthly  每月的第一天的0时0分0秒触发
    * @yearly   每年的第一个月的第一天的0时0分0秒触发
    
* 时间格式
> go版本的cron任务管理器默认支持用 5 个空格分隔的域来表示时间。如果要支持用 6 个空格分隔的域来表示时间，则需要添加如下代码：

```go
c := cron.New(cron.WithSeconds())
```

* 时间格式顺序：
> 秒(0~59) 分(0~59) 小时(0~23) 天(1~31) 月(1~12) 周(0~6)

* cron特定符号说明
    * 星号（*） - 表示可以匹配任意值
    * 斜线（/） - 表示步长，秒（*/3）表示秒数为0~59内间隔3秒触发一次；秒（20-59/2）表示在秒数在20~59范围内，间隔2秒触发一次。(20-59/2) 等价于 (20/2)。
    * 逗号（,） - 用于指定多个值，天（1,3,5）表示1号、3号、5号。
    * 连接符（-） - 表示范围
    * 问号（?） - 只用于 天 和 周，表示任意值，可以用于代替 * 。
    
* 示例
```go
package main

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"time"
)

func main() {

	var flag bool

	c := cron.New(cron.WithSeconds())
	
	addFunc, _ := c.AddFunc("*/5 * * * * *", func() {
		fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
		flag = true
	})

	go func() {
		for {
			if flag {
				time.Sleep(time.Second * 2)
				entry := c.Entry(addFunc)

				// 删除指定Job
				c.Remove(addFunc)

				// 获取某个Job的Entry
				c.Entry(addFunc)

				// 添加Job，函数方式
				c.AddFunc()

				// 可以自己实现 Job 接口，这样在执行的Job的时候可以有状态
				c.AddJob()

				// 返回时区信息
				fmt.Println(c.Location().String())

				// 返回所有的Entries
				c.Entries()

				// 启动任务管理器
				c.Start()

				// 返回任务上一次执行的实际
				fmt.Println(entry.Prev.Format("2006-01-02 15-04-05"))
				// 返回这个任务下一次执行的时间
				fmt.Println(entry.Next.Format("2006-01-02 15-04-05"))
				// 立即执行任务
				entry.Job.Run()
				// 检查任务是否正常
				fmt.Println(entry.Valid())
				// 返回任务 EntryId
				fmt.Println(entry.ID)
			}
		}
	}()

	c.Start()
	select {}
}
```
    
    
    
    
    
    
    
    
    
    
    