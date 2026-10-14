/*
 * Copyright © 2019 – 2026 Red Hat Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

// Package config implements helper functions to manage toolbx containers persistent config
// parameters, such as injected environment variables
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/sirupsen/logrus"
	podenv "go.podman.io/podman/v6/pkg/env"
)

type Container interface {
	Env() map[string]string // Get container environment map
	ListEnv() []string      // Get serialized container environment
	Save() error            // Save updated container config
	Delete() error          // Delete config file
}

type containerConfig struct {
	Environment map[string]string `json:"env,omitempty"` // Environment var map

	fpath string // config location
	gid   int
	uid   int
}

func (cfg *containerConfig) Env() map[string]string {
	return cfg.Environment
}

func (cfg *containerConfig) ListEnv() []string {
	var str strings.Builder
	return slices.Collect(func(yield func(string) bool) {
		for k, v := range cfg.Environment {
			str.Grow(len(v) + 1) // key=val

			str.WriteString(k)
			str.WriteRune('=')
			str.WriteString(v)

			if !yield(str.String()) {
				return
			}
			str.Reset()
		}
	})
}

func (cfg *containerConfig) Save() error {
	f, err := os.OpenFile(cfg.fpath, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		logrus.Debugf("Save container config: error opening config file %s: %s", cfg.fpath, err)
		return fmt.Errorf("error opening file config file %s: %w", cfg.fpath, err)
	}
	defer f.Close()
	if cfg.uid != -1 && cfg.gid != -1 {
		f.Chown(cfg.uid, cfg.gid)
	}

	enc := json.NewEncoder(f)
	// 4-space pretty indentation
	enc.SetIndent("", "    ")

	return enc.Encode(cfg)
}

func (cfg *containerConfig) Delete() (err error) {
	err = os.Remove(cfg.fpath)
	if err != nil {
		logrus.Debugf("Delete container config: error deleting config %s: %s", cfg.fpath, err)
	}

	return
}

func GetContainerConfig(name string) (Container, error) {
	return loadContainerConfig(name, "", -1, -1)
}

func GetRunningContainerConfig(configDir string, uid, gid int) (Container, error) {
	env, err := parseContainerEnv()
	if err != nil {
		logrus.Debugf("Load container config: failed to load running container env: %s", err)
		return nil, fmt.Errorf("failed to parse running container env: %w", err)
	}
	return loadContainerConfig(env.Name, configDir, uid, gid)
}

func loadContainerConfig(name, configDir string, uid, gid int) (Container, error) {
	if configDir == "" {
		dir, err := os.UserConfigDir()
		if err != nil {
			logrus.Debugf("Load container config: failed to get the user config directory: %s", err)
			return nil, fmt.Errorf("failed to get the user config directory: %w", err)
		}
		configDir = dir
	}

	fpath := filepath.Join(configDir, "toolbox", name+"-config.json")
	cfg := containerConfig{fpath: fpath, uid: uid, gid: uid, Environment: make(map[string]string)}

	f, err := os.OpenFile(fpath, os.O_RDWR, 0644)
	if os.IsNotExist(err) {
		// Empty config in case file not exist (older container without config)
		logrus.Debugf("Load container config: container config file %s not found", fpath)
		return &cfg, nil
	} else if err != nil {
		logrus.Debugf("Load container config: failed to get the user config directory: %s", err)
		return nil, fmt.Errorf("error opening config %s: %w", cfg.fpath, err)
	}
	defer f.Close()

	err = json.NewDecoder(f).Decode(&cfg)
	if err != nil {
		logrus.Debugf("Load container config: failed decoding container config %s: %s", fpath, err)
		return nil, fmt.Errorf("error decoding config %s: %w", fpath, err)
	}

	return &cfg, nil
}

type podmanEnv struct {
	ImageID string
	Image   string
	Name    string
	ID      string
}

// parseContainerEnv parses podman's .containerenv file information
func parseContainerEnv() (*podmanEnv, error) {
	data, err := podenv.ParseFile("/run/.containerenv")
	if err != nil {
		return nil, err
	}

	return &podmanEnv{
		ImageID: strings.Trim(data["imageid"], "\""),
		Image:   strings.Trim(data["image"], "\""),
		Name:    strings.Trim(data["name"], "\""),
		ID:      strings.Trim(data["id"], "\""),
	}, nil
}
