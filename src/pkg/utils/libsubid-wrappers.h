/*
 * Copyright © 2023 Red Hat Inc.
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

#pragma once

#include <shadow/subid.h>

extern const char *TOOLBOX_LIBSUBID;

extern const char *TOOLBOX_LIBSUBID_INIT;
extern const char *TOOLBOX_SUBID_INIT;

extern const char *TOOLBOX_SUBID_GET_GID_RANGES_SYMBOL;
extern const char *TOOLBOX_SUBID_GET_UID_RANGES_SYMBOL;

void  toolbox_subid_init           (void *subid_init_func);

int   toolbox_subid_get_id_ranges  (void *subid_get_id_ranges_func, const char *owner, struct subid_range **ranges);
