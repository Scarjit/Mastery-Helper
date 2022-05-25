package helper

import (
	"golang.org/x/sys/windows"
	"syscall"
	"unicode/utf16"
	"unsafe"
)

// unsafe.Sizeof(windows.ProcessEntry32{})
const processEntrySize = 568

// GetProcesses returns a list of all processes.
func GetProcesses() []windows.ProcessEntry32 {

	h, e := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPPROCESS, 0)
	if e != nil {
		return []windows.ProcessEntry32{}
	}
	defer func(handle windows.Handle) {
		_ = windows.CloseHandle(handle)
	}(h)
	p := windows.ProcessEntry32{Size: processEntrySize}
	processes := make([]windows.ProcessEntry32, 0)
	for {
		e = windows.Process32Next(h, &p)
		if e != nil {
			break
		}
		processes = append(processes, p)
	}
	return processes
}

const ProcessVmRead = 0x00000010

// GetProcessPath returns the path of the executable
func GetProcessPath(pid uint32) (string, error) {
	handle, err := windows.OpenProcess(windows.PROCESS_QUERY_INFORMATION|ProcessVmRead, false, pid)
	if err != nil || handle == 0 {
		return "", err
	}

	var sysproc = syscall.MustLoadDLL("psapi.dll").MustFindProc("GetModuleFileNameExW")
	b := make([]uint16, syscall.MAX_PATH)
	r, _, err := sysproc.Call(uintptr(handle), 0, uintptr(unsafe.Pointer(&b[0])), uintptr(len(b)))
	n := uint32(r)
	if n == 0 {
		return "", err
	}
	name := string(utf16.Decode(b[0:n]))

	_ = windows.CloseHandle(handle)
	return name, err
}
