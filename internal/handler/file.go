// Code generated by https://github.com/go-dev-frame/sponge

package handler

import (
	"context"
	"github.com/go-dev-frame/sponge/pkg/gin/middleware"
	"io"
	"mime/multipart"
	"os"
	//"github.com/go-dev-frame/sponge/pkg/gin/middleware"

	schoolV1 "school/api/school/v1"
)

var _ schoolV1.FileLogicer = (*fileHandler)(nil)

type fileHandler struct{}

// NewFileHandler create a handler
func NewFileHandler() schoolV1.FileLogicer {
	return &fileHandler{}
}

// CreateFile ......
func (h *fileHandler) CreateFile(ctx context.Context, req *schoolV1.UploadFileRequest) (*schoolV1.UploadFileResponse, error) {
	c, _ := middleware.AdaptCtx(ctx)
	fh, err := c.FormFile("file")
	if err != nil {
		return nil, err
	}
	file, err := fh.Open()
	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)
	targetFile, err := os.OpenFile("uploads/"+fh.Filename, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return nil, err
	}
	defer func(targetFile *os.File) {
		err := targetFile.Close()
		if err != nil {
		}
	}(targetFile)
	_, err = io.Copy(targetFile, file)
	if err != nil {
		return nil, err
	}
	fileinfo, err := os.Stat(targetFile.Name())
	if err != nil {
		return nil, err
	}
	return &schoolV1.UploadFileResponse{
		FileId:   fileinfo.Name(),
		FileName: fileinfo.Name(),
		FileSize: fileinfo.Size(),
	}, nil
}
