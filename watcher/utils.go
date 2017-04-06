package watcher

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
)

func getEnvWithDefault(key, defaultVal string) string {
	var val string
	if val = os.Getenv(key); val == "" {
		val = defaultVal
	}
	return val
}

func removeContents(dir string) error {
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

func getPidFromPidfile(pidFile string) int {
	data, err := ioutil.ReadFile(pidFile)
	if err != nil {
		return 0
	}
	pid, _ := strconv.Atoi(string(data))
	return pid
}

func checkProcessAlive(pid int) error {
	var (
		process *os.Process
		err     error
	)
	if process, err = os.FindProcess(pid); err != nil {
		return err
	}
	return process.Signal(0)
}
