/*
 * Copyright © 2020 – 2021 Red Hat Inc.
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


#include <signal.h>


#if defined __aarch64__
__asm__(".symver pthread_sigmask,pthread_sigmask@GLIBC_2.17");
#elif defined __arm__
__asm__(".symver pthread_sigmask,pthread_sigmask@GLIBC_2.4");
#elif defined __i386__
__asm__(".symver pthread_sigmask,pthread_sigmask@GLIBC_2.0");
#elif defined __powerpc64__ && _CALL_ELF == 2 /* ppc64le */
__asm__(".symver pthread_sigmask,pthread_sigmask@GLIBC_2.17");
#elif defined __s390x__
__asm__(".symver pthread_sigmask,pthread_sigmask@GLIBC_2.2");
#elif defined __x86_64__
__asm__(".symver pthread_sigmask,pthread_sigmask@GLIBC_2.2.5");
#else
#error "Please specify symbol version for pthread_sigmask"
#endif


int
__wrap_pthread_sigmask (int how, const sigset_t *set, sigset_t *oldset)
{
  return pthread_sigmask (how, set, oldset);
}
