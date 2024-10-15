package firehydrant

import (
	"fmt"
	"runtime/debug"
)

// Version value gets defined at build-time, stamped by GoReleaser.
var Version = "dev"

type BuildInfo struct {
	Version string
	Commit  string
}

func GetBuildInfo() *BuildInfo {
	result := &BuildInfo{
		Version: Version,
		Commit:  "unknown",
	}

	bi, ok := debug.ReadBuildInfo()
	if !ok {
		return result
	}

	for _, i := range bi.Settings {
		if i.Key == "vcs.revision" {
			result.Commit = i.Value
			break
		}
	}
	return result
}

func (b *BuildInfo) String() string {
	return fmt.Sprintf("%s@%s", b.Version, b.Commit)
}
