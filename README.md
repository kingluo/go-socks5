# go-socks5
go socks5 server

# Build and Run

```
go install ./...
socks5_proxy -listen 127.0.0.1:8080

# test
curl -sS --proxy socks5://127.0.0.1:20002 http://www.baidu.com
# dns via socks5
curl -sS --proxy socks5h://127.0.0.1:20002 http://www.baidu.com
```

# lib usage

```go
for {
	conn, err := ssock.Accept()
	if err != nil {
		fmt.Println("accept", err.Error())
		os.Exit(1)
	}
	// wrap conn if you need to do custom encode/decode
	go socks5.RunSocks5Server(conn)
}
```
