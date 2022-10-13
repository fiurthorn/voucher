darwin-arm64:
	fyne-cross darwin -app-id fiurtorn.voucher -arch arm64 cmd/voucher.go

darwin-amd64:
	fyne-cross darwin -app-id fiurtorn.voucher -arch amd64 cmd/voucher.go

release-darwin-amd64:
	fyne-cross darwin -app-id fiurtorn.voucher -arch amd64 cmd/voucher.go -release

release-darwin-arm64:
	fyne-cross darwin -app-id fiurtorn.voucher -arch arm64 cmd/voucher.go -release