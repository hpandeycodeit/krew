// Copyright 2019 The Kubernetes Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"os"

	"github.com/golang/glog"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"sigs.k8s.io/krew/pkg/constants"
	"sigs.k8s.io/krew/pkg/environment"
	"sigs.k8s.io/krew/pkg/version"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show krew version and diagnostics",
	Long: `Show version information and diagnostics about krew itself.

Remarks:
  - IsPlugin is true if krew is executed as a plugin
  - ExecutedVersion is the version of the currently executed binary. This is detected through the path.
  - GitTag describes the release name krew is built from.
  - GitCommit describes the git revision ID which krew is built from.
  - IndexURI is the URI where the index is updated from.
  - BasePath is the root directory for krew installation.
  - IndexPath is the directory that stores the local copy of the index git repository.
  - InstallPath is the directory for plugin installations.
  - DownloadPath is the directory for temporarily downloading plugins.
  - BinPath is the directory for the symbolic links to the installed plugin executables.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		selfPath, err := os.Executable()
		if err != nil {
			glog.Fatalf("failed to get the own executable path")
		}

		executedVersion, runningAsPlugin, err := environment.GetExecutedVersion(paths.InstallPath(), selfPath, environment.Realpath)
		if err != nil {
			return errors.Wrap(err, "failed to find current krew version")
		}

		conf := [][]string{
			{"IsPlugin", fmt.Sprintf("%v", runningAsPlugin)},
			{"ExecutedVersion", executedVersion},
			{"GitTag", version.GitTag()},
			{"GitCommit", version.GitCommit()},
			{"IndexURI", constants.IndexURI},
			{"BasePath", paths.BasePath()},
			{"IndexPath", paths.IndexPath()},
			{"InstallPath", paths.InstallPath()},
			{"DownloadPath", paths.DownloadPath()},
			{"BinPath", paths.BinPath()},
		}
		return printTable(os.Stdout, []string{"OPTION", "VALUE"}, conf)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
