package internal

/*
 * define a watcher to watch folders and files,
 * if there's any change the command will reload.
 */

import (
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/yeqown/go-watcher/internal/command"
	"github.com/yeqown/go-watcher/internal/log"
	"github.com/yeqown/go-watcher/utils"

	"github.com/howeyc/fsnotify"
)

// WatcherOption ...
type WatcherOption struct {
	D                 int      `yaml:"duration"`
	IncludedFiletypes []string `yaml:"included_filetypes"`
	ExcludedRegexps   []string `yaml:"excluded_regexps"`
}

// Watcher ... to watch and run command
type Watcher struct {
	fsWatcher *fsnotify.Watcher // *fsnotify.Watcher the actual file-wacher
	exit      chan<- bool       // exit channel to watcher notify main goroutine quit
	d         time.Duration     // watch duration
	cmd       *command.Command  // command to reload

	watchingPaths []string // paths those are wanted to be watched
	// watchingFiletype  []string // filetyps those are wanted to be watched
	// unWatchingRegular []string // regular expressions those are not wanted to be watched

	filetypeIncludeChecker IncludeChecker
	regularExcludeChecker  ExcludeChecker

	evtTime map[string]int64 // record watching files modified time
}

// NewWatcher ...
func NewWatcher(paths []string, exit chan<- bool, opt *WatcherOption) (*Watcher, error) {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	watcher := &Watcher{
		fsWatcher: w,
		exit:      exit,
		d:         time.Duration(opt.D) * time.Millisecond,

		watchingPaths: paths,
		// watchingFiletype:  watchingFiletype,
		// unWatchingRegular: unWatchingRegular,

		filetypeIncludeChecker: newIncludeChecker(opt.IncludedFiletypes, filetypeIncludePre),
		regularExcludeChecker:  newExcludeChecker(opt.ExcludedRegexps, regularExcludePre),

		evtTime: make(map[string]int64),
	}

	return watcher, nil
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
				// record modified time and
				// controls event notify would not repeated over times
				mt := utils.GetFileModTime(evt.Name)
				if t, ok := w.evtTime[evt.Name]; ok &&
					utils.UnixTimeDuration(mt, t) <= w.d {
					continue
				}

				// skip un-target file event
				if !w.filetypeIncludeChecker.Include(evt.Name, filetypeIncludeJudge) {
					log.Infof("(%s) is skipped, not target filetype\n", evt.Name)
					continue
				}
				if w.regularExcludeChecker.Exclude(evt.Name, regularExcludeJudge) {
					log.Infof("(%s) is skipped, not target file\n", evt.Name)
					continue
				}

				// if event is target files event
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
			log.Errorf("failed to watch dir (%s)\n", path)
			continue
		}
		log.Infof("directory (%s) is under watching\n", path)
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

var (
	_ PreFunc   = filetypeIncludePre
	_ PreFunc   = regularExcludePre
	_ JudgeFunc = filetypeIncludeJudge
	_ JudgeFunc = regularExcludeJudge

	errInvalidFiltypeString = errors.New("invalid filetype string to parse")
)

// only save file extend name into map
// ".go" save "go"
func filetypeIncludePre(filetype string) (string, error) {
	parsed := strings.Split(filetype, ".")
	if len(parsed) == 1 {
		return "", errInvalidFiltypeString
	}
	return parsed[len(parsed)-1], nil
}

// ...
func filetypeIncludeJudge(s string, m map[string]bool) bool {
	parsed := strings.Split(s, ".")
	_, ok := m[parsed[len(parsed)-1]]
	return ok
}

func regularExcludePre(regular string) (string, error) {
	return regular, nil
}

func regularExcludeJudge(s string, m map[string]bool) bool {
	for regular := range m {
		m, err := regexp.Match(regular, []byte(s))
		if err == nil && m {
			return true
		}
	}
	return false
}
