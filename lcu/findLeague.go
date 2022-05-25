package lcu

import (
	"AramHelper/helper"
	"golang.org/x/sys/windows"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var defaultPath = filepath.FromSlash("C:\\Riot Games\\League of Legends\\lockfile")

// FindLeague finds the League of Legends installation path or retrieves the path from a running League of Legends process
func FindLeague() (path string, exists bool) {

	if helper.FileExists(defaultPath) {
		return defaultPath, true
	}

	return getLeagueFromProcesses()
}

// getLeagueFromProcesses retrieves the League of Legends installation path from a running League of Legends process
func getLeagueFromProcesses() (path string, open bool) {
	for _, processEntry32 := range helper.GetProcesses() {
		s := windows.UTF16ToString(processEntry32.ExeFile[:])
		if s == "LeagueClient.exe" {
			p, err := helper.GetProcessPath(processEntry32.ProcessID)
			if err != nil {
				return filepath.FromSlash(p), true
			}
		}
	}
	return "", false
}

// LockFile provides the password and port required to interact with the LCU
type LockFile struct {
	Process  string
	PID      uint64
	Port     uint64
	Password string
	Protocol string
}

// ReadLockFile reads the lockfile from the League of Legends installation path
func ReadLockFile() *LockFile {
	path, exist := FindLeague()
	if !exist {
		return nil
	}

	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	fileContent := string(bytes)
	splits := strings.Split(fileContent, ":")
	pid, _ := strconv.ParseUint(splits[1], 10, 64)
	port, _ := strconv.ParseUint(splits[2], 10, 64)
	return &LockFile{
		Process:  splits[0],
		PID:      pid,
		Port:     port,
		Password: splits[3],
		Protocol: splits[4],
	}
}
