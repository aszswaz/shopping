package config

import (
	"flag"
	"os"
	"path"
)

var optCache *Options

type Options struct {
	ConfigFile string
	Address    string
	Port       uint16
	Static     string
}

func GetOptions() *Options {
	if optCache != nil {
		return optCache
	}

	optCache = new(Options)
	flag.StringVar(&optCache.ConfigFile, "c", "", "yml configuration file path. (default shopping.xml in the working directory)")
	flag.StringVar(&optCache.Address, "a", "", "bind address. (default 127.0.0.1)")
	port := flag.Uint("p", 0, "Bind tcp port. (default 8080)")
	flag.StringVar(&optCache.Static, "s", "", "static file directory. (default static directory in the working directory.)")
	flag.Parse()

	optCache.Port = uint16(*port)
	optCache.setDefault()
	return optCache
}

func (opt *Options) setDefault() {
	if opt.ConfigFile == "" {
		if cwd, err := os.Getwd(); err == nil {
			opt.ConfigFile = path.Join(cwd, "shopping.yml")
		}
	}
}
