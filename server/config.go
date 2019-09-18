package server

import (
	"fmt"
	"io/ioutil"
	"micro/util"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/golang/glog"
	"github.com/luci/go-render/render"
	yaml "gopkg.in/yaml.v2"
)

var (
	koalaConf = &KoalaConf{
		Port: 8080,
		Prometheus: PrometheusConf{
			SwitchOn: true,
			Port:     8081,
		},
		ServiceName: "koala_server",
		Register: RegisterConf{
			SwitchOn: false,
		},
		Log: LogConf{
			Level: "debug",
			Dir:   "./logs/",
		},
		Limit: LimitConf{
			SwitchOn: true,
			QPSLimit: 50000,
		},
	}
)

type LimitConf struct {
	QPSLimit int  `yaml:"qps"`
	SwitchOn bool `yaml:"switch_on"`
}

type KoalaConf struct {
	Port        int            `yaml:"port"`
	Prometheus  PrometheusConf `yaml:"prometheus"`
	ServiceName string         `yaml:"service_name"`
	Register    RegisterConf   `yaml:"register"`
	Log         LogConf        `yaml:"log"`
	Limit       LimitConf      `yaml:"limit"`

	ConfigDir  string `yaml:"-"`
	RootDir    string `yaml:"-"`
	ConfigFile string `yaml:"-"`
}

type PrometheusConf struct {
	SwitchOn bool `yaml:"switch_on"`
	Port     int  `yaml:"port"`
}

type RegisterConf struct {
	SwitchOn bool `yaml:"switch_on"`
}

type LogConf struct {
	Level string `yaml:"level"`
	Dir   string `yaml:"path"`
}

func initDir(serviceName string) (err error) {
	exeFilePath, err := filepath.Abs(os.Args[0])

	if err != nil {
		return
	}

	lastIdx := strings.LastIndex(exeFilePath, "/")
	if lastIdx < 0 {
		err = fmt.Errorf("invalid path :%v", exeFilePath)
		return
	}

	koalaConf.RootDir = path.Join(exeFilePath[0:lastIdx], "..")
	koalaConf.ConfigDir = path.Join(koalaConf.RootDir, "./conf/", util.GetEnv())
	koalaConf.ConfigFile = path.Join(koalaConf.ConfigDir, fmt.Sprintf("%s.yaml", serviceName))
	return
}

func InitConfig(serviceName string) (err error) {
	err = initDir(serviceName)
	if err != nil {
		return
	}

	data, err := ioutil.ReadFile(koalaConf.ConfigFile)
	if err != nil {
		return
	}

	err = yaml.Unmarshal(data, &koalaConf)
	if err != nil {
		return
	}

	glog.Info("init koala conf succ ", render.Render(koalaConf))
	return
}

func GetConfigDir() string {
	return koalaConf.ConfigDir
}

func GetRootDir() string {
	return koalaConf.RootDir
}

func GetServerPort() int {
	return koalaConf.Port
}

func GetConf() *KoalaConf {
	return koalaConf
}
