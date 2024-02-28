/*
 * Copyright © 2022 – 2024 Red Hat Inc.
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

package utils

import (
	"errors"
	"fmt"
	"os/user"
	"unsafe"
)

/*
#cgo LDFLAGS: -ldl

#include <stdlib.h>

#include <dlfcn.h>
#include <shadow/subid.h>

#include "libsubid-wrappers.h"
*/
import "C"

func validateSubIDRange(user *user.User, libsubid unsafe.Pointer, cSubidGetIDRangesSymbol *C.char) (bool, error) {
	subid_get_id_ranges := C.dlsym(libsubid, cSubidGetIDRangesSymbol)
	if subid_get_id_ranges == nil {
		subidGetIDRangesSymbol := C.GoString(cSubidGetIDRangesSymbol)
		return false, fmt.Errorf("cannot dlsym(3) %s", subidGetIDRangesSymbol)
	}

	cUsername := C.CString(user.Username)
	defer C.free(unsafe.Pointer(cUsername))

	var cRanges *C.struct_subid_range
	defer C.free(unsafe.Pointer(cRanges))

	nRanges := C.toolbox_subid_get_id_ranges(subid_get_id_ranges, cUsername, &cRanges)
	if nRanges <= 0 {
		cUid := C.CString(user.Uid)
		defer C.free(unsafe.Pointer(cUid))

		nRanges = C.toolbox_subid_get_id_ranges(subid_get_id_ranges, cUid, &cRanges)
	}

	if nRanges <= 0 {
		return false, errors.New("cannot read subids")
	}

	return true, nil
}

func ValidateSubIDRanges(user *user.User) (bool, error) {
	if IsInsideContainer() {
		panic("cannot validate subordinate IDs inside container")
	}

	if user == nil {
		panic("cannot validate subordinate IDs when user is nil")
	}

	if user.Username == "ALL" {
		return false, errors.New("username ALL not supported")
	}

	libsubid := C.dlopen(C.TOOLBOX_LIBSUBID, C.RTLD_LAZY)
	if libsubid == nil {
		filename := C.GoString(C.TOOLBOX_LIBSUBID)
		return false, fmt.Errorf("cannot dlopen(3) %s", filename)
	}

	defer C.dlclose(libsubid)

	subid_init := C.dlsym(libsubid, C.TOOLBOX_SUBID_INIT)
	if subid_init == nil {
		subidInitSymbol := C.GoString(C.TOOLBOX_SUBID_INIT)
		return false, fmt.Errorf("cannot dlsym(3) %s", subidInitSymbol)
	}

	C.toolbox_subid_init(subid_init)

	if _, err := validateSubIDRange(user, libsubid, C.TOOLBOX_SUBID_GET_GID_RANGES_SYMBOL); err != nil {
		return false, err
	}

	if _, err := validateSubIDRange(user, libsubid, C.TOOLBOX_SUBID_GET_UID_RANGES_SYMBOL); err != nil {
		return false, err
	}

	return true, nil
}
