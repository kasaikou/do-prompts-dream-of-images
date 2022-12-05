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

	"github.com/dietsche/rfsnotify"
	"gopkg.in/fsnotify.v1"
)

func main() {

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	watcher, err := rfsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	exitWatcher := make(chan *sync.WaitGroup)
	go func() {
		parallels := parallelProcesses{}

		for {

		selectbreak:
			select {
			case wg := <-exitWatcher:
				defer wg.Done()
				time.Sleep(3 * time.Second)
				return
			case event := <-watcher.Events:
				// check executing
				if parallels.numRunning() > 0 {
					log.Printf("file changed detected, process is running: %s", event.Name)
					break selectbreak
				}

				if isTargetFile(event) {
					events := map[string]struct{}{
						event.Name: {},
					}

					timer := time.NewTimer(time.Second)

				timeoutLoop:
					for {
						select {
						case <-timer.C:
							break timeoutLoop
						case event := <-watcher.Events:
							if isTargetFile(event) {
								timer = time.NewTimer(time.Second)
								events[event.Name] = struct{}{}
							}
						}
					}

					for path := range events {
						parallels.add(func() {
							cmd := exec.Command(`jupytext`, `--set-formats`, `@/ipynb,docs//md:markdown,py:percent`, path)
							cmd.Stdout = os.Stdout
							cmd.Stderr = os.Stderr

							log.Printf("execute command: %s", cmd.String())
							if err := cmd.Run(); err != nil {
								log.Printf("process error: %s", err.Error())
							}
						})
					}

				}
			}
		}
	}()

	if err := watcher.AddRecursive("."); err != nil {
		log.Fatal("cannot add fswatcher: ", err.Error())
	}

	if err := watcher.RemoveRecursive(".git"); err != nil {
		log.Fatal("cannot add fswatcher: ", err.Error())
	}

	sig := <-signals
	log.Printf("signal called: %s", sig.String())
}

type parallelProcesses struct {
	lock       sync.RWMutex
	numProcess int
}

func (pp *parallelProcesses) add(fn func()) {
	func() {
		pp.lock.Lock()
		defer pp.lock.Unlock()
		pp.numProcess++
	}()

	go func() {
		fn()
		pp.lock.Lock()
		defer pp.lock.Unlock()
		pp.numProcess--
	}()
}

func (pp *parallelProcesses) numRunning() int {
	pp.lock.RLock()
	defer pp.lock.RUnlock()
	return pp.numProcess
}

func isTargetFile(event fsnotify.Event) bool {

	switch path.Ext(event.Name) {
	case ".ipynb", ".py":
	default:
		log.Printf("ignore file event: %s", event.Name)
		return false
	}

	return true
}
