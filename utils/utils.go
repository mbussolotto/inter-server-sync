package utils

import (
	"bufio"
	"fmt"
	"github.com/rs/zerolog/log"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

//ReverseArray reverses the array
func ReverseArray(slice interface{}) {
	size := reflect.ValueOf(slice).Len()
	swap := reflect.Swapper(slice)
	for i, j := 0, size-1; i < j; i, j = i+1, j-1 {
		swap(i, j)
	}
}

// Contains is a helper method to check if a string element exist in the string slice
func Contains(slice []string, elementToFind string) bool {
	for _, element := range slice {
		if elementToFind == element {
			return true
		}
	}
	return false
}

func GetAbsPath(path string) string{
	result := path
	if filepath.IsAbs(path) {
		result, _ = filepath.Abs(path)
	} else {
		homedir, err := os.UserHomeDir()
		if err != nil {
			log.Fatal().Msg("Couldn't determine the home directory")
		}
		if strings.HasPrefix(path, "~") {
			result = strings.Replace(path, "~", homedir, -1)
		}
	}
	return result
}

func GetCurrentServerVersion() string {
	var version string
	// var uversion string
	path := "/usr/share/rhn/config-defaults/rhn_web.conf"

	f, err := os.Open(path)
	if err != nil {
		log.Fatal().Msg("Couldn't find rhn_web.conf")
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "web.version = ") {
			version = scanner.Text()
			splits := strings.Split(version, "= ")
			version = splits[len(splits)-1]
			fmt.Printf("Found string: %s", version)
		}
		/*
		if strings.Contains(scanner.Text(), "web.version.uyuni = ") {
			uversion = scanner.Text()
			splits := strings.Split(uversion, "= ")
			uversion = splits[len(splits)-1]
			fmt.Printf("Found string: %s", uversion)
		}
		*/
	}
	return version
}



