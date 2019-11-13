package errors

import (
	"fmt"

	"golang.org/x/xerrors"
)

// 定义通用错误
var (
	New               = xerrors.New
	Is                = xerrors.Is
	As                = xerrors.As
	Unwrap            = xerrors.Unwrap // 获取内层错误
	DbMsg             = "服务器错误"
	Casbin            = "权限系统错误"
	Code              = "健壮验证错误"
	ErrForbidden      = New("禁止访问")
	ErrNotFound       = New("资源不存在")
	ErrBadRequest     = New("请求无效")
	ErrUnauthorized   = New("未授权")
	ErrInternalServer = New("服务器错误")
	ErrDatabase       = New("数据库请求错误")
	ErrExistsfail     = New("扩展存储错误")
	ErrCodefail       = New("编码错误")
)

// NewBadRequestError 创建请求无效错误
func NewBadRequestError(msg ...string) error {
	return newMessageError(ErrBadRequest, 1000, xerrors.Caller(1), msg...)
}

// NewUnauthorizedError 创建未授权错误
func NewUnauthorizedError(msg ...string) error {
	return newMessageError(ErrUnauthorized, 2000, xerrors.Caller(1), msg...)
}

// NewForbiddenError 创建资源禁止访问错误
func NewForbiddenError(msg ...string) error {
	return newMessageError(ErrForbidden, 3000, xerrors.Caller(1), msg...)
}

// NewNotFoundError 创建资源不存在错误
func NewNotFoundError(msg ...string) error {
	return newMessageError(ErrNotFound, 4000, xerrors.Caller(1), msg...)
}

// NewInternalServerError 创建服务器错误
func NewInternalServerError(msg ...string) error {
	return newMessageError(ErrInternalServer, 5000, xerrors.Caller(1), msg...)
}

// ErrDatabase 创建数据库错误
func NewInternalDatabaseError(msg ...string) error {
	return newMessageError(ErrDatabase, 6000, xerrors.Caller(1), msg...)
}

//NewFileErr 文件错误
func NewFileErr(msg ...string) error {
	return newMessageError(ErrNotFound, 7000, xerrors.Caller(1), msg...)
}

//NewCodeErr 编码错误
func NewCodeErr(msg ...string) error {
	return newMessageError(ErrCodefail, 8000, xerrors.Caller(1), msg...)
}

//NewDataErr 数据获取错误
func NewDataErr(code int, msg ...string) error {
	return newMessageError(ErrExistsfail, code, xerrors.Caller(1), msg...)
}

// NewMessageError 创建自定义消息错误
func NewMessageError(parent error, code int, msg ...string) (err error) {
	if code == 0 {
		code = 999
	}
	err = newMessageError(parent, code, xerrors.Caller(1), msg...)
	return
}

// newMessageError 创建自定义消息错误
func newMessageError(parent error, code int, frame xerrors.Frame, msg ...string) (err error) {
	if parent == nil {
		return nil
	}

	m := parent.Error()
	if len(msg) > 0 {
		m = msg[0]
	}
	err = &MessageError{parent, code, m, frame}

	return
}

// MessageError 自定义消息错误
type MessageError struct {
	err   error
	code  int
	msg   string
	frame xerrors.Frame
}

//func (e *MessageError) MarshalJSON() ([]byte, error) {
//	return []byte(fmt.Sprintf(`"%s"`,e.Error())), nil
//}
func (e *MessageError) Error() string {
	return fmt.Sprint(e)
}

func (e *MessageError) Format(f fmt.State, c rune) {
	xerrors.FormatError(e, f, c)
}

func (e *MessageError) Unwrap() error {
	return e.Parent()
}

func (e *MessageError) FormatError(p xerrors.Printer) (next error) {
	p.Printf("[%d] %s", e.code, e.msg)
	e.frame.Format(p)
	next = e.Parent()
	return
}

// Parent 父级错误
func (e *MessageError) Parent() error {
	return e.err
}
