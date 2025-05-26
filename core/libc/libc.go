package libc

// #include <utmp.h>
import "C"
import (
	"fmt"
	"time"
	"unsafe"
)

const (
	// Empty - 0, No valid user accounting information
	Empty = uint16(C.EMPTY)

	// RunLvl - 1, The system's runlevel
	RunLvl = uint16(C.RUN_LVL)
	// BootTime - 2, Time of system boot
	BootTime = uint16(C.BOOT_TIME)
	// NewTime - 3, Time after system clock changed
	NewTime = uint16(C.NEW_TIME)
	// OldTime - 4, Time when system clock changed
	OldTime = uint16(C.OLD_TIME)

	// InitProcess - 5, Process spawned by the init process
	InitProcess = uint16(C.INIT_PROCESS)
	// LoginProcess - 6, Session leader of a logged in user
	LoginProcess = uint16(C.LOGIN_PROCESS)
	// UserProcess - 7, Normal process
	UserProcess = uint16(C.USER_PROCESS)
	// DeadProcess - 8, Terminated process
	DeadProcess = uint16(C.DEAD_PROCESS)
)

type ExitStatus struct {
	Termination int16
	Exit        int16
}

type UTMP struct {
	Type    uint16
	Pid     int32
	Line    string
	Id      string
	User    string
	Host    string
	Exit    ExitStatus
	Session int32
	Sec     int32
	USec    int32 // microseconds

	// Golang fields
	Time time.Time
}

func TypeToString(type0 uint16) string {
	switch type0 {
	case Empty:
		return "empty"
	case RunLvl:
		return "system"
	case BootTime:
		return "boot"
	case NewTime:
		return "new_time"
	case OldTime:
		return "old_time"
	case InitProcess:
		return "init"
	case LoginProcess:
		return "login"
	case UserProcess:
		return "user"
	case DeadProcess:
		return "dead"
	default:
		return fmt.Sprintf("%d", type0)
	}
}

// CGoString - конвертировать строку в go из массива фиксированной длинны.
// Отдаёт минимальную строку из двух вариантов: \0 терминированная строка,
// или ограниченная по длинне.
func CGoString(cstr *C.char, max C.int) string {
	str0 := C.GoStringN(cstr, max)
	l0 := len(str0)
	str1 := C.GoString(cstr)
	l1 := len(str1)

	if l0 < l1 {
		return str0
	}
	return str1
}

// Getutent0 - extern struct utmp *getutent (void) __THROW;
func Getutent0() []*UTMP {
	result := make([]*UTMP, 0)
	for {
		utmp0 := C.getutent()
		utmp1 := uintptr(unsafe.Pointer(utmp0))
		if utmp1 == 0 {
			break
		}

		line := CGoString(&utmp0.ut_line[0], C.UT_LINESIZE)
		id := CGoString(&utmp0.ut_id[0], 4)
		user := CGoString(&utmp0.ut_user[0], C.UT_NAMESIZE)
		host := CGoString(&utmp0.ut_host[0], C.UT_HOSTSIZE)

		exit := ExitStatus{
			Termination: int16(utmp0.ut_exit.e_termination),
			Exit:        int16(utmp0.ut_exit.e_exit),
		}
		utmp2 := &UTMP{
			Type:    uint16(utmp0.ut_type),
			Pid:     int32(utmp0.ut_pid),
			Line:    line,
			Id:      id,
			User:    user,
			Host:    host,
			Exit:    exit,
			Session: int32(utmp0.ut_session),
			Sec:     int32(utmp0.ut_tv.tv_sec),
			USec:    int32(utmp0.ut_tv.tv_usec),
		}

		time0 := time.Unix(int64(utmp2.Sec), int64(1000*utmp2.USec))
		utmp2.Time = time0

		result = append(result, utmp2)
	}
	return result
}
