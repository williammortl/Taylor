package plugins

import (
	"plugin"

	"github.com/WilliamMortl/src/starlightd/pkg/starlightdint"
)

type pluginRecord struct {
	starlightdint.TaylorPluginInfo
	Filename string
	Handle   *plugin.Plugin
}
