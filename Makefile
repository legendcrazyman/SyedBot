all: syed

syed:
	go build -o syed syed.go

clean:
	rm -f syed