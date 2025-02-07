all: project
clean:
	rm -f package *.o project *.o hello.pdf image.pdf temp.jpg
project: package.go
	go build -o project package.go
run: 
	@echo "Running with arguments: $(args)"
	./project $(args)
