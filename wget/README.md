go-tests / wget
========

Port of wget unix command

	wget https://github.com/friedrichg/go-tests/raw/master/wget/wget.go
	go build wget.go
	./wget http://www.google.com

And perhaps cross-compile ;) 

	GOOS=windows GOARCH=386 go build -o wget.exe wget.go 
	
Look & run make.bash in your go src dir to make this last thing work
