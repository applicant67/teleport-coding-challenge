# Testing
1. Basic start/stop/query/stream
2. Test start with 
    1. Basic commands
    2. Running process (simulate by infinite do while loop)
3. Stream to multiple clients (concurrency + real time updates)
    1. Able to stream to single client and output to should be similar to `docker -f`
    2. Test logs that outputs non-UTF-8 encoded content.
    3. Should be able to stream to multiple clients.
    4. Simulate a slow client (add a delay in the client's processing loop). There should be no degradation to other clients.
    5. Simulate the failure of the Linux process whose logs are being streamed. Verify that the server handles this gracefully, such as stopping the log stream and notifying the clients of the exit code.
    6. Test with large log messages.
4. Refer [cgroup docs](cgroup.md)
5. Coverage 
    1. Basic unit tests
    2. Test with `-race` flag
    3. govulncheck
