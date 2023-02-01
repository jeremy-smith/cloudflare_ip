package main

import (
	"errors"
	"flag"
	"io"
	"log"
	"net/http"
	"regexp"

	"github.com/tidwall/gjson"

	"cloudflare/cloudflare"
	"cloudflare/config"
)

const (
	defaultConfigFile = "config.yml"
)

func getExternalIP(jsonIPService, jsonQuery string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, jsonIPService, nil)
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

	jsonRes := gjson.Get(string(b), jsonQuery)

	if ok, _ := regexp.MatchString(`^((25[0-5]|(2[0-4]|1\d|[1-9]|)\d)\.?\b){4}$`, jsonRes.String()); !ok {
		return "", errors.New("could not get the ip from the IP service")
	}

	return jsonRes.String(), nil
}

func main() {
	configFile := flag.String("c", "config.yml", "config file")
	flag.Parse()

	if *configFile == "" {
		*configFile = defaultConfigFile
	}

	// read Cloudflare config from yaml file
	bs, err := config.ReadConfig(*configFile)
	if err != nil {
		log.Fatalln(err)
	}
	conf, err := config.ParseConfig(bs)
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Calling %s", conf.JsonIPService)
	ip, err := getExternalIP(conf.JsonIPService, conf.JsonQuery)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("Got ip: %s", ip)

	cf := cloudflare.NewCloudflare(conf.AccessToken)

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
		log.Printf("Calling ListDNSRecords {zoneID: %s, recordName: %s, recordType: %s}",
			zoneID, conf.RecordName, conf.RecordType)
		dnsID, err = cf.ListDNSRecords(zoneID, conf.RecordName, conf.RecordType)
		if err != nil {
			log.Fatalln(err)
		}

		if dnsID == "" {
			log.Printf("Calling CreateDNSRecord {zoneID: %s, recordName: %s, recordType: %s, ip: %s}",
				zoneID, conf.RecordName, conf.RecordType, ip)
			createdDnsId, err := cf.CreateDNSRecord(zoneID, conf.RecordName, conf.RecordType, ip)
			if err != nil {
				log.Fatalln(err)
			}
			log.Printf("Created %s record with dnsId: %s", conf.RecordType, createdDnsId)
		}
	}

	if dnsID != "" {
		log.Printf("Calling UpdateDNSRecord {zoneID: %s, dnsID: %s, recordName: %s, recordType: %s, ip: %s}",
			zoneID, dnsID, conf.RecordName, conf.RecordType, ip)
		_, err = cf.UpdateDNSRecord(zoneID, dnsID, conf.RecordName, conf.RecordType, ip)
		if err != nil {
			log.Fatalln(err)
		}
	}
}
