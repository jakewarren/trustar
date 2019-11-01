package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/GitbookIO/diskache"
	"github.com/fatih/color"
	"github.com/jakewarren/trustar-golang"
	"github.com/rodaine/table"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// command to list available enclaves
// https://docs.trustar.co/api/v13/enclaves/get_enclaves.html
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists enclaves",
	Run: func(cmd *cobra.Command, args []string) {
		// get Enclaves
		enclaves, err := c.GetEnclaves()
		if err != nil {
			fmt.Println(err)
		}

		// TODO: add option to output as json?

		prettyPrintEnclaves(enclaves)
	},
}

func prettyPrintEnclaves(enclaves []trustar.Enclave) {
	headerFmt := color.New(color.FgCyan, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()
	tbl := table.New("id", "name", "type", "read", "update", "create")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

	for _, e := range enclaves {
		tbl.AddRow(e.ID, e.Name, e.Type, e.Read, e.Update, e.Create)
	}
	tbl.Print()
}

func setupEnclaveCache() error {
	var err error

	// open up a disk cache
	dir := filepath.Join(os.TempDir(), "trustar-enclave-cache")

	log.Debug().Str("dir", dir).Msg("setting up enclave")
	opts := diskache.Opts{
		Directory: fmt.Sprintf("%s", dir),
	}

	enclaveCache, err = diskache.New(&opts)
	if err != nil {
		return err
	}

	// get the number of files in the directory
	files, _ := ioutil.ReadDir(dir)
	if len(files) == 0 {
		return fillEnclaveCache()
	}

	return nil
}

func fillEnclaveCache() error {
	log.Debug().Msg("filling up enclave")

	// get all available enclaves from Trustar
	enclaves, err := c.GetEnclaves()
	if err != nil {
		return err
	}

	// fill the cache associating the enclave id to the enclave name
	for _, e := range enclaves {
		err = enclaveCache.Set(e.ID, []byte(e.Name))
		if err != nil {
			return err
		}
	}

	return nil
}

func lookupEnclave(enclaveID string) (enclaveName string) {
	onlyOnce.Do(func() {
		err := setupEnclaveCache()
		if err != nil {
			log.Error().Err(err).Msg("error setting up cache")
		}
	})

	name, exists := enclaveCache.Get(enclaveID)

	if exists {
		enclaveName = string(name)
	} else {
		// if we can't resolve the id to an enclave name, just return the id
		enclaveName = enclaveID
	}

	return enclaveName
}
