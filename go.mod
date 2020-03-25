module github.com/hilonfot/network

go 1.13

replace (
	github.com/hilonfot/network/cmd_test => ./cmd_test
	github.com/hilonfot/network/conn => ./conn
	github.com/hilonfot/network/message => ./message
	github.com/hilonfot/network/server => ./server
	github.com/hilonfot/network/server/globalobj => ./server/globalobj
	github.com/hilonfot/network/utils => ./utils
	github.com/hilonfot/network/utils/catch => ./utils/catch
	github.com/hilonfot/network/utils/log => ./utils/log

)
