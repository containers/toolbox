/*
 * Copyright © 2019 – 2021 Red Hat Inc.
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
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/containers/toolbox/pkg/shell"
	"github.com/containers/toolbox/pkg/utils"
	"github.com/fsnotify/fsnotify"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
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
		{"/etc/machine-id", "/run/host/etc/machine-id", "ro"},
		{"/run/libvirt", "/run/host/run/libvirt", ""},
		{"/run/systemd/journal", "/run/host/run/systemd/journal", ""},
		{"/run/systemd/resolve", "/run/host/run/systemd/resolve", ""},
		{"/run/udev/data", "/run/host/run/udev/data", ""},
		{"/tmp", "/run/host/tmp", "rslave"},
		{"/var/lib/flatpak", "/run/host/var/lib/flatpak", "ro"},
		{"/var/lib/libvirt", "/run/host/var/lib/libvirt", ""},
		{"/var/lib/systemd/coredump", "/run/host/var/lib/systemd/coredump", "ro"},
		{"/var/log/journal", "/run/host/var/log/journal", "ro"},
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
		"Create a user inside the toolbox container whose numerical group ID is GID")

	flags.StringVar(&initContainerFlags.home,
		"home",
		"",
		"Create a user inside the toolbox container whose login directory is HOME")
	initContainerCmd.MarkFlagRequired("home")

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
		false,
		"Ensure that certain configuration files inside the toolbox container are in sync with the host")

	flags.StringVar(&initContainerFlags.shell,
		"shell",
		"",
		"Create a user inside the toolbox container whose login shell is SHELL")
	initContainerCmd.MarkFlagRequired("shell")

	flags.IntVar(&initContainerFlags.uid,
		"uid",
		0,
		"Create a user inside the toolbox container whose numerical user ID is UID")
	initContainerCmd.MarkFlagRequired("uid")

	flags.StringVar(&initContainerFlags.user,
		"user",
		"",
		"Create a user inside the toolbox container whose login name is USER")
	initContainerCmd.MarkFlagRequired("user")

	initContainerCmd.SetHelpFunc(initContainerHelp)
	rootCmd.AddCommand(initContainerCmd)
}

func initContainer(cmd *cobra.Command, args []string) error {
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

	if initContainerFlags.monitorHost {
		logrus.Debug("Monitoring host")

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

			if _, err := os.Readlink("/etc/resolv.conf"); err != nil {
				if err := redirectPath("/etc/resolv.conf",
					"/run/host/etc/resolv.conf",
					false); err != nil {
					return err
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

	if _, err := user.Lookup(initContainerFlags.user); err != nil {
		if err := configureUsers(initContainerFlags.uid,
			initContainerFlags.user,
			initContainerFlags.home,
			initContainerFlags.shell,
			initContainerFlags.homeLink,
			false); err != nil {
			return err
		}
	} else {
		if err := configureUsers(initContainerFlags.uid,
			initContainerFlags.user,
			initContainerFlags.home,
			initContainerFlags.shell,
			initContainerFlags.homeLink,
			true); err != nil {
			return err
		}
	}

	if utils.PathExists("/etc/krb5.conf.d") && !utils.PathExists("/etc/krb5.conf.d/kcm_default_ccache") {
		logrus.Debug("Setting KCM as the default Kerberos credential cache")

		kcmConfigString := `# Written by Toolbox
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
			return errors.New("failed to set KCM as the defult Kerberos credential cache")
		}
	}

	logrus.Debug("Setting up daily ticker")

	daily, err := time.ParseDuration("24h")
	if err != nil {
		panicMsg := fmt.Sprintf("failed to parse duration: %v", err)
		panic(panicMsg)
	}

	tickerDaily := time.NewTicker(daily)
	defer tickerDaily.Stop()

	logrus.Debug("Setting up watches for file system events")

	watcherForHost, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	defer watcherForHost.Close()

	if err := watcherForHost.Add("/run/host/etc"); err != nil {
		return err
	}

	logrus.Debug("Finished initializing container")

	uidString := strconv.Itoa(initContainerFlags.uid)
	targetUser, err := user.LookupId(uidString)
	if err != nil {
		return fmt.Errorf("failed to lookup user ID %s: %w", uidString, err)
	}

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
		case event := <-watcherForHost.Events:
			handleFileSystemEvent(event)
		case err := <-watcherForHost.Errors:
			logrus.Warnf("Received an error from the file system watcher: %v", err)
		}
	}

	return nil
}

func initContainerHelp(cmd *cobra.Command, args []string) {
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

	if err := showManual("toolbox-init-container"); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		return
	}
}

func configureUsers(targetUserUid int,
	targetUser, targetUserHome, targetUserShell string,
	homeLink, targetUserExists bool) error {
	if homeLink {
		if err := redirectPath("/home", "/var/home", true); err != nil {
			return err
		}
	}

	sudoGroup, err := utils.GetGroupForSudo()
	if err != nil {
		return fmt.Errorf("failed to get group for sudo: %w", err)
	}

	if targetUserExists {
		logrus.Debugf("Modifying user %s with UID %d:", targetUser, targetUserUid)

		usermodArgs := []string{
			"--append",
			"--groups", sudoGroup,
			"--home", targetUserHome,
			"--shell", targetUserShell,
			"--uid", fmt.Sprint(targetUserUid),
			targetUser,
		}

		logrus.Debug("usermod")
		for _, arg := range usermodArgs {
			logrus.Debugf("%s", arg)
		}

		if err := shell.Run("usermod", nil, nil, nil, usermodArgs...); err != nil {
			return fmt.Errorf("failed to modify user %s with UID %d", targetUser, targetUserUid)
		}
	} else {
		logrus.Debugf("Adding user %s with UID %d:", targetUser, targetUserUid)

		useraddArgs := []string{
			"--groups", sudoGroup,
			"--home-dir", targetUserHome,
			"--no-create-home",
			"--shell", targetUserShell,
			"--uid", fmt.Sprint(targetUserUid),
			targetUser,
		}

		logrus.Debug("useradd")
		for _, arg := range useraddArgs {
			logrus.Debugf("%s", arg)
		}

		if err := shell.Run("useradd", nil, nil, nil, useraddArgs...); err != nil {
			return fmt.Errorf("failed to add user %s with UID %d", targetUser, targetUserUid)
		}
	}

	logrus.Debugf("Removing password for user %s", targetUser)

	if err := shell.Run("passwd", nil, nil, nil, "--delete", targetUser); err != nil {
		return fmt.Errorf("failed to remove password for user %s", targetUser)
	}

	logrus.Debug("Removing password for user root")

	if err := shell.Run("passwd", nil, nil, nil, "--delete", "root"); err != nil {
		return errors.New("failed to remove password for root")
	}

	return nil
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
	} else if fileMode.IsRegular() {
		logrus.Debugf("Creating regular file %s", containerPath)

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
// configuration files to their counterparts on the host's filesystem.
//
// containerPath and target must be absolute paths
//
// If the target itself is a symbolic link, redirectPath will evaluate the
// link. If it's valid then redirectPath will link containerPath to target.
// If it's not, then redirectPath will still proceed with the linking in two
// different ways depending whether target is an absolute or a relative link:
//
//   * absolute - target's destination will be prepended with /run/host, and
//                containerPath will be linked to the resulting path. This is
//                an attempt to unbreak things, but if it doesn't work then
//                it's the user's responsibility to fix it up.
//
//                This is meant to address the common case where
//                /etc/resolv.conf on the host (ie., /run/host/etc/resolv.conf
//                inside the container) is an absolute symbolic link to
//                something like /run/systemd/resolve/stub-resolv.conf. The
//                container's /etc/resolv.conf will then get linked to
//                /run/host/run/systemd/resolved/resolv-stub.conf.
//
//                This is why properly configured hosts should use relative
//                symbolic links, because they don't need to be adjusted in
//                such scenarios.
//
//   * relative - containerPath will be linked to the invalid target, and it's
//                the user's responsibility to fix it up.
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
