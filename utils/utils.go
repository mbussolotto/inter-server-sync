package utils

import (
	"bufio"
	"fmt"
	"github.com/rs/zerolog/log"
	"io"
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

func GetAbsPath(path string) string {
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

func GetCurrentServerVersion() (string, string) {

	rhndefault := "/etc/rhn/rhn.conf"
	webpath := "/usr/share/rhn/config-defaults/rhn_web.conf"
	altpath := "/usr/share/rhn/config-defaults/rhn.conf"

	files := []string{rhndefault, webpath, altpath}
	property := []string{"product_name", "web.product_name"}
	product := "SUSE Manager"
	p, err := getProperty(files, property)
	if err == nil {
		product = p
	}

	propertyVersion := []string{"web.version"}
	if product != "SUSE Manager" {
		propertyVersion = []string{"web.version.uyuni"}
		product = "uyuni"
	}
	version, err := getProperty(files, propertyVersion)
	if err != nil {
		log.Fatal().Msgf("No version found for product %s", product)
	}
	return version, product
}

func getProperty(filePaths []string, names []string) (string, error) {
	for _, path := range filePaths {
		for _, search := range names {
			p, err := ScannerFunc(path, search)
			if err == nil {
				return p, nil
			}
		}
	}
	return "", fmt.Errorf("String not found!")
}

func ScannerFunc(path string, search string) (string, error) {
	var output string
	f, err := os.Open(path)
	if err != nil {
		log.Fatal().Msg("Couldn't open file")
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), search) {
			splits := strings.Split(scanner.Text(), "=")
			output = splits[1]
			if output == " SUSE Manager" {
				output = strings.Replace(output, " SUSE Manager", "SUSE Manager", 1)
			} else {
				splits = strings.Split(output, " ")
				output = splits[len(splits)-1]
			}
			return output, nil
		}
	}
	return "", fmt.Errorf("String not found!")
}

func ValidateExportFolder(outputFolderAbs string) {
	outputFolder, err := os.Open(outputFolderAbs)
	defer outputFolder.Close()
	if err != nil {
		if os.IsNotExist(err) {
			err := os.MkdirAll(outputFolderAbs, 0755)
			if err != nil {
				log.Fatal().Err(err).Msg("Error creating dir")
			}
			outputFolder, _ = os.Open(outputFolderAbs)
		} else {
			log.Fatal().Err(err).Msg("Error getting output foulder")
		}

	}
	folderInfo, err := outputFolder.Stat()
	if err != nil {
		log.Fatal().Err(err).Msg("Error getting folder info")
	}

	if !folderInfo.IsDir() {
		log.Fatal().Err(err).Msg(fmt.Sprintf("export location is not a directory: %s", outputFolderAbs))
	}
	_, errEmpty := outputFolder.Readdirnames(1) // Or f.Readdir(1)
	if errEmpty != io.EOF {
		log.Fatal().Msg(fmt.Sprintf("export location is not empty: %s", outputFolderAbs))

	}
}
