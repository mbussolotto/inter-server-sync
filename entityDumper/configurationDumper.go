package entityDumper

import (
	"bufio"
	"database/sql"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/uyuni-project/inter-server-sync/dumper"
	"github.com/uyuni-project/inter-server-sync/schemareader"
	"os"
)

func ConfigTableNames() []string {
	return []string {
		"rhnconfigchanneltype",
		"rhnconfigfile",
		"rhnconfigfilename",
		"rhnconfigrevision",
		"rhnconfigcontent",
		"rhnconfigchannel",
		"rhnconfigfilestate",
		"rhnregtokenconfigchannels",
		"rhnserverconfigchannel",
		"rhnsnapshotconfigchannel",
		"susestaterevisionconfigchannel",
		"rhnconfiginfo",
		"rhnconfigfilefailure",
	}
}

func DumpConfigs(options ChannelDumperOptions) {
	var outputFolderAbs = options.GetOutputFolderAbsPath()
	db := schemareader.GetDBconnection(options.ServerConfig)
	defer db.Close()
	file, err := os.Create(outputFolderAbs + "/configurations.sql")
	if err != nil {
		log.Fatal().Err(err).Msg("error creating sql file")
		panic(err)
	}
	defer file.Close()
	bufferWriter := bufio.NewWriter(file)
	defer bufferWriter.Flush()

	bufferWriter.WriteString("BEGIN;\n")
	processAndInsertConfigs(db, bufferWriter)
	// processAndInsertConfigChannels(db, bufferWriter, loadChannelsToProcess(db, options), options)

	bufferWriter.WriteString("COMMIT;\n")
}

func processAndInsertConfigs(db *sql.DB, writer *bufio.Writer) {
	schemaMetadata := schemareader.ReadTablesSchema(db, ConfigTableNames())
	startingTables := []schemareader.Table{schemaMetadata["rhnconfigchannel"]}
	var whereFilterClause = func(table schemareader.Table) string {
		filterOrg := ""
		// if _, ok := table.ColumnIndexes["org_id"]; ok {

			//filterOrg = " where org_id is null"
		//}
		return filterOrg
	}

	dumper.DumpAllTablesData(db, writer, schemaMetadata, startingTables, whereFilterClause, onlyIfParentExistsTables )
}


func processAndInsertConfigChannels(db *sql.DB, writer *bufio.Writer, channels[]string, options ChannelDumperOptions) {
	log.Info().Msg(fmt.Sprintf("%d channels to process", len(channels)))
	schemaMetadata := schemareader.ReadTablesSchema(db, ConfigTableNames())
	log.Debug().Msg("channel schema metadata loaded")
	configChannels, err := os.Create(options.GetOutputFolderAbsPath() + "/exportedConfigChannels.txt")
	if err != nil {
		log.Fatal().Err(err).Msg("error creating exportedConfigChannel file")
		panic(err)
	}
	defer configChannels.Close()
	bufferWriterChannels := bufio.NewWriter(configChannels)
	defer bufferWriterChannels.Flush()

	count := 0
	for _, channelLabel := range channels {
		count++
		log.Info().Msg(fmt.Sprintf("Processing channel [%d/%d] %s", count, len(channels), channelLabel))
		processChannel(db, writer, channelLabel, schemaMetadata, options)
		writer.Flush()
		bufferWriterChannels.WriteString(fmt.Sprintf("%s\n", channelLabel))
	}

}

