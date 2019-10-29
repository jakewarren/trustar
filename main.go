package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/briandowns/spinner"
	"github.com/jakewarren/trustar-golang"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	c *trustar.Client

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
	}{}

	s = spinner.New(spinner.CharSets[14], 100*time.Millisecond)
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
}

func main() {
	var err error
	c, err = trustar.NewClient(viper.GetString("TRUSTAR_CLIENT_ID"), viper.GetString("TRUSTAR_CLIENT_SECRET"), trustar.APIBaseLive)
	if err != nil {
		log.Fatal().Err(err).Msg("error creating client")
	}

	// TODO: add option for user to log to file
	c.SetLog(ioutil.Discard)

	_, err = c.GetAccessToken()
	if err != nil {
		log.Fatal().Err(err).Msg("error while getting access token")
	}

	if err = rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(0)
	}
}
