package service

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"

	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"github.com/takama/daemon"
	"github.com/virink/vulwarning/common"
	"github.com/virink/vulwarning/model"
	"github.com/virink/vulwarning/plugins"
)

// Service has embedded daemon
type Service struct {
	daemon.Daemon
}

// Manage -
func (service *Service) Manage() (string, error) {
	usage := fmt.Sprintf("Usage: %s install | remove | start | restart | stop | status | initdb | config", os.Args[0])
	if len(os.Args) > 1 {
		command := os.Args[1]
		switch command {
		case "install":
			return service.Install()
		case "remove":
			return service.Remove()
		case "start":
			return service.Start()
		case "stop":
			return service.Stop()
		case "restart":
			if x, err := service.Stop(); err != nil {
				return x, err
			}
			return service.Start()
		case "status":
			return service.Status()
		case "initdb":
			return initDb()
		case "config":
			return echoConfig()
		default:
			return usage, nil
		}
	}
	return serviceDaemon()
}

var (
	configFile string
	workDir    string
	config     = common.Conf
	logger     = common.Logger

	err error
)

func init() {
	debugEnv := os.Getenv("DEBUG")
	if debugEnv != "" && debugEnv != "0" && debugEnv != "false" {
		common.DebugMode = true
	}

	exePath, err := os.Executable()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	workDir = filepath.Dir(exePath)
	configFile = filepath.Join(workDir, common.ConfigFile)
}

// Entry -
func Entry() {
	kind := daemon.UserAgent
	if runtime.GOOS != "darwin" {
		kind = daemon.SystemDaemon
	}
	srv, err := daemon.New(common.ServiceName, common.Description, kind)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	service := &Service{srv}
	status, err := service.Manage()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	fmt.Fprintln(os.Stdout, status)
}

func echoConfig() (string, error) {
	if data := common.TemplateConfig(); data != nil {
		fmt.Println(string(data))
	}
	return "", nil
}

func baseInit() (string, error) {
	// Load Config
	if config, err = common.LoadConfig(configFile); err != nil {
		return "Load Config", err
	}
	level := logrus.InfoLevel
	if config.Server.Debug || common.DebugMode {
		common.DebugMode = true
		level = logrus.DebugLevel
	}

	// Init Logger
	logger = common.InitLogger(filepath.Join(workDir, common.LogFile), level)

	// Connect Database
	if _, err = model.InitConnect(config, common.DebugMode); err != nil {
		return "Connect Database", err
	}

	if config.Server.Migrate {
		model.AutoMigrate()
	}

	return "", nil
}

func initDb() (string, error) {

	if res, err := baseInit(); err != nil {
		return res, err
	}

	// Init Database
	model.InitTable()

	logger.Println("Crawl Vul But not push message in first time")
	plugins.DoJob(false)

	return "Init Database Success", nil
}

func serviceDaemon() (string, error) {

	if res, err := baseInit(); err != nil {
		return res, err
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, os.Kill, syscall.SIGTERM)

	c := cron.New()
	c.AddFunc(config.Server.Spec, func() {
		logger.Println("Start Job...")
		plugins.DoJob(true)
	})
	c.Start()

	sig := <-interrupt
	if sig == os.Interrupt {
		return "Daemon was interruped by system signal", nil
	}
	return "Daemon was killed", nil
}
