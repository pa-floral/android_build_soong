// Copyright 2015 Google Inc. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package android

import (
	"os"
	"strings"

	"android/soong/shared"
)

// This file supports dependencies on environment variables.  During build manifest generation,
// any dependency on an environment variable is added to a list.  During the singleton phase
// a JSON file is written containing the current value of all used environment variables.
// The next time the top-level build script is run, it uses the soong_env executable to
// compare the contents of the environment variables, rewriting the file if necessary to cause
// a manifest regeneration.

var originalEnv map[string]string
var SdclangEnv map[string]string
var QiifaBuildConfig string


func InitEnvironment(envFile string) {
	var err error
	originalEnv, err = shared.EnvFromFile(envFile)
	if err != nil {
		panic(err)
	}
	// TODO(b/182662723) move functionality as necessary
	SdclangEnv = make(map[string]string)
	for _, env := range os.Environ() {
		idx := strings.IndexRune(env, '=')
		if idx != -1 {
			SdclangEnv[env[:idx]] = env[idx+1:]
		}
	}
        QiifaBuildConfig = originalEnv["QIIFA_BUILD_CONFIG"]
}

func EnvSingleton() Singleton {
	return &envSingleton{}
}

type envSingleton struct{}

func (c *envSingleton) GenerateBuildActions(ctx SingletonContext) {
	envDeps := ctx.Config().EnvDeps()

	envFile := PathForOutput(ctx, "soong.environment.used")
	if ctx.Failed() {
		return
	}

	data, err := shared.EnvFileContents(envDeps)
	if err != nil {
		ctx.Errorf(err.Error())
	}

	err = WriteFileToOutputDir(envFile, data, 0666)
	if err != nil {
		ctx.Errorf(err.Error())
	}

	ctx.AddNinjaFileDeps(envFile.String())
}
