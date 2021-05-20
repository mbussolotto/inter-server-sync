package cmd

import (
	"bufio"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/uyuni-project/inter-server-sync/utils"
	"os"
	"os/exec"
	"strings"
)

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import data to server",
	Run: runImport,
}

var importDir string

func init() {

	importCmd.Flags().StringVar(&importDir, "importDir", ".", "Location import data from")
	rootCmd.AddCommand(importCmd)
}

func runImport(cmd *cobra.Command, args []string) {
	absImportDir := utils.GetAbsPath(importDir)
	log.Info().Msg(fmt.Sprintf("starting import from dir %s", absImportDir))
	fversion, fproduct , fuyuni := getImportVersion(absImportDir)
	sversion, sproduct, suyuni := utils.GetCurrentServerVersion()
	if fversion != sversion || fproduct != sproduct || fuyuni != suyuni {
		log.Fatal().Msg("Wrong version detected")
	}
	validateFolder(absImportDir)
	runPackageFileSync(absImportDir)
	runImportSql(absImportDir)
	log.Info().Msg("import finished")
}

func getImportVersion(path string) (string, string, string) {
	var versionfile string
	var version, product, uyuni string
	versionfile = path + "/version.txt"
	f, err := os.Open(versionfile)
	if err != nil {
		log.Error().Msg("version.txt not found.")
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "product_name=") {
			splits := strings.Split(line, "=")
			product = splits[1]
		}
		if strings.Contains(line, "version=") && strings.HasPrefix(line, "v") {
			splits := strings.Split(line, "=")
			version = splits[1]
		}
		if strings.Contains(line, "uyuni_version=") {
			splits := strings.Split(line, "=")
			uyuni = splits[1]
		}
		}
			log.Debug().Msgf("Import Product: %s; Version: %s; Uyuni: %s", product, version, uyuni)
	return version, product , uyuni
	}


func validateFolder(absImportDir string) {
	_, err := os.Stat(fmt.Sprintf("%s/sql_statements.sql", absImportDir))
	if os.IsNotExist(err) {
		log.Fatal().Err(err).Msg("sql file doesn't exists on import directory.")
	}
}

func runPackageFileSync(absImportDir string) {
	cmd := exec.Command("rsync", "-og", "--chown=wwwrun:www", "-r",
		fmt.Sprintf("%s/packages/", absImportDir),
		"/var/spacewalk/packages/")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	log.Info().Msg("starting importing package files")
	err := cmd.Run()
	if err != nil {
		log.Fatal().Err(err).Msg("error importing package files")
	}
}

func runImportSql(absImportDir string) {
	cmd := exec.Command("spacewalk-sql", fmt.Sprintf("%s/sql_statements.sql", absImportDir))
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	log.Info().Msg("starting sql import")
	err := cmd.Run()
	if err != nil {
		log.Fatal().Err(err).Msg("error running the sql script")
	}
}