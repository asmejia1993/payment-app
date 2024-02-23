package async

import (
	"testing"
	"time"

	"github.com/asmejia1993/payment-app/db/util"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockQueuer struct {
	mock.Mock
}

func (m *MockQueuer) Enqueue(payload AuditLogEntry, taskType string) error {
	args := m.Called(payload, taskType)
	return args.Error(0)
}

func (m *MockQueuer) Close() {
	m.Called()
}

func TestEnqueueSerializeError(t *testing.T) {
	payload := AuditLogEntry{
		Actor:   "user1",
		Action:  "create",
		Module:  "transaction",
		When:    time.Now(),
		Details: "details",
	}
	config := util.Config{
		RedisAddr: "localhost:6379",
	}
	log := *logrus.New()
	client := NewAsynqClient(config, &log)

	err := client.Enqueue(payload, TypeNewTransaction)
	assert.NoError(t, err)
}
