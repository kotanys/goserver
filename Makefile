SOURCES = $(wildcard *.go)

STDOUT = go-server-stdout.txt
STDERR = go-server-stderr.txt
PID = /tmp/.go-server.pid

start :
	@ go run $(SOURCES) > $(STDOUT) 2> $(STDERR) & echo $$! > $(PID)
	@ echo PID: `cat $(PID)`

stop :
	@ touch $(PID)
	@ pkill -2 -P `cat $(PID)` 2>/dev/null
	@ kill -2 `cat $(PID)` 2>/dev/null
	@ rm $(PID)

.PHONY : start stop
