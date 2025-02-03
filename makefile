VERSION = 0.0.7

upload:
	git tag -a v$(VERSION) -m "Version $(VERSION)"
	git push origin v$(VERSION)

.PHONY: upload
