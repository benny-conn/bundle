package cli

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/bennycio/bundle/api"
	"github.com/bennycio/bundle/cli/intfile"
	"github.com/bennycio/bundle/cli/term"
	"github.com/bennycio/bundle/cli/uploader"
	"github.com/bennycio/bundle/internal"
	"github.com/bennycio/bundle/internal/gate"
	"github.com/c-bata/go-prompt"
	"github.com/c-bata/go-prompt/completer"
	. "github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
)

var rdmeFileCompleter = completer.FilePathCompleter{
	IgnoreCase: true,
	Filter: func(fi os.FileInfo) bool {
		if fi.IsDir() {
			return true
		}
		if strings.HasSuffix(fi.Name(), ".md") {
			return true
		}
		return false
	},
}

// uploadCmd represents the upload command
var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload your plugin as specified in bundle-make.yml to the official Bundle Repository",
	Long: `Will upload the jar specified under JarPath into the official Bundle Repository, allowing public access
	to your plugin. Version must be unique per upload and name must be unique globally for the initial upload`,
	RunE: func(cmd *cobra.Command, args []string) error {

		if !internal.IsValidPath(args[0]) {
			log.Fatal(errors.New("invalid path").Error())
		}

		path := args[0]

		user := credentialsPrompt()

		plugin := &api.Plugin{}

		isReadme := strings.HasSuffix(path, "README.md")

		if isReadme {

			plugin = pluginInfoPrompt()

		} else {
			result, err := intfile.ParsePluginYml(path)

			if err != nil {
				return err
			}

			plugin.Name = result.Name
			plugin.Version = result.Version
			plugin.Description = result.Description
			plugin.Category = api.Category(result.Category)
			plugin.Conflicts = result.Conflicts
		}

		gs := gate.NewGateService("localhost", "8020")

		dbPl, err := gs.GetPlugin(plugin)
		isUpdating := err == nil

		term.Print(fmt.Sprintf("Uploading to Bundle Repository From: %s\n", path))

		fi, err := os.Open(path)
		if err != nil {
			return err
		}
		defer fi.Close()

		upl := uploader.New(user, fi, plugin).WithReadme(isReadme)

		err = upl.Upload()
		if err != nil {
			return err
		}
		if !isReadme {
			term.Println(Green("Successfully uploaded plugin! :)").Bold())

			term.Println("Would you like to upload a README as well? [Y/n]")

			rdmeToo := prompt.Input(">> ", yesOrNoCompleter)

			if strings.EqualFold(rdmeToo, "y") || strings.EqualFold(rdmeToo, "yes") {
				term.Println("Please specify a path to your readme file or press enter to scan for readme in close directories.")
				p := prompt.Input(">> ", rdmeFileCompleter.Complete, prompt.OptionCompletionWordSeparator(completer.FilePathCompletionSeparator))

				var rdme *os.File
				if p == "" {
					wlk := uploader.NewFileWalker("README.md", plugin.Name)
					f, err := wlk.Walk()
					if err != nil {
						return err
					}
					rdme = f
				} else {
					f, err := os.Open(path)
					if err != nil {
						return err
					}
					rdme = f
				}
				rdmeUpl := uploader.New(user, rdme, plugin).WithReadme(true)
				err = rdmeUpl.Upload()
				if err != nil {
					return err
				}
				term.Println(Green("Successfully uploaded README! :)").Bold())
			}
		} else {
			term.Println(Green("Successfully uploaded README! :)").Bold())
		}

		if isUpdating {
			if err = makeChangelog(dbPl.Id, plugin.Version); err != nil {
				return err
			} else {
				term.Println(Green("Successfully uploaded Changelog! :)").Bold())
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(uploadCmd)
}

func pluginInfoPrompt() *api.Plugin {
	fmt.Println("Enter plugin name: ")
	pluginName := prompt.Input(">> ", nilCompleter)

	plugin := &api.Plugin{
		Name:    pluginName,
		Version: "README",
	}

	return plugin
}

func yesOrNoCompleter(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: "Y"},
		{Text: "Yes"},
		{Text: "n"},
		{Text: "no"},
	}
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

func makeChangelog(pluginId, version string) error {

	gs := gate.NewGateService("localhost", "8020")
	addedList := []string{}
	updatedList := []string{}
	removedList := []string{}

	term.Println("What did you add in this version?")
	term.Println(Gray(12, "Press enter on an empty line to continue"))
	for {
		added := prompt.Input(">> ", nilCompleter)
		if strings.Trim(strings.TrimSpace(added), "\n") == "" {
			break
		}
		addedList = append(addedList, added)
	}
	term.Println("What did you remove in this version?")
	term.Println(Gray(12, "Press enter on an empty line to continue"))
	for {
		removed := prompt.Input(">> ", nilCompleter)
		if strings.Trim(strings.TrimSpace(removed), "\n") == "" {
			break
		}
		removedList = append(removedList, removed)
	}
	term.Println("What did you update in this version?")
	term.Println(Gray(12, "Press enter on an empty line to continue"))
	for {
		updated := prompt.Input(">> ", nilCompleter)
		if strings.Trim(strings.TrimSpace(updated), "\n") == "" {
			break
		}
		updatedList = append(updatedList, updated)
	}

	changelog := &api.Changelog{
		PluginId: pluginId,
		Version:  version,
		Added:    addedList,
		Removed:  removedList,
		Updated:  updatedList,
	}

	term.Println("Is this correct? [Y/n]")
	fmt.Println(Green("Added: ").Bold())
	for _, v := range addedList {
		fmt.Printf("\t - %s\n", Green(v))
	}
	fmt.Println(Red("Removed: ").Bold())
	for _, v := range removedList {
		fmt.Printf("\t - %s\n", Red(v))
	}
	fmt.Println(Blue("Updated: ").Bold())
	for _, v := range updatedList {
		fmt.Printf("\t - %s\n", Blue(v))
	}

	correct := prompt.Input(">> ", yesOrNoCompleter)

	if strings.EqualFold(correct, "y") || strings.EqualFold(correct, "yes") {
		err := gs.InsertChangelog(changelog)
		if err != nil {
			return err
		}
	} else {
		err := makeChangelog(pluginId, version)
		if err != nil {
			return err
		}
	}
	return nil
}
