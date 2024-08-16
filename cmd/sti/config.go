package main

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Speedtest []struct {
		ServerID string `mapstructure:"serverID"`
	} `mapstructure:"speedtest"`
	Iperf []struct {
		ServerIP   string `mapstructure:"serverIP"`
		ServerPort int    `mapstructure:"serverPort"`
		Streams    int    `mapstructure:"streams"`
		Seconds    int    `mapstructure:"seconds"`
		TCP        bool   `mapstructure:"tcp"`
		Bandwidth  string `mapstructure:"bandwidth"`
		Reverse    bool   `mapstructure:"reverse"`
	} `mapstructure:"iperf"`
	Ping []struct {
		Target string `mapstructure:"target"`
		Count  int    `mapstructure:"count"`
	} `mapstructure:"ping"`

	InfluxURL      string `mapstructure:"influxURL"`
	InfluxToken    string `mapstructure:"influxToken"`
	InfluxOrg      string `mapstructure:"influxOrg"`
	InfluxBucket   string `mapstructure:"influxBucket"`
	InfluxHostname string `mapstructure:"influxHostname"`

	Interval time.Duration `mapstructure:"interval"`
	Verbose  bool          `mapstructure:"verbose"`
}

var (
	defaultConfig = map[string]any{
		"speedtest": []map[string]any{},
		"iperf":     []map[string]any{},
		"ping":      []map[string]any{},

		"influxURL":      "http://localhost:8086",
		"influxToken":    "",
		"influxOrg":      "",
		"influxBucket":   "",
		"influxHostname": "",
		"interval":       60 * time.Second,
		"verbose":        false,
	}
)

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.config/sti")
	viper.AddConfigPath("/etc/sti")
	viper.AutomaticEnv()
	viper.SetDefault("speedtest", defaultConfig["speedtest"])
	viper.SetDefault("iperf", defaultConfig["iperf"])
	viper.SetDefault("ping", defaultConfig["ping"])
	viper.SetDefault("influxURL", defaultConfig["influxURL"])
	viper.SetDefault("influxToken", defaultConfig["influxToken"])
	viper.SetDefault("influxOrg", defaultConfig["influxOrg"])
	viper.SetDefault("influxBucket", defaultConfig["influxBucket"])
	viper.SetDefault("influxHostname", defaultConfig["influxHostname"])
	viper.SetDefault("interval", defaultConfig["interval"])
}
