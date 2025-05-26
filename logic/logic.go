package logic

import (
	"errors"
	"github.com/pnsafonov/inact/core/libc"
	"log"
	"os/exec"
	"runtime"
	"time"
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
	now := time.Now()
	before := now.Add(-1 * 24 * time.Hour * 7) // 7 days before
	if ctx.Mins > 0 {
		before = now.Add(-1 * time.Minute * time.Duration(ctx.Mins))
	}
	if ctx.Mins == 0 && ctx.Days > 0 {
		before = now.Add(-1 * 24 * time.Hour * time.Duration(ctx.Days))
	}

	utmps := libc.Getutent0()

	l0 := len(utmps)
	count := 0
	for i := 0; i < l0; i++ {
		utp0 := utmps[i]

		if ctx.Verbose {
			printUTMP(utp0, i)
		}

		if utp0.Type != libc.UserProcess && utp0.Type != libc.DeadProcess {
			continue
		}
		if utp0.Time.Before(before) {
			continue
		}
		count++
	}

	log.Printf("DoRun, recent logins count = %d\n", count)
	if count != 0 {
		log.Printf("DoRun, where is recent logins, no shutdown\n")
		return nil
	}

	var args []string
	//if runtime.GOOS == "freebsd" {
	//	args = []string{"-p", "now"}
	//} else { // linux
	//	args = []string{"-h", "now"}
	//}
	//cmd := exec.Command("shutdown", args...)

	if runtime.GOOS == "freebsd" {
		args = []string{"-alh"}
	} else { // linux
		args = []string{"-h"}
	}
	cmd := exec.Command("ls", args...)

	log.Printf("DoRun, do shutdown\n")
	if errors.Is(cmd.Err, exec.ErrDot) {
		cmd.Err = nil
	}
	err := cmd.Run()
	if err != nil {
		log.Printf("DoRun, shutdown err = %v\n", err)
		return err
	}

	return nil
}

func printUTMP(utp0 *libc.UTMP, i int) {
	log.Printf("i = %d, ut_type = %s, tv_sec = %v, ut_id = %s, ut_pid = %d, ut_user = %s, ut_line = %s, ut_host = %s\n",
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
