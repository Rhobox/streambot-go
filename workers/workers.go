package workers

import (
	"github.com/aricodes-oss/std"
	"sync"

	"streambot/discord"
)

var log = std.Logger
var discordClient = discord.Session

func LaunchAll(wg *sync.WaitGroup) {
	go StreamsWorker(wg)
	go CleanChannelsWorker(wg)
	go PostMessagesWorker(wg)
	go ScrubDBWorker(wg)
}
