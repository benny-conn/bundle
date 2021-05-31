package cli

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/bennycio/bundle/api"
	"github.com/bennycio/bundle/cli/file"
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

		user, err := getCurrentUser()
		if err != nil {
			return err
		}

		plugin := &api.Plugin{}

		plfile, err := os.Open(path)
		if err != nil {
			return err
		}
		defer plfile.Close()

		info, err := plfile.Stat()
		if err != nil {
			return err
		}

		result, err := file.ParsePluginYml(plfile, info.Size())

		if err != nil {
			return err
		}

		plugin.Name = result.Name
		plugin.Version = result.Version
		plugin.Description = result.Description
		plugin.Category = api.Category(result.Category)
		plugin.Metadata = &api.PluginMetadata{
			Conflicts: result.Conflicts,
		}

		gs := gate.NewGateService("localhost", "8020")

		dbPl, err := gs.GetPlugin(plugin)
		isUpdating := err == nil

		fi, err := os.Open(path)
		if err != nil {
			return err
		}
		defer fi.Close()

		upl := &uploader.Uploader{
			PluginFile: fi,
			Plugin:     plugin,
			User:       user,
		}

		term.Println(Green("Queued Plugin for Upload! :)! :)").Bold())

		term.Println("Would you like to upload a README as well? [Y/n]")

		rdmeToo := prompt.Input(">> ", yesOrNoCompleter)

		if strings.EqualFold(rdmeToo, "y") || strings.EqualFold(rdmeToo, "yes") || rdmeToo == "" {
			term.Println("Please specify a path to your readme file or press enter to scan for readme in close directories.")
			p := prompt.Input(">> ", rdmeFileCompleter.Complete, prompt.OptionCompletionWordSeparator(completer.FilePathCompletionSeparator))

			if p == "" {
				wlk := uploader.NewFileWalker("README.md", plugin.Name)
				f, err := wlk.Walk()
				if err != nil {
					return err
				}
				defer f.Close()
				buf := &bytes.Buffer{}

				_, err = io.Copy(buf, f)
				if err != nil {
					return err
				}

				rdme := &api.Readme{
					Plugin: plugin,
					Text:   buf.String(),
				}
				upl.Readme = rdme
			} else {
				f, err := os.Open(p)
				if err != nil {
					return err
				}
				defer f.Close()
				buf := &bytes.Buffer{}

				_, err = io.Copy(buf, f)
				if err != nil {
					return err
				}

				rdme := &api.Readme{
					Plugin: plugin,
					Text:   buf.String(),
				}
				upl.Readme = rdme
			}
			term.Println(Green("Queued Readme for Upload! :)").Bold())
		}

		if isUpdating {
			if ch, err := makeChangelog(dbPl.Id, plugin.Version); err != nil {
				return err
			} else {
				upl.Changelog = ch
				term.Println(Green("Queued Changelog for Upload! :)").Bold())
			}
		}

		if err = upl.Upload(); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(uploadCmd)
}

func makeChangelog(pluginId, version string) (*api.Changelog, error) {
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
		fmt.Printf("  - %s\n", Green(v))
	}
	fmt.Println(Red("Removed: ").Bold())
	for _, v := range removedList {
		fmt.Printf("  - %s\n", Red(v))
	}
	fmt.Println(Blue("Updated: ").Bold())
	for _, v := range updatedList {
		fmt.Printf("  - %s\n", Blue(v))
	}

	correct := prompt.Input(">> ", yesOrNoCompleter)

	if strings.EqualFold(correct, "y") || strings.EqualFold(correct, "yes") || correct == "" {
		return changelog, nil
	} else {
		c, err := makeChangelog(pluginId, version)
		if err != nil {
			return nil, err
		}
		changelog = c
	}
	return changelog, nil
}
