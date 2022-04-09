package plugins

type TaylorPluginInfo struct {
	Name            string
	Version         string
	Language        string
	LanguageVersion string
}

type TaylorVariable struct {
	Type string
	Data []byte
}

type TaylorInterop interface{}

/*
	ThreadSpawn()
	VariableLoad(name string) []byte
	MutexLock(name string)
	MutexUnloack(name string)
}
*/

// Note: InitPlugin is *spiritually* a member of this interface, but isn't really.
// 	It is the method that bootstraps the connection between starlightd and the language
// 	plugin
type TaylorPluginInterop interface {
	// InitPlugin(taylor *TaylorInterop) *TaylorPluginInfo, *TaylorPluginInterop
	RunFile(filename string, args string) (string, error)
	ResumeRun(filename string, offset uint64) error
}
