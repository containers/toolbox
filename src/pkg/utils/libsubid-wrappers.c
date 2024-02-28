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


#include <stdbool.h>
#include <stdio.h>

#include "libsubid-wrappers.h"

#ifndef SUBID_ABI_VERSION
#define SUBID_ABI_VERSION 3.0.0
#endif

#if SUBID_ABI_MAJOR < 4
#define subid_init libsubid_init
#define subid_get_gid_ranges get_subgid_ranges
#define subid_get_uid_ranges get_subuid_ranges
#endif

#define TOOLBOX_STRINGIZE_HELPER(s) #s
#define TOOLBOX_STRINGIZE(s) TOOLBOX_STRINGIZE_HELPER (s)


typedef bool (*ToolboxSubidInitFunc) (const char *progname, FILE *logfd);
typedef int (*ToolboxSubidGetRangesFunc) (const char *owner, struct subid_range **ranges);

const char *TOOLBOX_LIBSUBID = "libsubid.so." TOOLBOX_STRINGIZE (SUBID_ABI_VERSION);

const char *TOOLBOX_SUBID_INIT = TOOLBOX_STRINGIZE (subid_init);

const char *TOOLBOX_SUBID_GET_GID_RANGES_SYMBOL = TOOLBOX_STRINGIZE (subid_get_gid_ranges);
const char *TOOLBOX_SUBID_GET_UID_RANGES_SYMBOL = TOOLBOX_STRINGIZE (subid_get_uid_ranges);


void
toolbox_subid_init (void *subid_init_func)
{
  (* (ToolboxSubidInitFunc) subid_init_func) (NULL, stderr);
}


int
toolbox_subid_get_id_ranges (void *subid_get_id_ranges_func, const char *owner, struct subid_range **ranges)
{
  int ret_val = 0;

  ret_val = (* (ToolboxSubidGetRangesFunc) subid_get_id_ranges_func) (owner, ranges);
  return ret_val;
}
