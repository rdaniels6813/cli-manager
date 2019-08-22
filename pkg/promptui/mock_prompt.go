// Code generated by MockGen. DO NOT EDIT.
// Source: prompt.go

// Package promptui is a generated GoMock package.
package promptui

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockPrompter is a mock of Prompter interface
type MockPrompter struct {
	ctrl     *gomock.Controller
	recorder *MockPrompterMockRecorder
}

// MockPrompterMockRecorder is the mock recorder for MockPrompter
type MockPrompterMockRecorder struct {
	mock *MockPrompter
}

// NewMockPrompter creates a new mock instance
func NewMockPrompter(ctrl *gomock.Controller) *MockPrompter {
	mock := &MockPrompter{ctrl: ctrl}
	mock.recorder = &MockPrompterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockPrompter) EXPECT() *MockPrompterMockRecorder {
	return m.recorder
}

// PromptString mocks base method
func (m *MockPrompter) PromptString(message string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PromptString", message)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PromptString indicates an expected call of PromptString
func (mr *MockPrompterMockRecorder) PromptString(message interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PromptString", reflect.TypeOf((*MockPrompter)(nil).PromptString), message)
}

// PromptPassword mocks base method
func (m *MockPrompter) PromptPassword(message string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PromptPassword", message)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PromptPassword indicates an expected call of PromptPassword
func (mr *MockPrompterMockRecorder) PromptPassword(message interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PromptPassword", reflect.TypeOf((*MockPrompter)(nil).PromptPassword), message)
}
