package logic

import (
	"fmt"
	libc0 "modernc.org/libc"
	"unsafe"
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

	tls := libc0.NewTLS()
	defer tls.Close()

	for i := 0; true; i++ {
		utp0 := libc0.Xgetutent(tls)
		if utp0 == 0 {
			break
		}
		utmp1 := unsafe.Pointer(utp0)
		utmp2 := (*libc0.Tutmpx)(utmp1)

		fmt.Printf("i = %d, ut_type = %d, tv_sec = %d, ut_id = %s, ut_pid = %d, ut_user = %s, ut_line = %s\n",
			i,
			utmp2.Fut_type,
			utmp2.Fut_tv.Ftv_sec,
			utmp2.Fut_id,
			utmp2.Fut_pid,
			utmp2.Fut_user,
			utmp2.Fut_line,
		)

		//utmp2.
	}

	return nil
}
