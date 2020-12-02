/*
 * Copyright © 2019 – 2020 Red Hat Inc.
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
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/containers/toolbox/pkg/podman"
	"github.com/containers/toolbox/pkg/shell"
	"github.com/containers/toolbox/pkg/toolbox"
	"github.com/containers/toolbox/pkg/utils"
	"github.com/godbus/dbus/v5"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
)

const (
	alpha    = `abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ`
	num      = `0123456789`
	alphanum = alpha + num
)

var (
	createFlags struct {
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

	enterCommand string
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new toolbox container",
	RunE:  create,
}

func init() {
	flags := createCmd.Flags()

	flags.StringVarP(&createFlags.container,
		"container",
		"c",
		"",
		"Assign a different name to the toolbox container.")

	flags.StringVarP(&createFlags.distro,
		"distro",
		"d",
		"",
		"Use different distribution for the toolbox")

	flags.StringVarP(&createFlags.image,
		"image",
		"i",
		"",
		"Change the name of the base image used to create the toolbox container.")

	flags.StringVarP(&createFlags.release,
		"release",
		"r",
		"",
		"Create a toolbox container for a different operating system release than the host.")

	createCmd.SetHelpFunc(createHelp)
	rootCmd.AddCommand(createCmd)
}

func create(cmd *cobra.Command, args []string) error {
	var image toolbox.Image
	var container, containerArg string
	var err error

	if utils.IsInsideContainer() {
		if !utils.IsInsideToolboxContainer() {
			return errors.New("this is not a toolbox container")
		}

		if _, err = utils.ForwardToHost(); err != nil {
			return err
		}

		return nil
	}

	if cmd.Flag("distro").Changed && cmd.Flag("image").Changed {
		fmt.Fprintf(os.Stderr, "Options --distro and --image cannot be used together\n")
		return err
	}

	if len(args) != 0 {
		container = args[0]
		containerArg = "CONTAINER"
	} else if createFlags.container != "" {
		container = createFlags.container
		containerArg = "--container"
	}

	if container != "" {
		if _, err = utils.IsContainerNameValid(container); err != nil {
			var builder strings.Builder
			fmt.Fprintf(&builder, "invalid argument for '%s'\n", containerArg)
			fmt.Fprintf(&builder, "Container names must match '%s'\n", utils.ContainerNameRegexp)
			fmt.Fprintf(&builder, "Run '%s --help' for usage.", executableBase)

			errMsg := builder.String()
			return errors.New(errMsg)
		}
	}

	if cmd.Flag("distro").Changed {
		if !toolbox.IsSystemSupported(createFlags.distro) {
			var builder strings.Builder
			fmt.Fprintf(&builder, "Toolbox does not support using distribution %s for toolboxes\n", createFlags.distro)
			fmt.Fprintf(&builder, "Run '%s create --help' to see a list of supported distributions.", executableBase)

			errMsg := builder.String()
			return errors.New(errMsg)
		}

		image, _ = toolbox.GetImageForSystem(createFlags.distro)
	} else {
		if toolbox.IsHostSystemSupported() {
			hostID, err := utils.GetHostID()
			if err != nil {
				return fmt.Errorf("There was an error while getting host's ID: %v", err)
			}

			image, _ = toolbox.GetImageForSystem(hostID)
		} else {
			image = toolbox.GetFallbackImage()
		}
	}

	if cmd.Flag("release").Changed {
		image.Tag = createFlags.release
	}

	if cmd.Flag("image").Changed {
		image.SetImageURI(createFlags.image)
	}

	if container == "" {
		container = image.CreateContainerName()
	}

	// Convenience warning in case the used image is "non-standard".
	if !image.IsImageKnown() {
		fmt.Fprintf(os.Stderr, "Used image is not officially supported. It might not work correctly.")
	}

	enterCommand = getEnterCommand(cmd, container)

	if err = createContainer(container, image, true); err != nil {
		return err
	}

	return nil
}

func createContainer(container string, image toolbox.Image, showCommandToEnter bool) error {
	if container == "" {
		panic("container not specified")
	}

	logrus.Debugf("Checking if container %s already exists", container)

	if exists, _ := podman.ContainerExists(container); exists {
		var builder strings.Builder
		fmt.Fprintf(&builder, "container %s already exists\n", container)
		fmt.Fprintf(&builder, "Enter with: %s\n", enterCommand)
		fmt.Fprintf(&builder, "Run '%s --help' for usage.", executableBase)

		errMsg := builder.String()
		return errors.New(errMsg)
	}

	pulled, err := pullImage(image)
	if err != nil {
		return err
	}
	if !pulled {
		return nil
	}

	imageFull, err := getFullyQualifiedImageName(image)
	if err != nil {
		return err
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

	logrus.Debug("Checking if 'podman create' supports '--ulimit host'")

	var ulimitHost []string

	if podman.CheckVersion("1.5.0") {
		logrus.Debug("'podman create' supports '--ulimit host'")
		ulimitHost = []string{"--ulimit", "host"}
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

	homeDirEvaled, err := filepath.EvalSymlinks(currentUser.HomeDir)
	if err != nil {
		return fmt.Errorf("failed to canonicalize %s", currentUser.HomeDir)
	}

	logrus.Debugf("%s canonicalized to %s", currentUser.HomeDir, homeDirEvaled)
	homeDirMountArg := homeDirEvaled + ":" + homeDirEvaled + ":rslave"

	usrMountFlags := "ro"
	isUsrReadWrite, err := isUsrReadWrite()
	if err != nil {
		return err
	}
	if isUsrReadWrite {
		usrMountFlags = "rw"
	}

	usrMountArg := "/usr:/run/host/usr:" + usrMountFlags + ",rslave"

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

	logrus.Debug("Checking if /mnt is a symbolic link to /var/mnt")

	var mntLink []string
	var mntMount []string

	mntPath, _ := filepath.EvalSymlinks("/mnt")
	if mntPath == "/var/mnt" {
		logrus.Debug("/mnt is a symbolic link to /var/mnt")
		mntLink = []string{"--mnt-link"}
	} else {
		mntMount = []string{"--volume", "/mnt:/mnt:rslave"}
	}

	var runMediaMount []string

	if utils.PathExists("/run/media") {
		runMediaMount = []string{"--volume", "/run/media:/run/media:rslave"}
	}

	logrus.Debug("Looking for toolbox.sh")

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
		"toolbox", "--verbose",
		"init-container",
		"--home", currentUser.HomeDir,
	}

	entryPoint = append(entryPoint, slashHomeLink...)
	entryPoint = append(entryPoint, mediaLink...)
	entryPoint = append(entryPoint, mntLink...)

	entryPoint = append(entryPoint, []string{
		"--monitor-host",
		"--shell", userShell,
		"--uid", currentUser.Uid,
		"--user", currentUser.Username,
	}...)

	createArgs := []string{
		"--log-level", logLevelString,
		"create",
		"--dns", "none",
		"--env", toolboxPathEnvArg,
	}

	createArgs = append(createArgs, xdgRuntimeDirEnv...)

	createArgs = append(createArgs, []string{
		"--hostname", "toolbox",
		"--ipc", "host",
		"--label", "com.github.containers.toolbox=true",
		"--label", "com.github.debarshiray.toolbox=true",
	}...)

	createArgs = append(createArgs, devPtsMount...)

	createArgs = append(createArgs, []string{
		"--name", container,
		"--network", "host",
		"--no-hosts",
		"--pid", "host",
		"--privileged",
		"--security-opt", "label=disable",
	}...)

	createArgs = append(createArgs, ulimitHost...)

	createArgs = append(createArgs, []string{
		"--userns", usernsArg,
		"--user", "root:root",
		"--volume", "/boot:/run/host/boot:rslave",
		"--volume", "/etc:/run/host/etc",
		"--volume", "/dev:/dev:rslave",
		"--volume", "/run:/run/host/run:rslave",
		"--volume", "/tmp:/run/host/tmp:rslave",
		"--volume", "/var:/run/host/var:rslave",
		"--volume", dbusSystemSocketMountArg,
		"--volume", homeDirMountArg,
		"--volume", toolboxPathMountArg,
		"--volume", usrMountArg,
		"--volume", runtimeDirectoryMountArg,
	}...)

	createArgs = append(createArgs, avahiSocketMount...)
	createArgs = append(createArgs, kcmSocketMount...)
	createArgs = append(createArgs, mediaMount...)
	createArgs = append(createArgs, mntMount...)
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

	s := spinner.New(spinner.CharSets[9], 500*time.Millisecond)

	stdoutFd := os.Stdout.Fd()
	stdoutFdInt := int(stdoutFd)
	if logLevel := logrus.GetLevel(); logLevel < logrus.DebugLevel && terminal.IsTerminal(stdoutFdInt) {
		s.Prefix = fmt.Sprintf("Creating container %s: ", container)
		s.Writer = os.Stdout
		s.Start()
		defer s.Stop()
	}

	if err := shell.Run("podman", nil, nil, nil, createArgs...); err != nil {
		return fmt.Errorf("failed to create container %s", container)
	}

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
			fmt.Fprintf(os.Stderr, "Error: this is not a toolbox container\n")
			return
		}

		if _, err := utils.ForwardToHost(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err)
			return
		}

		return
	}

	if err := utils.ShowManual("toolbox-create"); err != nil {
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
		return "", errors.New("failed to resolve the path to the D-Bus system socket")
	}

	return pathEvaled, nil
}

func getEnterCommand(cmd *cobra.Command, container string) string {
	var enterCommand string

	if container == utils.ContainerNameDefault {
		enterCommand = fmt.Sprintf("%s enter", executableBase)
	} else if cmd.Flag("release").Changed {
		enterCommand = fmt.Sprintf("%s enter --release %s", executableBase, createFlags.release)
	} else {
		enterCommand = fmt.Sprintf("%s enter %s", executableBase, container)
	}

	return enterCommand
}

func getFullyQualifiedImageName(img toolbox.Image) (string, error) {
	var image string = img.GetImageURI()
	var imageFull string

	logrus.Debugf("Resolving fully qualified name for image %s", image)

	if utils.ImageReferenceHasDomain(image) {
		imageFull = image
	} else {
		info, err := podman.Inspect("image", image)
		if err != nil {
			return "", fmt.Errorf("failed to inspect image %s", image)
		}

		if info["RepoTags"] == nil {
			return "", fmt.Errorf("missing RepoTag for image %s", image)
		}

		repoTags := info["RepoTags"].([]interface{})
		if len(repoTags) == 0 {
			return "", fmt.Errorf("empty RepoTag for image %s", image)
		}

		imageFull = repoTags[0].(string)
	}

	logrus.Debugf("Resolved image %s to %s", image, imageFull)

	return imageFull, nil
}

func getServiceSocket(serviceName string, unitName string) (string, error) {
	logrus.Debugf("Resolving path to the %s socket", serviceName)

	connection, err := dbus.SystemBus()
	if err != nil {
		return "", errors.New("failed to connect to the D-Bus system instance")
	}

	unitNameEscaped := systemdPathBusEscape(unitName)
	unitPath := dbus.ObjectPath("/org/freedesktop/systemd1/unit/" + unitNameEscaped)
	unit := connection.Object("org.freedesktop.systemd1", unitPath)
	call := unit.Call("org.freedesktop.DBus.Properties.GetAll", 0, "")

	var result map[string]dbus.Variant
	err = call.Store(&result)
	if err != nil {
		errMsg := fmt.Sprintf("failed to get the properties of %s", unitName)
		return "", errors.New(errMsg)
	}

	listenVariant, listenFound := result["Listen"]
	if !listenFound {
		errMsg := fmt.Sprintf("failed to find the Listen property of %s", unitName)
		return "", errors.New(errMsg)
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

	errMsg := fmt.Sprintf("failed to find a SOCK_STREAM socket for %s", unitName)
	return "", errors.New(errMsg)
}

func isUsrReadWrite() (bool, error) {
	logrus.Debug("Checking if /usr is mounted read-only or read-write")

	mountPoint, err := utils.GetMountPoint("/usr")
	if err != nil {
		return false, fmt.Errorf("failed to get the mount-point of /usr: %s", err)
	}

	logrus.Debugf("Mount-point of /usr is %s", mountPoint)

	mountFlags, err := utils.GetMountOptions(mountPoint)
	if err != nil {
		return false, fmt.Errorf("failed to get the mount options of %s: %s", mountPoint, err)
	}

	logrus.Debugf("Mount flags of /usr on the host are %s", mountFlags)

	if !strings.Contains(mountFlags, "ro") {
		return true, nil
	}

	return false, nil
}

func pullImage(image toolbox.Image) (bool, error) {
	var registries []string = []string{"localhost", "registry.fedoraproject.org"}
	var imageURI string = image.GetImageURI()
	var err error

	if _, err := utils.ImageReferenceCanBeID(imageURI); err == nil {
		logrus.Debugf("Looking for image %s", imageURI)

		if _, err := podman.ImageExists(imageURI); err == nil {
			return true, nil
		}
	}

	if !utils.ImageReferenceHasDomain(imageURI) {
		imageFound := false

		oldRegistry := image.Registry
		for _, registry := range registries {
			image.Registry = registry
			imageURI = image.GetImageURI()
			logrus.Debugf("Looking for image %s", imageURI)

			if _, err = podman.ImageExists(imageURI); err == nil {
				imageFound = true
				break
			}
		}

		if !imageFound {
			image.Registry = oldRegistry
			return false, fmt.Errorf("image %s does not exist", image.GetImageURI())
		}
	}

	registry := image.Registry
	if registry == "" {
		panicMsg := fmt.Sprintf("failed to get domain from %s", imageURI)
		panic(panicMsg)
	}

	promptForDownload := true
	var shouldPullImage bool

	if rootFlags.assumeYes || registry == "localhost" {
		promptForDownload = false
		shouldPullImage = true
	}

	if promptForDownload {
		fmt.Println("Image required to create toolbox container.")

		prompt := fmt.Sprintf("Download %s (500MB)? [y/N]:", imageURI)
		shouldPullImage = utils.AskForConfirmation(prompt)
	}

	if !shouldPullImage {
		return false, nil
	}

	logrus.Debugf("Pulling image %s", imageURI)

	stdoutFd := os.Stdout.Fd()
	stdoutFdInt := int(stdoutFd)
	if logLevel := logrus.GetLevel(); logLevel < logrus.DebugLevel && terminal.IsTerminal(stdoutFdInt) {
		s := spinner.New(spinner.CharSets[9], 500*time.Millisecond)
		s.Prefix = fmt.Sprintf("Pulling %s: ", imageURI)
		s.Writer = os.Stdout
		s.Start()
		defer s.Stop()
	}

	if err := podman.Pull(imageURI); err != nil {
		return false, fmt.Errorf("failed to pull image %s", imageURI)
	}

	return true, nil
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
