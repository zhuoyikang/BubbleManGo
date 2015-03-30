b:
	cd agent/; go build; cd ../
	cd bubble/; go build; make g; cd ../
	cd main/bubble;  go run main.go

