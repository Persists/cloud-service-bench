package config

import (
	"flag"
	"fmt"
	"time"
)

type Flags struct {
	StartAt      time.Time
	InstanceName string
	Zone         string
}

func GetFlags() (Flags, error) {
	startAtStr := flag.String("start-at", "", "Time to start the load generator in RFC3339 format")
	instanceName := flag.String("instance-name", "", "The name of the instance")
	zone := flag.String("zone", "europe-west3-c", "The zone of the instance")
	flag.Parse()

	if *startAtStr == "" {
		return Flags{}, fmt.Errorf("start-at flag is not set")
	}

	if *instanceName == "" {
		return Flags{}, fmt.Errorf("instance-name flag is not set")
	}

	startAt, err := time.Parse(time.RFC3339, *startAtStr)
	if err != nil {
		return Flags{}, fmt.Errorf("invalid startAt time format, please use RFC3339 format")
	}

	return Flags{
		StartAt:      startAt,
		InstanceName: *instanceName,
		Zone:         *zone,
	}, nil
}
