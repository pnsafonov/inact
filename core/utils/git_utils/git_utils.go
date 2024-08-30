package git_utils

import (
	_ "embed"
	"runtime/debug"
)

const (
	none = "none"
)

var (
	version string
	commit  string = none
)

func GetGitHash() string {
	if commit != none {
		return commit
	}

	if info, ok := debug.ReadBuildInfo(); ok {
		for _, setting := range info.Settings {
			if setting.Key == "vcs.revision" {
				return setting.Value
			}
		}
	}
	return ""
}

func GetGitHashShort() string {
	hash := GetGitHash()
	l0 := len(hash)
	if l0 >= 7 {
		return hash[:7]
	}
	return hash
}
