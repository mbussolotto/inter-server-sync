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

	//var uversion string
	//vpath := "/usr/share/rhn/config-defaults/rhn_web.conf"
	//ppath := "/usr/share/rhn/config-defaults/rhn.conf"
	vpath := "rhn_web.conf"
	ppath := "norm.conf"

	f, err := os.Open(vpath)
	if err != nil {
		log.Fatal().Msg("Couldn't find rhn_web.conf")
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "web.version = ") {
			version = scanner.Text()
			splits := strings.Split(version, "= ")
			version = splits[1]
		}
		if strings.Contains(scanner.Text(), "web.version.uyuni = ") {
			uyuni = scanner.Text()
			splits := strings.Split(uyuni, "= ")
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
			product = scanner.Text()
			splits := strings.Split(product, "= ")
			product = splits[1]
		}
	}
	return version, product, uyuni
}



