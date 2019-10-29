package main

import (
	"fmt"
	"net/url"
	"os"

	"github.com/fatih/color"
	"github.com/jakewarren/trustar-golang"
	"github.com/rodaine/table"
	"github.com/spf13/cobra"
)

var reportCmd = &cobra.Command{
	Use:   "reports",
	Short: "Interact with reports",
}

var reportSearchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search reports",
	Run: func(cmd *cobra.Command, args []string) {
		var (
			reports      trustar.ReportResponse
			err          error
			numOfReports int
		)

		headerFmt := color.New(color.FgCyan, color.Underline).SprintfFunc()
		columnFmt := color.New(color.FgYellow).SprintfFunc()
		tbl := table.New("id", "title", "created", "updated")
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

			for _, r := range reports.Reports {
				tbl.AddRow(r.ID, r.Title, formatTime(r.Created), formatTime(r.Updated))
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
	Use:   "open",
	Short: "Open the specified report(s) in your browser",
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
