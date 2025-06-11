//go:build windows
// +build windows

package main

import (
	"log/slog"
	"net"
	"syscall"

	"golang.org/x/sys/windows"
)

// SetQuickAck sets TCP_NODELAY on Windows as an equivalent optimization.
// Windows doesn't have TCP_QUICKACK, but TCP_NODELAY provides similar latency benefits.
func SetQuickAck(conn net.Conn) error {
	// Use Go's built-in SetNoDelay method for cross-platform compatibility
	if tcpConn, ok := conn.(*net.TCPConn); ok {
		return tcpConn.SetNoDelay(true)
	}
	return nil
}

// SetLinger configures socket linger behavior on Windows.
func SetLinger(conn net.Conn) error {
	tcpConn, ok := conn.(*net.TCPConn)
	if !ok {
		return nil
	}

	rawConn, err := tcpConn.SyscallConn()
	if err != nil {
		return err
	}

	return rawConn.Control(func(fd uintptr) {
		linger := windows.Linger{
			Onoff:  1,
			Linger: 0,
		}
		err := windows.SetsockoptLinger(windows.Handle(fd), windows.SOL_SOCKET, windows.SO_LINGER, &linger)
		if err != nil {
			slog.Error("failed to set SO_LINGER", "error", err)
		}
	})
}

// GetTCPControlWithFastOpen provides Windows-specific TCP optimizations.
// Windows doesn't support TCP_FASTOPEN like Linux, but we can apply other optimizations.
func GetTCPControlWithFastOpen() func(network, address string, c syscall.RawConn) error {
	return func(network, _ string, c syscall.RawConn) error {
		return c.Control(func(fd uintptr) {
			handle := windows.Handle(fd)

			// Enable SO_REUSEADDR (Windows equivalent for some SO_REUSEPORT scenarios)
			// This allows binding to addresses that are in TIME_WAIT state
			err := windows.SetsockoptInt(handle, windows.SOL_SOCKET, windows.SO_REUSEADDR, 1)
			if err != nil {
				slog.Error("failed to set SO_REUSEADDR", "error", err)
			}

			// Enable TCP_NODELAY for better performance (disable Nagle's algorithm)
			err = windows.SetsockoptInt(handle, windows.IPPROTO_TCP, windows.TCP_NODELAY, 1)
			if err != nil {
				slog.Error("failed to set TCP_NODELAY", "error", err)
			}

			// Set send buffer size for better performance
			err = windows.SetsockoptInt(handle, windows.SOL_SOCKET, windows.SO_SNDBUF, 64*1024)
			if err != nil {
				slog.Error("failed to set SO_SNDBUF", "error", err)
			}

			// Set receive buffer size for better performance
			err = windows.SetsockoptInt(handle, windows.SOL_SOCKET, windows.SO_RCVBUF, 64*1024)
			if err != nil {
				slog.Error("failed to set SO_RCVBUF", "error", err)
			}

			// Disable socket lingering for faster connection cleanup
			linger := windows.Linger{
				Onoff:  1,
				Linger: 0,
			}
			err = windows.SetsockoptLinger(handle, windows.SOL_SOCKET, windows.SO_LINGER, &linger)
			if err != nil {
				slog.Error("failed to set SO_LINGER", "error", err)
			}
		})
	}
}
