package core

import "time"

type SshConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type DocConfig struct {
	Enabled bool   `mapstructure:"enabled"`
	Host    string `mapstructure:"host"`
	Port    int    `mapstructure:"port"`
	Dir     string `mapstructure:"dir"`
}

type PluginConfig struct {
	Dir string `mapstructure:"dir"`
}

type Config struct {
	Server struct {
		Interval time.Duration `mapstructure:"interval"`
    LogLevel string        `mapstructure:"log_level"`
		Plugin   PluginConfig  `mapstructure:"plugin"`
		Ssh      SshConfig     `mapstructure:"ssh"`
		Doc      DocConfig     `mapstructure:"doc"`
		Version  string
	} `mapstructure:"server"`
	Plugins map[string]map[string]string `mapstructure:"plugins"`
}
