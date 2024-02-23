package async

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/asmejia1993/payment-app/db/util"
	"github.com/hibiken/asynq"
	"github.com/sirupsen/logrus"
)

const (
	TypeNewTransaction = "trx:create"
	TypeRefund         = "trx:refund"
	TypeNewUser        = "user:new"
	TypeLoginUser      = "user:login"
)

type AsyncLog struct {
	asynq *asynq.Client
	log   *logrus.Logger
}

type AuditLogEntry struct {
	Actor   string      `json:"actor"`
	Action  string      `json:"action"`
	Module  string      `json:"module"`
	When    time.Time   `json:"when"`
	Details interface{} `json:"details"`
}

// NewAsynqClient creates a new asynchronous client using the provided configuration.
func NewAsynqClient(config util.Config, log *logrus.Logger) Queuer {
	return AsyncLog{
		asynq: asynq.NewClient(asynq.RedisClientOpt{Addr: config.RedisAddr}),
		log:   log,
	}
}

// Enqueue new event log.
func (a AsyncLog) Enqueue(payload AuditLogEntry, taskType string) error {
	req, err := json.Marshal(payload)
	if err != nil {
		msg := fmt.Errorf("failed to serialize payload: %v", err)
		a.log.Error(err)
		return msg
	}
	task := asynq.NewTask(taskType, req)
	info, err := a.asynq.Enqueue(task)
	if err != nil {
		msg := fmt.Errorf("could not enqueue task: %v", err)
		a.log.Error(msg)
		return msg
	}
	a.log.Infof("enqueued task: id=%s queue=%s", info.ID, info.Queue)
	return nil
}

// Close queues and redis connection
func (a AsyncLog) Close() {
	a.log.Info("closing asynq gracefully ...")
	a.asynq.Close()
}
