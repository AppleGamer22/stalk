package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
)

var command string

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
		} else {
			return errors.New("at list 1 file must be specified")
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

		fullCommand := strings.Split(command, " ")
		command := fullCommand[0]
		arguments := make([]string, 0, len(fullCommand)-1)
		if len(fullCommand) > 1 {
			arguments = append(arguments, fullCommand[1:]...)
		}

		lastEventTime := time.Unix(0, 0)
		go func() {
			for {
				select {
				case event := <-watcher.Events:
					if event.Name == "" || (event.Op != fsnotify.Write && event.Op != fsnotify.Create) || time.Since(lastEventTime) <= time.Second/10 {
						continue
					}

					lastEventTime = time.Now()
					log.Println(event)
					if err := start(command, arguments...); err != nil {
						errs <- err
						return
					}
				case err := <-watcher.Errors:
					if err != nil {
						log.Println(err)
					} else {
						return
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
					log.Println("watching", path)
				}
			}
		}()

		signals := make(chan os.Signal, 1)
		signal.Notify(signals, os.Interrupt, syscall.SIGINT, syscall.SIGQUIT)
		select {
		case <-signals:
			kill(true)
			fmt.Print("\r")
			return nil
		case err := <-errs:
			kill(true)
			errs <- nil
			return err
		}
	},
}

func init() {
	watchCommand.Flags().BoolVarP(&verbose, "verbose", "v", false, "log a list of watched files")
	watchCommand.Flags().StringVarP(&command, "command", "c", "", "command to run after each file change")
	watchCommand.MarkFlagRequired("command")
	RootCommand.AddCommand(watchCommand)
}
