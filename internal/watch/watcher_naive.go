//go:build !darwin
// +build !darwin

package watch

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/pkg/errors"

	"github.com/tilt-dev/fsnotify"
	"github.com/tilt-dev/tilt/internal/ospath"
	"github.com/tilt-dev/tilt/pkg/logger"
)

// A naive file watcher that uses the plain fsnotify API.
// Used on all non-Darwin systems (including Windows & Linux).
//
// All OS-specific codepaths are handled by fsnotify.
type naiveNotify struct {
	// Paths that we're watching that should be passed up to the caller.
	// Note that we may have to watch ancestors of these paths
	// in order to fulfill the API promise.
	//
	// We often need to check if paths are a child of a path in
	// the notify list. It might be better to store this in a tree
	// structure, so we can filter the list quickly.
	notifyList map[string]bool

	ignore PathMatcher
	log    logger.Logger

	isWatcherRecursive bool
	watcher            *fsnotify.Watcher
	events             chan fsnotify.Event
	wrappedEvents      chan FileEvent
	errors             chan error
	numWatches         int64
}

func (d *naiveNotify) Start() error {
	if len(d.notifyList) == 0 {
		return nil
	}

	pathsToWatch := []string{}
	for path := range d.notifyList {
		pathsToWatch = append(pathsToWatch, path)
	}

	pathsToWatch, err := greatestExistingAncestors(pathsToWatch)
	if err != nil {
		return err
	}
	if d.isWatcherRecursive {
		pathsToWatch = dedupePathsForRecursiveWatcher(pathsToWatch)
	}

	for _, name := range pathsToWatch {
		fi, err := os.Stat(name)
		if err != nil && !os.IsNotExist(err) {
			return errors.Wrapf(err, "notify.Add(%q)", name)
		}

		// if it's a file that doesn't exist,
		// we should have caught that above, let's just skip it.
		if os.IsNotExist(err) {
			continue
		} else if fi.IsDir() {
			err = d.watchRecursively(name)
			if err != nil {
				return errors.Wrapf(err, "notify.Add(%q)", name)
			}
		} else {
			err = d.add(filepath.Dir(name))
			if err != nil {
				return errors.Wrapf(err, "notify.Add(%q)", filepath.Dir(name))
			}
		}
	}

	go d.loop()

	return nil
}

func (d *naiveNotify) watchRecursively(dir string) error {
	if d.isWatcherRecursive {
		err := d.add(dir)
		if err == nil || os.IsNotExist(err) {
			return nil
		}
		return errors.Wrapf(err, "watcher.Add(%q)", dir)
	}

	return filepath.WalkDir(dir, func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			return nil
		}

		shouldSkipDir, err := d.shouldSkipDir(path)
		if err != nil {
			return err
		}

		if shouldSkipDir {
			return filepath.SkipDir
		}

		err = d.add(path)
		if err != nil {
			if os.IsNotExist(err) {
				return nil
			}
			return errors.Wrapf(err, "watcher.Add(%q)", path)
		}
		return nil
	})
}

func (d *naiveNotify) Close() error {
	numberOfWatches.Add(-d.numWatches)
	d.numWatches = 0
	return d.watcher.Close()
}

func (d *naiveNotify) Events() chan FileEvent {
	return d.wrappedEvents
}

func (d *naiveNotify) Errors() chan error {
	return d.errors
}

func (d *naiveNotify) loop() {
	defer close(d.wrappedEvents)
	for e := range d.events {
		// The Windows fsnotify event stream sometimes gets events with empty names
		// that are also sent to the error stream. Hmmmm...
		if e.Name == "" {
			continue
		}

		if e.Op&fsnotify.Create != fsnotify.Create {
			if !d.shouldNotify(e.Name) {
				continue
			}

			// Don't send events for directories when the modtime is being changed.
			//
			// This is a bit of a hack because every OS represents modtime updates
			// a bit differently and they don't map well to fsnotify events.
			//
			// On Windows, updating the modtime of a directory is a fsnotify.Write.
			// On Linux, it's a fsnotify.Chmod.
			isDirUpdateOnly := (e.Op == fsnotify.Write || e.Op == fsnotify.Chmod) &&
				ospath.IsDir(e.Name)
			if isDirUpdateOnly {
				continue
			}

			d.wrappedEvents <- FileEvent{e.Name}
			continue
		}

		if d.isWatcherRecursive {
			if !d.shouldNotify(e.Name) {
				continue
			}
			d.wrappedEvents <- FileEvent{e.Name}
			continue
		}

		// If the watcher is not recursive, we have to walk the tree
		// and add watches manually. We fire the event while we're walking the tree.
		// because it's a bit more elegant that way.
		//
		// TODO(dbentley): if there's a delete should we call d.watcher.Remove to prevent leaking?
		err := filepath.WalkDir(e.Name, func(path string, info fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if d.shouldNotify(path) {
				d.wrappedEvents <- FileEvent{path}
			}

			// TODO(dmiller): symlinks 😭

			shouldWatch := false
			if info.IsDir() {
				// watch directories unless we can skip them entirely
				shouldSkipDir, err := d.shouldSkipDir(path)
				if err != nil {
					return err
				}
				if shouldSkipDir {
					return filepath.SkipDir
				}

				shouldWatch = true
			} else {
				// watch files that are explicitly named, but don't watch others
				_, ok := d.notifyList[path]
				if ok {
					shouldWatch = true
				}
			}
			if shouldWatch {
				err := d.add(path)
				if err != nil && !os.IsNotExist(err) {
					d.log.Infof("Error watching path %s: %s", e.Name, err)
				}
			}
			return nil
		})
		if err != nil && !os.IsNotExist(err) {
			d.log.Infof("Error walking directory %s: %s", e.Name, err)
		}
	}
}

