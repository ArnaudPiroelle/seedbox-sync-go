package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"seedbox-sync/downloader"
	"seedbox-sync/model"
	"seedbox-sync/notifier"
	"seedbox-sync/provider"
	"seedbox-sync/task"
)

func main() {

	syncCommand := flag.NewFlagSet("sync", flag.ExitOnError)
	scheduleCommand := flag.NewFlagSet("schedule", flag.ExitOnError)

	var config string
	syncCommand.StringVar(&config, "c", "", "")
	syncCommand.StringVar(&config, "config", "", "")
	scheduleCommand.StringVar(&config, "c", "", "")
	scheduleCommand.StringVar(&config, "config", "", "")

	if len(os.Args) < 2 {
		flag.PrintDefaults()
		os.Exit(1)
	}

	var command = os.Args[1]
	switch command {
	case "sync":
		syncCommand.Parse(os.Args[2:])
	case "schedule":
		scheduleCommand.Parse(os.Args[2:])
	default:
		flag.PrintDefaults()
		os.Exit(1)
	}

	c, err := loadConfiguration(config)
	if err != nil {
		fmt.Println("Configuration loading error")
		os.Exit(1)
	}

	providerConfiguration := c.Provider
	downloaderConfiguration := c.Downloader
	hooks := c.Hooks

	p, err := retrieveProvider(providerConfiguration)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	d, err := retrieveDownloader(downloaderConfiguration)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	logger := notifier.NewLogger()
	console := notifier.NewConsole()
	hookNotifier := notifier.NewHookNotifier(hooks)
	allNotifiers := notifier.NewCompose(logger, console, hookNotifier)

	t, err := retrieveTask(command, c, p, d, allNotifiers)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	d.Connect()
	t.Execute()
	d.Disconnect()
}

func retrieveTask(command string, c model.Configuration, p provider.Provider, d downloader.Downloader, n notifier.Notifier) (t task.Task, err error) {
	switch command {
	case "sync":
		t = task.NewSync(c.Folders, p, d, n)
		return
	case "schedule":
		t = task.NewSchedule(c.Scheduler, c.Folders, p, d, n)
		return

	default:
		err = errors.New("unknown task")
		return
	}
}

func retrieveDownloader(downloaderConfiguration model.DownloaderConfiguration) (d downloader.Downloader, err error) {
	switch downloaderType := downloaderConfiguration.Type; downloaderType {
	case "ftp":
		d = downloader.NewFtp(
			downloaderConfiguration.Host,
			downloaderConfiguration.Port,
			downloaderConfiguration.Username,
			downloaderConfiguration.Password,
			downloaderConfiguration.Root,
		)
		return
	default:
		err = errors.New("unknown downloader type")
		return
	}
}

func retrieveProvider(providerConfiguration model.ProviderConfiguration) (p provider.Provider, err error) {
	switch providerType := providerConfiguration.Type; providerType {
	case "transmission":
		p = provider.NewTransmission(
			providerConfiguration.Url,
			providerConfiguration.Username,
			providerConfiguration.Password,
		)
		return
	default:
		err = errors.New("unknown provider")
		return
	}
}

func loadConfiguration(file string) (configuration model.Configuration, err error) {
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return
	}

	configuration = model.Configuration{}
	err = json.Unmarshal(bytes, &configuration)
	return
}
