package iplookup

import (
	"net"

	"github.com/fadhilyori/iplookup-go/internal/logger"
	"github.com/oschwald/maxminddb-golang"
)

var log = logger.GetLogger()

type IPLookup struct {
	db *maxminddb.Reader
}

func openDB(path string) *maxminddb.Reader {
	db, err := maxminddb.Open(path)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	if err := db.Verify(); err != nil {
		log.Fatalf("Failed to verify database: %v", err)
	}

	return db
}

func lookupIP(db *maxminddb.Reader, ip net.IP, data any) error {
	err := db.Lookup(ip, data)
	if err != nil {
		log.Errorf("Failed to lookup IP: %v", err)
		return err
	}

	return nil
}

func NewIPLookup(dbPath string) *IPLookup {
	return &IPLookup{
		db: openDB(dbPath),
	}
}

func (d *IPLookup) Lookup(ip net.IP, data any) error {
	return lookupIP(d.db, ip, data)
}
