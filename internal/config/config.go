package config

import (
	"sync"

	"github.com/fadhilyori/iplookup-go/internal/logger"
)

type Config struct {
	ListenAddress      string `mapstructure:"listen_address"`
	ListenPort         int    `mapstructure:"listen_port"`
	MMDBRegionFilePath string `mapstructure:"mmdb_region_file_path"`
	MMDBASNFilePath    string `mapstructure:"mmdb_asn_file_path"`
	EnableCache        bool   `mapstructure:"enable_cache"`
	RedisURL           string `mapstructure:"redis_url"`
	CacheTTLSec        int    `mapstructure:"cache_ttl_sec"`
	VerboseCount       int    `mapstructure:"verbose"`
}

var log = logger.GetLogger()

var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		instance = &Config{}
	})

	return instance
}

func (c *Config) SetupLogging() {
	switch instance.VerboseCount {
	case 0:
		log.SetLevel(logger.InfoLevel)
	case 1:
		log.SetLevel(logger.DebugLevel)
	case 2:
		log.SetLevel(logger.TraceLevel)
	default:
		log.SetLevel(logger.TraceLevel)
	}
	log.WithFields(logger.Fields{
		"LOG_LEVEL": log.GetLevel().String(),
	}).Infoln("Logging level set.")
}
