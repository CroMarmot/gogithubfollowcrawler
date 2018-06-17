run: main.go JsonSaver.go
	go run main.go JsonSaver.go

test: JsonSaver.go JsonSaver_test.go
	go test JsonSaver_test.go JsonSaver.go
