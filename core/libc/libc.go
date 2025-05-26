package libc

// #include <utmp.h>
import "C"
import (
	"time"
	"unsafe"
)

//type ExitStatus struct {
//	ETermination int32
//	EExit        int32
//}
//
//type UTTVType struct {
//	TVSec  int32
//	TVUsec int32
//}
//
//type UTMP struct {
//	UTType    uint16
//	UTPid     int32
//	UTLine    [32]byte
//	UTId      [4]byte
//	UTUser    [32]byte
//	UTHost    [256]byte
//	UTExit    ExitStatus
//	UTSession int32
//	UTTV      UTTVType
//	UTAddrV6  [4]int32
//	Reserved  [20]byte
//}

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
	USec    int32

	// Golang fields
	Time time.Time
}

// Getutent - extern struct utmp *getutent (void) __THROW;
func Getutent0() []*UTMP {
	//var (
	//	result0 uintptr
	//)
	//result0 = C.getutent()
	//if result0 == 0 {
	//	return
	//}

	result := make([]*UTMP, 0)
	for {
		//utmp0 := uintptr(0)
		// _Cfunc_getutent)() (value of type *_Ctype_struct_utmp)
		utmp0 := C.getutent()
		utmp1 := uintptr(unsafe.Pointer(utmp0))
		if utmp1 == 0 {
			break
		}

		//utmp1 := (*UTMP)(unsafe.Pointer(utmp0))
		//result = append(result, utmp1)
		//
		//i0 := int(utmp0.ut_type)
		//fmt.Printf("i0 = %d\n", i0)

		//line := C.GoStringN(utmp0.ut_line, C.UT_LINESIZE)
		line := C.GoStringN(&utmp0.ut_line[0], C.UT_LINESIZE)
		id := C.GoStringN(&utmp0.ut_id[0], 4)
		user := C.GoStringN(&utmp0.ut_user[0], C.UT_NAMESIZE)
		host := C.GoStringN(&utmp0.ut_host[0], C.UT_HOSTSIZE)

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
