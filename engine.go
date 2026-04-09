package main

import (
	"fmt"
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
	WatchPaths  []string
	TargetMap   map[string]string
	WaitTime    time.Duration
	Notifcation bool
	IsRunning   bool
	mu          sync.Mutex
	Rmu         sync.RWMutex
}

type MoveLog struct {
	FileName string
	Dest     string
	Time     string
}

type Organizer struct {
	Config          *AppConfig
	Watcher         *fsnotify.Watcher
	PrograssingFile sync.Map
	logCallback     func(string)
	RecentMoves     []MoveLog
}

func NewOrganizer() *Organizer {
	HomDir, _ := os.UserHomeDir()
	return &Organizer{
		Config: &AppConfig{
			WatchPaths:  []string{filepath.Join(HomDir, "Downloads")},
			Notifcation: true,
			WaitTime:    3 * time.Second,
			TargetMap: map[string]string{
				".jpg": "Image", ".png": "Image", ".jpeg": "Image",
				".gif": "Image", ".svg": "Image", ".webp": "Image",
				".mp3": "Audio", ".wav": "Audio", ".ogg": "Audio",
				".mp4": "Video", ".mkv": "Video",
				".zip": "Compressed", ".7z": "Compressed", ".rar": "Compressed",
				".tar": "Compressed", ".gz": "Compressed", ".bz2": "Compressed",
				".ydk": "Yugioh",
				".pdf": "Documents", ".docx": "Documents", ".txt": "Documents", ".log": "Documents"},
		},
	}
}

func (O *Organizer) AddPath(Path string) {
	O.Config.mu.Lock()
	defer O.Config.mu.Unlock()

	for _, p := range O.Config.WatchPaths {
		if p == Path {
			return
		}
	}

	O.Config.WatchPaths = append(O.Config.WatchPaths, Path)
	if O.Watcher != nil {
		O.Watcher.Add(Path)
	}
}

func (O *Organizer) log(msg string) {
	ts := time.Now().Format("15:04:05")
	line := fmt.Sprintf("[%s] %s", ts, msg)
	fmt.Println(line)
	if O.logCallback != nil {
		O.logCallback(line)
	}
}

func (O *Organizer) Start() {
	O.Config.mu.Lock()
	if O.Config.IsRunning {
		O.Config.mu.Unlock()
		return
	}

	var err error
	O.Watcher, err = fsnotify.NewWatcher()
	if err != nil {
		O.log("ERROR creating watcher: " + err.Error())
		O.Config.mu.Unlock()
		return
	}
	O.Config.IsRunning = true
	O.Config.mu.Unlock()

	go func() {
		for {
			select {
			case event, ok := <-O.Watcher.Events:
				if !ok {
					return
				}
				// نراقب إنشاء الملفات أو الكتابة النهائية فيها
				if event.Has(fsnotify.Create) || event.Has(fsnotify.Write) {
					O.HandleFileEvent(event.Name)
				}
			case err, ok := <-O.Watcher.Errors:
				if !ok {
					return
				}
				O.log("Watcher error: " + err.Error())
			}
		}
	}()

	for _, path := range O.Config.WatchPaths {
		err = O.Watcher.Add(path)
		if err != nil {
			O.log("ERROR adding path: " + err.Error())
			continue
		}
		O.log("Watching: " + path)
	}
}

func (O *Organizer) Stop() {
	O.Config.mu.Lock()
	defer O.Config.mu.Unlock()

	if O.Watcher != nil {
		O.Watcher.Close()
		O.Watcher = nil
	}
	O.Config.IsRunning = false
	O.log("Organizer stopped.")
}

func (O *Organizer) HandleFileEvent(FilePath string) {
	// التحقق أن الملف ليس مجلداً وأنه موجود
	info, err := os.Stat(FilePath)
	if err != nil || info.IsDir() {
		return
	}

	extension := strings.ToLower(filepath.Ext(FilePath))

	O.Config.Rmu.RLock()
	TargetFolder, exists := O.Config.TargetMap[extension]
	O.Config.Rmu.RUnlock()

	if exists {
		// منع معالجة نفس الملف عدة مرات في نفس الوقت
		if _, busy := O.PrograssingFile.LoadOrStore(FilePath, true); !busy {
			go func() {
				defer O.PrograssingFile.Delete(FilePath)
				// انتظار بسيط للتأكد من اكتمال تحميل الملف
				time.Sleep(O.Config.WaitTime)
				O.ProcessMove(FilePath, TargetFolder)
			}()
		}
	}
}

func (O *Organizer) ProcessMove(FilePath string, TargetFolder string) {
	FileName := filepath.Base(FilePath)
	CurrentDir := filepath.Dir(FilePath)
	
	// محاولة فتح الملف للتأكد أن النظام انتهى من كتابته (للملفات الكبيرة)
	for i := 0; i < 5; i++ {
		file, err := os.OpenFile(FilePath, os.O_RDWR, 0644)
		if err == nil {
			file.Close()
			break
		}
		time.Sleep(1 * time.Second)
	}

	destDir := filepath.Join(CurrentDir, TargetFolder)
	os.MkdirAll(destDir, os.ModePerm)

	finalPath := UniqPath(filepath.Join(destDir, FileName))

	if err := os.Rename(FilePath, finalPath); err != nil {
		O.log("ERROR moving " + FileName + ": " + err.Error())
		return
	}

	O.log(fmt.Sprintf("Moved: %s → %s", FileName, TargetFolder))
	
	if O.Config.Notifcation {
		beeep.Notify("File Organizer", "Moved: "+FileName+" to "+TargetFolder, "")
	}

	// تحديث قائمة السجلات بشكل آمن
	O.Config.mu.Lock()
	newLog := MoveLog{FileName: FileName, Dest: TargetFolder, Time: time.Now().Format("15:04")}
	O.RecentMoves = append([]MoveLog{newLog}, O.RecentMoves...)
	if len(O.RecentMoves) > 5 {
		O.RecentMoves = O.RecentMoves[:5]
	}
	O.Config.mu.Unlock()
}

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
