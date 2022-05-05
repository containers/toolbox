/*
 * Copyright © 2022 Ondřej Míchal
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

package cmd

import (
	"errors"
	"fmt"

	"github.com/containers/toolbox/pkg/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type ResolveFunc func([]string) error

type RequestType struct {
	ArgsN   int
	Resolve ResolveFunc
}

var (
	testFlags struct {
		requestType string
	}

	requestTypes = map[string]RequestType{
		"config-key": {
			1,
			resolveConfigKey,
		},
		"default-container-name": {
			0,
			resolveDefaultContainerName,
		},
		"default-image": {
			0,
			resolveDefaultImage,
		},
	}
)

var testCmd = &cobra.Command{
	Use:               "__test",
	Short:             "List existing toolbox containers and images",
	RunE:              test,
	Hidden:            true,
	ValidArgsFunction: completionEmpty,
}

func init() {
	flags := testCmd.Flags()

	flags.StringVar(&testFlags.requestType,
		"type",
		"",
		"type of data to be processed/retrieved")

	rootCmd.AddCommand(testCmd)
}

func test(cmd *cobra.Command, args []string) error {
	if !cmd.Flag("type").Changed {
		return errors.New("flag --type has to be provided")
	}

	requestType, ok := requestTypes[testFlags.requestType]
	if !ok {
		return fmt.Errorf("request type %s is not known", testFlags.requestType)
	}

	if requestType.ArgsN != len(args) {
		return fmt.Errorf("request type %s requires %d arguments", testFlags.requestType, requestType.ArgsN)
	}

	logrus.Debugf("Resolving request %s with arguments: %s", testFlags.requestType, args)

	resolve := requestType.Resolve
	if err := resolve(args); err != nil {
		return fmt.Errorf("failed to resolve request %s: %s", testFlags.requestType, err)
	}

	return nil
}

func resolveConfigKey(args []string) error {
	fmt.Print(viper.GetString(args[0]))
	return nil
}

func resolveDefaultContainerName(_ []string) error {
	containerName, _, _, err := utils.ResolveContainerAndImageNames("", "", "", "")
	if err != nil {
		return err
	}

	fmt.Print(containerName)
	return nil
}

func resolveDefaultImage(_ []string) error {
	_, image, release, err := utils.ResolveContainerAndImageNames("", "", "", "")
	if err != nil {
		return err
	}

	image, err = utils.GetFullyQualifiedImageFromDistros(image, release)
	if err != nil {
		return err
	}

	fmt.Print(image)
	return nil
}
