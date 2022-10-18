darwin-arm64:
	fyne-cross darwin -app-id fiurtorn.voucher -arch arm64 .

darwin-amd64:
	fyne-cross darwin -app-id fiurtorn.voucher -arch amd64 .

release-darwin-amd64:
	fyne-cross darwin -app-id fiurtorn.voucher -arch amd64 . -release

release-darwin-arm64:
	fyne-cross darwin -app-id fiurtorn.voucher -arch arm64 . -release

zip-darwin-amd64: release-darwin-amd64
	cd fyne-cross/dist/darwin-amd64 && zip -r ../../voucher-darwin-amd64.zip voucher.app

zip-darwin-arm64: release-darwin-arm64
	cd fyne-cross/dist/darwin-arm64 && zip -r ../../voucher-darwin-arm64.zip voucher.app

release: zip-darwin-amd64 zip-darwin-arm64