func (d *naiveNotify) shouldNotify(path string) bool {
	ignore, err := d.ignore.Matches(path)
	if err != nil {
		d.log.Infof("Error matching path %q: %v", path, err)
	} else if ignore {
		return false
	}

	if _, ok := d.notifyList[path]; ok {
		// We generally don't care when directories change at the root of an ADD
		isDir := ospath.IsDirLstat(path)
		if isDir {
			return false
		}
		return true
	}

	for root := range d.notifyList {
		if ospath.IsChild(root, path) {
			return true
		}
	}
	return false
}

func (d *naiveNotify) shouldSkipDir(path string) (bool, error) {
	// If path is directly in the notifyList, we should always watch it.
	if d.notifyList[path] {
		return false, nil
	}

	skip, err := d.ignore.MatchesEntireDir(path)
	if err != nil {
		return false, errors.Wrap(err, "shouldSkipDir")
	}

	if skip {
		return true, nil
	}

	// Suppose we're watching
	// /src/.tiltignore
	// but the .tiltignore file doesn't exist.
	//
	// Our watcher will create an inotify watch on /src/.
	//
	// But then we want to make sure we don't recurse from /src/ down to /src/node_modules.
	//
	// To handle this case, we only want to traverse dirs that are:
	// - A child of a directory that's in our notify list, or
	// - A parent of a directory that's in our notify list
	//   (i.e., to cover the "path doesn't exist" case).
	for root := range d.notifyList {
		if ospath.IsChild(root, path) || ospath.IsChild(path, root) {
			return false, nil
		}
	}
	return true, nil
}

func (d *naiveNotify) add(path string) error {
	err := d.watcher.Add(path)
	if err != nil {
		return err
	}
	d.numWatches++
	numberOfWatches.Add(1)
	return nil
}

func newWatcher(paths []string, ignore PathMatcher, l logger.Logger) (*naiveNotify, error) {
	if ignore == nil {
		return nil, fmt.Errorf("newWatcher: ignore is nil")
	}

	fsw, err := fsnotify.NewWatcher()
	if err != nil {
		if strings.Contains(err.Error(), "too many open files") && runtime.GOOS == "linux" {
			return nil, fmt.Errorf("Hit OS limits creating a watcher.\n" +
				"Run 'sysctl fs.inotify.max_user_instances' to check your inotify limits.\n" +
				"To raise them, run 'sudo sysctl fs.inotify.max_user_instances=1024'")
		}
		return nil, errors.Wrap(err, "creating file watcher")
	}
	MaybeIncreaseBufferSize(fsw)

	err = fsw.SetRecursive()
	isWatcherRecursive := err == nil

	wrappedEvents := make(chan FileEvent)
	notifyList := make(map[string]bool, len(paths))
	if isWatcherRecursive {
		paths = dedupePathsForRecursiveWatcher(paths)
	}
	for _, path := range paths {
		path, err := filepath.Abs(path)
		if err != nil {
			return nil, errors.Wrap(err, "newWatcher")
		}
		notifyList[path] = true
	}

	wmw := &naiveNotify{
		notifyList:         notifyList,
		ignore:             ignore,
		log:                l,
		watcher:            fsw,
		events:             fsw.Events,
		wrappedEvents:      wrappedEvents,
		errors:             fsw.Errors,
		isWatcherRecursive: isWatcherRecursive,
	}

	return wmw, nil
}

var _ Notify = &naiveNotify{}

func greatestExistingAncestors(paths []string) ([]string, error) {
	result := []string{}
	for _, p := range paths {
		newP, err := greatestExistingAncestor(p)
		if err != nil {
			return nil, fmt.Errorf("Finding ancestor of %s: %v", p, err)
		}
		result = append(result, newP)
	}
	return result, nil
}
