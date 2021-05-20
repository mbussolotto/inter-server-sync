package utils

import (
	"bufio"
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

func GetCurrentServerVersion() (string, string, string) {
	var version,product,uyuni string
	vpath := "/usr/share/rhn/config-defaults/rhn_web.conf"
	ppath := "/usr/share/rhn/config-defaults/rhn.conf"


	f, err := os.Open(vpath)
	if err != nil {
		log.Fatal().Msg("Couldn't find rhn_web.conf")
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "web.version = ") {
			splits := strings.Split(line, "= ")
			version = splits[1]
		}
		if strings.Contains(line, "web.version.uyuni = ") {
			splits := strings.Split(line, "= ")
			uyuni = splits[1]
		}
	}
	f, err = os.Open(ppath)
	if err != nil {
		log.Fatal().Msg("Couldn't find rhn.conf")
	}
	defer f.Close()
	scanner = bufio.NewScanner(f)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "product_name = ") {
			splits := strings.Split(scanner.Text(), "= ")
			product = splits[1]
		}
	}
	return version, product, uyuni
}



