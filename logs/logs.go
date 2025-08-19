package logs

import (
	"bankSystem/types"
	"os"
	"time"
)

var (
	LogCh    = make(chan *types.Log, 100)
	logTypes = map[int]string{
		0: "the service is opened",
		1: "fatal shutdown",
		2: "account created",
		3: "account login",
		4: "account logout",
		5: "account deposit",
		6: "account withdraw",
		7: "account transfer money",
	}
)

func CreateLogObjects(kind int, username string) {
	var newLog *types.Log
	if kind < 0 || kind >= len(logTypes) {
		return
	}
	newLog = &types.Log{
		LogType:  logTypes[kind],
		Date:     time.Now().Format("2006-01-02 15:04:05"),
		Username: username,
	}
	LogCh <- newLog
}

func Logs() {
	go WriteLogsIntoLogFile(LogCh)
	go ClearLogs()
}

func WriteLogsIntoLogFile(ch chan *types.Log) {
	f, err := os.OpenFile(".log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer func() { _ = f.Close() }()
	for v := range ch {
		_, err = f.WriteString(v.LogType + " " + v.Date + " " + v.Username + "\n")
		if err != nil {
			return
		}
	}
}

func ClearLogs() {
	tick := time.NewTicker(time.Hour * 24)
	for {
		select {
		case <-tick.C:
			_ = os.Truncate(".log", 0)
		}
	}
}
