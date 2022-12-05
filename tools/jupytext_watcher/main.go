package main

import (
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path"
	"sync"
	"syscall"
	"time"

	"github.com/fsnotify/fsnotify"
)

func main() {

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	exitWatcher := make(chan *sync.WaitGroup)
	go func() {
		var prevCmd *exec.Cmd = nil

		for {

		selectbreak:
			select {
			case wg := <-exitWatcher:
				defer wg.Done()
				time.Sleep(3 * time.Second)
				return
			case event := <-watcher.Events:
				// check executing
				if prevCmd != nil {
					if prevCmd.ProcessState == nil {
						log.Printf("file changed detected, process is running: %s", event.Name)
						break selectbreak
					}
					if !prevCmd.ProcessState.Exited() {
						log.Printf("file changed detected, but process is running: %s", event.Name)
						break selectbreak
					}
				}

				switch path.Ext(event.Name) {
				case ".ipynb", ".py", ".md":
				default:
					log.Printf("ignore file event: %s", event.Name)
					break selectbreak
				}

				run := false
				if event.Has(fsnotify.Write) {
					log.Printf("file change detected: %s", event.Name)
					run = true
				} else if event.Has(fsnotify.Create) {
					log.Printf("file create detected: %s", event.Name)
					run = true
				} else if event.Has(fsnotify.Remove) {
					log.Printf("file delete detected: %s", event.Name)
					run = true
				} else if event.Has(fsnotify.Rename) {
					log.Printf("file rename detected: %s", event.Name)
					run = true
				}

				if run {
					cmd := exec.Command(`jupytext`, `--set-formats`, `@/ipynb,docs//md:markdown,py:percent`, event.Name)
					cmd.Stdout = os.Stdout
					cmd.Stderr = os.Stderr
					log.Printf("execute command: %s", cmd.String())
					if err := cmd.Start(); err != nil {
						prevCmd = nil
					} else {
						prevCmd = cmd
					}
				}
			}
		}
	}()

	if err := watcher.Add("."); err != nil {
		log.Fatal("cannot add fswatcher: ", err.Error())
	}

	sig := <-signals
	log.Printf("signal called: %s", sig.String())
	return
}
