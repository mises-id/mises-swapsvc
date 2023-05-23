#mises_alpha
truss:
	truss proto/swapsvc.proto  --pbpkg github.com/mises-id/mises-swapsvc/proto --svcpkg github.com/mises-id/mises-swapsvc --svcout . -v 
run:
	go run cmd/cli/main.go
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build ./cmd/cli/main.go
upload:
	scp ./main mises_alpha:/apps/mises-swapsvc/
replace:
	ssh mises_alpha "mv /apps/mises-swapsvc/main /apps/mises-swapsvc/mises-swapsvc"
restart:
	ssh mises_alpha "sudo supervisorctl restart mises-swapsvc"
deploy: build \
	upload \
	replace \
	restart
#mises_backup
upload-backup:
	scp ./main mises_backup:/apps/mises-swapsvc/
replace-backup:
	ssh mises_backup "mv /apps/mises-swapsvc/main /apps/mises-swapsvc/mises-swapsvc"
restart-backup:
	ssh mises_backup "sudo supervisorctl restart mises-swapsvc"
deploy-backup: build \
	upload-backup \
	replace-backup \
	restart-backup 