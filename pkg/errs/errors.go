package errs

import (
	"errors"
	"net/http"
)

var (
	ErrNotFound            = New(http.StatusNotFound, "未找到")
	ErrInternalServerError = New(http.StatusInternalServerError, "服务器内部错误")
	ErrBadRequest          = New(http.StatusBadRequest, "请求错误")
	ErrUnauthorized        = New(http.StatusUnauthorized, "未授权")
	ErrForbidden           = New(http.StatusForbidden, "无权限")
)

type CyBlogError struct {
	Code int
	Msg  string

	Err error
	//Stack []uintptr
}

//func captureStack() []uintptr {
//	const skip = 3 // 跳过 captureStack、NewBizError、调用者 三层
//	pcs := make([]uintptr, 32)
//	n := runtime.Callers(skip, pcs)
//	return pcs[:n]
//}

func (e *CyBlogError) Error() string {
	return e.Msg
}

//func (e *CyBlogError) Format(s fmt.State, verb rune) {
//	switch verb {
//	case 'v':
//		if s.Flag('+') {
//			fmt.Fprintf(s, "%s\n", e.Error())
//			// 打印调用栈
//			frames := runtime.CallersFrames(e.Stack)
//			for {
//				frame, more := frames.Next()
//				fmt.Fprintf(s, "  at %s:%d %s\n", frame.File, frame.Line, frame.Function)
//				if !more {
//					break
//				}
//			}
//			return
//		}
//		fallthrough
//	case 's':
//		fmt.Fprint(s, e.Error())
//	}
//}

func New(code int, msg string) *CyBlogError {
	return &CyBlogError{Code: code, Msg: msg}
}
func (e *CyBlogError) Unwrap() error {
	return e.Err
}

func Wrap(code int, msg string, err error) *CyBlogError {
	return &CyBlogError{Code: code, Msg: msg, Err: err}
}

func WrapWithMsg(code int, msg string, err error) *CyBlogError {
	return &CyBlogError{Code: code, Msg: msg + ": " + err.Error(), Err: err}
}
func Code(err error) int {
	if err == nil {
		return http.StatusOK
	}
	if e, ok := errors.AsType[*CyBlogError](err); ok {
		return e.Code
	}
	return http.StatusInternalServerError
}

// Cover 使用现成的CyBlogError包裹错误
func Cover(aboveErr *CyBlogError, err error) *CyBlogError {
	return &CyBlogError{Code: aboveErr.Code, Msg: aboveErr.Msg, Err: err}
}
