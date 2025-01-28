package models

import (
	"log"
	"net"

	"github.com/oschwald/maxminddb-golang"
)

type Lookup interface {
	IpLookupRegion(fileLocation string) interface{}
}

type Reader struct {
	*maxminddb.Reader
}

func NewReaderHandler(fileLocation string) (*Reader, error) {
	log.Print("Configuring Maxmind Database")
	reader, err := maxminddb.Open(fileLocation)
	if err != nil {
		log.Print("Maxmind Database Not Found")
		return nil, err
	}
	if err = reader.Verify(); err != nil {
		log.Panic("Maxmind Database Invalid")
		return nil, err
	}
	log.Print("Maxmind Database Configured")
	return &Reader{reader}, nil
}

func (db *Reader) IpLookupRegion(ipaddress string) interface{} {
	log.Print("Get request /" + ipaddress)
	ip := net.ParseIP(ipaddress)
	var record interface{}
	err := db.Lookup(ip, &record)
	if err != nil {
		return nil
	}

	return record
}
