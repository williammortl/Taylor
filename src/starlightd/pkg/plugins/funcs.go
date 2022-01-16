package plugins

import (
	"github.com/WilliamMortl/src/starlightd/pkg/starlightdint"
)

func AddPlugin(filename string) (*starlightdint.TaylorPluginInfo, error) {
	return nil, nil
}

func DeletePlugin(name string, version string) error {
	return nil
}

func DeleteAllPlugins() {

}

func GetPlugin(name string, version string) *starlightdint.TaylorPluginInfo {
	return nil
}
