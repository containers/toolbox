/*
 * Copyright © 2019 – 2025 Red Hat Inc.
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
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/containers/toolbox/pkg/podman"
	"github.com/containers/toolbox/pkg/shell"
	"github.com/containers/toolbox/pkg/skopeo"
	"github.com/containers/toolbox/pkg/term"
	"github.com/containers/toolbox/pkg/utils"
	"github.com/docker/go-units"
	"github.com/godbus/dbus/v5"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type promptForDownloadError struct {
	ImageSize string
}

const (
	alpha    = `abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ`
	num      = `0123456789`
	alphanum = alpha + num
)

var (
	createFlags struct {
		authFile  string
		container string
		distro    string
		image     string
		release   string
	}

	createToolboxShMounts = []struct {
		containerPath string
		source        string
	}{
		{"/etc/profile.d/toolbox.sh", "/etc/profile.d/toolbox.sh"},
		{"/etc/profile.d/toolbox.sh", "/usr/share/profile.d/toolbox.sh"},
	}
)

var createCmd = &cobra.Command{
	Use:               "create",
	Short:             "Create a new Toolbx container",
	RunE:              create,
	ValidArgsFunction: completionEmpty,
}

func init() {
	flags := createCmd.Flags()

	flags.StringVar(&createFlags.authFile,
		"authfile",
		"",
		"Path to a file with credentials for authenticating to the registry for private images")

	flags.StringVarP(&createFlags.container,
		"container",
		"c",
		"",
		"Assign a different name to the Toolbx container")

	flags.StringVarP(&createFlags.distro,
		"distro",
		"d",
		"",
		"Create a Toolbx container for a different operating system distribution than the host")

	flags.StringVarP(&createFlags.image,
		"image",
		"i",
		"",
		"Change the name of the base image used to create the Toolbx container")

	flags.StringVarP(&createFlags.release,
		"release",
		"r",
		"",
		"Create a Toolbx container for a different operating system release than the host")

	createCmd.SetHelpFunc(createHelp)

	if err := createCmd.RegisterFlagCompletionFunc("distro", completionDistroNames); err != nil {
		panicMsg := fmt.Sprintf("failed to register flag completion function: %v", err)
		panic(panicMsg)
	}

	if err := createCmd.RegisterFlagCompletionFunc("image", completionImageNames); err != nil {
		panicMsg := fmt.Sprintf("failed to register flag completion function: %v", err)
		panic(panicMsg)
	}

	rootCmd.AddCommand(createCmd)
}

func create(cmd *cobra.Command, args []string) error {
	if utils.IsInsideContainer() {
		if !utils.IsInsideToolboxContainer() {
			return errors.New("this is not a Toolbx container")
		}

		exitCode, err := utils.ForwardToHost()
		return &exitError{exitCode, err}
	}

	if cmd.Flag("distro").Changed && cmd.Flag("image").Changed {
		var builder strings.Builder
		fmt.Fprintf(&builder, "options --distro and --image cannot be used together\n")
		fmt.Fprintf(&builder, "Run '%s --help' for usage.", executableBase)

		errMsg := builder.String()
		return errors.New(errMsg)
	}

	if cmd.Flag("image").Changed && cmd.Flag("release").Changed {
		var builder strings.Builder
		fmt.Fprintf(&builder, "options --image and --release cannot be used together\n")
		fmt.Fprintf(&builder, "Run '%s --help' for usage.", executableBase)

		errMsg := builder.String()
		return errors.New(errMsg)
	}

	if cmd.Flag("authfile").Changed {
		if !utils.PathExists(createFlags.authFile) {
			var builder strings.Builder
			fmt.Fprintf(&builder, "file %s not found\n", createFlags.authFile)
			fmt.Fprintf(&builder, "'podman login' can be used to create the file.\n")
			fmt.Fprintf(&builder, "Run '%s --help' for usage.", executableBase)

			errMsg := builder.String()
			return errors.New(errMsg)
		}
	}

	var container string
	var containerArg string

	if len(args) != 0 {
		container = args[0]
		containerArg = "CONTAINER"
	} else if createFlags.container != "" {
		container = createFlags.container
		containerArg = "--container"
	}

	container, image, release, err := resolveContainerAndImageNames(container,
		containerArg,
		createFlags.distro,
		createFlags.image,
		createFlags.release)

	if err != nil {
		return err
	}

	if err := createContainer(container, image, release, createFlags.authFile, true); err != nil {
		return err
	}

	return nil
}

func createContainer(container, image, release, authFile string, showCommandToEnter bool) error {
	if container == "" {
		panic("container not specified")
	}

	if image == "" {
		panic("image not specified")
	}

	if release == "" {
		panic("release not specified")
	}

	enterCommand := getEnterCommand(container)

	logrus.Debugf("Checking if container %s already exists", container)

	if exists, _ := podman.ContainerExists(container); exists {
		var builder strings.Builder
		fmt.Fprintf(&builder, "container %s already exists\n", container)
		fmt.Fprintf(&builder, "Enter with: %s\n", enterCommand)
		fmt.Fprintf(&builder, "Run '%s --help' for usage.", executableBase)

		errMsg := builder.String()
		return errors.New(errMsg)
	}

	pulled, err := pullImage(image, release, authFile)
	if err != nil {
		return err
	}
	if !pulled {
		return nil
	}

	imageFull, err := podman.GetFullyQualifiedImageFromRepoTags(image)
	if err != nil {
		var errImage *podman.ImageError

		if errors.As(err, &errImage) {
			if errors.Is(err, podman.ErrImageRepoTagsEmpty) {
				logrus.Debugf("Image %s has empty RepoTags, likely because it is without a name", image)
				imageFull = image
			} else if errors.Is(err, podman.ErrImageRepoTagsMissing) {
				return fmt.Errorf("missing RepoTags for image %s", image)
			} else {
				panicMsg := fmt.Sprintf("unexpected %T: %s", err, err)
				panic(panicMsg)
			}
		} else {
			return err
		}
	}

	if !rootFlags.assumeYes {
		if isToolboxImage, err := podman.IsToolboxImage(imageFull); err != nil {
			return fmt.Errorf("failed to verify image compatibility: %w", err)
		} else if !isToolboxImage {
			prompt := fmt.Sprintf("Image '%s' is not a Toolbx image and may not work properly (see https://containertoolbx.org/doc/). Continue anyway? [y/N]:", imageFull)
			if !askForConfirmation(prompt) {
				return nil
			}
		}
	}

	var toolbxDelayEntryPointEnv []string

	if toolbxDelayEntryPoint, ok := os.LookupEnv("TOOLBX_DELAY_ENTRY_POINT"); ok {
		toolbxDelayEntryPointEnvArg := "TOOLBX_DELAY_ENTRY_POINT=" + toolbxDelayEntryPoint
		toolbxDelayEntryPointEnv = []string{"--env", toolbxDelayEntryPointEnvArg}
	}

	var toolbxFailEntryPointEnv []string

	if toolbxFailEntryPoint, ok := os.LookupEnv("TOOLBX_FAIL_ENTRY_POINT"); ok {
		toolbxFailEntryPointEnvArg := "TOOLBX_FAIL_ENTRY_POINT=" + toolbxFailEntryPoint
		toolbxFailEntryPointEnv = []string{"--env", toolbxFailEntryPointEnvArg}
	}

	toolboxPath := os.Getenv("TOOLBOX_PATH")
	toolboxPathEnvArg := "TOOLBOX_PATH=" + toolboxPath
	toolboxPathMountArg := toolboxPath + ":/usr/bin/toolbox:ro"

	var runtimeDirectory string
	var xdgRuntimeDirEnv []string

	if currentUser.Uid == "0" {
		runtimeDirectory, err = utils.GetRuntimeDirectory(currentUser)
		if err != nil {
			return err
		}
	} else {
		xdgRuntimeDir := os.Getenv("XDG_RUNTIME_DIR")
		xdgRuntimeDirEnvArg := "XDG_RUNTIME_DIR=" + xdgRuntimeDir
		xdgRuntimeDirEnv = []string{"--env", xdgRuntimeDirEnvArg}

		runtimeDirectory = xdgRuntimeDir
	}

	runtimeDirectoryMountArg := runtimeDirectory + ":" + runtimeDirectory

	logrus.Debug("Checking if 'podman create' supports '--mount type=devpts'")

	var devPtsMount []string

	if podman.CheckVersion("2.1.0") {
		logrus.Debug("'podman create' supports '--mount type=devpts'")
		devPtsMount = []string{"--mount", "type=devpts,destination=/dev/pts"}
	}

	var usernsArg string
	if currentUser.Uid == "0" {
		usernsArg = "host"
	} else {
		usernsArg = "keep-id"
	}

	dbusSystemSocket, err := getDBusSystemSocket()
	if err != nil {
		return err
	}

	dbusSystemSocketMountArg := dbusSystemSocket + ":" + dbusSystemSocket

	currentUserHomeDir := getCurrentUserHomeDir()
	homeDirEvaled, err := filepath.EvalSymlinks(currentUserHomeDir)
	if err != nil {
		return fmt.Errorf("failed to canonicalize %s", currentUserHomeDir)
	}

	logrus.Debugf("%s canonicalized to %s", currentUserHomeDir, homeDirEvaled)
	homeDirMountArg := homeDirEvaled + ":" + homeDirEvaled + ":rslave"

	var avahiSocketMount []string

	avahiSocket, err := getServiceSocket("Avahi", "avahi-daemon.socket")
	if err != nil {
		logrus.Debug(err)
	}
	if avahiSocket != "" {
		avahiSocketMountArg := avahiSocket + ":" + avahiSocket
		avahiSocketMount = []string{"--volume", avahiSocketMountArg}
	}

	var kcmSocketMount []string

	kcmSocket, err := getServiceSocket("KCM", "sssd-kcm.socket")
	if err != nil {
		logrus.Debug(err)
	}
	if kcmSocket != "" {
		kcmSocketMountArg := kcmSocket + ":" + kcmSocket
		kcmSocketMount = []string{"--volume", kcmSocketMountArg}
	}

	var pcscSocketMount []string

	pcscSocket, err := getServiceSocket("pcsc", "pcscd.socket")
	if err != nil {
		logrus.Debug(err)
	}
	if pcscSocket != "" {
		pcscSocketMountArg := pcscSocket + ":" + pcscSocket
		pcscSocketMount = []string{"--volume", pcscSocketMountArg}
	}

	var mediaLink []string
	var mediaMount []string

	if utils.PathExists("/media") {
		logrus.Debug("Checking if /media is a symbolic link to /run/media")

		mediaPath, _ := filepath.EvalSymlinks("/media")
		if mediaPath == "/run/media" {
			logrus.Debug("/media is a symbolic link to /run/media")
			mediaLink = []string{"--media-link"}
		} else {
			mediaMount = []string{"--volume", "/media:/media:rslave"}
		}
	}

	var mntLink []string
	var mntMount []string

	if utils.PathExists("/mnt") {
		logrus.Debug("Checking if /mnt is a symbolic link to /var/mnt")

		mntPath, _ := filepath.EvalSymlinks("/mnt")
		if mntPath == "/var/mnt" {
			logrus.Debug("/mnt is a symbolic link to /var/mnt")
			mntLink = []string{"--mnt-link"}
		} else {
			mntMount = []string{"--volume", "/mnt:/mnt:rslave"}
		}
	}

	var runMediaMount []string

	if utils.PathExists("/run/media") {
		runMediaMount = []string{"--volume", "/run/media:/run/media:rslave"}
	}

	logrus.Debug("Looking up toolbox.sh")

	var toolboxShMount []string

	for _, mount := range createToolboxShMounts {
		if utils.PathExists(mount.source) {
			logrus.Debugf("Found %s", mount.source)

			toolboxShMountArg := mount.source + ":" + mount.containerPath + ":ro"
			toolboxShMount = []string{"--volume", toolboxShMountArg}
			break
		}
	}

	logrus.Debug("Checking if /home is a symbolic link to /var/home")

	var slashHomeLink []string

	slashHomeEvaled, _ := filepath.EvalSymlinks("/home")
	if slashHomeEvaled == "/var/home" {
		logrus.Debug("/home is a symbolic link to /var/home")
		slashHomeLink = []string{"--home-link"}
	}

	logLevelString := podman.LogLevel.String()

	userShell := os.Getenv("SHELL")
	if userShell == "" {
		return errors.New("failed to get the current user's default shell")
	}

	entryPoint := []string{
		"toolbox", "--log-level", "debug",
		"init-container",
		"--gid", currentUser.Gid,
		"--home", currentUserHomeDir,
		"--shell", userShell,
		"--uid", currentUser.Uid,
		"--user", currentUser.Username,
	}

	entryPoint = append(entryPoint, slashHomeLink...)
	entryPoint = append(entryPoint, mediaLink...)
	entryPoint = append(entryPoint, mntLink...)

	createArgs := []string{
		"--log-level", logLevelString,
		"create",
		"--cgroupns", "host",
		"--dns", "none",
	}

	createArgs = append(createArgs, toolbxDelayEntryPointEnv...)
	createArgs = append(createArgs, toolbxFailEntryPointEnv...)

	createArgs = append(createArgs, []string{
		"--env", toolboxPathEnvArg,
	}...)

	createArgs = append(createArgs, xdgRuntimeDirEnv...)

	createArgs = append(createArgs, []string{
		"--hostname", "toolbx",
		"--ipc", "host",
		"--label", "com.github.containers.toolbox=true",
	}...)

	createArgs = append(createArgs, devPtsMount...)

	createArgs = append(createArgs, []string{
		"--name", container,
		"--network", "host",
		"--no-hosts",
		"--pid", "host",
		"--privileged",
		"--security-opt", "label=disable",
		"--ulimit", "host",
		"--userns", usernsArg,
		"--user", "root:root",
		"--volume", "/:/run/host:rslave",
		"--volume", "/dev:/dev:rslave",
		"--volume", dbusSystemSocketMountArg,
		"--volume", homeDirMountArg,
		"--volume", toolboxPathMountArg,
		"--volume", runtimeDirectoryMountArg,
	}...)

	createArgs = append(createArgs, avahiSocketMount...)
	createArgs = append(createArgs, kcmSocketMount...)
	createArgs = append(createArgs, mediaMount...)
	createArgs = append(createArgs, mntMount...)
	createArgs = append(createArgs, pcscSocketMount...)
	createArgs = append(createArgs, runMediaMount...)
	createArgs = append(createArgs, toolboxShMount...)

	createArgs = append(createArgs, []string{
		imageFull,
	}...)

	createArgs = append(createArgs, entryPoint...)

	logrus.Debugf("Creating container %s:", container)
	logrus.Debug("podman")
	for _, arg := range createArgs {
		logrus.Debugf("%s", arg)
	}

	s := spinner.New(spinner.CharSets[9], 500*time.Millisecond, spinner.WithWriterFile(os.Stdout))
	if logLevel := logrus.GetLevel(); logLevel < logrus.DebugLevel {
		s.Prefix = fmt.Sprintf("Creating container %s: ", container)
		s.Start()
		defer s.Stop()
	}

	if err := shell.Run("podman", nil, nil, nil, createArgs...); err != nil {
		return fmt.Errorf("failed to create container %s", container)
	}

	// The spinner must be stopped before showing the 'enter' hint below.
	s.Stop()

	if showCommandToEnter {
		fmt.Printf("Created container: %s\n", container)
		fmt.Printf("Enter with: %s\n", enterCommand)
	}

	return nil
}

func createHelp(cmd *cobra.Command, args []string) {
	if utils.IsInsideContainer() {
		if !utils.IsInsideToolboxContainer() {
			fmt.Fprintf(os.Stderr, "Error: this is not a Toolbx container\n")
			return
		}

		if _, err := utils.ForwardToHost(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err)
			return
		}

		return
	}

	if err := showManual("toolbox-create"); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		return
	}
}

func getDBusSystemSocket() (string, error) {
	logrus.Debug("Resolving path to the D-Bus system socket")

	address := os.Getenv("DBUS_SYSTEM_BUS_ADDRESS")
	if address == "" {
		address = "unix:path=/var/run/dbus/system_bus_socket"
	}

	addressSplit := strings.Split(address, "=")
	if len(addressSplit) != 2 {
		return "", errors.New("failed to get the path to the D-Bus system socket")
	}

	path := addressSplit[1]
	pathEvaled, err := filepath.EvalSymlinks(path)
	if err != nil {
		logrus.Debugf("Resolving path to the D-Bus system socket: failed to evaluate symbolic links in %s: %s",
			path,
			err)
		return "", errors.New("failed to resolve the path to the D-Bus system socket")
	}

	return pathEvaled, nil
}

func getEnterCommand(container string) string {
	var enterCommand string

	switch container {
	case utils.ContainerNameDefault:
		enterCommand = fmt.Sprintf("%s enter", executableBase)
	default:
		enterCommand = fmt.Sprintf("%s enter %s", executableBase, container)
	}

	return enterCommand
}

func getImageSizeFromRegistry(ctx context.Context, imageFull string) (string, error) {
	image, err := skopeo.Inspect(ctx, imageFull)
	if err != nil {
		return "", err
	}

	if image.LayersData == nil {
		return "", errors.New("'skopeo inspect' did not have LayersData")
	}

	var imageSizeFloat float64

	for _, layer := range image.LayersData {
		if layerSize, err := layer.Size.Float64(); err != nil {
			return "", err
		} else {
			imageSizeFloat += layerSize
		}
	}

	imageSizeHuman := units.HumanSize(imageSizeFloat)
	return imageSizeHuman, nil
}

func getImageSizeFromRegistryAsync(ctx context.Context, imageFull string) (<-chan string, <-chan error) {
	retValCh := make(chan string)
	errCh := make(chan error)

	go func() {
		imageSize, err := getImageSizeFromRegistry(ctx, imageFull)
		if err != nil {
			errCh <- err
			return
		}

		retValCh <- imageSize
	}()

	return retValCh, errCh
}

func getServiceSocket(serviceName string, unitName string) (string, error) {
	logrus.Debugf("Resolving path to the %s socket", serviceName)

	connection, err := dbus.SystemBus()
	if err != nil {
		logrus.Debugf("Resolving path to the %s socket: failed to connect to the D-Bus system instance: %s",
			serviceName,
			err)
		return "", errors.New("failed to connect to the D-Bus system instance")
	}

	unitNameEscaped := systemdPathBusEscape(unitName)
	unitPath := dbus.ObjectPath("/org/freedesktop/systemd1/unit/" + unitNameEscaped)
	unit := connection.Object("org.freedesktop.systemd1", unitPath)
	call := unit.Call("org.freedesktop.DBus.Properties.GetAll", 0, "")

	var result map[string]dbus.Variant
	err = call.Store(&result)
	if err != nil {
		logrus.Debugf("Resolving path to the %s socket: failed to get the properties of %s: %s",
			serviceName,
			unitName,
			err)
		return "", fmt.Errorf("failed to get the properties of %s", unitName)
	}

	listenVariant, listenFound := result["Listen"]
	if !listenFound {
		return "", fmt.Errorf("failed to find the Listen property of %s", unitName)
	}

	listenVariantSignature := listenVariant.Signature().String()
	if listenVariantSignature != "aav" {
		return "", errors.New("unknown reply from org.freedesktop.DBus.Properties.GetAll")
	}

	listenValue := listenVariant.Value()
	sockets := listenValue.([][]interface{})
	for _, socket := range sockets {
		if socket[0] == "Stream" {
			path := socket[1].(string)
			if !strings.HasPrefix(path, "/") {
				continue
			}

			pathEvaled, err := filepath.EvalSymlinks(path)
			if err != nil {
				continue
			}

			return pathEvaled, nil
		}
	}

	return "", fmt.Errorf("failed to find a SOCK_STREAM socket for %s", unitName)
}

func pullImage(image, release, authFile string) (bool, error) {
	if ok := utils.ImageReferenceCanBeID(image); ok {
		logrus.Debugf("Looking up image %s", image)
		if _, err := podman.ImageExists(image); err == nil {
			return true, nil
		}
	}

	hasDomain := utils.ImageReferenceHasDomain(image)

	if !hasDomain {
		imageLocal := "localhost/" + image
		logrus.Debugf("Looking up image %s", imageLocal)

		if _, err := podman.ImageExists(imageLocal); err == nil {
			return true, nil
		}
	}

	var imageFull string

	if hasDomain {
		imageFull = image
	} else {
		var err error
		imageFull, err = utils.GetFullyQualifiedImageFromDistros(image, release)
		if err != nil {
			return false, fmt.Errorf("image %s not found in local storage and known registries", image)
		}
	}

	logrus.Debugf("Looking up image %s", imageFull)
	if _, err := podman.ImageExists(imageFull); err == nil {
		return true, nil
	}

	domain := utils.ImageReferenceGetDomain(imageFull)
	if domain == "" {
		panicMsg := fmt.Sprintf("failed to get domain from %s", imageFull)
		panic(panicMsg)
	}

	promptForDownload := true
	var shouldPullImage bool

	if rootFlags.assumeYes || domain == "localhost" {
		promptForDownload = false
		shouldPullImage = true
	}

	if promptForDownload {
		if !term.IsTerminal(os.Stdin) || !term.IsTerminal(os.Stdout) {
			var builder strings.Builder
			fmt.Fprintf(&builder, "image required to create Toolbx container.\n")
			fmt.Fprintf(&builder, "Use option '--assumeyes' to download the image.\n")
			fmt.Fprintf(&builder, "Run '%s --help' for usage.", executableBase)

			errMsg := builder.String()
			return false, errors.New(errMsg)
		}

		shouldPullImage = showPromptForDownload(imageFull)
	}

	if !shouldPullImage {
		return false, nil
	}

	logrus.Debugf("Pulling image %s", imageFull)

	if logLevel := logrus.GetLevel(); logLevel < logrus.DebugLevel {
		s := spinner.New(spinner.CharSets[9], 500*time.Millisecond, spinner.WithWriterFile(os.Stdout))
		s.Prefix = fmt.Sprintf("Pulling %s: ", imageFull)
		s.Start()
		defer s.Stop()
	}

	if err := podman.Pull(imageFull, authFile); err != nil {
		var builder strings.Builder
		fmt.Fprintf(&builder, "failed to pull image %s\n", imageFull)
		fmt.Fprintf(&builder, "If it was a private image, log in with: podman login %s\n", domain)
		fmt.Fprintf(&builder, "Use '%s --verbose ...' for further details.", executableBase)

		errMsg := builder.String()
		return false, errors.New(errMsg)
	}

	return true, nil
}

func createPromptForDownload(imageFull, imageSize string) string {
	var prompt string
	if imageSize == "" {
		prompt = fmt.Sprintf("Download %s? [y/N]:", imageFull)
	} else {
		prompt = fmt.Sprintf("Download %s (%s)? [y/N]:", imageFull, imageSize)
	}

	return prompt
}

func showPromptForDownloadFirst(imageFull string) (bool, error) {
	prompt := createPromptForDownload(imageFull, " ... MB")

	parentCtx := context.Background()
	askCtx, askCancel := context.WithCancelCause(parentCtx)
	defer askCancel(errors.New("clean-up"))

	askCh, askErrCh := askForConfirmationAsync(askCtx, prompt, nil)

	imageSizeCtx, imageSizeCancel := context.WithCancelCause(parentCtx)
	defer imageSizeCancel(errors.New("clean-up"))

	imageSizeCh, imageSizeErrCh := getImageSizeFromRegistryAsync(imageSizeCtx, imageFull)

	var imageSize string
	var shouldPullImage bool

	select {
	case val := <-askCh:
		shouldPullImage = val
		cause := fmt.Errorf("%w: received confirmation without image size", context.Canceled)
		imageSizeCancel(cause)
	case err := <-askErrCh:
		shouldPullImage = false
		cause := fmt.Errorf("failed to ask for confirmation without image size: %w", err)
		imageSizeCancel(cause)
	case val := <-imageSizeCh:
		imageSize = val
		cause := fmt.Errorf("%w: received image size", context.Canceled)
		askCancel(cause)
	case err := <-imageSizeErrCh:
		cause := fmt.Errorf("failed to get image size: %w", err)
		askCancel(cause)
	}

	if imageSizeCtx.Err() != nil && askCtx.Err() == nil {
		cause := context.Cause(imageSizeCtx)
		logrus.Debugf("Show prompt for download: image size canceled: %s", cause)
		return shouldPullImage, nil
	}

	var done bool

	if imageSizeCtx.Err() == nil && askCtx.Err() != nil {
		select {
		case val := <-askCh:
			logrus.Debugf("Show prompt for download: received pending confirmation without image size")
			shouldPullImage = val
			done = true
		case err := <-askErrCh:
			logrus.Debugf("Show prompt for download: failed to ask for confirmation without image size: %s",
				err)
		}
	} else {
		panic("code should not be reached")
	}

	cause := context.Cause(askCtx)
	logrus.Debugf("Show prompt for download: ask canceled: %s", cause)

	if done {
		return shouldPullImage, nil
	}

	return false, &promptForDownloadError{imageSize}
}

func showPromptForDownloadSecond(imageFull string, errFirst *promptForDownloadError) bool {
	oldState, err := term.GetState(os.Stdin)
	if err != nil {
		logrus.Debugf("Show prompt for download: failed to get terminal state: %s", err)
		return false
	}

	defer term.SetState(os.Stdin, oldState)

	lockedState := term.NewStateFrom(oldState,
		term.WithVMIN(1),
		term.WithVTIME(0),
		term.WithoutECHO(),
		term.WithoutICANON())

	if err := term.SetState(os.Stdin, lockedState); err != nil {
		logrus.Debugf("Show prompt for download: failed to set terminal state: %s", err)
		return false
	}

	parentCtx := context.Background()
	discardCtx, discardCancel := context.WithCancelCause(parentCtx)
	defer discardCancel(errors.New("clean-up"))

	discardCh, discardErrCh := discardInputAsync(discardCtx)

	var prompt string
	if errors.Is(errFirst, context.Canceled) {
		prompt = createPromptForDownload(imageFull, errFirst.ImageSize)
	} else {
		prompt = createPromptForDownload(imageFull, "")
	}

	fmt.Printf("\r")

	askCtx, askCancel := context.WithCancelCause(parentCtx)
	defer askCancel(errors.New("clean-up"))

	var askForConfirmationPreFnDone bool
	askForConfirmationPreFn := func() error {
		defer discardCancel(errors.New("clean-up"))
		if askForConfirmationPreFnDone {
			return nil
		}

		// Erase to end of line
		fmt.Printf("\033[K")

		// Save the cursor position.
		fmt.Printf("\033[s")

		if err := term.SetState(os.Stdin, oldState); err != nil {
			return fmt.Errorf("failed to restore terminal state: %w", err)
		}

		cause := errors.New("terminal restored")
		discardCancel(cause)

		err := <-discardErrCh
		if !errors.Is(err, context.Canceled) {
			return fmt.Errorf("failed to discard input: %w", err)
		}

		logrus.Debugf("Show prompt for download: stopped discarding input: %s", err)

		discardTotal := <-discardCh
		logrus.Debugf("Show prompt for download: discarded input: %d bytes", discardTotal)

		if discardTotal == 0 {
			askForConfirmationPreFnDone = true
			return nil
		}

		if err := term.SetState(os.Stdin, lockedState); err != nil {
			return fmt.Errorf("failed to set terminal state: %w", err)
		}

		discardCtx, discardCancel = context.WithCancelCause(parentCtx)
		// A deferred call can't be used for this CancelCauseFunc,
		// because the 'discard' operation needs to continue running
		// until the next invocation of this function.  It relies on
		// the guarantee that AskForConfirmationAsync will always call
		// its askForConfirmationPreFunc as long as the function
		// returns errContinue.

		discardCh, discardErrCh = discardInputAsync(discardCtx)

		// Restore the cursor position
		fmt.Printf("\033[u")

		// Erase to end of line
		fmt.Printf("\033[K")

		fmt.Printf("...\n")
		return errContinue
	}

	askCh, askErrCh := askForConfirmationAsync(askCtx, prompt, askForConfirmationPreFn)
	var shouldPullImage bool

	select {
	case val := <-askCh:
		logrus.Debug("Show prompt for download: received confirmation with image size")
		shouldPullImage = val
	case err := <-askErrCh:
		logrus.Debugf("Show prompt for download: failed to ask for confirmation with image size: %s", err)
		shouldPullImage = false
	}

	return shouldPullImage
}

func showPromptForDownload(imageFull string) bool {
	fmt.Println("Image required to create Toolbx container.")

	shouldPullImage, err := showPromptForDownloadFirst(imageFull)
	if err == nil {
		return shouldPullImage
	}

	var errPromptForDownload *promptForDownloadError
	if !errors.As(err, &errPromptForDownload) {
		panicMsg := fmt.Sprintf("unexpected %T: %s", err, err)
		panic(panicMsg)
	}

	shouldPullImage = showPromptForDownloadSecond(imageFull, errPromptForDownload)
	return shouldPullImage
}

// systemdNeedsEscape checks whether a byte in a potential dbus ObjectPath needs to be escaped
func systemdNeedsEscape(i int, b byte) bool {
	// Escape everything that is not a-z-A-Z-0-9
	// Also escape 0-9 if it's the first character
	return strings.IndexByte(alphanum, b) == -1 ||
		(i == 0 && strings.IndexByte(num, b) != -1)
}

// systemdPathBusEscape sanitizes a constituent string of a dbus ObjectPath using the
// rules that systemd uses for serializing special characters.
func systemdPathBusEscape(path string) string {
	// Special case the empty string
	if len(path) == 0 {
		return "_"
	}
	n := []byte{}
	for i := 0; i < len(path); i++ {
		c := path[i]
		if systemdNeedsEscape(i, c) {
			e := fmt.Sprintf("_%x", c)
			n = append(n, []byte(e)...)
		} else {
			n = append(n, c)
		}
	}
	return string(n)
}

func (err *promptForDownloadError) Error() string {
	innerErr := err.Unwrap()
	errMsg := innerErr.Error()
	return errMsg
}

func (err *promptForDownloadError) Unwrap() error {
	if err.ImageSize == "" {
		return errors.New("failed to get image size")
	}

	return context.Canceled
}
