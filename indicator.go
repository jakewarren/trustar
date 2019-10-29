package main

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/jakewarren/trustar-golang"
	"github.com/rodaine/table"
	"github.com/spf13/cobra"
)

var indicatorCmd = &cobra.Command{
	Use:   "indicator",
	Short: "Manage indicators",
}

// https://docs.trustar.co/api/v13/indicators/search_indicators.html
var indicatorSearchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search indicators",
	Run: func(cmd *cobra.Command, args []string) {
		output := func(indicators []trustar.Indicator) {
			for _, i := range indicators {
				fmt.Println(i.GUID)
			}

			if len(indicators) == 0 {
				fmt.Println("0 indicators found.")
			}
		}

		indicators := make([]trustar.Indicator, 0)

		if len(args) == 0 {

			// run an empty search, essentially list all
			i, err := runIndicatorSearch("")
			if err != nil {
				fmt.Println(err)
				return
			}

			indicators = append(indicators, i.Items...)

		}
		for _, a := range args {
			s.Suffix = " Searching..."
			s.Start()
			i, err := runIndicatorSearch(a)
			s.Stop()
			if err != nil {
				fmt.Println(err)
				return
			}
			indicators = append(indicators, i.Items...)
		}

		output(indicators)
	},
}

func runIndicatorSearch(searchTerm string) (trustar.SearchIndicatorReponse, error) {
	query := url.Values{}
	if searchTerm != "" {
		query.Add("searchTerm", searchTerm)
	}

	if len(config.indicatorSearch.enclaveIDs) > 0 {
		query.Add("enclaveIds", strings.Join(config.indicatorSearch.enclaveIDs, ","))
	}

	query.Add("pageSize", strconv.Itoa(config.indicatorSearch.pageSize))

	return c.SearchIndicators(query)
}

// https://docs.trustar.co/api/v13/reports/find_correlated_reports.html
var indicatorFindCorrelatedReportsCmd = &cobra.Command{
	Use:   "find-reports",
	Short: "Find all correlated reports for an indicator",
	Run: func(cmd *cobra.Command, args []string) {
		for _, a := range args {
			// find correlated reports for a given indicator
			query := url.Values{}
			query.Add("indicators", a)
			correlatedReports, err := c.FindCorrelatedReports(query)
			if err != nil {
				fmt.Println(err)
				return
			}

			headerFmt := color.New(color.FgCyan, color.Underline).SprintfFunc()
			columnFmt := color.New(color.FgYellow).SprintfFunc()
			tbl := table.New("id", "title", "created", "updated")
			tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

			formatTime := func(e int64) string {
				t, _ := trustar.MsEpochToTime(e)
				return t.Format("2006-01-02 15:04:05 MST")
			}

			for _, r := range correlatedReports.Items {
				tbl.AddRow(r.ID, r.Title, formatTime(r.Created), formatTime(r.Updated))
			}

			if len(correlatedReports.Items) > 0 {
				tbl.Print()
			} else {
				fmt.Println("0 correlated reports found.")
			}
		}
	},
}
