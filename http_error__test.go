package httperror_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wanglihui/httperror"
)

func TestTryConstructHttpError(t *testing.T) {
	const (
		statusCode = 400
		msg        = "test error"
		code       = int64(1000)
	)
	err := httperror.New(statusCode, msg, code)
	if e, ok := httperror.TryConstructHttpError(err); !ok {
		t.FailNow()
	} else {
		assert.Equal(t, e.Code, code)
		assert.Equal(t, e.Message, msg)
		assert.Equal(t, e.StatusCode, statusCode)
	}
}

func TestTryConstructHttpErrorWithJson(t *testing.T) {
	const (
		statusCode = 400
		msg        = "test error"
		code       = int64(1000)
	)
	message := fmt.Sprintf(`{"msg": "%s", "code": %d, "http_code": %d}`, msg, code, statusCode)
	if e, ok := httperror.TryConstructHttpError(message); !ok {
		t.FailNow()
	} else {
		assert.Equal(t, e.Code, code)
		assert.Equal(t, e.Message, msg)
		assert.Equal(t, e.StatusCode, statusCode)
	}
}

func TestTryConstructHttpErrorWithString(t *testing.T) {
	const (
		msg        = "test error"
		code       = int64(1000)
		statusCode = 400
	)
	if e, ok := httperror.TryConstructHttpError(msg); ok {
		t.FailNow()
	} else {
		// assert.Equal(t, e.Code, code)
		assert.Empty(t, e)
		// assert.Equal(t, e.Message, msg)
		// assert.Equal(t, e.StatusCode, statusCode)
	}
}
