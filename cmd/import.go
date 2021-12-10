package cmd

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/uyuni-project/inter-server-sync/utils"
	"os"
	"os/exec"
)

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import data to server",
	Run:   runImport,
}

var importDir string

func init() {

	importCmd.Flags().StringVar(&importDir, "importDir", ".", "Location import data from")
	rootCmd.AddCommand(importCmd)
}

func runImport(cmd *cobra.Command, args []string) {
	absImportDir := utils.GetAbsPath(importDir)
	log.Info().Msg(fmt.Sprintf("starting import from dir %s", absImportDir))
	fversion, fproduct := getImportVersionProduct(absImportDir)
	sversion, sproduct := utils.GetCurrentServerVersion()
	if fversion != sversion || fproduct != sproduct {
		log.Fatal().Msgf("Wrong version detected. Fileversion = %s ; Serverversion = %s", fversion, sversion)
	}
	validateFolder(absImportDir)
	//runPackageFileSync(absImportDir)
	runImportSql(absImportDir)
	log.Info().Msg("import finished")
}

func getImportVersionProduct(path string) (string, string) {
	var versionfile string
	versionfile = path + "/version.txt"
	version, err := utils.ScannerFunc(versionfile, "version")
	if err != nil {
		log.Error().Msg("Version not found.")
	}
	product, err := utils.ScannerFunc(versionfile, "product_name")
	if err != nil {
		log.Fatal().Msg("Product not found")
	}
	log.Debug().Msgf("Import Product: %s; Version: %s", product, version)
	return version, product
}

func validateFolder(absImportDir string) {
	_, err := os.Stat(fmt.Sprintf("%s/sql_statements.sql", absImportDir))
	if os.IsNotExist(err) {
		_, err := os.Stat(fmt.Sprintf("%s/configurations.sql", absImportDir))
		if os.IsNotExist(err) {
			log.Fatal().Err(err).Msg("No usable .sql files found in import directory")
		}
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
	func() {
		_, err := os.Stat(fmt.Sprintf("%s/sql_statements.sql", absImportDir))
		if os.IsExist(err) {
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
	}()

	func() {
		_, err := os.Stat(fmt.Sprintf("%s/configurations.sql", absImportDir))
		if err == nil {
			cmd := exec.Command("spacewalk-sql", fmt.Sprintf("%s/configurations.sql", absImportDir))
			err = cmd.Run()
			if err != nil {
				log.Fatal().Err(err).Msg("error running the sql script (configurations)")
			}
		}

	}()

}
