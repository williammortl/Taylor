package main

import (
	"github.com/williammortl/Taylor/src/starlightd/pkg/starlightdint"
)

//export TaylorGetPluginInfo
func TaylorGetPluginInfo() starlightdint.TaylorPluginInfo {
	return starlightdint.TaylorPluginInfo{
		Name:            "GPython",
		Version:         "0.0.3",
		Language:        "Python",
		LanguageVersion: "3.4.0",
	}
}

// export TaylorInit
func TaylorInit() {

}
