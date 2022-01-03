package setting

import (
	"github.com/go-ini/ini"
	"log"
)

type System struct {
	Source string
}

var cfg *ini.File
var SystemSetting = &System{}

func Setup() {
	var err error
	cfg, err = ini.Load("config.ini")
	if err != nil {
		log.Fatalf("setting.Setup, fail to parse 'config.ini': %v", err)
	}
	mapTo("system", SystemSetting)

}

// mapTo map section
func mapTo(section string, v interface{}) {
	err := cfg.Section(section).MapTo(v)
	if err != nil {
		log.Fatalf("Cfg.MapTo %s err: %v", section, err)
	}
}
