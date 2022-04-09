package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"plugin"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/williammortl/Taylor/src/starlightd/pkg/plugins"
	telemetryLib "github.com/williammortl/Taylor/src/starlightd/pkg/telemetry"
	"github.com/williammortl/Taylor/src/starlightd/pkg/util"
)

const componentName string = "launcher handler"

type taylorLauncherPost struct {
	Bytecode        string `json:"bytecode"`
	BytecodeVersion string `json:"bytecodeVersion"`
	FilenameStart   string `json:"filenameStart"`
	Args            string `json:"args,omitempty"`
	RunInfo         *taylorLauncherGet
}

type taylorLauncherPostReturn struct {
	RunID string `json:"runID"`
}

type taylorLauncherGet struct {
	Status   string `json:"status"`
	ExitCode int    `json:"exitCode,omitempty"`
	Terminal string `json:"terminal,omitempty"`
}

var (
	currentRuns map[string]taylorLauncherPost = make(map[string]taylorLauncherPost)
)

func LauncherHandlerPost(telemetry *telemetryLib.Telemetry) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		// read JSON body
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			errorTitle := "Missing JSON"
			errorMsg := "JSON was not posted posted to this endpoint!"
			http.Error(w, errorMsg, http.StatusBadRequest)
			(*telemetry).LogError(componentName, errorTitle, fmt.Errorf(errorMsg+" IP: ", util.GetIP(r)))
			return
		}

		// get JSON from the body
		var post taylorLauncherPost
		err = json.Unmarshal(body, &post)
		if err != nil {
			errorTitle := "Malformed JSON"
			errorMsg := "Malformed JSON was posted posted to this endpoint!"
			http.Error(w, errorMsg, http.StatusBadRequest)
			(*telemetry).LogError(componentName, errorTitle, fmt.Errorf(errorMsg+" IP: ", util.GetIP(r)))
			return
		}

		// TODO: split based on bytecode and version
		post.RunInfo = &taylorLauncherGet{
			Status:   "In Progress",
			ExitCode: -1,
			Terminal: "",
		}
		postResponse := taylorLauncherPostReturn{
			RunID: uuid.New().String(),
		}
		currentRuns[postResponse.RunID] = post
		go goLaunchPython(postResponse.RunID, &post, telemetry)

		// respond
		responseJson, _ := json.Marshal(postResponse)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, string(responseJson))
		(*telemetry).LogInfo(componentName, "Started",
			fmt.Sprintf("Started run: %v", postResponse.RunID))
	}
}

func goLaunchPython(runID string, post *taylorLauncherPost, telemetry *telemetryLib.Telemetry) {

	// load gpython plugin
	p, err := plugin.Open("gpython.so")
	if err != nil {
		panic(err)
	}

	initPlugin, err := p.Lookup("InitPlugin")
	if err != nil {
		panic(err)
	}

	_, pluginInterop := initPlugin.(func(plugins.TaylorInterop) (*plugins.TaylorPluginInfo, plugins.TaylorPluginInterop))(nil)

	// run bytecode
	post.RunInfo.Terminal, err = pluginInterop.RunFile(post.FilenameStart, post.Args)
	post.RunInfo.ExitCode = 0
	if err != nil {
		post.RunInfo.ExitCode = 1
	}
	post.RunInfo.Status = "Completed"
}

func LauncherHandlerGet(telemetry *telemetryLib.Telemetry) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		// get {id} from the REST path
		runID := mux.Vars(r)["runID"]
		if runID == "" {
			errorTitle := "Missing runID"
			errorMsg := "The runID is missing!"
			http.Error(w, errorMsg, http.StatusBadRequest)
			(*telemetry).LogError(componentName, errorTitle, fmt.Errorf(errorMsg+" IP: ", util.GetIP(r)))
			return
		}

		// check for runID in the running jobs
		post, exist := currentRuns[runID]
		if !exist {
			errorTitle := "Invalid runID"
			errorMsg := "The runID is not valid!"
			http.Error(w, errorMsg, http.StatusBadRequest)
			(*telemetry).LogError(componentName, errorTitle, fmt.Errorf(errorMsg+" IP: ", util.GetIP(r)))
			return
		}

		// respond
		responseJson, _ := json.Marshal(post.RunInfo)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, string(responseJson))
		(*telemetry).LogInfo(componentName, "Qeuried", runID)
	}
}
