package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"sync"
	"time"

	"github.com/GitbookIO/diskache"
	"github.com/briandowns/spinner"
	"github.com/jakewarren/trustar-golang"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	c *trustar.Client

	// disk cache to hold enclave information
	enclaveCache *diskache.Diskache
	onlyOnce     sync.Once

	// global config struct to hold all the flags
	config = struct {
		whitelistDelete struct {
			indicatorType string
		}
		indicatorSearch struct {
			enclaveIDs   []string
			entityTypes  []string  // TODO: add support for this parameter
			from         time.Time // TODO: add support for this parameter
			to           time.Time // TODO: add support for this parameter
			tags         []string  // TODO: add support for this parameter
			excludedTags []string  // TODO: add support for this parameter
			pageSize     int
		}
		verbose  bool
		debugLog string
	}{}

	s = spinner.New(spinner.CharSets[14], 100*time.Millisecond)

	// build information set by ldflags
	appName   = "trustar"
	version   = "(ﾉ☉ヮ⚆)ﾉ ⌒*:･ﾟ✧"
	commit    = "(ﾉ☉ヮ⚆)ﾉ ⌒*:･ﾟ✧"
	buildDate = "(ﾉ☉ヮ⚆)ﾉ ⌒*:･ﾟ✧"
)

func init() {
	log.Logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger()

	rootCmd.AddCommand(completionCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(tokenCmd)
	rootCmd.AddCommand(whitelistCmd)
	rootCmd.AddCommand(indicatorCmd)
	rootCmd.AddCommand(reportCmd)
	rootCmd.AddCommand(quotaCmd)
	rootCmd.AddCommand(versionCmd)

	rootCmd.PersistentFlags().BoolVarP(&config.verbose, "verbose", "v", false, "enable verbose output")
	rootCmd.PersistentFlags().StringVar(&config.debugLog, "debug-log", "", "filename to write a debug log containing the requests to trustar (by default discards debug output)")

	indicatorCmd.AddCommand(indicatorSearchCmd)
	indicatorCmd.AddCommand(indicatorFindCorrelatedReportsCmd)
	indicatorSearchCmd.Flags().StringSliceVar(&config.indicatorSearch.enclaveIDs, "enclaveids", nil, "enclave ids to search")
	indicatorSearchCmd.Flags().IntVar(&config.indicatorSearch.pageSize, "page-size", 100, "the number of results per page")

	whitelistCmd.AddCommand(whitelistListCmd)
	whitelistCmd.AddCommand(whitelistAddCmd)
	whitelistCmd.AddCommand(whitelistDeleteCmd)
	whitelistDeleteCmd.Flags().StringVarP(&config.whitelistDelete.indicatorType, "type", "t", "", "type of indicator")

	reportCmd.AddCommand(reportSearchCmd)
	reportCmd.AddCommand(reportOpenCmd)

	// TODO: add support for a config file
	viper.AutomaticEnv()
	// configErr := viper.ReadInConfig()
	// if configErr != nil {
	//	//log.Fatal().Err(configErr).Msg("error reading config")
	//}
}

var rootCmd = &cobra.Command{
	Use:   "trustar",
	Short: "A CLI for working with Trustar",
	Long:  `trustar is a swiss army knife for working with Trustar`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if config.verbose {
			zerolog.SetGlobalLevel(zerolog.DebugLevel)
		} else {
			zerolog.SetGlobalLevel(zerolog.ErrorLevel)
		}

		// for all commands other than autocomplete, initialize the client
		if cmd.Use != "autocomplete" {
			var err error
			c, err = trustar.NewClient(viper.GetString("TRUSTAR_API_KEY"), viper.GetString("TRUSTAR_API_SECRET"), trustar.APIBaseLive)
			if err != nil {
				log.Fatal().Err(err).Msg("error creating client")
			}

			_, err = c.GetAccessToken()
			if err != nil {
				log.Fatal().Err(err).Msg("error while getting access token")
			}

			if config.debugLog != "" {

				f, logErr := os.OpenFile(config.debugLog, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
				if logErr != nil {
					log.Fatal().Err(logErr).Msg("error opening debug log")
				}

				c.SetLog(f)
			} else {
				c.SetLog(ioutil.Discard)
			}
		}
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(0)
	}
}
