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
	return e.Reason == UserServiceErrorReason_UNKNOWN_ERROR.String() && e.Code == 500
}

func ErrorUnknownError(format string, args ...interface{}) *errors.Error {
	return errors.New(500, UserServiceErrorReason_UNKNOWN_ERROR.String(), fmt.Sprintf(format, args...))
}

func IsRoomNotFound(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == UserServiceErrorReason_ROOM_NOT_FOUND.String() && e.Code == 404
}

func ErrorRoomNotFound(format string, args ...interface{}) *errors.Error {
	return errors.New(404, UserServiceErrorReason_ROOM_NOT_FOUND.String(), fmt.Sprintf(format, args...))
}

func IsRoomNotAccessible(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == UserServiceErrorReason_ROOM_NOT_ACCESSIBLE.String() && e.Code == 403
}

func ErrorRoomNotAccessible(format string, args ...interface{}) *errors.Error {
	return errors.New(403, UserServiceErrorReason_ROOM_NOT_ACCESSIBLE.String(), fmt.Sprintf(format, args...))
}

func IsCreateRoomFailed(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == UserServiceErrorReason_CREATE_ROOM_FAILED.String() && e.Code == 400
}

func ErrorCreateRoomFailed(format string, args ...interface{}) *errors.Error {
	return errors.New(400, UserServiceErrorReason_CREATE_ROOM_FAILED.String(), fmt.Sprintf(format, args...))
}

func IsGetRoomFailed(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == UserServiceErrorReason_GET_ROOM_FAILED.String() && e.Code == 500
}

func ErrorGetRoomFailed(format string, args ...interface{}) *errors.Error {
	return errors.New(500, UserServiceErrorReason_GET_ROOM_FAILED.String(), fmt.Sprintf(format, args...))
}

func IsJoinRoomFailed(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == UserServiceErrorReason_JOIN_ROOM_FAILED.String() && e.Code == 500
}

func ErrorJoinRoomFailed(format string, args ...interface{}) *errors.Error {
	return errors.New(500, UserServiceErrorReason_JOIN_ROOM_FAILED.String(), fmt.Sprintf(format, args...))
}

func IsCreateRoomSessionFailed(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == UserServiceErrorReason_CREATE_ROOM_SESSION_FAILED.String() && e.Code == 500
}

func ErrorCreateRoomSessionFailed(format string, args ...interface{}) *errors.Error {
	return errors.New(500, UserServiceErrorReason_CREATE_ROOM_SESSION_FAILED.String(), fmt.Sprintf(format, args...))
}

func IsUpdateRoomFailed(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == UserServiceErrorReason_UPDATE_ROOM_FAILED.String() && e.Code == 500
}

func ErrorUpdateRoomFailed(format string, args ...interface{}) *errors.Error {
	return errors.New(500, UserServiceErrorReason_UPDATE_ROOM_FAILED.String(), fmt.Sprintf(format, args...))
}

func IsDeleteRoomFailed(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == UserServiceErrorReason_DELETE_ROOM_FAILED.String() && e.Code == 500
}

func ErrorDeleteRoomFailed(format string, args ...interface{}) *errors.Error {
	return errors.New(500, UserServiceErrorReason_DELETE_ROOM_FAILED.String(), fmt.Sprintf(format, args...))
}

func IsGetRoomMemberFailed(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == UserServiceErrorReason_GET_ROOM_MEMBER_FAILED.String() && e.Code == 500
}

func ErrorGetRoomMemberFailed(format string, args ...interface{}) *errors.Error {
	return errors.New(500, UserServiceErrorReason_GET_ROOM_MEMBER_FAILED.String(), fmt.Sprintf(format, args...))
}
