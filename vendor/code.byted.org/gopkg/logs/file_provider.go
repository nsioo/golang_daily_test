package logs

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)


type SegDuration string

const (
	HourDur SegDuration = "Hour"
	DayDur  SegDuration = "Day"
	NoDur   SegDuration = "No"
)

const (
	FilePageCacheSize int64 = 16 * 1024 * 1024
	FullPageCacheSize int64 = FilePageCacheSize * 16
)

type FileProvider struct {
	sync.Mutex
	currentTimeSeg time.Time
	duration SegDuration


	fd            *os.File
	filename      string
	level         int
	fadvOffset    int64
	writeSize     int64
	fileKeepCount int
}

// NOTE(xiangchao): 由于基于size大小的切割方案可能会导致切割出来的文件在命名上存在一些问题，因此这里废弃基于size大小的切割方案
func NewFileProvider(filename string, dur SegDuration, size int64) *FileProvider {
	provider := &FileProvider{
		filename:     filename,
		level:        LevelDebug,
		duration: dur,
	}

	return provider
}

func (fp *FileProvider) Init() error {
	var (
		fd  *os.File
		err error
	)
	fp.currentTimeSeg = time.Now()
	realFile, err := fp.timeFilename()
	if err != nil {
		return err
	}
	_ = os.MkdirAll(filepath.Dir(realFile), 0755)
	if env := os.Getenv("IS_PROD_RUNTIME"); len(env) == 0 {
		fd, err = os.OpenFile(realFile, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	} else {
		fd, err = os.OpenFile(realFile, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	}
	fp.fd = fd
	_, err = os.Lstat(fp.filename)
	if err == nil || os.IsExist(err) {
		os.Remove(fp.filename)
	}
	os.Symlink("./"+filepath.Base(realFile), fp.filename)
	pos, err := fd.Seek(0, os.SEEK_CUR)
	if err != nil {
		fp.fadvOffset = -1
	} else {
		fp.fadvOffset = pos / FilePageCacheSize * FilePageCacheSize
		fp.writeSize = pos % FilePageCacheSize
	}
	err = fp.cleanElderFiles()
	if err != nil {
		fmt.Fprintf(os.Stderr, "clean file %s error: %s\n", fp.filename, err)
	}
	return nil
}

func (fp *FileProvider) doCheck(logTime time.Time) error {
	fp.Lock()
	defer fp.Unlock()

	needTruncate := false

	// assume that the application would log at least one message in each year
	switch fp.duration {
	case DayDur:
		if fp.currentTimeSeg.YearDay() != logTime.YearDay() {
			needTruncate = true
		}
	case HourDur:
		if fp.currentTimeSeg.Hour() != logTime.Hour() || fp.currentTimeSeg.YearDay() != logTime.YearDay() {
			needTruncate = true
		}
	}

	if needTruncate {
		if err := fp.truncate(); err != nil {
			fmt.Fprintf(os.Stderr, "truncate file %s error: %s\n", fp.filename, err)
			return err
		}
		if err := fp.cleanElderFiles(); err != nil {
			fmt.Fprintf(os.Stderr, "clean file %s error: %s\n", fp.filename, err)
		}
	}
	return nil
}

func (fp *FileProvider) SetLevel(l int) {
	fp.level = l
}

func (fp *FileProvider) WriteMsg(msg string, level int) error {
	if level < fp.level {
		return nil
	}
	// NOTE(xiangchao): 按照size切割已经被忽略
	fp.doCheck(time.Now())
	written, err := fmt.Fprint(fp.fd, msg)
	if (err == nil) && (fp.fadvOffset >= 0) {
		fp.writeSize += int64(written)
		if fp.writeSize >= FilePageCacheSize {
			go func(fd int, offset int64, length int64) {
				TryToDropFilePageCache(fd, offset, length)
				/* full drop page cache every FullPageCacheSize */
				if ((offset + FilePageCacheSize) % FullPageCacheSize) <= FilePageCacheSize {
					TryToDropFilePageCache(int(fp.fd.Fd()), 0, fp.fadvOffset)
				}
			}(int(fp.fd.Fd()), fp.fadvOffset, FilePageCacheSize)

			fp.fadvOffset += FilePageCacheSize
			fp.writeSize -= FilePageCacheSize
		}
	}
	return err
}

func (fp *FileProvider) Destroy() error {
	TryToDropFilePageCache(int(fp.fd.Fd()), 0, fp.fadvOffset+FilePageCacheSize)
	return fp.fd.Close()
}

func (fp *FileProvider) Flush() error {
	return fp.fd.Sync()
}

// 1: 拼接出新的日志文件的名字
// 2: 拷贝当前日志文件到新的文件
// 3: Truncate当前日志文件
func (fp *FileProvider) truncate() error {
	fp.fd.Close()
	return fp.Init()
}

func (fp *FileProvider) timeFilename() (string, error) {
	absPath, err := filepath.Abs(fp.filename)
	if err != nil {
		return "", err
	}
	return absPath + "." + fp.currentTimeSeg.Format("2006-01-02_15"), nil
}

func (fp *FileProvider) SetKeepFiles(count int) {
	if fp.fd == nil {
		fp.fileKeepCount = count
	}
}

func (fp *FileProvider) cleanElderFiles() error {
	if fp.fileKeepCount == 0 {
		return nil
	}
	absPath, err := filepath.Abs(fp.filename)
	if err != nil {
		return err
	}
	logDir := filepath.Dir(absPath)
	files, err := ioutil.ReadDir(logDir)
	if err != nil {
		return err
	}
	var logFiles []os.FileInfo
	for _, file := range files {
		if file.IsDir() || !strings.HasPrefix(file.Name(), filepath.Base(fp.filename)) {
			continue
		}
		logFiles = append(logFiles, file)
	}
	if len(logFiles) <= fp.fileKeepCount {
		return nil
	}
	sortableFiles := File(logFiles)

	sort.Sort(sortableFiles)

	for _, file := range logFiles[fp.fileKeepCount+1:] {
		relativePath := filepath.Join(filepath.Dir(fp.filename), file.Name())
		fullPath, err := filepath.Abs(relativePath)
		if err != nil {
			return err
		}
		err = os.Remove(fullPath)
		if err != nil {
			return err
		}
	}
	return nil
}

type File []os.FileInfo

func (a File) Len() int {
	return len(a)

}
func (a File) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]

}
func (a File) Less(i, j int) bool {
	return compareFileCreatedTime(a[i], a[j])
}
