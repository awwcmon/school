// Code generated by https://github.com/go-dev-frame/sponge

package ecode

import (
	"github.com/go-dev-frame/sponge/pkg/errcode"
)

// fileService business-level http error codes.
// the fileServiceNO value range is 1~100, if the same error code is used, it will cause panic.
var (
	fileServiceNO       = 19
	fileServiceName     = "fileService"
	fileServiceBaseCode = errcode.HCode(fileServiceNO)

	ErrUploadFileFileService   = errcode.NewError(fileServiceBaseCode+1, "failed to UploadFile "+fileServiceName)
	ErrGetFileInfoFileService  = errcode.NewError(fileServiceBaseCode+2, "failed to GetFileInfo "+fileServiceName)
	ErrDownloadFileFileService = errcode.NewError(fileServiceBaseCode+3, "failed to DownloadFile "+fileServiceName)
	ErrCreateFileFileService   = errcode.NewError(fileServiceBaseCode+4, "failed to CreateFile "+fileServiceName)
	ErrGetFileFileService      = errcode.NewError(fileServiceBaseCode+5, "failed to GetFile "+fileServiceName)
	ErrListFilesFileService    = errcode.NewError(fileServiceBaseCode+6, "failed to ListFiles "+fileServiceName)
	ErrDeleteFileFileService   = errcode.NewError(fileServiceBaseCode+7, "failed to DeleteFile "+fileServiceName)
	ErrCreateFileFile          = errcode.NewError(fileServiceBaseCode+8, "failed to CreateFile "+fileName)
	ErrDownloadFileFile        = errcode.NewError(fileServiceBaseCode+9, "failed to DownloadFile "+fileName)

	// error codes are globally unique, adding 1 to the previous error code
)
