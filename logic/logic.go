package logic

import (
	"fmt"
	"github.com/pnsafonov/inact/core/libc"
)

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

	utmps := libc.Getutent0()

	l0 := len(utmps)
	for i := 0; i < l0; i++ {
		utp0 := utmps[i]

		if ctx.Verbose {
			printUTMP(utp0, i)
		}

	}

	return nil
}

func printUTMP(utp0 *libc.UTMP, i int) {
	fmt.Printf("i = %d, ut_type = %s, tv_sec = %v, ut_id = %s, ut_pid = %d, ut_user = %s, ut_line = %s, ut_host = %s\n",
		i,
		libc.TypeToString(utp0.Type),
		utp0.Time,
		utp0.Id,
		utp0.Pid,
		utp0.User,
		utp0.Line,
		utp0.Host,
	)
}
