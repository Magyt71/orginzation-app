package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/joho/godotenv"
)

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

	HomDir, _ := os.UserHomeDir()
	thepath := filepath.Join(HomDir, "Downloads")

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				log.Println("event:", event)
				if event.Has(fsnotify.Write) || event.Has(fsnotify.Create) {
					log.Println("modified file:", event.Name)
					go MoveFile(event.Name)
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

func MoveFile(FilePath string) {
	extension := strings.ToLower(filepath.Ext(FilePath))

	HomDir, _ := os.UserHomeDir()

	Targetfile, exists := theMap[extension]
	if exists {
		time.Sleep(5 * time.Second)
		Newpath := filepath.Join(HomDir, "Downloads", Targetfile)
		os.MkdirAll(Newpath, os.ModePerm)

		filename := filepath.Base(FilePath)
		finalpath := filepath.Join(Newpath, filename)

		err := os.Rename(FilePath, finalpath)
		if err != nil {
			log.Println("the file failed", err)
		} else {
			log.Println("the file pass", Targetfile)
		}
	}
	if !exists {

	}
}
