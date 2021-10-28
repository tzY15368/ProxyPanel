package main

import (
	"context"
	"fmt"
	"log"

	"github.com/cloudflare/cloudflare-go"
)

const key = ""
const email = ""

func main() {
	// Construct a new API object
	api, err := cloudflare.NewWithAPIToken(key)
	if err != nil {
		log.Fatal(err)
	}

	// Most API calls require a Context
	ctx := context.Background()

	// Fetch the zone ID
	id, err := api.ZoneIDByName("fmagic.icu") // Assuming example.com exists in your Cloudflare account already
	if err != nil {
		log.Fatal(err)
	}

	// Fetch zone details
	// zone, err := api.ZoneDetails(ctx, id)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// // Print zone details
	// fmt.Println(zone)
	r, err := api.DNSRecords(ctx, id, cloudflare.DNSRecord{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v", r)

	res, err := api.CreateDNSRecord(ctx, id, cloudflare.DNSRecord{Type: "A", Name: "test.fmagic.icu", Content: "118.190.58.76", TTL: 61})

	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v", res)
}
