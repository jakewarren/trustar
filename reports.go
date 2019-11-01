package main

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/jakewarren/trustar-golang"
	"github.com/rodaine/table"
	"github.com/spf13/cobra"
)

var reportCmd = &cobra.Command{
	Use:   "reports",
	Short: "Interact with reports",
}

// https://docs.trustar.co/api/v13/reports/search_reports.html
var reportSearchCmd = &cobra.Command{
	Use:   "search <search term>",
	Short: "Search reports",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var (
			reports      trustar.ReportResponse
			err          error
			numOfReports int
		)

		headerFmt := color.New(color.FgCyan, color.Underline).SprintfFunc()
		columnFmt := color.New(color.FgYellow).SprintfFunc()
		tbl := table.New("id", "title", "created", "updated", "enclave")
		tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

		formatTime := func(e int64) string {
			t, _ := trustar.MsEpochToTime(e)
			return t.Format("2006-01-02 15:04:05 MST")
		}

		// search the reports
		for _, a := range args {
			s.Suffix = " Searching..."
			s.Start()
			query := url.Values{}
			query.Add("searchTerm", a)
			reports, err = c.SearchReports(query)
			s.Stop()
			if err != nil {
				fmt.Println(err)
				return
			}

			if len(reports.Reports) == 0 {
				fmt.Println("0 reports found.")
				continue
			}

			// the reports you get from the API are not de-duplicated so we have to do it ourselves :(
			reportsSeen := make(map[string]struct{})
			dedupedReports := make([]trustar.ReportDetails, 0)
			for _, r := range reports.Reports {
				// if the report has already been processed, update the associated enclaves and de-duplicate it
				_, exists := reportsSeen[r.ID]
				if exists {
					// append the enclave
					for j := range dedupedReports {
						if dedupedReports[j].ID == r.ID {
							dedupedReports[j].EnclaveIds = append(dedupedReports[j].EnclaveIds, r.EnclaveIds...)
						}
					}
				} else {
					reportsSeen[r.ID] = struct{}{}
					dedupedReports = append(dedupedReports, trustar.ReportDetails{
						Created:    r.Created,
						EnclaveIds: r.EnclaveIds,
						ID:         r.ID,
						Title:      r.Title,
						Updated:    r.Updated,
					})
				}
			}

			for _, r := range dedupedReports {
				associatedEnclaves := make([]string, 0)

				for _, e := range r.EnclaveIds {
					associatedEnclaves = append(associatedEnclaves, lookupEnclave(e))
				}

				tbl.AddRow(r.ID, r.Title, formatTime(r.Created), formatTime(r.Updated), strings.Join(associatedEnclaves, ", "))
				numOfReports++
			}

		}

		if numOfReports > 0 {
			tbl.Print()
			fmt.Printf("\n%d report(s) found.\n", numOfReports)
		}
	},
}

var reportOpenCmd = &cobra.Command{
	Use:   "open <report id>...",
	Short: "Open the specified report(s) in your browser",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		for _, a := range args {
			err := openBrowser(fmt.Sprintf("https://station.trustar.co/constellation/reports/%s", a))
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				return
			}
		}
	},
}
