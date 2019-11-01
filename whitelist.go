package main

import (
	"net/url"

	"github.com/fatih/color"
	"github.com/rodaine/table"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var whitelistCmd = &cobra.Command{
	Use:   "whitelist",
	Short: "Manage the whitelist",
}

// https://docs.trustar.co/api/v13/indicators/get_whitelist.html
var whitelistListCmd = &cobra.Command{
	Use:   "list",
	Short: "List items in the whitelist",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		query := url.Values{}
		pageNumber := 0

		headerFmt := color.New(color.FgCyan, color.Underline).SprintfFunc()
		columnFmt := color.New(color.FgYellow).SprintfFunc()
		tbl := table.New("type", "value")
		tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

	begin:
		query.Add("pageNumber", string(pageNumber))
		whitelist, err := c.GetWhitelist(nil)
		if err != nil {
			log.Fatal().Err(err).Msg("error while getting whitelist")
		}

		for _, e := range whitelist.Items {
			tbl.AddRow(e.IndicatorType, e.Value)
		}

		if whitelist.HasNext {
			pageNumber++
			goto begin
		}

		tbl.Print()
	},
}

// https://docs.trustar.co/api/v13/indicators/add_to_whitelist.html
var whitelistAddCmd = &cobra.Command{
	Use:   "add <indicator>...",
	Short: "Add items to the whitelist",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := c.WhitelistIndicators(args)
		if err != nil {
			log.Fatal().Err(err).Msg("error while adding to whitelist")
		}
	},
}

// https://docs.trustar.co/api/v13/indicators/delete_from_whitelist.html
var whitelistDeleteCmd = &cobra.Command{
	Use:   "delete <indicator>...",
	Short: "Delete items from the whitelist",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if config.whitelistDelete.indicatorType == "" {
			log.Fatal().Msg("indicator type not specified")
		}

		for _, a := range args {
			// delete indicator from the whitelist
			query := url.Values{}
			query.Add("indicatorType", config.whitelistDelete.indicatorType)
			query.Add("value", a)
			err := c.DeleteFromWhitelist(query)
			if err != nil {
				log.Fatal().Err(err).Msg("error while adding to whitelist")
			}
		}
	},
}
