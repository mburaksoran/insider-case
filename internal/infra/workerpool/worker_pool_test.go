package workerpool

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"testing"
	"time"
)

type mockTask struct {
	executeFn     func() error
	onFailureFn   func(err error)
	executeCalled bool
	onFailureErr  error
}

func (m *mockTask) Execute() error {
	m.executeCalled = true
	return m.executeFn()
}

func (m *mockTask) OnFailure(err error) {
	m.onFailureFn(err)
}

func TestWorkerPool(t *testing.T) {
	lgr := prepareLogger()
	t.Run("StartAndStop", func(t *testing.T) {
		pool, err := NewWorkerPool(2, 10, lgr)
		if err != nil {
			t.Fatalf("Error creating pool: %v", err)
		}

		pool.Start()
		pool.Stop()
	})

	t.Run("AddWorkAndExecute", func(t *testing.T) {
		pool, err := NewWorkerPool(2, 10, lgr)
		if err != nil {
			t.Fatalf("Error creating pool: %v", err)
		}

		pool.Start()
		defer pool.Stop()

		taskChan := make(chan struct{})
		taskExecuted := false

		mockTask := &mockTask{
			executeFn: func() error {
				taskExecuted = true
				close(taskChan)
				return nil
			},
			onFailureFn: func(err error) {},
		}

		pool.AddWork(mockTask)

		select {
		case <-taskChan:
		case <-time.After(2 * time.Second):
			t.Fatal("Task execution timeout")
		}

		if !taskExecuted {
			t.Fatal("Task was not executed")
		}

		if !mockTask.executeCalled {
			t.Fatal("Execute function was not called")
		}
	})

}

func prepareLogger() *zap.SugaredLogger {
	loggerConfig := zap.NewProductionConfig()
	loggerConfig.EncoderConfig.TimeKey = "timestamp"
	loggerConfig.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)

	logger, err := loggerConfig.Build()
	if err != nil {
		log.Fatal(err)
	}

	return logger.Sugar()
}
