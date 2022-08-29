package workers

import (
	"sync"
	"time"

	"github.com/aricodes-oss/std"

	"streambot/channels"
	"streambot/discord"
)

type Entrypoint func()

type ScheduledTask struct {
	Run   Entrypoint
	Delay time.Duration
}

var log = std.Logger
var discordClient = discord.Session

var allWorkers = []*ScheduledTask{
	{Run: StreamsWorker, Delay: 15 * time.Second},
	{Run: CleanChannelsWorker},
	{Run: PostMessagesWorker, Delay: 20 * time.Second},
	{Run: ScrubDBWorker},
	{Run: CleanStaleMessages, Delay: 10 * time.Second},
}

func Launch(wg *sync.WaitGroup, task *ScheduledTask) {
	defer wg.Done()

	if task.Delay == 0 {
		task.Delay = 30 * time.Second
	}

	log.Debugf("Launching task %v", task)

	ticker := time.NewTicker(task.Delay)
	defer ticker.Stop()

	for {
		select {
		case <-channels.Running:
			return
		case <-ticker.C:
			task.Run()
		}
	}
}

func LaunchAll(wg *sync.WaitGroup) {
	wg.Add(len(allWorkers))

	for _, taskDef := range allWorkers {
		go Launch(wg, taskDef)
	}

	log.Infof("Launched %v background workers", len(allWorkers))
}
