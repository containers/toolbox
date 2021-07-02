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


#include <pthread.h>
#include <signal.h>


/*
 * pthread_attr_getstacksize < GLIBC_2.34
 */


#if defined __aarch64__
__asm__(".symver pthread_attr_getstacksize,pthread_attr_getstacksize@GLIBC_2.17");
#elif defined __arm__
__asm__(".symver pthread_attr_getstacksize,pthread_attr_getstacksize@GLIBC_2.4");
#elif defined __i386__
__asm__(".symver pthread_attr_getstacksize,pthread_attr_getstacksize@GLIBC_2.1");
#elif defined __powerpc64__ && _CALL_ELF == 2 /* ppc64le */
__asm__(".symver pthread_attr_getstacksize,pthread_attr_getstacksize@GLIBC_2.17");
#elif defined __s390x__
__asm__(".symver pthread_attr_getstacksize,pthread_attr_getstacksize@GLIBC_2.2");
#elif defined __x86_64__
__asm__(".symver pthread_attr_getstacksize,pthread_attr_getstacksize@GLIBC_2.2.5");
#else
#error "Please specify symbol version for pthread_attr_getstacksize"
#endif


int
__wrap_pthread_attr_getstacksize (const pthread_attr_t *attr, size_t *stacksize)
{
  return pthread_attr_getstacksize (attr, stacksize);
}


/*
 * pthread_create < GLIBC_2.34
 */


#if defined __aarch64__
__asm__(".symver pthread_create,pthread_create@GLIBC_2.17");
#elif defined __arm__
__asm__(".symver pthread_create,pthread_create@GLIBC_2.4");
#elif defined __i386__
__asm__(".symver pthread_create,pthread_create@GLIBC_2.1");
#elif defined __powerpc64__ && _CALL_ELF == 2 /* ppc64le */
__asm__(".symver pthread_create,pthread_create@GLIBC_2.17");
#elif defined __s390x__
__asm__(".symver pthread_create,pthread_create@GLIBC_2.2");
#elif defined __x86_64__
__asm__(".symver pthread_create,pthread_create@GLIBC_2.2.5");
#else
#error "Please specify symbol version for pthread_create"
#endif


int
__wrap_pthread_create (pthread_t *thread, const pthread_attr_t *attr, void *(*start_routine) (void *), void *arg)
{
  return pthread_create(thread, attr, start_routine, arg);
}


/*
 * pthread_detach < GLIBC_2.34
 */


#if defined __aarch64__
__asm__(".symver pthread_detach,pthread_detach@GLIBC_2.17");
#elif defined __arm__
__asm__(".symver pthread_detach,pthread_detach@GLIBC_2.4");
#elif defined __i386__
__asm__(".symver pthread_detach,pthread_detach@GLIBC_2.0");
#elif defined __powerpc64__ && _CALL_ELF == 2 /* ppc64le */
__asm__(".symver pthread_detach,pthread_detach@GLIBC_2.17");
#elif defined __s390x__
__asm__(".symver pthread_detach,pthread_detach@GLIBC_2.2");
#elif defined __x86_64__
__asm__(".symver pthread_detach,pthread_detach@GLIBC_2.2.5");
#else
#error "Please specify symbol version for pthread_detach"
#endif


int
__wrap_pthread_detach (pthread_t thread)
{
  return pthread_detach (thread);
}


/*
 * pthread_sigmask < GLIBC_2.32
 */


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
