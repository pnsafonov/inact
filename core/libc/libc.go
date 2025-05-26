package libc

// #include <utmp.h>
import "C"
import (
	"fmt"
	"unsafe"
)

type ExitStatus struct {
	ETermination int32
	EExit        int32
}

type UTTVType struct {
	TVSec  int32
	TVUsec int32
}

type UTMP struct {
	UTType    uint16
	UTPid     int32
	UTLine    [32]byte
	UTId      [4]byte
	UTUser    [32]byte
	UTHost    [256]byte
	UTExit    ExitStatus
	UTSession int32
	UTTV      UTTVType
	UTAddrV6  [4]int32
	Reserved  [20]byte
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
		utmp2 := uintptr(unsafe.Pointer(utmp0))
		if utmp2 == 0 {
			break
		}

		utmp1 := (*UTMP)(unsafe.Pointer(utmp0))
		result = append(result, utmp1)

		i0 := int(utmp0.ut_type)
		fmt.Printf("i0 = %d\n", i0)
	}
	return result
}
