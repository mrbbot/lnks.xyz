build-sass:
	sass -s compressed styles/styles.sass public/css/styles.css
watch-sass:
	sass --watch styles/styles.sass public/css/styles.css
build: build-sass
	go build main.go
