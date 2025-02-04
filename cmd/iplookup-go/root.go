package main

import (
	"fmt"
	"github.com/fadhilyori/iplookup-go/internal/config"
	"github.com/fadhilyori/iplookup-go/internal/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"runtime"
)

var (
	appVersion = "dev"
	appCommit  = "none"
	appLicense = "MIT"
)

var log = logger.GetLogger()

var rootCmd = &cobra.Command{
	Use:   "iplookup-go",
	Short: "IP Lookup Service",
	Long:  `IP Lookup Service is a simple service that provides information about an IP address from MaxMind GeoLite2 database.`,
	Run:   runApp,
	Args:  cobra.NoArgs,
}

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	
	// Read configuration from .env file in the current directory
	// viper.SetConfigFile("./.env")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		log.WithField("error", err).Warnln("Failed to read configuration file, skipping.")
	}

	viper.AutomaticEnv()

	conf := config.GetConfig()

	viper.SetDefault("listen_address", "0.0.0.0")
	viper.SetDefault("listen_port", "3000")
	viper.SetDefault("mmdb_region_file_path", "")
	viper.SetDefault("mmdb_asn_file_path", "")
	viper.SetDefault("enable_cache", false)
	viper.SetDefault("redis_url", "")
	viper.SetDefault("cache_ttl_sec", 3600)
	viper.SetDefault("verbose", 0)

	if err := viper.Unmarshal(&conf); err != nil {
		log.Fatalf("Failed to unmarshal configuration: %v", err)
	}

	flags := rootCmd.PersistentFlags()

	flags.StringVar(&conf.ListenAddress, "host", conf.ListenAddress, "Host to listen on")
	flags.IntVar(&conf.ListenPort, "port", conf.ListenPort, "Port to listen on")
	flags.CountVarP(&conf.VerboseCount, "verbose", "v", "Increase verbosity of the output.")

	if err := viper.BindPFlags(flags); err != nil {
		log.WithField("error", err).Fatalln("Failed to bind flags.")
	}
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
