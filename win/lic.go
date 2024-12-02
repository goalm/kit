//go:build windows

package win

import (
	"fmt"
	"github.com/goalm/kit/sys"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

const (
	ExpDate = "2024-12-31"
)

type NamedUsers struct {
	UUIDs []string `yaml:"UUIDs"`
}

func LoadYamlLicense(vaultPath string, licFile string, expDate string) bool {
	// Load the license file
	UUID := GetUUID()
	licFile = vaultPath + "/lic/" + licFile

	// Check if the license is valid
	yamlFile, err := os.ReadFile(licFile)
	if err != nil {
		panic(err)
	}

	var users NamedUsers
	err = yaml.Unmarshal(yamlFile, &users)
	if err != nil {
		panic(err)
	}
	log.Println(expDate)
	// step 1: check if the UUID is in the list
	isInUserList := func() bool {
		for _, value := range users.UUIDs {
			if value == UUID {
				return true
			}
		}
		log.Println("Please reach out to Molly / Martin to register your machine ID.")
		fmt.Println("Please reach out to Molly / Martin to register your machine ID.")
		return false
	}()

	// step 2: check if the license is still valid
	isBeforeExp := sys.IsBefore(expDate)
	if !isBeforeExp {
		log.Println("Please reach out to Molly / Martin to request an updated version.")
		fmt.Println("Please reach out to Molly / Martin to request an updated version.")
		return false
	}

	// step 3: return the result
	if isInUserList {
		return true
	}

	return false
}
