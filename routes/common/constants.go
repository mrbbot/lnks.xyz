package common

import "time"

const (
	PathDashboard = "/"
	PathLink      = "/api/link"
	PathLogin     = "/login"
	PathLogout    = "/logout"
	PathRegister  = "/register"
)

const (
	RedisRegistrationCodeNamespace = "code:"
	RedisHostsNamespace            = "hosts:"
	RedisLinkNamespace             = "link:"
	RedisPasswordNamespace         = "password:"
	RedisUserLinksNamespace        = "userlinks:"
)

const (
	LastClickLayout = "15:04 on 02/01/06"
	CreatedLayout   = "02/01/06"
)

var (
	LocationLondon *time.Location
)

func init() {
	var err error
	LocationLondon, err = time.LoadLocation("Europe/London")
	if err != nil {
		panic(err)
	}
}
