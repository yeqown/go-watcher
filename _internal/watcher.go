/*
 * file watcher
 */
package _internal

import (
	"time"

	"github.com/howeyc/fsnotify"
	"github.com/silenceper/log"
)

var (
	shouldRestart bool
	evtTime       map[string]int64
	secDuration   = 1 * time.Second
)

// init to do some initialize work
// evtTime
func init() {
	evtTime = map[string]int64{}
}

func NewWatcher() (*fsnotify.Watcher, error) {
	return fsnotify.NewWatcher()
}

func StartWatch(w *fsnotify.Watcher, paths []string, exit chan bool) {
	go func() {
		for {
			select {
			case evt := <-w.Event:
				if checkFileRegexpExcluded(evt.Name, UnWatchRegExps) {
					log.Infof("skipped file [%s] cause: 'not in include regexp filetype'\n", evt.Name)
					continue
				}

				if !checkFileIncluded(evt.Name, WatchFiletypes) {
					log.Infof("skipped file [%s] cause: 'not in include filetype'\n", evt.Name)
					continue
				}

				mt := getFileModTime(evt.Name)
				if t, ok := evtTime[evt.Name]; ok && t == mt {
					// reacv one event but FILE_MOD_TIME no change
					continue
				} else if ok && UnixTimeDuration(mt, t) <= secDuration {
					// do not tigger `hotReload` while same file evt in 1s
					// log.Info("skipped trigger")
					continue
				}

				log.Infof("[%10s] changed", evt.Name)
				evtTime[evt.Name] = mt

				// hotReload
				go hotReload()

			case err := <-w.Error:
				log.Warnf("%s", err.Error())
				exit <- true
			}

			time.Sleep(secDuration)
		}
	}()

	// append paths
	for _, path := range paths {
		if err := w.Watch(path); err != nil {
			log.Errorf("Faild to watch dir [%s]\n", path)
			continue
		}
		log.Infof("Dir [%s] watched\n", path)
	}
}
