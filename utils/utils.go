package utils

import (
	"os"
	"path/filepath"
	"syscall"
)

func GetEnvWithDefault(key, defaultVal string) string {
	var val string
	if val = os.Getenv(key); val == "" {
		val = defaultVal
	}
	return val
}

func RemoveContents(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return nil
}

func CheckProcessAlive(pid int) error {
	var (
		process *os.Process
		err     error
	)
	if process, err = os.FindProcess(pid); err != nil {
		return err
	}
	return process.Signal(syscall.Signal(0x0))
}
