package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/gen2brain/beeep"
)

type AppConfig struct {
	WatchPath   string
	TargetMap   map[string]string
	WaitTime    time.Duration
	Notifcation bool
	IsRunning   bool
	mu          sync.Mutex
	Rmu         sync.RWMutex
}

type Organizer struct {
	Config          *AppConfig
	Watcher         *fsnotify.Watcher
	PrograssingFile sync.Map
}

func NewOrganizer() *Organizer {
	var HomDir, _ = os.UserHomeDir()
	return &Organizer{
		Config: &AppConfig{
			WatchPath:   filepath.Join(HomDir, "Downloads"),
			Notifcation: true,
			WaitTime:    4 * time.Second,
			TargetMap: map[string]string{
				".jpg": "Image", ".png": "Image", ".jpeg": "Image",
				".gif": "Image", ".svg": "Image", ".webp": "Image",
				".mp3": "Audio", ".wav": "Audio", ".ogg": "Audio",
				".mp4": "Video", ".mkv": "Video",
				".zip": "compressed", ".7z": "compressed", ".rar": "compressed",
				".tar": "compressed", ".gz": "compressed", ".bz2": "compressed",
				".ydk": "Yugioh",
				".pdf": "Text", ".docx": "Text", ".txt": "Text", ".log": "Text"},
		},
	}
}

func (O *Organizer) Start() {
	O.Config.mu.Lock()
	O.Config.IsRunning = true
	O.Config.mu.Unlock()

	var err error
	O.Watcher, err = fsnotify.NewWatcher()
	if err != nil {
		log.Fatalln("there an error here", err)
	}

	go func() {
		for {
			select {
			case event, ok := <-O.Watcher.Events:
				if !ok {
					return
				}

				if event.Has(fsnotify.Write) || event.Has(fsnotify.Create) {
					O.HandleFileEvent(event.Name)
				}
			case err, ok := <-O.Watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	err = O.Watcher.Add(O.Config.WatchPath)
	if err != nil {
		log.Println("error adding the path", err)
	}
	fmt.Println("the app working now:", O.Config.WatchPath)
}

func (O *Organizer) HandleFileEvent(FilePath string) {

	if filepath.Dir(FilePath) != O.Config.WatchPath {
		return
	}
	extension := strings.ToLower(filepath.Ext(FilePath))

	O.Config.Rmu.RLock()
	TargetFolder, exists := O.Config.TargetMap[extension]
	O.Config.Rmu.RUnlock()

	if exists {
		if _, busy := O.PrograssingFile.LoadOrStore(FilePath, true); !busy {
			go func() {
				O.ProcessMove(FilePath, TargetFolder)
				time.Sleep(2 * time.Second)
				O.PrograssingFile.Delete(FilePath)
			}()
		}
		defer O.PrograssingFile.Delete(FilePath)
	}

}

func (O *Organizer) ProcessMove(FilePath string, TargetFolder string) {

	FileName := filepath.Base(FilePath)
	if O.Config.Notifcation == true {
		beeep.Notify("we found a file to move it", FileName, "just waiting to make it done")
	}

	for {
		file, err := os.OpenFile(FilePath, os.O_RDWR, 0644)
		if err == nil {
			file.Close()
			break
		}
		time.Sleep(2 * time.Second)
	}

	destDir := filepath.Join(O.Config.WatchPath, TargetFolder)
	os.MkdirAll(destDir, os.ModePerm)

	finalPath := UniqPath(filepath.Join(destDir, FileName))

	if err := os.Rename(FilePath, finalPath); err == nil {
		if O.Config.Notifcation == true {
			beeep.Alert("the file has moved", FileName, "successfully")
		}
	}

}

var HomDir, err = os.UserHomeDir()

// for make sure it's don't delete anyfile with the same name
func UniqPath(FilePath string) string {
	Dir := filepath.Dir(FilePath)
	FileName := filepath.Base(FilePath)
	Ext := filepath.Ext(FilePath)
	NameOnly := strings.TrimSuffix(FileName, Ext)

	re := regexp.MustCompile(`\(\d+\)$`)
	CleanName := re.ReplaceAllString(NameOnly, "")

	FinalPath := FilePath
	count := 1

	for {
		if _, err := os.Stat(FinalPath); os.IsNotExist(err) {
			break
		}

		NewName := fmt.Sprintf("%s(%d)%s", CleanName, count, Ext)
		FinalPath = filepath.Join(Dir, NewName)
		count++
	}
	return FinalPath
}

func main() {
	Org := NewOrganizer()
	Org.Start()

	select {}
}
