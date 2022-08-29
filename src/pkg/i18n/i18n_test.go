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
package i18n

import (
	"os"
	"testing"
)

var testLocale = "en_US"

func TestTranslation(t *testing.T) {
	os.Setenv("LANG", testLocale)
	if err := LoadTranslations("test"); err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	result := T("test message string")
	if result != "foobar" {
		t.Errorf("expected: %s, saw: %s", "foobar", result)
	}
}

func TestTranslationPlural(t *testing.T) {
	os.Setenv("LANG", testLocale)
	if err := LoadTranslations("test"); err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	result := T("test plural message string", 3)
	if result != "there were 3 cars behind" {
		t.Errorf("expected: %s, saw: %s", "there were 3 cars behind", result)
	}

	result = T("test plural message string", 1)
	if result != "there was 1 car behind" {
		t.Errorf("expected: %s, saw: %s", "there was 1 car behind", result)
	}
}
