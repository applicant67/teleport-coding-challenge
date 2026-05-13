# Commands 

1. Start

```
go run service/client/main.go --command start "bash" "-c" "while true; do date '+%Y-%m-%d %H:%M:%S'; echo -n $'\\x00'; sleep 1; done | stdbuf -o0 cat"
# Job started with ID: bbb679e4-29ec-4f2f-b7d3-ae6121e55cbd
```

2. Query 
```
go run service/client/main.go --command query --jobid=bbb679e4-29ec-4f2f-b7d3-ae6121e55cbd
# Job query result: {
# 	"status": "Signaled",
# 	"pid": 169521,
# 	"signal": 15,
# 	"message": "Sucess!"
#   }
```
3. Stop
```
go run service/client/main.go --command stop --jobid=bbb679e4-29ec-4f2f-b7d3-ae6121e55cbd
# Job stopped: Stopped job!
```
4. Stream
```
go run service/client/main.go --command stream --jobid=bbb679e4-29ec-4f2f-b7d3-ae6121e55cbd
```