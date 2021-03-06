package cmd

import (
	"os"
	"strings"

	"text/template"

	cloudflare "github.com/cloudflare/cloudflare-go"
	"github.com/spf13/cobra"

	"github.com/sirupsen/logrus"
)

const accessApplicationTemplate = `
resource "cloudflare_access_application" "{{.App.ID}}" {
 	zone_id = "{{.Zone.ID}}"
 	name = "{{.App.Name}}"
 	domain = "{{.App.Domain}}"
 	session_duration = "{{.App.SessionDuration}}"
}
`

func init() {
	rootCmd.AddCommand(accessApplicationCmd)
}

var accessApplicationCmd = &cobra.Command{
	Use:   "access_application",
	Short: "Import Access Application data into Terraform",
	Run: func(cmd *cobra.Command, args []string) {
		log.Debug("Importing Access Application data")

		for _, zone := range zones {

			log.WithFields(logrus.Fields{
				"ID":   zone.ID,
				"Name": zone.Name,
			}).Debug("Processing zone")

			accessApplications, _, err := api.AccessApplications(zone.ID, cloudflare.PaginationOptions{
				Page:    1,
				PerPage: 1000,
			})

			if err != nil {
				if strings.Contains(err.Error(), "HTTP status 403") {

					log.WithFields(logrus.Fields{
						"ID": zone.ID,
					}).Info("Insufficient permissions to access zone")
					continue
				}
				log.Fatal(err)
				os.Exit(1)
			}

			for _, app := range accessApplications {

				accessApplicationParse(app, zone)
			}

		}
	},
}

func accessApplicationParse(app cloudflare.AccessApplication, zone cloudflare.Zone) {
	tmpl := template.Must(template.New("access_rule").Funcs(templateFuncMap).Parse(accessApplicationTemplate))
	tmpl.Execute(os.Stdout,
		struct {
			App  cloudflare.AccessApplication
			Zone cloudflare.Zone
		}{
			App:  app,
			Zone: zone,
		})
}
