// Code generated by protoc-gen-go-errors. DO NOT EDIT.

package v1

import (
	fmt "fmt"
	errors "github.com/go-kratos/kratos/v2/errors"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the kratos package it is being compiled against.
const _ = errors.SupportPackageIsVersion1

func IsUnknownError(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == NesgoWebAPIErrorReason_UNKNOWN_ERROR.String() && e.Code == 500
}

func ErrorUnknownError(format string, args ...interface{}) *errors.Error {
	return errors.New(500, NesgoWebAPIErrorReason_UNKNOWN_ERROR.String(), fmt.Sprintf(format, args...))
}

func IsLoginFailed(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == NesgoWebAPIErrorReason_LOGIN_FAILED.String() && e.Code == 500
}

func ErrorLoginFailed(format string, args ...interface{}) *errors.Error {
	return errors.New(500, NesgoWebAPIErrorReason_LOGIN_FAILED.String(), fmt.Sprintf(format, args...))
}

func IsUsernameConflict(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == NesgoWebAPIErrorReason_USERNAME_CONFLICT.String() && e.Code == 500
}

func ErrorUsernameConflict(format string, args ...interface{}) *errors.Error {
	return errors.New(500, NesgoWebAPIErrorReason_USERNAME_CONFLICT.String(), fmt.Sprintf(format, args...))
}

func IsRegisterFailed(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == NesgoWebAPIErrorReason_REGISTER_FAILED.String() && e.Code == 500
}

func ErrorRegisterFailed(format string, args ...interface{}) *errors.Error {
	return errors.New(500, NesgoWebAPIErrorReason_REGISTER_FAILED.String(), fmt.Sprintf(format, args...))
}

func IsOperationFailed(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == NesgoWebAPIErrorReason_OPERATION_FAILED.String() && e.Code == 500
}

func ErrorOperationFailed(format string, args ...interface{}) *errors.Error {
	return errors.New(500, NesgoWebAPIErrorReason_OPERATION_FAILED.String(), fmt.Sprintf(format, args...))
}
