package cmd

import (
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

		go func() {
			for {
				select {
				case event := <-watcher.Events:
					if event.Name == "" {
						continue
					}
					log.Println(event)
					return
				case err := <-watcher.Errors:
					if err != nil {
						errs <- err
						return
					}
				}
			}
		}()

		for _, path := range args {
			if watcher.Add(path); err != nil {
				return err
			}
			log.Println("watching", path)
		}

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
