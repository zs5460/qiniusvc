module github.com/zs5460/qiniusvc

go 1.15

replace (
	golang.org/x/crypto => github.com/golang/crypto v0.0.0-20190426145343-a29dc8fdc734
	golang.org/x/net => github.com/golang/net v0.0.0-20190424112056-4829fb13d2c6
	golang.org/x/sync => github.com/golang/sync v0.0.0-20190423024810-112230192c58
	golang.org/x/sys => github.com/golang/sys v0.0.0-20190429190828-d89cdac9e872
	golang.org/x/text => github.com/golang/text v0.3.2
	golang.org/x/tools => github.com/golang/tools v0.0.0-20190430004104-b9fed7929fc1

)

require (
	github.com/baa-middleware/recovery v0.0.0-20200227085107-3da4ea0df8b2
	github.com/baa-middleware/static v0.0.0-20200227085341-c941b03b1006
	github.com/go-baa/baa v1.2.32
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/json-iterator/go v1.1.10 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/qiniu/api.v7/v7 v7.5.0
	github.com/smartystreets/goconvey v1.6.4 // indirect
	github.com/stretchr/testify v1.6.1 // indirect
	github.com/timest/env v0.0.0-20180717050204-5fce78d35255
)
