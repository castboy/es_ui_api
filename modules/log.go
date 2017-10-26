package modules

import (
	"os"

	seelog "github.com/cihub/seelog"
)

func InitLog() {
	logger, err := seelog.LoggerFromConfigAsFile(os.Getenv("GOPATH") + "/src/github.com/castboy/es_ui_api/modules/seelog.xml")

	if err != nil {
		seelog.Critical("err parsing config log file", err)
		return
	}
	seelog.ReplaceLogger(logger)
}

func Log(level string, content string) {
	defer seelog.Flush()

	switch level {
	case "Err":
		seelog.Error(content)
	case "Debug":
		seelog.Debug(content)
	case "Info":
		seelog.Info(content)
	}
}
