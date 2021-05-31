package cli

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"syscall"
	"time"

	"github.com/bennycio/bundle/cli/file"
	"github.com/bennycio/bundle/cli/logger"
	"github.com/bennycio/bundle/cli/term"
	"github.com/bennycio/bundle/internal"
	"github.com/c-bata/go-prompt"
	"github.com/jlaffaye/ftp"
	. "github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	goterm "golang.org/x/term"
)

type anFtp struct {
	Name     string
	Host     string
	Port     string
	Username string
	Password string
	Conn     *ftp.ServerConn
}

var theFtp anFtp

var buFileCache file.BundleFile

var connectCommands []prompt.Suggest = []prompt.Suggest{
	{Text: "help", Description: "See command options"},
	{Text: "install", Description: "Install/Update plugins"},
	{Text: "init", Description: "Create a new bundle file"},
	{Text: "remove", Description: "Remove a plugin from bundle file"},
	{Text: "uninstall", Description: "Delete a plugin"},
	{Text: "status", Description: "Check for updates"},
	{Text: "exit", Description: "Disconnect from FTP instance"},
	{Text: "list", Description: "List Plugins in Bundle File"},
	{Text: "add", Description: "Add plugin to bundle file"},
}

// testCmd represents the test command
var ftpCmd = &cobra.Command{
	Use:   "ftp",
	Short: "Connect to an instance of an FTP server to run bundle commands from",
	RunE: func(cmd *cobra.Command, args []string) error {
		ftps := viper.GetStringMap("FTP")

		if len(ftps) < 1 {
			if err := newFtp(); err != nil {
				return err
			}
		} else {
			keys := []string{"new"}
			term.Println("Which connection would you like to use?")
			for k := range ftps {
				fmt.Printf(" - %s\n", k)
				keys = append(keys, k)
			}
			fmt.Println(" -", Gray(12, "new (create a new connection)"))
			result := prompt.Choose(">> ", keys)

			if result == "new" {
				if err := newFtp(); err != nil {
					return err
				}
			} else {

				theFtp.Name = result

				term.Println("What would you like to do with this connection? (connect or remove)")
				opts := []string{"connect", "remove", "cancel"}
				resOpt := prompt.Choose(">> ", opts)

				switch strings.ToLower(resOpt) {
				case "connect":
					resultInMap, ok := ftps[result].(map[string]interface{})
					if !ok {
						return errors.New("invalid ftp config format")
					}

					if host, ok := resultInMap["host"].(string); ok {
						theFtp.Host = host
					} else {
						return errors.New("no host specified")
					}
					if port, ok := resultInMap["port"].(string); ok {
						theFtp.Port = port
					} else {
						return errors.New("no port specified")
					}
					if username, ok := resultInMap["username"].(string); ok {
						theFtp.Username = username
					} else {
						return errors.New("no username specified")
					}
					if pass, ok := resultInMap["password"].(string); ok {
						theFtp.Password = pass
					} else {
						return errors.New("no password specified")
					}
				case "remove":
					delete(ftps, result)
					viper.Set("FTP", ftps)
					err := viper.WriteConfig()
					if err != nil {
						return err
					}
					os.Exit(1)
				default:
					os.Exit(1)

				}
			}
		}

		connection, err := ftp.Dial(fmt.Sprintf("%s:%s", theFtp.Host, theFtp.Port), ftp.DialWithDisabledEPSV(true), ftp.DialWithDisabledMLSD(true))
		if err != nil {
			return err
		}

		err = connection.Login(theFtp.Username, theFtp.Password)
		if err != nil {
			return err
		}

		defer connection.Quit()
		defer connection.Logout()

		fmt.Printf("%s\nType 'help' for commands\nType 'exit' to exit\n", fmt.Sprintf("%s %s", Green("Connected To").Bold(), Green(theFtp.Name).Bold()))

		theFtp.Conn = connection

		bu, err := file.GetBundleFtp(connection)
		if err == nil {
			buFileCache = bu
		}

		pr := prompt.New(connectedExecutor, connectedCompleter, prompt.OptionPrefix(">> "))
		pr.Run()

		return nil
	},
}

func init() {
	rootCmd.AddCommand(ftpCmd)
}

func connectedCompleter(d prompt.Document) []prompt.Suggest {
	args := strings.Split(d.TextBeforeCursor(), " ")
	if len(args) > 1 {
		if args[0] == "remove" || args[0] == "uninstall" {
			s := []prompt.Suggest{}

			for k := range buFileCache.Plugins {
				s = append(s, prompt.Suggest{Text: k})
			}
			return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
		}
	}
	return prompt.FilterHasPrefix(connectCommands, d.GetWordBeforeCursor(), true)
}

