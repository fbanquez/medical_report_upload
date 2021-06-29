package main

import (
	"github.com/spf13/viper"
)

// Configuration defines the structure of the YAML configuration file
type Configuration struct {
	Institution struct {
		Id   string `mapstructure:"id"`
	} `mapstructure:"institution"`
	Database struct {
		Host   string `mapstructure:"host"`
		Port   string `mapstructure:"port"`
		Db     string `mapstructure:"db"`
		User   string `mapstructure:"user"`
		Passwd string `mapstructure:"passwd"`
	} `mapstructure:"database"`
	Service struct {
		Uri      string `mapstructure:"uri"`
		Port     string `mapstructure:"port"`
		Endpoint string `mapstructure:"endpoint"`
		Auth     string `mapstructure:"auth"`
		Agent    string `mapstructure:"agent"`
		Active   bool   `mapstructure:"active"`
	} `mapstructure:"service"`
	SystemLog struct {
		Level int    `mapstructure:"level"`
		Path  string `mapstructure:"path"`
	} `mapstructure:"systemLog"`
	Proxy struct {
		Host   string `mapstructure:"host"`
		Port   string `mapstructure:"port"`
		User   string `mapstructure:"user"`
		Passwd string `mapstructure:"passwd"`
	} `mapstructure:"proxy"`
	Timer struct {
		GetProcedure int `mapstructure:"getprocedure"`
		UpProcedure  int `mapstructure:"upprocedure"`
		GetReport    int `mapstructure:"getreport"`
		UpReport     int `mapstructure:"upreport"`
	} `mapstructure:"timer"`
}

// config holds configurations
var config Configuration

// loadConfig parse the configuration file options inside the config variable
func loadConfigFile() (err error) {
	viper.AddConfigPath("./config")
	viper.SetConfigName("mru")
	viper.SetConfigType("yml")

	if err = viper.ReadInConfig(); err != nil {
		Error.Println("Unable to read config file. ", err)
		return
	}

	if err = viper.Unmarshal(&config); err != nil {
		Error.Println("Unable to decode configuration into struct. ", err)
		return
	}

	return
}

// init function set up some form of state on the program's startup
func init() {
	loadConfigFile()
	SetLogLevel()
}
