// Code generated by https://github.com/go-dev-frame/sponge

package ecode

import (
	"github.com/go-dev-frame/sponge/pkg/errcode"
)

// file business-level http error codes.
// the fileNO value range is 1~100, if the same error code is used, it will cause panic.
var (
	fileNO       = 48
	fileName     = "file"
	fileBaseCode = errcode.HCode(fileNO)

	ErrCreateFileFile   = errcode.NewError(fileBaseCode+1, "failed to CreateFile "+fileName)
	ErrDownloadFileFile = errcode.NewError(fileBaseCode+2, "failed to DownloadFile "+fileName)

	// error codes are globally unique, adding 1 to the previous error code
)
