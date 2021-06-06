package cmd

import (
	"bytes"
	"os"
	"os/signal"
	"text/template"
	"time"

	"github.com/Masterminds/sprig/v3"
	"github.com/pterm/pcli"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"

	"github.com/MarvinJWendt/gomark/internal"
)

var rootCmd = &cobra.Command{
	Use:     "gomark",
	Short:   "Gomark generates markdown documentation for Go packages",
	Long:    `You can use Gomark to generate markdown documentations for your Go packages.`,
	Version: "v0.0.1", // <---VERSION---> Updating this version, will also create a new GitHub release.
	RunE: func(cmd *cobra.Command, args []string) error {
		startedAt := time.Now()
		pathFlag, _ := cmd.Flags().GetString("path")
		outputFlag, _ := cmd.Flags().GetString("output")

		godoc, err := internal.GetGoDoc(pathFlag)
		if err != nil {
			return err
		}
		err = godoc.Parse()
		if err != nil {
			return err
		}

		t := template.New("godoc").Funcs(sprig.TxtFuncMap())
		t, err = t.Parse(internal.DefaultMarkdownTemplate)
		if err != nil {
			return err
		}

		var tpl bytes.Buffer
		// err = t.Execute(&tpl, internal.Package{})
		// err = t.Execute(&tpl, internal.GenerateTestPackage())
		err = t.Execute(&tpl, godoc.Package)
		if err != nil {
			return err
		}

		if outputFlag != "" {
			err := os.WriteFile(outputFlag, tpl.Bytes(), 0600)
			if err != nil {
				return err
			}
		} else {
			pterm.Printfln("%s", tpl.String())
		}

		if !pterm.RawOutput {
			pterm.Success.Printfln("Successfully generated docs for %s! %s", pterm.Magenta(godoc.Package.Name), pterm.Gray("("+time.Since(startedAt).String()+")"))
		}

		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	// Fetch user interrupt
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		pterm.Warning.Println("user interrupt")
		pcli.CheckForUpdates()
		os.Exit(0)
	}()

	// Execute cobra
	if err := rootCmd.Execute(); err != nil {
		pcli.CheckForUpdates()
		os.Exit(1)
	}

	pcli.CheckForUpdates()
}

func init() {
	// Adds global flags for PTerm settings.
	// Fill the empty strings with the shorthand variant (if you like to have one).
	rootCmd.PersistentFlags().BoolVarP(&pterm.PrintDebugMessages, "debug", "d", false, "enable debug messages")
	rootCmd.PersistentFlags().BoolVarP(&pterm.RawOutput, "raw", "", false, "print unstyled raw output (set it if output is written to a file)")
	rootCmd.PersistentFlags().BoolVarP(&pcli.DisableUpdateChecking, "disable-update-checks", "", false, "disables update checks")

	rootCmd.Flags().StringP("path", "p", ".", "path to search for go files")
	rootCmd.Flags().StringP("output", "o", "", "output path")

	// Use https://github.com/pterm/pcli to style the output of cobra.
	pcli.SetRepo("MarvinJWendt/gomark")
	pcli.SetRootCmd(rootCmd)
	pcli.Setup()

	// Change global PTerm theme
	pterm.ThemeDefault.SectionStyle = *pterm.NewStyle(pterm.FgCyan)
}
