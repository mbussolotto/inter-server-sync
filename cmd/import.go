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
	fileversion := getVersion(absImportDir)
	serverversion := utils.GetCurrentServerVersion()
	if fileversion == serverversion {
		log.Debug().Msg("Same version")
	} else {
		log.Debug().Msg("Wrong version")
	}
	validateFolder(absImportDir)
	runPackageFileSync(absImportDir)
	runImportSql(absImportDir)
	log.Info().Msg("import finished")
}

func getVersion(path string) string{
	var versionfile string
	var v string
	versionfile = path + "/version.txt"
	f, err := os.Open(versionfile)
	if err != nil {
		log.Error().Msg("version.txt not found.")
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
			v := scanner.Text()
			splits := strings.Split(v, "\n")
			v = splits[0]
			continue
	}
	return v
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