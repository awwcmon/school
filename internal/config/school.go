// code generated by https://github.com/go-dev-frame/sponge

package config

import (
	"github.com/go-dev-frame/sponge/pkg/conf"
)

var config *Config

func Init(configFile string, fs ...func()) error {
	config = &Config{}
	return conf.Parse(configFile, config, fs...)
}

func Show(hiddenFields ...string) string {
	return conf.Show(config, hiddenFields...)
}

func Get() *Config {
	if config == nil {
		panic("config is nil, please call config.Init() first")
	}
	return config
}

func Set(conf *Config) {
	config = conf
}

type Config struct {
	App      App      `yaml:"app" json:"app"`
	Database Database `yaml:"database" json:"database"`
	HTTP     HTTP     `yaml:"http" json:"http"`
	Jaeger   Jaeger   `yaml:"jaeger" json:"jaeger"`
	Logger   Logger   `yaml:"logger" json:"logger"`
	Moroto   Moroto   `yaml:"moroto" json:"moroto"`
	Redis    Redis    `yaml:"redis" json:"redis"`
	Snailjob Snailjob `yaml:"snailjob" json:"snailjob"`
}

type Jaeger struct {
	AgentHost string `yaml:"agentHost" json:"agentHost"`
	AgentPort int    `yaml:"agentPort" json:"agentPort"`
}

type App struct {
	CacheType            string  `yaml:"cacheType" json:"cacheType"`
	EnableCircuitBreaker bool    `yaml:"enableCircuitBreaker" json:"enableCircuitBreaker"`
	EnableHTTPProfile    bool    `yaml:"enableHTTPProfile" json:"enableHTTPProfile"`
	EnableLimit          bool    `yaml:"enableLimit" json:"enableLimit"`
	EnableMetrics        bool    `yaml:"enableMetrics" json:"enableMetrics"`
	EnableStat           bool    `yaml:"enableStat" json:"enableStat"`
	EnableTrace          bool    `yaml:"enableTrace" json:"enableTrace"`
	Env                  string  `yaml:"env" json:"env"`
	Host                 string  `yaml:"host" json:"host"`
	Name                 string  `yaml:"name" json:"name"`
	TracingSamplingRate  float64 `yaml:"tracingSamplingRate" json:"tracingSamplingRate"`
	Version              string  `yaml:"version" json:"version"`
}

type Redis struct {
	DialTimeout  int    `yaml:"dialTimeout" json:"dialTimeout"`
	Dsn          string `yaml:"dsn" json:"dsn"`
	ReadTimeout  int    `yaml:"readTimeout" json:"readTimeout"`
	WriteTimeout int    `yaml:"writeTimeout" json:"writeTimeout"`
}

type Moroto struct {
	DocsPath string `yaml:"docsPath" json:"docsPath"`
	Enable   bool   `yaml:"enable" json:"enable"`
}

type Database struct {
	Driver  string  `yaml:"driver" json:"driver"`
	Mongodb Mongodb `yaml:"mongodb" json:"mongodb"`
}

type Mongodb struct {
	Dsn string `yaml:"dsn" json:"dsn"`
}

type Logger struct {
	Format string `yaml:"format" json:"format"`
	IsSave bool   `yaml:"isSave" json:"isSave"`
	Level  string `yaml:"level" json:"level"`
}

type Snailjob struct {
	GroupName    string `yaml:"GroupName" json:"GroupName"`
	HostIP       string `yaml:"HostIP" json:"HostIP"`
	HostPort     string `yaml:"HostPort" json:"HostPort"`
	Level        string `yaml:"Level" json:"Level"`
	Namespace    string `yaml:"Namespace" json:"Namespace"`
	ReportCaller bool   `yaml:"ReportCaller" json:"ReportCaller"`
	ServerHost   string `yaml:"ServerHost" json:"ServerHost"`
	ServerPort   string `yaml:"ServerPort" json:"ServerPort"`
	Token        string `yaml:"Token" json:"Token"`
}

type HTTP struct {
	Port    int `yaml:"port" json:"port"`
	Timeout int `yaml:"timeout" json:"timeout"`
}
