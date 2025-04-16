/*
 * Copyright © 2019 – 2024 Red Hat Inc.
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
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/containers/toolbox/pkg/shell"
	"github.com/containers/toolbox/pkg/utils"
	"github.com/fsnotify/fsnotify"
	"github.com/google/renameio/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/sys/unix"
	"tags.cncf.io/container-device-interface/pkg/cdi"
	"tags.cncf.io/container-device-interface/specs-go"
)

var (
	initContainerFlags struct {
		gid         int
		home        string
		homeLink    bool
		mediaLink   bool
		mntLink     bool
		monitorHost bool
		shell       string
		uid         int
		user        string
	}

	initContainerMounts = []struct {
		containerPath string
		source        string
		flags         string
	}{
		{"/etc/machine-id", "/run/host/etc/machine-id", ""},
		{"/run/libvirt", "/run/host/run/libvirt", ""},
		{"/run/systemd/journal", "/run/host/run/systemd/journal", ""},
		{"/run/systemd/resolve", "/run/host/run/systemd/resolve", ""},
		{"/run/systemd/sessions", "/run/host/run/systemd/sessions", ""},
		{"/run/systemd/system", "/run/host/run/systemd/system", ""},
		{"/run/systemd/users", "/run/host/run/systemd/users", ""},
		{"/run/udev/data", "/run/host/run/udev/data", ""},
		{"/run/udev/tags", "/run/host/run/udev/tags", ""},
		{"/tmp", "/run/host/tmp", "rslave"},
		{"/var/lib/flatpak", "/run/host/var/lib/flatpak", ""},
		{"/var/lib/libvirt", "/run/host/var/lib/libvirt", ""},
		{"/var/lib/systemd/coredump", "/run/host/var/lib/systemd/coredump", ""},
		{"/var/log/journal", "/run/host/var/log/journal", ""},
		{"/var/mnt", "/run/host/var/mnt", "rslave"},
	}
)

var initContainerCmd = &cobra.Command{
	Use:    "init-container",
	Short:  "Initialize a running container",
	Hidden: true,
	RunE:   initContainer,
}

func init() {
	flags := initContainerCmd.Flags()

	flags.IntVar(&initContainerFlags.gid,
		"gid",
		0,
		"Create a user inside the Toolbx container whose numerical group ID is GID")

	flags.StringVar(&initContainerFlags.home,
		"home",
		"",
		"Create a user inside the Toolbx container whose login directory is HOME")
	if err := initContainerCmd.MarkFlagRequired("home"); err != nil {
		panic("Could not mark flag --home as required")
	}

	flags.BoolVar(&initContainerFlags.homeLink,
		"home-link",
		false,
		"Make /home a symbolic link to /var/home")

	flags.BoolVar(&initContainerFlags.mediaLink,
		"media-link",
		false,
		"Make /media a symbolic link to /run/media")

	flags.BoolVar(&initContainerFlags.mntLink, "mnt-link", false, "Make /mnt a symbolic link to /var/mnt")

	flags.BoolVar(&initContainerFlags.monitorHost,
		"monitor-host",
		true,
		"Deprecated, does nothing")
	if err := flags.MarkDeprecated("monitor-host", "it does nothing"); err != nil {
		panicMsg := fmt.Sprintf("cannot mark --monitor-host as deprecated: %s", err)
		panic(panicMsg)
	}

	flags.StringVar(&initContainerFlags.shell,
		"shell",
		"",
		"Create a user inside the Toolbx container whose login shell is SHELL")
	if err := initContainerCmd.MarkFlagRequired("shell"); err != nil {
		panic("Could not mark flag --shell as required")
	}

	flags.IntVar(&initContainerFlags.uid,
		"uid",
		0,
		"Create a user inside the Toolbx container whose numerical user ID is UID")
	if err := initContainerCmd.MarkFlagRequired("uid"); err != nil {
		panic("Could not mark flag --uid as required")
	}

	flags.StringVar(&initContainerFlags.user,
		"user",
		"",
		"Create a user inside the Toolbx container whose login name is USER")
	if err := initContainerCmd.MarkFlagRequired("user"); err != nil {
		panic("Could not mark flag --user as required")
	}

	initContainerCmd.SetHelpFunc(initContainerHelp)
	rootCmd.AddCommand(initContainerCmd)
}

func initContainer(cmd *cobra.Command, args []string) error {
	const factoryInitializedFlag string = "/var/lib/toolbox/factory-initialized"

	if !utils.IsInsideContainer() {
		var builder strings.Builder
		fmt.Fprintf(&builder, "the 'init-container' command can only be used inside containers\n")
		fmt.Fprintf(&builder, "Run '%s --help' for usage.", executableBase)

		errMsg := builder.String()
		return errors.New(errMsg)
	}

	if !cmd.Flag("gid").Changed {
		initContainerFlags.gid = initContainerFlags.uid
	}

	utils.EnsureXdgRuntimeDirIsSet(initContainerFlags.uid)

	logrus.Debug("Creating /run/.toolboxenv")

	toolboxEnvFile, err := os.Create("/run/.toolboxenv")
	if err != nil {
		return errors.New("failed to create /run/.toolboxenv")
	}

	defer toolboxEnvFile.Close()

	if toolbxDelayEntryPoint, ok := getDelayEntryPoint(); ok {
		delayString := toolbxDelayEntryPoint.String()
		logrus.Debugf("Adding a delay of %s", delayString)
		time.Sleep(toolbxDelayEntryPoint)
	}

	if toolbxFailEntryPoint, ok := getFailEntryPoint(); ok {
		var builder strings.Builder
		fmt.Fprintf(&builder, "TOOLBX_FAIL_ENTRY_POINT is set")
		if toolbxFailEntryPoint > 1 {
			fmt.Fprintf(&builder, "\n")
			fmt.Fprintf(&builder, "This environment variable should only be set when testing.")
		}

		errMsg := builder.String()
		return errors.New(errMsg)
	}

	if !utils.PathExists(factoryInitializedFlag) {
		if utils.PathExists("/run/host/etc") {
			logrus.Debug("Path /run/host/etc exists")

			if _, err := os.Readlink("/etc/host.conf"); err != nil {
				if err := redirectPath("/etc/host.conf",
					"/run/host/etc/host.conf",
					false); err != nil {
					return err
				}
			}

			if _, err := os.Readlink("/etc/hosts"); err != nil {
				if err := redirectPath("/etc/hosts",
					"/run/host/etc/hosts",
					false); err != nil {
					return err
				}
			}

			if localtimeTarget, err := os.Readlink("/etc/localtime"); err != nil ||
				localtimeTarget != "/run/host/etc/localtime" {
				if err := redirectPath("/etc/localtime",
					"/run/host/etc/localtime",
					false); err != nil {
					return err
				}
			}

			if err := updateTimeZoneFromLocalTime(); err != nil {
				return err
			}

			if resolvConfTarget, err := os.Readlink("/etc/resolv.conf"); err != nil ||
				resolvConfTarget != "/run/host/etc/resolv.conf" {
				if err := redirectPath("/etc/resolv.conf",
					"/run/host/etc/resolv.conf",
					false); err != nil {
					return err
				}
			}
		}

		if initContainerFlags.mediaLink {
			if _, err := os.Readlink("/media"); err != nil {
				if err = redirectPath("/media", "/run/media", true); err != nil {
					return err
				}
			}
		}

		if initContainerFlags.mntLink {
			if _, err := os.Readlink("/mnt"); err != nil {
				if err := redirectPath("/mnt", "/var/mnt", true); err != nil {
					return err
				}
			}
		}

		for _, mount := range initContainerMounts {
			if err := mountBind(mount.containerPath, mount.source, mount.flags); err != nil {
				return err
			}
		}

		if utils.PathExists("/sys/fs/selinux") {
			if err := mountBind("/sys/fs/selinux", "/usr/share/empty", ""); err != nil {
				return err
			}
		}

		if err := configureUsers(initContainerFlags.uid,
			initContainerFlags.user,
			initContainerFlags.home,
			initContainerFlags.shell,
			initContainerFlags.homeLink); err != nil {
			return err
		}

		if utils.PathExists("/etc/krb5.conf.d") && !utils.PathExists("/etc/krb5.conf.d/kcm_default_ccache") {
			logrus.Debug("Setting KCM as the default Kerberos credential cache")

			kcmConfigString := `# Written by Toolbx
# https://github.com/containers/toolbox
#
# # To disable the KCM credential cache, comment out the following lines.

[libdefaults]
    default_ccache_name = KCM:
`

			kcmConfigBytes := []byte(kcmConfigString)
			if err := ioutil.WriteFile("/etc/krb5.conf.d/kcm_default_ccache",
				kcmConfigBytes,
				0644); err != nil {
				return errors.New("failed to set KCM as the default Kerberos credential cache")
			}
		}

		if utils.PathExists("/usr/lib/rpm/macros.d") {
			logrus.Debug("Configuring RPM to ignore bind mounts")

			var builder strings.Builder
			fmt.Fprintf(&builder, "# Written by Toolbx\n")
			fmt.Fprintf(&builder, "# https://github.com/containers/toolbox\n")
			fmt.Fprintf(&builder, "\n")
			fmt.Fprintf(&builder, "%%_netsharedpath /dev:/media:/mnt:/proc:/sys:/tmp:/var/lib/flatpak:/var/lib/libvirt\n")

			rpmConfigString := builder.String()
			rpmConfigBytes := []byte(rpmConfigString)
			if err := ioutil.WriteFile("/usr/lib/rpm/macros.d/macros.toolbox",
				rpmConfigBytes,
				0644); err != nil {
				return fmt.Errorf("failed to configure RPM to ignore bind mounts: %w", err)
			}
		}

		logrus.Debugf("Creating factory initialization flag %s", factoryInitializedFlag)

		err := os.MkdirAll(path.Base(factoryInitializedFlag), 0700)
		if err != nil {
			return errors.New("failed to create toolbox data directory")
		}

		factoryInitializedFlagFile, err := os.Create(factoryInitializedFlag)
		if err != nil {
			return errors.New("failed to create initialization flag")
		}

		defer factoryInitializedFlagFile.Close()
	}

	uidString := strconv.Itoa(initContainerFlags.uid)
	targetUser, err := user.LookupId(uidString)
	if err != nil {
		return fmt.Errorf("failed to look up user ID %s: %w", uidString, err)
	}

	cdiFileForNvidia, err := getCDIFileForNvidia(targetUser)
	if err != nil {
		return err
	}

	logrus.Debugf("Loading Container Device Interface for NVIDIA from file %s", cdiFileForNvidia)

	cdiSpecForNvidia, err := loadCDISpecFrom(cdiFileForNvidia)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			logrus.Debugf("Loading Container Device Interface for NVIDIA: file %s not found",
				cdiFileForNvidia)
		} else {
			logrus.Debugf("Loading Container Device Interface for NVIDIA: failed: %s", err)
			return errors.New("failed to load Container Device Interface for NVIDIA")
		}
	}

	if cdiSpecForNvidia != nil {
		if err := applyCDISpecForNvidia(cdiSpecForNvidia); err != nil {
			return err
		}
	}

	logrus.Debug("Setting up daily ticker")

	tickerDaily := time.NewTicker(24 * time.Hour)
	defer tickerDaily.Stop()

	logrus.Debug("Setting up watches for file system events")

	var watcherForHostErrors chan error
	var watcherForHostEvents chan fsnotify.Event

	watcherForHost, err := fsnotify.NewWatcher()
	if err != nil {
		if errors.Is(err, unix.EMFILE) || errors.Is(err, unix.ENFILE) || errors.Is(err, unix.ENOMEM) {
			logrus.Debugf("Setting up watches for file system events: failed to create Watcher: %s", err)
		} else {
			return fmt.Errorf("failed to create Watcher: %w", err)
		}
	}

	if watcherForHost != nil {
		defer watcherForHost.Close()

		watcherForHostErrors = watcherForHost.Errors
		watcherForHostEvents = watcherForHost.Events

		if err := watcherForHost.Add("/run/host/etc"); err != nil {
			if errors.Is(err, unix.ENOMEM) || errors.Is(err, unix.ENOSPC) {
				logrus.Debugf("Setting up watches for file system events: failed to add path: %s", err)
			} else {
				return fmt.Errorf("failed to add path: %w", err)
			}
		}
	}

	logrus.Debug("Finished initializing container")

	toolboxRuntimeDirectory, err := utils.GetRuntimeDirectory(targetUser)
	if err != nil {
		return err
	}

	pid := os.Getpid()
	initializedStamp := fmt.Sprintf("%s/container-initialized-%d", toolboxRuntimeDirectory, pid)

	logrus.Debugf("Creating initialization stamp %s", initializedStamp)

	initializedStampFile, err := os.Create(initializedStamp)
	if err != nil {
		return errors.New("failed to create initialization stamp")
	}

	defer initializedStampFile.Close()

	if err := initializedStampFile.Chown(initContainerFlags.uid, initContainerFlags.gid); err != nil {
		return errors.New("failed to change ownership of initialization stamp")
	}

	logrus.Debug("Listening to file system and ticker events")

	go runUpdateDb()

	for {
		select {
		case event := <-tickerDaily.C:
			handleDailyTick(event)
		case event := <-watcherForHostEvents:
			handleFileSystemEvent(event)
		case err := <-watcherForHostErrors:
			logrus.Warnf("Received an error from the file system watcher: %v", err)
		}
	}

	// code should not be reached
}

func initContainerHelp(cmd *cobra.Command, args []string) {
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

	if err := showManual("toolbox-init-container"); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		return
	}
}

func applyCDISpecForNvidia(spec *specs.Spec) error {
	if spec == nil {
		panic("spec not specified")
	}

	logrus.Debug("Applying Container Device Interface for NVIDIA")

	for _, mount := range spec.ContainerEdits.Mounts {
		if err := (&cdi.Mount{Mount: mount}).Validate(); err != nil {
			logrus.Debugf("Applying Container Device Interface for NVIDIA: invalid mount: %s", err)
			return errors.New("invalid mount in Container Device Interface for NVIDIA")
		}

		if mount.Type == "" {
			mount.Type = "bind"
		}

		if mount.Type != "bind" {
			logrus.Debugf("Applying Container Device Interface for NVIDIA: unknown mount type %s",
				mount.Type)
			continue
		}

		flags := strings.Join(mount.Options, ",")
		hostPath := filepath.Join(string(filepath.Separator), "run", "host", mount.HostPath)
		if err := mountBind(mount.ContainerPath, hostPath, flags); err != nil {
			logrus.Debugf("Applying Container Device Interface for NVIDIA: %s", err)
			return errors.New("failed to apply mount from Container Device Interface for NVIDIA")
		}
	}

	for _, hook := range spec.ContainerEdits.Hooks {
		if err := (&cdi.Hook{Hook: hook}).Validate(); err != nil {
			logrus.Debugf("Applying Container Device Interface for NVIDIA: invalid hook: %s", err)
			return errors.New("invalid hook in Container Device Interface for NVIDIA")
		}

		if hook.HookName != cdi.CreateContainerHook {
			logrus.Debugf("Applying Container Device Interface for NVIDIA: unknown hook name %s",
				hook.HookName)
			continue
		}

		if len(hook.Args) >= 2 &&
			hook.Args[0] == "nvidia-cdi-hook" &&
			hook.Args[1] == "create-symlinks" {
			hookArgs := hook.Args[2:]
			if err := applyCDISpecForNvidiaHookCreateSymlinks(hookArgs); err != nil {
				logrus.Debugf("Applying Container Device Interface for NVIDIA: %s", err)
				return errors.New("failed to create symlinks for Container Device Interface for NVIDIA")
			}

			continue
		} else if len(hook.Args) >= 2 &&
			hook.Args[0] == "nvidia-cdi-hook" &&
			hook.Args[1] == "update-ldcache" {
			hookArgs := hook.Args[2:]
			if err := applyCDISpecForNvidiaHookUpdateLDCache(hookArgs); err != nil {
				logrus.Debugf("Applying Container Device Interface for NVIDIA: %s", err)
				return errors.New("failed to update ldcache for Container Device Interface for NVIDIA")
			}

			continue
		}

		logrus.Debug("Applying Container Device Interface for NVIDIA: unknown hook arguments:")
		for _, arg := range hook.Args {
			logrus.Debugf("%s", arg)
		}
	}

	return nil
}

func applyCDISpecForNvidiaHookCreateSymlinks(hookArgs []string) error {
	var linkFlag bool

	for _, hookArg := range hookArgs {
		if hookArg == "--link" {
			linkFlag = true
			continue
		}

		if linkFlag {
			linkFlag = false
			if linkParts := strings.Split(hookArg, "::"); len(linkParts) == 2 {
				existingTarget := linkParts[0]

				newLink := linkParts[1]
				if !filepath.IsAbs(newLink) {
					return fmt.Errorf("invalid --link argument: link %s is not an absolute path",
						newLink)
				}

				if err := createSymbolicLink(existingTarget, newLink); err != nil {
					return err
				}
			} else {
				return fmt.Errorf("invalid --link argument: %s not in '<target>::<link>' format",
					hookArg)
			}
		}
	}

	if linkFlag {
		return errors.New("missing --link argument")
	}

	return nil
}

func applyCDISpecForNvidiaHookUpdateLDCache(hookArgs []string) error {
	var folderFlag bool
	var folders []string

	for _, hookArg := range hookArgs {
		if hookArg == "--folder" {
			folderFlag = true
			continue
		}

		if folderFlag {
			folders = append(folders, hookArg)
		}

		folderFlag = false
	}

	if err := ldConfig("toolbx-nvidia.conf", folders); err != nil {
		return err
	}

	return nil
}

func configureUsers(targetUserUid int, targetUser, targetUserHome, targetUserShell string, homeLink bool) error {
	if homeLink {
		if err := redirectPath("/home", "/var/home", true); err != nil {
			return err
		}
	}

	sudoGroup, err := utils.GetGroupForSudo()
	if err != nil {
		return fmt.Errorf("failed to get group for sudo: %w", err)
	}

	if _, err := user.Lookup(targetUser); err != nil {
		logrus.Debugf("Adding user %s with UID %d:", targetUser, targetUserUid)

		useraddArgs := []string{
			"--groups", sudoGroup,
			"--home-dir", targetUserHome,
			"--no-create-home",
			"--password", "",
			"--shell", targetUserShell,
			"--uid", fmt.Sprint(targetUserUid),
			targetUser,
		}

		logrus.Debug("useradd")
		for _, arg := range useraddArgs {
			logrus.Debugf("%s", arg)
		}

		if err := shell.Run("useradd", nil, nil, nil, useraddArgs...); err != nil {
			return fmt.Errorf("failed to add user %s with UID %d: %w", targetUser, targetUserUid, err)
		}
	} else {
		logrus.Debugf("Modifying user %s with UID %d:", targetUser, targetUserUid)

		usermodArgs := []string{
			"--append",
			"--groups", sudoGroup,
			"--home", targetUserHome,
			"--password", "",
			"--shell", targetUserShell,
			"--uid", fmt.Sprint(targetUserUid),
			targetUser,
		}

		logrus.Debug("usermod")
		for _, arg := range usermodArgs {
			logrus.Debugf("%s", arg)
		}

		if err := shell.Run("usermod", nil, nil, nil, usermodArgs...); err != nil {
			return fmt.Errorf("failed to modify user %s with UID %d: %w", targetUser, targetUserUid, err)
		}
	}

	logrus.Debug("Removing password for user root")

	var stderr strings.Builder
	if err := shell.Run("passwd", nil, nil, &stderr, "--delete", "root"); err != nil {
		errString := stderr.String()
		logrus.Debugf("Removing password for user root: failed: %s", errString)
		return fmt.Errorf("failed to remove password for root: %w", err)
	}

	return nil
}

func createSymbolicLink(existingTarget, newLink string) error {
	logrus.Debugf("Creating symbolic link with target %s and link %s", existingTarget, newLink)

	newLinkDir := filepath.Dir(newLink)
	if err := os.MkdirAll(newLinkDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", newLinkDir, err)
	}

	if err := os.Symlink(existingTarget, newLink); err != nil {
		var errLink *os.LinkError
		if errors.As(err, &errLink) {
			if errors.Is(err, os.ErrExist) {
				logrus.Debugf("Creating symbolic link: file %s already exists", newLink)
				return nil
			}
		}

		return fmt.Errorf("failed to create symbolic link: %w", err)
	}

	return nil
}

func getDelayEntryPoint() (time.Duration, bool) {
	valueString := os.Getenv("TOOLBX_DELAY_ENTRY_POINT")
	if valueString == "" {
		return 0, false
	}

	if valueN, err := strconv.Atoi(valueString); valueN > 0 && err == nil {
		delay := time.Duration(valueN) * time.Second
		return delay, true
	}

	return 0, false
}

func getFailEntryPoint() (uint, bool) {
	valueString := os.Getenv("TOOLBX_FAIL_ENTRY_POINT")
	if valueString == "" {
		return 0, false
	}

	if valueN, err := strconv.Atoi(valueString); valueN > 0 && err == nil {
		return uint(valueN), true
	}

	return 0, false
}

func handleDailyTick(event time.Time) {
	eventString := event.String()
	logrus.Debugf("Handling daily tick %s", eventString)

	runUpdateDb()
}

func handleFileSystemEvent(event fsnotify.Event) {
	eventOpString := event.Op.String()
	logrus.Debugf("Handling file system event: operation %s on %s", eventOpString, event.Name)

	if event.Name == "/run/host/etc/localtime" {
		if err := updateTimeZoneFromLocalTime(); err != nil {
			logrus.Warnf("Failed to handle changes to the host's /etc/localtime: %v", err)
		}
	}
}

func ldConfig(configFileBase string, dirs []string) error {
	logrus.Debug("Updating dynamic linker cache")

	var args []string

	if !utils.PathExists("/etc/ld.so.cache") {
		logrus.Debug("Updating dynamic linker cache: no /etc/ld.so.cache found")
		args = append(args, "-N")
	}

	if utils.PathExists("/etc/ld.so.conf.d") {
		if len(dirs) > 0 {
			var builder strings.Builder
			builder.WriteString("# Written by Toolbx\n")
			builder.WriteString("# https://containertoolbx.org/\n")
			builder.WriteString("\n")

			configured := make(map[string]struct{})

			for _, dir := range dirs {
				if _, ok := configured[dir]; ok {
					continue
				}

				configured[dir] = struct{}{}
				builder.WriteString(dir)
				builder.WriteString("\n")
			}

			dirConfigString := builder.String()
			dirConfigBytes := []byte(dirConfigString)
			configFile := filepath.Join("/etc/ld.so.conf.d", configFileBase)
			if err := renameio.WriteFile(configFile, dirConfigBytes, 0644); err != nil {
				logrus.Debugf("Updating dynamic linker cache: failed to update configuration: %s", err)
				return errors.New("failed to update dynamic linker cache configuration")
			}
		}
	} else {
		logrus.Debug("Updating dynamic linker cache: no /etc/ld.so.conf.d found")
		args = append(args, dirs...)
	}

	if err := shell.Run("ldconfig", nil, nil, nil, args...); err != nil {
		logrus.Debugf("Updating dynamic linker cache: failed: %s", err)
		return errors.New("failed to update dynamic linker cache")
	}

	return nil
}

func loadCDISpecFrom(path string) (*specs.Spec, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	spec, err := cdi.ParseSpec(data)
	if err != nil {
		return nil, err
	}
	if spec == nil {
		return nil, errors.New("missing data")
	}

	return spec, nil
}

func mountBind(containerPath, source, flags string) error {
	fi, err := os.Stat(source)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}

		return fmt.Errorf("failed to stat %s", source)
	}

	fileMode := fi.Mode()

	if fileMode.IsDir() {
		logrus.Debugf("Creating directory %s", containerPath)
		if err := os.MkdirAll(containerPath, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", containerPath, err)
		}
	} else if fileMode.IsRegular() || fileMode&os.ModeSocket != 0 {
		logrus.Debugf("Creating regular file %s", containerPath)

		containerPathDir := filepath.Dir(containerPath)
		if err := os.MkdirAll(containerPathDir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", containerPathDir, err)
		}

		containerPathFile, err := os.Create(containerPath)
		if err != nil && !os.IsExist(err) {
			return fmt.Errorf("failed to create regular file %s: %w", containerPath, err)
		}

		defer containerPathFile.Close()
	}

	logrus.Debugf("Binding %s to %s", containerPath, source)

	args := []string{
		"--rbind",
	}

	if flags != "" {
		args = append(args, []string{"-o", flags}...)
	}

	args = append(args, []string{source, containerPath}...)

	if err := shell.Run("mount", nil, nil, nil, args...); err != nil {
		return fmt.Errorf("failed to bind %s to %s", containerPath, source)
	}

	return nil
}

// redirectPath serves for creating symbolic links for crucial system
// configuration files to their counterparts on the host's file system.
//
// containerPath and target must be absolute paths
//
// If the target itself is a symbolic link, redirectPath will evaluate the
// link. If it's valid then redirectPath will link containerPath to target.
// If it's not, then redirectPath will still proceed with the linking in two
// different ways depending whether target is an absolute or a relative link:
//
//   - absolute - target's destination will be prepended with /run/host, and
//     containerPath will be linked to the resulting path. This is an attempt
//     to unbreak things, but if it doesn't work then it's the user's
//     responsibility to fix it up.
//
//     This is meant to address the common case where /etc/resolv.conf on the
//     host (ie., /run/host/etc/resolv.conf inside the container) is an
//     absolute symbolic link to /run/systemd/resolve/stub-resolv.conf. The
//     container's /etc/resolv.conf will then get linked to
//     /run/host/run/systemd/resolved/resolv-stub.conf.
//
//     This is why properly configured hosts should use relative symbolic
//     links, because they don't need to be adjusted in such scenarios.
//
//   - relative - containerPath will be linked to the invalid target, and it's
//     the user's responsibility to fix it up.
//
// folder signifies if the target is a folder
func redirectPath(containerPath, target string, folder bool) error {
	if !filepath.IsAbs(containerPath) {
		panic("containerPath must be an absolute path")
	}

	if !filepath.IsAbs(target) {
		panic("target must be an absolute path")
	}

	logrus.Debugf("Preparing to redirect %s to %s", containerPath, target)
	targetSanitized := sanitizeRedirectionTarget(target)

	err := os.Remove(containerPath)
	if folder {
		if err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("failed to redirect %s to %s: %w", containerPath, target, err)
		}

		if err := os.MkdirAll(target, 0755); err != nil {
			return fmt.Errorf("failed to redirect %s to %s: %w", containerPath, target, err)
		}
	}

	logrus.Debugf("Redirecting %s to %s", containerPath, targetSanitized)

	if err := os.Symlink(targetSanitized, containerPath); err != nil {
		return fmt.Errorf("failed to redirect %s to %s: %w", containerPath, target, err)
	}

	return nil
}

func runUpdateDb() {
	if err := shell.Run("updatedb", nil, nil, nil); err != nil {
		logrus.Warnf("Failed to run updatedb(8): %v", err)
	}
}

func sanitizeRedirectionTarget(target string) string {
	if !filepath.IsAbs(target) {
		panic("target must be an absolute path")
	}

	fileInfo, err := os.Lstat(target)
	if err != nil {
		if os.IsNotExist(err) {
			logrus.Warnf("%s not found", target)
		} else {
			logrus.Warnf("Failed to lstat %s: %v", target, err)
		}

		return target
	}

	fileMode := fileInfo.Mode()
	if fileMode&os.ModeSymlink == 0 {
		logrus.Debugf("%s isn't a symbolic link", target)
		return target
	}

	logrus.Debugf("%s is a symbolic link", target)

	_, err = filepath.EvalSymlinks(target)
	if err == nil {
		return target
	}

	logrus.Warnf("Failed to resolve %s: %v", target, err)

	targetDestination, err := os.Readlink(target)
	if err != nil {
		logrus.Warnf("Failed to get the destination of %s: %v", target, err)
		return target
	}

	logrus.Debugf("%s points to %s", target, targetDestination)

	if filepath.IsAbs(targetDestination) {
		logrus.Debugf("Prepending /run/host to %s", targetDestination)
		targetGuess := filepath.Join("/run/host", targetDestination)
		return targetGuess
	}

	return target
}

func extractTimeZoneFromLocalTimeSymLink(path string) (string, error) {
	zoneInfoRoots := []string{
		"/run/host/usr/share/zoneinfo",
		"/usr/share/zoneinfo",
	}

	for _, root := range zoneInfoRoots {
		if !strings.HasPrefix(path, root) {
			continue
		}

		timeZone, err := filepath.Rel(root, path)
		if err != nil {
			return "", fmt.Errorf("failed to extract time zone: %w", err)
		}

		return timeZone, nil
	}

	return "", errors.New("/etc/localtime points to unknown location")
}

func updateTimeZoneFromLocalTime() error {
	localTimeEvaled, err := filepath.EvalSymlinks("/etc/localtime")
	if err != nil {
		if os.IsNotExist(err) {
			if err := writeTimeZone("UTC"); err != nil {
				return err
			}

			return nil
		}

		return fmt.Errorf("failed to resolve /etc/localtime: %w", err)
	}

	logrus.Debugf("Resolved /etc/localtime to %s", localTimeEvaled)

	timeZone, err := extractTimeZoneFromLocalTimeSymLink(localTimeEvaled)
	if err != nil {
		return err
	}

	if err := writeTimeZone(timeZone); err != nil {
		return err
	}

	return nil
}

func writeTimeZone(timeZone string) error {
	const etcTimeZone = "/etc/timezone"

	if err := os.Remove(etcTimeZone); err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("failed to remove old %s: %w", etcTimeZone, err)
		}
	}

	timeZoneBytes := []byte(timeZone + "\n")
	if err := ioutil.WriteFile(etcTimeZone, timeZoneBytes, 0664); err != nil {
		return fmt.Errorf("failed to create new %s: %w", etcTimeZone, err)
	}

	return nil
}