func connectedExecutor(s string) {

	conn := theFtp.Conn

	if s == "exit" {
		conn.Logout()
		conn.Quit()
		os.Exit(1)
	}

	args := strings.Split(s, " ")
	if len(args) < 1 {
		return
	}

	switch strings.ToLower(args[0]) {
	case "help":
		for _, v := range connectCommands {
			fmt.Printf("%s: %s\n", Green(v.Text).Bold(), v.Description)
		}
	case "install":
		if len(args) > 1 {

			plsToInstall := map[string]string{}
			for _, v := range args[1:] {
				spl := strings.Split(v, "@")
				if len(spl) < 2 {
					plsToInstall[spl[0]] = "latest"
				} else {
					plsToInstall[spl[0]] = spl[1]
				}
			}
			err := downloadAndInstall(plsToInstall, conn)
			if err != nil {
				logger.ErrLog.Print(err.Error())
				return
			}
		} else {
			result, err := file.GetBundleFtp(conn)
			if err != nil {
				logger.ErrLog.Print(err.Error())
				return
			}
			buFileCache = result
			err = downloadAndInstall(result.Plugins, conn)
			if err != nil {
				logger.ErrLog.Print(err.Error())
				return
			}
		}
		fmt.Println(Green("Successfully installed plugins!").Bold())
	case "init":

		names, err := conn.NameList(".")
		if err != nil {
			logger.ErrLog.Print(err.Error())
			return
		}
		if internal.Contains(names, file.BuFileName) {
			fmt.Println("bundle file already exists")
			return
		}

		pr, pw := io.Pipe()

		go func() {
			defer pw.Close()
			_, err := pw.Write([]byte(file.BuFile))
			if err != nil {
				logger.ErrLog.Print(err.Error())
				return
			}
		}()

		err = conn.Stor(file.BuFileName, pr)
		if err != nil {
			logger.ErrLog.Print(err.Error())
			return
		}

		fmt.Println(Green("Successfully initialized bundle file!").Bold())

	case "status":
		bufile, err := file.GetBundleFtp(conn)
		if err != nil {
			logger.ErrLog.Print(err.Error())
			return
		}
		buFileCache = bufile
		printStatus(bufile.Plugins, conn)
	case "uninstall":
		if len(args) < 2 {
			fmt.Println("Please specify a plugin to remove")
			return
		}
		bu, err := file.GetBundleFtp(conn)
		if err != nil {
			logger.ErrLog.Print(err.Error())
			return
		}
		for _, v := range args[1:] {
			for pl := range bu.Plugins {
				if strings.EqualFold(v, pl) {
					err = conn.Delete(fmt.Sprintf("plugins/%s.jar", pl))
					if err != nil {
						logger.ErrLog.Print(err.Error())
						return
					}
					time.Sleep(2 * time.Second)
					delete(bu.Plugins, pl)
				}
			}
		}
		err = file.WritePluginsToBundleFtp(conn, bu.Plugins)
		if err != nil {
			logger.ErrLog.Print(err.Error())
			return
		}
		new, err := file.GetBundleFtp(conn)
		if err != nil {
			logger.ErrLog.Print(err.Error())
			return
		}
		buFileCache = new
	case "remove":
		if len(args) < 2 {
			fmt.Println("Please specify a plugin to remove")
			return
		}
		bu, err := file.GetBundleFtp(conn)
		if err != nil {
			logger.ErrLog.Print(err.Error())
			return
		}
		for _, v := range args[1:] {
			for pl := range bu.Plugins {
				if strings.EqualFold(v, pl) {
					delete(bu.Plugins, pl)
				}
			}
		}
		err = file.WritePluginsToBundleFtp(conn, bu.Plugins)
		if err != nil {
			logger.ErrLog.Print(err.Error())
			return
		}
		new, err := file.GetBundleFtp(conn)
		if err != nil {
			logger.ErrLog.Print(err.Error())
			return
		}
		buFileCache = new
	case "list":
		bu, err := file.GetBundleFtp(conn)
		if err != nil {
			logger.ErrLog.Print(err.Error())
			return
		}
		buFileCache = bu
		i := 1
		for k, v := range bu.Plugins {
			fmt.Printf("%d. %s - %s\n", Green(i), Yellow(k).Bold(), Yellow(v))
			i += 1
		}
	case "add":
		if len(args) < 2 {
			fmt.Println("Please specify a plugin to add")
			return
		}
		bu, err := file.GetBundleFtp(conn)
		if err != nil {
			logger.ErrLog.Print(err.Error())
			return
		}
		if bu.Plugins == nil {
			bu.Plugins = map[string]string{}
		}
		for _, v := range args[1:] {
			spl := strings.Split(v, "@")
			if len(spl) > 1 {
				bu.Plugins[spl[0]] = spl[1]
			} else {
				bu.Plugins[v] = "latest"
			}
		}
		err = file.WritePluginsToBundleFtp(conn, bu.Plugins)
		if err != nil {
			logger.ErrLog.Print(err.Error())
			return
		}
		new, err := file.GetBundleFtp(conn)
		if err != nil {
			logger.ErrLog.Print(err.Error())
			return
		}
		buFileCache = new
	}
}

func newFtp() error {
	ftps := viper.GetStringMap("FTP")
	term.Println("Unique name for this FTP connection: ")
	theFtp.Name = prompt.Input(">> ", nilCompleter)
	term.Println("FTP Host: ")
	theFtp.Host = prompt.Input(">> ", nilCompleter)
	term.Println("FTP Port: ")
	theFtp.Port = prompt.Input(">> ", nilCompleter)
	term.Println("FTP Username: ")
	theFtp.Username = prompt.Input(">> ", nilCompleter)
	term.Println("FTP Password: ")
	bytePassword, err := goterm.ReadPassword(syscall.Stdin)
	if err != nil {
		log.Fatal(err.Error())
	}
	theFtp.Password = string(bytePassword)
	ftps[theFtp.Name] = map[string]string{
		"Host":     theFtp.Host,
		"Port":     theFtp.Port,
		"Username": theFtp.Username,
		"Password": theFtp.Password,
	}
	viper.Set("FTP", ftps)
	err = viper.WriteConfig()
	if err != nil {
		return err
	}
	return nil
}
