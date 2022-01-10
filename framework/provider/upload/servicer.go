package upload

import (
	"bufio"
	"errors"
	"fmt"
	"hade/framework"
	"hade/framework/util"
	"io"
	"os"
	"path"
	"time"
)

type HadeUpload struct {
	c           framework.Container // 容器
	folder      string              // 存储文件夹
	maxFileSize int64               // 最大文件大小
}

func NewHadeUpload(params ...interface{}) (interface{}, error) {
	container := params[0].(framework.Container)
	folder := params[1].(string)
	maxFileSize := params[2].(int64)

	hadeUpload := &HadeUpload{
		c:           container,
		folder:      folder,
		maxFileSize: maxFileSize,
	}

	return hadeUpload, nil
}

func (hadeUpload *HadeUpload) Upload(subFolder, fileName string, reader *io.Reader) error {
	uploadPath := path.Join(hadeUpload.folder, subFolder)

	if !util.Exists(uploadPath) {
		os.MkdirAll(uploadPath, os.ModePerm)
	}

	if err := upload(uploadPath, fileName, hadeUpload.maxFileSize, reader); err != nil {
		return err
	}

	return nil
}

func upload(uploadPath, fileName string, maxFileSize int64, reader *io.Reader) error {
	fi, err := os.Create(fmt.Sprintf("%s/%s-%s", uploadPath, fileName, time.Now().Format("hh-mm-ss")))
	if err != nil {
		return err
	}
	defer fi.Close()

	br := bufio.NewReader(*reader)
	if br.Size() > int(maxFileSize) {
		return errors.New("file is over maxFileSize")
	}

	if _, err := br.WriteTo(fi); err != nil {
		return err
	}

	return nil
}
