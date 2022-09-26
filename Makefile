build:
	@echo "Building Current Time App"
	env GOOS=linux GOARCH=amd64
	go build -o curr-time .
	chmod 777 curr-time
	@echo "Done!"

imgbuild:
	docker build -t lnanjangud653/curr-time:"${IMG_VER}" .

imgpush:
	docker push lnanjangud653/curr-time:${IMG_VER}
