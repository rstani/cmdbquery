package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

type Device struct {
	Hostname    string `yaml:"hostname"`
	Vendor      string `yaml:"vendor"`
	DeviceRole  string `yaml:"device_role"`
	Model       string `yaml:"model"`
	Environment string `yaml:"environment"`
}

func loadDataFromYAML() ([]Device, error) {
	data, err := os.ReadFile("cmdb.yml")
	if err != nil {
		return nil, err
	}

	var devices []Device
	err = yaml.Unmarshal(data, &devices)
	if err != nil {
		return nil, err
	}

	return devices, nil
}

func filterDevices(devices []Device, query string) []Device {
	keyValues := strings.Split(query, ",")
	var filtered []Device

	for _, device := range devices {
		matches := true
		for _, kv := range keyValues {
			parts := strings.Split(strings.TrimSpace(kv), "=")
			if len(parts) != 2 {
				fmt.Println("Invalid query format!")
				return nil
			}

			key, value := parts[0], parts[1]
			switch key {
			case "environment":
				if device.Environment != value {
					matches = false
				}
			case "vendor":
				if device.Vendor != value {
					matches = false
				}
			case "device_role":
				if device.DeviceRole != value {
					matches = false
				}
			case "model":
				if device.Model != value {
					matches = false
				}
			// Add other cases as needed.
			default:
				fmt.Println("Unsupported query key:", key)
				return nil
			}
		}
		if matches {
			filtered = append(filtered, device)
		}
	}
	return filtered
}

func main() {
	queryFlag := flag.String("q", "", "Query string")
	hFlag := flag.Bool("h", false, "Display full data")
	flag.Parse()

	if *queryFlag == "" {
		fmt.Println("Please provide a query using -q.")
		return
	}

	devices, err := loadDataFromYAML()
	if err != nil {
		fmt.Println("Error loading data:", err)
		return
	}

	filteredDevices := filterDevices(devices, *queryFlag)

	fmt.Printf("Total devices matched: %d\n", len(filteredDevices))
	if *hFlag {
		for _, d := range filteredDevices {
			fmt.Printf("Hostname: %s, Environment: %s, Vendor: %s\n", d.Hostname, d.Environment, d.Vendor)
		}
	}
}
