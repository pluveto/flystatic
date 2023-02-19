package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/alexflint/go-arg"
	"github.com/pluveto/flystatic/cmd/flystatic/app"
	"github.com/pluveto/flystatic/cmd/flystatic/conf"
	"github.com/pluveto/flystatic/pkg/logger"
	"github.com/sirupsen/logrus"
)

// main Entry point of the application
func main() {
	// main Entry point of the application
	var args app.Args
	var cnf conf.Conf
	var defaultConf = conf.GetDefaultConf()

	args = loadArgsValid()
	cnf = loadConfValid(args.Config, defaultConf, "config.toml")
	overrideConf(&cnf, args)
	app.InitLogger(cnf.Log, args.Verbose)
	logger.Debug("log level: ", logger.GetLevel())
	app.Run(cnf)
}

func overrideConf(cnf *conf.Conf, args app.Args) {
	if args.Verbose {
		cnf.Log.Level = logrus.DebugLevel.String()
	}
	if args.Port != 0 {
		cnf.Server.Port = args.Port
	}
	if args.Host != "" {
		cnf.Server.Host = args.Host
	}
}

func loadArgsValid() app.Args {
	var args app.Args
	arg.MustParse(&args)
	return args
}

func getAppDir() string {
	dir, err := os.Executable()
	if err != nil {
		logger.Fatal(err)
	}
	return filepath.Dir(dir)
}

func loadConfValid(path string, defaultConf conf.Conf, defaultConfPath string) conf.Conf {
	if path == "" {
		path = defaultConfPath
	}
	// app executable dir + config.toml has the highest priority
	preferredPath := filepath.Join(getAppDir(), path)
	if _, err := os.Stat(preferredPath); err == nil {
		path = preferredPath
	}
	_, err := toml.DecodeFile(path, &defaultConf)
	if err != nil {
		fmt.Println("failed to load config file: (" + err.Error() + ") using default config")
	}
	logger.WithField("conf", &defaultConf).Debug("configuration loaded")
	return defaultConf
}
