package internal

/*
 * define a watcher to watch folders and files,
 * if there's any change the command will reload.
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
	fsWatcher *fsnotify.Watcher // *fsnotify.Watcher the actual file-wacher
	exit      chan<- bool       // exit channel to watcher notify main goroutine quit
	d         time.Duration     // watch duration
	cmd       *command.Command  // command to reload

	watchingPaths     []string // paths those are wanted to be watched
	watchingFiletype  []string // filetyps those are wanted to be watched
	unWatchingRegular []string // regular expressions those are not wanted to be watched

	evtTime map[string]int64 // record watching files modified time
}

// NewWatcher ...
func NewWatcher(paths []string, exit chan<- bool, watchingFiletype, unWatchingRegular []string) (*Watcher, error) {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	watcher := &Watcher{
		fsWatcher:         w,
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

// Watching ... before watching, must call `SetCommand` at first
func (w *Watcher) Watching() {
	if w.cmd == nil {
		panic("call watcher.SetCommand at first")
	}

	go func() {
		for {
			select {
			case evt := <-w.fsWatcher.Event:
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
			case err := <-w.fsWatcher.Error:
				log.Warnf("%s", err.Error())
				w.exit <- true
			default:
				time.Sleep(w.d)
			}
		}
	}()

	// append paths
	for _, path := range w.watchingPaths {
		if err := w.fsWatcher.Watch(path); err != nil {
			log.Errorf("failed to watch dir [%s]\n", path)
			continue
		}
		log.Infof("directory [%s] is watched\n", path)
	}
}

// Exit ...
func (w *Watcher) Exit() {
	if w.cmd != nil {
		w.cmd.Exit()
	}
	if err := w.fsWatcher.Close(); err != nil {
		panic(err)
	}
	// stop watching
}
