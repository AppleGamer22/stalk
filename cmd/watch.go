package cmd

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
)

var watchCommand = &cobra.Command{
	Use:   "watch",
	Short: "watch a file for change",
	Long:  "watch a file for change",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 {
			for _, path := range args {
				_, err := os.Stat(path)
				if err != nil {
					return err
				}
			}
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			return err
		}
		defer watcher.Close()

		errs := make(chan error, 1)

		go func() {
			for {
				select {
				case event := <-watcher.Events:
					if event.Name == "" {
						continue
					}
					log.SetOutput(os.Stdout)
					log.Println(event)
					return
				case err := <-watcher.Errors:
					if err != nil {
						log.SetOutput(os.Stderr)
						log.Println(err)
					}
				}
			}
		}()

		go func() {
			for _, path := range args {
				if err := watcher.Add(path); err != nil {
					errs <- err
				}
				if verbose {
					log.SetOutput(os.Stdout)
					log.Println("watching", path)
				}
			}
		}()

		signals := make(chan os.Signal, 1)
		signal.Notify(signals, os.Interrupt, syscall.SIGINT, syscall.SIGQUIT)
		select {
		case <-signals:
			return nil
		case err := <-errs:
			return err
		}
	},
}

func init() {
	watchCommand.Flags().BoolVarP(&verbose, "verbose", "v", false, "log a list of watched files")
	RootCommand.AddCommand(watchCommand)
}
