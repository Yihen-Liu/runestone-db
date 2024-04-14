package main

import (
	"context"
	"galactic-monitor/config"
	"galactic-monitor/log"
	"galactic-monitor/models"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var (
	Version   string
	BuildTime string
	mainCmd   = &cobra.Command{
		Use:   os.Args[0],
		Short: "inscribe nft/ft monitor service",
		Long: `Description:
			when buyer submit new order into db, monitor will catch it and distribute NFT/FT to received address;`,
		Run: entryPoint,
	}
	cfgFile string
)

func exitSignals() context.Context {
	signalsToCatch := []os.Signal{
		os.Interrupt,
		os.Kill,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	}
	interceptor := make(chan os.Signal, 1)
	signal.Notify(interceptor, signalsToCatch...)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-interceptor
		cancel()
	}()
	return ctx
}

func entryPoint(cmd *cobra.Command, args []string) {
	signal := exitSignals()
	go StartMonitorNewOrder(signal)
	cronTask()
	select {
	case <-signal.Done():
		log.Info("Goodbye monitor transaction service. Go Back Home.")

	}
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the mainCmd.
func main() {
	err := mainCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	mainCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c",
		"./conf.json", "./monitor-darwin-amd64 --config conf.json")

	cobra.OnInitialize(
		func() {
			//setting.Setup()
			models.Setup()
			//log.InitLog(setting.AppSetting.LogLevel, setting.AppSetting.LogSavePath, log.Stdout)
			log.InitLog(config.AppConf.Logger.LogLevel, config.AppConf.Logger.LogFileDir, log.Stdout)
			log.Infof("version:%s, build time:%s", Version, BuildTime)
		})
}
