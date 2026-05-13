# Testing memory limits: Example 1

1. Set memory.max to 100 MB 
2. Pass the command such as
```json
	{
		"command_args": [
		   "--vm",
		   "2",
		   "--vm-bytes",
		   "90M",
		   "--timeout",
		   "120"
		],
		"command_name": "stress"
	}
```
Note: seems like the "stress" voluntarily exits before the OOM killer intervenes. It returns exit code as 2 for an internal error rather than SIGKILL.
3.  Check if job_<uuid> cgroup is generated and cleanedup
```sh
cd /sys/fs/cgroup
watch -n 1 "ls -ltr | grep 'job'"
```
4. Note: need to do `sudo swapoff -a` so that OOM happens

# Testing memory: Example 2

1. Set memory.max to 100 MB 
2. Pass the command such as
```json
{
    "command_args": [
        "/home/fedora/job-worker-service/memory_hog.py"
    ],
    "command_name": "python"
}
```
3. Run the following python program
```py
memory_hog = []
while True:
    memory_hog.append(' ' * 1024 * 1024 * 10)  # Allocate 10MB blocks
```
4. In this case, you should get signal 9


# Testing CPU
* Run stress with 
```json
// --cpu 1: Spawns one worker to consume CPU.
{
    "command_args": [
       "--cpu",
       "1",
       "--timeout",
       "60"
    ],
    "command_name": "stress"
}
```
* If you do a top, you should be able to see the process consuming around 50% of one core.
<img width="669" alt="Screenshot 2024-09-06 at 3 08 47 AM" src="https://github.com/user-attachments/assets/26321f1c-3c10-41cd-8d14-95a2e276671d">

# Testing IO.

1. Run the following command with priority of 500 
```
// copy block size of 4 MB 20 times i.e 80 MB file
// oflag=direct -> bypasses cache and writes directly
// dd if=/dev/zero of=/home/fedora/abc-1.txt bs=4M count=20 oflag=direct

{
    "command_args": [
       "if=/dev/zero",
       "of=/home/fedora/abc-1.txt",
       "bs=4M",
       "count=20"
       "oflag=direct"
    ],
    "command_name": "dd"
}
```

2. Concurrently run the same command with priority of 1000 
```
// copy block size of 20 MB 20 times i.e 400 MB file
dd if=/dev/zero of=/home/fedora/abc-1.txt bs=4M count=20 oflag=direct

{
    "command_args": [
       "if=/dev/zero",
       "of=/home/fedora/abc-1.txt",
       "bs=20M",
       "count=20"
       "oflag=direct"
    ],
    "command_name": "dd"
}
```
3. Observe completion time in both cases and throughput with iostat.
