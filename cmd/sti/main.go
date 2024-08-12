package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/willfantom/sti/pkg/influx"
	"github.com/willfantom/sti/pkg/iperf"
	"github.com/willfantom/sti/pkg/ping"
	"github.com/willfantom/sti/pkg/speedtest"
	"github.com/willfantom/sti/pkg/tester"
)

var (
	config Config

	rootCmd = &cobra.Command{
		Use:  "sti",
		Long: "Speed Test Influx (sti) performs a both an internet speed test and an iperf test to a specific site, reporting the results to an InfluxDB instance.",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// set config path to that set in flag
			if viper.GetString("config") != "" {
				viper.SetConfigFile(viper.GetString("config"))
			}

			// read in the config from file
			if err := viper.ReadInConfig(); err != nil {
				if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
					return fmt.Errorf("failed to read config file: %w", err)
				}
			}

			// put the config into a config struct
			if err := viper.Unmarshal(&config); err != nil {
				return fmt.Errorf("failed to parse config: %w", err)
			}

			// set verbose logging if flagged
			if config.Verbose {
				logrus.SetLevel(logrus.DebugLevel)
			} else {
				logrus.SetLevel(logrus.InfoLevel)
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			logrus.Infoln("starting sti")

			// make testers from configs
			tests := make([]tester.Test, 0)
			for _, speedtestConfig := range config.Speedtest {
				tests = append(tests, speedtest.New(speedtestConfig.ServerID))
			}
			for _, iperfConfig := range config.Iperf {
				tests = append(tests, iperf.New(iperfConfig.ServerIP, iperfConfig.ServerPort, iperfConfig.Streams, iperfConfig.Seconds, iperfConfig.TCP))
			}
			for _, pingConfig := range config.Ping {
				tests = append(tests, ping.New(pingConfig.Target, pingConfig.Count))
			}

			// exit if no tests?
			logrus.WithField("count", len(tests)).Infoln("configured tests")
			if len(tests) == 0 {
				logrus.WithError(fmt.Errorf("no tests configured")).Infoln("exiting")
				return nil
			}

			// shuffle tests (just for fun)
			rand.Shuffle(len(tests), func(i, j int) {
				tests[i], tests[j] = tests[j], tests[i]
			})

			// run
			for {
				for _, test := range tests {
					logrus.
						WithField("test", test.Name()).
						WithFields(test.Config()).
						Infoln("running test")
					testResult, err := test.RunTest() // run the test
					if err != nil {
						logrus.
							WithField("test", test.Name()).
							WithError(err).
							Errorln("failed to run test")
					} else {
						logrus.
							WithFields(testResult.Data).
							Infoln("test complete")
						logrus.
							WithField("url", config.InfluxURL).
							WithField("org", config.InfluxOrg).
							WithField("bucket", config.InfluxBucket).
							WithField("host", config.InfluxHostname).
							Infoln("writing data to influx")
						if err := influx.WriteData( // write the data to influx
							config.InfluxURL,
							config.InfluxOrg,
							config.InfluxBucket,
							config.InfluxToken,
							config.InfluxHostname,
							test.Name(),
							testResult,
						); err != nil {
							logrus.
								WithField("test", test.Name()).
								WithError(err).
								Errorln("failed to write data to influx")
						} else {
							logrus.
								WithField("test", test.Name()).
								Infoln("data written to influx")
						}
					}
					logrus.WithField("seconds", config.Interval.Seconds()).Infoln("pausing...")
					<-time.After(config.Interval)
				}
			}
		},
	}
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		logrus.WithError(err).Fatalln("sti failed to execute")
	}
}

func init() {
	rootCmd.PersistentFlags().String("config", "", "config file (default is $HOME/.config/sti/config.yaml)")
	viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "output debug logs")
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
}
