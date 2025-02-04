package config

import (
	"github.com/fadhilyori/iplookup-go/internal/logger"
	"testing"
)

func TestConfig_SetupLogging(t *testing.T) {
	type fields struct {
		ListenAddress      string
		ListenPort         int
		MMDBRegionFilePath string
		MMDBASNFilePath    string
		EnableCache        bool
		RedisURL           string
		CacheTTLSec        int
		VerboseCount       int
	}
	tests := []struct {
		name   string
		fields fields
		want   logger.Level
	}{
		{
			name: "Test SetupLogging, must set log level to debug",
			fields: fields{
				ListenAddress:      "0.0.0.0",
				ListenPort:         3000,
				MMDBRegionFilePath: "/path/to/GeoLite2-City.mmdb",
				MMDBASNFilePath:    "/path/to/GeoLite2-ASN.mmdb",
				EnableCache:        false,
				RedisURL:           "",
				CacheTTLSec:        3600,
				VerboseCount:       2,
			},
			want: logger.DebugLevel,
		},
		{
			name: "Test SetupLogging, must set log level to trace",
			fields: fields{
				ListenAddress:      "0.0.0.0",
				ListenPort:         3000,
				MMDBRegionFilePath: "/path/to/GeoLite2-City.mmdb",
				MMDBASNFilePath:    "/path/to/GeoLite2-ASN.mmdb",
				EnableCache:        false,
				RedisURL:           "",
				CacheTTLSec:        3600,
				VerboseCount:       3,
			},
			want: logger.TraceLevel,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := GetConfig()
			c.VerboseCount = tt.fields.VerboseCount
			c.SetupLogging()

			logInstance := logger.GetLogger()

			if logInstance.GetLevel() != tt.want {
				t.Errorf("SetupLogging() got = %v, want = %v", logInstance.GetLevel(), tt.want)
			}
		})
	}
}
