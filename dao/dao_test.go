package dao

import (
	"testing"
	"gotodo/proto"

	"github.com/sirupsen/logrus/hooks/test"

	"github.com/sirupsen/logrus"

	"github.com/stretchr/testify/assert"
)

func TestNewDao(t *testing.T) {
	assert.NotNil(t, NewDao(logrus.NewEntry(logrus.New())))
}

func TestCreate(t *testing.T) {
	testCases := []struct {
		testName string
		model    *model.Todo
	}{
		{
			testName: "nil argument",
			model:    nil,
		},
		{
			testName: "valid argument",
			model: &model.Todo{
				Text:      "text",
				Completed: false,
			},
		},
	}

	for _, tc := range testCases {
		// Setup
		mockLogger, hook := test.NewNullLogger()
		mockLogger.Level = logrus.TraceLevel
		dao := NewDao(mockLogger.WithField("test", "dao"))

		// Test
		result, err := dao.Create(tc.model)

		// Verify
		if tc.model == nil {
			// result should be nil
			assert.Nil(t, result, tc.testName)
			// err should not be nil
			assert.NotNil(t, err, tc.testName)
			// an error log message should have been created
			assert.Equal(t, logrus.ErrorLevel, hook.LastEntry().Level, tc.testName)
		} else {
			// result should not be nil
			assert.NotNil(t, result, tc.testName)
			// err should be nil
			assert.Nil(t, err, tc.testName)
			// the model's id should be set
			assert.Equal(t, "0", result.ID, tc.testName)
			// no warning or error logs should have been created
			for _, e := range hook.AllEntries() {
				if e.Level == logrus.WarnLevel || e.Level == logrus.ErrorLevel {
					assert.Fail(t, "No warning or error logs should have been created", tc.testName, e)
				}
			}
		}
	}
}
