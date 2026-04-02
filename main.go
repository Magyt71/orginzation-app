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

var HomDir, _ = os.UserHomeDir()

// the extension and where it go
var theMap = map[string]string{
	".jpg": "Image", ".png": "Image", ".jpeg": "Image",
	".gif": "Image", ".svg": "Image", ".webp": "Image",
	".mp3": "Audio", ".wav": "Audio", ".ogg": "Audio",
	".mp4": "Video", ".mkv": "Video",
	".zip": "compressed", ".7z": "compressed", ".rar": "compressed",
	".tar": "compressed", ".gz": "compressed", ".bz2": "compressed",
	".ydk": "Yugioh",
	".pdf": "Text", ".docx": "Text", ".txt": "Text", ".log": "Text"}

func main() {

	thepath := filepath.Join(HomDir, "Downloads")

	// Creating the watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	var PrograssingFile sync.Map

	// making loop for the watcher
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				log.Println("event:", event)
				if event.Has(fsnotify.Write) || event.Has(fsnotify.Create) {

					if _, IsProssing := PrograssingFile.LoadOrStore(event.Name, true); !IsProssing {

						log.Println("New file detected:", event.Name)

						go func(FileName string) {
							log.Println("modified file:", FileName)
							MoveFile(FileName)

							PrograssingFile.Delete(FileName)
						}(event.Name)
					}
				}

			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(thepath)
	if err != nil {
		log.Println(err)
	}
	fmt.Println("it's working")
	done := make(chan bool)
	<-done
}

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

// Moving the file to the target folder
func MoveFile(FilePath string) {

	if _, err := os.Stat(FilePath); os.IsNotExist(err) {
		return
	}

	extension := strings.ToLower(filepath.Ext(FilePath))
	Targetfile, exists := theMap[extension]

	if exists {
		time.Sleep(10 *time.Second)
		FileName := filepath.Base(FilePath)
		beeep.Notify("we found a file to move it", FileName, "just waiting to make it done")

		for {
			// for make sure it's done downloading
			TheFile, err := os.OpenFile(FilePath, os.O_RDWR, 0)
			if err == nil {
				TheFile.Close()
				break
			}
			log.Println("still downloading")
			time.Sleep(20 * time.Minute)
		}

		if _, err := os.Stat(FilePath); os.IsNotExist(err) {
			return
		}

		Newpath := filepath.Join(HomDir, "Downloads", Targetfile)
		os.MkdirAll(Newpath, os.ModePerm)

		filename := filepath.Base(FilePath)
		finalpath := UniqPath(filepath.Join(Newpath, filename))

		err := os.Rename(FilePath, finalpath)
		if err != nil {
			log.Println("the file failed", err)
		} else {
			log.Println("the file pass", Targetfile)
		}
		beeep.Alert("the file has moved", FileName, "successfully")
	}
}
