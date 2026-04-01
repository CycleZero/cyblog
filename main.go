package main

import (
	"cyblog/conf"
	"cyblog/pkg/log"
	"os"
	"os/signal"

	"go.uber.org/zap"
)

// TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>
func main() {

	vc := conf.GetConfig()
	logger := log.GetLogger()

	app := initApp(vc, logger)
	done := make(chan os.Signal)
	go func() {
		defer func() {
			done <- os.Interrupt
		}()
		logger.Info("服务已启动")
		err := app.StartServer()
		if err != nil {
			logger.Error("服务崩溃", zap.Error(err))
			return
		}
	}()

	signal.Notify(done, os.Interrupt)
	<-done
	logger.Info("服务退出")
	return
}
