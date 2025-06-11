//go:build windows
// +build windows

package main

// SetRLimitNoFile is a no-op on Windows as it doesn't have Unix-style rlimits.
// Windows handles file descriptors differently and doesn't require manual limit adjustment.
// Windows uses handles instead of file descriptors and manages them automatically.
func SetRLimitNoFile() error {
	// No-op on Windows - Windows doesn't use Unix-style rlimits
	// File handle limits are managed by the system and don't need manual adjustment
	return nil
}