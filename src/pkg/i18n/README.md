# Translations README

### General workflow should be:
- Wrap the strings for translation
- Create new (updated) template file
- Add a new language to toolbox
- Update `PO` files with latest `POT`
- Translate `PO` files, check statistics
- Compile `PO` files to `MO` binary

# Manage Translations

## Wrap the strings for translation

To mark the strings for translation, simply wrap them:

```go
import "github.com/containers/toolbox/pkg/i18n"
...
// a string
translated := i18n.T("Some english message here")

// plural string
translated := i18n.T("You had % values", values)

// an error
err := i18n.Error("Something unexpected happened")

// plural error
err := i18n.Error("%d unexpected events")
```

## Create new (updated) template file

Everytime new strings are marked for translation, we need to extract them to a `POT` file.

```console
make potfile
```
This should create latest `POT` at `src/pkg/i18n/locale`.

## Add a new language

To add a new language, use following command. Replace `en_US` with your locale `xx_YY`.

```console
make addlang locale=en_US
```
This will create required folders. And the `PO` from the `POT` file.

## Update PO files with latest POT

If we have new strings in the `POT` file, they should be merged with `PO` files.

```console
make updatepo
```
This will merge the latest `POT` with all existing `PO` files.

## Get the PO files translated

Edit the appropriate `toolbox.po` file, `poedit` is a popular open source tool
for translations.

## Get translation statistics

```console
make postat
```
This will show statistics of all the `PO` files.

## Compile PO files to MO binary

Once translations are done with `toolbox.po` file, generate the corresponding `toolbox.mo` file.

```console
make po2mo
```
This will compile all the `PO` files to `MO`.
