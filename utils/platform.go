package utils

import "runtime"

func IsMac() bool {
	if runtime.GOOS == "darwin" {
		return true
	}

	return false
}

func IsWindows() bool {
	if runtime.GOOS == "windows" {
		return true
	}

	return false
}

func IsLinux() bool {
	if runtime.GOOS == "linux" {
		return true
	}

	return false
}
