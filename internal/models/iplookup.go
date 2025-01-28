package models

import (
	"github.com/oschwald/maxminddb-golang"
	"log"
	"net"
)

type IPLookup struct {
	dbRegion *maxminddb.Reader
	dbASN    *maxminddb.Reader
}

func NewIPLookup(regionDBFilePath string, asnDBFilePath string) *IPLookup {

	regionDB, err := maxminddb.Open(regionDBFilePath)
	if err != nil {
		log.Fatal(err)
	}
	if err = regionDB.Verify(); err != nil {
		log.Fatal(err)
	}

	asnDB, err := maxminddb.Open(asnDBFilePath)
	if err != nil {
		log.Fatal(err)
	}
	if err = asnDB.Verify(); err != nil {
		log.Fatal(err)
	}

	return &IPLookup{
		dbRegion: regionDB,
		dbASN:    asnDB,
	}
}

func (d *IPLookup) LookupRegion(ip net.IP) (interface{}, error) {
	var record interface{}
	err := d.dbRegion.Lookup(ip, &record)
	if err != nil {
		return nil, err
	}
	return record, nil
}

func (d *IPLookup) LookupASN(ip net.IP) (interface{}, error) {
	var record interface{}
	err := d.dbASN.Lookup(ip, &record)
	if err != nil {
		return nil, err
	}
	return record, nil
}

func (d *IPLookup) Lookup(ip string) map[string]interface{} {
	ipAddress := net.ParseIP(ip)

	regionData, err := d.LookupRegion(ipAddress)
	if err != nil {
		log.Printf("Error looking up IP address %s: %s", ip, err)
	}

	asnData, err := d.LookupASN(ipAddress)
	if err != nil {
		log.Printf("Error looking up ASN for %s: %s", ip, err)
	}

	return map[string]interface{}{
		"region": regionData,
		"asn":    asnData,
	}
}
