git clone https://github.com/medical-blockchain/storj-go.git
#changed line https://github.com/medical-blockchain/storj-go/blob/f10e1162b8dde73f2c35a5daa5f414aa565a244a/app.go#L14 from "./storj" to "github.com/medical-blockchain/storj-go/storj"
go mod init github.com/medical-blockchain/storj-go
go build app.go
