BASEDIR  := ./../../../src
TRANSDIR := pkg/i18n/locale
TRANSPATH:= $(BASEDIR)/$(TRANSDIR)

GETTEXT_DOMAIN=toolbox
POT2PO_BIN=/usr/bin/msginit
PO2MO_BIN=/usr/bin/msgfmt
TRANS_UPDATE_BIN=/usr/bin/msgmerge
TRANS_EXTRACT_BIN=/usr/bin/go-xgettext

.PHONY: test
test:
	cd $(BASEDIR) && ls

.PHONY: potfile
potfile:
	cd $(BASEDIR) && $(TRANS_EXTRACT_BIN) -k=i18n.T cmd/* > $(TRANSDIR)/$(GETTEXT_DOMAIN).pot

.PHONY: addlang
addlang:
	@[ "$(locale)" ] || ( echo "Set locale as locale=en_US"; exit 1 )
	cd $(BASEDIR) && mkdir -p $(TRANSDIR)/$(GETTEXT_DOMAIN)/$(locale)/LC_MESSAGES
	cd $(BASEDIR) && $(POT2PO_BIN) -i $(TRANSDIR)/$(GETTEXT_DOMAIN).pot -o $(TRANSDIR)/$(GETTEXT_DOMAIN)/$(locale)/LC_MESSAGES/$(GETTEXT_DOMAIN).po -l $(locale).UTF-8

.PHONY: po2mo
po2mo:
	@for lang in $(shell ls $(TRANSPATH)/$(GETTEXT_DOMAIN));\
	do\
		$(PO2MO_BIN) -o $(TRANSPATH)/$(GETTEXT_DOMAIN)/$${lang}/LC_MESSAGES/$(GETTEXT_DOMAIN).mo $(TRANSPATH)/$(GETTEXT_DOMAIN)/$${lang}/LC_MESSAGES/$(GETTEXT_DOMAIN).po;\
		echo "po2mo for $${lang} done.";\
	done

.PHONY: updatepo
updatepo:
	@for lang in $(shell ls $(TRANSPATH)/$(GETTEXT_DOMAIN));\
	do\
		echo "updatepo for $${lang}";\
		$(TRANS_UPDATE_BIN) -U $(TRANSPATH)/$(GETTEXT_DOMAIN)/$${lang}/LC_MESSAGES/$(GETTEXT_DOMAIN).po $(TRANSPATH)/$(GETTEXT_DOMAIN).pot;\
	done

.PHONY: postat
postat:
	@for lang in $(shell ls $(TRANSPATH)/$(GETTEXT_DOMAIN));\
	do\
		echo "stats for $${lang}";\
		$(PO2MO_BIN) --statistics $(TRANSPATH)/$(GETTEXT_DOMAIN)/$${lang}/LC_MESSAGES/$(GETTEXT_DOMAIN).po -o /dev/null;\
	done
