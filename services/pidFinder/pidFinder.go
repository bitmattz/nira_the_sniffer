package services

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func FindPIDByPort(port int) ([]string, error) {
	targetHex := fmt.Sprintf("%04X", port)
	inodes := make(map[string]bool)

	files := []string{"/proc/net/tcp", "/proc/net/tcp6"}

	for _, file := range files {
		f, err := os.Open(file)
		if err != nil {
			continue
		}
		scanner := bufio.NewScanner(f)
		scanner.Scan() // skip header
		for scanner.Scan() {
			fields := strings.Fields(scanner.Text())
			localAddr := fields[1]
			inode := fields[9]

			parts := strings.Split(localAddr, ":")
			if len(parts) != 2 {
				continue
			}

			if parts[1] == targetHex {
				inodes[inode] = true
			}
		}
		f.Close()
	}

	var pids []string

	procDirs, _ := os.ReadDir("/proc")
	for _, dir := range procDirs {
		if !dir.IsDir() {
			continue
		}

		if _, err := strconv.Atoi(dir.Name()); err != nil {
			continue
		}

		fdPath := filepath.Join("/proc", dir.Name(), "fd")
		fds, err := os.ReadDir(fdPath)
		if err != nil {
			continue
		}

		for _, fd := range fds {
			link, err := os.Readlink(filepath.Join(fdPath, fd.Name()))
			if err != nil {
				continue
			}

			if strings.HasPrefix(link, "socket:[") {
				inode := strings.TrimSuffix(strings.TrimPrefix(link, "socket:["), "]")
				if inodes[inode] {
					pids = append(pids, dir.Name())
					break
				}
			}
		}
	}

	return pids, nil
}
