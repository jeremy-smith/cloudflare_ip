package main

import (
	"cloudflare/cloudflare"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
)

const (
	configFile  = "config.yml"
	IPService   = "http://ip-api.com/json/"
	logFileName = "cloudflare.log"
)

type IPResponse struct {
	Status string `json:"status"`
	IP     string `json:"query"`
}

func getExternalIP() (string, error) {
	req, err := http.NewRequest(http.MethodGet, IPService, nil)
	if err != nil {
		log.Fatalln(err)
	}
	c := http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var d IPResponse
	if err := json.Unmarshal(b, &d); err != nil {
		return "", err
	}
	if d.Status != "success" {
		return "", errors.New("error: status of request to IP service was not successful")
	}
	return d.IP, nil
}

func main() {

	// open log file and direct the output of log to it
	logFile, err := os.OpenFile(logFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("Error opening log file: %v", err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	log.Printf("Calling %s", IPService)
	ip, err := getExternalIP()
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("Got ip: %s", ip)

	// read Cloudflare config from file
	bs, err := readConfig(configFile)
	if err != nil {
		log.Fatalln(err)
	}
	conf, err := parseConfig(bs)
	if err != nil {
		log.Fatalln(err)
	}

	cf := cloudflare.NewCloudflare(conf.BearerToken)

	zoneID := conf.ZoneID
	if zoneID == "" {
		log.Printf("Calling ListZones {domain: %s}", conf.Domain)
		zoneID, err = cf.ListZones(conf.Domain)
		if err != nil {
			log.Fatalln(err)
		}
	}

	dnsID := conf.DnsID
	if dnsID == "" {
		log.Printf("Calling ListDNSRecords {zoneID: %s, recordNameToSet: %s, recordType: %s}",
			zoneID, conf.RecordNameToset, "A")
		dnsID, err = cf.ListDNSRecords(zoneID, conf.RecordNameToset, "A")
		if err != nil {
			log.Fatalln(err)
		}

		if dnsID == "" {
			log.Printf("Calling CreateDNSRecord {zoneID: %s, recordNameToSet: %s, ip: %s}",
				zoneID, conf.RecordNameToset, ip)
			_, err = cf.CreateDNSRecord(zoneID, conf.RecordNameToset, ip)
			if err != nil {
				log.Fatalln(err)
			}
		}
	}
	if dnsID != "" {
		log.Printf("Calling UpdateDNSRecord {zoneID: %s, dnsID: %s, recordNameToSet: %s, ip: %s}",
			zoneID, dnsID, conf.RecordNameToset, ip)
		_, err = cf.UpdateDNSRecord(zoneID, dnsID, conf.RecordNameToset, ip)
		if err != nil {
			log.Fatalln(err)
		}
	}
}
