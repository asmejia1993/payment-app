package async

type Queuer interface {
	Enqueue(payload AuditLogEntry, taskType string) error
	Close()
}
