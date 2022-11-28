/*
 * Copyright © 2019 – 2022 Red Hat Inc.
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

/*
#cgo LDFLAGS: -l subid
#include <shadow/subid.h>
#include <stdlib.h>
#include <stdio.h>
const char *Prog = "storage";
FILE *shadow_logfd = NULL;
struct subid_range get_range(struct subid_range *ranges, int i)
{
	shadow_logfd = stderr;
	return ranges[i];
}
#if !defined(SUBID_ABI_MAJOR) || (SUBID_ABI_MAJOR < 4)
# define subid_get_uid_ranges get_subuid_ranges
# define subid_get_gid_ranges get_subgid_ranges
#endif
*/
import "C"

import (
	"errors"
	"os/user"
	"unsafe"

	"github.com/sirupsen/logrus"
)

type subIDRange struct {
	Start  int
	Length int
}

type ranges []subIDRange

func validateSubIDRange(username string, isUser bool) (ranges, error) {
	var ret ranges

	queryType := ""

	if isUser {
		queryType = "subuid"
	} else {
		queryType = "subgid"
	}

	uidstr := ""

	if username == "ALL" {
		return nil, errors.New("username ALL not supported")
	}

	if u, err := user.Lookup(username); err == nil {
		uidstr = u.Uid
	}

	cUsername := C.CString(username)
	defer C.free(unsafe.Pointer(cUsername))

	cuidstr := C.CString(uidstr)
	defer C.free(unsafe.Pointer(cuidstr))

	var nRanges C.int
	var cRanges *C.struct_subid_range
	if isUser {
		nRanges = C.subid_get_uid_ranges(cUsername, &cRanges)
		if nRanges <= 0 {
			nRanges = C.subid_get_uid_ranges(cuidstr, &cRanges)
		}
	} else {
		nRanges = C.subid_get_gid_ranges(cUsername, &cRanges)
		if nRanges <= 0 {
			nRanges = C.subid_get_gid_ranges(cuidstr, &cRanges)
		}
	}
	if nRanges <= 0 {
		return nil, errors.New("cannot read subids")
	}
	defer C.free(unsafe.Pointer(cRanges))

	for i := 0; i < int(nRanges); i++ {
		r := C.get_range(cRanges, C.int(i))
		newRange := subIDRange{
			Start:  int(r.start),
			Length: int(r.count),
		}
		ret = append(ret, newRange)

		logrus.Debugf("Found %s range %d for %s: start [%d] length [%d]",
			queryType,
			i,
			username,
			r.start,
			r.count)
	}
	return ret, nil
}
