$env:GOOS="windows"
$env:GOARCH="amd64" 
$env:CGO_ENABLED="0" 
go build -o tcpulse.exe -trimpath -gcflags=all=-l=4
