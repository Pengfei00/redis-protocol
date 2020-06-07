# redis-protocol

## 安装
```
$ go get github.com/wnstar/redis-protocol
```

## replay
```
command := protocol.Command{}

error replay
replay := command.Error("error")

text replay
replay := command.Text("error")

array text replay
res := make([]interface{}, 2)
res[1] = "OK"
res[2] = 1
replay := command.Array(res)

integers replay
replay := command.Integers(1)

```

## receiver
```
var text = bufio.NewReader(bytes.NewReader([]byte("*1\r\n$7\r\nCOMMAND\r\n")))
cmd :=protocol.Command{}
fmt.Println(p.Receiver(text))
fmt.Println(p.Value)
```