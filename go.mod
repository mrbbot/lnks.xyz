// +heroku goVersion go1.12

module shortener

go 1.12

require (
	github.com/dimfeld/httptreemux v5.0.1+incompatible
	github.com/go-redis/redis v6.15.2+incompatible
	github.com/goincremental/negroni-sessions v0.0.0-20171223143234-40b49004abee
	github.com/gorilla/context v1.1.1 // indirect
	github.com/gorilla/sessions v1.2.0 // indirect
	github.com/joho/godotenv v1.3.0
	github.com/onsi/ginkgo v1.8.0 // indirect
	github.com/onsi/gomega v1.5.0 // indirect
	github.com/urfave/negroni v1.0.0
	golang.org/x/crypto v0.0.0-20190701094942-4def268fd1a4
)
