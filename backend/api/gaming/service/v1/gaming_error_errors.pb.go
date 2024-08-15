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

func IsRoomSessionNotAccessible(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == UserServiceErrorReason_ROOM_SESSION_NOT_ACCESSIBLE.String() && e.Code == 500
}

func ErrorRoomSessionNotAccessible(format string, args ...interface{}) *errors.Error {
	return errors.New(500, UserServiceErrorReason_ROOM_SESSION_NOT_ACCESSIBLE.String(), fmt.Sprintf(format, args...))
}

func IsCreateGameInstanceFailed(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == UserServiceErrorReason_CREATE_GAME_INSTANCE_FAILED.String() && e.Code == 500
}

func ErrorCreateGameInstanceFailed(format string, args ...interface{}) *errors.Error {
	return errors.New(500, UserServiceErrorReason_CREATE_GAME_INSTANCE_FAILED.String(), fmt.Sprintf(format, args...))
}

func IsGameFileNotFound(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == UserServiceErrorReason_GAME_FILE_NOT_FOUND.String() && e.Code == 500
}

func ErrorGameFileNotFound(format string, args ...interface{}) *errors.Error {
	return errors.New(500, UserServiceErrorReason_GAME_FILE_NOT_FOUND.String(), fmt.Sprintf(format, args...))
}

func IsSavedGameNotFound(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == UserServiceErrorReason_SAVED_GAME_NOT_FOUND.String() && e.Code == 500
}

func ErrorSavedGameNotFound(format string, args ...interface{}) *errors.Error {
	return errors.New(500, UserServiceErrorReason_SAVED_GAME_NOT_FOUND.String(), fmt.Sprintf(format, args...))
}

func IsUploadFileFailed(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == UserServiceErrorReason_UPLOAD_FILE_FAILED.String() && e.Code == 500
}

func ErrorUploadFileFailed(format string, args ...interface{}) *errors.Error {
	return errors.New(500, UserServiceErrorReason_UPLOAD_FILE_FAILED.String(), fmt.Sprintf(format, args...))
}

func IsInvalidGameFile(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == UserServiceErrorReason_INVALID_GAME_FILE.String() && e.Code == 500
}

func ErrorInvalidGameFile(format string, args ...interface{}) *errors.Error {
	return errors.New(500, UserServiceErrorReason_INVALID_GAME_FILE.String(), fmt.Sprintf(format, args...))
}

func IsListGameFailed(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == UserServiceErrorReason_LIST_GAME_FAILED.String() && e.Code == 500
}

func ErrorListGameFailed(format string, args ...interface{}) *errors.Error {
	return errors.New(500, UserServiceErrorReason_LIST_GAME_FAILED.String(), fmt.Sprintf(format, args...))
}

func IsDeleteGameFailed(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == UserServiceErrorReason_DELETE_GAME_FAILED.String() && e.Code == 500
}

func ErrorDeleteGameFailed(format string, args ...interface{}) *errors.Error {
	return errors.New(500, UserServiceErrorReason_DELETE_GAME_FAILED.String(), fmt.Sprintf(format, args...))
}

func IsOpenGameConnectionFailed(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == UserServiceErrorReason_OPEN_GAME_CONNECTION_FAILED.String() && e.Code == 500
}

func ErrorOpenGameConnectionFailed(format string, args ...interface{}) *errors.Error {
	return errors.New(500, UserServiceErrorReason_OPEN_GAME_CONNECTION_FAILED.String(), fmt.Sprintf(format, args...))
}

func IsGameConnectionNotFound(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == UserServiceErrorReason_GAME_CONNECTION_NOT_FOUND.String() && e.Code == 500
}

func ErrorGameConnectionNotFound(format string, args ...interface{}) *errors.Error {
	return errors.New(500, UserServiceErrorReason_GAME_CONNECTION_NOT_FOUND.String(), fmt.Sprintf(format, args...))
}

func IsSdpAnswerFailed(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == UserServiceErrorReason_SDP_ANSWER_FAILED.String() && e.Code == 500
}

func ErrorSdpAnswerFailed(format string, args ...interface{}) *errors.Error {
	return errors.New(500, UserServiceErrorReason_SDP_ANSWER_FAILED.String(), fmt.Sprintf(format, args...))
}

func IsIceCandidateFailed(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == UserServiceErrorReason_ICE_CANDIDATE_FAILED.String() && e.Code == 500
}

func ErrorIceCandidateFailed(format string, args ...interface{}) *errors.Error {
	return errors.New(500, UserServiceErrorReason_ICE_CANDIDATE_FAILED.String(), fmt.Sprintf(format, args...))
}

func IsGameInstanceNotAccessible(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == UserServiceErrorReason_GAME_INSTANCE_NOT_ACCESSIBLE.String() && e.Code == 500
}

func ErrorGameInstanceNotAccessible(format string, args ...interface{}) *errors.Error {
	return errors.New(500, UserServiceErrorReason_GAME_INSTANCE_NOT_ACCESSIBLE.String(), fmt.Sprintf(format, args...))
}

func IsOperationFailed(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == UserServiceErrorReason_OPERATION_FAILED.String() && e.Code == 500
}

func ErrorOperationFailed(format string, args ...interface{}) *errors.Error {
	return errors.New(500, UserServiceErrorReason_OPERATION_FAILED.String(), fmt.Sprintf(format, args...))
}

func IsSaveGameFailed(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == UserServiceErrorReason_SAVE_GAME_FAILED.String() && e.Code == 500
}

func ErrorSaveGameFailed(format string, args ...interface{}) *errors.Error {
	return errors.New(500, UserServiceErrorReason_SAVE_GAME_FAILED.String(), fmt.Sprintf(format, args...))
}

func IsLoadSaveFailed(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == UserServiceErrorReason_LOAD_SAVE_FAILED.String() && e.Code == 500
}

func ErrorLoadSaveFailed(format string, args ...interface{}) *errors.Error {
	return errors.New(500, UserServiceErrorReason_LOAD_SAVE_FAILED.String(), fmt.Sprintf(format, args...))
}

func IsRestartFailed(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == UserServiceErrorReason_RESTART_FAILED.String() && e.Code == 500
}

func ErrorRestartFailed(format string, args ...interface{}) *errors.Error {
	return errors.New(500, UserServiceErrorReason_RESTART_FAILED.String(), fmt.Sprintf(format, args...))
}
