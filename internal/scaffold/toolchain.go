package scaffold

import (
	"os/exec"
	"runtime"
	"strings"
)

// ToolchainGoVersion returns the version of the go command on PATH, such as "1.27.2".
// Generated projects target it so they build with the toolchain that generated them.
func ToolchainGoVersion() string {
	out, err := exec.Command("go", "env", "GOVERSION").Output()
	version := strings.TrimSpace(string(out))
	if err != nil || version == "" {
		version = runtime.Version()
	}
	version, _, _ = strings.Cut(strings.TrimPrefix(version, "go"), " ")
	return version
}
