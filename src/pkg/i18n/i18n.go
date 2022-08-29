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
	"archive/zip"
	"bytes"
	"embed"
	"errors"
	"fmt"

	"github.com/chai2010/gettext-go"
	"github.com/sirupsen/logrus"
)

//go:embed locale
var translations embed.FS

func LoadTranslations(gettextDomain string) error {
	locale := gettext.DefaultLanguage
	logrus.Debugf("Setting language to %s", locale)
	gettext.SetLanguage(locale)

	translationFiles := []string{
		fmt.Sprintf("%s/%s/LC_MESSAGES/%s.po", gettextDomain, locale, gettextDomain),
		fmt.Sprintf("%s/%s/LC_MESSAGES/%s.mo", gettextDomain, locale, gettextDomain),
	}

	buf := new(bytes.Buffer)
	w := zip.NewWriter(buf)

	for _, file := range translationFiles {
		filename := "locale/" + file
		f, err := w.Create(file)
		if err != nil {
			return err
		}
		data, err := translations.ReadFile(filename)
		if err != nil {
			return err
		}
		if _, err := f.Write(data); err != nil {
			return nil
		}
	}
	if err := w.Close(); err != nil {
		return err
	}

	gettext.BindLocale(gettext.New(gettextDomain, gettextDomain+".zip", buf.Bytes()))
	return nil
}

func T(defaultValue string, args ...int) string {
	if len(args) == 0 {
		return gettext.Gettext(defaultValue)
	}
	return fmt.Sprintf(gettext.NGettext(defaultValue, defaultValue+".plural", args[0]),
		args[0])
}

func Errorf(defaultValue string, args ...int) error {
	return errors.New(T(defaultValue, args...))
}
