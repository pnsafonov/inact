package logic

import (
	"errors"
	"github.com/pnsafonov/inact/core/libc"
	"log"
	"os"
	"os/exec"
	"runtime"
	"time"
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

	pid := os.Getpid()
	log.Printf("DoRun, pid = %d, mins = %d, days = %d, verbose = %v, now = %v, before = %v\n", pid, ctx.Mins, ctx.Days, ctx.Verbose, now, before)

	count := 0
	for i := 0; i < l0; i++ {
		utp0 := utmps[i]

		if utp0.Type != libc.UserProcess && utp0.Type != libc.DeadProcess {
			printUTMP(ctx.Verbose, utp0, i, "0 by type")
			continue
		}

		if utp0.Time.Before(before) {
			printUTMP(ctx.Verbose, utp0, i, "0 by time")
			continue
		}

		count++
		printUTMP(ctx.Verbose, utp0, i, "1")
	}

	log.Printf("DoRun, recent logins count = %d\n", count)
	if count != 0 {
		log.Printf("DoRun, where is recent logins, no shutdown\n")
		return nil
	}

	var args []string
	// fix for shutdown not found
	_ = os.Setenv("PATH", "PATH=/sbin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/usr/local/sbin")

	if runtime.GOOS == "freebsd" {
		args = []string{"-p", "now"}
	} else { // linux
		args = []string{"-h", "now"}
	}
	cmd := exec.Command("shutdown", args...)

	//if runtime.GOOS == "freebsd" {
	//	args = []string{"-alh"}
	//} else { // linux
	//	args = []string{"-h"}
	//}
	//cmd := exec.Command("ls", args...)

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

func printUTMP(verbose bool, utp0 *libc.UTMP, i int, msg string) {
	if !verbose {
		return
	}
	log.Printf("i = %d, ut_type = %s, tv_sec = %v, ut_id = %s, ut_pid = %d, ut_user = %s, ut_line = %s, ut_host = %s, is_user = %s\n",
		i,
		libc.TypeToString(utp0.Type),
		utp0.Time,
		utp0.Id,
		utp0.Pid,
		utp0.User,
		utp0.Line,
		utp0.Host,
		msg,
	)
}
