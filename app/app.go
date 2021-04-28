package app

var Name string

func Init(name, version, commit, buildDate, mode string) {
	Name = name

	IsDebugMode = mode == DevMode
	IsProdMode = mode == ProdMode
	IsTestDevMode = mode == TestDevMode

	Version = newVersion(version, commit, buildDate)
}
