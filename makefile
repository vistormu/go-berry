VERSION = 0.0.8

upload:
	git tag -a v$(VERSION) -m "Version $(VERSION)"
	git push origin v$(VERSION)

.PHONY: upload
