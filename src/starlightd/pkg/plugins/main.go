package plugins

import (
	"plugin"
)

type pluginRecord struct {
	TaylorPluginInfo
	Filename string
	Handle   *plugin.Plugin
}

type pluginRecordsKey struct {
	Name    string
	Version string
}

var pluginRecords map[*pluginRecordsKey]*pluginRecord = make(map[*pluginRecordsKey]*pluginRecord)

var pluginRecordsByFilename map[string]*pluginRecordsKey = make(map[string]*pluginRecordsKey)

func AddPlugin(filename string) (*TaylorPluginInfo, error) {
	return nil, nil
}

func DeletePlugin(name string, version string) error {
	return nil
}

func DeleteAllPlugins() error {
	return nil
}

func GetPlugin(name string, version string) *TaylorPluginInfo {
	return nil
}
