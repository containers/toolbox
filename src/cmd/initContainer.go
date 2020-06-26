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
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"strings"
	"syscall"

	"github.com/containers/toolbox/pkg/shell"
	"github.com/containers/toolbox/pkg/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	initContainerFlags struct {
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
		{"/var/lib/flatpak", "/run/host/var/lib/flatpak", "ro"},
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

	flags.StringVar(&initContainerFlags.home,
		"home",
		"",
		"Create a user inside the toolbox container whose login directory is HOME.")
	initContainerCmd.MarkFlagRequired("home")

	flags.BoolVar(&initContainerFlags.homeLink,
		"home-link",
		false,
		"Make /home a symbolic link to /var/home.")

	flags.BoolVar(&initContainerFlags.mediaLink,
		"media-link",
		false,
		"Make /media a symbolic link to /run/media.")

	flags.BoolVar(&initContainerFlags.mntLink, "mnt-link", false, "Make /mnt a symbolic link to /var/mnt.")

	flags.BoolVar(&initContainerFlags.monitorHost,
		"monitor-host",
		false,
		"Ensure that certain configuration files inside the toolbox container are in sync with the host.")

	flags.StringVar(&initContainerFlags.shell,
		"shell",
		"",
		"Create a user inside the toolbox container whose login shell is SHELL.")
	initContainerCmd.MarkFlagRequired("shell")

	flags.IntVar(&initContainerFlags.uid,
		"uid",
		0,
		"Create a user inside the toolbox container whose numerical user ID is UID.")
	initContainerCmd.MarkFlagRequired("uid")

	flags.StringVar(&initContainerFlags.user,
		"user",
		"",
		"Create a user inside the toolbox container whose login name is USER.")
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

	runtimeDirectory := os.Getenv("XDG_RUNTIME_DIR")
	if runtimeDirectory == "" {
		logrus.Debug("XDG_RUNTIME_DIR is unset")

		runtimeDirectory = fmt.Sprintf("/run/user/%d", initContainerFlags.uid)
		os.Setenv("XDG_RUNTIME_DIR", runtimeDirectory)

		logrus.Debugf("XDG_RUNTIME_DIR set to %s", runtimeDirectory)
	}

	logrus.Debug("Creating /run/.toolboxenv")

	toolboxEnvFile, err := os.Create("/run/.toolboxenv")
	if err != nil {
		return fmt.Errorf("failed to create /run/.toolboxenv: %w", err)
	}

	defer toolboxEnvFile.Close()

	logrus.Debug("Mounting tmpfs at /tmp")

	if err := syscall.Mount("tmpfs",
		"/tmp",
		"tmpfs",
		syscall.MS_NODEV|syscall.MS_STRICTATIME|syscall.MS_NOSUID,
		"mode=1777"); err != nil {
		return fmt.Errorf("failed to mount tmpfs at /tmp: %s", err)
	}

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

		if utils.PathExists("/run/host/monitor") {
			logrus.Debug("Path /run/host/monitor exists")

			if localtimeTarget, err := os.Readlink("/etc/localtime"); err != nil ||
				localtimeTarget != "/run/host/monitor/localtime" {
				if err := redirectPath("/etc/localtime",
					"/run/host/monitor/localtime", false); err != nil {
					return err
				}
			}

			if _, err := os.Readlink("/etc/timezone"); err != nil {
				if err := redirectPath("/etc/timezone",
					"/run/host/monitor/timezone",
					false); err != nil {
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
		if initContainerFlags.homeLink {
			if err := redirectPath("/home", "/var/home", true); err != nil {
				return err
			}
		}

		sudoGroup, err := utils.GetGroupForSudo()
		if err != nil {
			return fmt.Errorf("failed to add user %s: %w", initContainerFlags.user, err)
		}

		logrus.Debugf("Adding user %s with UID %d:", initContainerFlags.user, initContainerFlags.uid)

		useraddArgs := []string{
			"--home-dir", initContainerFlags.home,
			"--no-create-home",
			"--shell", initContainerFlags.shell,
			"--uid", fmt.Sprint(initContainerFlags.uid),
			"--groups", sudoGroup,
			initContainerFlags.user,
		}

		logrus.Debug("useradd")
		for _, arg := range useraddArgs {
			logrus.Debugf("%s", arg)
		}

		if err := shell.Run("useradd", nil, nil, nil, useraddArgs...); err != nil {
			return fmt.Errorf("failed to add user %s with UID %d",
				initContainerFlags.user,
				initContainerFlags.uid)
		}

		logrus.Debugf("Removing password for user %s", initContainerFlags.user)

		if err := shell.Run("passwd", nil, nil, nil, "--delete", initContainerFlags.user); err != nil {
			return fmt.Errorf("failed to remove password for user %s: %w", initContainerFlags.user, err)
		}

		logrus.Debug("Removing password for user root")

		if err := shell.Run("passwd", nil, nil, nil, "--delete", "root"); err != nil {
			return fmt.Errorf("failed to remove password for root: %w", err)
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
			return fmt.Errorf("failed to set KCM as the defult Kerberos credential cache: %w", err)
		}
	}

	logrus.Debug("Finished initializing container")

	toolboxRuntimeDirectory := runtimeDirectory + "/toolbox"
	logrus.Debugf("Creating runtime directory %s", toolboxRuntimeDirectory)

	if err := os.MkdirAll(toolboxRuntimeDirectory, 0700); err != nil {
		return fmt.Errorf("failed to create runtime directory %s: %w", toolboxRuntimeDirectory, err)
	}

	if err := os.Chown(toolboxRuntimeDirectory, initContainerFlags.uid, initContainerFlags.uid); err != nil {
		return fmt.Errorf("failed to change ownership of the runtime directory %s: %w",
			toolboxRuntimeDirectory, err)
	}

	pid := os.Getpid()
	initializedStamp := fmt.Sprintf("%s/container-initialized-%d", toolboxRuntimeDirectory, pid)

	logrus.Debugf("Creating initialization stamp %s", initializedStamp)

	initializedStampFile, err := os.Create(initializedStamp)
	if err != nil {
		return fmt.Errorf("failed to create initialization stamp: %w", err)
	}

	defer initializedStampFile.Close()

	if err := initializedStampFile.Chown(initContainerFlags.uid, initContainerFlags.uid); err != nil {
		return fmt.Errorf("failed to change ownership of initialization stamp: %w", err)
	}

	logrus.Debug("Going to sleep")

	sleepBinary, err := exec.LookPath("sleep")
	if err != nil {
		if errors.Is(err, exec.ErrNotFound) {
			return errors.New("sleep(1) not found")
		}

		return fmt.Errorf("failed to lookup sleep(1): %w", err)
	}

	sleepArgs := []string{"sleep", "+Inf"}
	env := os.Environ()

	if err := syscall.Exec(sleepBinary, sleepArgs, env); err != nil {
		return fmt.Errorf("failed to invoke sleep(1): %w", err)
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

	if err := utils.ShowManual("toolbox-init-container"); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		return
	}
}

func mountBind(containerPath, source, flags string) error {
	fi, err := os.Stat(source)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}

		return fmt.Errorf("failed to stat %s: %w", source, err)
	}

	if fi.IsDir() {
		logrus.Debugf("Creating %s", containerPath)
		if err := os.MkdirAll(containerPath, 0755); err != nil {
			return fmt.Errorf("failed to create %s: %w", containerPath, err)
		}
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
		return fmt.Errorf("failed to bind %s to %s: %w", containerPath, source, err)
	}

	return nil
}

func redirectPath(containerPath, target string, folder bool) error {
	logrus.Debugf("Redirecting %s to %s", containerPath, target)

	err := os.Remove(containerPath)
	if folder {
		if err != nil {
			return fmt.Errorf("failed to redirect %s to %s: %w", containerPath, target, err)
		}

		if err := os.MkdirAll(target, 0755); err != nil {
			return fmt.Errorf("failed to redirect %s to %s: %w", containerPath, target, err)
		}
	}

	if err := os.Symlink(target, containerPath); err != nil {
		return fmt.Errorf("failed to redirect %s to %s: %w", containerPath, target, err)
	}

	return nil
}
