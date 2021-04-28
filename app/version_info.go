package app

var Version *versionInfo

type versionInfo struct {
	Version   string
	Commit    string
	BuildDate string
}

func newVersion(version, commit, buildDate string) *versionInfo {
	return &versionInfo{
		BuildDate: buildDate,
		Commit:    commit,
		Version:   version,
	}

}
