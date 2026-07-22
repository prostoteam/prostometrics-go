package prostometrics

import "runtime/debug"

// VersionString is returned for local builds without module version metadata.
const VersionString = "devel"

// Version returns the client version string.
func Version() string {
	if info, ok := debug.ReadBuildInfo(); ok {
		for _, dependency := range info.Deps {
			if dependency.Path == "github.com/prostoteam/prostometrics-go" && dependency.Version != "" {
				return dependency.Version
			}
		}
		if info.Main.Path == "github.com/prostoteam/prostometrics-go" && info.Main.Version != "" && info.Main.Version != "(devel)" {
			return info.Main.Version
		}
	}
	return VersionString
}
