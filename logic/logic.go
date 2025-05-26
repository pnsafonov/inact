package logic

const (
	/* DefaultLogFile = "/var/log/inact/inact.log" */

	Version = "1.0.0"
)

type Context struct {
	Days    int
	Mins    int
	Verbose bool
}

func DoRun(ctx *Context) error {

	return nil
}
