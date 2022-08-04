package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
)

var waitCommand = &cobra.Command{
	Use:   "wait",
	Short: "wait a file for change",
	Long:  "wait a file for change",
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
		signals := make(chan os.Signal, 1)

		go func() {
			for {
				select {
				case event := <-watcher.Events:
					if event.Name == "" || (event.Op != fsnotify.Write && event.Op != fsnotify.Create) {
						continue
					}
					log.Println(event)
					signals <- syscall.Signal(0)
					return
				case err := <-watcher.Errors:
					if err != nil {
						errs <- err
						return
					}
				}
			}
		}()

		go func() {
			for _, path := range args {
				if err := watcher.Add(path); err != nil {
					errs <- err
					return
				}
				if verbose {
					log.Println("watching", path)
				}
			}
		}()

		signal.Notify(signals, os.Interrupt, syscall.SIGINT, syscall.SIGQUIT)
		select {
		case <-signals:
			fmt.Print("\r")
			return nil
		case err := <-errs:
			return err
		}
	},
}

func init() {
	waitCommand.Flags().BoolVarP(&verbose, "verbose", "v", false, "log a list of watched files")
	RootCommand.AddCommand(waitCommand)
}
