package errcode

import (
	"errors"

	"github.com/gin-gonic/gin"
)

var (
	// 请求错误
	ParamErr        = publicErr(10001, "参数错误")
	FileAnalysisErr = publicErr(11001, "文件解析错误")
	FileWriteErr    = privateErr(11002, "文件保存出错，请重试")
	FileReadErr     = privateErr(11003, "文件读取错误，请重试")

	// 服务器错误
	ServerErr = publicErr(20001, "服务器内部错误")
)

type Err struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
}

func bindErr(code int32, message string) *gin.Error {
	return &gin.Error{
		Err:  errors.New("error"),
		Type: gin.ErrorTypeBind,
		Meta: Err{
			Code:    code,
			Message: message,
		},
	}
}

func renderErr(code int32, message string) *gin.Error {
	return &gin.Error{
		Err:  errors.New("error"),
		Type: gin.ErrorTypeRender,
		Meta: Err{
			Code:    code,
			Message: message,
		},
	}
}

func publicErr(code int32, message string) *gin.Error {
	return &gin.Error{
		Err:  errors.New("error"),
		Type: gin.ErrorTypePublic,
		Meta: Err{
			Code:    code,
			Message: message,
		},
	}
}

func privateErr(code int32, message string) *gin.Error {
	return &gin.Error{
		Err:  errors.New("error"),
		Type: gin.ErrorTypePrivate,
		Meta: Err{
			Code:    code,
			Message: message,
		},
	}
}

func anyErr(code int32, message string) *gin.Error {
	return &gin.Error{
		Err:  errors.New("error"),
		Type: gin.ErrorTypeAny,
		Meta: Err{
			Code:    code,
			Message: message,
		},
	}
}
