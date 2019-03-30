package internal

/*
 * file watcher
 */

import (
	"time"

	"github.com/yeqown/go-watcher/internal/command"
	"github.com/yeqown/go-watcher/utils"

	"github.com/howeyc/fsnotify"
	"github.com/silenceper/log"
)

// Watcher ... to watch and run command
type Watcher struct {
	*fsnotify.Watcher
	exit              chan<- bool
	d                 time.Duration
	cmd               *command.Command
	watchingPaths     []string
	watchingFiletype  []string
	unWatchingRegular []string

	evtTime map[string]int64
}

// NewWatcher ...
func NewWatcher(paths []string, exit chan<- bool, watchingFiletype, unWatchingRegular []string) (*Watcher, error) {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	watcher := &Watcher{
		Watcher:           w,
		watchingPaths:     paths,
		exit:              exit,
		d:                 1 * time.Second,
		watchingFiletype:  watchingFiletype,
		unWatchingRegular: unWatchingRegular,
		evtTime:           make(map[string]int64),
	}

	return watcher, nil
}

// AppendWatchPaths ...
// TODO: ignore duliacted path
func (w *Watcher) AppendWatchPaths(paths ...string) {
	w.watchingPaths = append(w.watchingPaths, paths...)
}

// SetCommand ...
func (w *Watcher) SetCommand(cmdName string, cmdArgs, envs []string) {
	w.cmd = command.New(cmdName, cmdArgs, envs)
	w.cmd.Start()
}

// Watching ...
func (w *Watcher) Watching() {
	go func() {
		for {
			select {
			case evt := <-w.Event:
				if utils.CheckFileRegexpExcluded(evt.Name, w.unWatchingRegular) {
					log.Infof("skipped file [%s] cause: 'not in include regexp filetype'\n", evt.Name)
					continue
				}

				if !utils.CheckFileIncluded(evt.Name, w.watchingFiletype) {
					log.Infof("skipped file [%s] cause: 'not in include filetype'\n", evt.Name)
					continue
				}

				mt := utils.GetFileModTime(evt.Name)
				if t, ok := w.evtTime[evt.Name]; ok && t == mt {
					// reacv one event but FILE_MOD_TIME no change
					continue
				} else if ok && utils.UnixTimeDuration(mt, t) <= w.d {
					// do not tigger `hotReload` while same file evt in 1s
					// log.Info("skipped trigger")
					continue
				}

				log.Infof("[%10s] changed", evt.Name)
				w.evtTime[evt.Name] = mt

				// hotReload
				go w.cmd.HotReload()
			case err := <-w.Error:
				log.Warnf("%s", err.Error())
				w.exit <- true
			default:
				time.Sleep(w.d)
			}
		}
	}()

	// append paths
	for _, path := range w.watchingPaths {
		if err := w.Watch(path); err != nil {
			log.Errorf("failed to watch dir [%s]\n", path)
			continue
		}
		log.Infof("directory [%s] is watched\n", path)
	}
}

// Exit ...
func (w *Watcher) Exit() {
	w.cmd.Exit()
}
