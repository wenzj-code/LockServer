package Config

import (
	"flag"
	"io/ioutil"
	"os"

	log "github.com/Sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
)

var configOpt Option

//Option all the config
type Option struct {
	Addr string `yaml:"Addr"`

	LogFile    string `yaml:"LogFile"`
	LogLevel   string `yaml:"LogLevel"`
	SysLogAddr string `yaml:"SysLogAddr"`

	
}

func InitConfig() {
	configOpt = loadConfig()
}

//GetConfig create option
func GetConfig() *Option {
	return &configOpt
}

func loadConfig() (p Option) {
	var confName string
	flag.StringVar(&confName, "f", "config.yml", "config file of monitor")
	flag.Parse()
	f, err := os.Open(confName)
	data, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatal("read config file err: " + err.Error())
	}
	err = yaml.Unmarshal(data, &p)
	if err != nil {
		log.Fatal("unmarshal yaml config err: " + err.Error())
	}
	return p
}
