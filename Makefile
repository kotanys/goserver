SOURCES = $(wildcard src/*.go)

STDOUT = go-server-stdout.txt
STDERR = go-server-stderr.txt
PIDMASTER = "/tmp/.go-server-master.pid"
PIDSLAVE = "/tmp/.go-server-slave.pid"

start : start-master start-slave
start-master :
	@ go run $(SOURCES) -- config-master.json > $(STDOUT) 2> $(STDERR) & echo $$! > $(PIDMASTER)
	@ echo master PID: `cat $(PIDMASTER)`
start-slave :
	@ go run $(SOURCES) -- config-slave.json > $(STDOUT)2 2> $(STDERR)2 & echo $$! > $(PIDSLAVE)
	@ echo slave PID: `cat $(PIDSLAVE)`

stop : stop-master stop-slave
stop-master :
	@ pkill -2 -P `cat $(PIDMASTER)` 2>/dev/null
	@ kill -2 `cat $(PIDMASTER)` 2>/dev/null
	@ rm $(PIDMASTER)
stop-slave :
	@ pkill -2 -P `cat $(PIDSLAVE)` 2>/dev/null
	@ kill -2 `cat $(PIDSLAVE)` 2>/dev/null
	@ rm $(PIDSLAVE)

.PHONY : start stop start-master start-slave
