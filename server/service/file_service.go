package service

import (
	"fmt"
	"os"
	"path/filepath"

	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// go-kit services

type FileInfo struct {
	Name       string
	ModifyTime string
	Size       int64
}

type FileAccessor interface {
	GetFileList(string) ([]FileInfo, error)
}

type CommonFileAccessor struct{}

func (fa *CommonFileAccessor) GetFileList(dirName string) ([]FileInfo, error) {
	var files = []FileInfo{}
	relativePath := "./" + dirName + "/"
	err := filepath.Walk(relativePath, func(_ string, info os.FileInfo, err error) error {
		if info == nil {
			return err
		}
		if !info.IsDir() {
			fi := FileInfo{
				Name:       info.Name(),
				ModifyTime: info.ModTime().Format("2006-01-02 15:04:05"),
				Size:       info.Size(),
			}

			files = append(files, fi)
		}
		return err
	})
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	fmt.Println(files)
	return files, nil
}

// gin services
const (
	dirKey  = "dirname"
	fileKey = "filename"

	MAX_DIR_SIZE = 3 * 1024 * 1024 * 1024
)

func makeFileRelativePath(d, f string) string {
	return "./" + d + "/" + f
}

func checkFileExist(p string) bool {
	_, err := os.Stat(p)
	if err != nil {
		fmt.Printf("file[%v] doesn't exist\n", p)
		return false
	}
	return true
}

func getDirSize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if !info.IsDir() && info.Name() != "" {
			size += info.Size()
		}
		return err
	})
	return size, err
}

func FileRequestChecker() gin.HandlerFunc {
	return func(c *gin.Context) {
		dirname := c.Param(dirKey)
		fname := c.Param(fileKey)
		if fname == "" || dirname == "" {
			c.String(http.StatusBadRequest, "Invalid request")
			c.Abort()
		}

		p := makeFileRelativePath(dirname, fname)
		if !checkFileExist(p) {
			c.String(http.StatusNotFound, "Request file doesn't exist")
			c.Abort()
		}

		c.Set(dirKey, dirname)
		c.Set(fileKey, fname)

		c.Next()
	}
}

func DownloadFile(c *gin.Context) {
	dirname := c.MustGet(dirKey).(string)
	fname := c.MustGet(fileKey).(string)
	fpath := makeFileRelativePath(dirname, fname)

	c.File(fpath)
}

func DeleteFile(c *gin.Context) {
	dirname := c.MustGet(dirKey).(string)
	fname := c.MustGet(fileKey).(string)
	fpath := makeFileRelativePath(dirname, fname)

	err := os.Remove(fpath)
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("filename exist, but remove failed."))
		c.Abort()
		return
	}

	c.String(200, "file delete success")
}

func UploadFile(c *gin.Context) {
	// check request dirname
	dirname := c.Param(dirKey)
	fmt.Println("dirname: ", dirname)
	if dirname == "" || !checkFileExist(dirname) {
		c.String(http.StatusForbidden, "dir name doesn't exist")
		c.Abort()
		return
	}

	// get uploaded files
	form, err := c.MultipartForm()
	if err != nil {
		c.String(http.StatusForbidden, "invalid request")
		c.Abort()
		return
	}

	files := form.File["upload"]
	if len(files) == 0 {
		c.String(http.StatusForbidden, fmt.Sprintf("No file recieved, please select files."))
		c.Abort()
		return

	}

	// check rest file size
	var totalUploadSize int64
	for _, file := range files {
		if strings.Contains(file.Filename, "/") || len(file.Filename) > 100 {
			c.String(http.StatusForbidden, fmt.Sprintf("Filename is too long, max length is 100 characters."))
			return
		}
		totalUploadSize += file.Size
	}

	curDirSize, _ := getDirSize(dirname)
	restDirSize := MAX_DIR_SIZE - curDirSize

	if totalUploadSize > restDirSize {
		c.String(http.StatusForbidden, fmt.Sprintf("Storage not enough, please delete unused files"))
		c.Abort()
		return
	}

	// save uploaded files
	for _, file := range files {
		path := makeFileRelativePath(dirname, file.Filename)
		if checkFileExist(path) {
			c.String(http.StatusForbidden, fmt.Sprintf("Filename already exist, change the filename or delete exist file."))
			c.Abort()
			return
		}

		c.SaveUploadedFile(file, path)
	}

	c.String(http.StatusOK, fmt.Sprintf("%d files uploaded!", len(files)))
}
