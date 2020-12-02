package toolbox

import (
	"errors"
	"path/filepath"

	"github.com/containers/toolbox/pkg/utils"
	"github.com/sirupsen/logrus"
)

var (
	// knownImages holds all known toolbox images
	// These should work out of the box with all Toolbox commands
	knownImages = map[string]Image{
		"fedora-toolbox": {
			"registry.fedoraproject.org",
			"fedora-toolbox",
			"latest",
		},
		"ubi7": {
			"registry.access.redhat.com",
			"ubi7",
			"latest",
		},
		"ubi7-2": {
			"registry.access.redhat.com",
			"ubi7/ubi",
			"latest",
		},
		"ubi8": {
			"registry.access.redhat.com",
			"ubi8",
			"latest",
		},
		"ubi8-2": {
			"registry.access.redhat.com",
			"ubi8/ubi",
			"latest",
		},
	}

	// supportedSystems holds ids of known systems. IDs are taken from the ID
	// entry of 'os-release'. Every id has a toolbox image attached.
	supportedSystems = map[string]Image{
		"fedora": knownImages["fedora-toolbox"],
		"rhel":   knownImages["ubi8"],
	}
)

// GetFallbackImage returns a Fedora toolbox image
func GetFallbackImage() Image {
	return supportedSystems["fedora"]
}

// GetImageForSystem returns an image matching the given system ID.
//
// System ID should be taken from the ID entry of 'os-release'.
//
// https://www.freedesktop.org/software/systemd/man/os-release.html#ID=
func GetImageForSystem(systemID string) (Image, error) {
	var img Image

	if !IsSystemSupported(systemID) {
		return img, errors.New("unsupported system")
	}

	img = supportedSystems[systemID]

	return img, nil
}

// Image holds parts of a full URI of an image
type Image struct {
	Registry   string
	Repository string
	Tag        string
}

// CreateContainerName creates a name suitable for a toolbox container based on
// image.
func (img Image) CreateContainerName() string {
	var containerName string

	logrus.Debug("Resolving container name")

	containerName = filepath.Base(img.Repository)

	if img.Tag != "" {
		containerName = containerName + "-" + img.Tag
	}

	logrus.Debugf("Resolved container name to %s", containerName)

	return containerName
}

// GetImageURI assembles full uri that can be used to access an image through
// e.g. 'podman pull'
func (img Image) GetImageURI() string {
	var imageURI string

	if img.Registry != "" {
		imageURI += img.Registry + "/"
	}
	imageURI += img.Repository
	if img.Tag != "" {
		imageURI += ":" + img.Tag
	}

	return imageURI
}

// SetImageURI sets the image to target the
func (img Image) SetImageURI(imageURI string) {
	img.Registry = utils.ImageReferenceGetDomain(imageURI)
	img.Repository = utils.ImageReferenceGetRepository(imageURI)
	img.Tag = utils.ImageReferenceGetTag(imageURI)
}

// IsImageKnown returns a bool saying whether an image is a known Toolbox
// image.
func (img Image) IsImageKnown() bool {
	for _, knownImage := range knownImages {
		// Tag does not really affect the type of image, hence is not tested
		if img.Registry == knownImage.Registry && img.Repository == knownImage.Repository {
			return true
		}
	}

	return false
}
