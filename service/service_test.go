package service

import (
	"testing"
	"gotodo/proto"

	"github.com/pkg/errors"

	d "todo/dao"

	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"

	"github.com/stretchr/testify/assert"
)

func TestNewService(t *testing.T) {
	assert.NotNil(t, NewService(logrus.NewEntry(logrus.New()), d.MockDao{}))
}

func TestDelete(t *testing.T) {
	testCases := []struct {
		testName string
		model    *model.Task
		hasError bool
	}{
		{
			testName: "nil argument",
			model:    nil,
		},
		{
			testName: "valid argument",
			model:    &model.Task{},
		},
		{
			testName: "has error",
			model:    &model.Task{},
			hasError: true,
		},
	}

	for _, tc := range testCases {
		// Setup
		mockLogger, _ := test.NewNullLogger()
		mockLogger.Level = logrus.TraceLevel

		mockDao := d.MockDao{}
		if tc.hasError {
			mockDao.On("Delete", tc.model).Return(errors.New("error"))
		} else {
			mockDao.On("Delete", tc.model).Return(nil)
		}

		service := NewService(mockLogger.WithField("test", "service"), mockDao)

		// Test
		err := service.Delete(tc.model)

		// Verify
		if tc.hasError {
			// err should not be nil
			assert.NotNil(t, err, tc.testName)
		} else {
			// err should be nil
			assert.Nil(t, err, tc.testName)
		}
	}
}
