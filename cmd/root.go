package validate

import (
	"fmt"
	"github.com/SennaSemakula/nrobservgen/pkg/service"
	"github.com/spf13/cobra"
	"os"
)

var (
	cfgFile, svc, dest string
	tfFlag             bool
	red                = "\033[31m"
	green              = "\033[32m"
	reset              = "\033[0m"
	rootCmd            = &cobra.Command{
		Short: "Generate new relic monitoring template",
		Long:  "Observgen is a CLI tool used to automate the configuration management of your new relic observability stack",
	}
	templateCmd = &cobra.Command{
		Use:   "template",
		Short: "Generate new relic monitoring template",
		Long:  "Generates <service>-alerts.yaml file and .tf files for your newrelic alerts and dashboards.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var out string
			var err error

			if len(svc) <= 0 {
				return fmt.Errorf("missing --service flag")
			}

			if !tfFlag {
				out, err = service.NewYamlTemplate(svc, dest)
				if err != nil {
					return err
				}
				fmt.Printf(green+"template written to: %s\n"+reset, out)
			}

			newSvc := service.NewService(svc)
			if len(out) == 0 {
				out = svc + "-alerts.yaml"
			}
			newSvc.Load(out, newSvc)
			out, err = service.NewTFTemplate(*newSvc, dest)
			if err != nil {
				return err
			}
			fmt.Printf(green+"terraform template written to: %s\n"+reset, out)

			return nil
		},
	}
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print the version number of observgen",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("observgen v0.1.0")
		},
	}
	validateCmd = &cobra.Command{
		Use:   "validate",
		Short: "Check whether your <service>-alerts.yaml file is valid",
		Run: func(cmd *cobra.Command, args []string) {
			if len(cfgFile) <= 0 {
				fmt.Println(red + "Error: " + reset + "Did not specify config file")
				return
			}
			if err := service.Validate(cfgFile); err != nil {
				fmt.Println(red+"Error:"+reset, err)
				os.Exit(1)
			}
			fmt.Printf(green+"Success: %s is valid\n"+reset, cfgFile)
		},
	}
)

func init() {
	rootCmd.AddCommand(versionCmd, validateCmd, templateCmd)
	templateCmd.PersistentFlags().StringVar(&svc, "service", "", "microservice name e.g. myservice")
	templateCmd.PersistentFlags().StringVar(&dest, "dest", "", "folder destination for template. If empty defaults to current working directory.")
	validateCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file e.g. <service>-alerts.yaml")
	templateCmd.PersistentFlags().BoolVar(&tfFlag, "terraform", false, "generate tf files files only")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
