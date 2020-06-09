package modules

import (
	"github.com/pkg/errors"
	"github.com/radovskyb/watcher"
	"strings"
	"time"
	"yym/hydra_extension/hydra"
)

func init() {
	hydra.RegisterType("FileWatcher", newFileWatcher)
}

func newFileWatcher() (*fileWatcher, error) {
	watcher := watcher.New()

	return &fileWatcher{watcher: watcher}, nil
}

type fileWatcher struct {
	watcher *watcher.Watcher
}

func (self *fileWatcher) Add(path string) {
	self.watcher.Add(path)
}

func (self *fileWatcher) Close() {
	self.watcher.Close()
}

func (self *fileWatcher) Start() {
	go func() {
		self.watcher.Start(time.Second)
	}()
}

func (self *fileWatcher) Events() (string, string, error) {
	for {
		select {
		case event := <-self.watcher.Event:
			return strings.ToLower(event.Op.String()), event.Path, nil
		case err := <-self.watcher.Error:
			return "", "", err
		case <-self.watcher.Closed:
			return "", "", errors.Errorf("closed")
		}
	}
}
