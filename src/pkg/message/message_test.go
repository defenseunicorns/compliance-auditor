package message_test

import (
	"bytes"
	"testing"

	"github.com/defenseunicorns/lula/src/pkg/message"
	"github.com/stretchr/testify/assert"
)

func TestUseBuffer(t *testing.T) {
	message.DisableColor()
	message.NoProgress = true
	message.SetLogLevel(message.DebugLevel)

	var buf bytes.Buffer
	message.UseBuffer(&buf)

	message.Info("info msg")
	message.Debug("debug msg")
	message.Warn("warn msg")
	message.Success("success msg")
	message.Detail("detail msg")
	message.Fail("fail msg")
	message.Note("note msg")
	message.Printf("printf msg")
	message.Question("question msg")

	bufOut := buf.String()
	assert.Contains(t, bufOut, "info msg")
	assert.Contains(t, bufOut, "debug msg")
	assert.Contains(t, bufOut, "warn msg")
	assert.Contains(t, bufOut, "success msg")
	assert.Contains(t, bufOut, "detail msg")
	assert.Contains(t, bufOut, "fail msg")
	assert.Contains(t, bufOut, "note msg")
	assert.Contains(t, bufOut, "printf msg")
	assert.Contains(t, bufOut, "question msg")
}
