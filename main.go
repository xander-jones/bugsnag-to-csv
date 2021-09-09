package main

import (
	"flag"
	"fmt"

	"github.com/bugsnag/bugsnag-go/v2"
	"github.com/xander-jones/bugsnag-to-csv/pkg/common"
	"github.com/xander-jones/bugsnag-to-csv/pkg/daa"
)

func main() {
	bugsnag.Configure(bugsnag.Configuration{
		APIKey:       "575d0c15fae9fa9c865ede8258dea307",
		AppVersion:   common.PackageVersion,
		ReleaseStage: "development",
		ProjectPackages: []string{
			"main",
			"github.com/xander-jones/bugsnag-to-csv",
			"github.com/xander-jones/bugsnag-to-csv/common",
			"github.com/xander-jones/bugsnag-to-csv/daa",
		},
		Synchronous: true,
	})

	token := flag.String("token", "", "[REQUIRED][String] Your Bugsnag personal auth token.")
	getProjectIds := flag.Bool("show-project-ids", false, "[Flag] Use this flag to get a list of project IDs accessible with your token.")
	projectId := flag.String("project-id", "", "[String] The Project ID you wish to download from")
	errorId := flag.String("error-id", "", "[String] An error ID to download. If provided, downloads all events within filters for this error ID")
	events := flag.Bool("events", false, "[Flag] Download events rather than error groups when this flag is enabled (will connsume a lot more data)")
	//outputFilepath := flag.String("output-file", "", "[String] Filepath to store the downloaded CSV.")
	//filters := flag.String("filters", "", "A JSON string array of filters to apply")
	verbose := flag.Bool("verbose", false, "[Flag] Set the output to be verbose for debugging purposes.")
	flag.Parse()

	common.Verbose = *verbose
	common.PrintHeader()

	if *token == "" {
		common.ExitWithError(1, "Missing token. Please supply Bugsnag personal auth token with --token flag")
	} else {
		daa.PersonalAuthToken = *token
		if *getProjectIds {
			common.Print("Getting your project IDs with provided token")
			orgs := daa.GetUsersOrganizations(false, 30)
			for _, org := range orgs {
				common.Print("Organization: " + fmt.Sprint(org["name"]) + " [" + fmt.Sprint(org["id"]) + "]")
				projects := daa.GetOrganizationsProjects(org["id"].(string), 10)
				for _, proj := range projects {
					common.Print("  > " + fmt.Sprint(proj["name"]) + " [" + fmt.Sprint(proj["id"]) + "]")
				}
			}
		} else {
			if *projectId == "" {
				common.ExitWithError(1, "Missing Project ID. Please supply a Project ID with --project-id flag")
			} else {
				if *errorId == "" {
					if *events {
						common.Print("Downloading all events for projectId within filters")
						events := daa.GetProjectEvents(*projectId)
						for _, event := range events {
							common.Print(fmt.Sprint(event))
						}
					} else {
						common.Print("Downloading all errors from projectId within filters")
						errs := daa.GetProjectErrors(*projectId)
						for _, err := range errs {
							common.Print(fmt.Sprint(err))
						}
					}
				} else {
					if *events {
						common.Print("Downloading all events for projectId & errorId within filters")
						events := daa.GetErrorEvents(*projectId, *errorId)
						for _, event := range events {
							common.Print(fmt.Sprint(event))
						}
					} else {
						common.Print("Downloading error from projectId & errorId within filters")
						errs := daa.GetError(*projectId, *errorId)
						for _, err := range errs {
							common.Print(fmt.Sprint(err))
						}
					}
				}
			}
		}
	}
}
