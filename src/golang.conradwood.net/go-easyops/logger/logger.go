package logger

import ()

type AsyncLogQueue struct{}

func NewAsyncLogQueue(appname string, buildid, repoid uint64, group, namespace, deplid string) (*AsyncLogQueue, error) {
	panic("this logger no longer exists as part of go-easyops. sorry.")
}

func (alq *AsyncLogQueue) String() string {
	panic("this logger no longer exists as part of go-easyops. sorry.")
}
func (alq *AsyncLogQueue) SetStatus(status string) {
	panic("this logger no longer exists as part of go-easyops. sorry.")
}
func (alq *AsyncLogQueue) LogCommandStdout(line string, status string) error {
	panic("this logger no longer exists as part of go-easyops. sorry.")
}
func (alq *AsyncLogQueue) Write(buf []byte) (int, error) {
	panic("this logger no longer exists as part of go-easyops. sorry.")
}
func (alq *AsyncLogQueue) Log(status string, format string, a ...interface{}) {
	panic("this logger no longer exists as part of go-easyops. sorry.")
}

func (alq *AsyncLogQueue) Close(exitcode int) error {
	panic("this logger no longer exists as part of go-easyops. sorry.")
}

func (alq *AsyncLogQueue) SetStartupID(s string) {
	panic("this logger no longer exists as part of go-easyops. sorry.")
}

func (alq *AsyncLogQueue) Flush() error {
	panic("this logger no longer exists as part of go-easyops. sorry.")
}
