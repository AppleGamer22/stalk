package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
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

		arguments := strings.Split(command, " ")
		var process *exec.Cmd
		lastEventTime := time.Unix(0, 0)
		go func() {
			for {
				select {
				case event := <-watcher.Events:
					if event.Name == "" {
						continue
					} else if time.Since(lastEventTime) >= time.Second/10 {
						lastEventTime = time.Now()
						log.Println(event)
					} else {
						continue
					}

					if process != nil {
						process.Process.Kill()
					}

					if len(arguments) == 1 {
						process = exec.Command(command)
					} else if len(arguments) > 1 {
						process = exec.Command(arguments[0], arguments[1:]...)
					}

					process.Stdout = os.Stdout
					if err := process.Run(); err != nil {
						errs <- err
					}
				case err := <-watcher.Errors:
					if err != nil {
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
					log.Println("watching", path)
				}
			}
		}()

		signals := make(chan os.Signal, 1)
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
	watchCommand.Flags().BoolVarP(&verbose, "verbose", "v", false, "log a list of watched files")
	watchCommand.Flags().StringVarP(&command, "command", "c", "", "command to run after each file change")
	_ = waitCommand.MarkFlagRequired("command")
	RootCommand.AddCommand(watchCommand)
}
