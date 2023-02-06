/*
Copyright 2022 The Kubernetes Authors All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"context"
	"time"

	"golang.org/x/mod/semver"
	"k8s.io/klog/v2"

	"k8s.io/minikube/hack/update"
)

const (
	// default context timeout
	cxTimeout = 5 * time.Minute
)

var (
	schema = map[string]update.Item{
		"netlify.toml": {
			Replace: map[string]string{
				`HUGO_VERSION = .*`: `HUGO_VERSION = "{{.StableVersion}}"`,
			},
		},
	}
)

// Data holds stable Hugo version in semver format.
type Data struct {
	StableVersion string
}

func main() {
	// set a context with defined timeout
	ctx, cancel := context.WithTimeout(context.Background(), cxTimeout)
	defer cancel()

	// get Hugo stable version
	stable, err := hugoVersion(ctx, "gohugoio", "hugo")
	if err != nil {
		klog.Fatalf("Unable to get Hugo stable version: %v", err)
	}
	data := Data{StableVersion: stable}
	klog.Infof("Hugo stable version: %s", stable)

	update.Apply(schema, data)
}

// hugoVersion returns stable version in semver format.
func hugoVersion(ctx context.Context, owner, repo string) (string, error) {
	// get Hugo version from GitHub Releases
	stable, _, _, err := update.GHReleases(ctx, owner, repo)
	if err != nil || !semver.IsValid(stable.Tag) {
		return "", err
	}
	return stable.Tag, nil
}
