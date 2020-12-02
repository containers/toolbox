package toolbox

import (
	"github.com/containers/toolbox/pkg/utils"
	"github.com/sirupsen/logrus"
)

// IsSystemSupported takes a string identifying an operating system and
// compares it to an internal list of supported systems.
//
// The system ids are taken from the ID entry of 'os-release'.
//
// https://www.freedesktop.org/software/systemd/man/os-release.html#ID=
func IsSystemSupported(systemID string) bool {
	if _, ok := supportedSystems[systemID]; ok {
		return true
	}

	return false
}

// IsHostSystemSupported checks if the systemm where Toolbox is run, has
// a supported toolbox image (e.g. Fedora has fedora-toolbox).
//
// The compared information are: the system ID and an internal list of
// system IDs where ID is taken from 'os-release'.
//
// https://www.freedesktop.org/software/systemd/man/os-release.html#ID=
func IsHostSystemSupported() bool {
	hostID, err := utils.GetHostID()
	if err != nil {
		logrus.Warnf("There was an error while getting host's ID: %v", err)
		return false
	}

	if _, ok := supportedSystems[hostID]; ok {
		return true
	}

	return false
}
