# Reviewer feedback — distilled

One-off best-effort clustering of every captured reviewer comment across the public PR corpus (systems, fullstack, security-automation, sre). Themes are ranked by how often reviewers raised them. Quotes are verbatim from the captured `pr-review-comments.json` files; the `→` link points to the source thread on GitHub.

**Caveats.** Source dataset: only PR review threads that were public as of capture. Sample noise: a Haiku model did the initial clustering, so some quotes are placed in adjacent themes rather than the perfect one. Self-comments (the candidate replying to themselves) are filtered out; reviewer ↔ candidate threading is not preserved here — read the themed files under `tracks/<track>/corpus/reviewer-feedback/` if you want the curated narrative version.

**Coverage**: 2038 clustered quotes, 622 unclustered (LGTM/emoji noise), 21 themes.

## Themes by frequency

| Rank | Theme | Quote count |
|---|---|---|
| 1 | [`output-streaming`](#output-streaming) | 381 |
| 2 | [`auth-identity`](#auth-identity) | 365 |
| 3 | [`status-lifecycle`](#status-lifecycle) | 297 |
| 4 | [`error-handling`](#error-handling) | 147 |
| 5 | [`api-shape`](#api-shape) | 142 |
| 6 | [`code-style-nits`](#code-style-nits) | 139 |
| 7 | [`process-containment`](#process-containment) | 133 |
| 8 | [`testing`](#testing) | 110 |
| 9 | [`scope-cutting`](#scope-cutting) | 92 |
| 10 | [`design-doc-style`](#design-doc-style) | 64 |
| 11 | [`concurrency-and-locking`](#concurrency-and-locking) | 45 |
| 12 | [`implementation-details`](#implementation-details) | 23 |
| 13 | [`reproducibility-and-tooling`](#reproducibility-and-tooling) | 19 |
| 14 | [`command-invocation`](#command-invocation) | 16 |
| 15 | [`code-organization`](#code-organization) | 16 |
| 16 | [`auth-storage`](#auth-storage) | 15 |
| 17 | [`auth0-least-privilege`](#auth0-least-privilege) | 10 |
| 18 | [`tls-and-hashing`](#tls-and-hashing) | 9 |
| 19 | [`input-validation`](#input-validation) | 8 |
| 20 | [`setup-configuration`](#setup-configuration) | 5 |
| 21 | [`csrf-protection`](#csrf-protection) | 2 |

---

## output-streaming

_How new bytes are signaled and handled; managing multiple readers; slow-reader handling and write-after-close races._

**381 quotes** from `30` distinct reviewers across `47` candidate submissions.

**@rosstimothy** on `adalton/teleport-exercise` [→](https://github.com/adalton/teleport-exercise/pull/1#discussion_r767874390)

> Why ignore `closed` here?
> ```suggestion
> 			bufferSize, closed := b.buffer.waitForChange(nextByte)
> 
> 			// At this point the underlying buffer could be closed, there could
> 			// be new bytes in the buffer to process, or both.
> 
> 			if bufferSize == nextByte  && closed {
> ```

**@codingllama** on `adalton/teleport-exercise` [→](https://github.com/adalton/teleport-exercise/pull/1#discussion_r768030168)

> If we made `blockIo` a public struct we could take advantage of literal assignment and do away with a bunch of methods. Same goes for various other types.
> 
> Eg:
> 
> ```go
> _ = &BlockIO{ReadBPSDevice: x, WriteBPSDevice: y}
> ```
> 
> vs
> 
> ```go
> _ = NewBlockIoController(). SetReadBpsDevice(x). SetWriteBpsDevice(y)
> ```

**@codingllama** on `adalton/teleport-exercise` [→](https://github.com/adalton/teleport-exercise/pull/1#discussion_r768052216)

> ```suggestion
> 					// Create a copy here because we're reusing readBuffer here.
> ```

**@codingllama** on `adalton/teleport-exercise` [→](https://github.com/adalton/teleport-exercise/pull/1#discussion_r768054468)

> Hold read locks instead?

**@codingllama** on `adalton/teleport-exercise` [→](https://github.com/adalton/teleport-exercise/pull/1#discussion_r768058650)

> How is this scenario different from the one above? It looks to be the same test case (reading from end, buffer > remaining chars).

**@codingllama** on `adalton/teleport-exercise` [→](https://github.com/adalton/teleport-exercise/pull/1#discussion_r768692731)

> > That also changes the validation portion of the tests from semi-simple asserts into "stat this file and make sure it exists" and "read this file and make sure it has the expected content".
> 
> This is one of the greatest benefits of the change, IMO - we are not doing asserts on special structs anymore, just reading files (which is what the program is supposed to output).
> 
> > If you'd like, I can rework this to use the temp dir approach -- just let me know.
> 
> I wanted to hear your thoughts more than anything else. We can keep as-is.

**@codingllama** on `adalton/teleport-exercise` [→](https://github.com/adalton/teleport-exercise/pull/1#discussion_r768967316)

> You could start with a `var buckets [numGoroutines]bytes.Buffer` and skip the transformation below.

**@rosstimothy** on `adalton/teleport-exercise` [→](https://github.com/adalton/teleport-exercise/pull/2#discussion_r771430966)

> What if I want both stdout and stderr in one stream?

**@codingllama** on `adalton/teleport-exercise` [→](https://github.com/adalton/teleport-exercise/pull/2#discussion_r771438230)

> I'm finding this a bit hard to reason about. Even with the various b.Closed() checks it seems we still have a "soft" race between b.Close() and writing to the channel (line 97) - ie, Close() could happen at "t0" and we still get channel write at "t1".

**@rosstimothy** on `adalton/teleport-exercise` [→](https://github.com/adalton/teleport-exercise/pull/2#discussion_r771446123)

> You could also potentially utilize a channel instead or repeatedly attempting to dial

**@codingllama** on `adalton/teleport-exercise` [→](https://github.com/adalton/teleport-exercise/pull/2#discussion_r771460755)

> nit: don't send os.Kill - it looks like the signal makes a difference, but it doesn't - any number would do (right?).
> 
> ```suggestion
> 	stop <- 1
> ```

**@codingllama** on `adalton/teleport-exercise` [→](https://github.com/adalton/teleport-exercise/pull/2#discussion_r771468645)

> I wouldn't bother distinguishing "invalid" from "unknown", less code to write. Up to you.

**@codingllama** on `adalton/teleport-exercise` [→](https://github.com/adalton/teleport-exercise/pull/2#discussion_r772449077)

> Ditch getPort and return the unchanged addr? A bit less code to write.

**@rosstimothy** on `adalton/teleport-exercise` [→](https://github.com/adalton/teleport-exercise/pull/3#discussion_r773291195)

> nit: copy/paste error
> ```suggestion
> // StreamStderr invokes an RPC on the JobManager server to stream the standard
> ```

**@sclevine** on `benmoss/job-worker-service` [→](https://github.com/benmoss/job-worker-service/pull/1#discussion_r2669779942)

> How will you avoid race conditions like this?
>
> - Writer writes and signals
> - Reader reads to EOF
> - Writer writes and signals
> - Reader waits
> - Writer does not write again for 10 minutes
>
> Result: reader has output delayed for 10 minutes

**@nklaassen** on `benmoss/job-worker-service` [→](https://github.com/benmoss/job-worker-service/pull/1#discussion_r2669981937)

> I think this api will limit your ability to terminate the stream or clean up resources early if/when the client disconnects before the logs have been read to EOF

**@nklaassen** on `benmoss/job-worker-service` [→](https://github.com/benmoss/job-worker-service/pull/1#discussion_r2669996011)

> I would like to see mentioned how you will handle clients that disconnect before reading the whole output or before the job completes, and also how readers will know that the job has completed and there is no more output to wait for

**@nklaassen** on `benmoss/job-worker-service` [→](https://github.com/benmoss/job-worker-service/pull/1#discussion_r2670607692)

> is the condition variable you're talking about a sync.Cond? If so, how will context cancellation wake a Read call blocked on the Cond?

**@nklaassen** on `benmoss/job-worker-service` [→](https://github.com/benmoss/job-worker-service/pull/1#discussion_r2670609578)

> what does the Close method on the returned ReadCloser do? is it equivalent to context cancellation? if so, do you really need both?

**@sclevine** on `benmoss/job-worker-service` [→](https://github.com/benmoss/job-worker-service/pull/1#discussion_r2670644719)

> > I think I can keep a counter on the reader and the writer
> 
> Sounds good 👍 
> 
> > I can include this in the RFD if you think it's warranted, I think it might be too much of an implementation detail
> 
> I think a brief clarification that your signaling mechanism will track progress to ensure real-time output would be nice, but no need to go into more detail than that 🙂

**@nklaassen** on `benmoss/job-worker-service` [→](https://github.com/benmoss/job-worker-service/pull/3#discussion_r2683997432)

> seems aggressive for each outputWriter instance to try to remove the parent directory of the file it's handling and all other files in the directory, I wouldn't expect a Cleanup on a single outputWriter to affect every other outputWriter

**@nklaassen** on `benmoss/job-worker-service` [→](https://github.com/benmoss/job-worker-service/pull/3#discussion_r2684001826)

> I think the order of these should be switched, no? Or else you could still miss a final write + close that happened after r.file.Read returned and before getNotifyState got the lock

**@nklaassen** on `benmoss/job-worker-service` [→](https://github.com/benmoss/job-worker-service/pull/3#discussion_r2684021930)

> consider `w.wg.Go(func() { w.waitForProcess(entry) })`

**@dboslee** on `benmoss/job-worker-service` [→](https://github.com/benmoss/job-worker-service/pull/3#discussion_r2684531330)

> there are multiple places where this channel is closed. is it possible for a write->broadcast to come in after markDone? I think this would cause a panic.

**@sclevine** on `benmoss/job-worker-service` [→](https://github.com/benmoss/job-worker-service/pull/3#discussion_r2688696816)

> > Right now if they inherit stdout/stderr cmd.Wait will actually keep blocking until they exit, but if either close stdout/stderr or don't inherit it they will be orphaned.
>
> This behavior seems like it could be surprising for remote job execution.
>
> > do you want me to implement that?
>
> Waiting for child processes to terminate is not out-of-scope for the challenge 🙂

**@sclevine** on `benmoss/job-worker-service` [→](https://github.com/benmoss/job-worker-service/pull/3#discussion_r2695364059)

> There are a couple of strategies you could use to wait on the child processes and avoid the tight polling loop, but I'd recommend completing the challenge first and coming back to this if there is time 🙂

**@smallinsky** on `bill-rich/jobworker` [→](https://github.com/bill-rich/jobworker/pull/1#discussion_r1606456288)

> How will stream work? Where do you plan to store process logs, in memory or in a file? How streaming to multiple clients/connection will be supported ?

**@sclevine** on `bkneis/jobworker` [→](https://github.com/bkneis/jobworker/pull/2#discussion_r1657660848)

> It looks like this consumes an fd per concurrent reader. Are you concerned about exhausting the number of available descriptors?

**@dboslee** on `bkneis/jobworker` [→](https://github.com/bkneis/jobworker/pull/2#discussion_r1661344851)

> Agreed inotify/tail are probably more complex/hacky than its worth here. Thanks for the thought process though.
> 
> As for the reading blocking until the ReadCloser is closed, I think this does impact the API a bit where retrieving logs hangs until the client decides to end the stream instead of being able to close the stream server side with a good status code when the end of the log file is reached and there will be no more additional output.
> 
> Not a big deal if its to much work to fix this, but I am curious to hear your thoughts on how this could be improved.

**@tcsc** on `bkneis/jobworker` [→](https://github.com/bkneis/jobworker/pull/2#discussion_r1663362677)

> I have a personal phobia of `bool` parameters, as they're usually completely uninformative at the call site. Maybe a 2-state enum?
> 
> ```golang
> type OutputMode int
> const (
>   OutputModeNoFollow OutputMode = 0 // feel free to supply a better name
>   OutputModeFollow   OutputMode = 1
> )
> ```
> 
> ```suggestion
> func (job *Job) Output(mode OutputMode) (reader io.ReadCloser, err error) {
> ```
> 
> 
> which makes the call site go from 
> 
> ```golang
> job.Output(true)
> ```
> 
> to 
> 
> ```golang
> job.Output(OutputModeFollow)
> ```

**@tcsc** on `bkneis/jobworker` [→](https://github.com/bkneis/jobworker/pull/2#discussion_r1663392416)

> ```suggestion
> 	return os.WriteFile(filepath.Join(cg.groupPath(name), file), []byte(val), 0644)
> ```

**@tcsc** on `bkneis/jobworker` [→](https://github.com/bkneis/jobworker/pull/2#discussion_r1663398202)

> I think there is a mismatch between the comment and code here. If I understand it correctly, the `io.Reader()` returned by `Output(false)` will _not_ tail the job output and only return the logs until it hits the first `EOF`.

**@tcsc** on `bkneis/jobworker` [→](https://github.com/bkneis/jobworker/pull/2#discussion_r1663412236)

> As far as I can see you're not really asserting either of those things.
> 
> This test just stops counting after it gets to _n_ lines of log output. It counts getting any arbitrary number lines as a success, as long as all of those log lines contain "hello". You could have no output at all and this test would pass.
> 
> You probably want to add a line count to the output (see above) and rewrite the assertion to look something like:
> 
> ```golang
> go func(r io.ReadCloser) {
> 	scanner := bufio.NewScanner(r)
> 	logs := []string{}
> 	for scanner.Scan() {
> 		logs = append(logs, scanner.Text())
> 	}
> 
> 	if len(logs) != n {
> 		t.Failf("Expected %d lines, got %d", n, len(logs)
> 	}
> 
> 	for i, log := range logs {
> 		expected := fmt.Sprintf("%d: %s", i+1, echo)
> 		if log != expected {
> 			t.Failf("Line %d, expected %s, got %s", i, expected, log)
> 		}
> 	}
> }(reader)
> ```

**@sclevine** on `bkneis/jobworker` [→](https://github.com/bkneis/jobworker/pull/2#discussion_r1663624229)

> I'm not sure if adding an extra flag improves the UX much. The behavior I'd expect is like `docker logs -f`. If the command quits, the reader gets io.EOF. Otherwise, the reader waits for more output.
> 
> Without that UX, the use case for remote commend execution seems limited. E.g., I can't reliably start a command and capture its output to a file without waiting for it to quit somehow.

**@rosstimothy** on `bucknercd/jobworker` [→](https://github.com/bucknercd/jobworker/pull/1#discussion_r2266878347)

> What happens if the entire file is consumed and the process is still running and might write more output to the file?

**@rosstimothy** on `bucknercd/jobworker` [→](https://github.com/bucknercd/jobworker/pull/1#discussion_r2266886361)

> Are pipes needed to write the output to files? Can you think of a way to do this without copying in goroutines?

**@rosstimothy** on `bucknercd/jobworker` [→](https://github.com/bucknercd/jobworker/pull/1#discussion_r2269748350)

> I don't know that these two items add much value. There is a built in mechanism in gRPC to indicate when a stream has closed. By and large the clients shouldn't need to know about chunk sizing and have no recourse to do anything with the information.

**@rosstimothy** on `bucknercd/jobworker` [→](https://github.com/bucknercd/jobworker/pull/1#discussion_r2269776087)

> What benefit does the channel provide if this is already running in a dedicated goroutine?

**@zmb3** on `bucknercd/jobworker` [→](https://github.com/bucknercd/jobworker/pull/1#discussion_r2271494402)

> Probably no need to return the ID here since the client already specified it in the request.

**@zmb3** on `bucknercd/jobworker` [→](https://github.com/bucknercd/jobworker/pull/1#discussion_r2271531348)

> Maybe I'm just confused by the order in which this is written, but if you've deleted the cgroup directory in step 4 there won't be anything to write here.

**@zmb3** on `bucknercd/jobworker` [→](https://github.com/bucknercd/jobworker/pull/1#discussion_r2271540945)

> After reading through this a couple times I'm still not totally clear on the authorization model. Maybe having both roles _and_ user groups is more complicated than what we need here? I'd welcome any opportunity to simplify.

**@zmb3** on `bucknercd/jobworker` [→](https://github.com/bucknercd/jobworker/pull/1#discussion_r2271570630)

> ```suggestion
> - The `stream` command has `tail` like functionality. This means that if the job is still running, the stream will continue to follow the output in real-time. If the job has exited, the stdout and stderr file descriptors will be closed and the stream will end.
> ```

**@espadolini** on `chintamanil/job-worker` [→](https://github.com/chintamanil/job-worker/pull/1#discussion_r2822087956)

> This will block until there's new output (or the job terminates) even if the client disconnects, so a long running job with no output will accumulate goroutines all blocked on the waitgroup.

**@espadolini** on `chintamanil/job-worker` [→](https://github.com/chintamanil/job-worker/pull/1#discussion_r2822095110)

> Keeping stdout and stderr completely separate will result in output that is hard to understand if the process is sending some data to stdout and some data to stderr.

**@espadolini** on `chintamanil/job-worker` [→](https://github.com/chintamanil/job-worker/pull/1#discussion_r2822104412)

> Checking `stream.Context().Err()` does not solve this problem.

**@tigrato** on `chintamanil/job-worker` [→](https://github.com/chintamanil/job-worker/pull/1#discussion_r2822317503)

> Can we satisfy io.ReadCloser so io.Copy works?

**@tigrato** on `chintamanil/job-worker` [→](https://github.com/chintamanil/job-worker/pull/3#discussion_r2829707136)

> Can we design the API so that users don’t have to explicitly invoke a broadcast when the context is canceled? Understanding the broadcast concept forces users to understand how streaming is implemented,
>
> I think I already mentioned before, but could the API implement io.ReadCloser instead? It’s much easier to understand that you should close the reader when you’re no longer interested, rather than having to cancel the context and then call broadcast.

**@tigrato** on `chintamanil/job-worker` [→](https://github.com/chintamanil/job-worker/pull/2#discussion_r2832678354)

> If a reader calls `reader.Close()` without cancelling the context, will it have any real effect? Will it cause a deadlock?
>
> Close is also a dangerous no-op in this situation

**@espadolini** on `chintamanil/job-worker` [→](https://github.com/chintamanil/job-worker/pull/2#discussion_r2833520757)

> Similarly to JobStore, I don't think the user of rwmutex is warranted here, especially when every reader currently acquires the exclusive lock just to wait for output.

**@rosstimothy** on `clydotron/job_worker_service` [→](https://github.com/clydotron/job_worker_service/pull/1#discussion_r1599972359)

> Details are a bit light here. Where and how is output stored? How are channels used to provide data to clients? How are clients made aware that more data is available if they've already consumed all output?

**@rosstimothy** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/1#discussion_r1782726751)

> Does reading job data also include the ability to stop a job? Just want to make sure I understand the proposed scheme here. It sounds like this is an "owner" model where a job is only accessible to the identity that created the job. Is that correct?

**@codingllama** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/1#discussion_r1783300070)

> Could you add some details on how you intend to implement the broadcast?

**@espadolini** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/1#discussion_r1785885226)

> The requirements call for the full output to be streamed from the beginning of the process until the current time, using a broker that we have to then be careful to attach a subscriber to before anything is output seems like a roundabout way of achieving that.

**@codingllama** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/2#discussion_r1787738895)

> ```suggestion
> // StreamingMiddleware handles the authn from mtls for streaming endpoints.
> ```

**@codingllama** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/2#discussion_r1787744406)

> Should we move this outside of the main package, or are we going to write everything in here?

**@rosstimothy** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/4#discussion_r1790169186)

> If the user is authorized to open the stream, then shouldn't they be authorized to send and receive messages via the stream?

**@codingllama** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/4#discussion_r1790718635)

> nit: I have a personal crusade against the `srv` abbreviation:
> 
> ```suggestion
> 	server any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler,
> ```

**@codingllama** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/4#discussion_r1790747691)

> Since the middleware is querying the job already, should it "pass it down" to the handlers so we avoid the double-query?
> 
> Alternatively we can have the handlers invoke the middleware - more boilerplate, less magic (and a few different code tradeoffs).
> 
> WDYT?

**@codingllama** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/5#discussion_r1790767751)

> If the broadcaster accumulated the output you wouldn't need to expose Pause/Unpause.

**@codingllama** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/5#discussion_r1790769020)

> I get the usage, but I question whether we want Write to be part of the Broadcaster's public API?

**@codingllama** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/5#discussion_r1790784294)

> A few reasons:
> 
> 1. This is a block, it just doesn't like so visually
> 2. The label/goto idiom is rare in practice, which I expect would cause weirdness for most readers
> 3. This might encourage more sophisticated goto uses sneaking into the codebase, at which point we may get into real goto readability downsides
> 
> Note that this isn't a "goto bad" comment, but I also don't see the appeal in replacing the more usual `for {}` construct.

**@espadolini** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/5#discussion_r1791368721)

> `FreeOSMemory` already does a GC pass, are we doing two for a specific reason? (finalizers?)

**@espadolini** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/5#discussion_r1791389652)

> While this is correct for the current implementation of `StreamLogs` I wouldn't be surprised if a refactor changed the implementation of the returned reader to one that can return a nonzero `n` and an error at the same time (as the `io.Reader` contract allows), and that would result in losing some data here. We should send the data through the stream (if we got some) before checking the error.

**@espadolini** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/5#discussion_r1791425174)

> FWIW switching to a `strings.Builder` is a +3-4 diff and avoids copying the whole output on every `StreamLogs` call.

**@espadolini** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/5#discussion_r1791434240)

> This has the potential to actually block the command itself once the OS pipe it's using is full, right? Not just other readers.

**@rosstimothy** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/5#discussion_r1791833042)

> ```suggestion
> func (b *Broadcaster) Resume() { b.mu.Unlock() }
> ```

**@rosstimothy** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/5#discussion_r1791837918)

> I think there is a typo here?
> ```suggestion
> // If broadcaster is closed, do nothing.
> ```

**@rosstimothy** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/5#discussion_r1791841732)

> The Broadcaster API is quite nuanced and open to footguns. If you are set moving forward with this API it would be nice if at minimum the godocs better explained the intended use of the Broadcaster.

**@rosstimothy** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/5#discussion_r1791843327)

> This comment exposes implementations detail that shouldn't matter to consumers and offers little insight into the consequences of pausing a broadcast without resuming, or attempting to resume an unpaused broadcast.

**@rosstimothy** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/5#discussion_r1791863285)

> > Multiple concurrent clients should be supported.
> 
> This could  impact the ability to meet the above requirement of the challenge. If there is a large amount of historical output and existing clients are already streaming the output any new clients streaming said output would likely cause lag in consuming output for already subscribed clients while the broadcaster is paused and the output is copied.

**@rosstimothy** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/5#discussion_r1791931613)

> The strings.Builder might alleviate some of the problems at this level, however, the Broadcaster design is still going to be limited by https://github.com/creack/telepilot/pull/5/files#diff-ae467001c26644299f4da860817ff90008e892d897df53fe209cb4122037123cR52-R56

**@codingllama** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/5#discussion_r1792425301)

> I missed some of the "live" discussion here (and in the other thread), but I think it's safe to say we would rather be resilient to a slow consumer, despite what Docker may do. This is one of the more interesting parts of the challenge for us.

**@codingllama** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/5#discussion_r1792436929)

> ```suggestion
> // Surface the lock as Pause/Resume to freeze the broadcast.
> ```

**@codingllama** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/5#discussion_r1792437573)

> ```suggestion
> // Calling Pause() pauses all the clients as well as the use of the broadcast
> ```

**@codingllama** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/5#discussion_r1792445233)

> Document that clients==nil means that the Broadcaster is closed?

**@codingllama** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/5#discussion_r1792447925)

> Suggestion: less indirection, so slightly clearer:
> 
> ```suggestion
> 	j.cmd.Stderr = j.broadcaster // NOTE: Merge out/err for simplicity. Should split them for production.
> ```

**@codingllama** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/5#discussion_r1792453059)

> Suggestion: most manual unlocks come right before a return, suggesting we could use a defer (making the lock easier to reason about). The exception is the "j.mu.Unlock(); <-j.waitchan; return" part, but that has alternatives too.

**@codingllama** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/5#discussion_r1793588633)

> Make sure all broadcast goroutines stopped before returning from Close as well?
> 
> It would be nice if we could do the same in Unsubscribe, but only for that particular client, but I'll take a Close guarantee only.

**@codingllama** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/5#discussion_r1793614039)

> ```suggestion
> 			t.Fatal("Timeout waiting for read loop to end.")
> ```

**@rosstimothy** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/5#discussion_r1793684071)

> What ownership requirements does gRPC have for data sent via a stream?

**@rosstimothy** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/5#discussion_r1793689370)

> This means that `var b Broadcaster` is implicitly considered to be closed correct?

**@rosstimothy** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/5#discussion_r1793694358)

> If a client is slow to consume the 128 buffered entires, it will be unsubscribed, closed and receive no future data?

**@rosstimothy** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/5#discussion_r1793698045)

> Should this prevent additional writes to a "closed" writer?

**@rosstimothy** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/5#discussion_r1793894485)

> To consumers this is a io.WriteCloser though, and I can imagine it would be expected that writing to a closed writer shouldn't add more data.

**@rosstimothy** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/5#discussion_r1793900426)

> From the ServerStream [docs](https://pkg.go.dev/google.golang.org/grpc#ServerStream):
> ```
> // SendMsg does not wait until the message is received by the client. An
> // untimely stream closure may result in lost messages.
> //
> // It is safe to have a goroutine calling SendMsg and another goroutine
> // calling RecvMsg on the same stream at the same time, but it is not safe
> // to call SendMsg on the same stream in different goroutines.
> //
> // It is not safe to modify the message after calling SendMsg. Tracing
> // libraries and stats handlers may use the message lazily.
> SendMsg(m any) error
> ```
> 
> The most important and relevant bit there is the following:
> > It is not safe to modify the message after calling SendMsg

**@rosstimothy** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/5#discussion_r1793925835)

> Yeah I hear you, it likely wouldn't add much value no need to address here. Though being pedantic this a lockWriteCloser and not nopStringsBuilderWithLockCloser.

**@espadolini** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/6#discussion_r1795029157)

> Why not just write a static `+cpu +memory +io` or perhaps just blindly enable all the required ones? IIRC it's not an error to enable a controller that's already enabled.

**@espadolini** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/6#discussion_r1795068432)

> Is waitpid on the init of a pid namespace guaranteed to return after all the processes in the namespace are fully cleaned up? Even ignoring the potential for other things to be put in the cgroup by something else (which violates the cgroup single writer principle) I think that there can be some time between the termination of the main process and the cgroup becoming unpopulated.

**@codingllama** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/5#discussion_r1795623554)

> Is this somehow different from letting the empty output reach the MultiReader, or is it just an optimization? No need for changes, just curious.

**@codingllama** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/5#discussion_r1795657133)

> Let me see if I got it right:
> 
> 1. Closing the BufferedBroadcaster closes the clients
> 2. Closing a client closes the msgs channel
> 3. A closed msgs channel should make the Write goroutine exit
> 4. That was to happen in 500ms, otherwise the client's io.WriterCloser is closed, preventing further writes and making the Write goroutine abort
> 
> This seems like an easy avenue for a slow client to lose data once a job stops, as we are somewhat negating the advantages of the buffer in the msgs channel.
> 
> I think a simple solution is to simply not take ownership of the client's io.WriterClose - ie, take only an io.Writer when subscribing and let the client handle the closing.

**@rosstimothy** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/6#discussion_r1795686928)

> Use a buffered solution or a scanner instead of reading the entire content at once?

**@codingllama** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/6#discussion_r1795687046)

> Wrong file/write for freeze?

**@codingllama** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/6#discussion_r1795691297)

> I think the freeze block could benefit from being its own function. It makes the break/return dynamics simpler to read and we could move the error swallowing to the caller.

**@codingllama** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/6#discussion_r1795817101)

> ```suggestion
> 		t.Cleanup(func() { noError(t, <-ch, "Stream Logs") }) // TODO: Consider adding timeout.
> ```

**@codingllama** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/6#discussion_r1795876324)

> This test randomly hangs for me in the golang:1.23.2-alpine container, I think it's the pipe write.

**@espadolini** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/6#discussion_r1796874512)

> This early exit is not closing `w`, and nothing is closing `r`.

**@GavinFrazar** on `devnulled/runjob` [→](https://github.com/devnulled/runjob/pull/2#discussion_r1797319710)

> > The Jobs' buffer will support concurrent reads and writes via a mutex.
> > Subscribers can be created at any time and will be fed all of the contents in the buffer created since the start of the Job.
> 
> How will you ensure that all clients are able to read the output independently, i.e without waiting on each other?
> And how will you ensure that readers can't block writing new command output?
> 
> For example, suppose there is a long-running job and there are several subscribers that have already consumed all of the command output buffer. Then a new subscriber is added. This new subscriber is quite slow though (maybe their connection isn't great). While the new subscriber "catches up" reading the output buffer, more output is generated.
> How can we ensure that the slow reader doesn't prevent the others from seeing the latest output?

**@nklaassen** on `diptadas/job-worker` [→](https://github.com/diptadas/job-worker/pull/1#discussion_r645156572)

> Don't wrap the command with `bash` by default, please let the client send the exact command, with arguments. If the client wants to run bash, they can do so explicitly

**@Joerger** on `diptadas/job-worker` [→](https://github.com/diptadas/job-worker/pull/1#discussion_r645196347)

> Are permissions tied to specific jobs? Or can a Read+Write user stop/query/purge any job?

**@awly** on `diptadas/job-worker` [→](https://github.com/diptadas/job-worker/pull/1#discussion_r646774619)

> ```suggestion
>     - Read and write: Required for creating and stopping a job
> ```

**@awly** on `diptadas/job-worker` [→](https://github.com/diptadas/job-worker/pull/1#discussion_r647023573)

> nit: you will need some synchronization around this, if you want to forward both stdout and stderr to the same buffer

**@Joerger** on `diptadas/job-worker` [→](https://github.com/diptadas/job-worker/pull/2#discussion_r647908103)

> How will you handle synchronous writes to `allJobs`? maps alone aren't thread safe.

**@Joerger** on `diptadas/job-worker` [→](https://github.com/diptadas/job-worker/pull/2#discussion_r647909238)

> How will `bytes.Buffer` handle concurrent writes?

**@Joerger** on `diptadas/job-worker` [→](https://github.com/diptadas/job-worker/pull/2#discussion_r647913113)

> This is a race condition against the write to `job.Status` in `handleFinish`

**@Joerger** on `diptadas/job-worker` [→](https://github.com/diptadas/job-worker/pull/2#discussion_r647913755)

> Can you think of a simpler solution than a `sync.WaitGroup` for tracking when the job is done?

**@Joerger** on `diptadas/job-worker` [→](https://github.com/diptadas/job-worker/pull/2#discussion_r647924315)

> If you are going to write to `*job.Output` in this way, you'll need to consider synchronization (along with all other job read/writes). It also may become out of sync with `OutputBuffer`. What about returning a shallow copy with `Output` populated instead?
> ```suggestion
> 		jobStatus := *job
> 		jobStatus.Output = job.outputBuffer.String()
> 		return jobStatus, nil
> ```

**@Joerger** on `diptadas/job-worker` [→](https://github.com/diptadas/job-worker/pull/2#discussion_r652062286)

> What about concurrent reads/writes to the outputBuffer?

**@awly** on `diptadas/job-worker` [→](https://github.com/diptadas/job-worker/pull/2#discussion_r652906918)

> error returned by https://golang.org/pkg/os/exec/#Cmd.Wait

**@rosstimothy** on `DiscoRiver/jobber` [→](https://github.com/DiscoRiver/jobber/pull/1#discussion_r1989345076)

> Including the id is somewhat redundant since it must be provided in the request and all response across the stream will be for the same id.

**@zmb3** on `DiscoRiver/jobber` [→](https://github.com/DiscoRiver/jobber/pull/1#discussion_r1989347475)

> Could you provide more details on how streaming will be implemented? In particular, this channel-based API looks like it could suffer from slow consumers holding up the producer.

**@tigrato** on `DiscoRiver/jobber` [→](https://github.com/DiscoRiver/jobber/pull/1#discussion_r1989471165)

> Do we have guarantees the output buffer is utf-8 aligned data?

**@tigrato** on `DiscoRiver/jobber` [→](https://github.com/DiscoRiver/jobber/pull/1#discussion_r1989477996)

> what's the downside if a subscriber is slow and the the library can't send the logs to the channel?

**@zmb3** on `DiscoRiver/jobber` [→](https://github.com/DiscoRiver/jobber/pull/1#discussion_r1989633434)

> You can have the CLI always send the same set of hard coded limits on all jobs. 
>
> Same functionality, just less plumbing work.

**@rosstimothy** on `DiscoRiver/jobber` [→](https://github.com/DiscoRiver/jobber/pull/1#discussion_r1991731526)

> What will happen if a client consumes all content in the file and the process is still running and may write more output to the file in the future?

**@zmb3** on `DiscoRiver/jobber` [→](https://github.com/DiscoRiver/jobber/pull/1#discussion_r2063662703)

> No problem sticking with the channel so long as you ensure that one slow consumer can't hold up the producer. 
>
> As far as the buffer goes, I don't think waiting for 4KB of data before sending anything to the client is the right UX. Imagine a process that outputs a few bytes per second. It would take 30-60 minutes before the client sees any output. 
>
> Can we do something closer to real time so that the client gets a nice experience no matter how much (or how little) output the job produces?

**@zmb3** on `DiscoRiver/jobber` [→](https://github.com/DiscoRiver/jobber/pull/1#discussion_r2065163001)

> Gotcha. I found it a bit confusing to read about DB config in the design document for a program that doesn't use a DB.
>
> In any case, I'd simplify this as much as possible to save time. No need for flexible configuration handling or parsing from the environment.

**@tigrato** on `DiscoRiver/jobber` [→](https://github.com/DiscoRiver/jobber/pull/1#discussion_r2066258831)

> When a client reaches the end of file, Read will return io.EOF. How will the system handle new incoming data?
> How will the system wake up on new writes to the temporary file?

**@rosstimothy** on `dtom90/teleport-backend-challenge` [→](https://github.com/dtom90/teleport-backend-challenge/pull/1#discussion_r2052995048)

> Could you include some details on how you will handle the situation where all of the output in the file has been consumed by the user, the process is still running and may write more output to the file in the future? How is any new output being detected and provided to the user?

**@greedy52** on `ehsaniara/joblet` [→](https://github.com/ehsaniara/joblet/pull/2#discussion_r2115932934)

> thanks for the update
> 
> > Output should be from start of process execution.
> 
> how does the streaming client get the "history" outputs before updates from the subscription?

**@greedy52** on `ehsaniara/joblet` [→](https://github.com/ehsaniara/joblet/pull/2#discussion_r2121171523)

> sorry for the late reply. that's not what i meant. cgroup is for the job not for the client.
> 
> for clients, no matter how many subscribers to stream the job, they are all reading the same content at the end of the day. so i am just asking if any optimization can be done to avoid `50` buffer per subscriber.
> 
> it's just a question though, you don't have to derail from your current design.

**@zmb3** on `ehsaniara/joblet` [→](https://github.com/ehsaniara/joblet/pull/2#discussion_r2121644671)

> I don't quite understand the reasoning behind correlating the size of a job's output buffer to the amount of memory the job is allowed to use when running. I don't think we can assume there's any correlation between those two things. You can have a job that is very intensive but produces little output, or a job that uses little memory but is very chatty.

**@russjones** on `flychicken123/Job_Worker` [→](https://github.com/flychicken123/Job_Worker/pull/1#discussion_r2065390142)

> Can you share some more details about how you plan to implement output streaming? For example, tell me what interface you plan to implement in your library and have us 3-5 sentences on how you plan to implement output streaming.

**@GavinFrazar** on `flychicken123/Job_Worker` [→](https://github.com/flychicken123/Job_Worker/pull/1#discussion_r2067104177)

> nit: I would prefer to avoid coupling the arguments for these RPCs, i.e have separate messages `[Stop/Get/Stream]JobRequest`

**@rosstimothy** on `flychicken123/Job_Worker` [→](https://github.com/flychicken123/Job_Worker/pull/1#discussion_r2069327458)

> Isn't the actual messages in `data`? It seems that from the comment `stream` will either be stdout or stderr to indicate which output stream the `data` was from. Is that not correct?

**@eriktate** on `GevorgGal/jobworker` [→](https://github.com/GevorgGal/jobworker/pull/1#discussion_r2975120187)

> This line implies that the reader _pushes_ data somewhere rather than a caller pulling data from it. Is that accurate?

**@eriktate** on `GevorgGal/jobworker` [→](https://github.com/GevorgGal/jobworker/pull/1#discussion_r2975171300)

> Graceful exits aren't required, so you could consider reducing this to sending `SIGKILL`

**@eriktate** on `GevorgGal/jobworker` [→](https://github.com/GevorgGal/jobworker/pull/1#discussion_r2977735619)

> Ah okay, I think I misunderstood. So the `buffer.ReadFrom(offset)` is the only actual interaction with the output buffer. At first I thought this was some kind of reader provided by the library. This makes sense to me :+1:

**@Tener** on `GevorgGal/jobworker` [→](https://github.com/GevorgGal/jobworker/pull/2#discussion_r2980018065)

> This reads like an optimization: "only wait for more data if previous read returned nothing"; it should work just fine without that, though? Can you comment on your reasoning here?

**@Tener** on `GevorgGal/jobworker` [→](https://github.com/GevorgGal/jobworker/pull/2#discussion_r2980046770)

> My only concern is efficiency of streaming: if we have a large buffer already accumulated, each reader will immediately make a full copy of it, even though we'll be sending it in much smaller chunks.
>
> One way to it is by adding "max read size" parameter that limits the amount of data each individual `ReadFrom` can return. The other, canonical way is to pass explicit buffer to be filled; this avoids allocation by reusing the buffer while also specifying the max read size in a natural way.
>
> That would cost us API complexity, so personally I'm fine with this remaining a theoretical improvement.

**@Tener** on `GevorgGal/jobworker` [→](https://github.com/GevorgGal/jobworker/pull/2#discussion_r2980152151)

> Nit: `wait()` closes `j.done` after releasing the lock, so it is possible to (briefly) observe `JobStatusExited` before `j.done` is closed.
>
> ```suggestion
> // Done returns a channel that is closed after the job has fully terminated.
> ```

**@rosstimothy** on `GevorgGal/jobworker` [→](https://github.com/GevorgGal/jobworker/pull/2#discussion_r2982827357)

> Should this goroutine be added to the wait group created below?

**@rosstimothy** on `GevorgGal/jobworker` [→](https://github.com/GevorgGal/jobworker/pull/2#discussion_r2982832836)

> Optional suggestion: You can modernize this code.
>
> ```suggestion
> 	for i := range numReaders {
> 		wg.Go(func() {
>
> 			var collected []byte
> 			offset := 0
>
> 			for {
> 				data, done, ch := buf.ReadFrom(offset)
> 				if len(data) > 0 {
> 					collected = append(collected, data...)
> 					offset += len(data)
> 				}
> 				if done && len(data) == 0 {
> 					break
> 				}
> 				if len(data) == 0 {
> 					select {
> 					case <-ch:
> 					case <-time.After(5 * time.Second):
> 						t.Errorf("reader timed out at offset %d", offset)
> 						return
> 					}
> 				}
> 			}
>
> 			if got := string(collected); got != want {
> 				t.Errorf("reader got %q, want %q", got, want)
> 			}
> 		})
> 	}
> ```

**@Tener** on `GevorgGal/jobworker` [→](https://github.com/GevorgGal/jobworker/pull/2#discussion_r3001264840)

> Nit:
>
> ```suggestion
> 		end := min(len(b.data), offset+maxReadSize)
> ```

**@Tener** on `GevorgGal/jobworker` [→](https://github.com/GevorgGal/jobworker/pull/2#discussion_r3001281389)

> Might be worth noting that unlike `OutputBuffer`, actual `OutputReader` is not thread safe due to the shared offset field.

**@rosstimothy** on `GevorgGal/jobworker` [→](https://github.com/GevorgGal/jobworker/pull/3#discussion_r3003181189)

> If the files are overwritten, then how is this command idempotent? Should git report any differences if I run this command and the certs directory is already populated?

**@Tener** on `GevorgGal/jobworker` [→](https://github.com/GevorgGal/jobworker/pull/3#discussion_r3010523828)

> This is technically correct, as roleUnknown is exactly zero, but it looks rather  unusual. Returning `roleUnknown` directly would be an easy way to improve readability.

**@rosstimothy** on `GevorgGal/jobworker` [→](https://github.com/GevorgGal/jobworker/pull/4#discussion_r3018982083)

> There is a chance for a subtle race between a second signal on `sigCh` and closing `done`. Should we call signal.Stop prior to closing done so we always hit the graceful shutdown complete branch?

**@rosstimothy** on `GevorgGal/jobworker` [→](https://github.com/GevorgGal/jobworker/pull/4#discussion_r3018997271)

> Do we need both collectOutput and awaitExit? The two functions are almost identical.

**@Tener** on `GevorgGal/jobworker` [→](https://github.com/GevorgGal/jobworker/pull/4#discussion_r3020414003)

> The design stated literally:
>
> > _TODO: Graceful server shutdown — stop accepting new RPCs, drain active streams with a timeout, then exit._
>
> I agree that we are still short of "full graceful shutdown", and yet we are in contradiction with the design spec.
>
> It is rather surprising we are not trying to stop the jobs either, given all the effort to stop the RPCs. There are some elegant ways we could achieve that with minimal effort.

**@tigrato** on `gstelang/job-worker-service` [→](https://github.com/gstelang/job-worker-service/pull/1#discussion_r1735965512)

> Can you add details on how log streaming will work?
> I didn't saw any reference to it in this doc

**@strideynet** on `gstelang/job-worker-service` [→](https://github.com/gstelang/job-worker-service/pull/1#discussion_r1736026530)

> A little more detail would be helpful to ensure that your design will:
> - not degrade if one consumer of the stream is slow
> - that consumers will be caught up if they connect after logs have already been omitted

**@awly** on `hashsequence/Linux-Job-Worker` [→](https://github.com/hashsequence/Linux-Job-Worker/pull/2#discussion_r478686656)

> if you need to create the log folder before starting the command, you won't have a PID until the command is already running and writing output.
> instead, I recommend using some other _unique_ (note: PIDs are not unique) identifier for the jobs.

**@awly** on `hashsequence/Linux-Job-Worker` [→](https://github.com/hashsequence/Linux-Job-Worker/pull/2#discussion_r478690860)

> If I understand correctly:
> - the client opens connects to the server
> - the client sends StartRequest
> - the server starts the command and sends 1st StartResponse
> - the client receives StartResponse and waits, keeping the connection open
> - the command finishes and server sends a 2nd StartResponse
> - the client receives StartResponse, prints output and closes the connection
> 
> If my understanding is correct, let's change this so that client doesn't need to wait with an idle connection until the command is complete.
> As a user, I should be able to kick off multiple long-running commands in parallel and fetch their results asynchronously whenever I want.

**@andreiko** on `ilyazz/jobs` [→](https://github.com/ilyazz/jobs/pull/2#discussion_r1001124526)

> Would it be possible for a client to read the entire output of a job if it ends while the client is still reading?

**@andreiko** on `ilyazz/jobs` [→](https://github.com/ilyazz/jobs/pull/2#discussion_r1001159041)

> Alright. The closing of the channel when a job ends got me worried that the channel might not have enough buffer to fit all the remaining chunks when it's time to close it.

**@andreiko** on `ilyazz/jobs` [→](https://github.com/ilyazz/jobs/pull/2#discussion_r1001173304)

> Sounds good. Let's make sure to test this scenario (slow reader, job ends quickly after producing large output) when the implementation is ready.

**@sclevine** on `ilyazz/jobs` [→](https://github.com/ilyazz/jobs/pull/2#discussion_r1001275192)

> Is there a more idiomatic type than `<-chan []byte` that you could use for streaming logs in Go?

**@sclevine** on `ilyazz/jobs` [→](https://github.com/ilyazz/jobs/pull/2#discussion_r1001290452)

> If I want to gracefully terminate the server (e.g., by sending it SIGTERM or SIGINT), do I need to manually cleanup cgroups on the host?

**@jimbishopp** on `ilyazz/jobs` [→](https://github.com/ilyazz/jobs/pull/2#discussion_r1001925351)

> Questions about this: 
> - Is a new byte slice allocated for every send to the channel?
> - If the job has finished executing, does the channel only contain a single byte slice with all the output?
> - Can streaming output be accomplished without follow mode (2 modes)?

**@sclevine** on `ilyazz/jobs` [→](https://github.com/ilyazz/jobs/pull/2#discussion_r1002845437)

> How do I close the file descriptor that is used to read the log file?

**@andreiko** on `ilyazz/jobs` [→](https://github.com/ilyazz/jobs/pull/4#discussion_r1007433771)

> nit: would be nice to have a readme entry listing which tools and which plugins are required on the developer machine

**@andreiko** on `ilyazz/jobs` [→](https://github.com/ilyazz/jobs/pull/3#discussion_r1007480280)

> Would it actually be bad if the log file was deleted with some readers still active?

**@sclevine** on `ilyazz/jobs` [→](https://github.com/ilyazz/jobs/pull/3#discussion_r1008486990)

> > Here I assume our command is a good citizen, and waits for its children, especially if it supports graceful exit.
> 
> Exactly 🙂

**@rosstimothy** on `jaylane/job-scheduler` [→](https://github.com/jaylane/job-scheduler/pull/1#discussion_r1714213505)

> Can you provide more details on how streaming will work? Where is the output being stored?

**@codingllama** on `jaylane/job-scheduler` [→](https://github.com/jaylane/job-scheduler/pull/1#discussion_r1714331671)

> Either works, but my 2c is to keep the proto file. No need to undo what's already done.

**@tigrato** on `joshuarubin/teleport-job-worker` [→](https://github.com/joshuarubin/teleport-job-worker/pull/3#discussion_r1743551952)

> Can you please indicate more details on how the library will inform the reader that the stream ended and how it will also inform that new data was written to stderr or stdout?

**@rosstimothy** on `joshuarubin/teleport-job-worker` [→](https://github.com/joshuarubin/teleport-job-worker/pull/3#discussion_r1743924763)

> > will block until there is either more output to return or the job ends in which case `io.EOF` is returned.
> 
> How will Read be unblocked when more output is available?

**@rosstimothy** on `joshuarubin/teleport-job-worker` [→](https://github.com/joshuarubin/teleport-job-worker/pull/3#discussion_r1743978517)

> The Reader is fine, I was just curious what your plan was, but I don't want to get too lost in implementation details here.

**@rosstimothy** on `joshuarubin/teleport-job-worker` [→](https://github.com/joshuarubin/teleport-job-worker/pull/3#discussion_r1744472498)

> Might also be worth mentioning that if the user closes the ReadCloser returned from `JobOutput` it will also unblock any read operations and terminate the reading goroutine.
> 
> ```suggestion
> Internally, the library will maintain a buffer containing the job output. The job itself will be configured to use this buffer for `stdout` and `stderr` and the `Write()` method will be goroutine safe. When a client calls `JobOutput()` the library will create a new `io.ReadCloser` that independently, and with goroutine safety, reads through to the end of the output buffer. If the job has already completed at this time, `io.EOF` will be returned. If not, subsequent calls to `Read()` will block until there is either more output to return or the job ends in which case `io.EOF` is returned. When the job calls `Write()` it will signal connected readers that new data has been written to the output buffer. The readers will in turn be able to unblock and return their `Read()` calls to connected clients at that time.
> ```

**@codingllama** on `joshuarubin/teleport-job-worker` [→](https://github.com/joshuarubin/teleport-job-worker/pull/3#discussion_r1745574300)

> > If the user closes the `io.ReadCloser`, any pending `Read()` operations will unblock with `io.EOF` and any goroutines and other associated resources with the reader will be released.
> 
> nit: I would probably do an ErrClosed (or something like that) so callers can distinguish closed from a "true" EOF.

**@tigrato** on `joshuarubin/teleport-job-worker` [→](https://github.com/joshuarubin/teleport-job-worker/pull/4#discussion_r1746724188)

> What will happen if a reader is slow to consume the output?

**@tigrato** on `joshuarubin/teleport-job-worker` [→](https://github.com/joshuarubin/teleport-job-worker/pull/4#discussion_r1746733434)

> Why do we need to keep a full copy of the job output per reader?
> 
> If a job outputs 1Gb of data and we have 5 active readers, we will take 5GB of ram. Is there a way where we do not need to implement a solution with such bottleneck ?

**@tigrato** on `joshuarubin/teleport-job-worker` [→](https://github.com/joshuarubin/teleport-job-worker/pull/4#discussion_r1746735279)

> why closing the reader here?

**@tigrato** on `joshuarubin/teleport-job-worker` [→](https://github.com/joshuarubin/teleport-job-worker/pull/4#discussion_r1746737300)

> if `dataAvail` is only consumed in a single place, why do you need to close (broadcast) info everywhere?

**@rosstimothy** on `joshuarubin/teleport-job-worker` [→](https://github.com/joshuarubin/teleport-job-worker/pull/4#discussion_r1747536697)

> Is it possible for the final stdout/stderr write of the process to still be in flight when the jobDone channel is closed? If so wouldn't that cause data loss?

**@rosstimothy** on `joshuarubin/teleport-job-worker` [→](https://github.com/joshuarubin/teleport-job-worker/pull/4#discussion_r1747542043)

> Will this return the correct number of bytes read? If we've gotten into this case it's because we already read from the buffer once and received an EOF. However, we could have read _some_ data from the buffer the first time.

**@tigrato** on `joshuarubin/teleport-job-worker` [→](https://github.com/joshuarubin/teleport-job-worker/pull/4#discussion_r1747561903)

> It will also return io.EOF and mislead the client because we ended up the buffer but the job hasn't ended

**@codingllama** on `joshuarubin/teleport-job-worker` [→](https://github.com/joshuarubin/teleport-job-worker/pull/4#discussion_r1747645928)

> nit: Drop the assertion, usage in code should already guarantee conformance to the interface.

**@codingllama** on `joshuarubin/teleport-job-worker` [→](https://github.com/joshuarubin/teleport-job-worker/pull/4#discussion_r1747656865)

> nit:
> 
> ```suggestion
> 				slog.Error("Error writing to Reader buffer", "err", err)
> ```

**@codingllama** on `joshuarubin/teleport-job-worker` [→](https://github.com/joshuarubin/teleport-job-worker/pull/4#discussion_r1747658578)

> Why is this a field on Reader?

**@codingllama** on `joshuarubin/teleport-job-worker` [→](https://github.com/joshuarubin/teleport-job-worker/pull/4#discussion_r1747663882)

> Should this be a slice of Channels instead?

**@codingllama** on `joshuarubin/teleport-job-worker` [→](https://github.com/joshuarubin/teleport-job-worker/pull/4#discussion_r1747680315)

> Assert the content read?

**@codingllama** on `joshuarubin/teleport-job-worker` [→](https://github.com/joshuarubin/teleport-job-worker/pull/4#discussion_r1747683216)

> What is this block testing? `done` is already closed and `n` and `err` are leftovers, or am I missing something?

**@codingllama** on `joshuarubin/teleport-job-worker` [→](https://github.com/joshuarubin/teleport-job-worker/pull/4#discussion_r1750840229)

> I thought tooling was past that already, my bad.

**@codingllama** on `joshuarubin/teleport-job-worker` [→](https://github.com/joshuarubin/teleport-job-worker/pull/4#discussion_r1750843514)

> I meant `[]*safereader.Channels`.

**@rosstimothy** on `joshuarubin/teleport-job-worker` [→](https://github.com/joshuarubin/teleport-job-worker/pull/4#discussion_r1751053562)

> While there may be no concerns with clients preventing others from consuming output, couldn't a very verbose process increase the likelihood, frequency, and duration of stuttering? In such a scenario where the process produces a large amount of output faster than a reader may consume the contention on the lock could increase to the point where it may be impossible for the reader to hold the lock.

**@tigrato** on `joshuarubin/teleport-job-worker` [→](https://github.com/joshuarubin/teleport-job-worker/pull/4#discussion_r1752209252)

> To use a channel as signal, you don't need to close the channel. You just need to send data to it
> 
> `r.dataAvail <-struct{}{}` given that it's only consumed in the Read function so you don't need to close it every time. 
> 
> If you don't want the channel send to be a block operation, use a buffered channel

**@greedy52** on `kiakeshmiri/process-runner` [→](https://github.com/kiakeshmiri/process-runner/pull/1#discussion_r1765671546)

> could you provide some more details on the streaming part? 
> 1. how to transition the client from reading the log history to waiting on the channel?
> 2. what happens when the process is completed while streaming? (i see you mention in the edge case but it is missing details on how to handle it properly?)
> 3. what happens when multiple concurrent clients streaming log from the same process?

**@tigrato** on `kiakeshmiri/process-runner` [→](https://github.com/kiakeshmiri/process-runner/pull/1#discussion_r1766509283)

> ```suggestion
> func (so *outputLogs) Write(p []byte) (n int, err error) {
> ```

**@tigrato** on `kiakeshmiri/process-runner` [→](https://github.com/kiakeshmiri/process-runner/pull/1#discussion_r1766512720)

> what happens if the reader won't read from the channel?

**@tigrato** on `kkloberdanz/teleport-challenge` [→](https://github.com/kkloberdanz/teleport-challenge/pull/1#discussion_r2781987125)

> How each reader will track his position and how can we cancel a subscription?
>
> Can you please include the subscriber interface?

**@tigrato** on `kkloberdanz/teleport-challenge` [→](https://github.com/kkloberdanz/teleport-challenge/pull/7#discussion_r2812000945)

> l.cmd.Wait() can only be called once. Should we document that behaviour?

**@tigrato** on `kkloberdanz/teleport-challenge` [→](https://github.com/kkloberdanz/teleport-challenge/pull/8#discussion_r2821944779)

> you never call close send. this means the stream will be kept open if the server never cancels it

**@tigrato** on `kkloberdanz/teleport-challenge` [→](https://github.com/kkloberdanz/teleport-challenge/pull/8#discussion_r2832821022)

> This will notify the server the client is no longer interested. the server cancels the stream context causing the log output handler to halt and return without sending all the logs

**@tigrato** on `kkloberdanz/teleport-challenge` [→](https://github.com/kkloberdanz/teleport-challenge/pull/8#discussion_r2832838022)

> Should we wait for the complete output? The `echo client-stream` will be included in the command, but we won’t know whether the system actually received the command output

**@tigrato** on `kkloberdanz/teleport-challenge` [→](https://github.com/kkloberdanz/teleport-challenge/pull/8#discussion_r2832885969)

> if the client disconnected, shouldn't we send io.EOF?
> from the client perspective, he sent a "disconnect, I am no longer interested"

**@tigrato** on `kkloberdanz/teleport-challenge` [→](https://github.com/kkloberdanz/teleport-challenge/pull/8#discussion_r2835836699)

> I forgot this is a rpc to stream method. CloseSend is used to inform the server the client sent everything. In this case, CloseSend is used to signal the server we only sent a single message.
> 
> In this case, the way of releasing the gRPC stream is not by sending close but by cancelling the context which we don't do. 
> 
> For instance, if ` w.Write(resp.GetData())` fails, we never properly release the gRPC stream and return. We assume the caller context to be cancelled but we have no guarantee that it will ever happen.
> 
> We should also wait for the server to release the rpc before returning

**@tigrato** on `kkloberdanz/teleport-challenge` [→](https://github.com/kkloberdanz/teleport-challenge/pull/8#discussion_r2835838668)

> when I say wait for the complete output is to ensure we have `client-stream` twice in our output. One from the command itself and other from the command being executed. 
> 
> currently you only check if there is one entry but this flow doesn't test if the output is complete, i,e, 
> 
> ```
> $ echo client-stream 
> client-stream 
> ```

**@zmb3** on `kovyrin/teleport-exec` [→](https://github.com/kovyrin/teleport-exec/pull/1#discussion_r788222817)

> Similar note here. A client should be able to stream the output as many times as they'd like, so let's not delete the file when it may still be needed.

**@zmb3** on `kovyrin/teleport-exec` [→](https://github.com/kovyrin/teleport-exec/pull/1#discussion_r788223727)

> Could you add a couple sections to solution details? In particular, I'm interested in more information on:
> 
> - the streaming solution
> - the process lifecycle (how will the server/library handle starting a job and setting up the isolation)

**@zmb3** on `kovyrin/teleport-exec` [→](https://github.com/kovyrin/teleport-exec/pull/1#discussion_r788238842)

> Yeah if you can use the certs you already have I think that will be cleaner. 🙂

**@russjones** on `kovyrin/teleport-exec` [→](https://github.com/kovyrin/teleport-exec/pull/1#discussion_r788285770)

> You don't need to implement this as a choice, you can always assume stream until the command stops.

**@rosstimothy** on `kovyrin/teleport-exec` [→](https://github.com/kovyrin/teleport-exec/pull/4#discussion_r794920428)

> What will happen if you close a stream more than once?

**@rosstimothy** on `kovyrin/teleport-exec` [→](https://github.com/kovyrin/teleport-exec/pull/4#discussion_r794922752)

> So to close a `FileStream` do I call `Close` or send something via the `Done` channel?

**@zmb3** on `kovyrin/teleport-exec` [→](https://github.com/kovyrin/teleport-exec/pull/4#discussion_r794927733)

> This stutters a bit. The package is already named `file_stream`, so anyone referring to this from outside this package will have to type `file_stream.NewFileStream()`.
> 
> I would go with something more like `filestream.New()`

**@zmb3** on `kovyrin/teleport-exec` [→](https://github.com/kovyrin/teleport-exec/pull/4#discussion_r794930481)

> You're so close to implementing `io.Reader` here, which would give you compatibility and composability with a bunch of Go APIs. Worth considering, since you should really be returning the error anyways.

**@rosstimothy** on `kovyrin/teleport-exec` [→](https://github.com/kovyrin/teleport-exec/pull/4#discussion_r795072342)

> I'm not sure you want to expose this publicly. Doing so could do two things:
> 1) allow a caller to spawn an endless number of goroutines
> 2) cause `Read` to potentially miss events from `fsnotify`
> 
> ```suggestion
> func (s *FileStream) waitForChanges() (changed bool) {
> ```

**@rosstimothy** on `kovyrin/teleport-exec` [→](https://github.com/kovyrin/teleport-exec/pull/4#discussion_r795072685)

> If this goroutine is truly needed due to limitations in `fsnotify`, perhaps `FileStream` could own the `result` chan and just launch one goroutine that monitors the events in `New`. Then in `Read` you could just block on reading from the channel instead. Thoughts?

**@rosstimothy** on `kovyrin/teleport-exec` [→](https://github.com/kovyrin/teleport-exec/pull/4#discussion_r795809414)

> Seems like a good place for a read lock

**@zmb3** on `kovyrin/teleport-exec` [→](https://github.com/kovyrin/teleport-exec/pull/2#discussion_r795812509)

> The `ioutil` package was deprecated. Prefer `os` instead.
> 
> ```suggestion
> 		err := os.WriteFile(path, []byte(data), mode)
> ```

**@zmb3** on `kovyrin/teleport-exec` [→](https://github.com/kovyrin/teleport-exec/pull/2#discussion_r795822320)

> This is an inefficient way to wait. Can you think of another approach?

**@zmb3** on `kovyrin/teleport-exec` [→](https://github.com/kovyrin/teleport-exec/pull/2#discussion_r795830729)

> `Kill()` already checks if `c.Running()`, so there's no need to check twice.

**@rosstimothy** on `kovyrin/teleport-exec` [→](https://github.com/kovyrin/teleport-exec/pull/4#discussion_r795837395)

> I think you probably want to propagate the `io.EOF` here and let the user call `Close` themselves. That behavior would be more consistent with other `io.Reader` implementations.

**@nklaassen** on `kovyrin/teleport-exec` [→](https://github.com/kovyrin/teleport-exec/pull/2#discussion_r796203098)

> `isClosed` is protected by `readersLock`?

**@nklaassen** on `kovyrin/teleport-exec` [→](https://github.com/kovyrin/teleport-exec/pull/2#discussion_r796212989)

> but `waitForCompletion` will run in a separate goroutine, is there any point to this?

**@GavinFrazar** on `kurczynski/teleport-job-worker` [→](https://github.com/kurczynski/teleport-job-worker/pull/1#discussion_r1610796882)

> fetching/listing the jobs owned by a user is out of scope, you can get rid of this.
> The API only needs to allow a user to start, stop, stream output, and get the status of an individual job.

**@GavinFrazar** on `kurczynski/teleport-job-worker` [→](https://github.com/kurczynski/teleport-job-worker/pull/1#discussion_r1610804798)

> Some more detail here would be nice - how will concurrent output be handled? How will you stream the output from process start? What happens when the job is no longer running?

**@bernardjkim** on `kurczynski/teleport-job-worker` [→](https://github.com/kurczynski/teleport-job-worker/pull/1#discussion_r1613781347)

> Looks like tokens are no longer going to be used. Might want to remove the `-token` flags from the UX design.

**@GavinFrazar** on `kurczynski/teleport-job-worker` [→](https://github.com/kurczynski/teleport-job-worker/pull/1#discussion_r1614146199)

> last thing, but I won't block the approval over this detail:
> 
> with this setup the server and client certs are both CA certs
> 
> ```fish
> $ openssl x509 -in client-cert.pem -noout -text | grep CA
>                 CA:TRUE
> ```
> 
> This won't be secure for mTLS - let's update it so the client and server certs are not CAs

**@sclevine** on `kurczynski/teleport-job-worker` [→](https://github.com/kurczynski/teleport-job-worker/pull/2#discussion_r1621429603)

> This does not appear to support multiple simultaneous clients each streaming from the start of execution.

**@sclevine** on `kurczynski/teleport-job-worker` [→](https://github.com/kurczynski/teleport-job-worker/pull/2#discussion_r1621524093)

> You can use as much memory or disk as you need. It's okay for the output of an individual job to grow unbounded, as long as the output is cleaned up when its no longer useful.
> 
> The API service does not need to scale, but we will evaluate the efficiency of the output streaming solution. I would think carefully about the ideal structure for solving this problem.

**@bernardjkim** on `kurczynski/teleport-job-worker` [→](https://github.com/kurczynski/teleport-job-worker/pull/2#discussion_r1624964328)

> I think this will need to be addressed before you can work on streaming.

**@greedy52** on `m3talsmith/jobberthehut` [→](https://github.com/m3talsmith/jobberthehut/pull/1#discussion_r2499178968)

> for this challenge, it is expected that `tail` to be the default behavior. 
> ```
>     rpc StopJob(StopJobRequest) returns (StopJobResponse);
> ``` 
> how would this api handle tailing/streaming the output?

**@greedy52** on `m3talsmith/jobberthehut` [→](https://github.com/m3talsmith/jobberthehut/pull/1#discussion_r2499197965)

> could you provide more details on version of the cgroup and the controls/limits will be used?

**@greedy52** on `m3talsmith/jobberthehut` [→](https://github.com/m3talsmith/jobberthehut/pull/1#discussion_r2499205972)

> could you provide more details on the handling of the logs. In particular, we are interested in how concurrent log-tailing clients will be handled.

**@greedy52** on `m3talsmith/jobberthehut` [→](https://github.com/m3talsmith/jobberthehut/pull/1#discussion_r2499231586)

> correct me if i am wrong. i don't find any details on authorization like who can do what.

**@rosstimothy** on `m3talsmith/jobberthehut` [→](https://github.com/m3talsmith/jobberthehut/pull/1#discussion_r2500042011)

> Same sentiment as above, tailing output should be the default and only behavior for this challenge.

**@greedy52** on `m3talsmith/jobberthehut` [→](https://github.com/m3talsmith/jobberthehut/pull/1#discussion_r2500061658)

> first of all, the streaming log should return all logs regardless of whether the job is running or stopped.
> 
> what I mean originally, is how this gRPC can handle the tailing functionality since it's just a single API call?
> ```
> rpc StopJob(StopJobRequest) returns (StopJobResponse);
> ```

**@gabrielcorado** on `m3talsmith/jobberthehut` [→](https://github.com/m3talsmith/jobberthehut/pull/1#discussion_r2500544262)

> Just as reference, the standard error model does support an optional error message. 
> 
> > If an error occurs, gRPC returns one of its error status codes instead, with an optional string error message that provides further details about what happened. -- https://grpc.io/docs/guides/error/#standard-error-model
> 
> Would that be enough for the error reporting you're going to have?
> 
> > For example, if there was an exit other than 0, or a permissions issue, then I could provide additional information on what the error was.
> 
> Unrelated to the error handling here, but does that mean that this gRPC call will block until the process is completed? Based on the naming (`StartJob`), I would expect it to just start the command and return immediately. Is that the case?
> 
> If so, what would happen if this process takes a little longer (e.g., due to the cgroup setup) and the request gets canceled? What would happen then?

**@greedy52** on `MarkDHarris/teleport-systems-challenge` [→](https://github.com/MarkDHarris/teleport-systems-challenge/pull/1#discussion_r3016929999)

> using `sync.Cond` is definitely cool. otoh, "maintain and iterate a subscriber list" won't necessarily "block the writer on a slow consumer", if implemented correctly. just a random comment no action required.

**@rosstimothy** on `MarkDHarris/teleport-systems-challenge` [→](https://github.com/MarkDHarris/teleport-systems-challenge/pull/2#discussion_r3035838095)

> It might be worthwhile to separate the responsibilities of output management to another object.

**@rosstimothy** on `MarkDHarris/teleport-systems-challenge` [→](https://github.com/MarkDHarris/teleport-systems-challenge/pull/2#discussion_r3039688731)

> Is there a more idiomatic way to express this that doesn't rely on callbacks or channels?

**@rosstimothy** on `MarkDHarris/teleport-systems-challenge` [→](https://github.com/MarkDHarris/teleport-systems-challenge/pull/2#discussion_r3039697477)

> These provide little value, and also exposes writing output to anyone that can get a handle to the job - which we probably don't want to do. Can we interact with the outputBuffer API directly instead of indirectly via these two functions?
> ```suggestion
> ```

**@greedy52** on `MarkDHarris/teleport-systems-challenge` [→](https://github.com/MarkDHarris/teleport-systems-challenge/pull/2#discussion_r3039946833)

> it is a little unclear how `ctx` is handled. for example, lets say the client cancels the stream while still streaming historical data, what happens?
> 
> From the library perspective, the caller does pass in a `ctx` but `fn` may or may not be handling for the same `ctx`. IMHO,  library should at the minimum document the expectations, or better to be designed in a way that avoid problem like this in the first place.

**@greedy52** on `MarkDHarris/teleport-systems-challenge` [→](https://github.com/MarkDHarris/teleport-systems-challenge/pull/2#discussion_r3041601755)

> minor suggestion: do you really need two channels? also, they can be moved out of the go routine too to avoid this weird channel handshake.

**@rosstimothy** on `MarkDHarris/teleport-systems-challenge` [→](https://github.com/MarkDHarris/teleport-systems-challenge/pull/2#discussion_r3045000566)

> Should the output buffer care about the distinction?

**@rosstimothy** on `MarkDHarris/teleport-systems-challenge` [→](https://github.com/MarkDHarris/teleport-systems-challenge/pull/2#discussion_r3045083825)

> Does this test provide any value if we have other tests that validate reading the output matches what was written? Getting rid of it would also allow us to remove `bytesLen` and clean up the output buffer API.

**@rosstimothy** on `MarkDHarris/teleport-systems-challenge` [→](https://github.com/MarkDHarris/teleport-systems-challenge/pull/2#discussion_r3045122208)

> This function might be more resilient and cleaner if we deferred closing the channels
> ```suggestion
> 	go func() {
> 		r := b.newReader(context.Background())
> 		defer r.Close()
> 		defer close(errCh)
> 		defer close(received)
>
> 		buf := make([]byte, 256)
> 		for {
> 			n, err := r.Read(buf)
> 			if n > 0 {
> 				received <- string(buf[:n])
> 			}
> 			if errors.Is(err, io.EOF) {
> 				break
> 			}
> 			if err != nil {
> 				errCh <- err
> 				return
> 			}
> 		}
> ```

**@rosstimothy** on `MarkDHarris/teleport-systems-challenge` [→](https://github.com/MarkDHarris/teleport-systems-challenge/pull/2#discussion_r3045135530)

> Why do we need to close the channel here?

**@tigrato** on `mcampo84/teleport_challenge` [→](https://github.com/mcampo84/teleport_challenge/pull/1#discussion_r1910147420)

> How will this work with reading the existing output (start from the beginning of the log)?
> Can you expand a bit in this section to include more details? It looks like a bit sparse in terms of details

**@tigrato** on `mcampo84/teleport_challenge` [→](https://github.com/mcampo84/teleport_challenge/pull/1#discussion_r1915098129)

> If a client subscribed for notifications but is slow to consume them, can that client impact other client ability to consume the logs given this solution?
> Even if the channel is buffered, there is always a limit after which the info above can impact other clients
> 
> Is there a more resilient way of implementing this with channels?/Can other sync primitives be used for this purpose?

**@tigrato** on `mcampo84/teleport_challenge` [→](https://github.com/mcampo84/teleport_challenge/pull/1#discussion_r1915236314)

> Use a io.Read function. We do not need buffering

**@tigrato** on `mcampo84/teleport_challenge` [→](https://github.com/mcampo84/teleport_challenge/pull/5#discussion_r1920266604)

> being this publicly exposed, a user can call it multiple times which will cause panics because multiple channel close operations

**@tigrato** on `mcampo84/teleport_challenge` [→](https://github.com/mcampo84/teleport_challenge/pull/5#discussion_r1920267207)

> closing channels before finishing the job seems to hide important logs

**@tigrato** on `mcampo84/teleport_challenge` [→](https://github.com/mcampo84/teleport_challenge/pull/5#discussion_r1920268573)

> what if one of the channels is full or the reader stop reading the data?

**@fspmarshall** on `mcampo84/teleport_challenge` [→](https://github.com/mcampo84/teleport_challenge/pull/5#discussion_r1920404086)

> Log buffer should not be cleared, otherwise performing a start followed by trying to read output is a race condition.

**@fspmarshall** on `mcampo84/teleport_challenge` [→](https://github.com/mcampo84/teleport_challenge/pull/5#discussion_r1920407629)

> Closing log channels here is a race condition as `logOutput` may still be trying to send against them.

**@fspmarshall** on `mcampo84/teleport_challenge` [→](https://github.com/mcampo84/teleport_challenge/pull/5#discussion_r1920411417)

> Sending on the streamer while the lock is held may produce unacceptable contention.  According to the docs of the `OutputStreamer` interface, `OutputStreamer` is allowed to be backed by a GRPC stream.  In that case, this would be doing I/O under lock.
> 
> If you want to `Send` under lock, then `OutputStreamer` needs a send method that guarantees that that it won't block, otherwise you can end up hanging here with the lock held for long periods of time.

**@fspmarshall** on `mcampo84/teleport_challenge` [→](https://github.com/mcampo84/teleport_challenge/pull/5#discussion_r1920413642)

> Exiting as soon as `DoneChannel` is closed is a race since it can't be guaranteed that there isn't additional data available via `logChannel`.  The entire output of `logChannel` needs to be consumed before returning.
> 
> Any early-return that might happen before the full output has been sent should be an error return.

**@tigrato** on `mcampo84/teleport_challenge` [→](https://github.com/mcampo84/teleport_challenge/pull/5#discussion_r1925389477)

> this code can be concurrently called (check your function calls) and potentially panic with the close being called several times for the same channel when given the correct order and time of execution

**@tigrato** on `mikewurtz/taskman` [→](https://github.com/mikewurtz/taskman/pull/1#discussion_r2027187117)

> Can you describe how this channel system will handle slow consumers?

**@rosstimothy** on `mikewurtz/taskman` [→](https://github.com/mikewurtz/taskman/pull/1#discussion_r2027219735)

> > The task must have been created by the user in order for the server to return the task details
> 
> Is this accurate? Shouldn't the admin be able to retrieve the status for any task?

**@rosstimothy** on `mikewurtz/taskman` [→](https://github.com/mikewurtz/taskman/pull/1#discussion_r2027225884)

> > or is not associated with the use
> 
> Same question as above, shouldn't the admin be allowed to stream the task output?

**@rosstimothy** on `mikewurtz/taskman` [→](https://github.com/mikewurtz/taskman/pull/1#discussion_r2027230409)

> Should there also be a flag for setting the address of the server, with a sensible default?

**@tigrato** on `mikewurtz/taskman` [→](https://github.com/mikewurtz/taskman/pull/1#discussion_r2031548867)

> Imagine a program that writes one payload every millisecond. The client, however, is too slow to keep up, so its read position lags behind, and the channel eventually fills up. Once the client catches up and consumes all messages from the channel, the program stops writing because there's no new data. From the client's perspective, it appears to be caught up and in a "live" state—until a new write happens - imagine it takes several hours. At that point, the server realizes the client is actually far behind (based on its read position), so it pushes a burst of data to help it catch up. This causes the client to lag again, and the cycle repeats.
> 
> Now, if we already maintain the read position for each user, is it necessary to continue writing payloads to the channel and consuming memory until it fills up?
> 
> Let’s say the interface is something like `chan []byte`. If you allocate a buffer of 1024 messages per client, and each message is ~1024 bytes, then each client consumes ~1MB of memory. Scale that to many clients or larger messages, and it becomes easy for an attacker to overwhelm the server by inflating payload sizes and client counts.
> 
> Is there a better approach where we can leverage each client's read position to avoid sending data into a channel until it's truly needed—*and* without having to wait for a new write to detect laggy clients?

**@tigrato** on `mikewurtz/taskman` [→](https://github.com/mikewurtz/taskman/pull/1#discussion_r2031851100)

> That approach sounds way better to me. 
>
> Do you need a goroutine per stream if the gRPC server already spawns a goroutine for you?

**@rosstimothy** on `mikewurtz/taskman` [→](https://github.com/mikewurtz/taskman/pull/1#discussion_r2032319521)

> Isn't the gRPC stream handler running in its own goroutine already? Why does another goroutine need to be launched just to forward data to the stream?

**@tigrato** on `mikewurtz/taskman` [→](https://github.com/mikewurtz/taskman/pull/2#discussion_r2035474768)

> if you receive a context, you can cancel the stream once the context is cancelled by ctrl+c

**@timothyb89** on `mikewurtz/taskman` [→](https://github.com/mikewurtz/taskman/pull/4#discussion_r2040093626)

> Is a `sync.Once` needed here if only a single goroutine will ever try to close the channel?

**@tigrato** on `mikewurtz/taskman` [→](https://github.com/mikewurtz/taskman/pull/4#discussion_r2040100195)

> Should we consider moving authorization to the gRPC level? Does it make sense to validate the owner of the job so far away from the gRPC server where these details are defined?

**@tigrato** on `mikewurtz/taskman` [→](https://github.com/mikewurtz/taskman/pull/5#discussion_r2045177465)

> What if you wake up the sync.Cond when a client wants to leave. Check context.AfterFunc .
>
> That will also allow you to rewrite the streamer as an io.ReadCloser interface

**@tigrato** on `mikewurtz/taskman` [→](https://github.com/mikewurtz/taskman/pull/5#discussion_r2047025295)

> When I said I wanted to transform the streamer into an `io.ReadCloser`, I meant refactoring the entire call:
> 
> ```go
> err = s.taskManager.StreamTaskOutput(stream.Context(), req.TaskId, writer)
> ```
> 
> into something that behaves like this:
> 
> ```go
> jobStreamer, err := s.taskManager.StreamTaskOutput(ctx, req.TaskId)
> 
> context.AfterFunc(ctx, func() { jobStreamer.Close() })
> 
> for {
>     n, err := jobStreamer.Read(buf)
>     if errors.Is(err, io.EOF) {
>         // ...
>     } else if err != nil {
>         // ...
>     }
> 
>     err := stream.Send(&pb.StreamTaskOutputResponse{
>         Output: slices.Clone(buf[:n]),
>     })
>     // ...
> }
> ```
> 
> This design allows the caller to treat the result as a regular `io.Reader`, making it compatible with utilities like `io.Copy(os.Stdout, jobStreamer)` or `io.ReadAll(jobStreamer)` for more flexible consumption.

**@rosstimothy** on `mikewurtz/taskman` [→](https://github.com/mikewurtz/taskman/pull/5#discussion_r2047083583)

> What ownership guarantees does gRPC make about data sent over a stream?

**@oeric** on `moalf/teleport` [→](https://github.com/moalf/teleport/pull/1#discussion_r2594175718)

> Can you provide an outline of the directory structure of your implementation, and some more details on how many PRs you're planning on making + what each PR would add?

**@russjones** on `moalf/teleport` [→](https://github.com/moalf/teleport/pull/1#discussion_r2594336320)

> This looks good, but can you provide some more details in particular on this one?
> 
> > - pull requests and conversations resolutions before merging

**@russjones** on `moalf/teleport` [→](https://github.com/moalf/teleport/pull/1#discussion_r2594340052)

> In my view, MFA satisfies the following requirement of the challenge, no need to implement the others to save time.
> 
> > Enforce strong authentication for all web applications within tenant.
> 
> However, can you provide details on how you plan to enforce MFA?

**@nklaassen** on `MrChristianL/Teleport-Job-Worker-Service` [→](https://github.com/MrChristianL/Teleport-Job-Worker-Service/pull/1#discussion_r2879009292)

> I'm curious in what scenario a client may read a length header before the chunk of data has been committed to the log, specifically to see if there may be an issue with the client "missing" a `cond.Broadcast` and sleeping indefinitely. A sync.Mutex is always used along with a sync.Cond but I don't see it mentioned here. Could you describe this scenario more completely including which component holds the mutex at which times and what size of reads the client is doing?

**@nklaassen** on `MrChristianL/Teleport-Job-Worker-Service` [→](https://github.com/MrChristianL/Teleport-Job-Worker-Service/pull/1#discussion_r2879019352)

> I don't understand how context cancellation and the dedicated `done` channel relate. One question I have is, how will a client blocked on `cond.Wait()` wake up when the gRPC context cancels (when the user disconnects, for example).

**@nklaassen** on `MrChristianL/Teleport-Job-Worker-Service` [→](https://github.com/MrChristianL/Teleport-Job-Worker-Service/pull/1#discussion_r2879125415)

> I think this doesn't need to be part of the message and the stream RPC can just return nil when it's done

**@nklaassen** on `MrChristianL/Teleport-Job-Worker-Service` [→](https://github.com/MrChristianL/Teleport-Job-Worker-Service/pull/1#discussion_r2881098209)

> I don't think the write-event boundaries your library sees will be very meaningful. I don't think they necessarily correlate 1-1 to a write syscall from the job process, for example. There are probably a couple layers of buffering (depending what you set cmd.Stdout/cmd.Stderr to). In this context I don't believe "chunk-level correctness" has a well-defined meaning.
> 
> In short, I don't think inclusion of the length header matters, it's an implementation detail. If it makes implementation simpler you may keep it, if it makes it more complicated feel free to drop it.

**@eriktate** on `MrChristianL/Teleport-Job-Worker-Service` [→](https://github.com/MrChristianL/Teleport-Job-Worker-Service/pull/1#discussion_r2881195877)

> You might consider only including the job ID as output. This makes it easier to do things like start a job and stream log output:
> ```
> ./jobctl start python3 process_dataset.py | ./jobctl stream
> ```
> I do see below a `--follow` flag that appears to do the same thing, but I think it's still worthwhile to avoid the potential friction

**@eriktate** on `MrChristianL/Teleport-Job-Worker-Service` [→](https://github.com/MrChristianL/Teleport-Job-Worker-Service/pull/1#discussion_r2881207643)

> The user stories are helpful, and I can guess what `--follow` is doing here, but could you include some more specific documentation around what sub commands and flags you plan to support?

**@eriktate** on `MrChristianL/Teleport-Job-Worker-Service` [→](https://github.com/MrChristianL/Teleport-Job-Worker-Service/pull/1#discussion_r2886291489)

> Looks like you might be missing some content and/or formatting here. I also don't see anything for stopping a job and it's not clear if `GetJob()` is intended to power both getting the status and streaming output.

**@eriktate** on `MrChristianL/Teleport-Job-Worker-Service` [→](https://github.com/MrChristianL/Teleport-Job-Worker-Service/pull/1#discussion_r2892338369)

> > If we are looking to implement them as roles rather than their identities, that is no problem. It would be as simple as adding a secondary map that correlates identities to roles, and roles to permissions.
> 
> I don't think that's necessary, I was just trying to understand how many concepts we were working with and what was being mapped where. A hardcoded map of identities to permissions seems reasonable to me for this challenge
> 
> > I could also change the CLI's --role flag to something more fitting, such as --cert, since the user would be specifying choosing which certificate they would be using -- and by extention, which identity they would be running commands under -- rather than adding any sort of confusion with the current --role flag.
> 
> I think that would be clearer and more in line with what I would expect if this were meant to be a shippable CLI :+1:

**@nklaassen** on `MrChristianL/Teleport-Job-Worker-Service` [→](https://github.com/MrChristianL/Teleport-Job-Worker-Service/pull/1#discussion_r2897051897)

> thanks for the detailed explanation, it's more clear to me now

**@eriktate** on `MrChristianL/Teleport-Job-Worker-Service` [→](https://github.com/MrChristianL/Teleport-Job-Worker-Service/pull/1#discussion_r2898057965)

> nit: I disagree that `--role` is more accurate to the implementation. With identities separated from roles, you could have a cert issued to any identity with any of the supported roles (e.g. `christian.l@example.com` with the `admin` role). Since the flag is ultimately choosing which certificate to load I think `--cert` still makes the most sense.

**@eriktate** on `MrChristianL/Teleport-Job-Worker-Service` [→](https://github.com/MrChristianL/Teleport-Job-Worker-Service/pull/2#discussion_r2905546300)

> nit: considering the whole broker is meant for output streaming, this comment doesn't seem super helpful. Maybe something like:
> ```suggestion
> 	cond *sync.Cond // used to broadcast new writes to open readers
> ```

**@eriktate** on `MrChristianL/Teleport-Job-Worker-Service` [→](https://github.com/MrChristianL/Teleport-Job-Worker-Service/pull/2#discussion_r2905687588)

> Isn't this just an `io.Copy`?

**@eriktate** on `MrChristianL/Teleport-Job-Worker-Service` [→](https://github.com/MrChristianL/Teleport-Job-Worker-Service/pull/2#discussion_r2905879180)

> nit: prefer `wg.Go()` where possible
> ```suggestion
> 	wg.Go(func() {
> 		var buf bytes.Buffer
> 		readyWg.Done()
> 		broker.streamFromDisk(context.Background(), &buf)
> 		readerOutput <- buf.Bytes()
> 	})
> ```

**@rosstimothy** on `MrChristianL/Teleport-Job-Worker-Service` [→](https://github.com/MrChristianL/Teleport-Job-Worker-Service/pull/2#discussion_r2919465161)

> Given that the contents of the file are rather small, this goroutine and channel approach make this test more complicated than it needs to be. The test itself doesn't need to exercise concurrent readers to validate that that multiple clients can stream the same historic data from a closed broker.

**@rosstimothy** on `MrChristianL/Teleport-Job-Worker-Service` [→](https://github.com/MrChristianL/Teleport-Job-Worker-Service/pull/2#discussion_r2919518045)

> Could this test be simplified by passing in an already cancelled context?
> ```suggestion
> 	ctx, cancel := context.WithCancel(t.Context())
> 	cancel()
>
> 	
> 	err := job.StreamFromDisk(ctx, &bytes.Buffer{})
> ```

**@nklaassen** on `MrChristianL/Teleport-Job-Worker-Service` [→](https://github.com/MrChristianL/Teleport-Job-Worker-Service/pull/2#discussion_r2927365455)

> you will, however, get an extra spurious cond.Broadcast at the end of every Read RPC when the context is canceled (waking other readers unnecessarily).
>
> You could consider something like this, taken from the example on context.AfterFunc, which will avoid broadcasting after streamFromDisk returns, and may avoid a goroutine
> ```go
> stopf := context.AfterFunc(ctx, func() {
> 	b.mu.Lock()
> 	b.cond.Broadcast()
> 	b.mu.Unlock()
> })
> defer stopf()
> ```

**@hugoShaka** on `neildo/tjob` [→](https://github.com/neildo/tjob/pull/1#discussion_r1758934248)

> How will the client stream logs? This RPC only sends the response once. Should it be a `stream`? Should the client poll?

**@hugoShaka** on `neildo/tjob` [→](https://github.com/neildo/tjob/pull/1#discussion_r1759025001)

> The challenge asks for multiple clients to concurrently stream logs from the same job. It's not clear from the function signatures how new logs will be streamed, if clients will get notified, do polling, ...
> 
> For example, you should be able to run `tjob logs <jobid>` on an unfinished job `bash -c "while true; do date && sleep 10; done"` and see a new log line every 10 seconds.
> 
> As log streaming is more complex than one-shot log lookup, it would be useful to have a dedicated section in the RFD, like you did for cgroups and namespaces. This section could describe how the server collects the job logs and how it streams the new logs to the watching clients.

**@rosstimothy** on `neildo/tjob` [→](https://github.com/neildo/tjob/pull/1#discussion_r1759128947)

> As discussed out of band, I meant that streaming should be the default and only option, and that you could simplify by making that so and drop the `--tail` and `--follow` flags.

**@rosstimothy** on `neildo/tjob` [→](https://github.com/neildo/tjob/pull/1#discussion_r1759137026)

> The examples are nice, but could you describe in words what additional things need to be handled beyond setting these flags for true resource isolation?

**@rosstimothy** on `neildo/tjob` [→](https://github.com/neildo/tjob/pull/1#discussion_r1763583075)

> Expectations for a L5 are that output streaming can be achieved without relying on a periodic sleep and polling a file for more output. Could you please use inotify directly, or use some other means here instead.

**@rosstimothy** on `neildo/tjob` [→](https://github.com/neildo/tjob/pull/3#discussion_r1771894169)

> What ownership guarantees are required for payloads in a gRPC message? Is there any potential for the message to be buffered and the next read to overwrite this buffer?

**@smallinsky** on `P-A-R-U-S/Go-Job-Worker-Service` [→](https://github.com/P-A-R-U-S/Go-Job-Worker-Service/pull/1#discussion_r1730898720)

> How a client will be able to cancel ongoing stream ?

**@smallinsky** on `P-A-R-U-S/Go-Job-Worker-Service` [→](https://github.com/P-A-R-U-S/Go-Job-Worker-Service/pull/1#discussion_r1730905222)

> Could you elaborate bit more about: 
> How a client stream part will be notified about new data produced by a process and where data will be stored (memory/file) ?

**@greedy52** on `P-A-R-U-S/Go-Job-Worker-Service` [→](https://github.com/P-A-R-U-S/Go-Job-Worker-Service/pull/1#discussion_r1731340032)

> I would also have multiple clients stream the same job out simultaneously. And another stream after job is completed.

**@smallinsky** on `P-A-R-U-S/Go-Job-Worker-Service` [→](https://github.com/P-A-R-U-S/Go-Job-Worker-Service/pull/1#discussion_r1732271573)

> This will work on grpc client->server side, but what about server-side ->library usage?
> 
> Right now, the Job interface is defined as:
> ```
> func New(config JobConfig) *Job 
> func (*Job) Stream() OutputReader
> ```
> 
> Currently, there is no way for the library user to cancel the `Stream()` call.
> 
> What would you think about returning an `io.ReadCloser` object from `func (*Job) Stream() io.ReadCloser` and allow job `Stream` consumer to cancel the streaming in any time ?

**@greedy52** on `P-A-R-U-S/Go-Job-Worker-Service` [→](https://github.com/P-A-R-U-S/Go-Job-Worker-Service/pull/2#discussion_r1743727970)

> This doesn't look right to me. Closing the reader shouldn't close the buffer. 
> 
> Also, ideally closing the reader should stop any running `Read`. How do you plan to handle the case where a client stops streaming (while the process is still running)?

**@smallinsky** on `P-A-R-U-S/Go-Job-Worker-Service` [→](https://github.com/P-A-R-U-S/Go-Job-Worker-Service/pull/2#discussion_r1743760502)

> ```suggestion
> func (job *Job) Stream() io.ReadCloser {
> ```

**@smallinsky** on `P-A-R-U-S/Go-Job-Worker-Service` [→](https://github.com/P-A-R-U-S/Go-Job-Worker-Service/pull/2#discussion_r1745002764)

> In this design the Read and Close calls are not thread safe. 
> 
> So let say that consumer want to create a read loop that will block by r.Read call this a new output is available but  want to cancel the reeder loop 
> ```go 
> ctx, cancel := ctx.WithCancel(ctx)
> defer cancel()
> 
> r : NewOutputReadCloser()
> 
> go func() {
>    switch {
>     case <- ctx.Done:
>         r.Close()
>    } 
> }
> for  { 
> 
>   .., ... = r.Read(buff)
>   ... 
> }
> 
> ```
> 
> So the question is how a consumer of outputReadCloser can cancel ongoing stream in the Read/Close  if `outputReadCloser `  Read Close functions are not thread safe ?

**@smallinsky** on `P-A-R-U-S/Go-Job-Worker-Service` [→](https://github.com/P-A-R-U-S/Go-Job-Worker-Service/pull/3#discussion_r1749755721)

> Stream is blocking operation 
> 
> By  calling 
> ```
> 	s.mutex.Lock()
> 	defer s.mutex.Unlock()
> ```
> 
>  jobWorker will be blocked for all  users
> 
> 
> also the same behaviors applies to Stop gRPC handler

**@atburke** on `rajansandeep/teleport-job-worker-svc` [→](https://github.com/rajansandeep/job-worker-svc/pull/1#discussion_r3081000253)

> Do the cert flags need to be specified for every invocation of `jobworker-cli`?

**@rhammonds-teleport** on `rajansandeep/teleport-job-worker-svc` [→](https://github.com/rajansandeep/job-worker-svc/pull/1#discussion_r3087148403)

> ```
> Start(cmd string, args []string, output *OutputBuffer) (Process, error)
> ```
>
> It's likely that the implementation will only need to call `Write` and `Close` on `OutputBuffer`, right? Maybe we could use a well known abstraction from the standard library 😃

**@eriktate** on `rajansandeep/teleport-job-worker-svc` [→](https://github.com/rajansandeep/job-worker-svc/pull/1#discussion_r3088864013)

> Is there a way we could facilitate cancellation while still implementing the `io.Reader` or `io.ReadCloser` interface? I would also point out that there's currently no limit to the size of the returned `[]byte`. Considering gRPC has message size limits, it's worth considering how best to buffer output. Implementing `io.Reader` could help simplify that

**@eriktate** on `rajansandeep/teleport-job-worker-svc` [→](https://github.com/rajansandeep/job-worker-svc/pull/1#discussion_r3088975824)

> nit: I'd consider dropping the prefix `job_id:` to make it easier when piping the job ID. e.g.
> ```sh
> jobworker-cli start -- ls -la /tmp | jobworker-cli stream
> ```

**@rhammonds-teleport** on `rajansandeep/teleport-job-worker-svc` [→](https://github.com/rajansandeep/job-worker-svc/pull/2#discussion_r3100945468)

> I'm not a fan of `cancelOnCloseReader`, and `StreamOutput` doesn't really need a context parameter. The function itself does not block, and the caller can unblock itself from a `Read` by calling `Close` on the returned `ReadCloser`. The caller has two ways of stopping the reader which is confusing.
>
> Also, let's say we keep this code and a reader does end up getting closed via context cancellation. As it's written, I would expect a blocked call to `Read()` to return a `context.Canceled` or `context.DeadlineExceeded` error. The current implementation would just return an `io.EOF` indicating that all data has been read. The caller can't distinguish between a full/complete read or an interrupted/canceled read.

**@rhammonds-teleport** on `rajansandeep/teleport-job-worker-svc` [→](https://github.com/rajansandeep/job-worker-svc/pull/3#discussion_r3111236582)

> I think you could simplify this with:
> ```
> // Use AfterFunc to automatically close the reader on
> // context cancellation. This will break us out of the
> // read loop below.
> cancelStop := context.AfterFunc(stream.Context(), func() {
> 	_ = reader.Close()
> })
> defer cancelStop()
> ```

**@rhammonds-teleport** on `rajansandeep/teleport-job-worker-svc` [→](https://github.com/rajansandeep/job-worker-svc/pull/3#discussion_r3111369288)

> It seems like this entire chunk could be simplified to:
>
> ```
> if err = mapWorkerError(err); err != nil {
>    return err
> }
>
> // err is nil
> if ctxErr := stream.Context().Err(); ctxErr != nil {
> 	err = ctxErr
> }
> return err
> ```
>
> if `mapWorkerError` were extended to handle io.ErrClosedPipe and nil errors  with:
> ```
> case errors.Is(err, io.ErrClosedPipe):
>   return status.Error(codes.Canceled, "stream interrupted")
> case err == nil:
>   return nil
> ```
>
> I would say that although `io.ErrClosedPipe` isn't a custom error defined by your worker package, it _is_ part of `worker`'s API contract.
>
> I appreciate that you're programming defensively here and guarding against potential (unexpected) changes to that contract though!

**@atburke** on `rajansandeep/teleport-job-worker-svc` [→](https://github.com/rajansandeep/job-worker-svc/pull/4#discussion_r3120037609)

> Nit: naming the return args might make it easier to tell at a glance which return value is which. Up to you if you want to change it.
> ```suggestion
> func resolveTLSPaths(certFlag, keyFlag, caFlag string) (certPath, keyPath, caPath string, err error) {
> ```

**@zmb3** on `razzam21/job-worker-service` [→](https://github.com/razzam21/job-worker-service/pull/1#discussion_r2255188366)

> We don't really need to include the job ID in the response since it was specified in the request.
>
> (Same goes for output streaming, where it's even more beneficial to avoid sending unnecessary data on the wire)

**@zmb3** on `razzam21/job-worker-service` [→](https://github.com/razzam21/job-worker-service/pull/1#discussion_r2255190354)

> The challenge only requires that we offer the ability to stream output from the very beginning of process execution, so feel free to simplify by removing the `from_offset` field.
>
> Similarly, we can make `follow` the default behavior to further reduce scope.

**@zmb3** on `razzam21/job-worker-service` [→](https://github.com/razzam21/job-worker-service/pull/1#discussion_r2255196023)

> I could use a little more info on how the streaming will work. What will you use for efficient output discovery?

**@creack** on `razzam21/job-worker-service` [→](https://github.com/razzam21/job-worker-service/pull/1#discussion_r2257544958)

> What happens when a client stops consuming the stream?

**@zmb3** on `razzam21/job-worker-service` [→](https://github.com/razzam21/job-worker-service/pull/1#discussion_r2257557628)

> I see you're looking to use channels for streaming.
>
> The thing about channels is that a slow reader can prevent the writer from writing. How do you plan to handle this?

**@creack** on `razzam21/job-worker-service` [→](https://github.com/razzam21/job-worker-service/pull/1#discussion_r2257573337)

> Not sure to understand the goal of the notification channel. Could you expand a bit more?

**@zmb3** on `razzam21/job-worker-service` [→](https://github.com/razzam21/job-worker-service/pull/1#discussion_r2258489593)

> I don't expect the buffered channel to change much here as it doesn't solve the slow-reader problem, it merely delays it.
> 
> The 5-second timeout to disconnect slow readers works, if we're okay requiring that clients read quickly and terminating any clients that don't follow our rules.
> 
> Here's the thing - we're buffering all data in-memory anyway. The data is guaranteed to exist as long as the server stays running, so it really shouldn't matter how fast the reader wants to consume the data, right?

**@zmb3** on `razzam21/job-worker-service` [→](https://github.com/razzam21/job-worker-service/pull/1#discussion_r2258493346)

> So what would the Go code look like for an append operation on this buffer? That might help me conceptualize this a bit.

**@awly** on `renatoaguimaraes/golang-job-scheduler` [→](https://github.com/renatoaguimaraes/golang-job-scheduler/pull/2#discussion_r623379185)

> ```suggestion
> 				select {
> 				// send the event throught the channel
> 				case eventchan <- wrapLinuxEvent(raw.Mask):
> 				// stop the watcher
> 				case <-ctx.Done():
> 					return
> 				}
> 				// calculate the offset
> 				offset += unix.SizeofInotifyEvent + raw.Len
> ```

**@r0mant** on `renatoaguimaraes/golang-job-scheduler` [→](https://github.com/renatoaguimaraes/golang-job-scheduler/pull/5#discussion_r625200306)

> Why not use `stream.Context()` here?

**@awly** on `renatoaguimaraes/golang-job-scheduler` [→](https://github.com/renatoaguimaraes/golang-job-scheduler/pull/2#discussion_r625345531)

> nit: you probably don't need a buffered reader anymore

**@awly** on `renatoaguimaraes/golang-job-scheduler` [→](https://github.com/renatoaguimaraes/golang-job-scheduler/pull/6#discussion_r625935262)

> let's bind these to flags for now using golang.org/pkg/flag

**@awly** on `renatoaguimaraes/golang-job-scheduler` [→](https://github.com/renatoaguimaraes/golang-job-scheduler/pull/6#discussion_r625943821)

> when the stream ends server-side (when the job finishes), you should get `io.EOF` or a similar error
> don't output those errors

**@strideynet** on `rexposadas/teleport` [→](https://github.com/rexposadas/teleport/pull/1#discussion_r1790419045)

> A little more detail on how this will work would be good.

**@strideynet** on `rexposadas/teleport` [→](https://github.com/rexposadas/teleport/pull/1#discussion_r1793611384)

> What concurrency primitives will your design use to implement the requirements of the log streaming?

**@rosstimothy** on `RichyHBM/teleport-challenge` [→](https://github.com/RichyHBM/teleport-challenge/pull/1#discussion_r2048945602)

> Could you provide a high level explanation for how consumers of the output will be alerted that more output is available if they've already consumed all the historical output that was in the buffer?

**@tigrato** on `RichyHBM/teleport-challenge` [→](https://github.com/RichyHBM/teleport-challenge/pull/1#discussion_r2049668880)

> What happens if the client is slow to consume this data and the channel becomes full?

**@rosstimothy** on `RichyHBM/teleport-challenge` [→](https://github.com/RichyHBM/teleport-challenge/pull/1#discussion_r2050614374)

> Could you provide a bit more detail on what will be performed during the synchronous operation to actually terminate the jobs?

**@rosstimothy** on `RichyHBM/teleport-challenge` [→](https://github.com/RichyHBM/teleport-challenge/pull/2#discussion_r2056048698)

> Suggestion: here and above reduce the error scope where possible
> ```suggestion
> 		if err = grpcClient.Close();  err != nil {
> 			log.Printf("Unable to close gRPC channel %v", err)
> 		}
> ```

**@eriktate** on `RichyHBM/teleport-challenge` [→](https://github.com/RichyHBM/teleport-challenge/pull/2#discussion_r2058463002)

> I see `pflag` is an indirect dependency through `cliff`, but I don't think it's being used here?

**@rosstimothy** on `RichyHBM/teleport-challenge` [→](https://github.com/RichyHBM/teleport-challenge/pull/5#discussion_r2064101856)

> What ownership guarantees about the buffer need to be taken into consideration here?

**@rosstimothy** on `RichyHBM/teleport-challenge` [→](https://github.com/RichyHBM/teleport-challenge/pull/5#discussion_r2064107421)

> I think you could reduce the complexity a bit if this was inverted. Instead of registering an io.Writer the API provided a means for callers to get a unique io.Reader that tracked their progress through the buffer.

**@rosstimothy** on `RichyHBM/teleport-challenge` [→](https://github.com/RichyHBM/teleport-challenge/pull/5#discussion_r2064472100)

> If this doesn't return a copy of the bytes then it's a bit unclear what the ownership semantics of the returned byte slice are. If the mutex is meant to protect the buffer, then doesn't this currently bypass that?

**@eriktate** on `RichyHBM/teleport-challenge` [→](https://github.com/RichyHBM/teleport-challenge/pull/5#discussion_r2064946941)

> Is this true? This seems like a potentially dangerous UX if we force users to issue `stop` commands for cases other than stopping a running job. What additional details are missing?

**@eriktate** on `RichyHBM/teleport-challenge` [→](https://github.com/RichyHBM/teleport-challenge/pull/5#discussion_r2064984664)

> I think you might have a data race here. If you have a job continuously logging output, what happens to the data captured between `outputStream.GetBuffer()` and `job.outputStream.Connect()`?

**@timothyb89** on `rsteinkeXJ/teleport_interview` [→](https://github.com/rsteinkeXJ/teleport_interview/pull/1#discussion_r3121467668)

> These aren't required, you can drop them to save some time if you'd like. (And arguably, `AttachOutStream()` should function as a `WaitJob()`)

**@timothyb89** on `rsteinkeXJ/teleport_interview` [→](https://github.com/rsteinkeXJ/teleport_interview/pull/1#discussion_r3121530915)

> Hmm, I'm not sure this meets the challenge requirements:
> > * Library should be able to stream the output of a job. 
> >   * [...]
> >   * Output should be from start of process execution.
> >   * Multiple concurrent clients should be supported.
> >   * [...]
>
> It seems like attaching streams isn't quite the metaphor we're looking for?

**@timothyb89** on `rsteinkeXJ/teleport_interview` [→](https://github.com/rsteinkeXJ/teleport_interview/pull/1#discussion_r3121536242)

> Hmm, it seems like you could just always assume the streams should be attached, right? Clients don't necessarily need to configure this, and it'd probably simplify the implementation a bit.

**@tigrato** on `rsteinkeXJ/teleport_interview` [→](https://github.com/rsteinkeXJ/teleport_interview/pull/1#discussion_r3123765896)

> If it helps, feel free to merge stderr and stdout into the same stream instead of having multiple streams, otherwise you can't correctly represent intertwined stdout and stderr

**@tigrato** on `rsteinkeXJ/teleport_interview` [→](https://github.com/rsteinkeXJ/teleport_interview/pull/1#discussion_r3124173353)

> Could you include more detail on how the stream will operate?
>
> Please also consider how it will satisfy the challenge requirements

**@rosstimothy** on `rsteinkeXJ/teleport_interview` [→](https://github.com/rsteinkeXJ/teleport_interview/pull/1#discussion_r3138112216)

> Do you think we could improve the CLI UX a bit? The single character flags switching the behavior was a bit confusing to me at first. Then I realized k was for kill q was for query and w was for watch. Perhaps subcommands would make this more obvious?

**@rosstimothy** on `rsteinkeXJ/teleport_interview` [→](https://github.com/rsteinkeXJ/teleport_interview/pull/1#discussion_r3138146461)

> Can you provide more details on the output streaming model. How does output get from the command to the array of events from process state? How does this channel mechanism work?

**@tigrato** on `rsteinkeXJ/teleport_interview` [→](https://github.com/rsteinkeXJ/teleport_interview/pull/1#discussion_r3149024777)

> what if a channel is full because a client is slow? what will occur?
> will we block or will we skip the event?
>
> is there another sync primitive that doesn't require using channels, or if you prefer to use channels is there a way you could ensure that a slow consumer won't:
>
> - lose data
> - impact other consumers

**@tigrato** on `rsteinkeXJ/teleport_interview` [→](https://github.com/rsteinkeXJ/teleport_interview/pull/1#discussion_r3149035279)

> can you include mode details how the streaming mechanism will occur?
>
> We are interested in the sync primitives, where it can miss events and how the transition between past data and real time data occurs

**@zmb3** on `s-gruneberg/jobWorker` [→](https://github.com/s-gruneberg/jobWorker/pull/1#discussion_r2305449999)

> This section is good and covers a lot of the important details.
>
> I'm not requesting any specific changes, but would like to offer a general tip. This section reads as if it is describing code that already exists, which is a bit hard for me to process as a reader. It feels as if we (the author and the reader) are discussing code that you can see and I can't.
>
> If you reframe things a bit to focus less on code specifics (for example, whether or not you use a switch statement is not something we care about at this stage of the project) and more on how things will work and the key decisions you've made then you'll have a better quality document.

**@russjones** on `s-gruneberg/jobWorker` [→](https://github.com/s-gruneberg/jobWorker/pull/1#discussion_r2305715715)

> Can you provide some more details. For example, what verbs do you plan to use, expected response codes, request body, and response body formats.

**@russjones** on `s-gruneberg/jobWorker` [→](https://github.com/s-gruneberg/jobWorker/pull/1#discussion_r2305726016)

> Can you add details about what version of TLS you plan to use, what ciphersuites, etc.?

**@rosstimothy** on `sabernabil12/teleport-job-worker` [→](https://github.com/sabernabil12/teleport-job-worker/pull/1#discussion_r2164741540)

> How can clients read from the same channel without "stealing" data from each other?

**@rosstimothy** on `sabernabil12/teleport-job-worker` [→](https://github.com/sabernabil12/teleport-job-worker/pull/1#discussion_r2164743845)

> Is this accurate? Isn't there two goroutines copying data from the process pipes to a channel?

**@nklaassen** on `sabernabil12/teleport-job-worker` [→](https://github.com/sabernabil12/teleport-job-worker/pull/1#discussion_r2165151179)

> I don't think this is possible given the requirements and the gRPC streaming

**@nklaassen** on `sabernabil12/teleport-job-worker` [→](https://github.com/sabernabil12/teleport-job-worker/pull/1#discussion_r2165156151)

> output buffered in channels probably doesn't work for multiple clients streaming from the start

**@rosstimothy** on `sabernabil12/teleport-job-worker` [→](https://github.com/sabernabil12/teleport-job-worker/pull/1#discussion_r2166639666)

> Hrm, doesn't the proposed `channelWriter` implement `io.Writer` though? Do the pipes and goroutines aid at all in real time streaming? Their only purpose is to copy data from the process to your channelWriter. Isn't it the channelWriters responsibility to support real time streaming?

**@smallinsky** on `samschurter/teleport-challenge` [→](https://github.com/samschurter/teleport-challenge/pull/2#discussion_r698317503)

> I don't think that this is totally safe. For instance if the `safeBuffer` user would like to modify the `Bytes()`  result it will lead to race condition.

**@strideynet** on `shawon-crosen/teleport-challenge` [→](https://github.com/shawon-crosen/teleport-challenge/pull/1#discussion_r2124425849)

> It seems a little odd to me to leverage `tail` rather than just directly reading the file? Is there any benefit to your approach?

**@strideynet** on `shawon-crosen/teleport-challenge` [→](https://github.com/shawon-crosen/teleport-challenge/pull/1#discussion_r2125939748)

> It's not very clear to me what CLI these flags are for - I would guess perhaps openssl but I don't think I recognise these flags.

**@rosstimothy** on `shawon-crosen/teleport-challenge` [→](https://github.com/shawon-crosen/teleport-challenge/pull/1#discussion_r2126983374)

> When consuming the file to stream the data back to the user what happens if reading the file returns an io.EOF and the process is still running?

**@zmb3** on `shawon-crosen/teleport-challenge` [→](https://github.com/shawon-crosen/teleport-challenge/pull/1#discussion_r2129165255)

> Can you include more details on this? Is this what you're thinking?
> 
> - Read from file, first read succeeds with valid data
> - Read from file again, this read returns EOF since you reached the end of the file
> - Check whether process is still running, if not, return.
> - If process is still running, read again.

**@r0mant** on `supby/job-worker` [→](https://github.com/supby/job-worker/pull/1#discussion_r716789825)

> Could you add some information about how log streaming will work please?

**@smallinsky** on `supby/job-worker` [→](https://github.com/supby/job-worker/pull/2#discussion_r720147600)

> why you choose to use buffered channel of size 20 ?

**@r0mant** on `supby/job-worker` [→](https://github.com/supby/job-worker/pull/2#discussion_r721442082)

> Are multiple clients able to stream the logs? Also, how do clients get logs from the beginning?

**@rosstimothy** on `tjper/teleport` [→](https://github.com/tjper/teleport/pull/1#discussion_r821058198)

> You can omit this and always tail to simplify.

**@rosstimothy** on `tjper/teleport` [→](https://github.com/tjper/teleport/pull/1#discussion_r821061402)

> You can always `--tail` to simplify. Also can you describe how you will implement streaming?

**@rosstimothy** on `tjper/teleport` [→](https://github.com/tjper/teleport/pull/1#discussion_r821063672)

> How do you plan on stopping a job? I don't see any details about it below.

**@codingllama** on `tjper/teleport` [→](https://github.com/tjper/teleport/pull/1#discussion_r821707444)

> Let's assume the client never cancels, as they want to stream until the job is over. How do you detect/handle that?

**@codingllama** on `tjper/teleport` [→](https://github.com/tjper/teleport/pull/1#discussion_r821714823)

> Please mention the different streams (stdout and stderr).

**@codingllama** on `tjper/teleport` [→](https://github.com/tjper/teleport/pull/3#discussion_r829314963)

> Looking back at this, I'd roll job.Job's API entirely into JobService. Callers have to go through multiple hoops to get actions done, where we could have a simpler API
> 
> For example:
> 
> * Starting a job is job.New followed by jobService.StartJob (skipping the JobService creation). There's nothing useful to do with that job pointer before StartJob.
> * Streaming is jobService.FetchJob followed by job.StreamOutput
> * Status is jobService.FetchJob followed by job.Status
> * Stopping is just jobService.StopJob, which I think is better, but is an outlier compared to others.

**@codingllama** on `tjper/teleport` [→](https://github.com/tjper/teleport/pull/3#discussion_r829322437)

> A blocking StreamOutput function forces all callers to deal with concurrency for it to be useful. In this case I think making the function non-blocking would make for an API that is easier to use.
> 
> You may either want to return the channel for streaming or take ownership of it, so it can be closed when the stream ends. We'd get a usage idiom like this:
> 
> 
> ```go
> ch := make(chan []byte)
> if err := job.StreamOutput(ctx, ch, chunkSize); err != nil {
>   // Abort.
> }
> for _, chunk := range ch {
>   // Handle chunk
> }
> ```
> 
> If you take the change don't forget to update the godocs.

**@codingllama** on `tjper/teleport` [→](https://github.com/tjper/teleport/pull/3#discussion_r829329468)

> Flags in tests are a pain, in my experience. They force you to single out tests in CI so you can set the flags, which causes all sorts of complications. Better to use env variables.
> 
> In this case we should be using the equivalent of net.Listen("localhost:") to bring up a new server, instead of picking an explicit port number.

**@zmb3** on `tjper/teleport` [→](https://github.com/tjper/teleport/pull/3#discussion_r830294292)

> No change necessary, just a comment.
> 
> I'm of the mind that flag parsing and help messages are the responsibility of a program, not of a library, and should live closer to main.

**@zmb3** on `tjper/teleport` [→](https://github.com/tjper/teleport/pull/3#discussion_r830299455)

> Hmm, the bytes you're sending over the channel here are the same memory that the next iteration of the loop is going to try and populate.
> 
> I think you may have an issue here if the consumer of the channel is slow - especially because the channel is buffered. This means you can have multiple items buffered up in the channel, but they all actually refer to the same memory which is being continuously overwritten.

**@greedy52** on `xSudoNymx/job-manager` [→](https://github.com/xSudoNymx/job-manager/pull/1#discussion_r2704972091)

> There is chance that new output comes in where replay is still in progress. How do you plan to coordinate when to switch to the channel?

**@greedy52** on `xSudoNymx/job-manager` [→](https://github.com/xSudoNymx/job-manager/pull/1#discussion_r2704983272)

> > Given that fanout is best-effort, a slow subscriber may drop live frames once its buffer fills. If stronger guarantees are needed later, we can introduce per-subscriber replay offsets or explicit backpressure policies.
>
> Is it possible to design something that avoids dropping live frames for this challenge? not for later

**@greedy52** on `xSudoNymx/job-manager` [→](https://github.com/xSudoNymx/job-manager/pull/1#discussion_r2705000443)

> TLS 1.3 is great. could you provide some details on the certificates like how they are generated, verified, etc.

**@tigrato** on `xSudoNymx/job-manager` [→](https://github.com/xSudoNymx/job-manager/pull/1#discussion_r2709738410)

> >Note: This change amplifies disk reads. for the challenge it may be a non issue but in prod could produce risk if the number of subscribers are unbounded.
> 
> You can always drop disk usage and leverage memory if this is a concern to you. For us it's not a concern but memory is easier to handle.
> 
> >Sure! the next step would be to have each subscriber track its own read offsets and tail the file. I would also have to add a notification mechanism that wakes up the subscribers when append happens. So to summarize my thoughts each subscriber would have a goroutine that replays from offset 0 then continues reading the log as it grows.
> 
> In this case are you dropping the fanout mechanism?

**@tigrato** on `xSudoNymx/job-manager` [→](https://github.com/xSudoNymx/job-manager/pull/1#discussion_r2709753017)

> >@greedy52 Output is always persisted to the log file before any live fan-out occurs. When a client subscribes, it replays the log sequentially until EOF. Because replay reads the same log that output is appended to, any output produced while replay is in progress is guaranteed to appear either in the replay stream (if the reader has not yet passed that offset) or via live fan-out after replay completes. Reaching EOF marks the transition point from replay to live streaming.
> 
> If io.EOF is the trigger to switch from file read to fanout mechanism, is it possible to drop payloads during the switch?
> 
> Is there anything that prevents the system from writing to the file while the client is switching to the fanout mechanism?
> It seems possible that such writes could occur: once the client receives an io.EOF from the file read, it is no longer interested in the file and switches to fanout, which could cause the client to miss some events.

**@greedy52** on `xSudoNymx/job-manager` [→](https://github.com/xSudoNymx/job-manager/pull/1#discussion_r2714033109)

> what mechanism is used to signal the readers for log closure? note that a client may also start streaming after the job is exited.

**@tigrato** on `xSudoNymx/job-manager` [→](https://github.com/xSudoNymx/job-manager/pull/1#discussion_r2714479470)

> A subscriber might not miss the notification to wake up if it is parked, but if it's not yet parked - i.e. sending payload to the client - it's totally possible for a subscriber to miss a broadcast from the sync.Cond.
> In those cases do you re-check the committedLen before calling wait?

**@tigrato** on `xSudoNymx/job-manager` [→](https://github.com/xSudoNymx/job-manager/pull/1#discussion_r2714481758)

> How do we identify that a client drops mid stream?

**@tigrato** on `xSudoNymx/job-manager` [→](https://github.com/xSudoNymx/job-manager/pull/1#discussion_r2716485639)

> Just to be clear, if the notify channel has anything we skip the insertion right? since the subscriber will be already notified once he reads from the channel

**@tigrato** on `xSudoNymx/job-manager` [→](https://github.com/xSudoNymx/job-manager/pull/3#discussion_r2728215541)

> can we satisfy io.ReadCloser interface?
>
> that would allow us to io.Copy(os.Stdout, subscriber)

**@sclevine** on `zeeshanhaque21/JobWorkerService` [→](https://github.com/zeeshanhaque21/JobWorkerService/pull/1#discussion_r1686913061)

> Why are logs stored in chunks? What if a client wants to read less than one chunk? Are channels ideal for streaming data?

**@gabrielcorado** on `zeeshanhaque21/JobWorkerService` [→](https://github.com/zeeshanhaque21/JobWorkerService/pull/1#discussion_r1692429066)

> Will those channels receive all the logs or just the latest entries? How will consumers fetch all log entries since the beginning of the job execution?

**@sclevine** on `zeeshanhaque21/JobWorkerService` [→](https://github.com/zeeshanhaque21/JobWorkerService/pull/1#discussion_r1695471383)

> As a consumer of the library, I think I'd expect to be able to stream bytes with a standard io.Reader or io.ReadCloser, which allow me read an arbitrary number of bytes on each read. These may be individually buffered, or they may pull from shared byte storage.

**@espadolini** on `Zephan92/teleport-handson` [→](https://github.com/Zephan92/teleport-handson/pull/1#discussion_r2718759422)

> How does a reader get unblocked if its client stops the stream while blocked on the condvar?

**@espadolini** on `Zephan92/teleport-handson` [→](https://github.com/Zephan92/teleport-handson/pull/1#discussion_r2718770749)

> For the sake of streamlining the implementation, feel free to hardcode some resource limits that are applied to every cgroup without having to wire numbers through from the client CLI to the worker library.

**@espadolini** on `Zephan92/teleport-handson` [→](https://github.com/Zephan92/teleport-handson/pull/1#discussion_r2718775625)

> The pseudocode for `OutputBuffer` doesn't seem to be able to support interleaved streams for stdout and stderr and timestamps.

**@nklaassen** on `Zephan92/teleport-handson` [→](https://github.com/Zephan92/teleport-handson/pull/1#discussion_r2719088573)

> you mentioned the output is designed to be machine-readable where possible. the most likely case i'd think you'd want the output to be machine-readable is when streaming process output, and you probably wouldn't want superfluous timestamp injected in the middle of the output in this case. It's also not clear what these timestamps are... the time that line of output was read by the library? or the time the cli received it?

**@nklaassen** on `Zephan92/teleport-handson` [→](https://github.com/Zephan92/teleport-handson/pull/1#discussion_r2719116629)

> what is the `source io.Reader` going to be set to?
>
> Actually I'm being a bit coy with that question, i don't think there should be an io.Reader source, it can be avoided
>
> also whatever reads from this buffer is probably going to be more interesting than whatever writes to it, but I don't see much about that or any pseudo-code for the read side

**@nklaassen** on `Zephan92/teleport-handson` [→](https://github.com/Zephan92/teleport-handson/pull/1#discussion_r2728991505)

> what will happen if the client disconnects while blocked on the condition variable here? I think it could be sort of a goroutine leak if many clients stream output of a job that runs indefinitely without writing any more output

**@nklaassen** on `Zephan92/teleport-handson` [→](https://github.com/Zephan92/teleport-handson/pull/1#discussion_r2729754263)

> The problem with example code is I feel compelled to review it. Broadcasting without holding the lock almost always leads to a race, and this case is no exception. But, okay we can say this is pseudo code and hold off on detailed review until the implementation PR.

**@espadolini** on `Zephan92/teleport-handson` [→](https://github.com/Zephan92/teleport-handson/pull/4#discussion_r2745666663)

> You don't have to indirect and allocate the condvar separately, including it in`OutputBuffer` and setting `L` any time before the first use works just as well.

**@espadolini** on `Zephan92/teleport-handson` [→](https://github.com/Zephan92/teleport-handson/pull/4#discussion_r2745669029)

> What's the point of `Stdout()` if the `*OutputBuffer` does the same thing on its own?

**@espadolini** on `Zephan92/teleport-handson` [→](https://github.com/Zephan92/teleport-handson/pull/4#discussion_r2745743725)

> I sort of question the design, you need a goroutine for each reader and you end up with an API that still requires manual error checks on the reader side (to distinguish between the output being done and the context getting canceled).
>
> Why can't the reader of the stream call a fallible method on some object to get the next `LogEntry`, for example?

**@nklaassen** on `Zephan92/teleport-handson` [→](https://github.com/Zephan92/teleport-handson/pull/4#discussion_r2747982565)

> agreed, as is you will have two loops, the streamLoop in a goroutine writing to the channel and another loop at the caller reading from the channel. Edoardo's suggesting is good

**@nklaassen** on `Zephan92/teleport-handson` [→](https://github.com/Zephan92/teleport-handson/pull/4#discussion_r2747987121)

> don't you need a broadcast (while holding the lock) on context cancellation to avoid getting blocked in a cond.Wait(), as the comment says?

**@nklaassen** on `Zephan92/teleport-handson` [→](https://github.com/Zephan92/teleport-handson/pull/4#discussion_r2748516402)

> ```suggestion
> 	stream := buf.Stream(t.Context())
> ```

**@nklaassen** on `Zephan92/teleport-handson` [→](https://github.com/Zephan92/teleport-handson/pull/4#discussion_r2748529784)

> I don't really see the point of the goroutine here, if you want to test writes before or after the buf.Stream call you can do that, this just makes it non-deterministic. If you really want to test concurrency you'd probably need to do some more writes and maybe some pauses

---

## auth-identity

_Identity verification from certificates; handling Subject vs Serial vs VerifiedChains vs PeerCertificates; avoiding info leakage in errors._

**365 quotes** from `27` distinct reviewers across `45` candidate submissions.

**@codingllama** on `adalton/teleport-exercise` [→](https://github.com/adalton/teleport-exercise/pull/2#discussion_r771417047)

> Embed the certs? Keeps you free of hard-coded path dependencies.

**@codingllama** on `adalton/teleport-exercise` [→](https://github.com/adalton/teleport-exercise/pull/2#discussion_r771457573)

> Mentioned before, but if you embed the certs you can get away without this.

**@codingllama** on `adalton/teleport-exercise` [→](https://github.com/adalton/teleport-exercise/pull/2#discussion_r771458581)

> The server should be the one attaching the user credentials to the context, based on their cert. Plus, whatever is in this context doesn't go all the way, so I expect this to do nothing.

**@codingllama** on `adalton/teleport-exercise` [→](https://github.com/adalton/teleport-exercise/pull/2#discussion_r771464492)

> ```suggestion
> 		return nil, errors.New("cannot append ca cert to ca pool")
> ```
> 
> (Using a constant string.)

**@codingllama** on `adalton/teleport-exercise` [→](https://github.com/adalton/teleport-exercise/pull/2#discussion_r772439109)

> I have a few suggestions regarding this setup:
> 
> * Place all certificate constants in a single file? Seems simpler to manage.
> * Ditch the factory concept and just use the constants in pairs?

**@smallinsky** on `bill-rich/jobworker` [→](https://github.com/bill-rich/jobworker/pull/1#discussion_r1606446875)

> Could you elaborate how simple authorization scheme will look like ? 
> 
> How are you planing distinguish clients  ?

**@rosstimothy** on `bill-rich/jobworker` [→](https://github.com/bill-rich/jobworker/pull/1#discussion_r1608423258)

> Does Go allow you to chose cipher suites if using TLS 1.3?

**@AntonAM** on `bill-rich/jobworker` [→](https://github.com/bill-rich/jobworker/pull/1#discussion_r1608687287)

> Can you add a bit more details on how we go from mTLS authentication to actually allowing/not allowing to fx stop job from another user.

**@rosstimothy** on `bill-rich/jobworker` [→](https://github.com/bill-rich/jobworker/pull/1#discussion_r1610679199)

> https://go.dev/blog/tls-cipher-suites

**@rosstimothy** on `bill-rich/jobworker` [→](https://github.com/bill-rich/jobworker/pull/1#discussion_r1610681146)

> Can you think of any security concerns with this approach?

**@rosstimothy** on `bucknercd/jobworker` [→](https://github.com/bucknercd/jobworker/pull/1#discussion_r2269786660)

> I'm not certain that this mechanism will achieve the desired goal of avoid cgroup races. Are there any other ways that you can think of the achieve the desired behavior?

**@rosstimothy** on `bucknercd/jobworker` [→](https://github.com/bucknercd/jobworker/pull/1#discussion_r2269799183)

> Overall you've done a good job coming up with a mechanism to identify users and determine what they have access to. I question whether the command lists is a more complicated authorization scheme than is needed for this challenge. It will require parsing the provided commands and arguments, and could be worked around by a clever user. 
> 
> Is there any other scheme you can think of that is both simpler, and still allows for admins to perform some actions and users to perform other actions?

**@zmb3** on `bucknercd/jobworker` [→](https://github.com/bucknercd/jobworker/pull/1#discussion_r2271534095)

> How will the server know which client certificates to trust?
>
> Similarly, how will the CLI client know which client certificate to use?

**@zmb3** on `bucknercd/jobworker` [→](https://github.com/bucknercd/jobworker/pull/1#discussion_r2271537792)

> Up above you said
>
> > authorization done via CA. O=admin|user
>
> This seems to contradict with the `<role-group>` notation mentioned here.

**@zmb3** on `bucknercd/jobworker` [→](https://github.com/bucknercd/jobworker/pull/1#discussion_r2271564506)

> This section doesn't provide much value IMO. You've already discussed mTLS in other parts of the document, the gRPC API is covered in the protobuf spec, and authn/z has its own dedicated section too.

**@rosstimothy** on `bucknercd/jobworker` [→](https://github.com/bucknercd/jobworker/pull/1#discussion_r2276616890)

> Does the allow list here still apply? Is this left over from the previous authorization scheme?

**@rosstimothy** on `bucknercd/jobworker` [→](https://github.com/bucknercd/jobworker/pull/1#discussion_r2276626153)

> Can you help clarify the authorization scheme for me? 
> 
> From [above](https://github.com/bucknercd/jobworker/pull/1/files#diff-3dc5dd454e080eb849ee5efaf79df2585fbe0a06804d69249b4c81da05a63875R86-R89), the authorization scheme is said to be
> 
> > A user can **only** can perform operations such as stream output from any job as long as they have the job id 
> 
> Does that contradict what is being described here?

**@rosstimothy** on `bucknercd/jobworker` [→](https://github.com/bucknercd/jobworker/pull/1#discussion_r2277468330)

> No need to support chroot jail. For this exercise it is fine to run the command passed in verbatim without any extra handling or security considerations. All we are looking to see if that you can identify some of the risks that poses and document how to mitigate them in your design.

**@rosstimothy** on `clydotron/job_worker_service` [→](https://github.com/clydotron/job_worker_service/pull/1#discussion_r1599963437)

> > Do not use any other authentication protocols on top of mTLS
> 
> There should be no need to login, specify the user or add any user tokens in interceptors. Authentication and authorization should strictly rely on certificates exchanged via mTLS.

**@rosstimothy** on `clydotron/job_worker_service` [→](https://github.com/clydotron/job_worker_service/pull/1#discussion_r1599964155)

> Could you add a bit more detail about how you will create certs via openssl?

**@GavinFrazar** on `clydotron/job_worker_service` [→](https://github.com/clydotron/job_worker_service/pull/1#discussion_r1602543075)

> The steps you provided will produce two different key/cert pairs but both will have the same subject:
> ```sh
> $ openssl x509 -in client.pem -noout -subject; openssl x509 -in client_alex.pem -noout -subject
> subject=CN=localhost
> subject=CN=localhost
> ```

**@strideynet** on `clydotron/job_worker_service` [→](https://github.com/clydotron/job_worker_service/pull/1#discussion_r1603089781)

> What versions of TLS will be supported ?

**@espadolini** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/1#discussion_r1782721109)

> Just the one CA for both server and client authn?

**@rosstimothy** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/1#discussion_r1782736583)

> How will the address and certificates be provided?

**@codingllama** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/1#discussion_r1783305359)

> Could you mention your ciphersuites of choice and also the certificate algorithms?

**@codingllama** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/1#discussion_r1786288324)

> What part of the Subject, just for completeness?

**@espadolini** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/1#discussion_r1786408786)

> Keep in mind that the TLS 1.3 ciphersuites are not configurable in `crypto/tls`, so this is currently true but is subject to the will of the stdlib maintainers in the future.

**@espadolini** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/2#discussion_r1787689666)

> Does cfssl not have a way to make a cert that only has the CN in the subject?

**@codingllama** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/2#discussion_r1787747491)

> Could we make sure client certificates won't have these?

**@codingllama** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/2#discussion_r1787755308)

> nit:
> 
> ```suggestion
> 	sed 's/{{CN}}/ca/' make/csr.json | ${CFSSL_BIN} cfssl genkey -initca - | ${CFSSL_BIN} sh -c 'cd certs && cfssljson -bare ca'
> ```
> 
> (No need to follow up on this everywhere, the time spent is not worth it.)

**@codingllama** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/4#discussion_r1788341712)

> ```suggestion
> 				Usage: "Client user name. Cert and key expected in <certdir>.",
> ```

**@espadolini** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/4#discussion_r1789731624)

> The only reason for the `user` argument is to compose some filenames, why not have the certdir be dedicated to a single user instead, like the server certdir?

**@espadolini** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/4#discussion_r1789789168)

> Why reject cert chains?

**@espadolini** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/4#discussion_r1789798586)

> This seems pretty ad-hoc and only sort of works here because we have a small amount of RPCs to deal with, all authorized in mostly the same way, but how can this approach scale?

**@espadolini** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/4#discussion_r1789815241)

> This doesn't do anything, ciphersuites are ignored in TLS 1.3 (and I suspect that the awkwardness of expressing an order based on the availability of AES acceleration is one of the reasons).

**@codingllama** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/4#discussion_r1790700530)

> I _think_ ErrInvalidClientCerts is going to inherit this a godoc.

**@codingllama** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/4#discussion_r1790728896)

> Could we use certs from an unknown CA, so the test is more realistic?

**@rosstimothy** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/4#discussion_r1790807778)

> > Each message potentially contains it's own resource ID and the authorization would need to be checked each time.
> 
> That may be true for a more complicated system, but for this one, it doesn't so much apply. I see now that this is a requirement for authorization to be performed via gRPC middleware.

**@codingllama** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/5#discussion_r1792440751)

> FYI, we considered whether to use the slog.Attr wrappers but decided to go with the "free form" approach, except from certain key scenarios where the Attr actually saves a needless evaluation.

**@rosstimothy** on `devnulled/runjob` [→](https://github.com/devnulled/runjob/pull/2#discussion_r1797078296)

> Isn't server-command a flag on the client CLI?
> ```suggestion
> $ jobmanager serve
>   --server-host="localhost:8443" \
>   --cert-ca-path="certs/cert-ca.crt" \
>   --tls-cert-path="certs/user-tls.crt" \
>   --tls-key-path="certs/user-tls.key"
> ```

**@rosstimothy** on `devnulled/runjob` [→](https://github.com/devnulled/runjob/pull/2#discussion_r1797080716)

> Isn't the server command the argument to start?
> ```suggestion
> $ runjob start "/bin/bash" \
> -A "echo hello world" \
> --server-host="localhost:8443" \
> --cert-ca-path="certs/cert-ca.crt" \
> --tls-cert-path="certs/user-tls.crt" \
> --tls-key-path="certs/user-tls.key"
> ```
> 
> Also couldn't this be executed via instead of providing the shell?
> ```suggestion
> $ runjob start "echo hello world" \
> --server-host="localhost:8443" \
> --cert-ca-path="certs/cert-ca.crt" \
> --tls-cert-path="certs/user-tls.crt" \
> --tls-key-path="certs/user-tls.key"
> ```

**@GavinFrazar** on `devnulled/runjob` [→](https://github.com/devnulled/runjob/pull/2#discussion_r1797290741)

> nit: tls cert/key should be for the server too, not user

**@GavinFrazar** on `devnulled/runjob` [→](https://github.com/devnulled/runjob/pull/2#discussion_r1797310488)

> could you include a script or makefile that generates an example CA and client/server certs/keys?

**@awly** on `diptadas/job-worker` [→](https://github.com/diptadas/job-worker/pull/1#discussion_r646775968)

> ```suggestion
> - `job-worker client [create/stop/status] [--address] [--cacert] [--cert] [--key]` run client to communicate with the API server
> - `job-worker client create [--command]` request a new job
> ```

**@rosstimothy** on `DiscoRiver/jobber` [→](https://github.com/DiscoRiver/jobber/pull/1#discussion_r1989360626)

> Could you add details about the "simple authorization scheme" required for the challenge?

**@rosstimothy** on `dtom90/teleport-backend-challenge` [→](https://github.com/dtom90/teleport-backend-challenge/pull/1#discussion_r2052989776)

> > A subshell was chosen so that the hostname/IP address of the server is preserved across subsequent commands without having to re-enter it
> 
> You could achieve the same by having a flag with a reasonable default for the address of the server. Same for any certs that need to be specified.
> 
> I suggest that you chose the route which is easier to implement at the expense of forcing the users to supply the same flags for each command.

**@zmb3** on `ehsaniara/joblet` [→](https://github.com/ehsaniara/joblet/pull/2#discussion_r2114128551)

> I would keep this simple. I'm not sure how multiple fallback methods help us here. Either the cert is valid and contains the expected data or it doesn't.

**@zmb3** on `ehsaniara/joblet` [→](https://github.com/ehsaniara/joblet/pull/2#discussion_r2114133772)

> TLS 1.3 only sounds good. Note that Go's `crypto/tls` package doesn't allow configuration cipher suites for 1.3 since the defaults are all considered secure by today's standards.

**@russjones** on `flychicken123/Job_Worker` [→](https://github.com/flychicken123/Job_Worker/pull/1#discussion_r2065388702)

> Hrm, this feels like building another layer of authentication. For authorization you typically answer the question: "who can access what on this system?"

**@tigrato** on `flychicken123/Job_Worker` [→](https://github.com/flychicken123/Job_Worker/pull/1#discussion_r2066132703)

> A simple authorization mechanism can be user X can only access his jobs while admins can access every user jobs

**@GavinFrazar** on `flychicken123/Job_Worker` [→](https://github.com/flychicken123/Job_Worker/pull/1#discussion_r2067115090)

> Will we be able to prevent a client cert from being used as a server cert and prevent a server cert being used as a client cert?

**@GavinFrazar** on `flychicken123/Job_Worker` [→](https://github.com/flychicken123/Job_Worker/pull/1#discussion_r2067121583)

> Authentication and authorization are requirements - why is a trusted network assumed and why would that be necessary for authentication and authorization?

**@rosstimothy** on `flychicken123/Job_Worker` [→](https://github.com/flychicken123/Job_Worker/pull/1#discussion_r2069336657)

> As Gavin pointed out authentication and authorization is very much in scope and required for this exercise.

**@GavinFrazar** on `flychicken123/Job_Worker` [→](https://github.com/flychicken123/Job_Worker/pull/1#discussion_r2069373532)

> What if my command also takes flags like `--tls-cert` etc? Won't this be ambiguous to parse?

**@GavinFrazar** on `flychicken123/Job_Worker` [→](https://github.com/flychicken123/Job_Worker/pull/1#discussion_r2069376480)

> I'm fine with in-memory allow lists, depending on the details.
> Please update the design doc to include details about what the whitelist and authorization scheme will look like.

**@eriktate** on `GevorgGal/jobworker` [→](https://github.com/GevorgGal/jobworker/pull/1#discussion_r2975236773)

> Some reasonable, hardcoded defaults would be nice to see. Especially given that `--server-add` and `--ca-cert` will never be different outside of negative test cases

**@eriktate** on `GevorgGal/jobworker` [→](https://github.com/GevorgGal/jobworker/pull/3#discussion_r3003675373)

> ```suggestion
> 	// Non-admin roles can only access explicitly allowed methods.
> 	if r == roleAdmin || viewerAllowed[method] {
> 		return nil
> 	}
>
> 	return errPermissionDenied
> ```
> You might consider inverting the condition here and having your `authorize()` function fail closed instead of failing open. That way if this function evolves over time to be more complex it's harder to accidentally create a set of conditions where access is granted when it shouldn't be.

**@Tener** on `GevorgGal/jobworker` [→](https://github.com/GevorgGal/jobworker/pull/3#discussion_r3007996421)

> We are mixing concerns here. Among other concerns this leads to somewhat convoluted testing code.
>
> Separating authentication from authorization code would be helpful in this regard.

**@Tener** on `GevorgGal/jobworker` [→](https://github.com/GevorgGal/jobworker/pull/3#discussion_r3007996656)

> Testing authorization shouldn't require us to construct elaborate contexts.

**@Tener** on `GevorgGal/jobworker` [→](https://github.com/GevorgGal/jobworker/pull/4#discussion_r3020375372)

> That "h2" warrants a comment, if we need it at all: `credentials.NewTLS(tlsCfg)` already appends it automatically.

**@rosstimothy** on `ghostsquad/prototype-job-worker` [→](https://github.com/ghostsquad/prototype-job-worker/pull/1#discussion_r1531026418)

> Please elaborate on your chosen authorization scheme. What exactly is going to be hard coded?

**@rosstimothy** on `ghostsquad/prototype-job-worker` [→](https://github.com/ghostsquad/prototype-job-worker/pull/1#discussion_r1531029046)

> Please elaborate on how the client and server will be configured for mTLS. How will you be generating certificates to validate that mTLS is working? The challenge explicitly states to cover these things in more detail:
> 
> > Be sure to cover the following in your design: ... TLS setup (version, cipher suites, etc.)...

**@espadolini** on `ghostsquad/prototype-job-worker` [→](https://github.com/ghostsquad/prototype-job-worker/pull/1#discussion_r1537831939)

> I'd leave this as "owner name" or "owner username" and document how users are authenticated when you're defining how authn is going to work.

**@espadolini** on `ghostsquad/prototype-job-worker` [→](https://github.com/ghostsquad/prototype-job-worker/pull/1#discussion_r1537848337)

> Why is the client cert getting a SAN for localhost?

**@espadolini** on `ghostsquad/prototype-job-worker` [→](https://github.com/ghostsquad/prototype-job-worker/pull/1#discussion_r1537850229)

> If you use the same CA for both client and server authentication you have to be _really_ careful about your `extKeyUsage` to avoid malicious actors with user credentials potentially trick other users, or to avoid malicious actors with server credentials potentially use them as client credentials with other servers.

**@espadolini** on `ghostsquad/prototype-job-worker` [→](https://github.com/ghostsquad/prototype-job-worker/pull/1#discussion_r1537884097)

> You should spend a few words on what "current user" means and how authentication ends up attributing a user to a given request.

**@tigrato** on `gstelang/job-worker-service` [→](https://github.com/gstelang/job-worker-service/pull/1#discussion_r1732894851)

> How do you plan to do the authn?
> How will you identify each user and assign him Admin/User roles?

**@strideynet** on `gstelang/job-worker-service` [→](https://github.com/gstelang/job-worker-service/pull/1#discussion_r1735749298)

> A brief sentence justifying why only allowing TLS1.3 would be useful - the reasoning behind decisions is just as useful/interesting as the decision itself.

**@awly** on `hashsequence/Linux-Job-Worker` [→](https://github.com/hashsequence/Linux-Job-Worker/pull/2#discussion_r478700556)

> how will client and server verify each others certificates?

**@awly** on `hashsequence/Linux-Job-Worker` [→](https://github.com/hashsequence/Linux-Job-Worker/pull/4#discussion_r480282389)

> update this to match the `Mutual TLS` section below

**@sclevine** on `ilyazz/jobs` [→](https://github.com/ilyazz/jobs/pull/5#discussion_r1007378738)

> Sad that mkcert doesn't support ed25519

**@rosstimothy** on `jaylane/job-scheduler` [→](https://github.com/jaylane/job-scheduler/pull/1#discussion_r1714211024)

> > The client certificates will have these as configured extensions.
> 
> Which extensions will be used?

**@AntonAM** on `jaylane/job-scheduler` [→](https://github.com/jaylane/job-scheduler/pull/1#discussion_r1714302894)

> How these certificates will be used/configured by the client cli?

**@codingllama** on `jaylane/job-scheduler` [→](https://github.com/jaylane/job-scheduler/pull/1#discussion_r1714327509)

> > The client certificates will have these as X.509 v3 extensions
> 
> Could you be more specific here? Which extensions? How are the user and role determined from the cert?

**@tigrato** on `joshuarubin/teleport-job-worker` [→](https://github.com/joshuarubin/teleport-job-worker/pull/3#discussion_r1743548030)

> Is there another way where we do not need to use the certificate serial number? 
> If the client certificate needs to be re-issued because it's about to expire, will we need to keep the same serial number? Are there other fields we can use instead of relying on a field that's expected to be unique per certificate?

**@rosstimothy** on `joshuarubin/teleport-job-worker` [→](https://github.com/joshuarubin/teleport-job-worker/pull/3#discussion_r1743926955)

> What CA will the clients be configured with to trust the server's certs?

**@rosstimothy** on `joshuarubin/teleport-job-worker` [→](https://github.com/joshuarubin/teleport-job-worker/pull/3#discussion_r1743979130)

> But which CA will it be? Is it the same CA that signed the client certs?

**@codingllama** on `joshuarubin/teleport-job-worker` [→](https://github.com/joshuarubin/teleport-job-worker/pull/3#discussion_r1744192207)

> * Are we using the same CA for both server and clients?
> * Could you call out the certificate algorithm / key length?

**@codingllama** on `joshuarubin/teleport-job-worker` [→](https://github.com/joshuarubin/teleport-job-worker/pull/3#discussion_r1744194222)

> > the client certificate's subject
> 
> What part of the Subject?

**@codingllama** on `joshuarubin/teleport-job-worker` [→](https://github.com/joshuarubin/teleport-job-worker/pull/3#discussion_r1744224876)

> The cgroups flags make sense to me, but we probably don't need the certificate-related ones.

**@rosstimothy** on `joshuarubin/teleport-job-worker` [→](https://github.com/joshuarubin/teleport-job-worker/pull/3#discussion_r1744473761)

> Are there any security concerns from using the same CA for clients and servers?

**@greedy52** on `kiakeshmiri/process-runner` [→](https://github.com/kiakeshmiri/process-runner/pull/1#discussion_r1765682985)

> what TLS version/ciphers are we using? how certs are generated?

**@greedy52** on `kiakeshmiri/process-runner` [→](https://github.com/kiakeshmiri/process-runner/pull/1#discussion_r1765684272)

> I see authentication here but authorization is missing. How does the server identify which user is it? What authorization schema on the api permissions?

**@tigrato** on `kiakeshmiri/process-runner` [→](https://github.com/kiakeshmiri/process-runner/pull/1#discussion_r1766520628)

> Since we are doing authentication, can we also use authorization based on tls certificates?

**@tigrato** on `kiakeshmiri/process-runner` [→](https://github.com/kiakeshmiri/process-runner/pull/1#discussion_r1766526368)

> can you please include the TLS settings you plan to enforce?

**@greedy52** on `kiakeshmiri/process-runner` [→](https://github.com/kiakeshmiri/process-runner/pull/1#discussion_r1766873980)

> Like Tiago mentioned below, for this project, we would prefer to build authorization on the TLS certificate without introducing another authorization flow like OAuth2 token.
> 
> Also i only see the tool for authorization but it doesn't answer the question who will have access to what.

**@tigrato** on `kkloberdanz/teleport-challenge` [→](https://github.com/kkloberdanz/teleport-challenge/pull/1#discussion_r2782015222)

> can you give the details how you plan to generate these certs?

**@smallinsky** on `kkloberdanz/teleport-challenge` [→](https://github.com/kkloberdanz/teleport-challenge/pull/8#discussion_r2822395602)

> Right now, if the client presents a certificate with any other role than the expected first OU in the subject:
>  
> > // First OU from the certificate subject ("admin" or "client")
> 
>  or even a user cert with an empty OU  it may be able to bypass the intended client-role authorization and execute a job.

**@smallinsky** on `kkloberdanz/teleport-challenge` [→](https://github.com/kkloberdanz/teleport-challenge/pull/8#discussion_r2833193668)

> Personally, I would choose a longer validity period since the unit tests rely on the generated certificate files.
>
> Theoretical question: If I had more time and this were production code with a real CI pipeline how would you tackle this ?

**@zmb3** on `kovyrin/teleport-exec` [→](https://github.com/kovyrin/teleport-exec/pull/1#discussion_r788220480)

> Can you think of a way to do authorization without needing access tokens?

**@russjones** on `kovyrin/teleport-exec` [→](https://github.com/kovyrin/teleport-exec/pull/1#discussion_r788287028)

> How do you plan to configure TLS + algorithm for certs?

**@GavinFrazar** on `kurczynski/teleport-job-worker` [→](https://github.com/kurczynski/teleport-job-worker/pull/1#discussion_r1612531424)

> The token seems unnecessary even for this toy project though - users can only access a job that they created right? And we can verify their identity with mTLS.

**@GavinFrazar** on `kurczynski/teleport-job-worker` [→](https://github.com/kurczynski/teleport-job-worker/pull/1#discussion_r1612823096)

> mTLS only establishes identity, i.e. authentication, (as does a secret token), so the authorization scheme is really independent from both of these things.
> 
> And I think just enforcing that only a job creator can interact with the job is fine as a simple authorization scheme for the challenge

**@greedy52** on `m3talsmith/jobberthehut` [→](https://github.com/m3talsmith/jobberthehut/pull/1#discussion_r2499156071)

> How are the certificates generated?
>
> Also loading all client certificates is not sustainable when you have many clients. any way to improve it?

**@gabrielcorado** on `m3talsmith/jobberthehut` [→](https://github.com/m3talsmith/jobberthehut/pull/1#discussion_r2499560612)

> Could add a bit more context on the benefits of using the verification flow and if that somehow overlaps with gRPC mTLS?

**@gabrielcorado** on `m3talsmith/jobberthehut` [→](https://github.com/m3talsmith/jobberthehut/pull/1#discussion_r2499600542)

> Is there a reason for not using the gRPC standard error model (used when returning certificate errors) here? Same comment applies to the other response messages with the `error` field.

**@rosstimothy** on `m3talsmith/jobberthehut` [→](https://github.com/m3talsmith/jobberthehut/pull/1#discussion_r2500030105)

> The challenge requirements are that all jobs have their output monitored, stored, and available for consumption by clients. I suggest for simplicity that you omit this optional behavior.

**@rosstimothy** on `m3talsmith/jobberthehut` [→](https://github.com/m3talsmith/jobberthehut/pull/1#discussion_r2500078453)

> Does this mean that all certificates allowed to interact with the API need to exist in this directory prior to the server launching?

**@rosstimothy** on `m3talsmith/jobberthehut` [→](https://github.com/m3talsmith/jobberthehut/pull/1#discussion_r2500080442)

> Why does the server require an exhaustive list of certs ahead of time?

**@rosstimothy** on `m3talsmith/jobberthehut` [→](https://github.com/m3talsmith/jobberthehut/pull/1#discussion_r2500082518)

> What would happen if my certificate was reissued?

**@rosstimothy** on `m3talsmith/jobberthehut` [→](https://github.com/m3talsmith/jobberthehut/pull/1#discussion_r2500086961)

> Can we move to a system where jobs are tied to an identity rather than a concrete instance of a certificate?

**@rosstimothy** on `m3talsmith/jobberthehut` [→](https://github.com/m3talsmith/jobberthehut/pull/1#discussion_r2500094680)

> There is no need for managing certificates by the service. I do agree with Steve though, let's try to come up with a way that doesn't require a hardcoded list of certificates to be provided to the server.

**@rosstimothy** on `m3talsmith/jobberthehut` [→](https://github.com/m3talsmith/jobberthehut/pull/1#discussion_r2505250910)

> Is there any other means of identifying the owner than a hash of a single certificate? How will owners be able to interact with old jobs if they are forced to reissue their credentials?

**@rosstimothy** on `m3talsmith/jobberthehut` [→](https://github.com/m3talsmith/jobberthehut/pull/1#discussion_r2505252143)

> What does it mean to load a certificate on demand?

**@rosstimothy** on `m3talsmith/jobberthehut` [→](https://github.com/m3talsmith/jobberthehut/pull/1#discussion_r2505256705)

> You may use a UUID library to make this easier for yourself.

**@greedy52** on `MarkDHarris/teleport-systems-challenge` [→](https://github.com/MarkDHarris/teleport-systems-challenge/pull/1#discussion_r3009864353)

> I don't know if SSH or interactive shell is relevant to the goal of the challenge. Just stdin is not part of the requirement. Up to you though.

**@greedy52** on `MarkDHarris/teleport-systems-challenge` [→](https://github.com/MarkDHarris/teleport-systems-challenge/pull/1#discussion_r3009966473)

> Can a client use its own cert to pretend to be a server?

**@rosstimothy** on `MarkDHarris/teleport-systems-challenge` [→](https://github.com/MarkDHarris/teleport-systems-challenge/pull/1#discussion_r3010187451)

> Agree with Steve, the allowed viewers is a nice concept but does add some additional scope. If you want to omit it and only support owner + admin that would still satisfy the challenge requirements.

**@rosstimothy** on `MarkDHarris/teleport-systems-challenge` [→](https://github.com/MarkDHarris/teleport-systems-challenge/pull/1#discussion_r3010211260)

> Is there a particular format we are expecting for the CN and OU? Are we looking for anything specific in either? Are both required to be populated?

**@rosstimothy** on `MarkDHarris/teleport-systems-challenge` [→](https://github.com/MarkDHarris/teleport-systems-challenge/pull/1#discussion_r3015959913)

> Can we enforce EKU for this exercise?

**@rosstimothy** on `MarkDHarris/teleport-systems-challenge` [→](https://github.com/MarkDHarris/teleport-systems-challenge/pull/1#discussion_r3015970179)

> > create a self-signed CA and issue test certs.
>
> Which algorithms will be used when generating the certificates?

**@rosstimothy** on `MarkDHarris/teleport-systems-challenge` [→](https://github.com/MarkDHarris/teleport-systems-challenge/pull/2#discussion_r3035840899)

> I'm curious, what was your thought process that led you to using the callback function?

**@rosstimothy** on `MarkDHarris/teleport-systems-challenge` [→](https://github.com/MarkDHarris/teleport-systems-challenge/pull/2#discussion_r3041727967)

> If changing to a pull based model the io.Reader does map well. Modern Go would probably prefer using a seq.Iter over the callback approach that you initially used as well.

**@tigrato** on `mcampo84/teleport_challenge` [→](https://github.com/mcampo84/teleport_challenge/pull/1#discussion_r1907288571)

> Although this is nice addition, we don't expect you to return this info

**@tigrato** on `mcampo84/teleport_challenge` [→](https://github.com/mcampo84/teleport_challenge/pull/1#discussion_r1907290074)

> for simplicity, you can exclude the timestamp reference. Return the [log output] only

**@tigrato** on `mcampo84/teleport_challenge` [→](https://github.com/mcampo84/teleport_challenge/pull/1#discussion_r1907339666)

> How do you plan to identify the client? Are there simpler ways besides using serials? If the client renews the certificate, the serial will change and he won't be able to continue to use the application

**@tigrato** on `mcampo84/teleport_challenge` [→](https://github.com/mcampo84/teleport_challenge/pull/1#discussion_r1907440228)

> Yes and the job id so you can follow

**@tigrato** on `mcampo84/teleport_challenge` [→](https://github.com/mcampo84/teleport_challenge/pull/1#discussion_r1909346424)

> Can we leverage only data within the certificate otherwise a client could potentially impersonate other client.
> Is there any way of encoding/storing data into x509 certificate that could help us identifying the subject?

**@tigrato** on `mcampo84/teleport_challenge` [→](https://github.com/mcampo84/teleport_challenge/pull/1#discussion_r1915011271)

> with a single chan only one reader will have access to the output and there are no guarantees it's the same that receive the previous line. How do you plan to implement the muti-reader support?

**@tigrato** on `mcampo84/teleport_challenge` [→](https://github.com/mcampo84/teleport_challenge/pull/1#discussion_r1915094617)

> what happens if no new line is found? 
> 
> > ReadBytes reads until the first occurrence of delim in the input, returning a slice containing the data up to and including the delimiter. If ReadBytes encounters an error before finding a delimiter, it returns the data read before the error and the error itself (often io.EOF). ReadBytes returns err != nil if and only if the returned data does not end in delim. For simple uses, a Scanner may be more convenient.
> 
> Can you just forward whatever payload you receive?

**@tigrato** on `mikewurtz/taskman` [→](https://github.com/mikewurtz/taskman/pull/1#discussion_r2027171744)

> For errors are you using gRPC standard error codes?

**@tigrato** on `mikewurtz/taskman` [→](https://github.com/mikewurtz/taskman/pull/1#discussion_r2027179193)

> Can't you just pass the same io.Writer to Stderr and Stdout instead of combining them in a goroutine?

**@tigrato** on `mikewurtz/taskman` [→](https://github.com/mikewurtz/taskman/pull/1#discussion_r2027189594)

> can you make this field optional and populate it only when the program stops?

**@rosstimothy** on `mikewurtz/taskman` [→](https://github.com/mikewurtz/taskman/pull/1#discussion_r2027234354)

> For simplicity you may omit graceful termination altogether.

**@rosstimothy** on `mikewurtz/taskman` [→](https://github.com/mikewurtz/taskman/pull/1#discussion_r2027241438)

> What algorithms will be used to generate the certificates?

**@rosstimothy** on `mikewurtz/taskman` [→](https://github.com/mikewurtz/taskman/pull/1#discussion_r2028820740)

> Can you think of a way to achieve this without forcefully terminating slow clients?

**@timothyb89** on `mikewurtz/taskman` [→](https://github.com/mikewurtz/taskman/pull/1#discussion_r2028841853)

> Do you think it might be simpler to return a gRPC error on failure?

**@rosstimothy** on `mikewurtz/taskman` [→](https://github.com/mikewurtz/taskman/pull/2#discussion_r2034141411)

> Suggestion: use Signal.NotifyContext and plumb the context through instead of using context.Background in the client RPCs so that you honor both your timeout and the users wish to terminate.

**@rosstimothy** on `mikewurtz/taskman` [→](https://github.com/mikewurtz/taskman/pull/2#discussion_r2034150777)

> Just a thought, not at all requesting you change things, but you could alternatively [`embed`](https://pkg.go.dev/embed) the static certs included in the repo.

**@rosstimothy** on `mikewurtz/taskman` [→](https://github.com/mikewurtz/taskman/pull/2#discussion_r2034153330)

> Should it be an error if the auth info is not set, is not `credentials.TLSInfo`, or has no peer certificates?

**@tigrato** on `mikewurtz/taskman` [→](https://github.com/mikewurtz/taskman/pull/2#discussion_r2035476843)

> Given that certs holds this information, should we move these constants there as well?

**@rosstimothy** on `mikewurtz/taskman` [→](https://github.com/mikewurtz/taskman/pull/2#discussion_r2035626553)

> Ah sorry I missed that you were already embedding.

**@rosstimothy** on `mikewurtz/taskman` [→](https://github.com/mikewurtz/taskman/pull/2#discussion_r2036225079)

> Suggestion: move this into your main functions and use [ExecuteContext](https://pkg.go.dev/github.com/spf13/cobra#Command.ExecuteContext) and [cmd.Context()](https://pkg.go.dev/github.com/spf13/cobra#Command.Context) so that you don't need to duplicate this for each subcommand.

**@zmb3** on `mikewurtz/taskman` [→](https://github.com/mikewurtz/taskman/pull/2#discussion_r2036267692)

> Consider naming this `ClientIDKey` or similar.
> 
> "CN" is a concept of the certificate, but now that we've taken the identity out of the certificate and stuffed it in the context the rest of the code doesn't really care that it originally came from a cert's common name.

**@zmb3** on `mikewurtz/taskman` [→](https://github.com/mikewurtz/taskman/pull/2#discussion_r2036270523)

> I can see why you used the sync.Once to share a single server between all tests, but I'm not sure the extra complexity is worth it here.
> 
> One thing I don't love about this approach is it creates asymmetry between `startTestServer` (which is called by each test), and `stopServer` (which is called once).

**@tigrato** on `mikewurtz/taskman` [→](https://github.com/mikewurtz/taskman/pull/4#discussion_r2039950334)

> you can use https://github.com/olekukonko/tablewriter

**@tigrato** on `mikewurtz/taskman` [→](https://github.com/mikewurtz/taskman/pull/4#discussion_r2039973402)

> should you call os.RemoveAll
> If a directory isn't empty how does calling over and over again will solve it?

**@tigrato** on `mikewurtz/taskman` [→](https://github.com/mikewurtz/taskman/pull/4#discussion_r2040098359)

> can't you use errors.As api instead of calling `IsTaskError` that instead of returning a bool also returns the error?

**@rosstimothy** on `mikewurtz/taskman` [→](https://github.com/mikewurtz/taskman/pull/5#discussion_r2047081919)

> Food for thought, but not something I expect you to address here: I think you could probably achieve the same result without the additional goroutine if you implemented a custom context.Context.

**@oeric** on `moalf/teleport` [→](https://github.com/moalf/teleport/pull/1#discussion_r2594169484)

> Can you think about some additional controls that you'd enforce to further secure this?

**@oeric** on `moalf/teleport` [→](https://github.com/moalf/teleport/pull/1#discussion_r2594171215)

> Can you list the scopes you'd need for this challenge?

**@russjones** on `moalf/teleport` [→](https://github.com/moalf/teleport/pull/1#discussion_r2594335839)

> Nice that you are thinking about this, but not needed for this challenge. Feel free to drop to reduce scope.

**@oeric** on `moalf/teleport` [→](https://github.com/moalf/teleport/pull/1#discussion_r2599115113)

> I like the first approach, can you add that to your doc?

**@russjones** on `moalf/teleport` [→](https://github.com/moalf/teleport/pull/1#discussion_r2604438803)

> This is close, but you may need to flesh it out a bit more when you get to implementation and you've hammered down the exact endpoints you want to hit.

**@russjones** on `moalf/teleport` [→](https://github.com/moalf/teleport/pull/1#discussion_r2604442934)

> Right, but you can't enforce conversation resolution before merge on GitHub as far as I am aware.
> 
> How would you enforce code review before merge on GitHub?

**@oeric** on `moalf/teleport` [→](https://github.com/moalf/teleport/pull/1#discussion_r2604532691)

> Yep, there are some additional scopes you'll probably need. Can you take another look and add those to your doc?
> 
> It may be helpful to think about your entire implementation, and not just your user creation script.

**@russjones** on `moalf/teleport` [→](https://github.com/moalf/teleport/pull/2#discussion_r2612551939)

> Do you need the access token for tenant configuration?

**@russjones** on `moalf/teleport` [→](https://github.com/moalf/teleport/pull/2#discussion_r2612576290)

> This might work, I would have to spend some more time reading about post-login actions. I assume you're handling situations where users can click "No thanks" to bypass enrollment?
> 
> https://community.auth0.com/t/auth0-trigger-post-login-user-can-skip-webauth-enrollment-by-clicking-no-thanks-and-the-next-subsequent-action-is-called/169362

**@oeric** on `moalf/teleport` [→](https://github.com/moalf/teleport/pull/2#discussion_r2615132413)

> Strangely I'm unable to recreate this issue (both in an existing tenant, and a brand new Auth0 account + tenant).
>
> Are you starting from a brand new tenant?

**@oeric** on `moalf/teleport` [→](https://github.com/moalf/teleport/pull/2#discussion_r2620600950)

> Does this ensure that Webauthn is the only allowed factor?
>
> Also, about `recovery_code` - what are some of the benefits/risks of allowing that? Would you use that in a corporate environment if Auth0 was an IDP protecting sensitive resources?

**@r0mant** on `moalf/teleport` [→](https://github.com/moalf/teleport/pull/2#discussion_r2621368315)

> I'm not sure why this script is needed? `terraform init` will find and install the provider automatically for you if the provider block is present in the tf file, right?

**@oeric** on `moalf/teleport` [→](https://github.com/moalf/teleport/pull/2#discussion_r2625139958)

> Ah sorry, must've missed that in your design doc PR.
>
> I agree with your assessment, you can turn that off.

**@oeric** on `moalf/teleport` [→](https://github.com/moalf/teleport/pull/3#discussion_r2628643715)

> Curious why you added this?

**@russjones** on `moalf/teleport` [→](https://github.com/moalf/teleport/pull/3#discussion_r2628869210)

> Can you explain this a bit more. How does disabling the login form increase security?

**@nklaassen** on `MrChristianL/Teleport-Job-Worker-Service` [→](https://github.com/MrChristianL/Teleport-Job-Worker-Service/pull/1#discussion_r2879031767)

> A couple concrete things I'm looking for:
>
> What is the exported API of the log broker, how will your API server call it?
>
> What will you set cmd.Stdout and cmd.Stderr to?

**@nklaassen** on `MrChristianL/Teleport-Job-Worker-Service` [→](https://github.com/MrChristianL/Teleport-Job-Worker-Service/pull/1#discussion_r2879098571)

> Not sure I can agree that "SAN-based authorization is a modern X.509 standard of identity management" unless there is some standard I have not read
>
> But really, the SAN is still part of _authentication_ here, not _authorization_. The SAN ties the mTLS certificate to an identity/role, then your code does authorization based on that authenticated identity

**@eriktate** on `MrChristianL/Teleport-Job-Worker-Service` [→](https://github.com/MrChristianL/Teleport-Job-Worker-Service/pull/1#discussion_r2881188654)

> So the identity -> role map is hard coded into the application? Did you consider any other approaches for determining which role is used for authZ decisions?

**@eriktate** on `MrChristianL/Teleport-Job-Worker-Service` [→](https://github.com/MrChristianL/Teleport-Job-Worker-Service/pull/1#discussion_r2886297267)

> > In a large production-ready implementation, we would likely opt for something like UUIDs to identify the jobs rather than a home-brewed solution.
> 
> You're welcome to use UUIDs in the challenge if you think that would save you any time.
> 
> > Even though there are over 2 billion possible IDs, the function still checks to ensure any assigned ID is not already present in the job tracker before confirming the jobID.
> 
> What will happen if the new ID conflicts with an existing job?

**@eriktate** on `MrChristianL/Teleport-Job-Worker-Service` [→](https://github.com/MrChristianL/Teleport-Job-Worker-Service/pull/1#discussion_r2886373820)

> > 2. Writes and syncs to disk (still holding the lock in order to serialize writers)
> 
> Do you expect there to be multiple writers for a particular log file?

**@rosstimothy** on `MrChristianL/Teleport-Job-Worker-Service` [→](https://github.com/MrChristianL/Teleport-Job-Worker-Service/pull/1#discussion_r2890717401)

> Can you think of any way to reduce the number of concurrent writers to one?

**@rosstimothy** on `MrChristianL/Teleport-Job-Worker-Service` [→](https://github.com/MrChristianL/Teleport-Job-Worker-Service/pull/1#discussion_r2890772977)

> You may reduce scope by cutting this.

**@nklaassen** on `MrChristianL/Teleport-Job-Worker-Service` [→](https://github.com/MrChristianL/Teleport-Job-Worker-Service/pull/1#discussion_r2897009660)

> suggestion: this is so close to role-based access control which sounds much nicer than "a hardcoded identity map". If you just define `admin` and `user` as "roles" that are encoded in the user certificate that could theoretically be assigned to multiple users by whatever issues the certificates, it doesn't really even change your design but sounds better than hardcoded identities. That might even be what you're trying to say, but the use of "hardcoded identity" is throwing me off

**@nklaassen** on `MrChristianL/Teleport-Job-Worker-Service` [→](https://github.com/MrChristianL/Teleport-Job-Worker-Service/pull/1#discussion_r2897026273)

> nit: I don't really see SPIFFE/SPIRE as an alternative to your authorization scheme or a policy engine. It could potentially fit as an answer to "how do I securely bootstrap trust and distribute these X.509 certificates to servers and clients", which is out of scope for the challenge, but it doesn't really address authorization

**@eriktate** on `MrChristianL/Teleport-Job-Worker-Service` [→](https://github.com/MrChristianL/Teleport-Job-Worker-Service/pull/2#discussion_r2905604470)

> Your `go.mod` is missing dependencies (e.g. `google.golang.org/grpc`). You're also missing a `go.sum` file

**@eriktate** on `MrChristianL/Teleport-Job-Worker-Service` [→](https://github.com/MrChristianL/Teleport-Job-Worker-Service/pull/2#discussion_r2906144937)

> You could combine this and the previous call to `j.onceDone.Do()` into a single defer.

**@eriktate** on `MrChristianL/Teleport-Job-Worker-Service` [→](https://github.com/MrChristianL/Teleport-Job-Worker-Service/pull/2#discussion_r2906211680)

> nit: this is a bit confusing to parse at first glance and you could cut down on the number of places `test.wantErr` is evaluated
>
> ```suggestion
> 			if test.wantErr {
> 			    if err == nil {
> 					t.Errorf("newJob() error = %v, wantErr %v", err, test.wantErr)
> 				}
> 				return
> 			}
> 			if job == nil {
> 				t.Errorf("expected newJob() not to be nil, wantErr %v", test.wantErr)
> 			}
> 			if job.ID != test.id {
> 				t.Errorf("expected ID %q, got %q", test.id, job.ID)
> 			}
> ```

**@eriktate** on `MrChristianL/Teleport-Job-Worker-Service` [→](https://github.com/MrChristianL/Teleport-Job-Worker-Service/pull/2#discussion_r2906285930)

> Did you consider using `rand.Text()` for this?

**@eriktate** on `MrChristianL/Teleport-Job-Worker-Service` [→](https://github.com/MrChristianL/Teleport-Job-Worker-Service/pull/2#discussion_r2906304283)

> Looks like you committed a merge conflict

**@Joerger** on `MrChristianL/Teleport-Job-Worker-Service` [→](https://github.com/MrChristianL/Teleport-Job-Worker-Service/pull/2#discussion_r2915149234)

> I don't have a strong opinion here, I just like the idea of the system being communicative about what action actually occurred. For example, if you thought a job was frozen, but it actually completed when you went to stop it and retry, you would start the job again instead instead of realizing it was actually completed and just checking the final logs that you wanted.
> 
> Your justification is definitely valid though, and simplicity is key in this challenge, so no need to change it.

**@rosstimothy** on `MrChristianL/Teleport-Job-Worker-Service` [→](https://github.com/MrChristianL/Teleport-Job-Worker-Service/pull/3#discussion_r2924641780)

> Did you consider making this a `map[string]map[string]struct{}`? That would simplify `isAuthorized` to the following.
>
>
> ```go
> func isAuthorized(role string, method string) bool {
> 	allowedMethods, ok := rolePermissions[role]
> 	if !ok {
> 		return false // role has no permissions defined
> 	}
>
> 	_, allowed := allowedMethods[method]
> 	return allowed
> }
> ```

**@rosstimothy** on `MrChristianL/Teleport-Job-Worker-Service` [→](https://github.com/MrChristianL/Teleport-Job-Worker-Service/pull/3#discussion_r2924669559)

> Disclaimer: I'm not asking for changes, I'm only interested in your thoughts.
>
> The interceptor approach works, but is redundant. The client certificate will not change over the lifetime of the connection. However, the interceptor will perform the same operations to extract the identity from the certificate for every RPC. What could we do to extract the identity from the certificate a single time?

**@rosstimothy** on `MrChristianL/Teleport-Job-Worker-Service` [→](https://github.com/MrChristianL/Teleport-Job-Worker-Service/pull/3#discussion_r2924680509)

> How could you protect against forgetting to update the authorization model if a new RPC was added?

**@hugoShaka** on `neildo/tjob` [→](https://github.com/neildo/tjob/pull/1#discussion_r1758878917)

> What does the client validate when looking at the server cert? The SANs must match the domain the client is connecting to?
> 
> Is it possible for someone with a client cert to impersonate the server if they are signed by the same CA?

**@hugoShaka** on `neildo/tjob` [→](https://github.com/neildo/tjob/pull/1#discussion_r1758887090)

> Thank you for the proof of concepts/examples, this helps a lot.
> 
> You described how to create a cgroup and set its limits. Could you also describe how the `tjob` library will put the child process in the cgroup?

**@nklaassen** on `neildo/tjob` [→](https://github.com/neildo/tjob/pull/1#discussion_r1759260390)

> not looking for it at all right now but just a hint that we like the project to have a useful README with examples of how to run the project by the time you're done

**@nklaassen** on `neildo/tjob` [→](https://github.com/neildo/tjob/pull/1#discussion_r1759261277)

> ```suggestion
> For mutual authentication (mTLS), both client and server requires certificates from trusted certificate authority (CA). Instead of well-known CA like Verisign, OpenSSL can generate the RSA 256 certificates required for the CA, clients, and servers for the **scope of this prototype.**
> ```

**@nklaassen** on `neildo/tjob` [→](https://github.com/neildo/tjob/pull/1#discussion_r1759266557)

> can you please mention the trust relationships and any self-signed certs or CAs you will be using?
> 
> also you mention RSA 256, can you be more specific? actually, is there better key algorithm you can use than RSA?
> 
> we also asked for a "strong set of cipher suites for TLS", now this is pretty easy to get right in Go these days but it warrant some mention, and I also like to see a mention of which TLS versions you will support

**@nklaassen** on `neildo/tjob` [→](https://github.com/neildo/tjob/pull/1#discussion_r1759286041)

> I don't think that's right, what is an EOF and how could you split on it? the api should work well even if the job doesn't output any concept of "lines"

**@nklaassen** on `neildo/tjob` [→](https://github.com/neildo/tjob/pull/1#discussion_r1759287331)

> why is the mutex embedded, do you want Job to act as a mutex? and why as a pointer?

**@rosstimothy** on `neildo/tjob` [→](https://github.com/neildo/tjob/pull/1#discussion_r1761092126)

> Note that Go doesn't allow configuring cipher suites if using TLS 1.3: https://go.dev/blog/tls-cipher-suites

**@rosstimothy** on `neildo/tjob` [→](https://github.com/neildo/tjob/pull/1#discussion_r1761098514)

> Can you think of any alternatives to sleeping and trying again in the future?

**@rosstimothy** on `neildo/tjob` [→](https://github.com/neildo/tjob/pull/1#discussion_r1761103121)

> I'm a bit confused about what is being proposed here. `lsblk` is showing a number of disks, but io.max is only being configured for disk `8:0`. Do you plan on applying to all disks? Or are you saying this will only be applied to a single disk?

**@nklaassen** on `neildo/tjob` [→](https://github.com/neildo/tjob/pull/1#discussion_r1761499721)

> To be honest i'm not sure what you mean by AES 256 GSM certificates. You mention a TLS_AES_256_GCM_SHA384 cipher suite for TLS, but this is not part of the certificate. Can you mention which key and signature algorithm you will use for the CA certificates?

**@nklaassen** on `neildo/tjob` [→](https://github.com/neildo/tjob/pull/1#discussion_r1761779717)

> alright, ed25519 is a fine choice as long as you control the server and all clients, but it might be worth a note that many TLS clients (including all major web browsers) don't support Ed25519 certs

**@nklaassen** on `neildo/tjob` [→](https://github.com/neildo/tjob/pull/1#discussion_r1761781646)

> that warning against watching changes to a file doesn't really apply here, because you control the file and know it won't be moved/overwritten
> @rosstimothy i always forget, do we allow an import for fsnotify?

**@rosstimothy** on `neildo/tjob` [→](https://github.com/neildo/tjob/pull/1#discussion_r1761995191)

> I don't know that you need to import that entire library to use inotify. I do agree with Nic that for the purpose of this challenge that warning can be ignored.

**@rosstimothy** on `neildo/tjob` [→](https://github.com/neildo/tjob/pull/1#discussion_r1764116209)

> What is this proof of concept demonstrating? I still don't follow what the purpose of the section is? Is this supposed to be pseudocode for how the server will spawn a job? Could you provide some clarity here.

**@nklaassen** on `neildo/tjob` [→](https://github.com/neildo/tjob/pull/2#discussion_r1767822510)

> messing with global variables is generally racy and hard to test, and why set it to true here if you might just set it back to false further down?

**@nklaassen** on `neildo/tjob` [→](https://github.com/neildo/tjob/pull/2#discussion_r1767824735)

> panicking isn't very cool, could you return an error instead?

**@nklaassen** on `neildo/tjob` [→](https://github.com/neildo/tjob/pull/2#discussion_r1767825083)

> could you make the API prevent the possibility of starting the same job twice?

**@nklaassen** on `neildo/tjob` [→](https://github.com/neildo/tjob/pull/2#discussion_r1767897885)

> why do you need to add and remove the watch on each read? isn't this a bit racy, you could read EOF, and then more could be written and the file closed before you add the watch, and then reading from the inotify file would block

**@smallinsky** on `P-A-R-U-S/Go-Job-Worker-Service` [→](https://github.com/P-A-R-U-S/Go-Job-Worker-Service/pull/1#discussion_r1730890998)

> For simplify you can hardcore the limits on the server side.

**@smallinsky** on `P-A-R-U-S/Go-Job-Worker-Service` [→](https://github.com/P-A-R-U-S/Go-Job-Worker-Service/pull/1#discussion_r1732262361)

> nit: 
> 
> Could you add undefined/unspecified value: https://protobuf.dev/programming-guides/dos-donts/#unspecified-enum

**@smallinsky** on `P-A-R-U-S/Go-Job-Worker-Service` [→](https://github.com/P-A-R-U-S/Go-Job-Worker-Service/pull/2#discussion_r1744967038)

> Please use your personal judgment if you tested this manually and it works for purpose of this challangage test coverage is not needed. 
> 
> However, I believe you need to set:
> ```
> &syscall.SysProcAttr{
>   Setpgid: true
> }
>  ```
>  
> to make it work.

**@greedy52** on `P-A-R-U-S/Go-Job-Worker-Service` [→](https://github.com/P-A-R-U-S/Go-Job-Worker-Service/pull/3#discussion_r1750878162)

> TLS version should be set.

**@rhammonds-teleport** on `rajansandeep/teleport-job-worker-svc` [→](https://github.com/rajansandeep/job-worker-svc/pull/1#discussion_r3079881619)

> ```suggestion
> Every server should also present a valid certificate using `tls.RequireAndVerifyClientCert` to prevent anonymous client connections.
> ```
> Just a typo?

**@rhammonds-teleport** on `rajansandeep/teleport-job-worker-svc` [→](https://github.com/rajansandeep/job-worker-svc/pull/1#discussion_r3079920732)

> Looks like you're using a UUID for the job ID. Could you update the doc to mention which package/component is responsible for generating it?

**@eriktate** on `rajansandeep/teleport-job-worker-svc` [→](https://github.com/rajansandeep/job-worker-svc/pull/1#discussion_r3088728942)

> Feel free to use the google/uuid package if it saves you time. Generating UUIDs from scratch and complying with RFC 9562 isn't part of the challenge scope

**@eriktate** on `rajansandeep/teleport-job-worker-svc` [→](https://github.com/rajansandeep/job-worker-svc/pull/1#discussion_r3088878812)

> It looks like the worker manager is responsible for much more than that. You might consider a quick summary before jumping into the API below

**@eriktate** on `rajansandeep/teleport-job-worker-svc` [→](https://github.com/rajansandeep/job-worker-svc/pull/1#discussion_r3088887146)

> Feel free to omit this unless it's helpful for you during implementation.

**@eriktate** on `rajansandeep/teleport-job-worker-svc` [→](https://github.com/rajansandeep/job-worker-svc/pull/1#discussion_r3088943501)

> Can you include more information how you will approach certs? Will they be pre-generated? Which algorithms will they use? Any additional configuration?

**@eriktate** on `rajansandeep/teleport-job-worker-svc` [→](https://github.com/rajansandeep/job-worker-svc/pull/1#discussion_r3088960337)

> It seems that the CN is a role rather than a unique identifier. I'm not sure it's appropriate as the owner of a job considering multiple users could share the same role.

**@eriktate** on `rajansandeep/teleport-job-worker-svc` [→](https://github.com/rajansandeep/job-worker-svc/pull/1#discussion_r3093603276)

> It looks like `Owner` was dropped from the `job` definition. Do we still need to move away from roles in the CN? If we continue with `alice`, `bob`, and `charlie`, how will they map to the `admin` and `viewer` roles?

**@eriktate** on `rajansandeep/teleport-job-worker-svc` [→](https://github.com/rajansandeep/job-worker-svc/pull/1#discussion_r3095040305)

> It would also be totally possible to encode both the user identifier and the assigned role in the certificate without having to defer to a secondary mapping to figure out which role to apply, but I agree this would be a totally viable approach as well :+1:

**@rhammonds-teleport** on `rajansandeep/teleport-job-worker-svc` [→](https://github.com/rajansandeep/job-worker-svc/pull/2#discussion_r3100975818)

> `synctest` was added to the standard library in Go 1.25. I believe it could help you remove this sleep and the sleep below.
> https://go.dev/blog/synctest

**@atburke** on `rajansandeep/teleport-job-worker-svc` [→](https://github.com/rajansandeep/job-worker-svc/pull/2#discussion_r3101734503)

> [strings.Repeat](https://pkg.go.dev/strings#Repeat) (or its byte equivalent) may be of interest to you.

**@atburke** on `rajansandeep/teleport-job-worker-svc` [→](https://github.com/rajansandeep/job-worker-svc/pull/2#discussion_r3101788424)

> You might be able to get rid of `waitGroupOrTimeout()` with synctest too.

**@atburke** on `rajansandeep/teleport-job-worker-svc` [→](https://github.com/rajansandeep/job-worker-svc/pull/2#discussion_r3103402262)

> Nit: as of go 1.22 you can do this:
> ```suggestion
> 			for range totalLines {
> ```

**@atburke** on `rajansandeep/teleport-job-worker-svc` [→](https://github.com/rajansandeep/job-worker-svc/pull/2#discussion_r3103535246)

> Can you add a comment or update the test name to mention that? Just so future readers aren't confused by an assert-less test.

**@atburke** on `rajansandeep/teleport-job-worker-svc` [→](https://github.com/rajansandeep/job-worker-svc/pull/2#discussion_r3103705044)

> You can add a TODO for now. I suspect you'll want a Close() by the time you implement the grpc server so you'll be able to shut down gracefully.

**@atburke** on `rajansandeep/teleport-job-worker-svc` [→](https://github.com/rajansandeep/job-worker-svc/pull/4#discussion_r3120000677)

> Nit: `proto` appears to be the default make target. I suspect that you'll want it to be one of the build targets, or maybe a new target that builds everything. Also, an e2e target would be nice.

**@rhammonds-teleport** on `rajansandeep/teleport-job-worker-svc` [→](https://github.com/rajansandeep/job-worker-svc/pull/4#discussion_r3124527634)

> You also have the `type tlsPaths struct` a few lines down which would help disambiguate the returned strings.

**@atburke** on `rajansandeep/teleport-job-worker-svc` [→](https://github.com/rajansandeep/job-worker-svc/pull/5#discussion_r3126102585)

> Oh whoops, you're right! My brain misread the check as `info == nil`.

**@zmb3** on `razzam21/job-worker-service` [→](https://github.com/razzam21/job-worker-service/pull/1#discussion_r2255194182)

> This is a great overview of what the gRPC API will look like at the server layer.
>
> Could you provide a similar overview of what the Go API will look like at the library layer?

**@zmb3** on `razzam21/job-worker-service` [→](https://github.com/razzam21/job-worker-service/pull/1#discussion_r2255197463)

> Any concerns with using the same CA to issue both client and server certs?

**@zmb3** on `razzam21/job-worker-service` [→](https://github.com/razzam21/job-worker-service/pull/1#discussion_r2255200250)

> What do you think about eliminating TLS 1.2 fallback completely? We own the only client, so _broad support across systems_ is not a major concern.

**@zmb3** on `razzam21/job-worker-service` [→](https://github.com/razzam21/job-worker-service/pull/1#discussion_r2255201062)

> This looks good, though note that in Go the cipher suites for TLS 1.3 are not configurable.

**@zmb3** on `razzam21/job-worker-service` [→](https://github.com/razzam21/job-worker-service/pull/1#discussion_r2255204587)

> Note that Go 1.24 ships with `crypto/rand.Text()`, which simplifies this a bit.
>
> (Though `google/uuid` is also fine if you prefer UUIDs)

**@creack** on `razzam21/job-worker-service` [→](https://github.com/razzam21/job-worker-service/pull/1#discussion_r2257512345)

> Looks like you use specific version a bit everywhere, but not here? Is that expected?

**@zmb3** on `razzam21/job-worker-service` [→](https://github.com/razzam21/job-worker-service/pull/1#discussion_r2257555985)

> This is fine, though I wonder if you considered any alternatives?
>
> By adding owner to the public API of the library, you've let your authorization scheme dictate the API of the library. Ideally we'd want the authorization scheme to be something that's applied at the server, not in the library.

**@creack** on `razzam21/job-worker-service` [→](https://github.com/razzam21/job-worker-service/pull/1#discussion_r2257559111)

> How would we select which is the current user? Could you include an example?

**@awly** on `renatoaguimaraes/golang-job-scheduler` [→](https://github.com/renatoaguimaraes/golang-job-scheduler/pull/2#discussion_r623310046)

> optional: `github.com/renatoaguimaraes/job-scheduler` is more appropriate as module name, but you'd have to update all the imports too

**@awly** on `renatoaguimaraes/golang-job-scheduler` [→](https://github.com/renatoaguimaraes/golang-job-scheduler/pull/2#discussion_r623372897)

> you should be able to use just the `syscall` package for all the inotify calls

**@dboslee** on `renatoaguimaraes/golang-job-scheduler` [→](https://github.com/renatoaguimaraes/golang-job-scheduler/pull/2#discussion_r623388994)

> I think it would be better to let the logs persist after the process exits. Otherwise you would be racing to view the logs, especially for short lived jobs.

**@r0mant** on `renatoaguimaraes/golang-job-scheduler` [→](https://github.com/renatoaguimaraes/golang-job-scheduler/pull/2#discussion_r623539401)

> You could just pass the job object from Start here since you already have it, then this wouldn't be needed.

**@r0mant** on `renatoaguimaraes/golang-job-scheduler` [→](https://github.com/renatoaguimaraes/golang-job-scheduler/pull/2#discussion_r625190630)

> Hmm, I may be missing something, but it seems like `offset` is not really used anymore? You're not doing Seek anymore.

**@Joerger** on `renatoaguimaraes/golang-job-scheduler` [→](https://github.com/renatoaguimaraes/golang-job-scheduler/pull/5#discussion_r625366009)

> These tests don't really affirm that authorization is working from the top level. Could you add short authorization tests for the gRPC endpoints themselves instead?

**@awly** on `renatoaguimaraes/golang-job-scheduler` [→](https://github.com/renatoaguimaraes/golang-job-scheduler/pull/6#discussion_r625937259)

> ```suggestion
> 		os.Stderr.WriteString("you must pass a command")
> ```

**@strideynet** on `rexposadas/teleport` [→](https://github.com/rexposadas/teleport/pull/1#discussion_r1790414746)

> Could you provide the specification for the request/response messages as well? It'd be good to double-check the design of these.

**@strideynet** on `rexposadas/teleport` [→](https://github.com/rexposadas/teleport/pull/1#discussion_r1790415956)

> Which part of the certificate will identify the user?

**@strideynet** on `rexposadas/teleport` [→](https://github.com/rexposadas/teleport/pull/1#discussion_r1790417594)

> How will you place the process into the cgroup ?

**@strideynet** on `rexposadas/teleport` [→](https://github.com/rexposadas/teleport/pull/1#discussion_r1790422786)

> Please include what version of TLS, cipher suites etc, you will support and why

**@rosstimothy** on `rexposadas/teleport` [→](https://github.com/rexposadas/teleport/pull/1#discussion_r1790654161)

> What version of TLS will be enforced? What cipher suites will be used?

**@rosstimothy** on `rexposadas/teleport` [→](https://github.com/rexposadas/teleport/pull/1#discussion_r1790655421)

> Is there a single CA for both client and server certificates?

**@strideynet** on `rexposadas/teleport` [→](https://github.com/rexposadas/teleport/pull/1#discussion_r1793031064)

> Can you provide justification for choosing TLS 1.3, and which ciphers will be used

**@rosstimothy** on `rexposadas/teleport` [→](https://github.com/rexposadas/teleport/pull/1#discussion_r1793931947)

> Do note that Go does not allow specifying these if using TLS 1.3: https://go.dev/blog/tls-cipher-suites.

**@rosstimothy** on `rexposadas/teleport` [→](https://github.com/rexposadas/teleport/pull/1#discussion_r1793934126)

> Simply inspecting the CN seems like an extension of authentication and not an authorization scheme. Identifying users by the CN is totally valid for this challenge, however, there should also be some limiting factors applied based on that identity to satisfy the challenge authorization requirement.

**@rosstimothy** on `rexposadas/teleport` [→](https://github.com/rexposadas/teleport/pull/1#discussion_r1797434372)

> What algorithm will the premade certificates be using? How do you plan on generating them?

**@greedy52** on `rexposadas/teleport` [→](https://github.com/rexposadas/teleport/pull/2#discussion_r1803057868)

> could you justify the one second timeout?

**@rosstimothy** on `RichyHBM/teleport-challenge` [→](https://github.com/RichyHBM/teleport-challenge/pull/1#discussion_r2048932914)

> I assume this is 2048 RSA certificates? What do you think about using an elliptic-curve for this challenge instead?

**@rosstimothy** on `RichyHBM/teleport-challenge` [→](https://github.com/RichyHBM/teleport-challenge/pull/1#discussion_r2048937015)

> What proof will the server have the the username provided is in fact accurate? Can you think of a more secure and less falsifiable way of establishing the identity of the user?

**@rosstimothy** on `RichyHBM/teleport-challenge` [→](https://github.com/RichyHBM/teleport-challenge/pull/1#discussion_r2048939083)

> Could you include what limits you plan on hard cording, and which cgroups v2 controllers they will be applied to?

**@rosstimothy** on `RichyHBM/teleport-challenge` [→](https://github.com/RichyHBM/teleport-challenge/pull/1#discussion_r2049501095)

> That could work, though I think you mean the public key, not the private key?

**@rosstimothy** on `RichyHBM/teleport-challenge` [→](https://github.com/RichyHBM/teleport-challenge/pull/1#discussion_r2049509174)

> > These permissions will just be hard coded
> 
> What will the permissions be? You've specified how you will be identifying users for your scheme, but you've not articulated what the scheme will enforce based on those identities.

**@tigrato** on `RichyHBM/teleport-challenge` [→](https://github.com/RichyHBM/teleport-challenge/pull/1#discussion_r2049665512)

> How will you send the user to the server?

**@rosstimothy** on `RichyHBM/teleport-challenge` [→](https://github.com/RichyHBM/teleport-challenge/pull/1#discussion_r2050608172)

> It sounds like what you are proposing is an allow list of commands per user then?

**@rosstimothy** on `RichyHBM/teleport-challenge` [→](https://github.com/RichyHBM/teleport-challenge/pull/1#discussion_r2050609778)

> Please note in the document that you intend to enumerate and apply the limits to all devices then.

**@rosstimothy** on `RichyHBM/teleport-challenge` [→](https://github.com/RichyHBM/teleport-challenge/pull/2#discussion_r2056051824)

> Use a random port to prevent test failures if something else is listening on 5555?
> ```suggestion
> 	grpcServer, listener, err := createGrpcServer(0, certFile, keyFile, certAuthorityFile)
> ```

**@rosstimothy** on `RichyHBM/teleport-challenge` [→](https://github.com/RichyHBM/teleport-challenge/pull/2#discussion_r2056055589)

> Can you get the address form the listener directly?
> ```suggestion
> 	log.Printf("Starting up on %s", listener.Addr())
> ```

**@zmb3** on `RichyHBM/teleport-challenge` [→](https://github.com/RichyHBM/teleport-challenge/pull/2#discussion_r2056547744)

> Why embed these files instead of reading them from disk?
> 
> Does the build of the CLI client really need to contain the private key for the server's CA?
> 
> Also, does the wildcard here mean you're actually embedding `certs.go` in the binary too?

**@zmb3** on `RichyHBM/teleport-challenge` [→](https://github.com/RichyHBM/teleport-challenge/pull/2#discussion_r2056559087)

> I'm looking at https://github.com/grpc/grpc-go/tree/master/examples/features/gracefulstop and I see:
> 
> > It's crucial to call Server.Stop() with a timeout before calling GracefulStop(). This acts as a safety net, ensuring that the server eventually shuts down even if some in-flight RPCs don't complete within a reasonable timeframe. This prevents indefinite blocking.
> 
> I'm also not sure if you need to close the listener - it seems like the server should do that for you.

**@eriktate** on `RichyHBM/teleport-challenge` [→](https://github.com/RichyHBM/teleport-challenge/pull/2#discussion_r2058508702)

> Should the `--` always be required before providing the remote command? I think having support for it can increase readability, but I wonder if the UX would be better if you could omit it in cases where there's no ambiguity?

**@rosstimothy** on `RichyHBM/teleport-challenge` [→](https://github.com/RichyHBM/teleport-challenge/pull/3#discussion_r2060444798)

> Suggestion: invert some of the conditions so that you can return early and reduce the level of indentation here.

**@tigrato** on `RichyHBM/teleport-challenge` [→](https://github.com/RichyHBM/teleport-challenge/pull/3#discussion_r2060518274)

> should you delete all if the folder has items?

**@tigrato** on `RichyHBM/teleport-challenge` [→](https://github.com/RichyHBM/teleport-challenge/pull/4#discussion_r2060521994)

> why do you need the username?

**@tigrato** on `RichyHBM/teleport-challenge` [→](https://github.com/RichyHBM/teleport-challenge/pull/4#discussion_r2060528369)

> Why do you need `HasAnyAuthorization`?

**@rosstimothy** on `RichyHBM/teleport-challenge` [→](https://github.com/RichyHBM/teleport-challenge/pull/4#discussion_r2064079323)

> This is not a secure way to do authorization. Any user can impersonate any other user without any verifiable proof of their identity.

**@rosstimothy** on `RichyHBM/teleport-challenge` [→](https://github.com/RichyHBM/teleport-challenge/pull/5#discussion_r2064101257)

> Is this something that you need to handle manually?

**@eriktate** on `RichyHBM/teleport-challenge` [→](https://github.com/RichyHBM/teleport-challenge/pull/3#discussion_r2064260896)

> I think you can simplify this a bit and remove the `else` block entirely. Spreading the error handling across the top level `if` and `else` blocks seems a little harder to follow.
> ```suggestion
> 	if _, err := os.Mkdir(cgroupPath); err != nil {
> 		if !errors.Is(err, os.ErrExist) {
> 			return nil, err
> 		}
> 	}
> ```

**@rosstimothy** on `RichyHBM/teleport-challenge` [→](https://github.com/RichyHBM/teleport-challenge/pull/4#discussion_r2066180172)

> That's what was specified and approved in the design document as well. Did you take a look at the [`peer`](https://pkg.go.dev/google.golang.org/grpc/peer) package?

**@timothyb89** on `rsteinkeXJ/teleport_interview` [→](https://github.com/rsteinkeXJ/teleport_interview/pull/1#discussion_r3121474780)

> This is a reasonable enough description of how the process works generally, but how do you plan to fork+exec the child in Go in particular? A quick Go snippet would potentially be helpful.

**@timothyb89** on `rsteinkeXJ/teleport_interview` [→](https://github.com/rsteinkeXJ/teleport_interview/pull/1#discussion_r3121478699)

> Good call on avoiding PID reuse, but can you define how you plan to generate the session IDs?

**@timothyb89** on `rsteinkeXJ/teleport_interview` [→](https://github.com/rsteinkeXJ/teleport_interview/pull/1#discussion_r3121484424)

> A bit nitpicky, but the usual recommendation is to have RPC specific request and response messages to allow for non-breaking proto changes if you ever decided e.g. `KillJob()` should accept different params than `QueryJob()`, for example if you wanted to add a signal type.

**@timothyb89** on `rsteinkeXJ/teleport_interview` [→](https://github.com/rsteinkeXJ/teleport_interview/pull/1#discussion_r3121511016)

> Can you elaborate on how a client cert will be mapped to an identity?

**@tigrato** on `rsteinkeXJ/teleport_interview` [→](https://github.com/rsteinkeXJ/teleport_interview/pull/1#discussion_r3123854189)

> can you include how you will generate the certificates?

**@tigrato** on `rsteinkeXJ/teleport_interview` [→](https://github.com/rsteinkeXJ/teleport_interview/pull/1#discussion_r3124145005)

> Is this safe? What if the job ended between the time you executed the QueryAPI request and the time you kill the process?

**@timothyb89** on `rsteinkeXJ/teleport_interview` [→](https://github.com/rsteinkeXJ/teleport_interview/pull/1#discussion_r3127738191)

> You're right about the particular syscalls that need to be executed, however Go's runtime more or less guarantees that you can't safely fork/exec without reproducing a surprising amount of extremely carefully written stdlib code. Have you examined Go's standard library closely? Maybe there's tools that can accomplish this for you without having to manage the runtime safety yourself?

**@timothyb89** on `rsteinkeXJ/teleport_interview` [→](https://github.com/rsteinkeXJ/teleport_interview/pull/1#discussion_r3127766287)

> If you'd like, you can decide on some explicit hard-coded limits (200mb RAM, 0.5 CPUs, etc) and not plumb through configuration to the client. But you can also keep it configurable if you'd like, that's just more code/plumbing work.

**@timothyb89** on `rsteinkeXJ/teleport_interview` [→](https://github.com/rsteinkeXJ/teleport_interview/pull/1#discussion_r3127770396)

> Right, you can hard code values for the challenge, but *which* values? When you have an incoming client connection, how will you determine which hard-coded value (whatever that value may be) matches the incoming client?

**@timothyb89** on `rsteinkeXJ/teleport_interview` [→](https://github.com/rsteinkeXJ/teleport_interview/pull/1#discussion_r3127820780)

> I think it's okay to assume the processes don't have permission to escape the cgroup and no setuid is needed. That said, processes don't need any special permissions to change their pgid, so a kill signal to the group isn't necessarily sufficient.
>
> That said, given that we're already using cgroups, maybe there's a simpler way to ensure all child processes are reliably killed and cleaned up? Maybe there's a kernel mechanism you could take advantage of?

**@timothyb89** on `rsteinkeXJ/teleport_interview` [→](https://github.com/rsteinkeXJ/teleport_interview/pull/1#discussion_r3127825836)

> That's a reasonable strategy, but maybe more simply - what numbers will you use? Are the IDs random, sequential, encoded UUIDs, etc?

**@timothyb89** on `rsteinkeXJ/teleport_interview` [→](https://github.com/rsteinkeXJ/teleport_interview/pull/1#discussion_r3135276992)

> An initial note, for the challenge you don't need to gracefully kill child processes, immediate SIGKILL is fine.
>
> That said, I think there's still a few issues with the non-friendly strategy - for example, an app can still fork faster than you can iterate through the child PIDs and effectively prevent the cgroup from being killed. 
>
> It might be worth taking a look at the `cgroup.kill` file: https://www.kernel.org/doc/html/latest/admin-guide/cgroup-v2.html

**@rosstimothy** on `rsteinkeXJ/teleport_interview` [→](https://github.com/rsteinkeXJ/teleport_interview/pull/1#discussion_r3138125977)

> Could you elaborate on this a bit more? How will the admin user be encoded into the cert? How will ordinary users be encoded into the cert?

**@rosstimothy** on `rsteinkeXJ/teleport_interview` [→](https://github.com/rsteinkeXJ/teleport_interview/pull/1#discussion_r3138133950)

> > A timeout heap, for destroying completed sessions a reasonable time after they've finished, will be protected by the same lock.
>
> You can move this complexity out of scope.

**@russjones** on `s-gruneberg/jobWorker` [→](https://github.com/s-gruneberg/jobWorker/pull/1#discussion_r2305721219)

> How do you plan to implement stopping a job?

**@russjones** on `s-gruneberg/jobWorker` [→](https://github.com/s-gruneberg/jobWorker/pull/1#discussion_r2305727467)

> My feeling is you don't need all three.

**@russjones** on `s-gruneberg/jobWorker` [→](https://github.com/s-gruneberg/jobWorker/pull/1#discussion_r2305730220)

> Similar to what @zmb3 said, can you explain what your authorization scheme is. Ideally in a few sentences. From this it's hard for me to understand what you plan to build.

**@rosstimothy** on `sabernabil12/teleport-job-worker` [→](https://github.com/sabernabil12/teleport-job-worker/pull/1#discussion_r2164728061)

> Again, feel free to omit to simplify implementation if you'd like.

**@rosstimothy** on `sabernabil12/teleport-job-worker` [→](https://github.com/sabernabil12/teleport-job-worker/pull/1#discussion_r2164736898)

> Does the standard library provide you any means to redirect stdout and stderr to your desired location without pipes and goroutines just to copy the data?

**@rosstimothy** on `sabernabil12/teleport-job-worker` [→](https://github.com/sabernabil12/teleport-job-worker/pull/1#discussion_r2164746557)

> I think this is an extension of authentication, rather than an authorization scheme. We're looking for you to enforce some set of limiting factors based on the identity of the user.

**@rosstimothy** on `sabernabil12/teleport-job-worker` [→](https://github.com/sabernabil12/teleport-job-worker/pull/1#discussion_r2164749161)

> I don't see any mention for how stopping a job will work. Could you please update the design to include this?

**@nklaassen** on `sabernabil12/teleport-job-worker` [→](https://github.com/sabernabil12/teleport-job-worker/pull/1#discussion_r2165148017)

> TLS 1.3 cipher suites are not configurable in Go

**@nklaassen** on `sabernabil12/teleport-job-worker` [→](https://github.com/sabernabil12/teleport-job-worker/pull/1#discussion_r2165148901)

> what do you mean by this?

**@nklaassen** on `sabernabil12/teleport-job-worker` [→](https://github.com/sabernabil12/teleport-job-worker/pull/1#discussion_r2165149596)

> what do you mean by "isolated" here I don't think it's a requirement for the L4 challenge

**@nklaassen** on `sabernabil12/teleport-job-worker` [→](https://github.com/sabernabil12/teleport-job-worker/pull/1#discussion_r2165152374)

> you don't need to support env vars for the challenge

**@nklaassen** on `sabernabil12/teleport-job-worker` [→](https://github.com/sabernabil12/teleport-job-worker/pull/1#discussion_r2165153416)

> you can skip the timestamp to simplify

**@nklaassen** on `sabernabil12/teleport-job-worker` [→](https://github.com/sabernabil12/teleport-job-worker/pull/1#discussion_r2165157226)

> checking the client cert subject in an allowlist is basically "authentication", we're looking for a bit more for your "authorization" scheme

**@nklaassen** on `sabernabil12/teleport-job-worker` [→](https://github.com/sabernabil12/teleport-job-worker/pull/1#discussion_r2165160203)

> I'm a bit curious what you would do here, but really it's not necessary for the challenge, you can just try to execute what the authorized user sends

**@nklaassen** on `sabernabil12/teleport-job-worker` [→](https://github.com/sabernabil12/teleport-job-worker/pull/1#discussion_r2165161783)

> I guess the inclusion of `crypto/rsa` implies you will use RSA keys for everything? Could you mention the key size and signature algorithm?

**@nklaassen** on `sabernabil12/teleport-job-worker` [→](https://github.com/sabernabil12/teleport-job-worker/pull/1#discussion_r2165164435)

> i mentioned this in another comment but crypto/tls doesn't actually let you configure the cipher suite with TLS 1.3

**@smallinsky** on `samschurter/teleport-challenge` [→](https://github.com/samschurter/teleport-challenge/pull/1#discussion_r692317869)

> For simplicity you can remove `config` subcommand and hardcode cert path in code.

**@smallinsky** on `samschurter/teleport-challenge` [→](https://github.com/samschurter/teleport-challenge/pull/1#discussion_r692330636)

> Could you elaborate based on what the server will validate client certificate ?

**@bernardjkim** on `samschurter/teleport-challenge` [→](https://github.com/samschurter/teleport-challenge/pull/1#discussion_r692549720)

> Have you considered any security concerns of using an incrementing counter?

**@russjones** on `samschurter/teleport-challenge` [→](https://github.com/samschurter/teleport-challenge/pull/1#discussion_r693275278)

> I think UUIDv4 is fine, it's just a randomly generated large number anyway. Then you don't have to worry about leaking any information.

**@russjones** on `samschurter/teleport-challenge` [→](https://github.com/samschurter/teleport-challenge/pull/1#discussion_r693275780)

> How do you intend to stop a job?

**@r0mant** on `samschurter/teleport-challenge` [→](https://github.com/samschurter/teleport-challenge/pull/1#discussion_r694101560)

> What kind of authorization rules are you planning to implement? I.e. what will the server check the cert's Subject against?

**@r0mant** on `samschurter/teleport-challenge` [→](https://github.com/samschurter/teleport-challenge/pull/1#discussion_r694102002)

> FYI it's ok to use a http mux library so you don't have to implement this parsing yourself.

**@r0mant** on `samschurter/teleport-challenge` [→](https://github.com/samschurter/teleport-challenge/pull/1#discussion_r694102838)

> When you say "credentials will be hardcoded" - what do you mean exactly? You will put CA/cert/key in the source code?

**@strideynet** on `shawon-crosen/teleport-challenge` [→](https://github.com/shawon-crosen/teleport-challenge/pull/1#discussion_r2124424351)

> How will the user's identity be encoded into the certificate?

**@strideynet** on `shawon-crosen/teleport-challenge` [→](https://github.com/shawon-crosen/teleport-challenge/pull/1#discussion_r2124453791)

> Yup - it's in-scope to leverage the mTLS certs for user authentication.

**@rosstimothy** on `shawon-crosen/teleport-challenge` [→](https://github.com/shawon-crosen/teleport-challenge/pull/1#discussion_r2126987040)

> What about authorization? After you've identified a user, how and what will be enforced for the given identity?

**@zmb3** on `shawon-crosen/teleport-challenge` [→](https://github.com/shawon-crosen/teleport-challenge/pull/1#discussion_r2129158193)

> This describes how you will generate the server certificate, but that's TLS, not mTLS.
> 
> What about the client certificate(s)?

**@zmb3** on `shawon-crosen/teleport-challenge` [→](https://github.com/shawon-crosen/teleport-challenge/pull/1#discussion_r2129159609)

> > will be authenticated against the user in the certificate
> 
> Can you elaborate? Where in the certificate is the user stored?

**@zmb3** on `shawon-crosen/teleport-challenge` [→](https://github.com/shawon-crosen/teleport-challenge/pull/1#discussion_r2129166871)

> How are you generating this `shawon-ls-idshhyqm` ID?

**@zmb3** on `shawon-crosen/teleport-challenge` [→](https://github.com/shawon-crosen/teleport-challenge/pull/1#discussion_r2129169057)

> If any authenticated user can perform any action then you're not really implementing any form of authorization.
> 
> I'd like to see some sort of decision after authentication but before performing any action to check whether the action is permitted.

**@zmb3** on `shawon-crosen/teleport-challenge` [→](https://github.com/shawon-crosen/teleport-challenge/pull/1#discussion_r2129171846)

> I'm not sure I follow this. Are you saying that you will create an OS-level user for each API user?

**@smallinsky** on `supby/job-worker` [→](https://github.com/supby/job-worker/pull/1#discussion_r715456935)

> Could you elaborate what extension will be used to encode client's role ? 
> Also does userA will be able to read userB jobs output ?

**@smallinsky** on `supby/job-worker` [→](https://github.com/supby/job-worker/pull/1#discussion_r715458923)

> Could you briefly sketch the CLI interface ?

**@smallinsky** on `supby/job-worker` [→](https://github.com/supby/job-worker/pull/1#discussion_r715461023)

> could you elaborate briefly how stop logic will be implemented ?

**@nklaassen** on `supby/job-worker` [→](https://github.com/supby/job-worker/pull/1#discussion_r715930750)

> Why did you choose this cipher suite and how will you use it? I must admit this is a bit of a trick question because the Go standard library does not allow you to select the cipher suites when using TLS 1.3

**@nklaassen** on `supby/job-worker` [→](https://github.com/supby/job-worker/pull/1#discussion_r717791206)

> Good idea with the timeout but for this POC you can just send SIGKILL, that should definitely kill the process

**@nklaassen** on `supby/job-worker` [→](https://github.com/supby/job-worker/pull/1#discussion_r717792781)

> Thanks, yes choosing TLS 1.3 and letting Go choose the cipher suite is good enough

**@rosstimothy** on `tjper/teleport` [→](https://github.com/tjper/teleport/pull/1#discussion_r821063243)

> Can you elaborate on the process execution model? How exactly do you plan on starting a command and ensuring it runs in the correct cgroup?

**@rosstimothy** on `tjper/teleport` [→](https://github.com/tjper/teleport/pull/1#discussion_r821159687)

> I asked because it wasn't overly obvious if you just missed the arguments or if they were supposed to be passed along in the command string.

**@codingllama** on `tjper/teleport` [→](https://github.com/tjper/teleport/pull/1#discussion_r821694839)

> Please mention the certificate parameters too.

**@codingllama** on `tjper/teleport` [→](https://github.com/tjper/teleport/pull/1#discussion_r821698237)

> How do you pick the devices to limit?

**@codingllama** on `tjper/teleport` [→](https://github.com/tjper/teleport/pull/1#discussion_r821712065)

> Should we push cgroups handling entirely to job.Service? Or do you have particular actions in mind for the gRPC to do?

**@codingllama** on `tjper/teleport` [→](https://github.com/tjper/teleport/pull/1#discussion_r821718069)

> I would advise for separate fields, parsing args can get tricky when you factor in quotes.

**@codingllama** on `tjper/teleport` [→](https://github.com/tjper/teleport/pull/1#discussion_r821719856)

> Please mention the cgroups limits you intend to use.

**@codingllama** on `tjper/teleport` [→](https://github.com/tjper/teleport/pull/1#discussion_r822078636)

> We typically ask for a sleepless solution, a filewatcher would work (as you suggested). Ditto for checking if the job ended (as discussed in the other thread).

**@rosstimothy** on `tjper/teleport` [→](https://github.com/tjper/teleport/pull/2#discussion_r826947222)

> If the mount path of the `Service` was changed to  `t.TempDir()` would all these tests still need to run as root? You would also get the added benefit or not leaving any artifacts around in the hosts cgroups directory in the event there are issues with tests.

**@codingllama** on `tjper/teleport` [→](https://github.com/tjper/teleport/pull/2#discussion_r827444723)

> This is good, but you could cut corners like this for challenge. For example, struct initialization would be alright.
> 
> (No need to change anything, though.)

**@codingllama** on `tjper/teleport` [→](https://github.com/tjper/teleport/pull/2#discussion_r827448056)

> pkg/errors was archived a while back, most of its functionality is now a part of the stdlib. Avoid it if you can.

**@codingllama** on `tjper/teleport` [→](https://github.com/tjper/teleport/pull/2#discussion_r827450203)

> Don't defer the Close if you are writing to the file, close after the write and handle the errors. A Close error could mean you write didn't go through.

**@codingllama** on `tjper/teleport` [→](https://github.com/tjper/teleport/pull/2#discussion_r827453835)

> Ditto for deferred close, same elsewhere if applicable.
> 
> `os.WriteFile` could save you some work too.

**@codingllama** on `tjper/teleport` [→](https://github.com/tjper/teleport/pull/2#discussion_r827459713)

> My 2: pick a logger you like and go with it, no need for the indirection.

**@codingllama** on `tjper/teleport` [→](https://github.com/tjper/teleport/pull/2#discussion_r828087864)

> I like Zap too :)
> 
> The main idea here is that we don't need the configurability for the challenge, stick a real logger wherever you want it and we're good to go.

**@codingllama** on `tjper/teleport` [→](https://github.com/tjper/teleport/pull/2#discussion_r828092442)

> nit: technically, I think you want filepath.Join here. Small practical difference, though.

**@codingllama** on `tjper/teleport` [→](https://github.com/tjper/teleport/pull/3#discussion_r829288742)

> Can't you unmarshal directly from the file?

**@codingllama** on `tjper/teleport` [→](https://github.com/tjper/teleport/pull/3#discussion_r829331038)

> Suggestion: If you embed the certs you can save yourself the paths and file reads. It's fine for testing.

**@zmb3** on `tjper/teleport` [→](https://github.com/tjper/teleport/pull/3#discussion_r830291909)

> ```suggestion
> func NewServerTLSConfig(serverCert, serverKey, caCert string) (*tls.Config, error) {
> ```

**@zmb3** on `tjper/teleport` [→](https://github.com/tjper/teleport/pull/3#discussion_r830292080)

> ```suggestion
> 	b, err := os.ReadFile(caCert)
> ```

**@zmb3** on `tjper/teleport` [→](https://github.com/tjper/teleport/pull/3#discussion_r830296940)

> Doesn't seem like you need a whole separate file for one string constant.

**@zmb3** on `tjper/teleport` [→](https://github.com/tjper/teleport/pull/3#discussion_r830300411)

> You already have a mutex, I would probably just use a regular map here for the added type safety.

**@greedy52** on `xSudoNymx/job-manager` [→](https://github.com/xSudoNymx/job-manager/pull/1#discussion_r2704936934)

> > avoid leaking processes
>
> FYI "avoid leaking processes" is not part of L4 requirement. You are welcome to discuss it but not required to design for it (or implement it).

**@tigrato** on `xSudoNymx/job-manager` [→](https://github.com/xSudoNymx/job-manager/pull/1#discussion_r2705310664)

> There is no need to send a soft kill signal followed by sigkill. 
> You send a sig kill directly without waiting for the job to finish

**@tigrato** on `xSudoNymx/job-manager` [→](https://github.com/xSudoNymx/job-manager/pull/1#discussion_r2705329895)

> Can you include the protobuf spec for the API?

**@greedy52** on `xSudoNymx/job-manager` [→](https://github.com/xSudoNymx/job-manager/pull/1#discussion_r2709732915)

> > each output Append operation triggers a notification to wake any subscribers waiting at EOF.
>
> how do you plan track which subscribers are waiting at EOF and which are not? is there a list you have to update? same concern as before where we want to avoid race when new outputs are being written

**@greedy52** on `xSudoNymx/job-manager` [→](https://github.com/xSudoNymx/job-manager/pull/1#discussion_r2714061428)

> the current design assumes disk read by N readers is not an issue (and it is not). You can make similar assumptions for memory where the server has unlimited memory for this challenge. If so, any advantage using file for reading?

**@tigrato** on `xSudoNymx/job-manager` [→](https://github.com/xSudoNymx/job-manager/pull/1#discussion_r2714464827)

> I also don't believe you need a reader go routine per client. Each client runs in its own goroutine and I believe we don't need another for the purpose of reading the file/memory

**@greedy52** on `xSudoNymx/job-manager` [→](https://github.com/xSudoNymx/job-manager/pull/3#discussion_r2738546166)

> > m.dataDir = ""
> 
> i was thinking more on deleting the directory from the file system. it may not always been in the tmp dir. and even in tmp dir, may be good to clean up when no longer needed. not big deal though so leave it to you to decide

**@greedy52** on `xSudoNymx/job-manager` [→](https://github.com/xSudoNymx/job-manager/pull/5#discussion_r2738616576)

> could you share these certificates/keys or share the script generating them?

**@greedy52** on `xSudoNymx/job-manager` [→](https://github.com/xSudoNymx/job-manager/pull/5#discussion_r2742016690)

> i don't think you need `keyEncipherment` but it is not a big deal. rest of the certs look awesome!

**@tigrato** on `xSudoNymx/job-manager` [→](https://github.com/xSudoNymx/job-manager/pull/6#discussion_r2742619339)

> you can use notifyContext for the same purpose

**@sclevine** on `zeeshanhaque21/JobWorkerService` [→](https://github.com/zeeshanhaque21/JobWorkerService/pull/1#discussion_r1686906397)

> Great choice to use the cert for authz, but RBAC with multiple roles is out-of-scope of the challenge requirements.

**@gabrielcorado** on `zeeshanhaque21/JobWorkerService` [→](https://github.com/zeeshanhaque21/JobWorkerService/pull/1#discussion_r1692421004)

> Are you going to calculate the soft/hard limits based on this value?

**@gabrielcorado** on `zeeshanhaque21/JobWorkerService` [→](https://github.com/zeeshanhaque21/JobWorkerService/pull/1#discussion_r1692431294)

> Can you provide a similar flow overview for removing/cleaning up jobs? Could you describe which resources need to be cleared, whether this will be done through context cancellation, etc?

**@sclevine** on `zeeshanhaque21/JobWorkerService` [→](https://github.com/zeeshanhaque21/JobWorkerService/pull/3#discussion_r1702139875)

> Careful with ReadAt. It should always return an error if n < len(p): https://pkg.go.dev/io#ReaderAt
> But this isn't the behavior you'd want for Read

**@espadolini** on `Zephan92/teleport-handson` [→](https://github.com/Zephan92/teleport-handson/pull/1#discussion_r2718742116)

> It's very likely that even very normal things such as shell scripts will be affected by this race, so it's not acceptable to just make a note of it. What ways do you envision to have the process be contained in the cgroup right as it's starting?

**@espadolini** on `Zephan92/teleport-handson` [→](https://github.com/Zephan92/teleport-handson/pull/1#discussion_r2718747648)

> `TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384` is not a TLS 1.3 ciphersuite.

**@espadolini** on `Zephan92/teleport-handson` [→](https://github.com/Zephan92/teleport-handson/pull/1#discussion_r2718778964)

> The zero value for protobuf enums should be left unused and explicitly defined as UNSPECIFIED or UNKNOWN, since a deliberate zero value can't be distinguished from a field that was not known by the peer.

**@espadolini** on `Zephan92/teleport-handson` [→](https://github.com/Zephan92/teleport-handson/pull/1#discussion_r2727155321)

> I don't believe `crypto/tls` will let you customize the ciphersuites in TLS 1.3 mode.

**@espadolini** on `Zephan92/teleport-handson` [→](https://github.com/Zephan92/teleport-handson/pull/1#discussion_r2727156827)

> Where are you going to get the maj:min for the IO limits?

**@nklaassen** on `Zephan92/teleport-handson` [→](https://github.com/Zephan92/teleport-handson/pull/1#discussion_r2729746350)

> So the "Self-Jailing Shim" pattern you describe sounds like it would work, and the ptrace approach you have in the pseudo-code looks like it works, but in my opinion they are two completely different methods.
>
> In this context I would expect a self-jailing shim to be a small wrapper program that adds itself to the cgroup, then execs the target binary. Your code is more of a parent-jailing approach, where the parent process (your server) starts the target binary with ptrace enabled then adds it to the cgroup.
>
> Again, both seem acceptable but I'd expect the example code to match the description. Please update either the description or the example code so they match

**@nklaassen** on `Zephan92/teleport-handson` [→](https://github.com/Zephan92/teleport-handson/pull/1#discussion_r2729759197)

> I'd like to see somewhere in the design which mount(s) you will actually apply the limits to

**@espadolini** on `Zephan92/teleport-handson` [→](https://github.com/Zephan92/teleport-handson/pull/3#discussion_r2736656919)

> pdeathsig considers the end of the thread that spawned the child process, so waiting in a different goroutine like you're doing has the potential of just killing the child process.

**@espadolini** on `Zephan92/teleport-handson` [→](https://github.com/Zephan92/teleport-handson/pull/3#discussion_r2741894033)

> Canceling the context for the command will still potentially kill the pid via `kill` rather than via pidfd depending on which implementation ends up being picked, and if you are skipping this wait you are essentially leaking the `exec.Cmd` without calling `Wait` or `Release`, which will panic when the object is collected.

**@espadolini** on `Zephan92/teleport-handson` [→](https://github.com/Zephan92/teleport-handson/pull/3#discussion_r2741918569)

> nit: there's a minor TOCTOU here, the I/O limit is applied to the device holding the path when we are evaluating the limit but the process is spawned at a later time so it might end up on a different device.

**@espadolini** on `Zephan92/teleport-handson` [→](https://github.com/Zephan92/teleport-handson/pull/3#discussion_r2741950569)

> I would honestly just document this or check if `cgroup.kill` is available on some test cgroup (I don't think it's ever there in the root cgroup unfortunately), it's possible that a kernel has had a backport of the feature and you'd just be excluding it.
>
> At a minimum this should just warn IMO.

**@nklaassen** on `Zephan92/teleport-handson` [→](https://github.com/Zephan92/teleport-handson/pull/3#discussion_r2743491516)

> consider using `t.Context()` everywhere you currently have context.Background(), especially when deferring cancelation

**@nklaassen** on `Zephan92/teleport-handson` [→](https://github.com/Zephan92/teleport-handson/pull/4#discussion_r2748544790)

> you can use WaitGroup.Go now

**@nklaassen** on `Zephan92/teleport-handson` [→](https://github.com/Zephan92/teleport-handson/pull/4#discussion_r2748550585)

> for loops got fixed in 1.22 so you don't really need to do this, you can just capture i

**@espadolini** on `Zephan92/teleport-handson` [→](https://github.com/Zephan92/teleport-handson/pull/5#discussion_r2759695459)

> I don't think this is ever possible, and nil protobuf messages are essentially equivalent to default messages if you don't need to write.

**@espadolini** on `Zephan92/teleport-handson` [→](https://github.com/Zephan92/teleport-handson/pull/6#discussion_r2759959479)

> I think you might want to output this to stdout rather than stderr.

**@nklaassen** on `Zephan92/teleport-handson` [→](https://github.com/Zephan92/teleport-handson/pull/5#discussion_r2761595153)

> imo it would be nicer to just do all these tests with a real TLS connection and real certs. I don't think this case would be possible with RequireAndVerifyClientCert

---

## status-lifecycle

_Distinguishing process states (started/stopped/failed/running); exit codes; synchronous vs async stopping; reaper timing._

**297 quotes** from `28` distinct reviewers across `45` candidate submissions.

**@codingllama** on `adalton/teleport-exercise` [→](https://github.com/adalton/teleport-exercise/pull/1#discussion_r767969572)

> Suggestion: pull this outside, maybe in separate if statements? The semantics change a bit, but I think its alright.
> 
> ```suggestion
> 	if j.cmd.Process != nil {
> 		// capture PID
> 	}
> 	if j.cmd.ProcessState != nil {
> 		// capture SignalNum
> 	}
> ```

**@codingllama** on `adalton/teleport-exercise` [→](https://github.com/adalton/teleport-exercise/pull/1#discussion_r767970603)

> What about the ExitCode?

**@codingllama** on `adalton/teleport-exercise` [→](https://github.com/adalton/teleport-exercise/pull/2#discussion_r771448935)

> Please avoid testing RPCs by client server methods - the recommended way is to use a proper client to exercise it, so you get the "full" server running (with interceptors and whatnot).

**@nklaassen** on `benmoss/job-worker-service` [→](https://github.com/benmoss/job-worker-service/pull/1#discussion_r2669970836)

> could you distinguish between a job that has stopped naturally and one that has been killed?

**@dboslee** on `benmoss/job-worker-service` [→](https://github.com/benmoss/job-worker-service/pull/1#discussion_r2670122933)

> with the value `running: false` it could be unclear whether the job is not yet started or already finished running depending on the implementation. could you clarify the behavior/usage or suggest changes to the API to avoid this ambiguity?

**@nklaassen** on `benmoss/job-worker-service` [→](https://github.com/benmoss/job-worker-service/pull/3#discussion_r2695353565)

> however unlikely it is for cgFD.Close to fail, this branch puts you into a state where the job has been started but never waited on, and Start has returned only an error despite the fact the job is running, this ID is lost, and the job will stay in the "Running" state forever

**@nklaassen** on `benmoss/job-worker-service` [→](https://github.com/benmoss/job-worker-service/pull/3#discussion_r2695406681)

> this is still tripping me up a bit, because adding it to the cleanups slice communicates to me "close this if creating the job fails" but you also want to close it if creating the job succeeds. Seems like it should be unconditionally deferred, the only question is what to do with the error if cgFD.Close fails, which you need to answer for the case where the job is started anyway

**@smallinsky** on `bill-rich/jobworker` [→](https://github.com/bill-rich/jobworker/pull/1#discussion_r1606464916)

> > The status of any job can be queried to return whether it is running,
>       stopped, or completed
>       
>  
> Are the statuses running, stopped, and completed sufficient? If a user's job fails due to an error, it is quite critical for the UI to be able to identify failed jobs. I would like to ask you to keep this requirement in scope for this challenge.

**@smallinsky** on `bill-rich/jobworker` [→](https://github.com/bill-rich/jobworker/pull/1#discussion_r1606495138)

> How a new process will be started ? 
> 
> and when `setupNamespaces` will be called ?

**@smallinsky** on `bill-rich/jobworker` [→](https://github.com/bill-rich/jobworker/pull/1#discussion_r1609434234)

> The `error` is not need but  a user should be able to distinguish if job terminated successful or failed. Do you see any way to implement this ?

**@smallinsky** on `bill-rich/jobworker` [→](https://github.com/bill-rich/jobworker/pull/1#discussion_r1609438461)

> > Clients will be distinguished based on their certs
> 
> How exactly ? 
>  
>  
>  > The simple scheme used here will be a list of which users can run which commands. One option will be any or *.
>  
> In this design user A will be to list all jobs started from user B and get the status and stream output.

**@sclevine** on `bkneis/jobworker` [→](https://github.com/bkneis/jobworker/pull/1#discussion_r1648311841)

> How will this work for long-running processes that need to be stopped gracefully by the client?

**@tcsc** on `bkneis/jobworker` [→](https://github.com/bkneis/jobworker/pull/1#discussion_r1650242738)

> Github supports [mermaid](https://docs.github.com/en/get-started/writing-on-github/working-with-advanced-formatting/creating-diagrams) diagrams natively, if that's easier than rendering images from plantUML.
> 
> The mermaid sequence diagram syntax is pretty similar, but just different enough to trip you up :-)

**@sclevine** on `bkneis/jobworker` [→](https://github.com/bkneis/jobworker/pull/2#discussion_r1657678239)

> Instead of polling, it might be more responsive to observe `runningJob.cmd.Wait()` completing.

**@sclevine** on `bkneis/jobworker` [→](https://github.com/bkneis/jobworker/pull/2#discussion_r1657810477)

> How are child processes handled? What if the command spawns a long-running child process and quits?

**@tcsc** on `bkneis/jobworker` [→](https://github.com/bkneis/jobworker/pull/2#discussion_r1663373312)

> ```suggestion
> // AddProcess mutates the given cmd to instruct Go to add the PID of the started process to a given cgroup
> ```

**@rosstimothy** on `bucknercd/jobworker` [→](https://github.com/bucknercd/jobworker/pull/1#discussion_r2266884565)

> Does this imply that the process is added to the cgroup _after_ it has started? If so, doesn't that provide a window of time for a malicious process to evade limits altogether?

**@rosstimothy** on `bucknercd/jobworker` [→](https://github.com/bucknercd/jobworker/pull/1#discussion_r2269781418)

> I wasn't so much asking how you would handle things when a process terminated. I was trying to figure out how you would handle discovering new output. Per the challenge requirements:
> 
> > Discovering new output should be efficient, avoid busy-waiting or polling.
> 
> If the entire content of the file is consumed, but the process is still running, how will you efficiently discover new output?

**@zmb3** on `bucknercd/jobworker` [→](https://github.com/bucknercd/jobworker/pull/1#discussion_r2271524955)

> This means there will be a non-zero amount of time where the process is running but is not yet in the cgroup.
>
> Any concerns about that approach?

**@rosstimothy** on `bucknercd/jobworker` [→](https://github.com/bucknercd/jobworker/pull/1#discussion_r2276605316)

> Are both Status and GetStatusResponse needed? What is the difference between the two?

**@rosstimothy** on `clydotron/job_worker_service` [→](https://github.com/clydotron/job_worker_service/pull/1#discussion_r1599965469)

> It might be nice for users to differentiate a terminated job and a job manually stopped by a user

**@rosstimothy** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/1#discussion_r1782743020)

> Is it possible to discern a job that terminated on its own from a job that was stopped by a user request?

**@espadolini** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/1#discussion_r1782768867)

> Is there a real benefit to the field presence here?
> 
> Are we going to have the full exit status in the field, or just the exit code? If it's the latter, how do we know if the process has been terminated with a signal?
> 
> Should we keep track of whether or not we stopped (or attempted to stop) the job?

**@espadolini** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/1#discussion_r1782803481)

> (if it's the full numerical exit status, are we going to lean hard on this being a Linux-only project? aiui the WIFEXITED/WIFSIGNALED/WEXISTATUS/WTERMSIG macros are not portable)

**@codingllama** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/1#discussion_r1783298507)

> Can we distinguish "failed to start" from "exited with error"?

**@espadolini** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/4#discussion_r1789826322)

> Shouldn't this be `codes.Unauthenticated`? `codes.Unavailable` signals a condition that can be retried.

**@espadolini** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/5#discussion_r1791455465)

> Doesn't this race with the waitpid in the `exec.Cmd.Wait` running as part of the backgrounded `Job.wait`?

**@rosstimothy** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/5#discussion_r1791850571)

> Style Suggestion: log the error and shuffle code around to get rid of the defer?
> ```suggestion
> 	if err := j.cmd.Wait(); err != nil {
> 		slog.Warn("something happened", slog.Any("error", err))
> 	}
> 	
> 	j.mu.Lock()
> 	defer j.mu.Unlock()
> 	if j.status != pb.JobStatus_JOB_STATUS_STOPPED {
> 		j.status = pb.JobStatus_JOB_STATUS_EXITED
> 	}
> 	j.exitCode = j.cmd.ProcessState.ExitCode()
> 	close(j.waitChan)
> 	_ = j.broadcaster.Close() // Best effort.
> ```

**@codingllama** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/5#discussion_r1792449458)

> ```suggestion
> 	// Job started successfully, store it.
> ```

**@espadolini** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/6#discussion_r1795063136)

> Why is this exported? There's no way to use this safely without being `newJob` or `(*Job).wait` (as you need knowledge about `wait` having been called and having returned to avoid the race on `j.cmd.ProcessState`).

**@codingllama** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/6#discussion_r1795694826)

> Take the name, fd, then close synchronously (and log the error)?

**@codingllama** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/6#discussion_r1795879377)

> Check the status before other assertions? The `ip` command isn't present on all systems.

**@awly** on `diptadas/job-worker` [→](https://github.com/diptadas/job-worker/pull/1#discussion_r646771535)

> this is effectively the same status - the job is still running

**@awly** on `diptadas/job-worker` [→](https://github.com/diptadas/job-worker/pull/1#discussion_r646772858)

> If `Kill()` fails, return an error in the `StopJob` response but keep the status as `Running`

**@Joerger** on `diptadas/job-worker` [→](https://github.com/diptadas/job-worker/pull/1#discussion_r646820566)

> nit: This is the same as `j.status=Exited && j.error != nil`

**@Joerger** on `diptadas/job-worker` [→](https://github.com/diptadas/job-worker/pull/1#discussion_r646829064)

> Similarly to awly's comment below on `StopJob`, startup errors should be sent in the `StartJob` response. Otherwise users may need to check the status after each Start request to make sure it is running. This change may make the `Created` status unnecessary.

**@awly** on `diptadas/job-worker` [→](https://github.com/diptadas/job-worker/pull/2#discussion_r648574302)

> avoid using global state
> create a wrapper type instead so you can have independent job sets for tests

**@zmb3** on `DiscoRiver/jobber` [→](https://github.com/DiscoRiver/jobber/pull/1#discussion_r1989351782)

> Is this necessary if the command was started inside the cgroup?

**@rosstimothy** on `dtom90/teleport-backend-challenge` [→](https://github.com/dtom90/teleport-backend-challenge/pull/1#discussion_r2052992869)

> If these two things happen _after_ the process is started in the previous step is there a chance that output may be lost and that cgroup limits may be evaded?

**@russjones** on `dtom90/teleport-backend-challenge` [→](https://github.com/dtom90/teleport-backend-challenge/pull/1#discussion_r2053135940)

> How will you stop a running job?

**@smallinsky** on `dtom90/teleport-backend-challenge` [→](https://github.com/dtom90/teleport-backend-challenge/pull/1#discussion_r2055368356)

> Can a user distinguish between a successfully completed job and a job that was stopped via the Stop gRPC API?

**@greedy52** on `ehsaniara/joblet` [→](https://github.com/ehsaniara/joblet/pull/2#discussion_r2113912350)

> one concern with the flow is the process is started before adding to the cgroup. this may allow the process to bypass the limits briefly.

**@russjones** on `flychicken123/Job_Worker` [→](https://github.com/flychicken123/Job_Worker/pull/1#discussion_r2065386917)

> Are these the only two states a process could be in?

**@GavinFrazar** on `flychicken123/Job_Worker` [→](https://github.com/flychicken123/Job_Worker/pull/1#discussion_r2067123584)

> I like the mermaid diagram but can you also include some CLI UX example of running the server and the client commands?

**@rosstimothy** on `flychicken123/Job_Worker` [→](https://github.com/flychicken123/Job_Worker/pull/1#discussion_r2069332424)

> It would be nice to be able to differentiate a job that was stopped by a user from a job that ended on it's own. I don't think you need the other states you listed for this challenge though.

**@rosstimothy** on `GevorgGal/jobworker` [→](https://github.com/GevorgGal/jobworker/pull/1#discussion_r2975029563)

> I'm curious, how come you opted for a standalone boolean flag to represent this instead of adding a JOB_STATUS_STOPPED state to the JobStatus enum?

**@rosstimothy** on `GevorgGal/jobworker` [→](https://github.com/GevorgGal/jobworker/pull/1#discussion_r2975072089)

> Is signaling the process group safe? Does signaling the process group provide any guarantees that child processes are terminated?

**@eriktate** on `GevorgGal/jobworker` [→](https://github.com/GevorgGal/jobworker/pull/1#discussion_r2975193809)

> nit: consider removing the `Job started:` prefix. It makes it easier to pipe `jobctl start` into other commands, such as:
> `jobctl start -- ls -la /tmp | jobctl logs`

**@Tener** on `GevorgGal/jobworker` [→](https://github.com/GevorgGal/jobworker/pull/1#discussion_r2976937357)

> What if the process exits before it receives SIGKILL, but after the handler has written to the `stopped` field?

**@rosstimothy** on `GevorgGal/jobworker` [→](https://github.com/GevorgGal/jobworker/pull/2#discussion_r2988929524)

> Is there a chance for a race here? Should we check if stopping is already set too?
>
> ```suggestion
> 	if j.stopping || j.status != JobStatusRunning {
> ```

**@rosstimothy** on `GevorgGal/jobworker` [→](https://github.com/GevorgGal/jobworker/pull/2#discussion_r2989010499)

> Can we prevent this race? Won't this end up erroneously classifying the final job status?

**@Tener** on `GevorgGal/jobworker` [→](https://github.com/GevorgGal/jobworker/pull/2#discussion_r2993090046)

> > There's no way to atomically check liveness and kill as one OS operation.
> 
> This is on OS when delivering signal. We can also check if the process has exited due to signal or not.

**@Tener** on `GevorgGal/jobworker` [→](https://github.com/GevorgGal/jobworker/pull/4#discussion_r3020587348)

> A relatively small thing I haven't realized in prior reviews, but `GetStatus` will return on error:
>
> ```go
> if !ok {
> 		return 0, 0, ErrJobNotFound
> 	}
> ```
>
> Which looks fine, except that first zero is actually `JobStatusRunning`:
>
> ```go
> JobStatusRunning JobStatus = iota
> ```
>
> Here we properly catch that issue by discarding other results when `err` is non-nil; it is a potential for API misuse, however.

**@rosstimothy** on `ghostsquad/prototype-job-worker` [→](https://github.com/ghostsquad/prototype-job-worker/pull/1#discussion_r1531949303)

> The term job in the challenge is referring to an arbitrary Linux process and not a job like object similar to a Kubernetes job. There shouldn't need to be any named files required to start the process. 
>
> Why couldn't starting a job that executes a long running script be as simple as `pjw start /usr/bin/long_running_script.sh`?

**@tigrato** on `gstelang/job-worker-service` [→](https://github.com/gstelang/job-worker-service/pull/1#discussion_r1732907618)

> Could you return the fields separately so that the client can choose the appropriate display method? 
> Additionally, please include a way to differentiate whether the job was stopped manually or ended naturally, along with the exit code and other relevant details.

**@awly** on `hashsequence/Linux-Job-Worker` [→](https://github.com/hashsequence/Linux-Job-Worker/pull/2#discussion_r478687461)

> Can you make the client stateless?
> It shouldn't have to persist a queue of commands, the list of jobs can be fetched from the server.

**@awly** on `hashsequence/Linux-Job-Worker` [→](https://github.com/hashsequence/Linux-Job-Worker/pull/2#discussion_r478691779)

> ```suggestion
>     //if the command is running it will be 0, otherwise if command successfully completes it will be 1
> ```

**@awly** on `hashsequence/Linux-Job-Worker` [→](https://github.com/hashsequence/Linux-Job-Worker/pull/2#discussion_r478692231)

> also include the exit code of the process

**@awly** on `hashsequence/Linux-Job-Worker` [→](https://github.com/hashsequence/Linux-Job-Worker/pull/2#discussion_r478696609)

> what's a `variableTypeRequest` in gRPC terms? is it `oneof`?
> if so, please use separate RPCs for start/stop/querypid/queryrunningprocesses instead

**@awly** on `hashsequence/Linux-Job-Worker` [→](https://github.com/hashsequence/Linux-Job-Worker/pull/2#discussion_r478698850)

> use https://golang.org/pkg/os/#Process.Signal instead (`os.Process` is provided by `os/exec.Command`)

**@awly** on `hashsequence/Linux-Job-Worker` [→](https://github.com/hashsequence/Linux-Job-Worker/pull/4#discussion_r480272262)

> > if it fails to execute then a log called FAILED-\<endtimestamp\>.log will be created to indcate that the job failed to execute
> 
> Is this log necessary?
> The server can keep track of the status in its process table

**@awly** on `hashsequence/Linux-Job-Worker` [→](https://github.com/hashsequence/Linux-Job-Worker/pull/4#discussion_r480274245)

> there should be some kind of `status` field to mark failure directly

**@awly** on `hashsequence/Linux-Job-Worker` [→](https://github.com/hashsequence/Linux-Job-Worker/pull/4#discussion_r481305748)

> ```suggestion
> * The User should be able to stop the request based on the uuid
> 
> * When stopped the process should be killed 
> ```

**@awly** on `hashsequence/Linux-Job-Worker` [→](https://github.com/hashsequence/Linux-Job-Worker/pull/4#discussion_r481306339)

> what are the possible status codes?

**@jimbishopp** on `ilyazz/jobs` [→](https://github.com/ilyazz/jobs/pull/2#discussion_r1001908596)

> Will inspect include detail to determine whether the job exited on its own or was stopped?

**@sclevine** on `ilyazz/jobs` [→](https://github.com/ilyazz/jobs/pull/2#discussion_r1002848411)

> The shim process is replaced by the target process via `syscall.Exec`, so it will be running with PID 1 in the PID namespace. Will this interfere with your plans for termination signal handling?

**@jimbishopp** on `ilyazz/jobs` [→](https://github.com/ilyazz/jobs/pull/3#discussion_r1006290630)

> In looking through the different handlers, it seems like there isn't very much variation between the different states to warrant the complexity of this abstraction. Wouldn't it be easier to reason about if the state checks were included in the exported Job methods?

**@sclevine** on `ilyazz/jobs` [→](https://github.com/ilyazz/jobs/pull/3#discussion_r1007328614)

> Missing something?
> ```
> Running {id [id]} with PID=1 (uid:1000; gid:1000)
> uid=1000(stephen) gid=1000(stephen) groups=1000(stephen),0(root)
> ```

**@sclevine** on `ilyazz/jobs` [→](https://github.com/ilyazz/jobs/pull/3#discussion_r1007541817)

> What happens if child processes are still running (and gracefully shutting down) after `cmd`'s PID exits?

**@sclevine** on `ilyazz/jobs` [→](https://github.com/ilyazz/jobs/pull/6#discussion_r1008985815)

> SIGINT works well when running the server locally, but less well in, e.g., Docker or K8s.

**@AntonAM** on `jaylane/job-scheduler` [→](https://github.com/jaylane/job-scheduler/pull/1#discussion_r1714310837)

> Is it supposed to be using `Status` type instead of `string`?

**@AntonAM** on `jaylane/job-scheduler` [→](https://github.com/jaylane/job-scheduler/pull/1#discussion_r1714311325)

> What would status be if there was an error while starting the job?

**@AntonAM** on `jaylane/job-scheduler` [→](https://github.com/jaylane/job-scheduler/pull/1#discussion_r1714313935)

> Is it supposed to be using Status type instead of string?

**@rosstimothy** on `jaylane/job-scheduler` [→](https://github.com/jaylane/job-scheduler/pull/1#discussion_r1715308367)

> What happens if the machine running the jobs has no matching disk?

**@rosstimothy** on `jaylane/job-scheduler` [→](https://github.com/jaylane/job-scheduler/pull/1#discussion_r1715310564)

> If a job is added by piping its pid into the correct file, doesn't that imply that a job can be running unconstrained prior to that? Wouldn't that open the door for a malicious job to escape the cgroup entirely?

**@tigrato** on `joshuarubin/teleport-job-worker` [→](https://github.com/joshuarubin/teleport-job-worker/pull/3#discussion_r1743552332)

> can you add other status indicating the job was manually stopped?

**@tigrato** on `joshuarubin/teleport-job-worker` [→](https://github.com/joshuarubin/teleport-job-worker/pull/3#discussion_r1743558919)

> can you add another status indicating the job was manually stopped?

**@tigrato** on `joshuarubin/teleport-job-worker` [→](https://github.com/joshuarubin/teleport-job-worker/pull/3#discussion_r1743561398)

> should we inform clients that exit_code is optional by making it a pointer to int32?

**@rosstimothy** on `joshuarubin/teleport-job-worker` [→](https://github.com/joshuarubin/teleport-job-worker/pull/3#discussion_r1743919915)

> > with defaults that do not implement limits
> 
> The challenge requires that all jobs started have resource limits enforced. Perhaps instead of defaulting to no limits the flags can default to some arbitrarily chosen limits?

**@rosstimothy** on `joshuarubin/teleport-job-worker` [→](https://github.com/joshuarubin/teleport-job-worker/pull/3#discussion_r1743922551)

> What will the exit code of a stopped job be? Will the library be able to differentiate between the error and stopped states?

**@rosstimothy** on `joshuarubin/teleport-job-worker` [→](https://github.com/joshuarubin/teleport-job-worker/pull/3#discussion_r1743929440)

> Does this imply that StopJob is synchronous and will block until the job is terminated?

**@rosstimothy** on `joshuarubin/teleport-job-worker` [→](https://github.com/joshuarubin/teleport-job-worker/pull/3#discussion_r1744032448)

> No preference whether stop is synchronous or asynchronous for this challenge, I was just trying to reason about API from the response message. In keeping with the general theme of the challenge though, the simplest thing is often best.

**@codingllama** on `joshuarubin/teleport-job-worker` [→](https://github.com/joshuarubin/teleport-job-worker/pull/3#discussion_r1744132736)

> If the `child` sub-command is running then the binary must call StartJobChild? Is that correct?

**@codingllama** on `joshuarubin/teleport-job-worker` [→](https://github.com/joshuarubin/teleport-job-worker/pull/4#discussion_r1747494567)

> nit:
> 
> ```suggestion
>   JOB_STATUS_COMPLETED = 3; // the job completed successfully on its own
> ```
> 
> so it goes "not started", "completed" and "stopped".

**@rosstimothy** on `joshuarubin/teleport-job-worker` [→](https://github.com/joshuarubin/teleport-job-worker/pull/4#discussion_r1747530246)

> I'm wondering if a single mutex for a small amount of shared state would reduce the complexity of some of the job code below.

**@codingllama** on `joshuarubin/teleport-job-worker` [→](https://github.com/joshuarubin/teleport-job-worker/pull/4#discussion_r1747572848)

> Should an error here cause a state transition? Will the job eternally be considered to be running otherwise?

**@codingllama** on `joshuarubin/teleport-job-worker` [→](https://github.com/joshuarubin/teleport-job-worker/pull/4#discussion_r1747574326)

> Should we capture the exit code here?

**@codingllama** on `joshuarubin/teleport-job-worker` [→](https://github.com/joshuarubin/teleport-job-worker/pull/4#discussion_r1747589273)

> Replicating the previous comment:
> 
> ```suggestion
> 	StatusCompleted          // the job completed successfully on its own
> ```

**@codingllama** on `joshuarubin/teleport-job-worker` [→](https://github.com/joshuarubin/teleport-job-worker/pull/4#discussion_r1747630987)

> > // ExitCode will not exist (...) if it exited uncleanly, e.g. because it was killed by a signal.
> 
> Is this correct?

**@rosstimothy** on `joshuarubin/teleport-job-worker` [→](https://github.com/joshuarubin/teleport-job-worker/pull/4#discussion_r1751047760)

> Should a job in this state be considered NotStarted or Completed? Technically the underlying process was never started, however, the job was attempted to be started but failed and now the job is in a semi-completed state since a job may only be started once.

**@rosstimothy** on `joshuarubin/teleport-job-worker` [→](https://github.com/joshuarubin/teleport-job-worker/pull/4#discussion_r1751049155)

> It might be a good idea to note somewhere that if `j.done` is closed prior to setting the cmdErr or exitCode, then there will be a race when calling Error, ExitCode, and any other helpers that interrogate the job done channel.

**@greedy52** on `kiakeshmiri/process-runner` [→](https://github.com/kiakeshmiri/process-runner/pull/1#discussion_r1765691138)

> why is the status here `string` instead of `Status`?

**@greedy52** on `kiakeshmiri/process-runner` [→](https://github.com/kiakeshmiri/process-runner/pull/1#discussion_r1765691548)

> it would be hard to tell if the process exited on its own or stopped by the user.

**@greedy52** on `kiakeshmiri/process-runner` [→](https://github.com/kiakeshmiri/process-runner/pull/1#discussion_r1765691897)

> any other information could be useful for status?

**@tigrato** on `kiakeshmiri/process-runner` [→](https://github.com/kiakeshmiri/process-runner/pull/1#discussion_r1771339323)

> exited with error != crashed and there is no status to indicate a job has completed

**@tigrato** on `kkloberdanz/teleport-challenge` [→](https://github.com/kkloberdanz/teleport-challenge/pull/6#discussion_r2798153901)

> Exit code 0 is the default (when set to 0 and when unset) and to figure out when set we would need to check the status and infer which status could set the exit code to 0.
>
> Can we make clear when this field is set by making it optional?

**@smallinsky** on `kkloberdanz/teleport-challenge` [→](https://github.com/kkloberdanz/teleport-challenge/pull/7#discussion_r2811933363)

> nit
>
> ```suggestion
> 	return resp.GetStatus(), resp.GetExitCode(), nil
> ```

**@tigrato** on `kkloberdanz/teleport-challenge` [→](https://github.com/kkloberdanz/teleport-challenge/pull/7#discussion_r2811969074)

> Should we allow running the server without any limits?
> Warns are not errors and running the server without resource limits seems to be problematic

**@tigrato** on `kkloberdanz/teleport-challenge` [→](https://github.com/kkloberdanz/teleport-challenge/pull/7#discussion_r2812020869)

> Even if the status is killed, we should set the exist code.

**@smallinsky** on `kkloberdanz/teleport-challenge` [→](https://github.com/kkloberdanz/teleport-challenge/pull/7#discussion_r2812107281)

> It look like the exit code is never set when the job is killed. Could you add proper handing and make sure that exit code matches `128 + <signal_number>` even if the job was killed by  `func (l *localJob) Stop() error {` API ?

**@smallinsky** on `kkloberdanz/teleport-challenge` [→](https://github.com/kkloberdanz/teleport-challenge/pull/7#discussion_r2813585742)

> I think that this this is a bit overcomplicated. 
> 
> Before the flow just returned early by 
> ```
> 	// If Stop already set killed, leave status as killed.
> 	if l.status == StatusKilled {
> 		return
> 	}
> ```
> 
> 
> In the current version, a magic value 128 is introduced. Is there a way to simplify this flow and set the exit code correctly without relying on a hard-coded magic number?
> 
> Related to  https://github.com/kkloberdanz/teleport-challenge/pull/7#discussion_r2813591858

**@smallinsky** on `kkloberdanz/teleport-challenge` [→](https://github.com/kkloberdanz/teleport-challenge/pull/7#discussion_r2813591858)

> > ec := 128 + int(syscall.SIGKILL)
>
> Releted to https://github.com/kkloberdanz/teleport-challenge/pull/7/changes#r2813585742
>
> Do you see any way to simplify this to version without leveraging magic 128 value knowing that cmd.Wait should return correct exit code when it was terminated by a signal ?

**@smallinsky** on `kkloberdanz/teleport-challenge` [→](https://github.com/kkloberdanz/teleport-challenge/pull/7#discussion_r2815511055)

> Thanks for verifying his. , You are absolute  right. in case of Termination by signal the ` exitErr.ExitCode()` returned from go flow is just  -1 so my suggestion is invalid.

**@tigrato** on `kkloberdanz/teleport-challenge` [→](https://github.com/kkloberdanz/teleport-challenge/pull/8#discussion_r2832846855)

> should we move this after we change the process state?

**@zmb3** on `kovyrin/teleport-exec` [→](https://github.com/kovyrin/teleport-exec/pull/1#discussion_r788222426)

> I wouldn't do this. It could allow a user to "hide" their command from an admin by calling `Status` as soon as possible to purge any record of the job running.

**@russjones** on `kovyrin/teleport-exec` [→](https://github.com/kovyrin/teleport-exec/pull/1#discussion_r788284350)

> How will you distinguish stopped by signal or exited on own?

**@russjones** on `kovyrin/teleport-exec` [→](https://github.com/kovyrin/teleport-exec/pull/1#discussion_r788285328)

> Is `Status` for your own debugging? It's not required to implement, but if you are, I would re-use `CommandStateResponse` since this endpoint is missing things like `exit_code`.

**@rosstimothy** on `kovyrin/teleport-exec` [→](https://github.com/kovyrin/teleport-exec/pull/4#discussion_r795071421)

> Did you consider using`signal.NotifyContext` and passing in the timeoutCtx as the parent?

**@zmb3** on `kovyrin/teleport-exec` [→](https://github.com/kovyrin/teleport-exec/pull/2#discussion_r795814960)

> Are you running an older version of Go? This looks to have been fixed.

**@GavinFrazar** on `kurczynski/teleport-job-worker` [→](https://github.com/kurczynski/teleport-job-worker/pull/1#discussion_r1610797956)

> nitpick: singular form makes more sense for an enum, e.g. "running" is a kind of "Status" not a kind of "Statuses"
> ```suggestion
> enum Status{
> ```

**@bernardjkim** on `kurczynski/teleport-job-worker` [→](https://github.com/kurczynski/teleport-job-worker/pull/1#discussion_r1613776825)

> What additional info will this provide over `status`?

**@bernardjkim** on `kurczynski/teleport-job-worker` [→](https://github.com/kurczynski/teleport-job-worker/pull/1#discussion_r1613779982)

> This is a little confusing because `start_time` initially makes me think the `Command` can be scheduled to start at this time.
>  
> Might be better to move this from `Command` to `Info` as `command_start_time`. Or maybe track status transitions with timestamps.

**@sclevine** on `kurczynski/teleport-job-worker` [→](https://github.com/kurczynski/teleport-job-worker/pull/1#discussion_r1617670073)

> I would avoid global state in the library.
> 
> While the API needs to keep a store of jobs, I don't see a similar requirement for the library.

**@sclevine** on `kurczynski/teleport-job-worker` [→](https://github.com/kurczynski/teleport-job-worker/pull/2#discussion_r1621430948)

> Missing cleanup on quit

**@sclevine** on `kurczynski/teleport-job-worker` [→](https://github.com/kurczynski/teleport-job-worker/pull/2#discussion_r1621518950)

> Up to you, as long as it's documented.
> 
> SIGKILL by itself would make the job runner incompatible with long-running processes that need graceful shutdown.
> 
> Also, consider what happens to child processes in either case.

**@bernardjkim** on `kurczynski/teleport-job-worker` [→](https://github.com/kurczynski/teleport-job-worker/pull/2#discussion_r1624888718)

> Similar issue with `Start` and `Stop`. Should there be some checks in place to prevent the job from being started/stopped multiple times?

**@greedy52** on `m3talsmith/jobberthehut` [→](https://github.com/m3talsmith/jobberthehut/pull/1#discussion_r2499170146)

> fyi listing jobs is not required for this challenge. on the other hand, an api to get job status for a single job is required.

**@greedy52** on `m3talsmith/jobberthehut` [→](https://github.com/m3talsmith/jobberthehut/pull/1#discussion_r2499192319)

> any other useful info to represent the state of a job?

**@greedy52** on `m3talsmith/jobberthehut` [→](https://github.com/m3talsmith/jobberthehut/pull/1#discussion_r2499199555)

> where does the pid come from? does it mean the command is already started?

**@greedy52** on `m3talsmith/jobberthehut` [→](https://github.com/m3talsmith/jobberthehut/pull/1#discussion_r2499212903)

> the naming of `NewCertificate` is quite confusing when it has functions like `GetJob`.

**@greedy52** on `m3talsmith/jobberthehut` [→](https://github.com/m3talsmith/jobberthehut/pull/1#discussion_r2500084555)

> yes please. as stated in the challenge:
> > Use a simple authorization scheme.
> 
> I don't think no-authorization counts 😛

**@rosstimothy** on `m3talsmith/jobberthehut` [→](https://github.com/m3talsmith/jobberthehut/pull/1#discussion_r2500111471)

> @m3talsmith as discussed in the kickoff call and stated in the requirements we are looking for you to identify who the caller is, and based on that identity, what they are allowed to do. 
>
>
> https://github.com/gravitational/careers/blob/main/challenges/systems/challenge-1.md?plain=1#L221

**@rosstimothy** on `m3talsmith/jobberthehut` [→](https://github.com/m3talsmith/jobberthehut/pull/1#discussion_r2505254762)

> Which of these states would a job terminated by a stop request fall into? How would a user differentiate a stopped job from a job that ended on its own?

**@rosstimothy** on `m3talsmith/jobberthehut` [→](https://github.com/m3talsmith/jobberthehut/pull/1#discussion_r2505268765)

> >  When a job is started, the job is run in it through systemd-run with those constraints.
> 
> The use of systemd-run violates the [challenge requirements](https://github.com/gravitational/careers/blob/main/challenges/systems/challenge-1.md?plain=1#L107C28-L108C53) around using external binaries.

**@rosstimothy** on `MarkDHarris/teleport-systems-challenge` [→](https://github.com/MarkDHarris/teleport-systems-challenge/pull/1#discussion_r3010025697)

> > Stopped jobs terminate the entire process group (parent and all children) to prevent orphaned processes.
>
>
> This is _not_ a requirement for L4. I suggest moving it out of scope as it adds additional complexity to do correctly.

**@greedy52** on `MarkDHarris/teleport-systems-challenge` [→](https://github.com/MarkDHarris/teleport-systems-challenge/pull/1#discussion_r3016992795)

> just curious what is the concern here to mask the `Start()` error? if executable not found, the caller can figure it out by running some other command like `bash -c 'which some-command'` anyway?

**@rosstimothy** on `MarkDHarris/teleport-systems-challenge` [→](https://github.com/MarkDHarris/teleport-systems-challenge/pull/2#discussion_r3039708562)

> It might be more idiomatic to write this by selecting on a timer and a ticker rather than sleeping and manually checking the deadline.
>
> However, since we control the Job API, is there an alternate approach we could take that doesn't require polling the state?

**@tigrato** on `mcampo84/teleport_challenge` [→](https://github.com/mcampo84/teleport_challenge/pull/5#discussion_r1925381407)

> why having these wrappers?
> unlock does way more than unlock the mutex. 
> What if you only need to read a field like when you do the set/get status, it will wake all readers for something they are not interested in...

**@tigrato** on `mcampo84/teleport_challenge` [→](https://github.com/mcampo84/teleport_challenge/pull/5#discussion_r1925392003)

> why does the job running matters here?

**@tigrato** on `mcampo84/teleport_challenge` [→](https://github.com/mcampo84/teleport_challenge/pull/5#discussion_r1925950393)

> I don't think that would cause a panic at all. The exit code will be non zero but nothing panics

**@tigrato** on `mikewurtz/taskman` [→](https://github.com/mikewurtz/taskman/pull/1#discussion_r2027191424)

> Can you transform this bool into an enum where you can carry more information such as:
> - job started
> - job was stopped
> - job terminated normally
> - job exited with an error...

**@rosstimothy** on `mikewurtz/taskman` [→](https://github.com/mikewurtz/taskman/pull/1#discussion_r2028824042)

> Suggestion: maybe note that a stopped job, which will also have a non-zero exit code will not be marked as `JOB_STATUS_EXITED_ERROR`
> ```suggestion
>     // job exited with a non-zero status and was not stopped
> ```

**@timothyb89** on `mikewurtz/taskman` [→](https://github.com/mikewurtz/taskman/pull/1#discussion_r2028852856)

> To add to @rosstimothy's comment a bit, will this state apply to jobs stopped only by `taskman` via user request, or also to jobs killed via signal? (`kill -9`, OOM kill, etc).

**@rosstimothy** on `mikewurtz/taskman` [→](https://github.com/mikewurtz/taskman/pull/2#discussion_r2034134948)

> Can you think of a way to avoid using a static address? What if the host running the test has another process listening on 50055?

**@rosstimothy** on `mikewurtz/taskman` [→](https://github.com/mikewurtz/taskman/pull/2#discussion_r2034147593)

> Suggestion: go all in on cobra and use [`OutOrStderr`](https://pkg.go.dev/github.com/spf13/cobra#Command.OutOrStderr). Same suggestion applies to the other commands as well.
> ```suggestion
> 				fmt.Fprintf(cmd.OutOrStderr(), "failed to close manager: %v\n", err)
> ```

**@timothyb89** on `mikewurtz/taskman` [→](https://github.com/mikewurtz/taskman/pull/2#discussion_r2034171016)

> Depending where these timeouts are enforced, I'm not sure how realistic some of them are. On the server, starting should be nearly instantaneous to the point I might recommend against a start timeout at all - most of the expensive operations happen asynchronously.
> 
> Similarly, I know the actual implementation is still stubbed out, but I'd expect `GetTaskStatusTimeout` to operate entirely in memory. Could that RPC conceivably take 10 seconds to complete?
> 
> From the client's perspective, network conditions aside, I think all three calls ought to take a more or less equivalent amount of time. What do you think?

**@zmb3** on `mikewurtz/taskman` [→](https://github.com/mikewurtz/taskman/pull/2#discussion_r2036266380)

> ```suggestion
> 		return nil, errors.New("failed to append CA certificate")
> ```
> 
> When you don't have any format verbs and just need an error with a static message you should prefer `errors.New`.

**@tigrato** on `mikewurtz/taskman` [→](https://github.com/mikewurtz/taskman/pull/4#discussion_r2039955061)

> is there any guarantee the task is in termination when `Shutdown` is called or can this hang forever?

**@timothyb89** on `mikewurtz/taskman` [→](https://github.com/mikewurtz/taskman/pull/4#discussion_r2040043725)

> Something of a nit, but I'd recommending hoisting these instructions up to the README since they're a prerequisite for running the app/tests.

**@timothyb89** on `mikewurtz/taskman` [→](https://github.com/mikewurtz/taskman/pull/4#discussion_r2040064685)

> This applies to most of these cgroup controller writes, but I'd suggest adding a method to clean up hanging cgroups if one of these writes fails. For example, the `io` controller wasn't enabled on my box initially, so this write failed. The resulting cgroups don't get cleaned up in that case, and pile up pretty rapidly.

**@tigrato** on `mikewurtz/taskman` [→](https://github.com/mikewurtz/taskman/pull/4#discussion_r2040101804)

> What's the difference between `GetTaskStatus` and `GetTask`? Both receive a taskID and return a `*Task`.
> When should we use one vs the other?

**@tigrato** on `mikewurtz/taskman` [→](https://github.com/mikewurtz/taskman/pull/4#discussion_r2040109080)

> can we use `exec.CommandContext` to cancel the job once the context is terminated?

**@rosstimothy** on `mikewurtz/taskman` [→](https://github.com/mikewurtz/taskman/pull/4#discussion_r2040232584)

> This is more informational than a request to make changes - defer will not be executed when a program is terminated via os.Exit. So in the error case below this will not run. In the grand scheme of things it's likely not important for this challenge, though I thought it might be something that you'd like to know and be aware of.

**@rosstimothy** on `mikewurtz/taskman` [→](https://github.com/mikewurtz/taskman/pull/4#discussion_r2042605996)

> There's a lot of grabbing and releasing the task lock going on here. Is there any way to create the task status response while only grabbing the lock once?

**@oeric** on `moalf/teleport` [→](https://github.com/moalf/teleport/pull/1#discussion_r2594155783)

> It'd probably be easier to run Grafana locally for testing via Docker (just share your compose file or `docker run` command). Running it in EC2 is a bit outside of the scope of this challenge, and we don't want you to incur any AWS expenses.

**@oeric** on `moalf/teleport` [→](https://github.com/moalf/teleport/pull/1#discussion_r2594160666)

> See comment above re: running Grafana in EC2. You can also just share how you'd configure OIDC in Grafana via `docker run` or your compose file with placeholder values, no need to use the Grafana API.

**@russjones** on `moalf/teleport` [→](https://github.com/moalf/teleport/pull/1#discussion_r2594337885)

> If you switch to running OSS Grafana in Docker on your machine like @oeric mentioned you can drop the Grafana Terraform Provider requirement as well.

**@oeric** on `moalf/teleport` [→](https://github.com/moalf/teleport/pull/1#discussion_r2599132307)

> Can you add the specific [Auth0 API scopes](https://auth0.com/docs/get-started/apis/scopes/api-scopes) that your implementation will use to your doc? eg. `create:users`

**@russjones** on `moalf/teleport` [→](https://github.com/moalf/teleport/pull/2#discussion_r2612554365)

> Running `terraform plan` on a Pull Request could lead to a RCE vulnerability: https://alex.kaskaso.li/post/terraform-plan-rce

**@russjones** on `moalf/teleport` [→](https://github.com/moalf/teleport/pull/2#discussion_r2621331125)

> What do you think about not running `terraform plan` at all in the PR and only upon merge to `main`? Sure it's nice if you can see `terraform plan` in the PR, but it's not required.
> 
> What you proposed works, but that's much simpler. What do you think?

**@oeric** on `moalf/teleport` [→](https://github.com/moalf/teleport/pull/2#discussion_r2624999621)

> If we're only running TF plans on merges to main, do we still need `auth0/scripts/download-auth0.sh`?

**@oeric** on `moalf/teleport` [→](https://github.com/moalf/teleport/pull/2#discussion_r2625136334)

> Right, but if we're only running `terraform plan` and `terraform apply` on merges to main, that would mean that any changes would have to be reviewed and approved before anything runs. The RCE calls out concerns with `terraform plan` being ran on untrusted code - if code is reviewed, approved, and merged, I would say that it's sufficiently vetted and trusted.
>
> If you think there's still benefit to having your script pull in the providers, I'm open to leaving it in, just my two cents.

**@russjones** on `moalf/teleport` [→](https://github.com/moalf/teleport/pull/3#discussion_r2628876027)

> I think you can get rid of this block now because you're not running `terraform plan` on PR on line 43 right?

**@nklaassen** on `MrChristianL/Teleport-Job-Worker-Service` [→](https://github.com/MrChristianL/Teleport-Job-Worker-Service/pull/1#discussion_r2879046593)

> we do like to see a solution that is able to distinguish between jobs that _stopped naturally_ and jobs that were _killed by a user of the library_

**@nklaassen** on `MrChristianL/Teleport-Job-Worker-Service` [→](https://github.com/MrChristianL/Teleport-Job-Worker-Service/pull/1#discussion_r2879052816)

> it is okay to skip graceful shutdown and just send SIGKILL

**@eriktate** on `MrChristianL/Teleport-Job-Worker-Service` [→](https://github.com/MrChristianL/Teleport-Job-Worker-Service/pull/1#discussion_r2881141973)

> > Messages for FAILED jobs will include more relevant messages, such as "stopped by user", or "general error", etc.
>
> What error messages do you expect to include? Since stderr is written to disk along with stdout, I assume these would be mostly static messages defined ahead of time?

**@rosstimothy** on `MrChristianL/Teleport-Job-Worker-Service` [→](https://github.com/MrChristianL/Teleport-Job-Worker-Service/pull/1#discussion_r2890754952)

> I don't think the exit code is enough to discern how a process ended. It is entirely possible that a process can exit with a 137 exit code on its own.
> 
> Please update the design to disabiguate.

**@rosstimothy** on `MrChristianL/Teleport-Job-Worker-Service` [→](https://github.com/MrChristianL/Teleport-Job-Worker-Service/pull/1#discussion_r2890762964)

> Should a job stopped by a user be classified as FAILED? The process didn't actually fail, it was terminated intentionally by a human.

**@rosstimothy** on `MrChristianL/Teleport-Job-Worker-Service` [→](https://github.com/MrChristianL/Teleport-Job-Worker-Service/pull/1#discussion_r2890787093)

> Does this mean that stop is synchronous? Meaning it blocks until the process terminates? 
>
> What is the value of success if the job already completed?

**@eriktate** on `MrChristianL/Teleport-Job-Worker-Service` [→](https://github.com/MrChristianL/Teleport-Job-Worker-Service/pull/1#discussion_r2892590659)

> Is the `userStopped` bool necessary? Wouldn't `status == STOPPED` already capture that state?

**@eriktate** on `MrChristianL/Teleport-Job-Worker-Service` [→](https://github.com/MrChristianL/Teleport-Job-Worker-Service/pull/1#discussion_r2897480367)

> The design mentions that `Stop()` is synchronous, so it seems like you could allow `waitForFinish` to handle normal bookkeeping and then `Stop()` could override the relevant status fields before returning. That would prevent having to manage redundant state and help reduce some complexity with synchronizing changes. For example, your current design appears to enforce an ordering that won't ever mark a job as `STOPPED`:
> ```
> 1. `Stop()` sends SIGKILL via `process.Kill()`
> 2. `Stop()` blocks on `<-j.done`
> 3. `Stop()` marks userStopped as true
> 4. `watchForFinish()` (running in goroutine) is blocked on `cmd.Wait()`
> 5. Process exits from signal
> 6. `cmd.Wait()` returns, `watchForFinish()` updates state (since `userStopped` is true, state is marked as `STOPPED`) and closes `j.done`
> 7. `Stop()` unblocks and returns
> ```
> Because `Stop()` blocks on `j.done` before updating `userStopped`, so it will always be false when `watchForFinish` evaluates it.

**@eriktate** on `MrChristianL/Teleport-Job-Worker-Service` [→](https://github.com/MrChristianL/Teleport-Job-Worker-Service/pull/1#discussion_r2898015141)

> > The current design was such that steps 2 and 3 as you described it were flipped, so userStopped was marked as true while Stop() is holding the lock, prior to the wait.
> 
> Not quite. `Stop()` was setting `userStopped` after blocking on `j.done` which isn't closed until `watchForFinish()` completes. I'm assuming these mutations to the `Job` are guarded by the job's mutex, but there was no mention of locking here.
> 
> > Because of this, watchForFinish was still able to mark jobs as stopped, but it was fragile.
> 
> It's not that it was fragile, it's that `userStopped` would _never_ be `true` when `cmd.Wait()` returns. Meaning an assignment of a job's status to `STOPPED` was not possible. 
> 
> That said, the new design looks like it avoids this issue altogether :+1:

**@eriktate** on `MrChristianL/Teleport-Job-Worker-Service` [→](https://github.com/MrChristianL/Teleport-Job-Worker-Service/pull/2#discussion_r2905575085)

> nit: reduce redundancy in error messages
> ```suggestion
> 		return bytesWritten, fmt.Errorf("writing log file: %w", err)
> ```
> We know something went wrong since we're returning an error so the "failed to" doesn't really add anything

**@eriktate** on `MrChristianL/Teleport-Job-Worker-Service` [→](https://github.com/MrChristianL/Teleport-Job-Worker-Service/pull/2#discussion_r2906240846)

> This doesn't hurt anything, but isn't this unnecessary considering `job.Stop()` is synchronous?

**@nklaassen** on `MrChristianL/Teleport-Job-Worker-Service` [→](https://github.com/MrChristianL/Teleport-Job-Worker-Service/pull/2#discussion_r2906824137)

> Stop() won't update the status until watchForFinish closes j.done though

**@Joerger** on `MrChristianL/Teleport-Job-Worker-Service` [→](https://github.com/MrChristianL/Teleport-Job-Worker-Service/pull/2#discussion_r2907474373)

> Since we failed to start the job, isn't the job ID irrelevant? Especially if this is a user-facing error.

**@Joerger** on `MrChristianL/Teleport-Job-Worker-Service` [→](https://github.com/MrChristianL/Teleport-Job-Worker-Service/pull/2#discussion_r2907485228)

> suggestion: status already equals `Running (0)` since it's the default.
> ```suggestion
> ```

**@Joerger** on `MrChristianL/Teleport-Job-Worker-Service` [→](https://github.com/MrChristianL/Teleport-Job-Worker-Service/pull/2#discussion_r2907506324)

> Might be worth an error if the job is in any non-running state so that the user knows what happened without having to check the status afterwards.

**@Joerger** on `MrChristianL/Teleport-Job-Worker-Service` [→](https://github.com/MrChristianL/Teleport-Job-Worker-Service/pull/2#discussion_r2907555056)

> Is it possible that a client retrieving the status at this point in execution would incorrectly see "Failed"?

**@Joerger** on `MrChristianL/Teleport-Job-Worker-Service` [→](https://github.com/MrChristianL/Teleport-Job-Worker-Service/pull/2#discussion_r2915233417)

> Should we still set the exitCode is the job was stopped?

**@rosstimothy** on `MrChristianL/Teleport-Job-Worker-Service` [→](https://github.com/MrChristianL/Teleport-Job-Worker-Service/pull/2#discussion_r2919505486)

> What if the job is attempted to be started in parallel? As @Joerger pointed out above there is a race that we could detect in this test.

**@rosstimothy** on `MrChristianL/Teleport-Job-Worker-Service` [→](https://github.com/MrChristianL/Teleport-Job-Worker-Service/pull/3#discussion_r2924692597)

> Alternate suggestion for you to consider.
> ```suggestion
> 	t.Cleanup(func() {
> 		conn.Close()
> 		s.Stop()
> 	})
>
> 	return pb.NewJobServiceClient(conn)
> ```

**@nklaassen** on `MrChristianL/Teleport-Job-Worker-Service` [→](https://github.com/MrChristianL/Teleport-Job-Worker-Service/pull/2#discussion_r2927422286)

> by checking j.starting it looks like you have avoided concurrent calls to Start. But why would you want to allow any concurrent calls to methods on Job in the middle of Start? E.g. after the lock is released in Start it looks like a concurrent call to Stop could see j.status as running and try to access j.cmd before it has been set

**@nklaassen** on `MrChristianL/Teleport-Job-Worker-Service` [→](https://github.com/MrChristianL/Teleport-Job-Worker-Service/pull/2#discussion_r2927425742)

> yet j.Status remains Running

**@rosstimothy** on `neildo/tjob` [→](https://github.com/neildo/tjob/pull/1#discussion_r1759126564)

> The requirement for the challenge is that I can get information about a single job, not that I can list or get status from all jobs.

**@rosstimothy** on `neildo/tjob` [→](https://github.com/neildo/tjob/pull/1#discussion_r1759130955)

> The challenge requirements are that the API allows retrieving status information about a particular job.

**@rosstimothy** on `neildo/tjob` [→](https://github.com/neildo/tjob/pull/1#discussion_r1759133653)

> Are you guaranteed to get "lines" of output, that is output terminated by `\n` or `\r\n` from any arbitrary linux process?

**@nklaassen** on `neildo/tjob` [→](https://github.com/neildo/tjob/pull/1#discussion_r1759273095)

> ```suggestion
>    Timestamp stopped = 4;
> ```

**@nklaassen** on `neildo/tjob` [→](https://github.com/neildo/tjob/pull/1#discussion_r1759277349)

> nit: linux exit status codes are 8 bits, int64 seems larger than necessary here

**@rosstimothy** on `neildo/tjob` [→](https://github.com/neildo/tjob/pull/1#discussion_r1761097017)

> Does this provide any insight as to why a job might have exited? For instance can a user discern if a job stopped on its own or was terminated by a user stopping the job?

**@rosstimothy** on `neildo/tjob` [→](https://github.com/neildo/tjob/pull/1#discussion_r1764109386)

> If using a Timestamp for started, what do you think about using a Duration here instead?

**@smallinsky** on `P-A-R-U-S/Go-Job-Worker-Service` [→](https://github.com/P-A-R-U-S/Go-Job-Worker-Service/pull/1#discussion_r1730893182)

> Could you switch to an enum for the available process statuses?

**@smallinsky** on `P-A-R-U-S/Go-Job-Worker-Service` [→](https://github.com/P-A-R-U-S/Go-Job-Worker-Service/pull/1#discussion_r1730907126)

> How a ongoing job will be stopped ?

**@smallinsky** on `P-A-R-U-S/Go-Job-Worker-Service` [→](https://github.com/P-A-R-U-S/Go-Job-Worker-Service/pull/1#discussion_r1730913207)

> Great that you brought this up. For the purpose of this change, you don’t need to worry about cleanup so for simplify please keep this out of scope.

**@smallinsky** on `P-A-R-U-S/Go-Job-Worker-Service` [→](https://github.com/P-A-R-U-S/Go-Job-Worker-Service/pull/1#discussion_r1732273377)

> ```suggestion
>     STARTED   = 0;
>     RUNNING   = 1;
>     STOPPED   = 2;
>     FAILED    = 3;
>     COMPLETED = 4;
> ```
> 
> https://protobuf.dev/programming-guides/style/#enums

**@smallinsky** on `P-A-R-U-S/Go-Job-Worker-Service` [→](https://github.com/P-A-R-U-S/Go-Job-Worker-Service/pull/2#discussion_r1743749621)

> This could lead to a potential race condition when using the library:
> 
> ```
> status := job.Status()
> if status.State == JOB_STATUS_RUNNING { 
> }
> ```
> 
> Since status.State refers to the same memory address, it could be modified by other functions (e.g., job.status.State = JOB_STATUS_TERMINATED) without proper locking, leading to unpredictable behavior during job handling.

**@smallinsky** on `P-A-R-U-S/Go-Job-Worker-Service` [→](https://github.com/P-A-R-U-S/Go-Job-Worker-Service/pull/2#discussion_r1743757613)

> This error hanging is quite confusing Does internal server error should be propagate to the user  and mixed with job error ?

**@greedy52** on `P-A-R-U-S/Go-Job-Worker-Service` [→](https://github.com/P-A-R-U-S/Go-Job-Worker-Service/pull/2#discussion_r1743778764)

> while waiting, `job.Status()` is blocked.

**@greedy52** on `P-A-R-U-S/Go-Job-Worker-Service` [→](https://github.com/P-A-R-U-S/Go-Job-Worker-Service/pull/2#discussion_r1745610302)

> I think it depends on the purpose of the lock. If the lock is meant to prevent other API calls, then it's doing its job. If the lock is meant to protect internal members like status, it's not doing anything while waiting so might as well just unlock during that. I will leave it to you to judge that.

**@greedy52** on `P-A-R-U-S/Go-Job-Worker-Service` [→](https://github.com/P-A-R-U-S/Go-Job-Worker-Service/pull/2#discussion_r1745619046)

> nit: what happens if `Stop` is called after the cmd is terminated

**@greedy52** on `P-A-R-U-S/Go-Job-Worker-Service` [→](https://github.com/P-A-R-U-S/Go-Job-Worker-Service/pull/2#discussion_r1746140088)

> do you need lock for `job.isTerminated` and `job.isCompleted`?

**@smallinsky** on `P-A-R-U-S/Go-Job-Worker-Service` [→](https://github.com/P-A-R-U-S/Go-Job-Worker-Service/pull/2#discussion_r1746590529)

> >sorry, trying to understand your point about can be scheduled after the function exit
> 
> Sorry for confusion.
> 
> The question is: what lock protects the following `job.exitReason` write [operation](https://github.com/P-A-R-U-S/Go-Job-Worker-Service/pull/2/files/51b5bbd76770e793822084936875536f669a2fb3#diff-35dec03abce473a482ac214d8104767d523b60d5cdf7d1248f7e02e24a8ccd6aR194)?
> 
> ```go
> if err := ns.DeleteCGroup(job.getCGroupName()); err != nil {
> 	log.Printf("error closing cgroup: %s\n", err)
> 	job.exitReason = errors.Join(job.exitReason, fmt.Errorf("error closing cgroup: %w\n", err))
> }
> ```
> 
> The current implementation assumes that the lock acquired by:
> ```go
>   processState, err := job.cmd.Process.Wait()
>   job.mutex.Lock()
>   defer job.mutex.Unlock()
>   ```
> protects the job.exitReason write, but the lock can be released earlier, and the following:
> 
> ```
> job.exitReason = errors.Join(job.exitReason, fmt.Errorf("error closing cgroup: %w\n", err))
> ```
> could be written after the lock is released.
> 
> Imagine following scenarion where the goroutine
> ```
>   defer func() {
>     cleanCGroup <- true
>     unmount <- true
>   }()
>   ```
> 
> is delayed by GO  scheduler and in meantime a the job.Status() is called where exitReason is read under read lock `job.getExitReason(),`  
> 
> ```go
> func Start() {
>   job.mutex.Lock()
>   defer job.mutex.Unlock()
>  ...
>     go func() {
>         select {
>         case <-unmount:
>             if err := ns.UnmountProc(); err != nil {
>                 log.Printf("error unmounting /proc - %s\n", err)
>                 job.exitReason = errors.Join(job.exitReason, fmt.Errorf("error unmounting /proc - %w\n", err))
>             }
>         }
>     }()
>  ...
>     go func() {
>         processState, err := job.cmd.Process.Wait()
>         job.mutex.Lock()
>         defer job.mutex.Unlock()
>         job.isCompleted = true
> 
>         defer func() {
>             sleep(time.Second * 10)
>             cleanCGroup <- true
>             unmount <- true
>         }()
>         
>         if err = job.output.Close(); err != nil {
>             job.exitReason = errors.Join(job.exitReason, fmt.Errorf("error closing output: %w\n", err))
>         }
>         return
>     }()
>  ...
>  return
> }
> 
> ```
> 
> The current flow doesn't quarantine the that `job.exitReason = errors.Join(job.exitReason, fmt.Errorf("error closing output: %w\n", err))` is called under the lock

**@smallinsky** on `P-A-R-U-S/Go-Job-Worker-Service` [→](https://github.com/P-A-R-U-S/Go-Job-Worker-Service/pull/2#discussion_r1748059263)

> Race condition those writes are unprotected by any lock  and it the a user calls job.Status will cause data race.

**@smallinsky** on `P-A-R-U-S/Go-Job-Worker-Service` [→](https://github.com/P-A-R-U-S/Go-Job-Worker-Service/pull/2#discussion_r1748062473)

> > And it is quite tricky to design unit test to check this, but I've artificially slow down go-routines which
> 
> https://github.com/P-A-R-U-S/Go-Job-Worker-Service/pull/2/files#diff-35dec03abce473a482ac214d8104767d523b60d5cdf7d1248f7e02e24a8ccd6aR203
> 
> This seem to be very hack way to test things caused by overcomplicated code. 
> 
> If you need to call 
> ```go
>  if err {
>     cleanCGroup <- true
>      waitingCleanCGroupToCompleted 
> 	  defer waitingCleanCGroupToCompleted()
> 		return fmt.Errorf("error creating cgroup: %w", err)
> }
> ```
> why not just  move cgrup into function and do the cleanup synchronous  ?

**@greedy52** on `P-A-R-U-S/Go-Job-Worker-Service` [→](https://github.com/P-A-R-U-S/Go-Job-Worker-Service/pull/3#discussion_r1752025858)

> it's unclear what `s.rwmutex` is protecting here. since `s.rwmutex` belongs to server, it should protect its immediate members like `s.userJob[jobID]`. Currently there is a race when `s.userJobs` is read above. Job's status should be protected by the job itslef, not by the server.

**@rhammonds-teleport** on `rajansandeep/teleport-job-worker-svc` [→](https://github.com/rajansandeep/job-worker-svc/pull/1#discussion_r3087228632)

> Similar to my feedback above, this signature is so close to an abstraction in the standard library that is ubiquitous in Go code.

**@eriktate** on `rajansandeep/teleport-job-worker-svc` [→](https://github.com/rajansandeep/job-worker-svc/pull/1#discussion_r3088791944)

> Are we guaranteed that child processes are killed when sending `SIGKILL` to the group?
>
> Child process cleanup isn't part of the L4 requirements, so feel free to remove this section if it saves you time

**@eriktate** on `rajansandeep/teleport-job-worker-svc` [→](https://github.com/rajansandeep/job-worker-svc/pull/1#discussion_r3088998640)

> It's not clear to me where this interface is actually used. The `Worker` seems to have a `map[string]*Job` without a definition for `Job` and I see a `JobInfo` struct that seems specifically for reporting status. Can you clarify how all of these work together? And provide an example definition for `Job`?

**@eriktate** on `rajansandeep/teleport-job-worker-svc` [→](https://github.com/rajansandeep/job-worker-svc/pull/1#discussion_r3093517203)

> ```suggestion
> To handle processes, `SysProcAttr.Setpgid = true` would create a new process group. On Stop, send SIGKILL to the entire group using `syscall.Kill(-pgid, SIGKILL)`. It reliably terminates the main process and any descendants that remain in the original process group are signaled on a best-effort basis.
> ```
> nit: this currently reads as both the main process and descendants are reliably terminated

**@eriktate** on `rajansandeep/teleport-job-worker-svc` [→](https://github.com/rajansandeep/job-worker-svc/pull/1#discussion_r3093577964)

> Make sure you include any required setup instructions, such as running `make certs`, in a `README` during your implementation PRs.

**@atburke** on `rajansandeep/teleport-job-worker-svc` [→](https://github.com/rajansandeep/job-worker-svc/pull/2#discussion_r3101825574)

> My understanding is that `ExitCode` will be nil if and only if `State` != `JobStateCompleted`. You could ditch the pointer and only have `ExitCode` be meaningful when the job is completed (just like you have it in your protobuf spec).

**@atburke** on `rajansandeep/teleport-job-worker-svc` [→](https://github.com/rajansandeep/job-worker-svc/pull/2#discussion_r3101931924)

> How about tests for successfully started jobs?

**@eriktate** on `rajansandeep/teleport-job-worker-svc` [→](https://github.com/rajansandeep/job-worker-svc/pull/2#discussion_r3110946494)

> `ErrJobNotRunning` seems inaccurate when when `j.state == JobStateRunning` and `j.stopRequested` is true.

**@eriktate** on `rajansandeep/teleport-job-worker-svc` [→](https://github.com/rajansandeep/job-worker-svc/pull/2#discussion_r3111243219)

> There's a subtle race condition here where a call to `Stop()` can race with a natural exit and return no error while job state is marked as `JobStateCompleted` or `JobStateFailed`. If `syscall.Kill(-pid, syscall.SIGKILL)` executes after the process exits but before the `cmd.Wait()` resolves, the associated pid will be a zombie and result in `Stop()` returning nil. The waiter will then mark the job as either completed or failed, which I would argue is the correct behavior. I would also consider `Stop()` incorrectly returning no error in these cases fairly minor, but it could lead to some flaky tests in the future or inaccurate reporting. How could we go about resolving this?

**@eriktate** on `rajansandeep/teleport-job-worker-svc` [→](https://github.com/rajansandeep/job-worker-svc/pull/2#discussion_r3111317348)

> Could you include a couple of tests that exercise some real job state transitions?

**@creack** on `razzam21/job-worker-service` [→](https://github.com/razzam21/job-worker-service/pull/1#discussion_r2257551047)

> What is the difference between those states?

**@zmb3** on `razzam21/job-worker-service` [→](https://github.com/razzam21/job-worker-service/pull/1#discussion_r2257563085)

> I would rely on gRPC status codes, that's what they're for, after all.

**@awly** on `renatoaguimaraes/golang-job-scheduler` [→](https://github.com/renatoaguimaraes/golang-job-scheduler/pull/2#discussion_r623316097)

> will this goroutine keep running even after the job is finished?

**@awly** on `renatoaguimaraes/golang-job-scheduler` [→](https://github.com/renatoaguimaraes/golang-job-scheduler/pull/2#discussion_r623375242)

> do you need to close this? or does `InotifyRmWatch` do all the cleanup?

**@Joerger** on `renatoaguimaraes/golang-job-scheduler` [→](https://github.com/renatoaguimaraes/golang-job-scheduler/pull/2#discussion_r623529788)

> Aren't these the values for a currently running job, not a stopped job?

**@r0mant** on `renatoaguimaraes/golang-job-scheduler` [→](https://github.com/renatoaguimaraes/golang-job-scheduler/pull/2#discussion_r623538717)

> `return *job.Status, nil`. Also, why not return a pointer to Status?

**@r0mant** on `renatoaguimaraes/golang-job-scheduler` [→](https://github.com/renatoaguimaraes/golang-job-scheduler/pull/2#discussion_r623540920)

> If the job stops on its own at the same time, this will attempt to send a signal to a terminated process.

**@r0mant** on `renatoaguimaraes/golang-job-scheduler` [→](https://github.com/renatoaguimaraes/golang-job-scheduler/pull/2#discussion_r625193987)

> What does this do?
> 
> Previously, I was referring to the fact that you could simply check the job status here to avoid sending a signal to the finished job.

**@awly** on `renatoaguimaraes/golang-job-scheduler` [→](https://github.com/renatoaguimaraes/golang-job-scheduler/pull/2#discussion_r625365265)

> why does a non-zero status code mean that the process is still running?

**@awly** on `renatoaguimaraes/golang-job-scheduler` [→](https://github.com/renatoaguimaraes/golang-job-scheduler/pull/2#discussion_r625375558)

> return a copy of status to avoid data races
> ```suggestion
> func (w *worker) Query(jobID string) (Status, error) {
> 	w.mtx.RLock()
> 	defer w.mtx.RUnlock()
> 	job, err := w.getJob(jobID)
> 	if err != nil {
> 		return nil, err
> 	}
> 	return *job.Status, nil
> }
> ```

**@Joerger** on `renatoaguimaraes/golang-job-scheduler` [→](https://github.com/renatoaguimaraes/golang-job-scheduler/pull/2#discussion_r625413988)

> My point is that if the expected values of a stopped job are `exited = false` and `exitCode = -1`, it's indistinguishable from a running job.

**@Joerger** on `renatoaguimaraes/golang-job-scheduler` [→](https://github.com/renatoaguimaraes/golang-job-scheduler/pull/2#discussion_r625422670)

> Nevermind, I see now that exited isn't set until the process exits, I thought it would be false while the process runs. My bad!

**@dboslee** on `renatoaguimaraes/golang-job-scheduler` [→](https://github.com/renatoaguimaraes/golang-job-scheduler/pull/5#discussion_r625434051)

> What is the behavior here when a process has exited? From the clients perspective I would expect to receive all the output for a job and then for the connection to close without an error.

**@awly** on `renatoaguimaraes/golang-job-scheduler` [→](https://github.com/renatoaguimaraes/golang-job-scheduler/pull/6#discussion_r625937363)

> nit: exit with a non-0 status

**@awly** on `renatoaguimaraes/golang-job-scheduler` [→](https://github.com/renatoaguimaraes/golang-job-scheduler/pull/6#discussion_r625938017)

> also exit with a non-0 status code

**@awly** on `renatoaguimaraes/golang-job-scheduler` [→](https://github.com/renatoaguimaraes/golang-job-scheduler/pull/6#discussion_r625941428)

> a length check is needed somewhere
> right now, running `worker-client query` (without a 3rd argument) will panic

**@awly** on `renatoaguimaraes/golang-job-scheduler` [→](https://github.com/renatoaguimaraes/golang-job-scheduler/pull/6#discussion_r625942220)

> same arguments as above - increase timeout, add length check for args, indent error flow and exit with non-0 status on error

**@awly** on `renatoaguimaraes/golang-job-scheduler` [→](https://github.com/renatoaguimaraes/golang-job-scheduler/pull/6#discussion_r625942393)

> same arguments as above - increase timeout, add length check for args, indent error flow and exit with non-0 status on error

**@awly** on `renatoaguimaraes/golang-job-scheduler` [→](https://github.com/renatoaguimaraes/golang-job-scheduler/pull/5#discussion_r625967211)

> return a gRPC `PermissionDenied` status code

**@rosstimothy** on `rexposadas/teleport` [→](https://github.com/rexposadas/teleport/pull/1#discussion_r1790652395)

> How will a job be stopped?

**@strideynet** on `rexposadas/teleport` [→](https://github.com/rexposadas/teleport/pull/1#discussion_r1793596905)

> The Protobuf style guide states you should include a zero value and that this zero value should:
> 
> > The zero value enum should have the suffix UNSPECIFIED, because a server or application that gets an unexpected enum value will mark the field as unset in the proto instance. The field accessor will then return the default value, which for enum fields is the first enum value. For more information on the unspecified enum value, see [the Proto Best Practices page](https://protobuf.dev/programming-guides/dos-donts#unspecified-enum).

**@strideynet** on `rexposadas/teleport` [→](https://github.com/rexposadas/teleport/pull/2#discussion_r1802780413)

> Again, risk of panic here if GetStatus returns nil and an error.

**@strideynet** on `rexposadas/teleport` [→](https://github.com/rexposadas/teleport/pull/2#discussion_r1802784403)

> I feel like the Job status potentially suffers from the risk of concurrent read/writes. Does this need a mutex over it?

**@greedy52** on `rexposadas/teleport` [→](https://github.com/rexposadas/teleport/pull/2#discussion_r1803066215)

> looks like status is not being set after `job.Cmd.Wait()` fails

**@rosstimothy** on `RichyHBM/teleport-challenge` [→](https://github.com/RichyHBM/teleport-challenge/pull/1#discussion_r2048949137)

> Might there be more states a job could be in than running/not running that a user would be interested to know?

**@rosstimothy** on `RichyHBM/teleport-challenge` [→](https://github.com/RichyHBM/teleport-challenge/pull/1#discussion_r2048953335)

> Same question about additional states as above for the JobStartResponse. Also I suggest use the same terminology as the JobStartResponse. 
> ```suggestion
>   bool isActive = 1;
> ```

**@rosstimothy** on `RichyHBM/teleport-challenge` [→](https://github.com/RichyHBM/teleport-challenge/pull/1#discussion_r2049499549)

> Does that imply that stopping will be a synchronous operation? That's totally fine if that's what you intend, I'm just looking for some clarity here.

**@rosstimothy** on `RichyHBM/teleport-challenge` [→](https://github.com/RichyHBM/teleport-challenge/pull/1#discussion_r2049505075)

> What do you think about using an enum to convey which state of several a job might currently be in instead of two booleans? It might also help to include the exit code and let users infer whether a process ended cleanly in their opinion rather than trying to have the service make that decision.

**@eriktate** on `RichyHBM/teleport-challenge` [→](https://github.com/RichyHBM/teleport-challenge/pull/1#discussion_r2051209982)

> Can you explain this a bit more? Is the "current user" referring to the host user running the server executable? Or the user provided by the client certificate?

**@rosstimothy** on `RichyHBM/teleport-challenge` [→](https://github.com/RichyHBM/teleport-challenge/pull/1#discussion_r2052357263)

> If it simplifies anything for you, consider omitting this and requiring users to call Stop, then GetStatus to determine the same information.

**@tigrato** on `RichyHBM/teleport-challenge` [→](https://github.com/RichyHBM/teleport-challenge/pull/3#discussion_r2060497128)

> ```suggestion
> 	} 
> 
> 		// Listen for sig interrupts and properly cleanup
> 		go func() {
> 			interruptChan := make(chan os.Signal, 1)
> 			signal.Notify(interruptChan, os.Interrupt)
> 			<-interruptChan
> 			cleanup(grpcServer, parentCGroup)
> 		}()
> 	
> ```

**@rosstimothy** on `RichyHBM/teleport-challenge` [→](https://github.com/RichyHBM/teleport-challenge/pull/4#discussion_r2064068509)

> I don't quite follow what you meant in the above explanation. This should not be needed to convey the true identity of the user.

**@rosstimothy** on `RichyHBM/teleport-challenge` [→](https://github.com/RichyHBM/teleport-challenge/pull/5#discussion_r2064112488)

> Is this needed? Shouldn't the `RemoteJobRunner` be able to track the lifecycle of a process directly?

**@eriktate** on `RichyHBM/teleport-challenge` [→](https://github.com/RichyHBM/teleport-challenge/pull/3#discussion_r2064349075)

> This `.Close()` silently errors in the case that the server is interrupted since the cgroup has already been removed. You might consider only closing `parentCGroup` in this defer instead of also closing it in the `cleanup` function

**@eriktate** on `RichyHBM/teleport-challenge` [→](https://github.com/RichyHBM/teleport-challenge/pull/4#discussion_r2064386217)

> I think the biggest benefit is flexibility in expressing the reason auth failed in the future. You could imagine having slightly different handling for when the user isn't present in the auth list at all versus when a specific command isn't in their allow list. I'd also agree that a `bool` return versus an `error` isn't going to matter from a performance perspective in the vast majority of cases.

**@rosstimothy** on `RichyHBM/teleport-challenge` [→](https://github.com/RichyHBM/teleport-challenge/pull/5#discussion_r2064475887)

> I would suggest the approach you mentioned, but holding of to inspect cmd.ProcessState and cmd.Process until the process has exited. As written this is in violation of the challenge requirements that you don't depend on any external binaries.

**@timothyb89** on `rsteinkeXJ/teleport_interview` [→](https://github.com/rsteinkeXJ/teleport_interview/pull/1#discussion_r3121545325)

> Is it possible to tell if a job is still running or not based on the query response?

**@tigrato** on `rsteinkeXJ/teleport_interview` [→](https://github.com/rsteinkeXJ/teleport_interview/pull/1#discussion_r3124109165)

> can you include a state to differentiate processes that naturally stopped and the ones that were killed?

**@tigrato** on `rsteinkeXJ/teleport_interview` [→](https://github.com/rsteinkeXJ/teleport_interview/pull/1#discussion_r3124118235)

> Is it possible for a process to escape the cleanup if they fork themselves into a new process group?

**@timothyb89** on `rsteinkeXJ/teleport_interview` [→](https://github.com/rsteinkeXJ/teleport_interview/pull/1#discussion_r3127733398)

> A bit of a nit, but you did correctly note this earlier - enum zero values should be explicitly `FOO_UNSPECIFIED`.
>
> (And they do recommend using prefixes, e.g. `JOB_STATUS_FOO` since proto3 enums are unfortunately not scoped)

**@timothyb89** on `rsteinkeXJ/teleport_interview` [→](https://github.com/rsteinkeXJ/teleport_interview/pull/1#discussion_r3135159635)

> Ah, sorry, I misremembered - `syscall.ForkExec()` ought to work as well, but the `os/exec` package should also be sufficient for this usecase. The main thing to avoid is actually running `clone3()` or similar directly, which is inherently unsafe. (I always forget which members of the `syscall` package are thin wrappers or are reasonably safe 😅)

**@rosstimothy** on `rsteinkeXJ/teleport_interview` [→](https://github.com/rsteinkeXJ/teleport_interview/pull/1#discussion_r3142144270)

> Can we make that possible via the status RPC instead of the output RPC?

**@rosstimothy** on `rsteinkeXJ/teleport_interview` [→](https://github.com/rsteinkeXJ/teleport_interview/pull/1#discussion_r3154034702)

> Why do consumers of output care about exit code? Can we separate concerns of what the process emitted as output and how the process ended?

**@rosstimothy** on `sabernabil12/teleport-job-worker` [→](https://github.com/sabernabil12/teleport-job-worker/pull/1#discussion_r2164725211)

> Graceful shutdown is not a requirement for the challenge. Feel free to omit if it you'd like.

**@rosstimothy** on `sabernabil12/teleport-job-worker` [→](https://github.com/sabernabil12/teleport-job-worker/pull/1#discussion_r2164747623)

> What possible values will the status here and in other messages have? Is there any way to make this less opaque?

**@nklaassen** on `sabernabil12/teleport-job-worker` [→](https://github.com/sabernabil12/teleport-job-worker/pull/1#discussion_r2165150021)

> graceful shutdown is not required

**@smallinsky** on `samschurter/teleport-challenge` [→](https://github.com/samschurter/teleport-challenge/pull/1#discussion_r692314905)

> these fields are returned by `status` endpoint. Not sure why we need them here. For example jobID is already know by the caller.

**@smallinsky** on `samschurter/teleport-challenge` [→](https://github.com/samschurter/teleport-challenge/pull/1#discussion_r692321606)

> > A worker will be considered stale after it has been stopped for 1 hour and its output will be discarded
> 
> The scope of discarding can be omitted.

**@smallinsky** on `samschurter/teleport-challenge` [→](https://github.com/samschurter/teleport-challenge/pull/1#discussion_r693240108)

> In current design the status endpoint returns subset of `output` result so for the user perspective there is not reason to call the `status` endpoint if a user can just call output and get all the fields retuned by `status` endpoint + additional fields thus I think that `output` endpoint should be only responsible for providing job output.

**@r0mant** on `samschurter/teleport-challenge` [→](https://github.com/samschurter/teleport-challenge/pull/1#discussion_r694100189)

> You can still return bytes if you just have your /output/jobID endpoint return pure bytes instead of json. Then no encoding will be needed. All that other information in the /output endpoint is available to the client as a part of /status call.

**@smallinsky** on `samschurter/teleport-challenge` [→](https://github.com/samschurter/teleport-challenge/pull/2#discussion_r698318902)

> if ExitCode is called on `running` job this will panic.

**@smallinsky** on `samschurter/teleport-challenge` [→](https://github.com/samschurter/teleport-challenge/pull/2#discussion_r698326439)

> `j.cmd.ProcessState` is updated inside `func (c *Cmd) Wait() error ` library internal call   thus `j.mu` mutex doesn't protect the `cmd.ProcessState` for race condition.

**@r0mant** on `samschurter/teleport-challenge` [→](https://github.com/samschurter/teleport-challenge/pull/2#discussion_r698890440)

> How does the user distinguish b/w the job that's stopped on its own vs by an explicit "stop" command?

**@rosstimothy** on `shawon-crosen/teleport-challenge` [→](https://github.com/shawon-crosen/teleport-challenge/pull/1#discussion_r2126977658)

> > We will have a single user on the machine we are running processes on named `processrunner`.
> 
> Is this stating that the server will be launching jobs as the `processrunner` user? If so, feel free to omit this to simplify.

**@smallinsky** on `supby/job-worker` [→](https://github.com/supby/job-worker/pull/1#discussion_r715444420)

> Why server can't just return some gRPC status code different when stop call failed?

**@smallinsky** on `supby/job-worker` [→](https://github.com/supby/job-worker/pull/1#discussion_r715450748)

> Does this statuses cover all cases ? For how stoped job can be distinguished from exited job ?

**@smallinsky** on `supby/job-worker` [→](https://github.com/supby/job-worker/pull/1#discussion_r715452344)

> Does the Start gRPC endpoint will be synchronous ?

**@nklaassen** on `supby/job-worker` [→](https://github.com/supby/job-worker/pull/1#discussion_r715928896)

> What does the `UNKNOWN` status mean?

**@smallinsky** on `supby/job-worker` [→](https://github.com/supby/job-worker/pull/1#discussion_r716410727)

> > From your prospective, what is appropriate status code if job failed to run due to some error in linux process running?
> 
> in that case I would probably pick `codes.Internal`
> 
> In current approach some messages  uses `success` and `message` custom fields but others like `QueryStatus` -> `QueryStatusResponse` don't wonder how the `QueryStatus` enpoint will indicate that jobID is not found ?

**@smallinsky** on `supby/job-worker` [→](https://github.com/supby/job-worker/pull/2#discussion_r720138301)

> race condition. The job.Stop acquires the `j.mtx` to updated the job but `updateJobStatus` doesn't

**@smallinsky** on `supby/job-worker` [→](https://github.com/supby/job-worker/pull/2#discussion_r720145250)

> Query seems to be a bit vague function name for function that returns job.Status.

**@rosstimothy** on `tjper/teleport` [→](https://github.com/tjper/teleport/pull/1#discussion_r821058428)

> What exactly is a human-readable status?

**@rosstimothy** on `tjper/teleport` [→](https://github.com/tjper/teleport/pull/1#discussion_r821065720)

> What are the potential values for status? Perhaps an enum would be better suited here?

**@rosstimothy** on `tjper/teleport` [→](https://github.com/tjper/teleport/pull/1#discussion_r821119889)

> How will you handle reading from the file when the a read returns an EOF error, yet the job is still running and may produce more output in the future?

**@codingllama** on `tjper/teleport` [→](https://github.com/tjper/teleport/pull/2#discussion_r827459115)

> I wouldn't bother with the cleanup. Let's keep the scope small, this is a good chunk of code we don't really need.

**@codingllama** on `tjper/teleport` [→](https://github.com/tjper/teleport/pull/2#discussion_r827462447)

> If we are testing the cleanup (TestCleanup), then why the assertions along the way? Others tests should take care of those.

**@codingllama** on `tjper/teleport` [→](https://github.com/tjper/teleport/pull/2#discussion_r828084447)

> Let's say cgroups cleanup runs on server shutdown, which mostly happens after all jobs are stopped/cancelled. That means a simple removal would work most of the time. If we do leak a directory occasionally I wouldn't cry over it, specially for the challenge. My 2c, though.

**@codingllama** on `tjper/teleport` [→](https://github.com/tjper/teleport/pull/2#discussion_r828086374)

> I'd be OK with the tests running in "regular" folders, I think it's a reasonable simplification. Plus, it makes them much simpler to run.

**@codingllama** on `tjper/teleport` [→](https://github.com/tjper/teleport/pull/3#discussion_r829295565)

> nit: Push the status check into job.stop(), the responsibility better fits there. (Or maybe don't check at all, multiple cancels are fine.)

**@codingllama** on `tjper/teleport` [→](https://github.com/tjper/teleport/pull/3#discussion_r829317012)

> I would ditch the entire validator package for a switch:
> 
> ```suggestion
> 	switch {
> 	case req.Command == nil:
> 		return nil, status.Errorf(...)
> 	case req.Command.Name == "":
> 		return nil, status.Errorf(...)
> 	case req.Limits == nil:
> 		return nil, status.Errorf(...)
> 	}
> ```
> 
> As a small benefit we even get conditions that are not inverted.

**@codingllama** on `tjper/teleport` [→](https://github.com/tjper/teleport/pull/3#discussion_r829332376)

> Does this need an external server to be running?

**@codingllama** on `tjper/teleport` [→](https://github.com/tjper/teleport/pull/3#discussion_r829335262)

> As before, are we relying on an external server running?

**@greedy52** on `xSudoNymx/job-manager` [→](https://github.com/xSudoNymx/job-manager/pull/1#discussion_r2709757306)

> how to tell if the process is stopped or exited naturally?

**@tigrato** on `xSudoNymx/job-manager` [→](https://github.com/xSudoNymx/job-manager/pull/1#discussion_r2714488705)

> What if a client disconnects while parked in wait?
>
> How can we cleanup the subscription?

**@tigrato** on `xSudoNymx/job-manager` [→](https://github.com/xSudoNymx/job-manager/pull/3#discussion_r2728138038)

> should wait err only return if the program has exited?

**@greedy52** on `xSudoNymx/job-manager` [→](https://github.com/xSudoNymx/job-manager/pull/3#discussion_r2733430022)

> should `m.jobs` be cleared after shutdown? should `DataDir` be removed?

**@greedy52** on `xSudoNymx/job-manager` [→](https://github.com/xSudoNymx/job-manager/pull/3#discussion_r2733434441)

> why not `t.Cleanup`?

**@gabrielcorado** on `zeeshanhaque21/JobWorkerService` [→](https://github.com/zeeshanhaque21/JobWorkerService/pull/1#discussion_r1692415737)

> Similar question: Is the library responsible for deciding which retrieving logs method to use? Subscribe (for running jobs) or retrieve all logs (for done jobs). If yes, should we provide a single method for retrieving the logs?

**@gabrielcorado** on `zeeshanhaque21/JobWorkerService` [→](https://github.com/zeeshanhaque21/JobWorkerService/pull/1#discussion_r1692423256)

> Given the resource limit, there can be cases where the process is stopped/terminated due to resource exhaustion. Will this status cover these scenarios, too? If so, can you update the description?

**@sclevine** on `zeeshanhaque21/JobWorkerService` [→](https://github.com/zeeshanhaque21/JobWorkerService/pull/1#discussion_r1696246863)

> How are processes stopped? If a long-running process spawns child processes, how are all of the processes gracefully terminated? What if they don't quit after a period of time?

**@espadolini** on `Zephan92/teleport-handson` [→](https://github.com/Zephan92/teleport-handson/pull/1#discussion_r2718750280)

> Privilege reduction and capabilities are not part of the challenge text, you can assume that the server process is running with whatever privileges it needs.
>
> Was the idea to drop and regain privileges to set processes up? If so, how?

**@espadolini** on `Zephan92/teleport-handson` [→](https://github.com/Zephan92/teleport-handson/pull/1#discussion_r2718764431)

> This CLI doesn't quite seem to fit with the API definition that uses separate fields for each argument to the process.

**@espadolini** on `Zephan92/teleport-handson` [→](https://github.com/Zephan92/teleport-handson/pull/1#discussion_r2718784344)

> When is a job considered to be finished? And what happens in terms of clean up if the job terminates by itself as opposed to being stopped?

**@nklaassen** on `Zephan92/teleport-handson` [→](https://github.com/Zephan92/teleport-handson/pull/1#discussion_r2719092146)

> there's a lot of stuff here that is not the job output. It's kind of expected to only print the job output. There is a status command if the user wants the status

**@nklaassen** on `Zephan92/teleport-handson` [→](https://github.com/Zephan92/teleport-handson/pull/1#discussion_r2719126168)

> i don't think a panic implies that the actual job crashed, putting it into a failed state would be misleading if it's still actually running

**@espadolini** on `Zephan92/teleport-handson` [→](https://github.com/Zephan92/teleport-handson/pull/1#discussion_r2727130803)

> Using ptrace to hold the child process and letting it continue immediately after the execve is maybe a little wasteful considering that in theory we are in control of what is running precisely right up to the execve, but it will do. 👍

**@espadolini** on `Zephan92/teleport-handson` [→](https://github.com/Zephan92/teleport-handson/pull/3#discussion_r2736618122)

> That's not something that is allowed by the spec, is it? You can assume that there is enough memory on the machine for all the state and you must successfully respond for status or output on any job that has existed.

**@espadolini** on `Zephan92/teleport-handson` [→](https://github.com/Zephan92/teleport-handson/pull/3#discussion_r2736663742)

> What happens if the process was terminated by a signal?

**@nklaassen** on `Zephan92/teleport-handson` [→](https://github.com/Zephan92/teleport-handson/pull/3#discussion_r2737663443)

> what if the status is already StatusFailed?

**@espadolini** on `Zephan92/teleport-handson` [→](https://github.com/Zephan92/teleport-handson/pull/3#discussion_r2741883919)

> The WUNTRACED docs say
>
> > Status for traced children which have stopped is provided even if this option is not specified.
>
> What purpose does it serve here?

**@espadolini** on `Zephan92/teleport-handson` [→](https://github.com/Zephan92/teleport-handson/pull/5#discussion_r2759676816)

> I'd prefer using `ConnectionState.VerifiedChains[0][0]` rather than `PeerCertificates` so you're safeguarded against accidental misconfigurations that let any client cert through (like using the wrong `tls.Config.ClientAuth`).

**@nklaassen** on `Zephan92/teleport-handson` [→](https://github.com/Zephan92/teleport-handson/pull/6#discussion_r2761632420)

> kind of a nit: you probably want to check if your client-side context was cancelled, or else if something internal is canceled on the server the client won't see the error. Alternately, you probably want to return an error (or at least a non-zero exit code) when killed by Ctrl+C.

---

## error-handling

_Error wrapping and propagation; structured error types; grpc status codes; avoiding panic; not swallowing errors._

**147 quotes** from `23` distinct reviewers across `37` candidate submissions.

**@codingllama** on `adalton/teleport-exercise` [→](https://github.com/adalton/teleport-exercise/pull/1#discussion_r768044592)

> Use uuid.MustParse instead of swallowing the error.

**@codingllama** on `adalton/teleport-exercise` [→](https://github.com/adalton/teleport-exercise/pull/2#discussion_r771420854)

> Use signal.NotifyContext instead? It's usually a bit tidier as you get a context to pass forward.

**@codingllama** on `adalton/teleport-exercise` [→](https://github.com/adalton/teleport-exercise/pull/2#discussion_r771423383)

> Let's not panic here - could we capture the error and give it back to RunJobmanagerServer callers?

**@codingllama** on `adalton/teleport-exercise` [→](https://github.com/adalton/teleport-exercise/pull/2#discussion_r771423627)

> Ditto not panicking/return the error to the caller.

**@codingllama** on `adalton/teleport-exercise` [→](https://github.com/adalton/teleport-exercise/pull/2#discussion_r771442279)

> Suggestion: instead of full-blown error types you could probably get away with a few "constant" errors and `errors.Is`.

**@codingllama** on `adalton/teleport-exercise` [→](https://github.com/adalton/teleport-exercise/pull/2#discussion_r771445900)

> Suggestion: you could stick this in an interceptor and get "free" error translations for all your RPCs.

**@codingllama** on `adalton/teleport-exercise` [→](https://github.com/adalton/teleport-exercise/pull/2#discussion_r771449602)

> ```suggestion
> 	assert.Error(t, err)
> ```

**@codingllama** on `adalton/teleport-exercise` [→](https://github.com/adalton/teleport-exercise/pull/2#discussion_r771449936)

> Please check the first error too.

**@codingllama** on `adalton/teleport-exercise` [→](https://github.com/adalton/teleport-exercise/pull/2#discussion_r771450348)

> Please assert the error.

**@codingllama** on `adalton/teleport-exercise` [→](https://github.com/adalton/teleport-exercise/pull/2#discussion_r771452386)

> Let's make sure errors like this are actual permission denied errors (and not something else). We could do the same for the other possible canonical codes, but let's cover at least the authn/z ones.

**@codingllama** on `adalton/teleport-exercise` [→](https://github.com/adalton/teleport-exercise/pull/2#discussion_r771452890)

> Please check the errors. Same throughout.

**@codingllama** on `adalton/teleport-exercise` [→](https://github.com/adalton/teleport-exercise/pull/2#discussion_r771457228)

> grpc conns are lazy, iirc, so I wouldn't expect you to have this problem (as long as you dial to a port that is open).

**@codingllama** on `adalton/teleport-exercise` [→](https://github.com/adalton/teleport-exercise/pull/2#discussion_r771461780)

> What is the error we expect here?

**@codingllama** on `adalton/teleport-exercise` [→](https://github.com/adalton/teleport-exercise/pull/2#discussion_r771461981)

> Same as above, what is the expected error?

**@codingllama** on `adalton/teleport-exercise` [→](https://github.com/adalton/teleport-exercise/pull/2#discussion_r771480518)

> Added in other comment, but as long as you dial to a port that is open you should not have server/client startup issues - grpc conns are lazy and tend to do the right thing.

**@codingllama** on `adalton/teleport-exercise` [→](https://github.com/adalton/teleport-exercise/pull/2#discussion_r772452942)

> Suggestion: Use RunE, return errors from runServer instead of panicking?

**@rosstimothy** on `adalton/teleport-exercise` [→](https://github.com/adalton/teleport-exercise/pull/2#discussion_r772505642)

> Pass the command context here and remove the `context.Background()` call in `runServer`
> ```suggestion
> 		runServer(cmd.Context())
> ```

**@codingllama** on `adalton/teleport-exercise` [→](https://github.com/adalton/teleport-exercise/pull/3#discussion_r773215462)

> This should do the right thing:
> 
> ```suggestion
> 			return response.Context().Error()
> ```

**@rosstimothy** on `adalton/teleport-exercise` [→](https://github.com/adalton/teleport-exercise/pull/3#discussion_r773291365)

> Should we handle the error here?

**@rosstimothy** on `adalton/teleport-exercise` [→](https://github.com/adalton/teleport-exercise/pull/3#discussion_r773292781)

> Why not just make this an `io.Closer`?
> ```suggestion
> func (c *Client) Close() error {
> 	return c.conn.Close()
> }
> ```

**@sclevine** on `benmoss/job-worker-service` [→](https://github.com/benmoss/job-worker-service/pull/1#discussion_r2670653149)

> Does this imply that the context can be canceled by the JobWorker after `Logs` returns?

**@sclevine** on `benmoss/job-worker-service` [→](https://github.com/benmoss/job-worker-service/pull/3#discussion_r2684495858)

> *nit*: what is your motivation for wrapping these arguments in a struct?

**@nklaassen** on `benmoss/job-worker-service` [→](https://github.com/benmoss/job-worker-service/pull/3#discussion_r2695323940)

> I wouldn't expect to have to call the cleanup function if createLogFile returns an error, especially if that isn't documented. And in fact, at the call site at worker.Start, the cleanup function currently isn't called when an error is returned

**@nklaassen** on `benmoss/job-worker-service` [→](https://github.com/benmoss/job-worker-service/pull/3#discussion_r2695374970)

> the challenge spec does have "Ensure error handling and error reporting is consistent" under areas of focus, asynchronous errors should be reported in some way. Not necessarily via a global logger

**@tcsc** on `bkneis/jobworker` [→](https://github.com/bkneis/jobworker/pull/2#discussion_r1663376174)

> Do we know what happens to these FDs when the process exits? Or fails to start? I just want to make sure we're not leaking FDs.

**@tcsc** on `bkneis/jobworker` [→](https://github.com/bkneis/jobworker/pull/2#discussion_r1663380990)

> What does `0x200000` mean in this context?

**@espadolini** on `chintamanil/job-worker` [→](https://github.com/chintamanil/job-worker/pull/2#discussion_r2833501326)

> With this definition of `Job`, how is the grpc server supposed to request termination for a job?

**@espadolini** on `chintamanil/job-worker` [→](https://github.com/chintamanil/job-worker/pull/4#discussion_r2839188311)

> A name like `needsCleanup` might be better, or maybe checking the return value for non-nil errors?

**@GavinFrazar** on `clydotron/job_worker_service` [→](https://github.com/clydotron/job_worker_service/pull/1#discussion_r1602558509)

> ```suggestion
>   or an error: (unknown user, invalid job_id, access denied)
> ```

**@espadolini** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/4#discussion_r1789729129)

> This will exit with no error, won't it? `ShowSubcommandHelpAndExit` might be more appropriate.

**@espadolini** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/4#discussion_r1789845622)

> Should this be a hard error? The job ID in the proto is just a string, it would not be a breaking change to use a different job ID format from the point of view of the server, so the client hard erroring out violates forwards compatibility here.

**@codingllama** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/4#discussion_r1790640270)

> nit: I believe dialing is lazy by default (which is for the best!), so this is unlikely to be a connect error.
> 
> - https://pkg.go.dev/google.golang.org/grpc#NewClient

**@codingllama** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/4#discussion_r1790642108)

> Return the grpc client directly? The wrapper adds a bunch of boilerplate and tweaks the interface for little gain in itself.

**@codingllama** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/4#discussion_r1790654809)

> Similar to the comment on apiclient, I find that there's little reason to wrap grpc.Server instances. You could return the gRPC server directly and save yourself the boilerplate.
> 
> It's also generally relevant that 1 grpc.Server can bind multiple gRPC "services" (gRPC nomenclature can be a tad confusing here). This isn't the case for the challenge, so it's just an observation, but a long-lived project would do well to consider that from the start.

**@codingllama** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/4#discussion_r1790716131)

> ```suggestion
> 			// NOTE: The only possible error at the moment is 'not found'. Return PermissionDenied
> ```

**@rosstimothy** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/5#discussion_r1791826766)

> Perhaps out of paranoia the error should be checked here as well? What happens if the middleware changes and this invariant no longer holds?

**@codingllama** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/5#discussion_r1792447207)

> Add some context to the message? Ie, mention that it's a cmd.Wait error?

**@codingllama** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/5#discussion_r1795627939)

> nit: t.Error would probably be nicer to other cleanups:
> 
> ```suggestion
> 			t.Error("Timeout waiting for cleanup.")
> ```

**@codingllama** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/6#discussion_r1795792097)

> nit: Close without defer, handle the error.

**@codingllama** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/6#discussion_r1795800040)

> Why is this an error?

**@espadolini** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/6#discussion_r1796876434)

> I wonder if it's worth to check if the error is EBUSY and only retry if it is.

**@nklaassen** on `diptadas/job-worker` [→](https://github.com/diptadas/job-worker/pull/1#discussion_r645161518)

> What will `Error` be? A string, integer, something else?

**@Joerger** on `diptadas/job-worker` [→](https://github.com/diptadas/job-worker/pull/1#discussion_r645197447)

> Can CreateJob return an error?

**@Joerger** on `diptadas/job-worker` [→](https://github.com/diptadas/job-worker/pull/1#discussion_r646818859)

> ```suggestion
> - `Error error` startup error, termination error, or exit error
> ```

**@awly** on `diptadas/job-worker` [→](https://github.com/diptadas/job-worker/pull/2#discussion_r648578863)

> ```suggestion
> 		t.Fatalf("found error: %v", err)
> ```

**@awly** on `diptadas/job-worker` [→](https://github.com/diptadas/job-worker/pull/2#discussion_r651966211)

> nit: handle the error returned here

**@zmb3** on `DiscoRiver/jobber` [→](https://github.com/DiscoRiver/jobber/pull/1#discussion_r2065158920)

> Fair enough. 
>
> There is of course some information that should not be leaked, but this does not mean you can't provide helpful error messages to users when something goes wrong.

**@smallinsky** on `dtom90/teleport-backend-challenge` [→](https://github.com/dtom90/teleport-backend-challenge/pull/1#discussion_r2055362393)

> Why not use a build in gRPC error message mechanism  ?

**@smallinsky** on `dtom90/teleport-backend-challenge` [→](https://github.com/dtom90/teleport-backend-challenge/pull/1#discussion_r2055363609)

> Why not use build in gRPC error mechanism ?

**@smallinsky** on `dtom90/teleport-backend-challenge` [→](https://github.com/dtom90/teleport-backend-challenge/pull/1#discussion_r2055365200)

> Same here, Can you leverage the build in grpc error mechanism  ?

**@smallinsky** on `dtom90/teleport-backend-challenge` [→](https://github.com/dtom90/teleport-backend-challenge/pull/1#discussion_r2055367046)

> Why is the success flag needed?
> 
> If the gRPC endpoint returns no error, wouldn’t that already imply the stop call succeeded? Could you elaborate on scenarios where success = false but no gRPC error is returned?

**@zmb3** on `ehsaniara/joblet` [→](https://github.com/ehsaniara/joblet/pull/2#discussion_r2114137917)

> What exactly does this mean? What active monitoring and leak prevention measures are you planning on implementing?

**@rosstimothy** on `GevorgGal/jobworker` [→](https://github.com/GevorgGal/jobworker/pull/2#discussion_r2988999659)

> This is something that we should solve instead of moving to a TODO. The gRPC layer is going to have a fixed size of data that it can transmit before running into errors. Can we avoid shifting that complexity outside of the library. As @Tener has mentioned previously, Go does have a canonical way to solve this problem.

**@eriktate** on `GevorgGal/jobworker` [→](https://github.com/GevorgGal/jobworker/pull/4#discussion_r3018549816)

> Will this ever produce a type error?

**@rosstimothy** on `GevorgGal/jobworker` [→](https://github.com/GevorgGal/jobworker/pull/4#discussion_r3018949314)

> What if we get a different error?

**@Tener** on `GevorgGal/jobworker` [→](https://github.com/GevorgGal/jobworker/pull/4#discussion_r3020443823)

> This kind of helper is not that great in practice as it discards the line information; every log entry will appear to originate in that helper, and not in the respective caller. `op` parameter gets us some of that back, but imperfectly.
>
> Instead of passing `op` and `jobID` as strings it would be cleaner to pass custom logger. For uniformity we can pass `err` as well, so we get:
>
> ```go
> logger := s.logger.With("op", op, "job_id", jobID, "error", err)
> ```
>
> Finally, there is a duplication in message being repeated twice: once for log call, second time for `status.Error` call.
>
> A different matter is the tight coupling between errors returned from the worker library and the returned status. We can already see this is fragile: `Start` is not using `mapError` because it has to handle different set of errors.
>
> Finally, making it a standalone function rather than a method of `Server` would improve encapsulation.

**@Tener** on `GevorgGal/jobworker` [→](https://github.com/GevorgGal/jobworker/pull/4#discussion_r3020540843)

> Why "reason" and not "error"?

**@rosstimothy** on `ghostsquad/prototype-job-worker` [→](https://github.com/ghostsquad/prototype-job-worker/pull/1#discussion_r1531022527)

> This looks like an REST API. The challenge specifically states that the server should be using gRPC.
> 
> > [GRPC](https://grpc.io/) API to start/stop/get status/stream output of a running process.
> 
> Could you please include the proto spec for a gRPC service instead?

**@awly** on `hashsequence/Linux-Job-Worker` [→](https://github.com/hashsequence/Linux-Job-Worker/pull/2#discussion_r478699439)

> gRPC can handle pipelining  for you - the same gRPC client can be used for multiple concurrent requests.
> no need to complicate the Execute API

**@awly** on `hashsequence/Linux-Job-Worker` [→](https://github.com/hashsequence/Linux-Job-Worker/pull/4#discussion_r481307013)

> if these represent the status of request processing (success/failure), then they should be returned as gRPC errors instead

**@jimbishopp** on `ilyazz/jobs` [→](https://github.com/ilyazz/jobs/pull/3#discussion_r1006274170)

> Can this be written differently to capture an error on close?

**@tigrato** on `joshuarubin/teleport-job-worker` [→](https://github.com/joshuarubin/teleport-job-worker/pull/4#discussion_r1746710064)

> Is it a good design to have the job depending on the parent context? If the parent context is the RPC context, won't it be canceled as soon as the job start RPC ends?

**@codingllama** on `joshuarubin/teleport-job-worker` [→](https://github.com/joshuarubin/teleport-job-worker/pull/4#discussion_r1747571521)

> I agree - there's not much use for this context. We don't want the job under a deadline and we have a dedicated Stop method.

**@codingllama** on `joshuarubin/teleport-job-worker` [→](https://github.com/joshuarubin/teleport-job-worker/pull/4#discussion_r1747627221)

> nit: The uppercase error will trip various linters.
> 
> - https://go.dev/wiki/CodeReviewComments#error-strings

**@codingllama** on `joshuarubin/teleport-job-worker` [→](https://github.com/joshuarubin/teleport-job-worker/pull/4#discussion_r1747631845)

> Drop the context.Context from the signature? It seems to serve no purpose.
> 
> Same for others.

**@codingllama** on `joshuarubin/teleport-job-worker` [→](https://github.com/joshuarubin/teleport-job-worker/pull/4#discussion_r1747639856)

> Is it worth adding context to this error, like the file we are writing?

**@codingllama** on `joshuarubin/teleport-job-worker` [→](https://github.com/joshuarubin/teleport-job-worker/pull/4#discussion_r1747640230)

> nit: Add context to the error?

**@codingllama** on `joshuarubin/teleport-job-worker` [→](https://github.com/joshuarubin/teleport-job-worker/pull/4#discussion_r1747640290)

> nit: Add context to the error?

**@codingllama** on `joshuarubin/teleport-job-worker` [→](https://github.com/joshuarubin/teleport-job-worker/pull/4#discussion_r1747654554)

> To be safe:
> 
> ```suggestion
> 	if errors.Is(err, io.EOF) || r.isClosed() {
> ```

**@codingllama** on `joshuarubin/teleport-job-worker` [→](https://github.com/joshuarubin/teleport-job-worker/pull/4#discussion_r1747655800)

> Should we error on multiple closes?
> 
> Document the chosen behavior?

**@smallinsky** on `kkloberdanz/teleport-challenge` [→](https://github.com/kkloberdanz/teleport-challenge/pull/6#discussion_r2797927791)

> ```suggestion
> 	 t.Close( grpcServer.Stop() )
> ```

**@smallinsky** on `kkloberdanz/teleport-challenge` [→](https://github.com/kkloberdanz/teleport-challenge/pull/7#discussion_r2811980408)

> nit: add error handling or at logging path.

**@smallinsky** on `kkloberdanz/teleport-challenge` [→](https://github.com/kkloberdanz/teleport-challenge/pull/7#discussion_r2811981233)

> nit: add error handling or at logging path.

**@tigrato** on `kkloberdanz/teleport-challenge` [→](https://github.com/kkloberdanz/teleport-challenge/pull/8#discussion_r2821969799)

> can we satisfy io.ReadClose interface instead of receiving the context.Context. This will allow io.Copy variants

**@tigrato** on `kkloberdanz/teleport-challenge` [→](https://github.com/kkloberdanz/teleport-challenge/pull/8#discussion_r2822003649)

> Should we move IdentityFromContext into a grpc interceptor to ensure we always have them available?

**@tigrato** on `kkloberdanz/teleport-challenge` [→](https://github.com/kkloberdanz/teleport-challenge/pull/8#discussion_r2832822370)

> ```suggestion
> 		if errors.Is(err,io.EOF) {
> ```

**@rosstimothy** on `kovyrin/teleport-exec` [→](https://github.com/kovyrin/teleport-exec/pull/4#discussion_r794920239)

> It is a common convention for the `context.Context` to be the first argument of a function and to return an error instead of panicing
> ```suggestion
> func NewFileStream(ctx context.Context, filename string) (*FileStream, error) {
> ```

**@zmb3** on `kovyrin/teleport-exec` [→](https://github.com/kovyrin/teleport-exec/pull/2#discussion_r795819603)

> `TearDown()` will not be called in case of error. This is because you're using `log.Fatalln` everywhere, which calls `os.Exit`, skipping your deferred calls.

**@zmb3** on `kovyrin/teleport-exec` [→](https://github.com/kovyrin/teleport-exec/pull/2#discussion_r795831104)

> These operations can fail, right? Should you be returning an error?

**@nklaassen** on `kovyrin/teleport-exec` [→](https://github.com/kovyrin/teleport-exec/pull/2#discussion_r796198932)

> you're ignoring the error return here

**@greedy52** on `m3talsmith/jobberthehut` [→](https://github.com/m3talsmith/jobberthehut/pull/1#discussion_r2499156259)

> why `HttpNotAuthorized` when using gRPC?

**@rosstimothy** on `m3talsmith/jobberthehut` [→](https://github.com/m3talsmith/jobberthehut/pull/1#discussion_r2500032861)

> Can we use the built in gRPC error handling mechanism instead of defining an error in band for each response?

**@greedy52** on `m3talsmith/jobberthehut` [→](https://github.com/m3talsmith/jobberthehut/pull/1#discussion_r2500036916)

> code is good. but why `Http` though instead of using gRPC codes?

**@rosstimothy** on `m3talsmith/jobberthehut` [→](https://github.com/m3talsmith/jobberthehut/pull/1#discussion_r2500104455)

> Same feedback here, let's avoid HttpInternalServerError and use just what gRPC provides for signaling errors.

**@greedy52** on `MarkDHarris/teleport-systems-challenge` [→](https://github.com/MarkDHarris/teleport-systems-challenge/pull/1#discussion_r3017013619)

> what happens when `CancelJob` fails?

**@tigrato** on `mcampo84/teleport_challenge` [→](https://github.com/mcampo84/teleport_challenge/pull/1#discussion_r1907297857)

> use gRPC native errors instead of strings. If it's a program execution error, it should be part of the output log

**@tigrato** on `mcampo84/teleport_challenge` [→](https://github.com/mcampo84/teleport_challenge/pull/5#discussion_r1920262799)

> why is recovery necessary?

**@tigrato** on `mcampo84/teleport_challenge` [→](https://github.com/mcampo84/teleport_challenge/pull/5#discussion_r1925382247)

> how can the job execution panic?

**@timothyb89** on `mikewurtz/taskman` [→](https://github.com/mikewurtz/taskman/pull/4#discussion_r2040008334)

> I think one or both of these needs to be `require` rather than `assert` - the assertion is failing on my box, then panicking due to a nil pointer deref on the 2nd assert call.
> 
> (FWIW in practice, I tend to prefer starting with `require` for this reason)

**@timothyb89** on `mikewurtz/taskman` [→](https://github.com/mikewurtz/taskman/pull/4#discussion_r2040015977)

> ```suggestion
> 	require.NotNil(t, resp)
> ```
> 
> same issue here, this will panic if the previous call fails; non exhaustive, there's several assertions like this

**@timothyb89** on `mikewurtz/taskman` [→](https://github.com/mikewurtz/taskman/pull/4#discussion_r2040082834)

> Additional prereq maybe worth noting in the README (or maybe auto detecting) - this will fail if the system does not have a `/dev/sda`.

**@rosstimothy** on `mikewurtz/taskman` [→](https://github.com/mikewurtz/taskman/pull/4#discussion_r2040228830)

> In what situations would the above output write fail and this one will succeed?

**@Joerger** on `MrChristianL/Teleport-Job-Worker-Service` [→](https://github.com/MrChristianL/Teleport-Job-Worker-Service/pull/2#discussion_r2915133170)

> Oh, I actually meant that we should check that the error is `io.ErrClosedPipe` as expected.

**@nklaassen** on `MrChristianL/Teleport-Job-Worker-Service` [→](https://github.com/MrChristianL/Teleport-Job-Worker-Service/pull/2#discussion_r2927387825)

> consider errors.Is in case the error gets wrapped

**@rosstimothy** on `neildo/tjob` [→](https://github.com/neildo/tjob/pull/2#discussion_r1768548327)

> Suggestion: handle the happy case and return early to unindent this code.
> 
> https://go.dev/wiki/CodeReviewComments#indent-error-flow

**@smallinsky** on `P-A-R-U-S/Go-Job-Worker-Service` [→](https://github.com/P-A-R-U-S/Go-Job-Worker-Service/pull/2#discussion_r1745054096)

> Nit: https://go.dev/wiki/CodeReviewComments#receiver-names

**@atburke** on `rajansandeep/teleport-job-worker-svc` [→](https://github.com/rajansandeep/job-worker-svc/pull/5#discussion_r3120521510)

> This will panic if info is nil.

**@zmb3** on `razzam21/job-worker-service` [→](https://github.com/razzam21/job-worker-service/pull/1#discussion_r2255185400)

> Do we need to communicate success/failure in the message like this?
>
> If we do, then why do we do it only in the stop response and not in the start response?

**@Joerger** on `renatoaguimaraes/golang-job-scheduler` [→](https://github.com/renatoaguimaraes/golang-job-scheduler/pull/2#discussion_r623528105)

> Why does this say "Fail to read events"?

**@r0mant** on `renatoaguimaraes/golang-job-scheduler` [→](https://github.com/renatoaguimaraes/golang-job-scheduler/pull/2#discussion_r623538036)

> > Fail to read events
> 
> Which events?

**@r0mant** on `renatoaguimaraes/golang-job-scheduler` [→](https://github.com/renatoaguimaraes/golang-job-scheduler/pull/5#discussion_r625200913)

> Check lengths to avoid panics.

**@awly** on `renatoaguimaraes/golang-job-scheduler` [→](https://github.com/renatoaguimaraes/golang-job-scheduler/pull/2#discussion_r625346700)

> ```suggestion
> 			log.Printf("Fail to close file descriptor: %v", fderr)
> ```

**@awly** on `renatoaguimaraes/golang-job-scheduler` [→](https://github.com/renatoaguimaraes/golang-job-scheduler/pull/6#discussion_r625937715)

> https://github.com/golang/go/wiki/CodeReviewComments#indent-error-flow

**@awly** on `renatoaguimaraes/golang-job-scheduler` [→](https://github.com/renatoaguimaraes/golang-job-scheduler/pull/6#discussion_r625940620)

> https://github.com/golang/go/wiki/CodeReviewComments#indent-error-flow

**@awly** on `renatoaguimaraes/golang-job-scheduler` [→](https://github.com/renatoaguimaraes/golang-job-scheduler/pull/5#discussion_r625960012)

> nit: check the length to avoid a panic

**@strideynet** on `rexposadas/teleport` [→](https://github.com/rexposadas/teleport/pull/2#discussion_r1802779809)

> Will this panic if the call to start returns a nil pointer and an error ?

**@nklaassen** on `sabernabil12/teleport-job-worker` [→](https://github.com/sabernabil12/teleport-job-worker/pull/1#discussion_r2165153131)

> i think success/fail could be determined by returning an error from the RPC rather than including it in the response

**@r0mant** on `samschurter/teleport-challenge` [→](https://github.com/samschurter/teleport-challenge/pull/2#discussion_r698888118)

> If we ignore the error, how would we know the reason for the failure?

**@codingllama** on `tjper/teleport` [→](https://github.com/tjper/teleport/pull/2#discussion_r827461042)

> nit: Write clearer failure messages, specially for methods that are actually under tests (and are not just setup).

**@codingllama** on `tjper/teleport` [→](https://github.com/tjper/teleport/pull/3#discussion_r829331517)

> Strange, I would expect a connection failure error instead of a timeout.

**@greedy52** on `xSudoNymx/job-manager` [→](https://github.com/xSudoNymx/job-manager/pull/4#discussion_r2738579486)

> suggestion: avoid `panic`

**@espadolini** on `Zephan92/teleport-handson` [→](https://github.com/Zephan92/teleport-handson/pull/1#discussion_r2718761440)

> What is liable to panic, exactly?

**@espadolini** on `Zephan92/teleport-handson` [→](https://github.com/Zephan92/teleport-handson/pull/4#discussion_r2745680320)

> I'd prefer returning `io.ErrClosedPipe` or `os.ErrClosed` rather than a stringy error.

**@russjones** on `atburke/teleport_interview` [→](https://github.com/atburke/teleport_interview/pull/1#discussion_r643452892)

> Don't drop errors here (and everywhere else) return the error to the caller and let the caller make the decision.

**@alex-kovoy** on `atburke/teleport_interview` [→](https://github.com/atburke/teleport_interview/pull/1#discussion_r644255061)

> lets add a proper error handling, right now all API errors get swallowed by the `api.tsx` so I suggest to propagate an actual error and display it to the user.

**@alex-kovoy** on `atburke/teleport_interview` [→](https://github.com/atburke/teleport_interview/pull/1#discussion_r644259522)

> either return a promise as is and let the caller to handle it, or you can so something like:
> 
> ```
> return { result, error }
> ```
> 
> so it could be used like so:
> ```ts
> const { result, error } = api.login(xxxx)
> const { error } = api.login(xxxx)
> if (error) {
>  ///
> }
> ```

**@russjones** on `atburke/teleport_interview` [→](https://github.com/atburke/teleport_interview/pull/1#discussion_r645900998)

> Instead of passing in `time.Now()` to each function,. the common Go pattern is to create a fake clock, something like the following:
> 
> ```go
> import (
>    "github.com/jonboulle/clockwork"
> )
> 
> type MySqlDatabase struct {
>    driver *sql.DB
>    clock clockwork.Clock
> }
> 
> func NewMySqlDatabase(username, password, databaseName string, clock clockwork.Clock) (*MySqlDatabase, error) {
>    [...]
>    if clock == nil {
>      clock = clockwork.NewRealClock()
>    }
>    return &MySqlDatabase{
>        [...]
>        clock: clock,
>    }
> }
> ```
> 
> Now the caller does not have to worry about time and you can also fake time in tests.

**@r0mant** on `atburke/teleport_interview` [→](https://github.com/atburke/teleport_interview/pull/1#discussion_r645902899)

> Curious why doesn't this return an error? If it did, it'd make it compatible with `io.Closer`.

**@r0mant** on `atburke/teleport_interview` [→](https://github.com/atburke/teleport_interview/pull/1#discussion_r645903394)

> Can this return an error?

**@russjones** on `atburke/teleport_interview` [→](https://github.com/atburke/teleport_interview/pull/1#discussion_r645904087)

> I don't think you need to split this out, just consolidate both checks and into `crypto.IsCorrectPassword` and return an error if either one fails.

**@alex-kovoy** on `ibeckermayer/teleport-interview` [→](https://github.com/ibeckermayer/teleport-interview/pull/2#discussion_r552384894)

> @ibeckermayer you are right, the server needs to return either requested file (if found) or index.html. Client side will handle the rest by showing `404` error component. This should be OK for this project.

**@alex-kovoy** on `ibeckermayer/teleport-interview` [→](https://github.com/ibeckermayer/teleport-interview/pull/2#discussion_r552385765)

> It's better not to leak any information about server file system via error messages. In this particular case it would be enough to return a generic error message ("Invalid request" for example) and log the actual error message.

**@alex-kovoy** on `ibeckermayer/teleport-interview` [→](https://github.com/ibeckermayer/teleport-interview/pull/2#discussion_r552387324)

> Here we are allowing an attacker to know more about our system by leaking this error text. A generic "Access Denied" would suffice (+ log the actual error of course)

**@alex-kovoy** on `ibeckermayer/teleport-interview` [→](https://github.com/ibeckermayer/teleport-interview/pull/2#discussion_r552390267)

> This error message handling seems redundant, just wonder if the Decode error text is too encryptic that it requires additional handling.
> 
> Hm...would this be enough? 
> 
> ```golang
> if err != nil {
>    return &malformedRequest{status: http.StatusBadRequest, msg: err}
> }
> ```

**@alex-kovoy** on `ibeckermayer/teleport-interview` [→](https://github.com/ibeckermayer/teleport-interview/pull/2#discussion_r552394853)

> kudos for trying to use native APIs (fetch). The only issue with it is that 
> > the Promise returned from fetch() won’t reject on HTTP error status even if the response is an HTTP 404 or 500. Instead, it will resolve normally (with ok status set to false), and it will only reject on network failure or if anything prevented the request from completing.
> 
> You can add [missing error handling yourself](https://github.com/github/fetch#handling-http-error-statuses) or use axios or similar library instead.

**@alex-kovoy** on `ibeckermayer/teleport-interview` [→](https://github.com/ibeckermayer/teleport-interview/pull/2#discussion_r552904071)

> I would say that a client (in this case web ui) does not need to know these details, some generic message "Unable to process request" would suffice and we can log detailed Decode error message as is on the server side. 
> 
> If API itself was a product (as `stripe` for example) then I would of course give more attention to error messages returning to the client, but in this particular case it is not worth it. 
> 
> For example, would we want to show a user this message in UI?
> ```
> Request body contains an invalid value for the "xxx" field (at position 454)
> ```
> 
> Yeah, lets remove this for now (i do agree that Decode messages might be hard to read but I would not try to optimize it). 
> 
> One thing that I do not like is that the author logic will break if internals of the Decode method change, like this string comparsion `strings.HasPrefix(err.Error(), "json: unknown field ")`

**@russjones** on `ibeckermayer/teleport-interview` [→](https://github.com/ibeckermayer/teleport-interview/pull/2#discussion_r553698473)

> You don't want to `log.Fatal` here. Return an error here, you can `log.Fatal` in `main` if you want.

**@awly** on `ibeckermayer/teleport-interview` [→](https://github.com/ibeckermayer/teleport-interview/pull/3#discussion_r555979035)

> ```suggestion
> func (db *Database) CreateAccount(accountID, email, password string) error {
> ```

**@awly** on `ibeckermayer/teleport-interview` [→](https://github.com/ibeckermayer/teleport-interview/pull/3#discussion_r555983707)

> nit: handle and log the error

**@alex-kovoy** on `ibeckermayer/teleport-interview` [→](https://github.com/ibeckermayer/teleport-interview/pull/3#discussion_r556064980)

> Hm...looks a bit strange but as an option you can always use promises since each `then` is asynchronous :-) 
> 
> ```ts
>   api.post('/login', {
>     email: e.target.email.value,
>     password: e.target.password.value,
>   }).then(session => {
>     setSession(session);
>   })
>   .then(() => {
>     history.push('/dashboard');
>   })
>   .catch(error => {
>     console.error(error);
>   })
> ```

**@awly** on `ibeckermayer/teleport-interview` [→](https://github.com/ibeckermayer/teleport-interview/pull/3#discussion_r556142085)

> nit: when you write to http.ResponseWriter for the first time, it automatically writes the status code first (200 OK, unless WriteHeader was called before).
> so you can't write a different status after you've started writing the response body.
> 
> it's ok to just log the error but not send anything new in response.

**@alex-kovoy** on `ibeckermayer/teleport-interview` [→](https://github.com/ibeckermayer/teleport-interview/pull/5#discussion_r556993344)

> Lets introduce a wrapper function and use it to wrap protected handlers. Something like below 
> 
> ```golang
> func withAuth(handler ProtectedHandler) httprouter.Handle {
>     return func(w http.ResponseWriter, r *http.Request) {
>         session, err := authenticate(w, r)
>         if err != nil  {
>             http.Error(w, "Forbidden", http.StatusForbidden)
>             return
>         }
> 
>         // Pass a session 
>         handler.ServeHTTP(w, r, session)
>     }
> }
> 
> // and then
> srv.router.Handle("/api/logout", withAuth(logoutHandler)).Methods("DELETE")
> srv.router.Handle("/api/plan", withAuth(upgradePlanHandler)).Methods("PUT")
> ```

**@alex-kovoy** on `ibeckermayer/teleport-interview` [→](https://github.com/ibeckermayer/teleport-interview/pull/5#discussion_r557055542)

> Using request.Context is actually a good idea to pass a userId parameter. This way you do not need to have a custom handler methods to pass the session.
> 
> ```golang
> 
> func withAuth(handler http.Handler) httprouter.Handle {
>     return func(w http.ResponseWriter, r *http.Request) {
>         userId, err := authenticate(w, r)
>         if err != nil  {
>             http.Error(w, "Forbidden", http.StatusForbidden)
>             return
>         }
>  
>         r = r.WithContext( // add userId to the context)
>         handler.ServeHTTP(w, r)
>     }
> }
> ```

**@awly** on `ibeckermayer/teleport-interview` [→](https://github.com/ibeckermayer/teleport-interview/pull/7#discussion_r557704204)

> you could get a panic by forgetting to wrap a handler with something that puts the context key in.
> while sometimes programming errors warrant a panic, in this case it's detached far enough that a regular error would work better.

**@alex-kovoy** on `ibeckermayer/teleport-interview` [→](https://github.com/ibeckermayer/teleport-interview/pull/7#discussion_r557746211)

> you do not want to shutdown a service this way (it assumes that getting a websession from the context is the most important operation there is ignoring metrics collector, etc). Instead you can do something like:
> 
> ```go
> func (lh *LogoutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
>         // you can add FromContext method to your session manager
>         sessionId, err := sm.FromContext(r.Context())
>         if err != nil {
>           // return error + log
>         }
> 
>         sm.DeleteSession(sessionId)
>         w.WriteHeader(http.StatusNoContent)	
> }
> ```

**@russjones** on `ibeckermayer/teleport-interview` [→](https://github.com/ibeckermayer/teleport-interview/pull/7#discussion_r557811638)

> Why panic instead of returning an error here?

**@alex-kovoy** on `tobocop/go-teleport-directory-browser` [→](https://github.com/tobocop/go-teleport-directory-browser/pull/2#discussion_r698897950)

> (nit) With type inference it looks cleaner
> ```tsx
>  const [isLoading, setIsLoading] = useState(false);
>  const [error, setError] = useState('');
> ```

**@alex-kovoy** on `tobocop/go-teleport-directory-browser` [→](https://github.com/tobocop/go-teleport-directory-browser/pull/2#discussion_r698899271)

> instead of hardcoded numbers you can use `http.StatusInternalServerError`.

**@russjones** on `tobocop/go-teleport-directory-browser` [→](https://github.com/tobocop/go-teleport-directory-browser/pull/2#discussion_r698906967)

> 1) Use of a constant would be preferred here, so `http.StatusBadRequest` instead of `400`.
> 
> 2) This suggestion is minor, you don't need to make this change, but I usually use `http.Error` for this because it will set the response body with a human readable message as well.
> 
> ```go
> http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
> ```

**@alex-kovoy** on `tobocop/go-teleport-directory-browser` [→](https://github.com/tobocop/go-teleport-directory-browser/pull/3#discussion_r699741476)

> ```tsx
> catch(() => setError('Network error occurred. Please check your connection or try again later')
> ```
> 
> this will report network error even if error happens in the `then` part of the promise. Ideally this component should just report errors as is and `api` layer would need to provide this specific message.
> 
> ```tsx
>   api.authenticate(username, password)
>       .then((r) => {
>         // code
>       })
>       .catch(err => {
>         setError(err.message);
>       })
>       .finally(() => {|
>         setIsLoading(false);
>       })
> ```

**@alex-kovoy** on `tobocop/go-teleport-directory-browser` [→](https://github.com/tobocop/go-teleport-directory-browser/pull/3#discussion_r699742395)

> I would return an error message from the server together with a status code so the client can display it as is.

**@alex-kovoy** on `tobocop/go-teleport-directory-browser` [→](https://github.com/tobocop/go-teleport-directory-browser/pull/3#discussion_r699746185)

> This message is used in several places now `Network error occurred. Please check your connection or try again later` does it make sense to throw it from `api` layer?

**@alex-kovoy** on `zship/teleport-challenge` [→](https://github.com/zship/teleport-challenge/pull/1#discussion_r718771332)

> I would add an error message as well, this will greatly simplify client side error handling.

**@kimlisa** on `zship/teleport-challenge` [→](https://github.com/zship/teleport-challenge/pull/4#discussion_r725394917)

> if i'm following your code correctly, the http code for this will default to `500` since it's not stated, which seems incorrect b/c a missing password in the request is a client side error?

**@alex-kovoy** on `zship/teleport-challenge` [→](https://github.com/zship/teleport-challenge/pull/4#discussion_r725677710)

> I would introduce a new Error class (in Javascript) and use it instead of this parsing code.

**@alex-kovoy** on `zship/teleport-challenge` [→](https://github.com/zship/teleport-challenge/pull/5#discussion_r726692929)

> ```js
> if ('code' in json && typeof json.code === 'string') {
>       code = json.code as ErrorCode;
>     }
> 
> let message = '';
> if ('message' in json && typeof json.message === 'string') {
>       message = json.message;
> }
> ```
> 
> I assume you can trust your backend API so instead of ^^^ you can do something like below (note that `response.ok` flag as handled)
> 
> ```ts
> if (!response.ok) {
>  throw new ApiError(response) 
> }
> ```

**@alex-kovoy** on `zship/teleport-challenge` [→](https://github.com/zship/teleport-challenge/pull/5#discussion_r726703188)

> why do we need to check the instance type? Does the code may throw anything else besides error? If so, should this check be used in every `catch` block?

---

## api-shape

_Library importability; public functions returning private types; io.Reader vs callback APIs; protobuf types in public API; snapshot atomicity._

**142 quotes** from `24` distinct reviewers across `42` candidate submissions.

**@codingllama** on `adalton/teleport-exercise` [→](https://github.com/adalton/teleport-exercise/pull/1#discussion_r768035067)

> Public funcs returning private types are a bit off - we can have assign variables and rely on type inference, but we can't explicitly type our vars, which is a bit weird.
> 
> Same for other occurrences.

**@codingllama** on `adalton/teleport-exercise` [→](https://github.com/adalton/teleport-exercise/pull/1#discussion_r768049189)

> Should this be configurable at the library level? Maybe we could we default to the current binary instead?

**@codingllama** on `adalton/teleport-exercise` [→](https://github.com/adalton/teleport-exercise/pull/1#discussion_r768878383)

> I think this one is worth having, but we can do it as a follow up to the library PR if you prefer. Your call.

**@codingllama** on `adalton/teleport-exercise` [→](https://github.com/adalton/teleport-exercise/pull/2#discussion_r771455502)

> +1, it's important to be able to run the suite without relying on specific ports being available.

**@dboslee** on `benmoss/job-worker-service` [→](https://github.com/benmoss/job-worker-service/pull/1#discussion_r2670121828)

> could you add go doc comments for the API? some of the fields/funcs are not immediately clear how they are intended to be used eg. `Job.Env`, `Job.Dir`,  `JobWorker.Logs(uint64, bool)`

**@nklaassen** on `benmoss/job-worker-service` [→](https://github.com/benmoss/job-worker-service/pull/1#discussion_r2670612883)

> And how will you keep track of the owner of each job? would this be handled at the library layer or the api layer? I don't think I can see how it will be tracked in the current version of this doc

**@nklaassen** on `benmoss/job-worker-service` [→](https://github.com/benmoss/job-worker-service/pull/3#discussion_r2683948413)

> it's a bit awkward to add the cgFD close to the cleanups slice here. Why not add it immediately after getting cgFD? Or, why not just directly call syscall.Close here instead of creating an unnecessary anonymous function? It's also a bit awkward to have the syscall.Close call duplicated below on the happy path.
>
> In general for callers of cg.fd() it seems easy to forget to close the file descriptor in some path. I would love to see the API for cg.fd() make it very difficult to forget to close it

**@nklaassen** on `benmoss/job-worker-service` [→](https://github.com/benmoss/job-worker-service/pull/3#discussion_r2684028835)

> it's not going to be tied to a specific API call so it seems the best you can do is log it

**@dboslee** on `benmoss/job-worker-service` [→](https://github.com/benmoss/job-worker-service/pull/3#discussion_r2684469961)

> does this need to be exported?

**@sclevine** on `benmoss/job-worker-service` [→](https://github.com/benmoss/job-worker-service/pull/3#discussion_r2684556329)

> Should these be exported?

**@sclevine** on `benmoss/job-worker-service` [→](https://github.com/benmoss/job-worker-service/pull/3#discussion_r2684586756)

> *nit*: document exported identifiers

**@smallinsky** on `bill-rich/jobworker` [→](https://github.com/bill-rich/jobworker/pull/1#discussion_r1606445480)

> Could you provide example `proto` API definitions?

**@sclevine** on `bkneis/jobworker` [→](https://github.com/bkneis/jobworker/pull/1#discussion_r1648011884)

> Can you be more specific about how your library will "tail the contents of the file"?
> 
> Can you add a section justifying why you've chosen to store logs on disk instead of in memory? Remember that we are not looking for a memory-optimized solution.

**@sclevine** on `bkneis/jobworker` [→](https://github.com/bkneis/jobworker/pull/1#discussion_r1648300806)

> How does the CLI find the API server?

**@sclevine** on `bkneis/jobworker` [→](https://github.com/bkneis/jobworker/pull/1#discussion_r1648305266)

> Not sure I understand the motivation for this interface. Are these methods also implemented on JobWorker? If so, why provide an interface?

**@sclevine** on `bkneis/jobworker` [→](https://github.com/bkneis/jobworker/pull/1#discussion_r1648305835)

> Does the library need this centralized struct that stores jobs? Would this functionality make more sense in the API server?

**@sclevine** on `bkneis/jobworker` [→](https://github.com/bkneis/jobworker/pull/1#discussion_r1648308852)

> Allowing these values to be specified as strings in user-facing UX (e.g., CLI flags) makes sense, but I’m on the fence about whether it’s the best interface for the gRPC API / Go library. Curious how about your reasoning.

**@tcsc** on `bkneis/jobworker` [→](https://github.com/bkneis/jobworker/pull/1#discussion_r1650235798)

> How are you planning on embedding the job `init` script in your worker library?

**@sclevine** on `bkneis/jobworker` [→](https://github.com/bkneis/jobworker/pull/1#discussion_r1651348210)

> Same question here about interfaces. What's the motivation for wrapping jobs in an interface? Do you plan to have multiple implementations of this?

**@sclevine** on `bkneis/jobworker` [→](https://github.com/bkneis/jobworker/pull/1#discussion_r1651477736)

> `cmd` is already unexported, so the interface does not change this

**@tcsc** on `bkneis/jobworker` [→](https://github.com/bkneis/jobworker/pull/2#discussion_r1663390955)

> Ideally the controller interface filenames should be in constants, but not a big deal for this exercise.

**@zmb3** on `bucknercd/jobworker` [→](https://github.com/bucknercd/jobworker/pull/1#discussion_r2271529485)

> Similarly, how will the low/med/high class from your proto map to `rbps` and `wbps`?

**@zmb3** on `bucknercd/jobworker` [→](https://github.com/bucknercd/jobworker/pull/1#discussion_r2271548157)

> I'd love to see what the Go API looks like at the library level, too.

**@strideynet** on `bucknercd/jobworker` [→](https://github.com/bucknercd/jobworker/pull/1#discussion_r2272851064)

> Nit: The Protobuf best practices advises against the zero value of the enum being semantically meaningful and suggest the inclusion of an explicit `UNSPECIFIED` value - https://protobuf.dev/best-practices/dos-donts/

**@espadolini** on `chintamanil/job-worker` [→](https://github.com/chintamanil/job-worker/pull/2#discussion_r2833505636)

> If these fields should be guarded by the mutex, should they really be exported?

**@strideynet** on `clydotron/job_worker_service` [→](https://github.com/clydotron/job_worker_service/pull/1#discussion_r1600234555)

> Will this only be the messages from STDOUT or STDERR? Or will both be written to the same field? It might be nice to be able to distinguish between them.

**@strideynet** on `clydotron/job_worker_service` [→](https://github.com/clydotron/job_worker_service/pull/1#discussion_r1600239914)

> Could a proto enum be used rather than a string ? This gives much more type safety.

**@strideynet** on `clydotron/job_worker_service` [→](https://github.com/clydotron/job_worker_service/pull/1#discussion_r1603092214)

> I find that proto files are an ideal way to document APIs, comments on the methods and the message fields would be helpful in understanding their behaviour. E.g Does `StopJob` wait for the job to stop before responding ?

**@rosstimothy** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/1#discussion_r1782731358)

> For the sake of simplicity you may hard code a static set of limits on the server and omit them from the API.

**@codingllama** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/1#discussion_r1783297360)

> Why the in-band message instead of using the gRPC error side-channel?
> 
> How can one tell whether the response is successful? Is a non-empty message always a failure? Or an empty job_id?

**@espadolini** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/2#discussion_r1787687089)

> Key encipherment is only needed for non-PFS RSA, but thankfully it's 2024 and we don't have to do that anymore (especially seeing as we're going to rightfully use ECDSA); from what I can tell from the cfssl docs, "signing" and "digital signature" are mapped to the same thing?
> 
> ```suggestion
>         "digital signature",
> ```

**@espadolini** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/2#discussion_r1787687380)

> ```suggestion
>           "digital signature",
> ```

**@espadolini** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/2#discussion_r1787698116)

> I personally really like `--go(-grpc)_opt=module=go.creack.net/telepilot` rather than `paths=source_relative` to define where the root of my go module is and have the files generated based on their go package name rather than their proto file path; it leads (or should lead) to the same result here, regardless.

**@codingllama** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/2#discussion_r1787760234)

> Make JobManager a struct and ditch the interface?
> 
> * https://go.dev/wiki/CodeReviewComments#interfaces

**@espadolini** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/4#discussion_r1789779243)

> This doesn't actually do any connection, and there's no way to have a connection in grpc-go without calling a RPC, using deprecated APIs or using experimental ones. 😅
> 
> Maybe "returns a client prepared to connect to the server" or something along those lines would be more appropriate.

**@rosstimothy** on `devnulled/runjob` [→](https://github.com/devnulled/runjob/pull/2#discussion_r1795662377)

> You may omit this from the API for time/simplicity of this challenge and hardcode a static set of limits that are applied to all jobs on the server.

**@smallinsky** on `devnulled/runjob` [→](https://github.com/devnulled/runjob/pull/2#discussion_r1799944899)

> nit: The name response is quite vague for a gRPC response message field in the context of
>  Request/Response messages:
>  
> ```suggestion
>     bytes output = 1;
> ```

**@awly** on `diptadas/job-worker` [→](https://github.com/diptadas/job-worker/pull/1#discussion_r645738171)

> Let's keep the design doc high-level.
> We can sort out function signatures during code review.
> The important things here are: how are jobs tracked, where does the output go, what failures are possible (fail to start, preemptive stop, etc)

**@zmb3** on `DiscoRiver/jobber` [→](https://github.com/DiscoRiver/jobber/pull/1#discussion_r1989319489)

> ```suggestion
> - A discrete CLI program for API interface.
> ```

**@zmb3** on `DiscoRiver/jobber` [→](https://github.com/DiscoRiver/jobber/pull/1#discussion_r1990183316)

> How does this change impact the API? Will you still have a channel based API here, or does the move to temp files allow you to present a different API?
> 
> What will your notification mechanism be with files? Once a reader catches up to current time, how will they be notified that new data has been written?

**@russjones** on `dtom90/teleport-backend-challenge` [→](https://github.com/dtom90/teleport-backend-challenge/pull/1#discussion_r2053136083)

> Can you also share your `.proto` file?

**@russjones** on `dtom90/teleport-backend-challenge` [→](https://github.com/dtom90/teleport-backend-challenge/pull/1#discussion_r2053137045)

> * How will you know new output has been written to the file?
> * Can you tell me the interface you plan to expose in your Go program?

**@smallinsky** on `dtom90/teleport-backend-challenge` [→](https://github.com/dtom90/teleport-backend-challenge/pull/1#discussion_r2055361379)

> Does the command’s arguments also get forwarded via the command field?
> 
> Does  it make sense to extend the message to include a repeated string args field, so that the command and its arguments can be structured separately?

**@zmb3** on `ehsaniara/joblet` [→](https://github.com/ehsaniara/joblet/pull/2#discussion_r2114129947)

> I might suggest separate request and response messages for each RPC so that they can grow and change independently of each other.

**@russjones** on `flychicken123/Job_Worker` [→](https://github.com/flychicken123/Job_Worker/pull/1#discussion_r2065389097)

> Can you share your `.proto` file?

**@tigrato** on `flychicken123/Job_Worker` [→](https://github.com/flychicken123/Job_Worker/pull/1#discussion_r2066124992)

> what @russjones meant was how the library will stop the process.

**@tigrato** on `flychicken123/Job_Worker` [→](https://github.com/flychicken123/Job_Worker/pull/1#discussion_r2066125837)

> You can use uuid library for this

**@rosstimothy** on `flychicken123/Job_Worker` [→](https://github.com/flychicken123/Job_Worker/pull/1#discussion_r2069333587)

> What is HashMap referring to here? Is that provided to you by the Go standard library or something you intend to implement yourself?

**@eriktate** on `GevorgGal/jobworker` [→](https://github.com/GevorgGal/jobworker/pull/1#discussion_r2975158872)

> Can you include which library you plan on using to generate UUIDs?

**@Tener** on `GevorgGal/jobworker` [→](https://github.com/GevorgGal/jobworker/pull/1#discussion_r2976832985)

> The design sounds good to me. A careful API design should help with the locking requirement.
>
> I'm curious to hear if you have considered any other designs and what led to your current choice?

**@rosstimothy** on `GevorgGal/jobworker` [→](https://github.com/GevorgGal/jobworker/pull/2#discussion_r2982803516)

> > func (b *OutputBuffer) ReadFrom(offset int) (data []byte, done bool, notify <-chan struct{}) {
>
> This API is clunky, inefficient (as Krzysztof pointed out), and requires a lot of extra effort from callers to get the output data. Can you think of any ways to better encapsulate some of this from end users to improve the API and remove potential footguns?

**@rosstimothy** on `GevorgGal/jobworker` [→](https://github.com/GevorgGal/jobworker/pull/2#discussion_r2982813129)

> Why does this information need to be exposed in the public API?

**@Tener** on `GevorgGal/jobworker` [→](https://github.com/GevorgGal/jobworker/pull/2#discussion_r2987465586)

> We can make an upper cap on the amount of data returned per call. This keeps the interface intact and yet avoids potential OOM.

**@strideynet** on `gstelang/job-worker-service` [→](https://github.com/gstelang/job-worker-service/pull/4#discussion_r1740685812)

> I think for as long as your interfaces are defined alongside the implementation, you end up with fairly weak abstractions that are tightly coupled to that specific implementation.

**@sclevine** on `ilyazz/jobs` [→](https://github.com/ilyazz/jobs/pull/2#discussion_r1001273484)

> Could you provide the `.proto` file that defines the gRPC API?

**@sclevine** on `ilyazz/jobs` [→](https://github.com/ilyazz/jobs/pull/2#discussion_r1001280015)

> *nit:* This API looks great for gRPC client/server communications, but it feels complex for non-server consumers of the library. If I just want to start/stop job processes in Go, what value do I get from the `Supervisor`? Are contexts useful for all of these operations?
> 
> *smaller nit:* https://github.com/golang/go/wiki/CodeReviewComments#interfaces

**@sclevine** on `ilyazz/jobs` [→](https://github.com/ilyazz/jobs/pull/2#discussion_r1002840385)

> > interface is used here only as a shorter way to show a list of API methods, and also New() should return an object of course, not an interface.
> 
> > The general purpose API for running jobs would be close to exec.Cmd with cgroup/ns related config options. Actually, the current design has JobHandler object inside, it is not a part of API now, but could be.
> 
> Would you like to proceed with the library public API as described in the RFD, or update it to be more general purpose?

**@sclevine** on `ilyazz/jobs` [→](https://github.com/ilyazz/jobs/pull/2#discussion_r1002842903)

> > to be able to make last minute changes in client/server protocol without discussions. :)
> 
> It's okay if you need to make API changes as you're implementing the challenge.

**@sclevine** on `ilyazz/jobs` [→](https://github.com/ilyazz/jobs/pull/2#discussion_r1004016769)

> > Sure, I can export internal Job implementation as library API.
> 
> Sounds good. I want to make sure you have time to get through the implementation, and this RFD looks good either way.
> 
> >  I'm not sure is such a good idea to duplicate the client process all the time.
> 
> This is a great point! I think relying on `/proc/self/exe` as you've proposed is a reasonable compromise for this challenge.
> 
> If you're concerned about conflicting with an unknown library client's `main`/`init`, maybe you could require the client to invoke a specific library function within `main` or `init` to more explicitly "enable" the library? Then indirect imports, etc., shouldn't cause problems.

**@rosstimothy** on `jaylane/job-scheduler` [→](https://github.com/jaylane/job-scheduler/pull/1#discussion_r1714210168)

> Feel free to reference the .proto file included instead of duplicating the information here.

**@AntonAM** on `jaylane/job-scheduler` [→](https://github.com/jaylane/job-scheduler/pull/1#discussion_r1714301193)

> I think you can include proto definitions right in this document.

**@AntonAM** on `jaylane/job-scheduler` [→](https://github.com/jaylane/job-scheduler/pull/1#discussion_r1714309533)

> ```suggestion
> message StartJobResponse {
> ```

**@tigrato** on `joshuarubin/teleport-job-worker` [→](https://github.com/joshuarubin/teleport-job-worker/pull/4#discussion_r1746718910)

> does this need to be atomic?

**@codingllama** on `joshuarubin/teleport-job-worker` [→](https://github.com/joshuarubin/teleport-job-worker/pull/4#discussion_r1747501128)

> You've placed the entire API inside internal/, so that makes it not importable for projects outside of teleport-job-worker. The intent for the "library" part of the project is that it could be imported and used outside of the project, so we should move the relevant parts out of internal/.

**@codingllama** on `joshuarubin/teleport-job-worker` [→](https://github.com/joshuarubin/teleport-job-worker/pull/4#discussion_r1747605439)

> Could we ditch the interface and expose JobWorker directly?
> 
> The interface is useful to express the API in the design phase, but I don't think the API actually calls for it.
> 
> - https://go.dev/wiki/CodeReviewComments#interfaces

**@tigrato** on `kiakeshmiri/process-runner` [→](https://github.com/kiakeshmiri/process-runner/pull/1#discussion_r1771346824)

> the fact that the code lives in internal package, forbids any library to import it

**@smallinsky** on `kkloberdanz/teleport-challenge` [→](https://github.com/kkloberdanz/teleport-challenge/pull/6#discussion_r2797782607)

> In the current form it looks more like helper method rather than public client interface. 
> 
> For instance  each call will  need to create new client and trigger a new TLS handshake and this will be pretty inefficient and cumbersome for the library usage where for each call the address needs to be provided. 
> 
> I understand this may be just scaffolding, but my concern is that we’ll end up previning the same functionality multiple times during follow-up refactors.

**@smallinsky** on `kkloberdanz/teleport-challenge` [→](https://github.com/kkloberdanz/teleport-challenge/pull/7#discussion_r2811931264)

> I’m a bit concerned about coupling the client library’s Go domain types directly to the protobuf transport-layer types. Would it make sense to keep the transport types internal to the client, and expose explicit library-level types instead (for example job.Status insead of `pb.JobStatus` in this case)?
> 
> That way, we keep a clean API boundary.

**@tigrato** on `kkloberdanz/teleport-challenge` [→](https://github.com/kkloberdanz/teleport-challenge/pull/8#discussion_r2832873839)

> does this type need to be exported?

**@zmb3** on `kovyrin/teleport-exec` [→](https://github.com/kovyrin/teleport-exec/pull/1#discussion_r788224317)

> Could you also include an example of what the API of the library might look like?

**@russjones** on `kovyrin/teleport-exec` [→](https://github.com/kovyrin/teleport-exec/pull/1#discussion_r789267212)

> :+1: to creating own library.

**@zmb3** on `kovyrin/teleport-exec` [→](https://github.com/kovyrin/teleport-exec/pull/1#discussion_r790012254)

> I disagree - the challenge has 3 parts:
> 
> 1. a library with an API for working with jobs
> 2. a GRPC server that uses the API of the library
> 3. a GRPC client
> 
> In this sense, the library API *is* public, and a well designed API is important in order to ensure that you can cleanly integrate the library functionality into your GRPC server.

**@kpumuk** on `kovyrin/teleport-exec` [→](https://github.com/kovyrin/teleport-exec/pull/2#discussion_r796981763)

> Define constants or params for the library?

**@kpumuk** on `kovyrin/teleport-exec` [→](https://github.com/kovyrin/teleport-exec/pull/2#discussion_r797003835)

> Wonder if it would make more sense to spit error messages into stderr...

**@sclevine** on `kurczynski/teleport-job-worker` [→](https://github.com/kurczynski/teleport-job-worker/pull/1#discussion_r1612586062)

> This method implies that the library stores a global map of jobs. Is this necessary?

**@sclevine** on `kurczynski/teleport-job-worker` [→](https://github.com/kurczynski/teleport-job-worker/pull/1#discussion_r1612590223)

> Note that you can use the google.protobuf types if you wish: https://protobuf.dev/reference/protobuf/google.protobuf/#timestamp

**@sclevine** on `kurczynski/teleport-job-worker` [→](https://github.com/kurczynski/teleport-job-worker/pull/1#discussion_r1617673471)

> Looks like the library is missing configuration for resource controls?

**@sclevine** on `kurczynski/teleport-job-worker` [→](https://github.com/kurczynski/teleport-job-worker/pull/1#discussion_r1617674569)

> Library should be separate from the API server / not coupled to gRPC.
> 
> (Assuming `job` is a generated protobuf package.)

**@bernardjkim** on `kurczynski/teleport-job-worker` [→](https://github.com/kurczynski/teleport-job-worker/pull/2#discussion_r1624878239)

> Do these functions `SetResourceLimits` and `SetCommand` need to be exported? Could there be any problems if these are called multiple times?

**@greedy52** on `MarkDHarris/teleport-systems-challenge` [→](https://github.com/MarkDHarris/teleport-systems-challenge/pull/1#discussion_r3017031689)

> any reasoning validating invalid args at the gRPC level versus in client and the library?

**@rosstimothy** on `MarkDHarris/teleport-systems-challenge` [→](https://github.com/MarkDHarris/teleport-systems-challenge/pull/2#discussion_r3035834962)

> How come the id and owner are exported? Are there any concerns with that?

**@rosstimothy** on `MarkDHarris/teleport-systems-challenge` [→](https://github.com/MarkDHarris/teleport-systems-challenge/pull/2#discussion_r3035844036)

> It seems that by only checking the context here we may provide more output to the callback function for a client that has disconnected.

**@eriktate** on `MrChristianL/Teleport-Job-Worker-Service` [→](https://github.com/MrChristianL/Teleport-Job-Worker-Service/pull/1#discussion_r2881227815)

> How will each job be managed (created, queried, etc.) in your library? What will the exported API look like?

**@nklaassen** on `MrChristianL/Teleport-Job-Worker-Service` [→](https://github.com/MrChristianL/Teleport-Job-Worker-Service/pull/1#discussion_r2896980826)

> suggestion: `func([]byte) error` looks almost like `io.Writer.Write`, I'd consider either having this function accept an io.Writer or return an io.Reader to make it more composable

**@rosstimothy** on `neildo/tjob` [→](https://github.com/neildo/tjob/pull/1#discussion_r1759131649)

> Would this API be able to scale with a large number of jobs?

**@nklaassen** on `neildo/tjob` [→](https://github.com/neildo/tjob/pull/1#discussion_r1759274770)

> nit: i suggest lower_snake_case for field names https://protobuf.dev/programming-guides/style/#message-field-names

**@nklaassen** on `neildo/tjob` [→](https://github.com/neildo/tjob/pull/1#discussion_r1759288904)

> this seems odd but i'll avoid doing code review in a design doc. it's probably fine to leave out the unexported fields from this design doc and just describe the exported API

**@nklaassen** on `neildo/tjob` [→](https://github.com/neildo/tjob/pull/2#discussion_r1767823195)

> nit: exported identifiers (constants, variables, types, functions, methods) should have GoDoc comments

**@rosstimothy** on `neildo/tjob` [→](https://github.com/neildo/tjob/pull/2#discussion_r1768522227)

> What happens if either of these times are written to concurrently while consumers of the library call these functions?

**@rosstimothy** on `neildo/tjob` [→](https://github.com/neildo/tjob/pull/2#discussion_r1768537631)

> This might aid local testing, however, there is no requirement for this challenge that the system works on macOS. Consider removing it instead of having to support two separate implementations.

**@rosstimothy** on `neildo/tjob` [→](https://github.com/neildo/tjob/pull/2#discussion_r1769131957)

> So if a user of the library reads and receives an ErrReadAgain are they supposed to ignore the error and call Read again? Should they do so immediately? Should they wait first?
> 
> That contract seems a bit unidiomatic. Should this just try reading again itself?

**@smallinsky** on `P-A-R-U-S/Go-Job-Worker-Service` [→](https://github.com/P-A-R-U-S/Go-Job-Worker-Service/pull/2#discussion_r1743742364)

> This will result in the process using the library terminating.
> Why not simply return an error and allow the caller to decide how to handle it?

**@rhammonds-teleport** on `rajansandeep/teleport-job-worker-svc` [→](https://github.com/rajansandeep/job-worker-svc/pull/1#discussion_r3079945419)

> The library will also support "early exit" for readers, right?

**@eriktate** on `rajansandeep/teleport-job-worker-svc` [→](https://github.com/rajansandeep/job-worker-svc/pull/1#discussion_r3093548921)

> nit: `job` is unexported in the definition below

**@atburke** on `rajansandeep/teleport-job-worker-svc` [→](https://github.com/rajansandeep/job-worker-svc/pull/2#discussion_r3101765447)

> Do reading and writing need separate wait groups?

**@atburke** on `rajansandeep/teleport-job-worker-svc` [→](https://github.com/rajansandeep/job-worker-svc/pull/2#discussion_r3101958147)

> If I'm calling this library, what do I do when I want to shut down my Worker?

**@zmb3** on `razzam21/job-worker-service` [→](https://github.com/razzam21/job-worker-service/pull/1#discussion_r2257564662)

> As above, authentication and authorization should part of the server, not the library.

**@awly** on `renatoaguimaraes/golang-job-scheduler` [→](https://github.com/renatoaguimaraes/golang-job-scheduler/pull/2#discussion_r623373792)

> `syscall.InotifyInit` should return a non-nil error based on `errno`, so no need to check `filed` separately

**@r0mant** on `renatoaguimaraes/golang-job-scheduler` [→](https://github.com/renatoaguimaraes/golang-job-scheduler/pull/2#discussion_r623529658)

> Do not panic in a library function. Alternatively, if this is only used in tests, consider just using `t.TempDir()` instead of this NewConfig method.

**@r0mant** on `renatoaguimaraes/golang-job-scheduler` [→](https://github.com/renatoaguimaraes/golang-job-scheduler/pull/5#discussion_r625202229)

> Don't use log.Fatal in the library code, return an error and let the caller (server cmd package in this case) handle it.

**@awly** on `renatoaguimaraes/golang-job-scheduler` [→](https://github.com/renatoaguimaraes/golang-job-scheduler/pull/5#discussion_r625963526)

> nit: does this need to be exported?

**@rosstimothy** on `rsteinkeXJ/teleport_interview` [→](https://github.com/rsteinkeXJ/teleport_interview/pull/1#discussion_r3127513826)

> I suggest changing this to a unary RPC. It will separate concerns and also be easier on consumers. The API as is introduces a fair amount of complexity. We have to consume the first reply, assert the type is the expected one, parse the id out of the oneof payload.

**@smallinsky** on `samschurter/teleport-challenge` [→](https://github.com/samschurter/teleport-challenge/pull/1#discussion_r692311588)

> It it not clear to me why worker library uses IPC mechanism.

**@bernardjkim** on `samschurter/teleport-challenge` [→](https://github.com/samschurter/teleport-challenge/pull/1#discussion_r693079545)

> On the client side how will the output be displayed if receiving two separate fields?

**@smallinsky** on `samschurter/teleport-challenge` [→](https://github.com/samschurter/teleport-challenge/pull/1#discussion_r693236968)

> We would rather prefer simplicity so if there is no technical reason why the worker library should expose its API via IPC to the HTTP server there is no reason to make the design complex. The worker library can be designed independently without using IPC.

**@zmb3** on `shawon-crosen/teleport-challenge` [→](https://github.com/shawon-crosen/teleport-challenge/pull/1#discussion_r2129153313)

> Protobuf `string`s must be valid UTF-8, and we can't make assumptions about the process output. I would use `bytes` here instead.

**@zmb3** on `shawon-crosen/teleport-challenge` [→](https://github.com/shawon-crosen/teleport-challenge/pull/1#discussion_r2129161999)

> Nit: `builtin` in Go refers to the builtin funcs (append, make, len, cap, etc). I would 
> ```suggestion
> The library will primarily utilize the standard library's os/exec package, with some use of the os/user and os packages.
> ```

**@codingllama** on `tjper/teleport` [→](https://github.com/tjper/teleport/pull/2#discussion_r827452225)

> nit: push the interface closer to where it is used, not where the types are declared.

**@codingllama** on `tjper/teleport` [→](https://github.com/tjper/teleport/pull/3#discussion_r829317835)

> I think we should roll this into the library, seems like it's a responsibility it should have.

**@tigrato** on `xSudoNymx/job-manager` [→](https://github.com/xSudoNymx/job-manager/pull/1#discussion_r2705317215)

> feel free to mid both stderr and stdout. there is no requirement for keeping them separate

**@sclevine** on `zeeshanhaque21/JobWorkerService` [→](https://github.com/zeeshanhaque21/JobWorkerService/pull/1#discussion_r1686910138)

> Why should the library care about users or Job IDs? These seem like gRPC API server concerns.

**@sclevine** on `zeeshanhaque21/JobWorkerService` [→](https://github.com/zeeshanhaque21/JobWorkerService/pull/1#discussion_r1693302505)

> I'm not sure I understand the motivation for the complexity and limitations of this API design. I would expect a consumer of the library to want to consume an arbitrary number of bytes on each read, or less if less data is available.

**@sclevine** on `zeeshanhaque21/JobWorkerService` [→](https://github.com/zeeshanhaque21/JobWorkerService/pull/1#discussion_r1696241908)

> ```suggestion
> func (l *Logger) GetLogReader() io.Reader
> ```
> Avoid pointer to interface:
> https://go.dev/wiki/CodeReviewComments#pass-values

**@espadolini** on `Zephan92/teleport-handson` [→](https://github.com/Zephan92/teleport-handson/pull/5#discussion_r2759696825)

> nit: prefer getters when reading from a protobuf message
>
> ```suggestion
> 	slog.Info("Start request", "cmd", req.GetCommand(), "dir", req.GetWorkingDirectory())
> ```

**@alex-kovoy** on `ibeckermayer/teleport-interview` [→](https://github.com/ibeckermayer/teleport-interview/pull/1#discussion_r551433914)

> Yes, you can use `rand.Reader.Read(32)` to generate a random string with crypto-strong pseudo random generator. This will give you an opaque token that you can use for session/api keys.

**@alex-kovoy** on `ibeckermayer/teleport-interview` [→](https://github.com/ibeckermayer/teleport-interview/pull/1#discussion_r551612973)

> I would also add HTML headers to APIs such as:
> 
> ```
> X-XSS-Protection
> X-Frame-Options
> Strict-Transport-Security
> ```

**@russjones** on `ibeckermayer/teleport-interview` [→](https://github.com/ibeckermayer/teleport-interview/pull/2#discussion_r552314452)

> Shouldn't you just return a 404 and only `index.html` if it's requested? @alex-kovoy thoughts?

**@alex-kovoy** on `ibeckermayer/teleport-interview` [→](https://github.com/ibeckermayer/teleport-interview/pull/2#discussion_r552391932)

> I would actually keep the index route for the dashboard and move the Login to `/login` route.
> 
> Also, consider adding Authenticated (or similar) to secure private routes and handle redirects to the login screen when access token is invalid.
> 
> ```tsx
>       <Switch>
>         <Route exact path="/">
>           <Login />
>         </Route>
>         <Route path="/">
>            <Authenticated>
>               <Dashboard />
>             <Authenticated/>
>         </Route>
>       </Switch>
> ```

**@alex-kovoy** on `ibeckermayer/teleport-interview` [→](https://github.com/ibeckermayer/teleport-interview/pull/2#discussion_r552396121)

> also consider creating some wrappers for get|post requests to incapsulate and reuse the code that initializes the HTTP request. 
> 
> ```ts
> api.get("/xxx")
> api.post('/xxx', json)
> ```

**@awly** on `ibeckermayer/teleport-interview` [→](https://github.com/ibeckermayer/teleport-interview/pull/2#discussion_r552918563)

> you're vendoring the entire godoc.org implementation for a relatively simple function.
> you can replace it with `if !strings.HasPrefix(r.Header.Get("Content-Type", "application/json"))` and drop a large dependency

**@awly** on `ibeckermayer/teleport-interview` [→](https://github.com/ibeckermayer/teleport-interview/pull/2#discussion_r552923709)

> This seems to create a tight coupling between frontend and backend logic.
> But I don't really write web applications these days, so maybe I just need to get on with the times.

**@awly** on `ibeckermayer/teleport-interview` [→](https://github.com/ibeckermayer/teleport-interview/pull/2#discussion_r552924499)

> ```suggestion
> 	r.Body = http.MaxBytesReader(w, r.Body, 1048576) // 1MB
> ```

**@alex-kovoy** on `ibeckermayer/teleport-interview` [→](https://github.com/ibeckermayer/teleport-interview/pull/3#discussion_r555537683)

> instead of HOC I would use a regular component as it's easier to read and follow, especially since we need only 1 context provider per application.
> 
> ```ts
> const AppContext = props => {
>   // 
> }
> ```
> and then
> 
> ```tsx
>  <AppContext>
>     <Router>
>       <Switch>
>         <Route exact path="/">
>           <Login />
>         </Route>
>         <Route path="/dashboard">
>           <Dashboard />
>         </Route>
>       </Switch>
>     </Router>
> </AppContext>
> ```

**@alex-kovoy** on `ibeckermayer/teleport-interview` [→](https://github.com/ibeckermayer/teleport-interview/pull/3#discussion_r556049201)

> > But as long as you're ok with the quicker-to-implement getSession solution for this project I'll go with that.
> 
> Yeah, i am fine with this performance overhead in favor of a better decoupling between components and ease of use. Also, I would not worry about performance optimizations at such level at this stage.

**@alex-kovoy** on `tobocop/go-teleport-directory-browser` [→](https://github.com/tobocop/go-teleport-directory-browser/pull/2#discussion_r698895986)

> I would have a separate route for `/login` and the rest treat as authenticated

**@alex-kovoy** on `tobocop/go-teleport-directory-browser` [→](https://github.com/tobocop/go-teleport-directory-browser/pull/2#discussion_r698896575)

> `createContext<ApiClient>(null);`

**@russjones** on `tobocop/go-teleport-directory-browser` [→](https://github.com/tobocop/go-teleport-directory-browser/pull/2#discussion_r698905664)

> Wouldn't the JSON parser catch this?
> 
> Not sure if someone could stage an attack here, but I could be wrong? Any links?

**@alex-kovoy** on `tobocop/go-teleport-directory-browser` [→](https://github.com/tobocop/go-teleport-directory-browser/pull/3#discussion_r699743220)

> will it work instead?
> ```tsx
> <AuthenticatedRoute path={`${Routes.ROOT}`}>
> ```

**@alex-kovoy** on `tobocop/go-teleport-directory-browser` [→](https://github.com/tobocop/go-teleport-directory-browser/pull/3#discussion_r699751271)

> will there be a situation when on rerender this hook returns a new instance of `api` object ?

**@alex-kovoy** on `tobocop/go-teleport-directory-browser` [→](https://github.com/tobocop/go-teleport-directory-browser/pull/3#discussion_r699752405)

> does this component need to deal with http codes?

**@alex-kovoy** on `tobocop/go-teleport-directory-browser` [→](https://github.com/tobocop/go-teleport-directory-browser/pull/3#discussion_r699795225)

> Below example will ensure that client instance is preserved during re-renders (it does not rerender in your case but it would if ContextProvider is placed inside BrowserRouter).
> 
> ```tsx
> export const AppContextProvider = ({ children }: { children: ReactNode }) => {
>   const context = useMemo(() => new AppContext());
>   return (
>   <AppContext.Provider value={context}>
>     {children}
>   </AppContext.Provider>
> );
> ```

**@alex-kovoy** on `zship/teleport-challenge` [→](https://github.com/zship/teleport-challenge/pull/1#discussion_r718767697)

> Are you trying to say that all requests that do not match `/api` route will return an index.html? If so I would clarify this message as it seems that `index.html` will be returned for a directory path requests that does not exist.

**@alex-kovoy** on `zship/teleport-challenge` [→](https://github.com/zship/teleport-challenge/pull/1#discussion_r718775250)

> What are the advantages of using `Redux` for this app?

**@alex-kovoy** on `zship/teleport-challenge` [→](https://github.com/zship/teleport-challenge/pull/1#discussion_r718775450)

> How will the routes be handled?

**@alex-kovoy** on `zship/teleport-challenge` [→](https://github.com/zship/teleport-challenge/pull/1#discussion_r718907189)

> Considering that Redux is loosing it's popularity (mainly because after React Context API and hooks being released) I suggest to stick with built-in state management.

**@russjones** on `zship/teleport-challenge` [→](https://github.com/zship/teleport-challenge/pull/1#discussion_r719878749)

> Doesn't the login endpoint change state?

**@alex-kovoy** on `zship/teleport-challenge` [→](https://github.com/zship/teleport-challenge/pull/1#discussion_r719887949)

> I was looking at something like the following 
> 
> ```
> /dir/../../../`
> ```
> 
>  or say a couple of words about encoding/decoding (if you are not planning to use react router for that)

**@alex-kovoy** on `zship/teleport-challenge` [→](https://github.com/zship/teleport-challenge/pull/1#discussion_r719913213)

> as per service requirements (if i am not mistaken), a read directory needs to be used of JSON.

**@kimlisa** on `zship/teleport-challenge` [→](https://github.com/zship/teleport-challenge/pull/4#discussion_r725399414)

> will there ever be a case where given a custom `code`, the http `statusCode` can be different? if not, i think it would be better to send in the custom `code` only and then assign the appropriate http code

**@russjones** on `tedmist1/challenge-auth0` [→](https://github.com/tedmist1/challenge-auth0/pull/1#discussion_r2488248176)

> Interesting, I didn't even realize that PowerShell works on GHA, but it does!
> 
> However, do you mind writing a Bash script instead? We don't use PowerShell at Teleport.

**@russjones** on `tedmist1/challenge-auth0` [→](https://github.com/tedmist1/challenge-auth0/pull/1#discussion_r2488253755)

> How is API access locked down?

**@sclevine** on `Chili-Man/teleport-sre-challenge` [→](https://github.com/Chili-Man/teleport-sre-challenge/pull/1#discussion_r1416180663)

> Here's a hint that we just merged: https://github.com/gravitational/careers/pull/92/files#diff-9ba8060c44ad19c05beb8348cb46de0e94a77e927a44fd3351f258e37f057924R92-R94

---

## code-style-nits

_Naming conventions; formatting; idiomatic Go (defer/error flow, package placement, type inference)._

**139 quotes** from `23` distinct reviewers across `34` candidate submissions.

**@codingllama** on `adalton/teleport-exercise` [→](https://github.com/adalton/teleport-exercise/pull/1#discussion_r767997338)

> nit: The more idiomatic Go way, IMO, is to move the Job interface closer to where it's used (Manager type), return concrete types from New* methods (concreteJob and mockJob).
> 
> https://github.com/golang/go/wiki/CodeReviewComments#interfaces

**@codingllama** on `adalton/teleport-exercise` [→](https://github.com/adalton/teleport-exercise/pull/1#discussion_r768005039)

> nit: Go conventions
> 
> ```suggestion
> 	jobsByUserByJobID   map[string]map[string]Job // userId->jobId->job
> 	jobsByUserByJobName map[string]map[string]Job // userId->jobName->job
> 	allJobsByJobID      map[string]Job            // jobId->job
> ```
> 
> https://github.com/golang/go/wiki/CodeReviewComments#initialisms

**@codingllama** on `adalton/teleport-exercise` [→](https://github.com/adalton/teleport-exercise/pull/1#discussion_r768009426)

> nit: Ditch the else?
> 
> https://github.com/golang/go/wiki/CodeReviewComments#indent-error-flow

**@codingllama** on `adalton/teleport-exercise` [→](https://github.com/adalton/teleport-exercise/pull/1#discussion_r768700334)

> You can move the logic to the cgroups package without changing the existing flow. This is more of a "which package should have this responsibility" question - the "when" of the setup is good.

**@codingllama** on `adalton/teleport-exercise` [→](https://github.com/adalton/teleport-exercise/pull/2#discussion_r771466737)

> nit: it's more usual, in Go, to use guard-like returns so we keep the "happy" path in minimal indentation.
> 
> https://github.com/golang/go/wiki/CodeReviewComments#indent-error-flow

**@codingllama** on `adalton/teleport-exercise` [→](https://github.com/adalton/teleport-exercise/pull/2#discussion_r772442974)

> FYI: Proto enums are package/namespace-scoped when compiled to certain languages, so typical style is to avoid overly generic names like "UNSET" and go for more verbose options like `MY_ENUM_NAME_UNSPECIFIED = 0`.

**@dboslee** on `benmoss/job-worker-service` [→](https://github.com/benmoss/job-worker-service/pull/3#discussion_r2688679545)

> I think a comment documenting the behavior would be good enough

**@sclevine** on `bkneis/jobworker` [→](https://github.com/bkneis/jobworker/pull/1#discussion_r1648302078)

> ```suggestion
> 	CPUWeight int    // `cpu.weight`
> ```
> https://go.dev/wiki/CodeReviewComments#initialisms

**@sclevine** on `bkneis/jobworker` [→](https://github.com/bkneis/jobworker/pull/2#discussion_r1657680581)

> https://go.dev/wiki/CodeReviewComments#mixed-caps

**@tcsc** on `bkneis/jobworker` [→](https://github.com/bkneis/jobworker/pull/2#discussion_r1663364790)

> See comment at `Job.Output()` implementation site.

**@zmb3** on `bucknercd/jobworker` [→](https://github.com/bucknercd/jobworker/pull/1#discussion_r2271496386)

> Does `true` mean:
> a) stdout _and_ stderr, or
> b) stderr only
>
> Depending on your answer, maybe either renaming the field or using an enum instead of a bool could have made it more clear.

**@zmb3** on `bucknercd/jobworker` [→](https://github.com/bucknercd/jobworker/pull/1#discussion_r2271523699)

> This seems lower-level than necessary. The syscall package docs even say:
>
> >  The primary use of syscall is inside other packages that provide a more portable interface to the system, such as "os", "time" and "net". Use those packages rather than this one if you can. 
>
>
> Can we use `os/exec` instead?

**@rosstimothy** on `bucknercd/jobworker` [→](https://github.com/bucknercd/jobworker/pull/1#discussion_r2273695472)

> You may use a UUID package if you wish instead of hand rolling your identifier.

**@strideynet** on `bucknercd/jobworker` [→](https://github.com/bucknercd/jobworker/pull/1#discussion_r2276112724)

> Nit: Should STDOUT and STDERR be prefixed with the enum name as well as you have done in JobStatus?

**@espadolini** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/1#discussion_r1782736231)

> See comment below about /proc

**@codingllama** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/2#discussion_r1787742328)

> nit: Document lack of lis.Close() (or close it anyway).
> 
> ```suggestion
> 		}
> 		// s.Serve takes ownership of lis.
> ```

**@codingllama** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/4#discussion_r1790702154)

> nit/Suggestion: indent error flow.
> 
> - https://go.dev/wiki/CodeReviewComments#indent-error-flow

**@codingllama** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/5#discussion_r1790891450)

> Thanks, an updated comment does help.

**@codingllama** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/5#discussion_r1795613162)

> nit:
> 
> ```suggestion
> 		return ctx.Err()
> ```

**@codingllama** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/6#discussion_r1795677263)

> nit: godoc, particularly for args as it's always a toss what args[0] is.

**@codingllama** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/6#discussion_r1795800869)

> nit: Unnecessary slog.With.

**@codingllama** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/6#discussion_r1795821945)

> nit: Move the cleanup right after StartJob?

**@awly** on `diptadas/job-worker` [→](https://github.com/diptadas/job-worker/pull/2#discussion_r648571547)

> nit: could you run `go mod tidy`?
> this `go.sum` seems to have a lot of stuff not actually used from `go.mod`

**@Tener** on `GevorgGal/jobworker` [→](https://github.com/GevorgGal/jobworker/pull/2#discussion_r2980008631)

> nit: we could probably drop offset and use `len(collected)` instead; both styles are fine, however.

**@sclevine** on `ilyazz/jobs` [→](https://github.com/ilyazz/jobs/pull/3#discussion_r1007302130)

> *nit*: https://github.com/golang/go/wiki/CodeReviewComments#indent-error-flow

**@codingllama** on `joshuarubin/teleport-job-worker` [→](https://github.com/joshuarubin/teleport-job-worker/pull/3#discussion_r1746054940)

> If you are looking for existing idioms I would refer to net.ErrClosed or http.ErrServerClosed, both of which are used to indicate interruption of some operation by an extraneous Close-like call.
> 
> This is only a nit comment though, so don't waste too much effort on it.

**@tigrato** on `joshuarubin/teleport-job-worker` [→](https://github.com/joshuarubin/teleport-job-worker/pull/4#discussion_r1746706012)

> Is it required to add a prefix? Aren't generally available UUIDs good enough?
> 
> 
> also add comments/godoc

**@codingllama** on `joshuarubin/teleport-job-worker` [→](https://github.com/joshuarubin/teleport-job-worker/pull/4#discussion_r1747504255)

> A small nit: config.Config stutters.

**@codingllama** on `joshuarubin/teleport-job-worker` [→](https://github.com/joshuarubin/teleport-job-worker/pull/4#discussion_r1747513909)

> A few comments here:
> 
> 1. There seems to be a lot in config.Config that Jobs don't care about, should we better select what we take in?
> 2. The way we capture the config lets it be changed externally, meaning it could change during the lifetime of the job. Is that expected?

**@codingllama** on `joshuarubin/teleport-job-worker` [→](https://github.com/joshuarubin/teleport-job-worker/pull/4#discussion_r1750908690)

> Nice, I frequently ask for this comment.

**@smallinsky** on `kkloberdanz/teleport-challenge` [→](https://github.com/kkloberdanz/teleport-challenge/pull/7#discussion_r2812065969)

> nit: 
> ```suggestion
> type Status int
> ```
>
> To make this more Go-idomatic by  avoiding repeating the package name in the type.

**@smallinsky** on `kkloberdanz/teleport-challenge` [→](https://github.com/kkloberdanz/teleport-challenge/pull/8#discussion_r2833180672)

> nit: I would suggest to moving the role  names to const

**@smallinsky** on `kkloberdanz/teleport-challenge` [→](https://github.com/kkloberdanz/teleport-challenge/pull/8#discussion_r2833232832)

> Nit: This probably should be just a `func (s Status) String` method  in job package.

**@zmb3** on `kovyrin/teleport-exec` [→](https://github.com/kovyrin/teleport-exec/pull/4#discussion_r794927137)

> nit: `_` is not idiomatic in package names

**@zmb3** on `kovyrin/teleport-exec` [→](https://github.com/kovyrin/teleport-exec/pull/2#discussion_r795810383)

> nit: prefer `filepath.Join()` when constructing filesystem paths.

**@nklaassen** on `kovyrin/teleport-exec` [→](https://github.com/kovyrin/teleport-exec/pull/2#discussion_r796198518)

> These comments are not very helpful. It's common to have godoc comments on exported types and methods https://go.dev/blog/godoc

**@sclevine** on `kurczynski/teleport-job-worker` [→](https://github.com/kurczynski/teleport-job-worker/pull/1#discussion_r1610845611)

> Would it make sense to give these names that better indicate how they apply controls?
> 
> Could you add a comment documenting each field?

**@sclevine** on `kurczynski/teleport-job-worker` [→](https://github.com/kurczynski/teleport-job-worker/pull/1#discussion_r1612586674)

> ```suggestion
> func (j *Job) ID() string
> ```
> https://go.dev/wiki/CodeReviewComments#initialisms

**@greedy52** on `m3talsmith/jobberthehut` [→](https://github.com/m3talsmith/jobberthehut/pull/1#discussion_r2500051357)

> sorry i missed it in my first pass. so if i get this right i will be something like this ?
> - first: `1762450000`
> - second: `111762460000`
> - third: `211762460000`
> 
> any reasons for this format?

**@rosstimothy** on `MarkDHarris/teleport-systems-challenge` [→](https://github.com/MarkDHarris/teleport-systems-challenge/pull/2#discussion_r3035831188)

> So we can tell which value is unknown?
> ```suggestion
> 		return fmt.Sprintf("UNKNOWN[%d]", s)
> ```

**@rosstimothy** on `MarkDHarris/teleport-systems-challenge` [→](https://github.com/MarkDHarris/teleport-systems-challenge/pull/2#discussion_r3039702699)

> If all branches of this function can only return non-nil errors, then do we need to return an error at all?
> ```suggestion
> // Cancel attempts to stop the job
> func (j *Job) Cancel() {
> ```

**@greedy52** on `MarkDHarris/teleport-systems-challenge` [→](https://github.com/MarkDHarris/teleport-systems-challenge/pull/2#discussion_r3041508823)

> minor suggestion, no need to check again. 
> ```suggestion
> ```

**@rosstimothy** on `MarkDHarris/teleport-systems-challenge` [→](https://github.com/MarkDHarris/teleport-systems-challenge/pull/2#discussion_r3045097508)

> ```suggestion
> 	_, err := b.Write([]byte("line 1\n"))
> 	if err != nil {
> ```

**@tigrato** on `mcampo84/teleport_challenge` [→](https://github.com/mcampo84/teleport_challenge/pull/5#discussion_r1920269666)

> why doesn't NewJob defines the UIID if we are passing it to the initializer. 
> If we don't need it in the initializer just skip sending it to Job and leave it only in JobManager

**@rosstimothy** on `mikewurtz/taskman` [→](https://github.com/mikewurtz/taskman/pull/1#discussion_r2028825245)

> Suggestion: reword/rename this since any binary in the PATH should be runnable as well.

**@rosstimothy** on `mikewurtz/taskman` [→](https://github.com/mikewurtz/taskman/pull/2#discussion_r2034142037)

> Suggestion: collapse the indirect blocks into a single block.
> ```suggestion
> 	gopkg.in/yaml.v3 v3.0.1 // indirect
> ```

**@rosstimothy** on `mikewurtz/taskman` [→](https://github.com/mikewurtz/taskman/pull/2#discussion_r2034143173)

> Suggestion: reduce the scope of the error.
> ```suggestion
> 			if err := cmd.Usage();  err != nil {
> ```

**@rosstimothy** on `mikewurtz/taskman` [→](https://github.com/mikewurtz/taskman/pull/2#discussion_r2034144943)

> Suggestion: reduce the error scope.
> ```suggestion
> 			if err := cmd.Usage(); err != nil {
> ```

**@rosstimothy** on `mikewurtz/taskman` [→](https://github.com/mikewurtz/taskman/pull/2#discussion_r2034145221)

> Suggestion: reduce the scope of the error.
> ```suggestion
> 			if err := cmd.Usage(); err != nil {
> ```

**@rosstimothy** on `mikewurtz/taskman` [→](https://github.com/mikewurtz/taskman/pull/2#discussion_r2034145520)

> Suggestion: reduce the scope of the error.
> ```suggestion
> 			if err := cmd.Usage(); err != nil {
> ```

**@zmb3** on `mikewurtz/taskman` [→](https://github.com/mikewurtz/taskman/pull/2#discussion_r2036264538)

> Clients will refer to this function as `server.NewServer`. WDYT about `server.New` instead?
> 
> (Same comment applies to `client.NewClient`)

**@rosstimothy** on `mikewurtz/taskman` [→](https://github.com/mikewurtz/taskman/pull/2#discussion_r2037393050)

> Suggestion: return the stop function instead of relying on the global `stopServer`.

**@tigrato** on `mikewurtz/taskman` [→](https://github.com/mikewurtz/taskman/pull/4#discussion_r2040105574)

> ```suggestion
> 	errC := make(chan error, 1)
> ```

**@rosstimothy** on `mikewurtz/taskman` [→](https://github.com/mikewurtz/taskman/pull/5#discussion_r2044570751)

> Suggestion: prefer errors.Is
> ```suggestion
> 		if errors.Is(err, io.EOF) {
> ```

**@rosstimothy** on `mikewurtz/taskman` [→](https://github.com/mikewurtz/taskman/pull/5#discussion_r2044583066)

> Suggestion: prefer errors.Is
> ```suggestion
> 		if errors.Is(err, io.EOF) {
> ```

**@timothyb89** on `mikewurtz/taskman` [→](https://github.com/mikewurtz/taskman/pull/5#discussion_r2045534876)

> nit: looks like there's an unhandled error here
> 
> (there's also one up a bit that I can't leave a comment on)

**@timothyb89** on `mikewurtz/taskman` [→](https://github.com/mikewurtz/taskman/pull/5#discussion_r2047839510)

> one final nit for the road, but looks like another unhandled `error` here

**@oeric** on `moalf/teleport` [→](https://github.com/moalf/teleport/pull/1#discussion_r2594165486)

> How would users sign on initially after their account is added to the tenant?

**@nklaassen** on `MrChristianL/Teleport-Job-Worker-Service` [→](https://github.com/MrChristianL/Teleport-Job-Worker-Service/pull/1#discussion_r2896987987)

> nit: the `http://` prefix seems unnecessary/confusing

**@eriktate** on `MrChristianL/Teleport-Job-Worker-Service` [→](https://github.com/MrChristianL/Teleport-Job-Worker-Service/pull/2#discussion_r2905854136)

> ```suggestion
> 	for i := range numClients {
> ```

**@eriktate** on `MrChristianL/Teleport-Job-Worker-Service` [→](https://github.com/MrChristianL/Teleport-Job-Worker-Service/pull/2#discussion_r2906068345)

> ```suggestion
> 		return nil, errors.New("job ID cannot be empty")
> ```

**@eriktate** on `MrChristianL/Teleport-Job-Worker-Service` [→](https://github.com/MrChristianL/Teleport-Job-Worker-Service/pull/2#discussion_r2906070962)

> ```suggestion
> 		return nil, errors.New("command cannot be empty")
> ```

**@nklaassen** on `MrChristianL/Teleport-Job-Worker-Service` [→](https://github.com/MrChristianL/Teleport-Job-Worker-Service/pull/2#discussion_r2906737392)

> suggestion: return on non-EOF error, continue the loop on err == nil, and unindent the rest of this branch (bring it out of the if)

**@Joerger** on `MrChristianL/Teleport-Job-Worker-Service` [→](https://github.com/MrChristianL/Teleport-Job-Worker-Service/pull/2#discussion_r2915166733)

> Not used anymore outside of tests
> ```suggestion
> ```

**@Joerger** on `MrChristianL/Teleport-Job-Worker-Service` [→](https://github.com/MrChristianL/Teleport-Job-Worker-Service/pull/2#discussion_r2915204722)

> suggestion: feel free to ignore
> ```suggestion
> 			// if broker is not closed, but we have read all avialable content, reader waits
> 			for b.totalWritten == offset {
> 				if ctx.Err() != nil {
> 					b.mu.Unlock()
> 					return ctx.Err()
> 				}
> 				if b.closed {
> 					// if broker is closed and all data has been read, exit
> 					b.mu.Unlock()
> 				    return nil
> 				}
> 				b.cond.Wait()
> 			}
>
> 			b.mu.Unlock()
> ```

**@rosstimothy** on `MrChristianL/Teleport-Job-Worker-Service` [→](https://github.com/MrChristianL/Teleport-Job-Worker-Service/pull/2#discussion_r2919439872)

> Suggestion: Invert the condition and continue so the rest of this logic can be unindented.

**@rosstimothy** on `MrChristianL/Teleport-Job-Worker-Service` [→](https://github.com/MrChristianL/Teleport-Job-Worker-Service/pull/2#discussion_r2919474089)

> Use [fmt.Fprintf](https://pkg.go.dev/fmt#Fprintf) instead?
> ```suggestion
> 		fmt.Fprintf(broker, "chunk-%d", i)
> ```

**@nklaassen** on `MrChristianL/Teleport-Job-Worker-Service` [→](https://github.com/MrChristianL/Teleport-Job-Worker-Service/pull/2#discussion_r2927432100)

> mentioned in another comment, but I think this could access j.cmd before it has been set in Start

**@rosstimothy** on `neildo/tjob` [→](https://github.com/neildo/tjob/pull/1#discussion_r1764112558)

> minor nit: fields should be marked as `optional`, not `option`
> ```suggestion
>   optional string job_id = 1;
> ```

**@rosstimothy** on `neildo/tjob` [→](https://github.com/neildo/tjob/pull/2#discussion_r1769133232)

> Suggestion: use t.TempDir instead.

**@greedy52** on `P-A-R-U-S/Go-Job-Worker-Service` [→](https://github.com/P-A-R-U-S/Go-Job-Worker-Service/pull/1#discussion_r1731338332)

> nit: `Expected behavior User2 should receive error.` feels like a negative case?

**@smallinsky** on `P-A-R-U-S/Go-Job-Worker-Service` [→](https://github.com/P-A-R-U-S/Go-Job-Worker-Service/pull/1#discussion_r1732274298)

> Does the use need to know the id format ? 
> 
> ```suggestion
>     string  id = 1;
> ```

**@smallinsky** on `P-A-R-U-S/Go-Job-Worker-Service` [→](https://github.com/P-A-R-U-S/Go-Job-Worker-Service/pull/2#discussion_r1745060619)

> nit: https://google.github.io/styleguide/go/decisions#constant-names

**@eriktate** on `rajansandeep/teleport-job-worker-svc` [→](https://github.com/rajansandeep/job-worker-svc/pull/1#discussion_r3088901371)

> nit: will these ever be significantly different enough to warrant two fields?

**@eriktate** on `rajansandeep/teleport-job-worker-svc` [→](https://github.com/rajansandeep/job-worker-svc/pull/2#discussion_r3110929642)

> nit: `ErrMsg` was originally `Error` in the design. No preference for either name, but consider making these consistent.

**@atburke** on `rajansandeep/teleport-job-worker-svc` [→](https://github.com/rajansandeep/job-worker-svc/pull/3#discussion_r3112329977)

> Nit: Check out [strings.LastIndex](https://pkg.go.dev/strings#LastIndex)

**@atburke** on `rajansandeep/teleport-job-worker-svc` [→](https://github.com/rajansandeep/job-worker-svc/pull/3#discussion_r3112486009)

> Nit: we have this form too :)
> ```suggestion
> 	for i := range iterations {
> ```

**@timothyb89** on `razzam21/job-worker-service` [→](https://github.com/razzam21/job-worker-service/pull/1#discussion_r2261143622)

> I know it's probably splitting hairs but this definition still feels a bit odd? The app shouldn't be examining bytes at all so bringing up byte boundaries in the first place feels a bit misleading.

**@awly** on `renatoaguimaraes/golang-job-scheduler` [→](https://github.com/renatoaguimaraes/golang-job-scheduler/pull/2#discussion_r623310774)

> nit: why not export the `logFolder` field in the struct?

**@awly** on `renatoaguimaraes/golang-job-scheduler` [→](https://github.com/renatoaguimaraes/golang-job-scheduler/pull/2#discussion_r623311997)

> nit: don't use named return values unless necessary

**@awly** on `renatoaguimaraes/golang-job-scheduler` [→](https://github.com/renatoaguimaraes/golang-job-scheduler/pull/2#discussion_r623312685)

> nit: rename `extctx` to `ctx`

**@awly** on `renatoaguimaraes/golang-job-scheduler` [→](https://github.com/renatoaguimaraes/golang-job-scheduler/pull/2#discussion_r623313504)

> ```suggestion
> 		cancel()
> 		file.Close()
> ```

**@awly** on `renatoaguimaraes/golang-job-scheduler` [→](https://github.com/renatoaguimaraes/golang-job-scheduler/pull/2#discussion_r623377340)

> ```suggestion
> 				return
> ```

**@awly** on `renatoaguimaraes/golang-job-scheduler` [→](https://github.com/renatoaguimaraes/golang-job-scheduler/pull/2#discussion_r623377948)

> is this more accurate?
> ```suggestion
> 			for offset < uint32(nbytes) {
> ```

**@awly** on `renatoaguimaraes/golang-job-scheduler` [→](https://github.com/renatoaguimaraes/golang-job-scheduler/pull/2#discussion_r623381753)

> same as above, don't use naked return values and prefer not to use named return parameters

**@awly** on `renatoaguimaraes/golang-job-scheduler` [→](https://github.com/renatoaguimaraes/golang-job-scheduler/pull/2#discussion_r623382474)

> ```suggestion
> 	if err = cmd.Start(); err != nil {
> 		w.logger.Remove(jobID)
> ```

**@awly** on `renatoaguimaraes/golang-job-scheduler` [→](https://github.com/renatoaguimaraes/golang-job-scheduler/pull/2#discussion_r623385783)

> ```suggestion
> 	assert.Equal(t, "Job not found", err.Error())
> ```

**@dboslee** on `renatoaguimaraes/golang-job-scheduler` [→](https://github.com/renatoaguimaraes/golang-job-scheduler/pull/2#discussion_r623401151)

> Looks like I got beat to submitting my review feel free to mark any redundant comments as resolved.

**@Joerger** on `renatoaguimaraes/golang-job-scheduler` [→](https://github.com/renatoaguimaraes/golang-job-scheduler/pull/2#discussion_r623530632)

> ```suggestion
> 	// Start returns
> ```

**@r0mant** on `renatoaguimaraes/golang-job-scheduler` [→](https://github.com/renatoaguimaraes/golang-job-scheduler/pull/2#discussion_r623531098)

> Information about this error is lost.

**@r0mant** on `renatoaguimaraes/golang-job-scheduler` [→](https://github.com/renatoaguimaraes/golang-job-scheduler/pull/2#discussion_r623531328)

> This error information is lost.

**@awly** on `renatoaguimaraes/golang-job-scheduler` [→](https://github.com/renatoaguimaraes/golang-job-scheduler/pull/2#discussion_r625344755)

> ```suggestion
> 	if err != nil {
> ```

**@awly** on `renatoaguimaraes/golang-job-scheduler` [→](https://github.com/renatoaguimaraes/golang-job-scheduler/pull/6#discussion_r625939738)

> nit: this timeout is very aggressive, would cause problems when connecting to a server over a long distance or slow link

**@awly** on `renatoaguimaraes/golang-job-scheduler` [→](https://github.com/renatoaguimaraes/golang-job-scheduler/pull/5#discussion_r625958967)

> this is likely because the job finished, not an error
> ```suggestion
> 				return nil
> ```

**@awly** on `renatoaguimaraes/golang-job-scheduler` [→](https://github.com/renatoaguimaraes/golang-job-scheduler/pull/5#discussion_r625963335)

> ```suggestion
> 	}
> 	defer lis.Close()
> ```

**@awly** on `renatoaguimaraes/golang-job-scheduler` [→](https://github.com/renatoaguimaraes/golang-job-scheduler/pull/5#discussion_r625968538)

> nit: all of this server setup code can be extracted into a helper func and reused in other tests

**@dboslee** on `renatoaguimaraes/golang-job-scheduler` [→](https://github.com/renatoaguimaraes/golang-job-scheduler/pull/6#discussion_r626148295)

> ```suggestion
> ```
> No need for a return after calling os.Exit

**@strideynet** on `rexposadas/teleport` [→](https://github.com/rexposadas/teleport/pull/1#discussion_r1793031241)

> ```suggestion
> ### TLS
> ```

**@rosstimothy** on `rexposadas/teleport` [→](https://github.com/rexposadas/teleport/pull/1#discussion_r1793930692)

> ```suggestion
> #### Trade offs and the Production environment
> ```

**@rosstimothy** on `rexposadas/teleport` [→](https://github.com/rexposadas/teleport/pull/1#discussion_r1795658708)

> An ownership model like that would be fine. Please include this information in the design.

**@greedy52** on `rexposadas/teleport` [→](https://github.com/rexposadas/teleport/pull/2#discussion_r1803056847)

> ```suggestion
> 		Args:  cobra.ExactArgs(1),
> ```

**@rosstimothy** on `RichyHBM/teleport-challenge` [→](https://github.com/RichyHBM/teleport-challenge/pull/1#discussion_r2049502012)

> Which device major:minor will be in scope? All devices? A particular device?

**@rosstimothy** on `RichyHBM/teleport-challenge` [→](https://github.com/RichyHBM/teleport-challenge/pull/1#discussion_r2050613910)

> Sounds reasonable. Please update the design to include this information.

**@rosstimothy** on `RichyHBM/teleport-challenge` [→](https://github.com/RichyHBM/teleport-challenge/pull/2#discussion_r2056050127)

> Suggestion: Write directly to stdout instead of converting to a string

**@rosstimothy** on `RichyHBM/teleport-challenge` [→](https://github.com/RichyHBM/teleport-challenge/pull/2#discussion_r2056058593)

> Suggestion: defer to the `UnimplementedJobsServiceServer` and return a not implemented error instead of responses with a not implemented error message.

**@rosstimothy** on `RichyHBM/teleport-challenge` [→](https://github.com/RichyHBM/teleport-challenge/pull/3#discussion_r2060433810)

> Suggestion: use filepath.Join here and all other places paths are being concatenated

**@rosstimothy** on `RichyHBM/teleport-challenge` [→](https://github.com/RichyHBM/teleport-challenge/pull/3#discussion_r2060448506)

> Suggestion: drop the else clause an unindent this block of code. The error is already handled and the function will return early if needed.

**@rosstimothy** on `RichyHBM/teleport-challenge` [→](https://github.com/RichyHBM/teleport-challenge/pull/3#discussion_r2060449970)

> Suggestion: drop the else clause and reduce indentation. The error is already handled above and the function will return early if needed.

**@rosstimothy** on `RichyHBM/teleport-challenge` [→](https://github.com/RichyHBM/teleport-challenge/pull/4#discussion_r2064069999)

> Why is this function needed? Can't the two call sites below interact with authData directly?
> ```suggestion
> ```

**@rosstimothy** on `RichyHBM/teleport-challenge` [→](https://github.com/RichyHBM/teleport-challenge/pull/4#discussion_r2064071698)

> Suggestion: consider using a map of sets instead of a map of slices
> ```suggestion
> 	authData  map[string]map[string]struct{} = map[string]map[string]struct{
> ```

**@rosstimothy** on `RichyHBM/teleport-challenge` [→](https://github.com/RichyHBM/teleport-challenge/pull/5#discussion_r2064108982)

> Same as previous comments: Prefer dropping the else clause when the previous condition is handled by returning early.

**@rosstimothy** on `RichyHBM/teleport-challenge` [→](https://github.com/RichyHBM/teleport-challenge/pull/5#discussion_r2064110814)

> Same as previous comments: Prefer dropping the else clause when the previous condition is handled by returning early.

**@rosstimothy** on `RichyHBM/teleport-challenge` [→](https://github.com/RichyHBM/teleport-challenge/pull/4#discussion_r2066178584)

> The suggestion was to optimize lookups of the current data more so than a suggestion to change what was being stored.

**@zmb3** on `s-gruneberg/jobWorker` [→](https://github.com/s-gruneberg/jobWorker/pull/1#discussion_r2305423039)

> ```suggestion
>         - Will use Google's UUID package for generating job IDs 
> ```

**@rosstimothy** on `shawon-crosen/teleport-challenge` [→](https://github.com/shawon-crosen/teleport-challenge/pull/1#discussion_r2126985823)

> Where will this information be included?

**@smallinsky** on `supby/job-worker` [→](https://github.com/supby/job-worker/pull/1#discussion_r715437433)

> ```suggestion
>     bytes output = 1;
> ```
> 
> Not all command output contain valid UTF or ASCII for intense terminal control sequences thus byte type seems to  be more suitable here.

**@smallinsky** on `supby/job-worker` [→](https://github.com/supby/job-worker/pull/1#discussion_r715451600)

> I think that pid is redundant information from a client perspective.

**@nklaassen** on `supby/job-worker` [→](https://github.com/supby/job-worker/pull/1#discussion_r715929834)

> spelling nit:
> ```suggestion
> ### Authentication
> ```

**@rosstimothy** on `tjper/teleport` [→](https://github.com/tjper/teleport/pull/2#discussion_r826953177)

> ```suggestion
> 		var stats unix.Stat_t
> ```

**@codingllama** on `tjper/teleport` [→](https://github.com/tjper/teleport/pull/2#discussion_r827463951)

> nit: Let's not rely on private functions during tests.

**@codingllama** on `tjper/teleport` [→](https://github.com/tjper/teleport/pull/2#discussion_r828085246)

> It may be the old man in me, but I prefer avoiding goto unless the readability gains are _very_ significant.

**@codingllama** on `tjper/teleport` [→](https://github.com/tjper/teleport/pull/3#discussion_r829289100)

> nit: cmdFD, contFD, so on.

**@codingllama** on `tjper/teleport` [→](https://github.com/tjper/teleport/pull/3#discussion_r829310119)

> nit: job.Job stutters. Not hugely important.

**@zmb3** on `tjper/teleport` [→](https://github.com/tjper/teleport/pull/3#discussion_r830293477)

> ```suggestion
> 		_, _ = fmt.Fprintf(b, "\nNotice: %s", text) 
> ```
> 
> `strings.Builder` implements `io.Writer` so this can be simplified.

**@rcanderson23** on `vincenthlam/teleport-interview` [→](https://github.com/vincenthlam/teleport-interview/pull/1#discussion_r1611929548)

> ```suggestion
> ## Controller - Worker(s) Configuration
> ```

**@rcanderson23** on `vincenthlam/teleport-interview` [→](https://github.com/vincenthlam/teleport-interview/pull/1#discussion_r1611931509)

> Looks like some of the markdown formatting is off here.

**@rcanderson23** on `vincenthlam/teleport-interview` [→](https://github.com/vincenthlam/teleport-interview/pull/1#discussion_r1611959403)

> Do we need information here regarding the amount of resources required for the job?

**@greedy52** on `xSudoNymx/job-manager` [→](https://github.com/xSudoNymx/job-manager/pull/3#discussion_r2727967033)

> minor suggestion: check length of `args` first before assignment

**@tigrato** on `xSudoNymx/job-manager` [→](https://github.com/xSudoNymx/job-manager/pull/3#discussion_r2728142914)

> ```suggestion
> 	return slices.Clone(j,args)
> ```

**@tigrato** on `xSudoNymx/job-manager` [→](https://github.com/xSudoNymx/job-manager/pull/4#discussion_r2738412074)

> ```suggestion
> 	if req.GetJobId() == "" {
> ```
>
> GetJobId checks for receiver nilness

**@gabrielcorado** on `zeeshanhaque21/JobWorkerService` [→](https://github.com/zeeshanhaque21/JobWorkerService/pull/1#discussion_r1692421371)

> nit: This message doesn't seem to be used anywhere.

**@nklaassen** on `Zephan92/teleport-handson` [→](https://github.com/Zephan92/teleport-handson/pull/1#discussion_r2719082176)

> nit: the timestamp feels out of place here, this is primary process output not a log

**@nklaassen** on `Zephan92/teleport-handson` [→](https://github.com/Zephan92/teleport-handson/pull/1#discussion_r2719141766)

> no need for a TTL, jobs can run indefinitely

**@espadolini** on `Zephan92/teleport-handson` [→](https://github.com/Zephan92/teleport-handson/pull/3#discussion_r2736619900)

> Please use structured logging with `log/slog` rather than string formatting.

**@espadolini** on `Zephan92/teleport-handson` [→](https://github.com/Zephan92/teleport-handson/pull/3#discussion_r2736649340)

> ```suggestion
> 	id := JobID(uuid.NewString())
> ```

**@nklaassen** on `Zephan92/teleport-handson` [→](https://github.com/Zephan92/teleport-handson/pull/6#discussion_r2761639056)

> this format string is approaching a valid usecase for text/template. It also doesn't really help to have it defined so far away (in a different file) from where it is used

**@nklaassen** on `Zephan92/teleport-handson` [→](https://github.com/Zephan92/teleport-handson/pull/6#discussion_r2761642426)

> nit: I don't see the benefit of using constants in another file for these strings, it kind of just makes this harder to read

**@espadolini** on `Zephan92/teleport-handson` [→](https://github.com/Zephan92/teleport-handson/pull/5#discussion_r2765811199)

> typo
>
> ```suggestion
> ```

**@alex-kovoy** on `ibeckermayer/teleport-interview` [→](https://github.com/ibeckermayer/teleport-interview/pull/2#discussion_r552398258)

> kudos for not using Create React App boilerplate :) (there is nothing wrong with it aside the fact that it uses a lot of dependencies this is why I also prefer a minimalistic custom setup when possible)

---

## process-containment

_Managing process groups, cgroups, pgid, pdeathsig, CLONE_INTO_CGROUP; ensuring spawned processes are properly contained and reaped._

**133 quotes** from `18` distinct reviewers across `38` candidate submissions.

**@rosstimothy** on `adalton/teleport-exercise` [→](https://github.com/adalton/teleport-exercise/pull/1#discussion_r767872569)

> Will this kill all spawned processes?

**@rosstimothy** on `adalton/teleport-exercise` [→](https://github.com/adalton/teleport-exercise/pull/1#discussion_r767893785)

> ```suggestion
> 	if err := command.Cgexec(os.Args); err != nil {
> ```

**@codingllama** on `adalton/teleport-exercise` [→](https://github.com/adalton/teleport-exercise/pull/1#discussion_r767953467)

> ```suggestion
> You can build the `cgexec` binary using `make cgexec`.  The resulting binary
> ```

**@codingllama** on `adalton/teleport-exercise` [→](https://github.com/adalton/teleport-exercise/pull/1#discussion_r768025777)

> The directory name is `pkg/group/v1`, but the package name `cgroup` - this is easy to escape and kind of looks like a Go Module version. It can also get confusing because it could be the cgroups v1 API or it could be the first version of the "cgroups" package.
> 
> For the sake of clarity / not making reviewers worry about Modules, I'd be tempted to rename the dir and package to `cgroupv1` (or something like it), to it 1) matches the dir and 2) looks like a "normal" package.

**@codingllama** on `adalton/teleport-exercise` [→](https://github.com/adalton/teleport-exercise/pull/1#discussion_r768042938)

> I wonder if we could push group joining to the cgroups package, instead of having the command responsible for it. Seems clearer in terms of design.
> 
> WDYT?

**@codingllama** on `adalton/teleport-exercise` [→](https://github.com/adalton/teleport-exercise/pull/1#discussion_r768047853)

> Check before writing cgroups?

**@rosstimothy** on `adalton/teleport-exercise` [→](https://github.com/adalton/teleport-exercise/pull/1#discussion_r769686101)

> This panics if Stop is called fast enough. `j.cmd.Process` will be nil until `exec.Cmd.Start` has returned.

**@codingllama** on `adalton/teleport-exercise` [→](https://github.com/adalton/teleport-exercise/pull/2#discussion_r771425060)

> Could we add the cgexec behavior to the server, so we don't depend on an external binary?
> 
> It seems like it's easy to miss this and have a server that doesn't really work. Let's also ditch the possible panic below, if we can.

**@sclevine** on `benmoss/job-worker-service` [→](https://github.com/benmoss/job-worker-service/pull/3#discussion_r2683955684)

> What about child processes?

**@smallinsky** on `bill-rich/jobworker` [→](https://github.com/bill-rich/jobworker/pull/1#discussion_r1606443013)

> Could you elaborate a bit more on job execution?
> 
> How is a job started, and when and how is isolation applied?

**@smallinsky** on `bill-rich/jobworker` [→](https://github.com/bill-rich/jobworker/pull/1#discussion_r1606450397)

> Does process will add itself to cgroup  `os.Getpid())` ? 
> 
> 
> Are you going to use cgroup v1 or cgroups v2 ?

**@rosstimothy** on `bill-rich/jobworker` [→](https://github.com/bill-rich/jobworker/pull/1#discussion_r1608425873)

> How is the process killed? Will this also terminate any possible child processes spawned by job?

**@rosstimothy** on `bill-rich/jobworker` [→](https://github.com/bill-rich/jobworker/pull/1#discussion_r1608429391)

> Can you provide some more detail about how you will leverage cgroups to restrict cpu, memory and disk i/o? What will your predefined resource limits be?

**@rosstimothy** on `bill-rich/jobworker` [→](https://github.com/bill-rich/jobworker/pull/1#discussion_r1608432222)

> Why do we need the distinction between start and run? The requirements for the challenge are that all jobs that are spawned by the server have resource control and isolation via cgroups and namespaces.

**@tcsc** on `bkneis/jobworker` [→](https://github.com/bkneis/jobworker/pull/1#discussion_r1650235375)

> Possibly naïve question: Will the `cgroup` (or any other) cleanup have any impact on clients asking for the command output long after the job has finished?

**@tcsc** on `bkneis/jobworker` [→](https://github.com/bkneis/jobworker/pull/1#discussion_r1650239657)

> Just clarifying my understanding of the cgroup boundary. 
> 
> Once everything is up and going (that is, ignoring the init sequence where the `job` adds itself to the `cgroup`), the process tree will look something like this will look like this, yes?
> 
> ```mermaid
> flowchart TD
>    JW[jobworker]
>    A[job]
>    B[job]
>    C[job]
> 
>    JW --- A
>    JW --- B
>    JW --- C
> 
>    subgraph cgroup boundary
>       A --- CMD[cmd]
>       CMD --- CAT[cat]
>       CMD --- PS[ps]
>    end
> ```

**@tcsc** on `bkneis/jobworker` [→](https://github.com/bkneis/jobworker/pull/1#discussion_r1650240631)

> Have you thought about how you might test that the cgroup limits are correctly enforced?

**@sclevine** on `bkneis/jobworker` [→](https://github.com/bkneis/jobworker/pull/1#discussion_r1651353340)

> Instead of assuming mega-units, could we accept the actual cgroup value? Then you could provide consts in Go to make the API easier to use (similar to `time.Second`).

**@sclevine** on `bkneis/jobworker` [→](https://github.com/bkneis/jobworker/pull/1#discussion_r1651759842)

> This adds a dependency on a shell. That's okay, but would be good to document this requirement carefully. 
> 
> Alternatively, can you think of a way to execve a process that lets you run additional code, without requiring a script or extra binary?

**@sclevine** on `bkneis/jobworker` [→](https://github.com/bkneis/jobworker/pull/1#discussion_r1653041608)

> That works in recent versions of Go, given that cgroups don't need much setup. Much cleaner solution 🙂 
> 
> Note that it's also possible to re-exec the current binary to avoid using a shell, for more complex use cases.

**@sclevine** on `bkneis/jobworker` [→](https://github.com/bkneis/jobworker/pull/2#discussion_r1663222019)

> Does that child process stop gracefully, even if the main process quits before it does?

**@tcsc** on `bkneis/jobworker` [→](https://github.com/bkneis/jobworker/pull/2#discussion_r1663383803)

> We only ever have one `exec.Cmd` per CGroup, yes? Or, to put it another way, is it worth doing anything smart to re-use cgroup FDs for multiple `Cmd`s?

**@tcsc** on `bkneis/jobworker` [→](https://github.com/bkneis/jobworker/pull/2#discussion_r1663386168)

> Is there any protection against deleting someone else's CGroup? Or is that handled at the job execution level?

**@tcsc** on `bkneis/jobworker` [→](https://github.com/bkneis/jobworker/pull/2#discussion_r1663390624)

> ~Can you add a comment explaining why is this a thing you would want to do?~
> 
> ```suggestion
> // updateController sets the content of the controller interface file for a
> // given resource controller within a CGroup (e.g. "memory.high", etc.)  
> ```

**@rosstimothy** on `bucknercd/jobworker` [→](https://github.com/bucknercd/jobworker/pull/1#discussion_r2266857437)

> > then starts a job with cgroup limits, or without
> 
> One of the requirements of the challenge is that _all_ jobs have limits enforced.

**@rosstimothy** on `bucknercd/jobworker` [→](https://github.com/bucknercd/jobworker/pull/1#discussion_r2266881961)

> What guarantees do we have that the child didn't change it's process group? Is there a more reliable way to terminate child processes?

**@zmb3** on `bucknercd/jobworker` [→](https://github.com/bucknercd/jobworker/pull/1#discussion_r2271521018)

> Feel free to skip straight to `cgroup.kill` in order to simplify the implementation.

**@rosstimothy** on `bucknercd/jobworker` [→](https://github.com/bucknercd/jobworker/pull/1#discussion_r2273693661)

> Is there any way to launch the child process directly in the cgroup? That seems like the most secure way to avoid any opportunity to escape the cgroup and enforced limits.

**@espadolini** on `chintamanil/job-worker` [→](https://github.com/chintamanil/job-worker/pull/1#discussion_r2822029721)

> You mention that cgroups are going to be used for process management and cleanup, if that's the case why are you relying on pgid signaling to kill processes?
>
> It's technically possible to signal the process group sanely, but that requires some extra hoops to avoid racing against waitid, so I would advise against it and use something that lets you take advantage of the cgroups instead.

**@espadolini** on `chintamanil/job-worker` [→](https://github.com/chintamanil/job-worker/pull/1#discussion_r2822066536)

> What's the point of a fallback to a much less reliable mode of operation? Unprivileged processes can setpgid, but they can't easily change cgroups.

**@espadolini** on `chintamanil/job-worker` [→](https://github.com/chintamanil/job-worker/pull/1#discussion_r2822079873)

> The challenge text doesn't mention efficiently listing jobs so what's the point of the per-user maps? I also don't think that custom job names would be necessary, so keeping a flat namespace with UUIDs assigned by the server would do just fine.

**@tigrato** on `chintamanil/job-worker` [→](https://github.com/chintamanil/job-worker/pull/1#discussion_r2822201086)

> To clarify, process isolation, such as PIDs, host, and network separation, is not part of this challenge
>
> My point was that it’s acceptable for the process to run without internet access (for example, in a VM with no external connectivity) rather than moving the process into a separate network namespace

**@tigrato** on `chintamanil/job-worker` [→](https://github.com/chintamanil/job-worker/pull/1#discussion_r2822222048)

> cgroups do not provide process tree isolation

**@tigrato** on `chintamanil/job-worker` [→](https://github.com/chintamanil/job-worker/pull/1#discussion_r2822235514)

> What you mean by escaping their cgroup?
>
> cgroup provide resource constrains and as long as the process has the necessary permissions (user), it can move itself to different cgroups or interact with other processes. It doesn't provide any isolation itself.
>
> Both cases are out of scope

**@strideynet** on `clydotron/job_worker_service` [→](https://github.com/clydotron/job_worker_service/pull/1#discussion_r1600236665)

> Which version of cgroups will be used? Why? Which cgroup controllers will need to be used?

**@espadolini** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/1#discussion_r1782738055)

> Considering that the requirements don't mention user namespaces at all I'd say that this is fine 😅

**@espadolini** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/1#discussion_r1782742706)

> Remounting proc in the new mount namespace would fix this even without a new root mount, right?

**@codingllama** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/1#discussion_r1783307292)

> Which cgroups toggles do you intend to use?

**@rosstimothy** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/1#discussion_r1786193820)

> There is still little information on the process execution lifecycle. How is the process going to be placed in the cgroup? More importantly, when will the process be placed in the cgroup?

**@rosstimothy** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/1#discussion_r1786195051)

> I don't see any details on how resource isolation will be achieved.

**@espadolini** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/1#discussion_r1786405889)

> For the sake of easier cleanup I'd recommend putting the job cgroups in a subtree like `telepilot/<uuid>` or something along those lines.

**@codingllama** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/2#discussion_r1787763309)

> This exposes all the data and intricacies of exec.Cmd. Should we limit ourselves to the prescribed actions only?

**@espadolini** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/6#discussion_r1795060022)

> I'd move this to the same place where you assign to `CgroupFD` or you might end up trying to use stdin as a cgroupfd if this gets forgotten with a refactor.

**@codingllama** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/6#discussion_r1795793388)

> ```suggestion
> 	// Make sure the base cgroup exists.
> ```

**@codingllama** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/6#discussion_r1795816396)

> It's worth noting that, with "true" isolation, this wouldn't be a path on the job filesystem.

**@espadolini** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/6#discussion_r1796039299)

> I would just unconditionally print out the error and have it show up as an `exec: <nil>` rather than panicking if we got no error after returning from execve.

**@rosstimothy** on `devnulled/runjob` [→](https://github.com/devnulled/runjob/pull/2#discussion_r1797084072)

> What will the hard coded limits be? Can you include them here? Also which cgroup controllers will the limits be applied to?

**@rosstimothy** on `devnulled/runjob` [→](https://github.com/devnulled/runjob/pull/2#discussion_r1797090350)

> I don't see any mention of how and when the process will be placed in the cgroup below.

**@nklaassen** on `diptadas/job-worker` [→](https://github.com/diptadas/job-worker/pull/1#discussion_r645159099)

> How will you kill the process after timeout, and how can this fail?

**@Joerger** on `diptadas/job-worker` [→](https://github.com/diptadas/job-worker/pull/2#discussion_r647916948)

> This log will also occur when the job is killed. Is there a way to skip this log in favor of `log.Infof("job %v: terminated", job.ID)` in `StopJob`?

**@zmb3** on `DiscoRiver/jobber` [→](https://github.com/DiscoRiver/jobber/pull/1#discussion_r1989355661)

> - What cgroup controllers will these options correspond to?
> - We also have a requirement to limit disk I/O.

**@rosstimothy** on `DiscoRiver/jobber` [→](https://github.com/DiscoRiver/jobber/pull/1#discussion_r1989356467)

> Could you add some more detail to this section about _how_ groups will be used? Which controllers are being used to apply limits? How will a process be added to the cgroup?

**@tigrato** on `DiscoRiver/jobber` [→](https://github.com/DiscoRiver/jobber/pull/1#discussion_r1989469447)

> Start the command inside which cgroup?
> The process adds itself to the cgroup.procs file when it starts? This section is not clear to me. 
> 
> Can you please include more details?

**@zmb3** on `DiscoRiver/jobber` [→](https://github.com/DiscoRiver/jobber/pull/1#discussion_r1990185228)

> I preferred the original approach of starting the job directly in a cgroup.
> 
> This approach of starting a job and then moving it into the cgroup at some point in the future leaves a window of time where the job is not subject to resource controls.

**@zmb3** on `DiscoRiver/jobber` [→](https://github.com/DiscoRiver/jobber/pull/1#discussion_r2063511870)

> If you start the command and then at some later time move the process into the cgroup, doesn't that create a window of time where the process runs unconstrained? Any issues with that?

**@rosstimothy** on `DiscoRiver/jobber` [→](https://github.com/DiscoRiver/jobber/pull/1#discussion_r2063520917)

> No need to adjust if you only plan on only supporting cgroups v2 - which is fine for this challenge.

**@zmb3** on `DiscoRiver/jobber` [→](https://github.com/DiscoRiver/jobber/pull/1#discussion_r2065159870)

> If you made some updates here I might be missing them. How will you start the process inside the cgroup?

**@russjones** on `dtom90/teleport-backend-challenge` [→](https://github.com/dtom90/teleport-backend-challenge/pull/1#discussion_r2053135418)

> Try not to use a external tool like `cgexec`, can you implement it yourself in Go?

**@greedy52** on `ehsaniara/joblet` [→](https://github.com/ehsaniara/joblet/pull/2#discussion_r2113891168)

> it's unclear if client is setting these cgroup limits or they are controlled/returned by the server. Same for `pid` and `cgroupPath`.

**@zmb3** on `ehsaniara/joblet` [→](https://github.com/ehsaniara/joblet/pull/2#discussion_r2121705379)

> I'm still seeing `pid` and `cgroupPath` in the proto spec.

**@greedy52** on `ehsaniara/joblet` [→](https://github.com/ehsaniara/joblet/pull/2#discussion_r2124462702)

> > my other idea wouldnt be in pure Go, like: to create a tiny init process that sets up environments and cgroups then the init process execs the real command
> 
> that works

**@russjones** on `flychicken123/Job_Worker` [→](https://github.com/flychicken123/Job_Worker/pull/1#discussion_r2065387971)

> What specific cgroups do you plan to use?

**@russjones** on `flychicken123/Job_Worker` [→](https://github.com/flychicken123/Job_Worker/pull/1#discussion_r2065389552)

> Can you provide some more details on how your process execution lifecycle. For example, how do you plan to place processes in cgroups?

**@rosstimothy** on `flychicken123/Job_Worker` [→](https://github.com/flychicken123/Job_Worker/pull/1#discussion_r2066226200)

> If the process is moved into the cgroup after it has started it allows the process time to fork an evade cgroups entirely. How can you protect against this?

**@eriktate** on `GevorgGal/jobworker` [→](https://github.com/GevorgGal/jobworker/pull/1#discussion_r2975147863)

> Attempting to clean up child processes isn't required, so feel free to drop this if it saves you some time.

**@Tener** on `GevorgGal/jobworker` [→](https://github.com/GevorgGal/jobworker/pull/4#discussion_r3020561983)

> This is fragile by design. For example, it will fail to catch:
>
> ```
> 	if c.Path == "" && c.Err == nil && c.lookPathErr == nil {
> 		c.Err = errors.New("exec: no command")
> 	}
> ```

**@tigrato** on `gstelang/job-worker-service` [→](https://github.com/gstelang/job-worker-service/pull/1#discussion_r1732901205)

> Could you please provide the following details:
> 
> - The TLS settings you intend to use.
> - An explanation of how the log stream will function.
> - Details on user identification.
> - The cgroup settings you plan to enforce and the method for adding a process to a cgroup.
> - Information on how process logs will be stored on the server (e.g., in memory, files, etc.).
> - The procedure for handling process termination (stop).

**@tigrato** on `gstelang/job-worker-service` [→](https://github.com/gstelang/job-worker-service/pull/1#discussion_r1734517762)

> >Am I allowed to use a library such as https://github.com/containerd/cgroups, correct?
> 
> unfortunately no. Candidates aren't allowed to use external libraries to implement the cgroup work.
> 
> > The cgroup settings you plan to enforce - what do you mean by this? Is this something I enforce or allow clients to pass as limit when starting a job?
> 
> The limit can be enforced by the server. Client's don't need to send it.

**@tigrato** on `gstelang/job-worker-service` [→](https://github.com/gstelang/job-worker-service/pull/1#discussion_r1735963957)

> Who is responsible for adding the process to the cgroup?
> How will you add PIDs to cgroups?

**@strideynet** on `gstelang/job-worker-service` [→](https://github.com/gstelang/job-worker-service/pull/1#discussion_r1736024226)

> Is there a possibility that the technique you will use to move the process into the cgroup will allow the process to run unconstrained by the cgroup for a short period? Is there a way that you can ensure that will not happen.

**@awly** on `hashsequence/Linux-Job-Worker` [→](https://github.com/hashsequence/Linux-Job-Worker/pull/2#discussion_r478698663)

> ```suggestion
> * killing of pid can be implemented through sigkill via syscall
> ```

**@sclevine** on `ilyazz/jobs` [→](https://github.com/ilyazz/jobs/pull/2#discussion_r1001286851)

> What happens to child processes during a graceful stop?

**@sclevine** on `ilyazz/jobs` [→](https://github.com/ilyazz/jobs/pull/3#discussion_r1007352481)

> A few lifecycle issues:
> 1. If the parent process (e.g., `test` binary) is killed, the child process keeps running (and cgroups are leaked)
> 2. It's impossible to kill the child process with SIGTERM on graceful stop because it has PID 1

**@sclevine** on `ilyazz/jobs` [→](https://github.com/ilyazz/jobs/pull/3#discussion_r1007545068)

> It's good practice to create a restricted parent cgroup, but if we go without that, I would at least make sure that cgroups created under the root path have the necessary controllers available (via `/sys/fs/cgroup/cgroup.subtree_control`)

**@rosstimothy** on `jaylane/job-scheduler` [→](https://github.com/jaylane/job-scheduler/pull/1#discussion_r1714212896)

> Can you provide mode detail about the process execution lifecycle of a job? How and when is the process added to the cgroup? How will a job be terminated?

**@rosstimothy** on `jaylane/job-scheduler` [→](https://github.com/jaylane/job-scheduler/pull/1#discussion_r1715315540)

> For simplicity feel free to drop SIGTERM and rely solely on SIGKILL. Will a SIGKILL alone guarantee that all subprocesses spawned by the job will be terminated?

**@rosstimothy** on `jaylane/job-scheduler` [→](https://github.com/jaylane/job-scheduler/pull/1#discussion_r1716069478)

> What do you mean be destroy the cgroup?

**@tigrato** on `joshuarubin/teleport-job-worker` [→](https://github.com/joshuarubin/teleport-job-worker/pull/3#discussion_r1743554482)

> by `exec` the job, you mean spinning another cmd.Exec instance or fully replace the execution with the job binary?

**@codingllama** on `joshuarubin/teleport-job-worker` [→](https://github.com/joshuarubin/teleport-job-worker/pull/3#discussion_r1744196558)

> This is more of an hypothetical, but how could we make the NET namespace more useful?

**@rosstimothy** on `joshuarubin/teleport-job-worker` [→](https://github.com/joshuarubin/teleport-job-worker/pull/3#discussion_r1744202232)

> Is there any additional work that needs do be done besides setting clone and unshare flags for proper isolation?

**@codingllama** on `joshuarubin/teleport-job-worker` [→](https://github.com/joshuarubin/teleport-job-worker/pull/4#discussion_r1747585565)

> Suggestion:
> 
> ```suggestion
> 		err := j.cmd.Process.Kill()
> 		if errors.Is(err, os.ErrProcessDone) {
> ```

**@greedy52** on `kiakeshmiri/process-runner` [→](https://github.com/kiakeshmiri/process-runner/pull/1#discussion_r1765679364)

> how does the process gets added to the cgroups when starting the command?

**@greedy52** on `kiakeshmiri/process-runner` [→](https://github.com/kiakeshmiri/process-runner/pull/1#discussion_r1765681656)

> what values/controls are we setting for these cgroup options? are those pre-defined?
> 
> Also, for simplicity, it is ok to set these values by the server instead of specifying them by the client.

**@tigrato** on `kiakeshmiri/process-runner` [→](https://github.com/kiakeshmiri/process-runner/pull/1#discussion_r1766516410)

> Is it possible for a program to temporarily escape cgroups?

**@tigrato** on `kkloberdanz/teleport-challenge` [→](https://github.com/kkloberdanz/teleport-challenge/pull/1#discussion_r2781992745)

> For simplicity, feel free to force kill the process.
>
> What happens if the process forked itself before receiving the sigkill? Is there a way of ensuring all processes will be terminated?

**@tigrato** on `kkloberdanz/teleport-challenge` [→](https://github.com/kkloberdanz/teleport-challenge/pull/1#discussion_r2782005832)

> Which cgroups controllers and restrictions will you enforce?

**@smallinsky** on `kkloberdanz/teleport-challenge` [→](https://github.com/kkloberdanz/teleport-challenge/pull/1#discussion_r2782075482)

> Does this approach guarantees that all jobs children process will be killed ?

**@smallinsky** on `kkloberdanz/teleport-challenge` [→](https://github.com/kkloberdanz/teleport-challenge/pull/7#discussion_r2811958815)

> Does the cgroup clanup should be also called in this error path ?

**@tigrato** on `kkloberdanz/teleport-challenge` [→](https://github.com/kkloberdanz/teleport-challenge/pull/7#discussion_r2811973760)

> ```suggestion
> 	cgroupMgr, err := resources.NewManager()
> ```
>
> What's the reasoning for declaring cgroupMgr, declaring mgr and then assign mgr to the cgroupMgr?
>
> This seems to have the same exact behaviour

**@smallinsky** on `kkloberdanz/teleport-challenge` [→](https://github.com/kkloberdanz/teleport-challenge/pull/7#discussion_r2812088620)

> What would you say by injecting noop cgroupMgr instead of doing if checks in the flow ?

**@tigrato** on `kkloberdanz/teleport-challenge` [→](https://github.com/kkloberdanz/teleport-challenge/pull/7#discussion_r2818548274)

> Same discussion as before. namespaces here serve to ensure the job is "isolated". 
> If not available, it should trigger an error

**@zmb3** on `kovyrin/teleport-exec` [→](https://github.com/kovyrin/teleport-exec/pull/1#discussion_r788224012)

> Can you please include the `remote_exec.proto` file? It seems to be missing.

**@rosstimothy** on `kovyrin/teleport-exec` [→](https://github.com/kovyrin/teleport-exec/pull/1#discussion_r788259562)

> Could you add the exact cgroups you plan to use below?

**@russjones** on `kovyrin/teleport-exec` [→](https://github.com/kovyrin/teleport-exec/pull/1#discussion_r788286518)

> Just to be clear, you're writing `container_exec` and not re-using something right?

**@russjones** on `kovyrin/teleport-exec` [→](https://github.com/kovyrin/teleport-exec/pull/1#discussion_r789266564)

> Taking an array of strings and feeding into `exec.Command` sounds great.

**@russjones** on `kovyrin/teleport-exec` [→](https://github.com/kovyrin/teleport-exec/pull/1#discussion_r789267036)

> You can do that, but is it reliable? I could exit with `137` but not really be killed by SIGKILL right?

**@zmb3** on `kovyrin/teleport-exec` [→](https://github.com/kovyrin/teleport-exec/pull/2#discussion_r795828186)

> 1. Why create the `exec.Command` here instead of in `NewCommand`? What happens if somebody creates a command with `NewCommand()`, doesn't `Start` it, and then calls `ResultCode`? (looks like it would panic)
> 2. This check should probably be done with the mutex held. Two concurrent calls to start could both see `executor` being nil and then sequentially try to start the command.

**@kpumuk** on `kovyrin/teleport-exec` [→](https://github.com/kovyrin/teleport-exec/pull/2#discussion_r797006281)

> nvm, found how you reset the env via reexec.

**@sclevine** on `kurczynski/teleport-job-worker` [→](https://github.com/kurczynski/teleport-job-worker/pull/1#discussion_r1610847711)

> How will you ensure the cgroup applies to this PID for the entirety of its execution?

**@sclevine** on `kurczynski/teleport-job-worker` [→](https://github.com/kurczynski/teleport-job-worker/pull/1#discussion_r1612583289)

> >  The job will fork from the service, the job will apply the cgroups, and then execute the process given by the user.
> 
> Awesome. Might want to make this a little more clear. It's clear to me now, re-reading:
> 
> > Jobs will be forked from the service's process and contain exactly one process that it manages.
> 
> But for a while I thought your service was managing the command process directly.

**@sclevine** on `kurczynski/teleport-job-worker` [→](https://github.com/kurczynski/teleport-job-worker/pull/2#discussion_r1621428644)

> How does the binary ensure that `job.Start()` is called on reexec?

**@sclevine** on `kurczynski/teleport-job-worker` [→](https://github.com/kurczynski/teleport-job-worker/pull/2#discussion_r1621429903)

> What happens if the job doesn't stop when it gets SIGTERM?

**@gabrielcorado** on `m3talsmith/jobberthehut` [→](https://github.com/m3talsmith/jobberthehut/pull/1#discussion_r2499476716)

> Additionally, how will users provide the cgroup control and limit values?

**@rosstimothy** on `m3talsmith/jobberthehut` [→](https://github.com/m3talsmith/jobberthehut/pull/1#discussion_r2500067557)

> How will a job be executed in a cgroup? How is output redirected to the file on disk?

**@rosstimothy** on `m3talsmith/jobberthehut` [→](https://github.com/m3talsmith/jobberthehut/pull/1#discussion_r2500072348)

> Will this terminate _all_ child processes of the job?

**@greedy52** on `m3talsmith/jobberthehut` [→](https://github.com/m3talsmith/jobberthehut/pull/1#discussion_r2500073820)

> if the PID is obtained from starting the process, won't the process run without any limit before adding to the cgroup?

**@rosstimothy** on `MarkDHarris/teleport-systems-challenge` [→](https://github.com/MarkDHarris/teleport-systems-challenge/pull/1#discussion_r3015931672)

> Remove the mention of process group termination from the design and move it to non-goals.
>
> ```suggestion
> - A `context.CancelFunc` used by `Stop()` to trigger process termination.
>   A standalone context is created per job with `context.WithCancel`. A dedicated
>   goroutine blocks on `<-ctx.Done()` and calls `cmd.Process.Kill()` when it
>   fires. 
> ```

**@rosstimothy** on `MarkDHarris/teleport-systems-challenge` [→](https://github.com/MarkDHarris/teleport-systems-challenge/pull/1#discussion_r3015938944)

> >  I wanted to use this pattern instead of exec.CommandContext so that I could sends the kill to the process group instead of only the parent process.
> 
> Why doesn't exec.CommandContext apply without process group termination? What benefits do we get from handling this manually instead of using exec.CommandContext?

**@tigrato** on `mcampo84/teleport_challenge` [→](https://github.com/mcampo84/teleport_challenge/pull/1#discussion_r1907293081)

> Pick just a single cgroup settings. Expanding to other cases is trivial and can be done later.

**@tigrato** on `mikewurtz/taskman` [→](https://github.com/mikewurtz/taskman/pull/4#discussion_r2040094275)

> Are we sure the system has the cgroup controllers enabled?

**@eriktate** on `MrChristianL/Teleport-Job-Worker-Service` [→](https://github.com/MrChristianL/Teleport-Job-Worker-Service/pull/1#discussion_r2886253817)

> > The server extracts these identities and compares them to a hard-coded map that translates the identities to allowed RPC methods
> 
> Does this mean that the identities themselves map directly to which RPCs they're permitted to use and roles don't actually exist?
> 
> > I did consider other options for implementing this authorization. I originally considered an external config file (YAML/JSON) to store the mappings, which would allow for runtime changes. However, for the sake of limiting the overall number of files for this challenge and complexity of use, this was opted against. I also looked into an external policy engine option (e.g. Open Policy Agent), but that felt overkill for this challenge and would have required adding a third-party dependency.
> 
> I agree these add more complexity than needed for this challenge :+1:

**@nklaassen** on `MrChristianL/Teleport-Job-Worker-Service` [→](https://github.com/MrChristianL/Teleport-Job-Worker-Service/pull/2#discussion_r2906698361)

> nit: RSA4096 is probably a bit overkill and slower than necessary for a user key, ECDSA is generally better these days

**@greedy52** on `P-A-R-U-S/Go-Job-Worker-Service` [→](https://github.com/P-A-R-U-S/Go-Job-Worker-Service/pull/2#discussion_r1743772378)

> leftover cgroups on these failure cases.

**@smallinsky** on `P-A-R-U-S/Go-Job-Worker-Service` [→](https://github.com/P-A-R-U-S/Go-Job-Worker-Service/pull/2#discussion_r1743787719)

> What if the process has children. Does the job.cmd.Process.Signal will also stop/terminate all process children's  ?

**@greedy52** on `P-A-R-U-S/Go-Job-Worker-Service` [→](https://github.com/P-A-R-U-S/Go-Job-Worker-Service/pull/2#discussion_r1745583607)

> ```suggestion
> 		cleanCGroup <- true
> ```

**@awly** on `renatoaguimaraes/golang-job-scheduler` [→](https://github.com/renatoaguimaraes/golang-job-scheduler/pull/2#discussion_r623383624)

> ```suggestion
> 	return job.Cmd.Process.Signal(syscall.SIGTERM)
> ```

**@rosstimothy** on `rexposadas/teleport` [→](https://github.com/rexposadas/teleport/pull/1#discussion_r1793929264)

> Will this terminate any and all child processes that may have been spawned by the root process of the job?

**@rosstimothy** on `rexposadas/teleport` [→](https://github.com/rexposadas/teleport/pull/1#discussion_r1793929862)

> Which cgroup controllers will these limits be applied to?

**@rosstimothy** on `RichyHBM/teleport-challenge` [→](https://github.com/RichyHBM/teleport-challenge/pull/3#discussion_r2064063863)

> Suggestion: add some entropy to the cgroup name to prevent any bugs in tests from leaking items which may impact future runs.

**@tigrato** on `rsteinkeXJ/teleport_interview` [→](https://github.com/rsteinkeXJ/teleport_interview/pull/1#discussion_r3123850545)

> feel free to sigkill

**@tigrato** on `rsteinkeXJ/teleport_interview` [→](https://github.com/rsteinkeXJ/teleport_interview/pull/1#discussion_r3149013057)

> does pidfd prevent pid recycling?
> if we use cgroup.kill, do we even care about the pid anymore? can we remove all the references?

**@smallinsky** on `samschurter/teleport-challenge` [→](https://github.com/samschurter/teleport-challenge/pull/1#discussion_r693127208)

> I think that this is overkill. Why worker library can't just expose library Go public interface ?

**@nklaassen** on `supby/job-worker` [→](https://github.com/supby/job-worker/pull/1#discussion_r715927345)

> What will the library do if the process does not die after SIGTERM is sent?

**@rosstimothy** on `tjper/teleport` [→](https://github.com/tjper/teleport/pull/2#discussion_r826932728)

> Shouldn't this be static for the life of a `Cgroup`? Would it make more sense for `path` to just be a member of `Cgroup` and assigned when the `Service` creates it?

**@codingllama** on `tjper/teleport` [→](https://github.com/tjper/teleport/pull/2#discussion_r827460543)

> +1 for Tim's suggestion, cgroups is mostly writing files, the tests shouldn't care about root or need special permissions.

**@codingllama** on `tjper/teleport` [→](https://github.com/tjper/teleport/pull/3#discussion_r829304244)

> Use *Job so invocations better match with job.New()
> 
> ```suggestion
> func (s *Service) StartJob(_ context.Context, job *Job, options ...cgroup.CgroupOption) error {
> ```

**@espadolini** on `Zephan92/teleport-handson` [→](https://github.com/Zephan92/teleport-handson/pull/1#discussion_r2718737882)

> Signaling process groups without races is pretty awkward to do, is there a better way to do it?

**@nklaassen** on `Zephan92/teleport-handson` [→](https://github.com/Zephan92/teleport-handson/pull/1#discussion_r2719101900)

> using the cgroup.kill file sounds good to me

**@espadolini** on `Zephan92/teleport-handson` [→](https://github.com/Zephan92/teleport-handson/pull/3#discussion_r2736706068)

> The contract of `os.Process` demands that nothing else reaps the process, this has the potential to do so and then the subsequent Kill and Wait might hit a completely unrelated process.

**@espadolini** on `Zephan92/teleport-handson` [→](https://github.com/Zephan92/teleport-handson/pull/3#discussion_r2736726696)

> Does waiting for the main process to exit guarantee that the cgroup is unpopulated? If we hit this fast enough I'm pretty sure we're going to always have leftover cgroups around.

**@nklaassen** on `Zephan92/teleport-handson` [→](https://github.com/Zephan92/teleport-handson/pull/3#discussion_r2737563840)

> I don't love mutating a package-level variable for this, it prevents parallel tests, there are other ways to make the cgroup root configurable

**@espadolini** on `Zephan92/teleport-handson` [→](https://github.com/Zephan92/teleport-handson/pull/3#discussion_r2741871323)

> I'm not sure this fallback will ever work if the intention is to get rid of any nested cgroups, won't it error out when removing any of the cgroup files?

**@nklaassen** on `Zephan92/teleport-handson` [→](https://github.com/Zephan92/teleport-handson/pull/3#discussion_r2743482393)

> I don't see a goto in the wild that often, and it's impossible to take the !seenMultiple branch
> ```suggestion
> 	seenMultiple := false
> 	for !seenMultiple {
> 		select {
> 		case <-timeout:
> 			t.Fatal("Timeout waiting for process tree to populate")
> 		case <-ticker.C:
> 			content, err := os.ReadFile(cgroupProcs)
> 			if err == nil {
> 				// PIDs are newline separated
> 				pids := strings.Fields(string(content))
> 				if len(pids) >= 2 {
> 					seenMultiple = true
> 				}
> 			}
> 		}
> 	}
> ```

---

## testing

_Test hermiticity (temp dirs, no hardcoded paths); t.Context; real vs mocked listeners; avoiding sleeps and flakiness; race conditions in test setup._

**110 quotes** from `21` distinct reviewers across `31` candidate submissions.

**@codingllama** on `adalton/teleport-exercise` [→](https://github.com/adalton/teleport-exercise/pull/1#discussion_r768012202)

> My 2c about this: create a temp dir for tests and let the code work in it. One could argue that the purity of unit tests is compromised, but as long as they are fast and don't leave garbage around I'm perfectly happy - more so by the fact that I don't have to worry about mocks misbehaving in subtle ways.

**@codingllama** on `adalton/teleport-exercise` [→](https://github.com/adalton/teleport-exercise/pull/1#discussion_r768013420)

> Suggestion: A type alias might be good here. Same for GetpidMock.

**@codingllama** on `adalton/teleport-exercise` [→](https://github.com/adalton/teleport-exercise/pull/1#discussion_r768013986)

> Why the special case? Shouldn't we let tests handle this?

**@codingllama** on `adalton/teleport-exercise` [→](https://github.com/adalton/teleport-exercise/pull/1#discussion_r768020397)

> FYI: This reminds me of [fs.FS](https://pkg.go.dev/io/fs#FS), so I did a quick search and it looks like [fstest.MapFS](https://pkg.go.dev/testing/fstest#MapFS) could fit the bill.

**@codingllama** on `adalton/teleport-exercise` [→](https://github.com/adalton/teleport-exercise/pull/1#discussion_r768068424)

> The limits tests look alright, but I wish they could 1) assert the results (and fail accordingly) and 2) were `go test`-able, so they could run in automated fashion.

**@codingllama** on `adalton/teleport-exercise` [→](https://github.com/adalton/teleport-exercise/pull/1#discussion_r768696899)

> I think we could get away with most simple types being aliases / defined in terms of other types, as we are also somewhat sure they won't change/increase in scope.
> 
> Eg:
> 
> ```go
> type EnvironMock []string
> 
> type GetpidMock int
> 
> // so on
> ```
> 
> Feel free to leave as-is, though.

**@codingllama** on `adalton/teleport-exercise` [→](https://github.com/adalton/teleport-exercise/pull/1#discussion_r768707202)

> I think we should at least try to have the binaries return a failed exit code, so it could be reasonably automated. This means they have more knowledge of what is going on and are capable of checking their output, so to speak.
> 
> As for making it `go test`-able, there are a few ways we can go about this: Benchmark tests, testing.Short(), setting a build tag, env var + t.Skip() - whatever seems more appropriate. Even if we don't, a quick test that asserts the happy path of the library would be great to have.

**@rosstimothy** on `adalton/teleport-exercise` [→](https://github.com/adalton/teleport-exercise/pull/2#discussion_r771428749)

> Perhaps you could pass a `net.Listener` into the `RunJobmanagerServer` function. Then in each test you could listen on `:0`. To remove the `Sleep` you could attempt to dial the address of the `net.Listner` until you get a successful connection or until some defined timeout/number of tries has been reached.

**@codingllama** on `adalton/teleport-exercise` [→](https://github.com/adalton/teleport-exercise/pull/2#discussion_r771451221)

> Are these values from the mock job? Could we refer to the mock job, then, instead of repeating the hardcoded values here?

**@codingllama** on `adalton/teleport-exercise` [→](https://github.com/adalton/teleport-exercise/pull/2#discussion_r771454460)

> Let's test with a proper client, we should not be mocking this.

**@codingllama** on `adalton/teleport-exercise` [→](https://github.com/adalton/teleport-exercise/pull/2#discussion_r771461251)

> Could we refactor the server/client setup so it's not repeated by various tests?

**@codingllama** on `adalton/teleport-exercise` [→](https://github.com/adalton/teleport-exercise/pull/2#discussion_r771462338)

> Is this testing anything new? Won't the other tests cover it?

**@codingllama** on `adalton/teleport-exercise` [→](https://github.com/adalton/teleport-exercise/pull/2#discussion_r771472278)

> Could we add a test for correct authorization? Ie, user trying to interact with someone else's job. I'd toss in an admin test as well.

**@codingllama** on `adalton/teleport-exercise` [→](https://github.com/adalton/teleport-exercise/pull/2#discussion_r771485359)

> I added a bunch of comments like this, so I'll elaborate why I ask for it:
> 
> - Team practice (we tend to not let errors slide, even in tests)
> - Having Start pass is important for this test - in other languages you'd get an early exit via exception, in Go you won't. You don't need to go overboard for asserts that are better covered in other tests, but an `err != nil` is fine.
> - It's annoying when tests swallow errors and fail mysteriously in some further line. I've lost a few minutes of my life tracking stuff like this.
> - Finally, less cognitive load: as a reviewer/reader I don't have to try and figure out _why_ this wasn't handled. If we end up handling a few impossible errors because of this, that's OK.
> 
> I get the philosophical aspects of it, but that's where I ended up (for now, at least :)).

**@rosstimothy** on `adalton/teleport-exercise` [→](https://github.com/adalton/teleport-exercise/pull/2#discussion_r771485578)

> I would argue that ignoring errors and assuming things are properly tested and verified by another test is a bit of a footgun. What if that other test which was checking the errors gets removed or has a bug?

**@codingllama** on `adalton/teleport-exercise` [→](https://github.com/adalton/teleport-exercise/pull/2#discussion_r771657289)

> I get the conceptual distinction, but I wouldn't recommend reaching for the server implementation directly in tests regardless - it tends to become brittle and leave important behavior out. If you want to run a "lighter" stack, there are alternatives like [bufconn.Listener](https://pkg.go.dev/google.golang.org/grpc@v1.43.0/test/bufconn) (although I would do a normal Dial in most cases). It's also sometimes possible to reuse the same server for all tests in a package, which amortizes the cost of setup.
> 
> We have other integration-like tests from you, so feel free to do as you prefer here.

**@codingllama** on `adalton/teleport-exercise` [→](https://github.com/adalton/teleport-exercise/pull/2#discussion_r772444658)

> nit: Either make sure you defer-Close the listener on errors or explicitly have RunJobmanagerServer take ownership of it. If panic because RunJobmanagerServer failed (around line 44) the listener leaks.

**@codingllama** on `adalton/teleport-exercise` [→](https://github.com/adalton/teleport-exercise/pull/2#discussion_r772448656)

> FYI, Go's testing doesn't like goroutines calling t.Fail, so this would cause problems on failures. t.Error variants are fair game.

**@sclevine** on `benmoss/job-worker-service` [→](https://github.com/benmoss/job-worker-service/pull/3#discussion_r2684511799)

> These subtests are largely independent. Why combine all of them into a single test?

**@tcsc** on `bkneis/jobworker` [→](https://github.com/bkneis/jobworker/pull/2#discussion_r1663349615)

> +1 for enabling the race detector 👍

**@tcsc** on `bkneis/jobworker` [→](https://github.com/bkneis/jobworker/pull/2#discussion_r1663353577)

> I can understand using indices into `argv` for this test program, even though it makes my maintenance brain wince :-)
> 
> There's always [flag](https://pkg.go.dev/flag) if you need anything more complex.

**@tcsc** on `bkneis/jobworker` [→](https://github.com/bkneis/jobworker/pull/2#discussion_r1663367272)

> [`T.TempDir()`](https://pkg.go.dev/testing#T.TempDir) is your friend here:
> 
> ```suggestion
> 	testDir := t.TempDir()
> ```
> 
> ... and the cleanup code below

**@tcsc** on `bkneis/jobworker` [→](https://github.com/bkneis/jobworker/pull/2#discussion_r1663371706)

> I'm generally one-feature-per-test person (in order to give the most signal in CI), but I get that it's tricky to do when you have a sequence of stuff to set up and tear down. 
> 
> Outside of this exercise, I would prefer to see this split up into multiple tests, or multiple subtests using [`T.Run()`](https://pkg.go.dev/testing#T.Run) - but I'll settle for a `TODO:` for now. 😄

**@tcsc** on `bkneis/jobworker` [→](https://github.com/bkneis/jobworker/pull/2#discussion_r1663406502)

> ```suggestion
> 	args := []string{"-c", fmt.Sprintf("for run in {1..%d}; do echo \\${run}: %s; sleep 0.1; done", n, echo)}
> ```

**@tcsc** on `bkneis/jobworker` [→](https://github.com/bkneis/jobworker/pull/2#discussion_r1663415654)

> In the real world I'd probably want some sort of short timeout on this, so a deadlocked test doesn't gum up CI until the overall test timeout expires. 
> 
> A quick `TODO` about how you'd approach this will do for now.

**@tcsc** on `bkneis/jobworker` [→](https://github.com/bkneis/jobworker/pull/2#discussion_r1663422223)

> Given that this can potentially take up to `STOP_GRACE_PERIOD`, is it worth taking a `Context` argument so that the caller can abort if they want?

**@tcsc** on `bkneis/jobworker` [→](https://github.com/bkneis/jobworker/pull/2#discussion_r1663424601)

> Similar comments as Concurrent Reader tests above.

**@rosstimothy** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/1#discussion_r1786187330)

> Perhaps it would be useful to at least include tests that validate the appropriate files are being written and that their contents match expectations.

**@espadolini** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/4#discussion_r1789822572)

> I'd say that it's sort of needed if you ever want to run the test with `-count 10000` or so.

**@rosstimothy** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/4#discussion_r1790173199)

> It'd also be a good idea to test that the hardcoded policies are enforcing the expected behavior for this exercise.

**@espadolini** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/5#discussion_r1791366637)

> ```suggestion
> 			time.Sleep(5 * time.Second)
> ```

**@codingllama** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/5#discussion_r1793611525)

> Maybe:
> 
> ```suggestion
> func assertChanOnce(ctx context.Context, t *testing.T, expect string, ch <-chan []byte, msg string) {
> ```
> 
> so it's clear that it's an assertion over a single receive?

**@codingllama** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/5#discussion_r1793639461)

> The is fine, but I would be tempted to just run a shell loop with a short sleep (or some other similar-behaving program). You get to do the same assertions but cut all the pipe logic.
> 
> (No need to change the tests.)

**@codingllama** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/6#discussion_r1795824642)

> Does anything use the test scripts? Are they meant for manual use?

**@rosstimothy** on `devnulled/runjob` [→](https://github.com/devnulled/runjob/pull/2#discussion_r1797091648)

> For simplicity you may omit graceful termination via SIGQUIT.

**@awly** on `diptadas/job-worker` [→](https://github.com/diptadas/job-worker/pull/2#discussion_r648579512)

> why is this sleep needed?

**@awly** on `diptadas/job-worker` [→](https://github.com/diptadas/job-worker/pull/2#discussion_r651966000)

> @diptadas this race condition is still present, since you're modifying the same `*Job` pointer

**@awly** on `diptadas/job-worker` [→](https://github.com/diptadas/job-worker/pull/2#discussion_r651966898)

> this is still a race condition - `outputBuffer` is the same pointer and is modified concurrently when the job is running

**@awly** on `diptadas/job-worker` [→](https://github.com/diptadas/job-worker/pull/2#discussion_r651967494)

> the job is `sleep 300`, it will not finish in 1s

**@awly** on `diptadas/job-worker` [→](https://github.com/diptadas/job-worker/pull/2#discussion_r651969208)

> > StopJob fails when job ID not found or, it fails to terminate the job
> 
> Job ID should be valid (because `CreateJob` succeeded above) and the code should be able to terminate `sleep 300` without any issues.
> I think that this test should call `t.Fatal` if `StopJob` fails.

**@zmb3** on `DiscoRiver/jobber` [→](https://github.com/DiscoRiver/jobber/pull/1#discussion_r2063502350)

> I'd think twice about exporting these fields. Allowing other packages to mutate them can result in data races or just break things.

**@greedy52** on `ehsaniara/joblet` [→](https://github.com/ehsaniara/joblet/pull/2#discussion_r2115926351)

> very interesting approach. one down side i see is that programs can detect if they are being traced and may react to it.

**@rosstimothy** on `flychicken123/Job_Worker` [→](https://github.com/flychicken123/Job_Worker/pull/1#discussion_r2069328777)

> Feel free to omit support for graceful termination to reduce scope.

**@eriktate** on `GevorgGal/jobworker` [→](https://github.com/GevorgGal/jobworker/pull/1#discussion_r2975249844)

> Will all gRPC integration happen in these cmd packages?

**@rosstimothy** on `GevorgGal/jobworker` [→](https://github.com/GevorgGal/jobworker/pull/2#discussion_r2988941894)

> Suggestion: Prefer using t.Context() in tests 
> ```suggestion
> 	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
> ```

**@Tener** on `GevorgGal/jobworker` [→](https://github.com/GevorgGal/jobworker/pull/4#discussion_r3020304387)

> I'm not sure this helps us in any significant way. The race is still there between `srv.GracefulStop()` and `signal.Stop(sigCh)`? Seems superfluous.

**@Tener** on `GevorgGal/jobworker` [→](https://github.com/GevorgGal/jobworker/pull/4#discussion_r3020318539)

> We miss "server shutting down" log entry for non-err outputs, which can happen with `Stop()/GracefulStop()`:
>
> > Serve will return a non-nil error unless Stop or GracefulStop is called.
>
> Also, should we return non-zero exit code when we were forced to call `Stop()`?

**@Tener** on `GevorgGal/jobworker` [→](https://github.com/GevorgGal/jobworker/pull/4#discussion_r3020665019)

> This should be mock or we are forced to spawn actual jobs for testing purposes, exposing us unnecessarily to edge conditions of underlying library and operating system. It also makes it so testing complex error conditions is complex too, instead of simply returning the error that triggers the behavior that we want tested.

**@Tener** on `GevorgGal/jobworker` [→](https://github.com/GevorgGal/jobworker/pull/4#discussion_r3020689084)

> This is rather convoluted; can we test it in a more direct way?

**@Tener** on `GevorgGal/jobworker` [→](https://github.com/GevorgGal/jobworker/pull/4#discussion_r3020782032)

> 700+ lines of hand-crafted testing code is quite impressive, but I'd appreciate if we could condense the test cases in some meaningful way, without throwing away the coverage. For example: having a common, shared set of known states for jobs that we are testing against would make it easier to see if anything is missing or not.
> 
> Also, some tests are table-driven while others are not; mixing styles worsens the readability and I think you can convert most, if not all of the tests, into table-driven ones.

**@sclevine** on `ilyazz/jobs` [→](https://github.com/ilyazz/jobs/pull/6#discussion_r1008985013)

> This gracefully stops the jobs, but not the gRPC server itself.

**@codingllama** on `joshuarubin/teleport-job-worker` [→](https://github.com/joshuarubin/teleport-job-worker/pull/4#discussion_r1747577248)

> I agree, this shouldn't be susceptible to races. Looping here is a bit strange.

**@tigrato** on `kiakeshmiri/process-runner` [→](https://github.com/kiakeshmiri/process-runner/pull/1#discussion_r1766511197)

> is it possible to live without sleep?

**@tigrato** on `kkloberdanz/teleport-challenge` [→](https://github.com/kkloberdanz/teleport-challenge/pull/6#discussion_r2798168935)

> use t.Context() instead of background

**@tigrato** on `kkloberdanz/teleport-challenge` [→](https://github.com/kkloberdanz/teleport-challenge/pull/6#discussion_r2798173121)

> this is ignoring the compiled protobuf boilerplate. This means go test will fail
>
> Feel free to include the generated stubs

**@smallinsky** on `kkloberdanz/teleport-challenge` [→](https://github.com/kkloberdanz/teleport-challenge/pull/8#discussion_r2833186051)

> nit: This seem to be duplicated with https://github.com/kkloberdanz/teleport-challenge/pull/8/changes#diff-9e543ce49053876900bc837a21c4f5e1c2cb80ea46d58a02346b7dd1fab8015cR141 
>
> Probably helper func and table test will  allow to reduce repetition.

**@zmb3** on `kovyrin/teleport-exec` [→](https://github.com/kovyrin/teleport-exec/pull/4#discussion_r794928587)

> Tests should clean up after themselves - it looks like this will leave a file behind each time it runs.

**@zmb3** on `kovyrin/teleport-exec` [→](https://github.com/kovyrin/teleport-exec/pull/4#discussion_r794929533)

> `strings.Repeat("x", 100)` would make more sense here.
> 
> I _assume_ there's 100 `x`s there, but I'm not very interested in counting them if this test starts to fail 😅

**@zmb3** on `kovyrin/teleport-exec` [→](https://github.com/kovyrin/teleport-exec/pull/4#discussion_r794931673)

> This is a bit of a heavy dependency for a small number of unit tests, and doesn't play well with standard Go tests. Fine for this project, but if it were me I would choose something much simpler.

**@nklaassen** on `kovyrin/teleport-exec` [→](https://github.com/kovyrin/teleport-exec/pull/4#discussion_r796191196)

> you can use `t.TempDir()` to create a portable temp dir which will automatically be cleaned up when your test exits

**@rhammonds-teleport** on `MarkDHarris/teleport-systems-challenge` [→](https://github.com/MarkDHarris/teleport-systems-challenge/pull/2#discussion_r3039789639)

> Small sleeps in unit tests can add up over time. Can any of these tests be reworked to avoid them?

**@greedy52** on `MarkDHarris/teleport-systems-challenge` [→](https://github.com/MarkDHarris/teleport-systems-challenge/pull/2#discussion_r3039973852)

> similar concern to comments above. `time.Sleep` is not reliable. `wg.Wait()` also get stuck potentially (this is a test so we cannot always assume implementation is correct)

**@rosstimothy** on `MarkDHarris/teleport-systems-challenge` [→](https://github.com/MarkDHarris/teleport-systems-challenge/pull/2#discussion_r3045038294)

> A way to know when a job has completed on demand could very much be useful in production settings. Please update to replace polling with a more efficient mechanism. Also, a sync.Cond might be a bit heavy handed for this use case.

**@rosstimothy** on `MarkDHarris/teleport-systems-challenge` [→](https://github.com/MarkDHarris/teleport-systems-challenge/pull/2#discussion_r3045093352)

> Is there any benefit to using 20 writes instead of a single write?
>
> ```suggestion
> 	if _, err := b.Write([]byte(strings.Repeat("x", 20))); err != nil {
> 		t.Fatalf("Write() error = %v", err)
> 	}
> ```

**@rosstimothy** on `MarkDHarris/teleport-systems-challenge` [→](https://github.com/MarkDHarris/teleport-systems-challenge/pull/2#discussion_r3045101182)

> Does this test need N Writes? Could we write all of the content in a single Write without affecting the test?

**@rosstimothy** on `MarkDHarris/teleport-systems-challenge` [→](https://github.com/MarkDHarris/teleport-systems-challenge/pull/2#discussion_r3045109598)

> Would a single write suffice here too?
> ```suggestion
> 	_, err := b.Write(bytes.Repeat(want, 5))
> 	if err != nil {
> 		t.Fatalf("Write() error = %v", err)
> 	}
> ```

**@fspmarshall** on `mcampo84/teleport_challenge` [→](https://github.com/mcampo84/teleport_challenge/pull/5#discussion_r1920406360)

> Halting consumption of command output before `EOF` is reached is a race condition.  Output must continue to be consumed until all output has been seen.

**@tigrato** on `mikewurtz/taskman` [→](https://github.com/mikewurtz/taskman/pull/1#discussion_r2027174410)

> what will be the grace period?

**@timothyb89** on `mikewurtz/taskman` [→](https://github.com/mikewurtz/taskman/pull/2#discussion_r2034227973)

> It might be even more robust to use something like [`require.Eventually()`](https://pkg.go.dev/github.com/stretchr/testify/require#Eventually) to both save time in the best case (it should usually take well under 100ms to spawn the server) and avoid failing in the worst case (e.g. a very slow CI runner)

**@rosstimothy** on `mikewurtz/taskman` [→](https://github.com/mikewurtz/taskman/pull/2#discussion_r2036221058)

> Is CGO required for tests?
> ```suggestion
> 	go test -v -race -cover ./...
> ```

**@tigrato** on `mikewurtz/taskman` [→](https://github.com/mikewurtz/taskman/pull/4#discussion_r2040104206)

> Given that the fields that the mutex guards are public, are we confident we don't have a data race here?

**@rosstimothy** on `mikewurtz/taskman` [→](https://github.com/mikewurtz/taskman/pull/4#discussion_r2040239261)

> Is the mutex meant to protect the Task and the map of tasks or just the map of tasks? Returning a pointer to the Task will make things very susceptible to data races.

**@eriktate** on `MrChristianL/Teleport-Job-Worker-Service` [→](https://github.com/MrChristianL/Teleport-Job-Worker-Service/pull/2#discussion_r2905806094)

> I don't think this test case actually exercises anything new compared to `TestLargeChunks`

**@eriktate** on `MrChristianL/Teleport-Job-Worker-Service` [→](https://github.com/MrChristianL/Teleport-Job-Worker-Service/pull/2#discussion_r2905892421)

> Is it possible that `10ms` would ever be too short to capture the intended test case? Is there a better way we can know whether or not the reader is caught up?

**@eriktate** on `MrChristianL/Teleport-Job-Worker-Service` [→](https://github.com/MrChristianL/Teleport-Job-Worker-Service/pull/2#discussion_r2906006211)

> Will this test actually produce the expected failure without the offset check?

**@eriktate** on `MrChristianL/Teleport-Job-Worker-Service` [→](https://github.com/MrChristianL/Teleport-Job-Worker-Service/pull/2#discussion_r2906263980)

> Do we need to sleep here?

**@Joerger** on `MrChristianL/Teleport-Job-Worker-Service` [→](https://github.com/MrChristianL/Teleport-Job-Worker-Service/pull/2#discussion_r2907574159)

> If these are called sequentially, I don't see how the caller will be able to avoid race conditions between each call. Likely will be an issue in the next PR.

**@rosstimothy** on `MrChristianL/Teleport-Job-Worker-Service` [→](https://github.com/MrChristianL/Teleport-Job-Worker-Service/pull/2#discussion_r2919508210)

> Is there any value in attempting these in parallel?

**@rosstimothy** on `MrChristianL/Teleport-Job-Worker-Service` [→](https://github.com/MrChristianL/Teleport-Job-Worker-Service/pull/3#discussion_r2924699806)

> Is there any way to not rely on relative paths? This makes it impossible to `go test -c` and the run the resulting binary from an alternate location.

**@rosstimothy** on `MrChristianL/Teleport-Job-Worker-Service` [→](https://github.com/MrChristianL/Teleport-Job-Worker-Service/pull/5#discussion_r2924866641)

> Suggestion: If using cobra, then go all in on cobra. Same suggestion elsewhere.
>
> ```suggestion
> 		// Start the job
> 		startCtx, startCancel := context.WithTimeout(cmd.Context(), 10*time.Second)
> 		defer startCancel()
>
> 		jobID, err := c.StartJob(startCtx, args)
> 		if err != nil {
> 			return fmt.Errorf("Error: %s", nicerErrors(err))
> 		}
>
> 		fmt.Fprintf(cmd.OutOrStdout(), "%s\n", jobID)
> 		return nil
> 	},
> ```

**@nklaassen** on `neildo/tjob` [→](https://github.com/neildo/tjob/pull/2#discussion_r1767899852)

> This is something that would be nice to have a unit test for

**@rosstimothy** on `neildo/tjob` [→](https://github.com/neildo/tjob/pull/2#discussion_r1769136529)

> Could this test be extended, or another be added to cover other edge cases? The challenge spec specifically calls out having tests for both the happy and unhappy path.

**@smallinsky** on `P-A-R-U-S/Go-Job-Worker-Service` [→](https://github.com/P-A-R-U-S/Go-Job-Worker-Service/pull/2#discussion_r1744979108)

> This seems to be a potential race condition. The entire function `func (job *Job) Start() error {` is called under the lock, but the go func() { can be scheduled after the function exit, and the lock is released by the defer `job.mutex.Unlock()` call.

**@smallinsky** on `P-A-R-U-S/Go-Job-Worker-Service` [→](https://github.com/P-A-R-U-S/Go-Job-Worker-Service/pull/2#discussion_r1748059052)

> This is dirty approach where the time Sleep is used to artificially delay operation to force  particular condition for UT.

**@smallinsky** on `P-A-R-U-S/Go-Job-Worker-Service` [→](https://github.com/P-A-R-U-S/Go-Job-Worker-Service/pull/3#discussion_r1749748803)

> This seems to be broken. 
> 
> the log.Fatalf will be never printed

**@atburke** on `rajansandeep/teleport-job-worker-svc` [→](https://github.com/rajansandeep/job-worker-svc/pull/2#discussion_r3101779835)

> What is this test asserting? Is it just that calling `Close()` twice doesn't panic?

**@atburke** on `rajansandeep/teleport-job-worker-svc` [→](https://github.com/rajansandeep/job-worker-svc/pull/3#discussion_r3112572895)

> Nit: it'd be really nice to have a test target too.

**@atburke** on `rajansandeep/teleport-job-worker-svc` [→](https://github.com/rajansandeep/job-worker-svc/pull/5#discussion_r3120515131)

> This target does not work without the integration build tag.

**@awly** on `renatoaguimaraes/golang-job-scheduler` [→](https://github.com/renatoaguimaraes/golang-job-scheduler/pull/2#discussion_r623308892)

> ```suggestion
> test:
> 	go test -race -v -timeout 30s -failfast -cover goteleport.com/...
> ```

**@awly** on `renatoaguimaraes/golang-job-scheduler` [→](https://github.com/renatoaguimaraes/golang-job-scheduler/pull/2#discussion_r623309059)

> ```suggestion
> ## Test
> 
> ```sh
> $ make test
> ```

**@Joerger** on `renatoaguimaraes/golang-job-scheduler` [→](https://github.com/renatoaguimaraes/golang-job-scheduler/pull/2#discussion_r623526305)

> ```suggestion
> 	jobID, err := w.Start(Command{Name: "sleep", Args: []string{"1"}})
> ```

**@awly** on `renatoaguimaraes/golang-job-scheduler` [→](https://github.com/renatoaguimaraes/golang-job-scheduler/pull/5#discussion_r625957594)

> ```suggestion
> 		log.Fatalf("fail to start server, %v", err)
> ```

**@awly** on `renatoaguimaraes/golang-job-scheduler` [→](https://github.com/renatoaguimaraes/golang-job-scheduler/pull/5#discussion_r625965044)

> here and below: use `require` instead of `assert` (different package from testify) to terminate the test on a failure; also use the appropriate helper functions for the check
> ```suggestion
> 	require.NoError(t, err)
> ```

**@awly** on `renatoaguimaraes/golang-job-scheduler` [→](https://github.com/renatoaguimaraes/golang-job-scheduler/pull/5#discussion_r625965637)

> is this necessary?
> the listener should already be accepting connections after `CreateServer` returns

**@rosstimothy** on `RichyHBM/teleport-challenge` [→](https://github.com/RichyHBM/teleport-challenge/pull/2#discussion_r2056053466)

> Suggetion: limit to localhost only for this challenge?
> ```suggestion
> 	listener, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
> ```

**@zmb3** on `RichyHBM/teleport-challenge` [→](https://github.com/RichyHBM/teleport-challenge/pull/2#discussion_r2056555590)

> Should we be closing the `listener` here to avoid leaking its fd?

**@eriktate** on `RichyHBM/teleport-challenge` [→](https://github.com/RichyHBM/teleport-challenge/pull/2#discussion_r2058497383)

> This test passes, but passing `-h` to the subcommand in practice leads to an `incorrect amount of arguments found` error

**@tigrato** on `RichyHBM/teleport-challenge` [→](https://github.com/RichyHBM/teleport-challenge/pull/5#discussion_r2066759917)

> Is there any viable option that doesn't require sleep?

**@rosstimothy** on `shawon-crosen/teleport-challenge` [→](https://github.com/shawon-crosen/teleport-challenge/pull/1#discussion_r2130918054)

> Per the challenge requirements:
> 
> > Discovering new output should be efficient, avoid busy-waiting or polling.

**@rosstimothy** on `tjper/teleport` [→](https://github.com/tjper/teleport/pull/1#discussion_r822155103)

> +1 to watching the file instead of sleeping

**@codingllama** on `tjper/teleport` [→](https://github.com/tjper/teleport/pull/2#discussion_r827461665)

> ```suggestion
> 		t.Fatalf("%+v", err)
> ```

**@codingllama** on `tjper/teleport` [→](https://github.com/tjper/teleport/pull/2#discussion_r827463648)

> How do we know the controllers actually did something? Point in case - "no options" passes the same test as the other controllers.

**@codingllama** on `tjper/teleport` [→](https://github.com/tjper/teleport/pull/3#discussion_r829333202)

> As soon as the test order changes this breaks.

**@codingllama** on `tjper/teleport` [→](https://github.com/tjper/teleport/pull/3#discussion_r829338591)

> Can we do this without the sleep? This may make the test flaky.

**@greedy52** on `xSudoNymx/job-manager` [→](https://github.com/xSudoNymx/job-manager/pull/3#discussion_r2733423668)

> where is `EndedAt` being used besides the test that tests it? same question for other similar functions like `job.Args()`, `log.CommittedLen()`

**@espadolini** on `Zephan92/teleport-handson` [→](https://github.com/Zephan92/teleport-handson/pull/1#discussion_r2718757596)

> The challenge text doesn't mention gracefully shutting down the server.

**@russjones** on `atburke/teleport_interview` [→](https://github.com/atburke/teleport_interview/pull/1#discussion_r643582075)

> I don't think you want to put everything in your main package (`main.go` and `database.go`). Your `main` package (and file) should basically just parse flags and call some library to do the actual work. This will also make it easier for you to write tests because you can target just that functionality.
> 
> For example you may want to do something like below.
> 
> ```
> /pkg/server/server.go
> /pkg/server/server_test.go
> /pkg/database/database.go
> /pkg/database/database_test.go
> /cmd
> /cmd/server.go
> ```
> 
> Then `/cmd/server.go` which will be your main package would look something like the following:
> 
> ```go
> package main
> 
> import (
>     "example.com/pkg/server"
> )
> 
> func main() {
>    var hostport := flag.String("hostport", ":8080", "help message for flag")
> 
>    server, err := server.New(&server.Config{
>       HostPort: *hostport,
>    })
>    if err != nil {
>      log.Fatal(err)
>    }
>    if err := server.Start(); err != nil {
>      log.Fatal(err)
>    }
> }
> ```
> 
> You don't have to use this exact structure, but having everything in the main package isn't a good idea.

**@alex-kovoy** on `ibeckermayer/teleport-interview` [→](https://github.com/ibeckermayer/teleport-interview/pull/3#discussion_r555525452)

> HOCs are usually prefixed with`with`, for example  createStore -> `withStore`

**@awly** on `ibeckermayer/teleport-interview` [→](https://github.com/ibeckermayer/teleport-interview/pull/3#discussion_r555981090)

> ```suggestion
> 		return nil, err
> 	}
> 
> 	// force a connection and test that it worked
> 	if err = sqlxdb.Ping(); err != nil {
> 		return nil, err
> 	}
> 
> 	db := &Database{cfg, sqlxdb}
> 	if err = db.init(); err != nil {
> 		return nil, err
> ```

**@kimlisa** on `zship/teleport-challenge` [→](https://github.com/zship/teleport-challenge/pull/1#discussion_r718715246)

> Testing was stated under server, so making sure that we don't forget including tests for the react comps

---

## scope-cutting

_Features reviewers recommend removing to narrow scope: interactive shell, multi-user ACLs, complex cleanup, Docker in favor of Go._

**92 quotes** from `20` distinct reviewers across `39` candidate submissions.

**@codingllama** on `adalton/teleport-exercise` [→](https://github.com/adalton/teleport-exercise/pull/1#discussion_r767953285)

> This is fine, I wouldn't worry unless you have a bunch of spare time.

**@codingllama** on `adalton/teleport-exercise` [→](https://github.com/adalton/teleport-exercise/pull/1#discussion_r767959926)

> Suggestion: simplify?
> 
> ```suggestion
> 	j.cmd.Env = nil // Do not pass along our environment
> ```

**@codingllama** on `adalton/teleport-exercise` [→](https://github.com/adalton/teleport-exercise/pull/1#discussion_r768040577)

> Could we do a RemoveAll in the job dir, or does that fail for the cgroups FS?

**@rosstimothy** on `adalton/teleport-exercise` [→](https://github.com/adalton/teleport-exercise/pull/2#discussion_r771462807)

> Perhaps you could make this a `const` in the `job_mock` so this doesn't seem so magical

**@codingllama** on `adalton/teleport-exercise` [→](https://github.com/adalton/teleport-exercise/pull/2#discussion_r771463043)

> Commented by others, but just as a reminder, let's find a way to remove the sleep.

**@codingllama** on `adalton/teleport-exercise` [→](https://github.com/adalton/teleport-exercise/pull/2#discussion_r771463921)

> Could we dismantle grpcutil and move the functions to other places? I'm not a fan of "util" packages, they tend to get too many responsibilities.

**@nklaassen** on `benmoss/job-worker-service` [→](https://github.com/benmoss/job-worker-service/pull/1#discussion_r2670600579)

> these are nice to have but it's fine to drop them to cut scope

**@nklaassen** on `benmoss/job-worker-service` [→](https://github.com/benmoss/job-worker-service/pull/3#discussion_r2683933786)

> to me there's no point putting the job into the jobs map before starting it, if you're going to immediately remove it if Start fails and not return its ID. Might as well try Start before adding it to w.jobs, what do you think?

**@rosstimothy** on `bill-rich/jobworker` [→](https://github.com/bill-rich/jobworker/pull/1#discussion_r1610680955)

> So the server will reexcute itself and call `run`? Could you please add more detail on how this will work and what the process execution lifecycle within the server will be.

**@tcsc** on `bkneis/jobworker` [→](https://github.com/bkneis/jobworker/pull/1#discussion_r1650248257)

> What happens if the process ignores the `SIGTERM`?

**@rosstimothy** on `bucknercd/jobworker` [→](https://github.com/bucknercd/jobworker/pull/1#discussion_r2266851662)

> Let's simplify this to always tail from the beginning of the output, no need to make clients aware of offsets or supporting multiple modes.

**@rosstimothy** on `bucknercd/jobworker` [→](https://github.com/bucknercd/jobworker/pull/1#discussion_r2269766747)

> Storing metadata on disk is totally fine, but, did you consider keeping this information in memory? Would doing so simplify your implementation any?

**@rosstimothy** on `bucknercd/jobworker` [→](https://github.com/bucknercd/jobworker/pull/1#discussion_r2269783268)

> Feel free to omit graceful termination to simplify things.

**@zmb3** on `bucknercd/jobworker` [→](https://github.com/bucknercd/jobworker/pull/1#discussion_r2271511182)

> By putting the command and its arguments into a single string like this, your program becomes responsible for parsing the string into a command and a set of arguments.
>
> This is not necessarily bad, but feels like unnecessary complexity. We're basically asking the user to go out of their way to make sure the shell _doesn't_ do this parsing (by quoting the whole string) and then doing it ourselves.
>
> Any opportunity to simplify?

**@rosstimothy** on `bucknercd/jobworker` [→](https://github.com/bucknercd/jobworker/pull/1#discussion_r2273680817)

> I think that the use of `/usr/local/bin/init-wrapper` would be in violation of the challenge requirements. Per the challenge spec:
> 
> > The server should also not rely on any shell scripts, external binaries or use containers to execute jobs.

**@rosstimothy** on `bucknercd/jobworker` [→](https://github.com/bucknercd/jobworker/pull/1#discussion_r2273689233)

> Was this reply intended for another thread? I was referring to dropping SIGTERM and relying solely on cgroup.kill here to cut scope.

**@rosstimothy** on `bucknercd/jobworker` [→](https://github.com/bucknercd/jobworker/pull/1#discussion_r2276607582)

> No need to worry about PATH or other environment variables. Feel free to omit to simplify and reduce scope.

**@rosstimothy** on `bucknercd/jobworker` [→](https://github.com/bucknercd/jobworker/pull/1#discussion_r2276608714)

> No need to worry about chroot. Feel free to omit to simplify and reduce scope.

**@espadolini** on `chintamanil/job-worker` [→](https://github.com/chintamanil/job-worker/pull/1#discussion_r2822099242)

> Not required by the challenge text.

**@strideynet** on `clydotron/job_worker_service` [→](https://github.com/clydotron/job_worker_service/pull/1#discussion_r1600233378)

> Persistence is out of scope for this task

**@codingllama** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/1#discussion_r1783280780)

> Feel free to cut this to save time/effort.

**@codingllama** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/1#discussion_r1783299116)

> nit: We require the job_id in the request, so technically we don't need it here.

**@codingllama** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/1#discussion_r1785259482)

> If it stops making sense after you combine Create and Start feel free to cut it altogether.

**@codingllama** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/1#discussion_r1785260774)

> SIGKILL is typically fine, although we see other strategies from time to time.

**@espadolini** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/6#discussion_r1795057879)

> This doesn't seem to be referenced anywhere else - crucially, it doesn't get closed either.

**@espadolini** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/6#discussion_r1795061237)

> Won't `RemoveAll` try to delete the special files inside of the cgroup directory?

**@codingllama** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/5#discussion_r1795619647)

> Is this used by prod code, or could it be removed?

**@rosstimothy** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/6#discussion_r1795684378)

> Should this directory be cleaned up or removed if any of the operations below result in returning an error?

**@rosstimothy** on `devnulled/runjob` [→](https://github.com/devnulled/runjob/pull/2#discussion_r1795672170)

> Per challenge requirements, the library may not use an external executable to run jobs.
> 
> > The server should also not rely on any shell scripts, external binaries or use containers to execute jobs.
> https://github.com/gravitational/careers/blob/main/challenges/systems/challenge-1.md#dependencies

**@rosstimothy** on `devnulled/runjob` [→](https://github.com/devnulled/runjob/pull/2#discussion_r1797075550)

> Perhaps it would be clearer to drop all references to the external `jobrunner` binary and be more specific that `jobmanager` will be reexecuted instead?

**@awly** on `diptadas/job-worker` [→](https://github.com/diptadas/job-worker/pull/1#discussion_r646774144)

> nit: this response doesn't seem useful
> if `StopJob` succeeds, respond with `200 OK` and no body

**@zmb3** on `DiscoRiver/jobber` [→](https://github.com/DiscoRiver/jobber/pull/1#discussion_r1989318013)

> ```suggestion
> - A GRPC-backed API server, which interacts with its host system only, allowing execution of arbitrary Linux commands.
> ```

**@zmb3** on `DiscoRiver/jobber` [→](https://github.com/DiscoRiver/jobber/pull/1#discussion_r1989327662)

> This is not required, feel free to remove it to reduce scope.

**@rosstimothy** on `DiscoRiver/jobber` [→](https://github.com/DiscoRiver/jobber/pull/1#discussion_r1989333151)

> Neither of these are required for the challenge, feel free to omit to simplify things.

**@rosstimothy** on `DiscoRiver/jobber` [→](https://github.com/DiscoRiver/jobber/pull/1#discussion_r1991730782)

> You may consider log rotation out of scope for this challenge.

**@russjones** on `dtom90/teleport-backend-challenge` [→](https://github.com/dtom90/teleport-backend-challenge/pull/1#discussion_r2053134987)

> No need for an interactive prompt like this, feel free to cut scope and just support 4 subcommands like you have above.

**@smallinsky** on `dtom90/teleport-backend-challenge` [→](https://github.com/dtom90/teleport-backend-challenge/pull/1#discussion_r2053963651)

> Please feel free to cut the code and design and skip the cleanup for the process log output files.

**@greedy52** on `ehsaniara/joblet` [→](https://github.com/ehsaniara/joblet/pull/2#discussion_r2113886097)

> job listing is out of scope so feel free to drop it.

**@zmb3** on `ehsaniara/joblet` [→](https://github.com/ehsaniara/joblet/pull/2#discussion_r2114131803)

> These two seem like internal implementation details that don't need to be exposed to the client.

**@rosstimothy** on `flychicken123/Job_Worker` [→](https://github.com/flychicken123/Job_Worker/pull/1#discussion_r2066220903)

> Supporting environment variables is not required for the challenge.
> ```suggestion
> Start/stop jobs with configurable commands, arguments.
> ```

**@rosstimothy** on `GevorgGal/jobworker` [→](https://github.com/GevorgGal/jobworker/pull/1#discussion_r2975018882)

> There is no requirement to list jobs. I suggest cutting it from scope to reduce the amount of work you have for this exercise.

**@rosstimothy** on `GevorgGal/jobworker` [→](https://github.com/GevorgGal/jobworker/pull/1#discussion_r2975035739)

> As mentioned above I suggest cutting this to reduce scope. 
>
>
> Disclaimer: I'm not asking you to implement this, I'm only interested in your thoughts and ideas in how to solve the problem.
>
> How would you handle a situation where the system is tracking a large enough set of jobs that they cannot all fit into a single ListJobsResponse message?

**@rosstimothy** on `GevorgGal/jobworker` [→](https://github.com/GevorgGal/jobworker/pull/1#discussion_r2975060308)

> You may simplify and omit graceful termination via SIGTERM.

**@rosstimothy** on `ghostsquad/prototype-job-worker` [→](https://github.com/ghostsquad/prototype-job-worker/pull/1#discussion_r1531953971)

> > My plan was to show that multiple users can have different permissions.
>
> The point of this design is to document your plans. Please include an example of users and their permissions.

**@awly** on `hashsequence/Linux-Job-Worker` [→](https://github.com/hashsequence/Linux-Job-Worker/pull/4#discussion_r480276128)

> No need to prefix every RPC with `Execute`. This can be just `Start`.
> Same for other RPCs below

**@andreiko** on `ilyazz/jobs` [→](https://github.com/ilyazz/jobs/pull/2#discussion_r1001145732)

> Is it possible to extend the authorization capabilities a little to support some sort of simple ACL? Hardcoding rules for that ACL in the server code is fine.

**@sclevine** on `ilyazz/jobs` [→](https://github.com/ilyazz/jobs/pull/2#discussion_r1003828016)

> It can, but processes with PID 1 don't have default handlers for SIGTERM (nothing happens). Graceful termination will take some work 🙂

**@andreiko** on `ilyazz/jobs` [→](https://github.com/ilyazz/jobs/pull/3#discussion_r1007491688)

> Totally optional, but If you're curious and have time, try it out — could be an opportunity to simplify things.

**@sclevine** on `ilyazz/jobs` [→](https://github.com/ilyazz/jobs/pull/3#discussion_r1007567848)

> The job can execute any arbitrary command, which may run additional processes.

**@codingllama** on `jaylane/job-scheduler` [→](https://github.com/jaylane/job-scheduler/pull/1#discussion_r1714368451)

> What if the job ignores SIGTERM?

**@rosstimothy** on `joshuarubin/teleport-job-worker` [→](https://github.com/joshuarubin/teleport-job-worker/pull/3#discussion_r1744034734)

> SIGINT + SIGKILL will work, but feel free to simplify things and solely use SIGKILL and make note of the SIGINT approach as an alternative.

**@rosstimothy** on `joshuarubin/teleport-job-worker` [→](https://github.com/joshuarubin/teleport-job-worker/pull/3#discussion_r1744198337)

> Feel free to drop the timeout and SIGINT in favor of always SIGKILLing the process.

**@codingllama** on `joshuarubin/teleport-job-worker` [→](https://github.com/joshuarubin/teleport-job-worker/pull/4#discussion_r1747520243)

> nit: We are using a lot of different concurrency control structures here. While this works (and we do need dedicated fields for some), I wonder if we could simplify and have less concurrency-related fields as a result.

**@codingllama** on `joshuarubin/teleport-job-worker` [→](https://github.com/joshuarubin/teleport-job-worker/pull/4#discussion_r1747657918)

> There are a lot of concurrency control structures here too, I wonder if we could simplify.

**@tigrato** on `kkloberdanz/teleport-challenge` [→](https://github.com/kkloberdanz/teleport-challenge/pull/6#discussion_r2798162447)

> you can register t.Cleanup and they will be executed when the test completes

**@tigrato** on `kkloberdanz/teleport-challenge` [→](https://github.com/kkloberdanz/teleport-challenge/pull/7#discussion_r2818570713)

> Does this remove all files within the directories?
>
> godoc:
> >Remove removes the named file or (empty) directory. If there is an error, it will be of type [*PathError](https://pkg.go.dev/os#PathError).

**@zmb3** on `kovyrin/teleport-exec` [→](https://github.com/kovyrin/teleport-exec/pull/1#discussion_r788239604)

> Agree that keeping it forever opens the server up to DoS attacks, but as long as you mention that in the doc I think it's fine. 
>
> Explicit delete like docker is a fine solution but that seems like busy work and is not required as part of the challenge.

**@russjones** on `kovyrin/teleport-exec` [→](https://github.com/kovyrin/teleport-exec/pull/1#discussion_r788284114)

> You are planning to execute through bash? Why not just directly execute processes?

**@russjones** on `kovyrin/teleport-exec` [→](https://github.com/kovyrin/teleport-exec/pull/1#discussion_r789266785)

> Yes, should cut scope quite a bit, what do you think?

**@rosstimothy** on `kovyrin/teleport-exec` [→](https://github.com/kovyrin/teleport-exec/pull/4#discussion_r795071520)

> The `log.Fatal` here will exit the process before `log.Print` can be executed

**@GavinFrazar** on `kurczynski/teleport-job-worker` [→](https://github.com/kurczynski/teleport-job-worker/pull/1#discussion_r1610799710)

> ~~I think user-specified resource limits is out of scope - you can hardcode some set of resource limits instead of letting users choose for each job.~~

**@sclevine** on `kurczynski/teleport-job-worker` [→](https://github.com/kurczynski/teleport-job-worker/pull/1#discussion_r1610839307)

> @GavinFrazar Not out of scope for this challenge:
> 
> > Add resource control for CPU, Memory and Disk IO **per job** using cgroups.

**@rosstimothy** on `m3talsmith/jobberthehut` [→](https://github.com/m3talsmith/jobberthehut/pull/1#discussion_r2500040130)

> I suggest simplifying this to retrieving a single job instead of all jobs. Also, depending on how many jobs exist this RPC will start to run into gRPC max message size issues.

**@rosstimothy** on `m3talsmith/jobberthehut` [→](https://github.com/m3talsmith/jobberthehut/pull/1#discussion_r2500089061)

> Note: There is no requirement to list all jobs. The requirement is to retrieve information about a single job.

**@rosstimothy** on `MarkDHarris/teleport-systems-challenge` [→](https://github.com/MarkDHarris/teleport-systems-challenge/pull/1#discussion_r3010240619)

> As mentioned above, this adds additional scope and is not required to satisfy the challenge requirements. I suggest removing it to simplify things for yourself.

**@greedy52** on `MarkDHarris/teleport-systems-challenge` [→](https://github.com/MarkDHarris/teleport-systems-challenge/pull/2#discussion_r3041543572)

> i am not against it, but at the same time curious whether `Closer` and input context are both required. Seems they are doing the same thing.

**@fspmarshall** on `mcampo84/teleport_challenge` [→](https://github.com/mcampo84/teleport_challenge/pull/1#discussion_r1907426950)

> Feel free to omit this It isn't required by the challenge, and implementing a command whitelist tends to be a pretty iffy undertaking at the best of times. The line between "restricted to the point of being a toy" and "unrestricted enough to allow anyone informed and motivated to do practically anything" is vanishingly small.

**@rosstimothy** on `mikewurtz/taskman` [→](https://github.com/mikewurtz/taskman/pull/1#discussion_r2027237461)

> > A rate limiter will be used to throttle the frequency of messages sent
> 
> Feel free to omit this for simplicity.

**@rosstimothy** on `mikewurtz/taskman` [→](https://github.com/mikewurtz/taskman/pull/2#discussion_r2035629102)

> I think it'd be fine for this challenge to omit timeouts and only honor terminating requests in response to receiving a signal from the user (i.e. they ctrl+c).

**@rosstimothy** on `MrChristianL/Teleport-Job-Worker-Service` [→](https://github.com/MrChristianL/Teleport-Job-Worker-Service/pull/1#discussion_r2890779168)

> I suggest dropping this to reduce scope.

**@eriktate** on `MrChristianL/Teleport-Job-Worker-Service` [→](https://github.com/MrChristianL/Teleport-Job-Worker-Service/pull/2#discussion_r2905821674)

> Not a requirement, but it could be useful to at least log potential errors here?

**@rosstimothy** on `MrChristianL/Teleport-Job-Worker-Service` [→](https://github.com/MrChristianL/Teleport-Job-Worker-Service/pull/5#discussion_r2924874439)

> Suggestion: simplify and drop support for stdin.

**@rosstimothy** on `neildo/tjob` [→](https://github.com/neildo/tjob/pull/1#discussion_r1759125374)

> Feel free to simplify things for this challenge and ditch graceful termination. It's fine to SIGKILL only for this challenge.

**@greedy52** on `P-A-R-U-S/Go-Job-Worker-Service` [→](https://github.com/P-A-R-U-S/Go-Job-Worker-Service/pull/1#discussion_r1731738513)

> sounds good. `stderr` is important so I just want to make sure it's included. ideally, the client can tell them (`stdout` vs `stderr`) apart but it's not required for this project.

**@zmb3** on `razzam21/job-worker-service` [→](https://github.com/razzam21/job-worker-service/pull/1#discussion_r2255205601)

> We're fans of testify, feel free to use it in your tests to simplify the assertions.

**@rosstimothy** on `sabernabil12/teleport-job-worker` [→](https://github.com/sabernabil12/teleport-job-worker/pull/1#discussion_r2164724420)

> Process isolation is not a requirement for level 4.

**@rosstimothy** on `sabernabil12/teleport-job-worker` [→](https://github.com/sabernabil12/teleport-job-worker/pull/1#discussion_r2164733397)

> As mentioned above, this is not required for the challenge.

**@rosstimothy** on `sabernabil12/teleport-job-worker` [→](https://github.com/sabernabil12/teleport-job-worker/pull/1#discussion_r2164746796)

> This is not required for the challenge.

**@smallinsky** on `samschurter/teleport-challenge` [→](https://github.com/samschurter/teleport-challenge/pull/1#discussion_r692320517)

> There isn't any requirement about `take` functionality in job task requirements thus endpoint can be remove in order to reduce scope.

**@bernardjkim** on `samschurter/teleport-challenge` [→](https://github.com/samschurter/teleport-challenge/pull/1#discussion_r693075968)

> List API call not required. Can cut to reduce scope.

**@rosstimothy** on `shawon-crosen/teleport-challenge` [→](https://github.com/shawon-crosen/teleport-challenge/pull/1#discussion_r2126986236)

> This is not required for this challenge.

**@greedy52** on `xSudoNymx/job-manager` [→](https://github.com/xSudoNymx/job-manager/pull/1#discussion_r2705003008)

> FYI `list` is not required for this challenge so feel free to drop it.

**@tigrato** on `xSudoNymx/job-manager` [→](https://github.com/xSudoNymx/job-manager/pull/1#discussion_r2705324829)

> I would expect the client subscription to be clear once the client disconnects. That will simplify the job and avoid memory leaks

**@greedy52** on `xSudoNymx/job-manager` [→](https://github.com/xSudoNymx/job-manager/pull/1#discussion_r2709756567)

> FYI, `env` is not required.

**@sclevine** on `zeeshanhaque21/JobWorkerService` [→](https://github.com/zeeshanhaque21/JobWorkerService/pull/1#discussion_r1686900871)

> Scheduling is not a challenge requirement

**@sclevine** on `zeeshanhaque21/JobWorkerService` [→](https://github.com/zeeshanhaque21/JobWorkerService/pull/1#discussion_r1686902955)

> Utilization reporting is not a challenge requirement

**@sclevine** on `zeeshanhaque21/JobWorkerService` [→](https://github.com/zeeshanhaque21/JobWorkerService/pull/1#discussion_r1686903398)

> Listing jobs is not a challenge requirement

**@espadolini** on `Zephan92/teleport-handson` [→](https://github.com/Zephan92/teleport-handson/pull/1#discussion_r2718734713)

> The challenge text does not mention using namespaces for isolation, I would recommend sticking to the requirements and leaving anything extra for the end if there's time left.

**@russjones** on `atburke/teleport_interview` [→](https://github.com/atburke/teleport_interview/pull/1#discussion_r645902017)

> Three things:
> 
> 1. Do you need this function, or can you update your `FetchSession` query to not fetch expired sessions then let your cleanup loop remove it whenever it's expired.
> 2. The way you are setting `ExpireIdle` and `ExpireAbs` I think this will always trigger on `ExpireIdle` because you are not doing any renewal as far as I can tell.
> 3. Renewal is not in scope, so just set a maximum session life time for your session. Upon logout just remove it. No need for idle timeouts.

**@russjones** on `atburke/teleport_interview` [→](https://github.com/atburke/teleport_interview/pull/1#discussion_r645903634)

> If you update your query to not return expired sessions, you don't need this function.

**@r0mant** on `atburke/teleport_interview` [→](https://github.com/atburke/teleport_interview/pull/1#discussion_r645904014)

> Wouldn't this always delete any session that's existed for longer than 30 minutes (current expire_idle duration) - even if it's still "active"? I couldn't see expire_idle field being updated anywhere.

**@alex-kovoy** on `ibeckermayer/teleport-interview` [→](https://github.com/ibeckermayer/teleport-interview/pull/1#discussion_r551415615)

> > and set the idle timeout to 5 minutes and the absolute timeout to 1 hour.
> 
> Lets ignore the idle timeout. I do not think that it's required. It would be strange to logout a user after any period of inactivity for this app. I would consider implementing timeout mechanism for banking apps but not for this one.

---

## design-doc-style

_Design doc readability and structure; explaining workflow and architecture; clarifying unclear descriptions._

**64 quotes** from `22` distinct reviewers across `32` candidate submissions.

**@sclevine** on `benmoss/job-worker-service` [→](https://github.com/benmoss/job-worker-service/pull/3#discussion_r2688686278)

> 👍
>
> Might be helpful to add a few comments in some of the longer functions explaining the workflow

**@smallinsky** on `bill-rich/jobworker` [→](https://github.com/bill-rich/jobworker/pull/1#discussion_r1609429516)

> Sorry but based on this description the flow is still quite unclear for me. 
> 
> Could you explain what "a new helper jobworker is started"  means.

**@rosstimothy** on `bucknercd/jobworker` [→](https://github.com/bucknercd/jobworker/pull/1#discussion_r2266853290)

> Let's omit the TTL and keep all job information for the life of the server.

**@rosstimothy** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/2#discussion_r1787023514)

> There appears to be some drift here between initial POC code and the approved design document.

**@rosstimothy** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/2#discussion_r1787023786)

> This doesn't match what's in the design document either.

**@espadolini** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/6#discussion_r1795015158)

> nit because we're specifically on Linux here and the two are equivalent: file paths should be manipulated with `path/filepath` rather than `path`

**@Joerger** on `diptadas/job-worker` [→](https://github.com/diptadas/job-worker/pull/1#discussion_r645188971)

> Can you specify the type for each of the Job object's fields?

**@zmb3** on `DiscoRiver/jobber` [→](https://github.com/DiscoRiver/jobber/pull/1#discussion_r1989321940)

> ```suggestion
>     And Mark should not receive any job information
> ```

**@rosstimothy** on `DiscoRiver/jobber` [→](https://github.com/DiscoRiver/jobber/pull/1#discussion_r1989353523)

> Could you explain a bit more how termination will work with a cancel function?

**@rosstimothy** on `dtom90/teleport-backend-challenge` [→](https://github.com/dtom90/teleport-backend-challenge/pull/1#discussion_r2054023060)

> Or specify a default value for the flag.

**@zmb3** on `ehsaniara/joblet` [→](https://github.com/ehsaniara/joblet/pull/2#discussion_r2114139258)

> What does "1000 lines" mean? I'm not sure we should design based on the concept of a "line" because that assumes to much about the format of the output.

**@zmb3** on `ehsaniara/joblet` [→](https://github.com/ehsaniara/joblet/pull/2#discussion_r2121709230)

> I don't think you should force user commands to run via bash, if they want a shell, they can specify the shell as the command they want to run.
> 
> I'm thinking there's a better way to structure the CLI interface so that it doesn't see the command + args as a single string.

**@rosstimothy** on `flychicken123/Job_Worker` [→](https://github.com/flychicken123/Job_Worker/pull/1#discussion_r2066226760)

> Please include that information in the design document.

**@rosstimothy** on `flychicken123/Job_Worker` [→](https://github.com/flychicken123/Job_Worker/pull/1#discussion_r2069335577)

> How will you achieve that in Go? Please update the design document to detail said process.

**@GavinFrazar** on `flychicken123/Job_Worker` [→](https://github.com/flychicken123/Job_Worker/pull/1#discussion_r2069374697)

> Please add any EKUs you plan to use (and why) to the design doc

**@rosstimothy** on `GevorgGal/jobworker` [→](https://github.com/GevorgGal/jobworker/pull/3#discussion_r3002844725)

> Should we specify a config file for consistent results?

**@rosstimothy** on `GevorgGal/jobworker` [→](https://github.com/GevorgGal/jobworker/pull/4#discussion_r3018987867)

> The [design](https://github.com/GevorgGal/jobworker/blob/main/docs/design.md?plain=1#L28) specified that graceful shutdown would be a future improvement and not something that is implemented for the challenge? How come you deviated from the design here?

**@fspmarshall** on `hashsequence/Linux-Job-Worker` [→](https://github.com/hashsequence/Linux-Job-Worker/pull/2#discussion_r478724516)

> Could you elaborate on the trust model here a bit?  Why or why not would a specific server trust a specific client or vise-versa?

**@codingllama** on `joshuarubin/teleport-job-worker` [→](https://github.com/joshuarubin/teleport-job-worker/pull/4#discussion_r1747501832)

> What is the expected format for these?

**@codingllama** on `joshuarubin/teleport-job-worker` [→](https://github.com/joshuarubin/teleport-job-worker/pull/4#discussion_r1747502850)

> What is the expected format / units for these?
> 
> (Also, as Tiago hinted, should we use stronger types?)

**@codingllama** on `joshuarubin/teleport-job-worker` [→](https://github.com/joshuarubin/teleport-job-worker/pull/4#discussion_r1750834583)

> Yep, I know, we keep hitting stuff like this in our platform-specific code too.

**@greedy52** on `kiakeshmiri/process-runner` [→](https://github.com/kiakeshmiri/process-runner/pull/1#discussion_r1766856966)

> My fault if the question isn't clear. What are the specific controller and values you would set to limit the resources? for example, are we setting `memory.max`? what value for it?

**@greedy52** on `kiakeshmiri/process-runner` [→](https://github.com/kiakeshmiri/process-runner/pull/1#discussion_r1767125761)

> we are just interested in the format of the ID, like how the uniq number is generated

**@greedy52** on `kiakeshmiri/process-runner` [→](https://github.com/kiakeshmiri/process-runner/pull/1#discussion_r1767167185)

> Again, it's ok for this project to take values from client or just hardcoded on the server side. But the project should demonstrate that
> - how to apply a value to a specific control
> - what controls are we setting when we say limit CPU/memory/IO

**@zmb3** on `kovyrin/teleport-exec` [→](https://github.com/kovyrin/teleport-exec/pull/4#discussion_r794926853)

> This is a bit inefficient, as every chunk of bytes you receive needs to be copied into a string before it can be printed. Further, the `fmt` package is best suited towards _formatted_ output, but you don't really care about formatting, you just want to write the bytes verbatim.

**@sclevine** on `kurczynski/teleport-job-worker` [→](https://github.com/kurczynski/teleport-job-worker/pull/1#discussion_r1612586075)

> ```suggestion
> // NewJob Create a new job to run the specified command.
> func NewJob(command string, args ...string) *Job
> ```

**@GavinFrazar** on `kurczynski/teleport-job-worker` [→](https://github.com/kurczynski/teleport-job-worker/pull/1#discussion_r1612843191)

> For this challenge I think it's ok to just pre-generate the key pairs for client and server, but can you provide explanation for how you will generate all the key pairs necessary to interact with the service, like specific commands and settings you will use?
> 
> Preferably as a script (you can link to it instead of copying it inline)

**@bernardjkim** on `kurczynski/teleport-job-worker` [→](https://github.com/kurczynski/teleport-job-worker/pull/1#discussion_r1614127541)

> Gotcha. The name is okay. It could be helpful to add this context to the comments.

**@greedy52** on `MarkDHarris/teleport-systems-challenge` [→](https://github.com/MarkDHarris/teleport-systems-challenge/pull/1#discussion_r3009926354)

> I am a little confused on `context.CancelFunc`. Is it canceling a context passed in to `exec.Cmd`? What are the things actually happen when `Stop` is called?

**@rosstimothy** on `MarkDHarris/teleport-systems-challenge` [→](https://github.com/MarkDHarris/teleport-systems-challenge/pull/1#discussion_r3023078632)

> Suggestion: Consider using context.AfterFunc here as well.

**@rosstimothy** on `MarkDHarris/teleport-systems-challenge` [→](https://github.com/MarkDHarris/teleport-systems-challenge/pull/2#discussion_r3045069027)

> Suggestion: Use testing.T.Context() instead of a background context.
> ```suggestion
> 	r := job.OutputReader(t.Context())
> ```

**@tigrato** on `mikewurtz/taskman` [→](https://github.com/mikewurtz/taskman/pull/2#discussion_r2035451691)

> should we receive `context.Context`?
> 
> same for all methods

**@tigrato** on `mikewurtz/taskman` [→](https://github.com/mikewurtz/taskman/pull/4#discussion_r2039951564)

> context is commonly the first arg

**@rosstimothy** on `mikewurtz/taskman` [→](https://github.com/mikewurtz/taskman/pull/4#discussion_r2040230610)

> It's idiomatic for context to be the first argument in a function.

**@russjones** on `moalf/teleport` [→](https://github.com/moalf/teleport/pull/1#discussion_r2604766251)

> Looks good to me, please update design document.

**@eriktate** on `MrChristianL/Teleport-Job-Worker-Service` [→](https://github.com/MrChristianL/Teleport-Job-Worker-Service/pull/1#discussion_r2892345209)

> > When the IDs are generated, there is a for loop that checks if that ID exists in the job tracker's map already. If it does, it generates another random ID. This occurs at the time of generation, so no job ever has the same ID.
> 
> Just to clarify, are we scanning all jobs to detect collisions? Or is the for loop a retry?

**@eriktate** on `MrChristianL/Teleport-Job-Worker-Service` [→](https://github.com/MrChristianL/Teleport-Job-Worker-Service/pull/2#discussion_r2905787568)

> nit: prefer `t.Context()` over `context.Background()` in tests

**@rosstimothy** on `MrChristianL/Teleport-Job-Worker-Service` [→](https://github.com/MrChristianL/Teleport-Job-Worker-Service/pull/2#discussion_r2919537710)

> Could this goroutine be eliminated if closing of resources was explicit and not tied to the context lifespan?

**@rosstimothy** on `neildo/tjob` [→](https://github.com/neildo/tjob/pull/1#discussion_r1761093581)

> This was resolved but I don't see any updates to the document to expose an address.

**@rosstimothy** on `neildo/tjob` [→](https://github.com/neildo/tjob/pull/1#discussion_r1761104607)

> This was resolved, but I'm not sure that I see any updates to the document that address the request here.

**@smallinsky** on `P-A-R-U-S/Go-Job-Worker-Service` [→](https://github.com/P-A-R-U-S/Go-Job-Worker-Service/pull/1#discussion_r1732293692)

> I agree, keeping the data in memory is totally acceptable. Please update the documentation to reflect this.

**@eriktate** on `rajansandeep/teleport-job-worker-svc` [→](https://github.com/rajansandeep/job-worker-svc/pull/1#discussion_r3095059492)

> Not a blocking comment, but I would consider if there's a way we could handle this automatically rather than depending on the caller to always remember to handle context cancellations.

**@rhammonds-teleport** on `rajansandeep/teleport-job-worker-svc` [→](https://github.com/rajansandeep/job-worker-svc/pull/2#discussion_r3100957626)

> Since we don't use it, can we drop the context from the signature?

**@atburke** on `rajansandeep/teleport-job-worker-svc` [→](https://github.com/rajansandeep/job-worker-svc/pull/3#discussion_r3112287477)

> Is there ever a situation where `serveReturned` is closed where the context isn't also canceled?

**@awly** on `renatoaguimaraes/golang-job-scheduler` [→](https://github.com/renatoaguimaraes/golang-job-scheduler/pull/6#discussion_r625936446)

> Enforce TLS 1.3 here, per your design doc

**@strideynet** on `rexposadas/teleport` [→](https://github.com/rexposadas/teleport/pull/1#discussion_r1790416804)

> I think we need a section here explaining what Authorization model will be used

**@rosstimothy** on `rexposadas/teleport` [→](https://github.com/rexposadas/teleport/pull/1#discussion_r1799134445)

> I don't see any changes to the document to address this yet.

**@rosstimothy** on `RichyHBM/teleport-challenge` [→](https://github.com/RichyHBM/teleport-challenge/pull/1#discussion_r2048907106)

> Slight clarification on the requirements: If a job has already completed users should still be able to retrieve it's output.

**@r0mant** on `samschurter/teleport-challenge` [→](https://github.com/samschurter/teleport-challenge/pull/2#discussion_r698889639)

> I didn't spot a deadline set on the context passed here anywhere so this potentially can be stuck for a long time.

**@smallinsky** on `supby/job-worker` [→](https://github.com/supby/job-worker/pull/1#discussion_r715442233)

> what is the purpose of the message field in context of StartResponse ?

**@smallinsky** on `supby/job-worker` [→](https://github.com/supby/job-worker/pull/1#discussion_r715442595)

> Same question. What is the purpose of the message field in context of StopResponse ?

**@espadolini** on `Zephan92/teleport-handson` [→](https://github.com/Zephan92/teleport-handson/pull/1#discussion_r2727149270)

> I would just document the requirement for Linux 5.14+.

**@nklaassen** on `Zephan92/teleport-handson` [→](https://github.com/Zephan92/teleport-handson/pull/3#discussion_r2737623200)

> avoid storing a context in a struct and prefer testing the external API https://go.dev/blog/context-and-structs
>
> i don't see it actually referenced in any tests either

**@espadolini** on `Zephan92/teleport-handson` [→](https://github.com/Zephan92/teleport-handson/pull/4#discussion_r2745701368)

> ```go
> <-ctx.Done()
> ...
> ```
>
> is better written out as a `context.AfterFunc`

**@nklaassen** on `Zephan92/teleport-handson` [→](https://github.com/Zephan92/teleport-handson/pull/4#discussion_r2748498178)

> ```suggestion
> 	ctx := t.Context()
> ```

**@espadolini** on `Zephan92/teleport-handson` [→](https://github.com/Zephan92/teleport-handson/pull/6#discussion_r2759972512)

> These execution functions are already taking a lot of parameters, I think I'd rather have them be a method of `Config` that takes context and output than take a `*cobra.Command` and everything else as individual parameters.

**@awly** on `ibeckermayer/teleport-interview` [→](https://github.com/ibeckermayer/teleport-interview/pull/3#discussion_r555981773)

> document what effect `env` has on the database client

**@alex-kovoy** on `zship/teleport-challenge` [→](https://github.com/zship/teleport-challenge/pull/1#discussion_r718764746)

> I think that JWT makes it a bit more complicated as it will be harder to implement a proper logout and enforcing 1h session TTL seems as a bad UX. I would revisit this design choice as part of the challenge is to implement a logout mechanism.

**@oeric** on `nick-hinds/teleport_challenge` [→](https://github.com/nick-hinds/teleport_challenge/pull/1#discussion_r2362774454)

> Can you provide more info on how you plan to implement this? (roughly) how many PRs, what they'll contain, etc.

**@russjones** on `nick-hinds/teleport_challenge` [→](https://github.com/nick-hinds/teleport_challenge/pull/1#discussion_r2366295079)

> Can you provide some more details here?

**@russjones** on `nick-hinds/teleport_challenge` [→](https://github.com/nick-hinds/teleport_challenge/pull/1#discussion_r2366295332)

> In general the structure looks fine, but can you provide some more details on what each one will do?

**@russjones** on `tedmist1/challenge-auth0` [→](https://github.com/tedmist1/challenge-auth0/pull/1#discussion_r2488244768)

> Can you share how you plan to start Grafana via Docker?

**@russjones** on `tedmist1/challenge-auth0` [→](https://github.com/tedmist1/challenge-auth0/pull/1#discussion_r2488252817)

> Can you provide 3-5 sentences on how you plan to configure the Auth0 tenant?

**@russjones** on `tedmist1/challenge-auth0` [→](https://github.com/tedmist1/challenge-auth0/pull/1#discussion_r2488253315)

> Can you provide 3-5 sentences of prose on how you plan to do this? For example, what endpoint do you intend to use?

---

## concurrency-and-locking

_Lock scoping and contention; sync.Once vs manual tracking; goroutine lifecycle; avoiding write-after-close and double-stop patterns._

**45 quotes** from `14` distinct reviewers across `21` candidate submissions.

**@rosstimothy** on `adalton/teleport-exercise` [→](https://github.com/adalton/teleport-exercise/pull/1#discussion_r767879798)

> You could use a `sync.Once` instead of tracking this manually

**@codingllama** on `adalton/teleport-exercise` [→](https://github.com/adalton/teleport-exercise/pull/1#discussion_r767962245)

> Suggestion: Merge this lockedOperation with the one below? It's easy to reason about the non-defer here, plus they are sequential.

**@codingllama** on `adalton/teleport-exercise` [→](https://github.com/adalton/teleport-exercise/pull/1#discussion_r768064312)

> Suggestion: print numbers in order, then make sure the goroutines see the expected output.

**@codingllama** on `adalton/teleport-exercise` [→](https://github.com/adalton/teleport-exercise/pull/1#discussion_r768707556)

> Yeah, I don't see the point for grabbing the lock twice. I'd either defer everything, or more likely push it all to the end, since we the function is short and doesn't have multiple exits.

**@dboslee** on `benmoss/job-worker-service` [→](https://github.com/benmoss/job-worker-service/pull/3#discussion_r2684545176)

> the entry mutex is never released here

**@espadolini** on `chintamanil/job-worker` [→](https://github.com/chintamanil/job-worker/pull/2#discussion_r2833514624)

> No need for a rwmutex if there's no need for multiple readers to block while holding the read lock and there's no contention as proven by a profile, sync.mutex is quite a bit lighter on the bookkeeping.

**@espadolini** on `chintamanil/job-worker` [→](https://github.com/chintamanil/job-worker/pull/2#discussion_r2833523494)

> No need to make the condvar a separate allocation
> ```suggestion
> 	cond      sync.Cond
> 	content   []byte
> 	closed    bool
> 	closeOnce sync.Once
> }
>
> // NewMemoryBuffer creates a new MemoryBuffer.
> func NewMemoryBuffer() *MemoryBuffer {
> 	mb := &MemoryBuffer{
> 		content: make([]byte, 0, 4096),
> 	}
> 	mb.cond.L = &mb.mu
> ```

**@codingllama** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/2#discussion_r1787743279)

> Wait for the s.Serve goroutine to gracefully end?

**@tigrato** on `DiscoRiver/jobber` [→](https://github.com/DiscoRiver/jobber/pull/1#discussion_r1989454559)

> Why using sync.Map and not a map protected by a mutex?

**@greedy52** on `ehsaniara/joblet` [→](https://github.com/ehsaniara/joblet/pull/2#discussion_r2121184613)

> what do you mean by `Sync pipe`? not sure how that avoids bypassing the limit.

**@awly** on `hashsequence/Linux-Job-Worker` [→](https://github.com/hashsequence/Linux-Job-Worker/pull/2#discussion_r478697452)

> Use a regular map with a mutex instead, per `sync.Map` docs:
> 
> > The Map type is specialized. Most code should use a plain Go map instead, with separate locking or coordination, for better type safety and to make it easier to maintain other invariants along with the map content.

**@rosstimothy** on `jaylane/job-scheduler` [→](https://github.com/jaylane/job-scheduler/pull/1#discussion_r1715313296)

> One of the requirements of the challenge is that the system can serve output to multiple concurrent clients, without dropping any output, or having two clients impact the others ability to retrieve job output.
> 
> Please consider the above when using channels to distribute process output.

**@codingllama** on `joshuarubin/teleport-job-worker` [→](https://github.com/joshuarubin/teleport-job-worker/pull/4#discussion_r1747636407)

> Is the sync.Once necessary here?

**@greedy52** on `kiakeshmiri/process-runner` [→](https://github.com/kiakeshmiri/process-runner/pull/1#discussion_r1766840951)

> +1 on the locking/sleeping concern.
> 
> however, this is really an implementation detail IMO (which should be leaved out from RFD ideally). Your description above (where each stream will have its own index/goroutine plus some protection on the log in memory) is good enough for my original question from design perspective.

**@rosstimothy** on `kovyrin/teleport-exec` [→](https://github.com/kovyrin/teleport-exec/pull/4#discussion_r794917457)

> I don't think you really need the channel and the goroutine here

**@rosstimothy** on `kovyrin/teleport-exec` [→](https://github.com/kovyrin/teleport-exec/pull/4#discussion_r795808495)

> I realize the convenience this provides, however it requires this function to acquire the lock twice in the event the request to disable should be processed.

**@zmb3** on `kovyrin/teleport-exec` [→](https://github.com/kovyrin/teleport-exec/pull/2#discussion_r795824562)

> This looks racy. You're taking out two separate and independent locks on the mutex.
> 
> `c.running` could be true on the initial check, and get set to false before you re-acquire the lock.

**@zmb3** on `kovyrin/teleport-exec` [→](https://github.com/kovyrin/teleport-exec/pull/2#discussion_r795841052)

> This is another case of you acquiring the lock multiple times and actually not getting an atomic operation out of it.
> 
> In this case, tail can be enabled on your first check, and then disabled by a concurrent call to `DisableTail()` before you re-acquire the mutex, resulting in a double close of `s.logComplete`, which will panic.

**@rosstimothy** on `MarkDHarris/teleport-systems-challenge` [→](https://github.com/MarkDHarris/teleport-systems-challenge/pull/2#discussion_r3035832497)

> This permits double close and alerts the readers on each subsequent close.

**@rosstimothy** on `MarkDHarris/teleport-systems-challenge` [→](https://github.com/MarkDHarris/teleport-systems-challenge/pull/2#discussion_r3035842086)

> Why do we need to acquire the lock again if we are only going to return in this branch? Won't that block all future attempts to acquire the lock?

**@rosstimothy** on `MarkDHarris/teleport-systems-challenge` [→](https://github.com/MarkDHarris/teleport-systems-challenge/pull/2#discussion_r3039649825)

> Ah apologies I missed the defer unlock above. This might be a sign that the locking employed by this function should be reworked though.

**@rosstimothy** on `MarkDHarris/teleport-systems-challenge` [→](https://github.com/MarkDHarris/teleport-systems-challenge/pull/2#discussion_r3045048915)

> Optional suggestion: We could modernize this code a bit.
> ```suggestion
> 		wg.Go(func() {
> 			id, err := m.CreateJob("mark", []string{"/bin/echo", fmt.Sprintf("concurrent-%d", i)})
> 			if err != nil {
> 				errs <- err
> 				return
> 			}
> 			ids <- id
> 		})
> ```

**@tigrato** on `mikewurtz/taskman` [→](https://github.com/mikewurtz/taskman/pull/4#discussion_r2042336791)

> id, clientid, pid, startTime and done do not need to be gated behind a mutex because these fields are immutable

**@rosstimothy** on `MrChristianL/Teleport-Job-Worker-Service` [→](https://github.com/MrChristianL/Teleport-Job-Worker-Service/pull/1#discussion_r2896145660)

> Does the mutex provide any value? What guarantees does Go make about exec.Cmd.Stdout and exe.Cmd.Stderr?

**@eriktate** on `MrChristianL/Teleport-Job-Worker-Service` [→](https://github.com/MrChristianL/Teleport-Job-Worker-Service/pull/2#discussion_r2906085923)

> Is there a reason not to `defer j.mu.Unlock()` here?

**@eriktate** on `MrChristianL/Teleport-Job-Worker-Service` [→](https://github.com/MrChristianL/Teleport-Job-Worker-Service/pull/2#discussion_r2906096388)

> Couldn't this be a `j.mu.RLock()`?

**@eriktate** on `MrChristianL/Teleport-Job-Worker-Service` [→](https://github.com/MrChristianL/Teleport-Job-Worker-Service/pull/2#discussion_r2906114409)

> Missing a lock here?

**@eriktate** on `MrChristianL/Teleport-Job-Worker-Service` [→](https://github.com/MrChristianL/Teleport-Job-Worker-Service/pull/2#discussion_r2906121587)

> No lock acquired while `j.cmd` is accessed.

**@nklaassen** on `MrChristianL/Teleport-Job-Worker-Service` [→](https://github.com/MrChristianL/Teleport-Job-Worker-Service/pull/2#discussion_r2906806646)

> could multiple concurrent calls to Stop be a problem?

**@Joerger** on `MrChristianL/Teleport-Job-Worker-Service` [→](https://github.com/MrChristianL/Teleport-Job-Worker-Service/pull/2#discussion_r2907519068)

> Do we expect much lock contention on the job to warrant the added algorithmic complexity of an RWMutex?

**@Joerger** on `MrChristianL/Teleport-Job-Worker-Service` [→](https://github.com/MrChristianL/Teleport-Job-Worker-Service/pull/2#discussion_r2915212689)

> If we don't hold the lock between checking j.cmd and setting j.cmd, then it's possible to get into a situation where `Start` is called twice and both get past this check. Should we avoid this?

**@rosstimothy** on `MrChristianL/Teleport-Job-Worker-Service` [→](https://github.com/MrChristianL/Teleport-Job-Worker-Service/pull/2#discussion_r2919425958)

> Does the lock need to be held while doing IO to uphold the guarantees mentioned in this comment?

**@rosstimothy** on `MrChristianL/Teleport-Job-Worker-Service` [→](https://github.com/MrChristianL/Teleport-Job-Worker-Service/pull/2#discussion_r2924618378)

> Would it be better to acquire and release the lock twice rather than holding it for the duration of the IO?

**@nklaassen** on `neildo/tjob` [→](https://github.com/neildo/tjob/pull/2#discussion_r1767829139)

> isn't this kinda racy, couldn't two concurrent calls to job.Start() both get past this line?

**@nklaassen** on `neildo/tjob` [→](https://github.com/neildo/tjob/pull/2#discussion_r1767832112)

> not very clear here which feels are/need to be protected by this mutex

**@smallinsky** on `P-A-R-U-S/Go-Job-Worker-Service` [→](https://github.com/P-A-R-U-S/Go-Job-Worker-Service/pull/2#discussion_r1743771396)

> Due to job.mutex.Lock()  call  In this case the job operation will be locked for 10 seconds.

**@creack** on `razzam21/job-worker-service` [→](https://github.com/razzam21/job-worker-service/pull/1#discussion_r2257562726)

> Just currious. Why use atomic.Value here and sync.RWMutex earlier?

**@r0mant** on `samschurter/teleport-challenge` [→](https://github.com/samschurter/teleport-challenge/pull/2#discussion_r698884795)

> This is being done under lock which means any client trying to schedule a new job will be stuck while the job is stopping (which may be a while).

**@zmb3** on `tjper/teleport` [→](https://github.com/tjper/teleport/pull/3#discussion_r830301460)

> This goroutine may leak. If the context is canceled before an error can be written to `errc` then it will just hang forever.

**@tigrato** on `xSudoNymx/job-manager` [→](https://github.com/xSudoNymx/job-manager/pull/3#discussion_r2728163326)

> requestStop is only invoked here so we can easily drop the lock unlock and the Stop holds the lock.

**@espadolini** on `Zephan92/teleport-handson` [→](https://github.com/Zephan92/teleport-handson/pull/1#discussion_r2718747005)

> Is `b.L` a different lock than the one embedded in `b`? If so, why?

**@espadolini** on `Zephan92/teleport-handson` [→](https://github.com/Zephan92/teleport-handson/pull/3#discussion_r2736609144)

> Does the benefit of RWMutex outweigh the significantly bigger overhead costs or the added complexity in logic?

**@espadolini** on `Zephan92/teleport-handson` [→](https://github.com/Zephan92/teleport-handson/pull/3#discussion_r2736689924)

> Does this have the potential to block while we're holding the job lock?

**@espadolini** on `Zephan92/teleport-handson` [→](https://github.com/Zephan92/teleport-handson/pull/4#discussion_r2745647772)

> Do we expect significant contention on this? I can't imagine having more than a handful of concurrent readers for any given output.

**@espadolini** on `Zephan92/teleport-handson` [→](https://github.com/Zephan92/teleport-handson/pull/4#discussion_r2745750869)

> Call this `flushEntriesRLocked` or something, to signify that the method must only be called while holding the read lock.

---

## implementation-details

_Specific questions about implementation choices, unused code paths, field purposes, and error handling._

**23 quotes** from `12` distinct reviewers across `11` candidate submissions.

**@tigrato** on `mcampo84/teleport_challenge` [→](https://github.com/mcampo84/teleport_challenge/pull/1#discussion_r1907296445)

> PIDs are useful, but they can be reused. For job metadata that persists beyond program execution, relying on them isn't ideal. Alternatives like UUIDs are better suited for this purpose.

**@rosstimothy** on `mikewurtz/taskman` [→](https://github.com/mikewurtz/taskman/pull/1#discussion_r2027213960)

> Is the full path actually required? What would happen if I ran `taskman start --user-id 5 "ls -lah /home/alice"`?

**@zmb3** on `mikewurtz/taskman` [→](https://github.com/mikewurtz/taskman/pull/2#discussion_r2036271488)

> This isn't performance-critical code, but worth pointing out that Sprintf is going to be significantly slower than the `+` operator for simple concatenation like this.

**@rosstimothy** on `mikewurtz/taskman` [→](https://github.com/mikewurtz/taskman/pull/5#discussion_r2044585912)

> What would happen if Close was called more than once?

**@tigrato** on `mikewurtz/taskman` [→](https://github.com/mikewurtz/taskman/pull/5#discussion_r2044888307)

> What happens if the client closes the connection and no-new data is written to the output since then? 
> Will the session terminate, will the goroutine become locked until a new write occurs?

**@eriktate** on `MrChristianL/Teleport-Job-Worker-Service` [→](https://github.com/MrChristianL/Teleport-Job-Worker-Service/pull/2#discussion_r2905724973)

> Could this be represented as the condition of the outer loop? Either way the continue below isn't really doing anything, I'd suggest removing it

**@eriktate** on `MrChristianL/Teleport-Job-Worker-Service` [→](https://github.com/MrChristianL/Teleport-Job-Worker-Service/pull/2#discussion_r2905869642)

> What happens if all of the readers produce the same incorrect output?

**@rosstimothy** on `MrChristianL/Teleport-Job-Worker-Service` [→](https://github.com/MrChristianL/Teleport-Job-Worker-Service/pull/3#discussion_r2924714805)

> What server is this connecting to? What happens if something else is listening on that port which allows TLS 1.2?

**@hugoShaka** on `neildo/tjob` [→](https://github.com/neildo/tjob/pull/1#discussion_r1758852623)

> What does this field do?

**@nklaassen** on `neildo/tjob` [→](https://github.com/neildo/tjob/pull/1#discussion_r1759278508)

> why repeated? is this stdout or stderr?

**@nklaassen** on `neildo/tjob` [→](https://github.com/neildo/tjob/pull/2#discussion_r1767822702)

> why print, isn't this going to interfere with your job output?

**@nklaassen** on `neildo/tjob` [→](https://github.com/neildo/tjob/pull/2#discussion_r1767824221)

> why return this as an error?

**@nklaassen** on `neildo/tjob` [→](https://github.com/neildo/tjob/pull/2#discussion_r1767898316)

> what could this error be and what does it mean to return it from your JobReader?

**@rosstimothy** on `neildo/tjob` [→](https://github.com/neildo/tjob/pull/2#discussion_r1769119466)

> What happens if two goroutines call Start and Logs on a single job instance at the same time?

**@atburke** on `rajansandeep/teleport-job-worker-svc` [→](https://github.com/rajansandeep/job-worker-svc/pull/3#discussion_r3112510703)

> What does the extra function do here?

**@awly** on `renatoaguimaraes/golang-job-scheduler` [→](https://github.com/renatoaguimaraes/golang-job-scheduler/pull/2#discussion_r623314633)

> why is this cancellation wrapper needed?

**@awly** on `renatoaguimaraes/golang-job-scheduler` [→](https://github.com/renatoaguimaraes/golang-job-scheduler/pull/2#discussion_r623377512)

> why `uint32` and not `int`?

**@dboslee** on `renatoaguimaraes/golang-job-scheduler` [→](https://github.com/renatoaguimaraes/golang-job-scheduler/pull/2#discussion_r623384851)

> why not use job.Cmd.Process?

**@tigrato** on `RichyHBM/teleport-challenge` [→](https://github.com/RichyHBM/teleport-challenge/pull/5#discussion_r2066764421)

> what happens if the write operation takes ages because the consumer is slow?

**@zmb3** on `shawon-crosen/teleport-challenge` [→](https://github.com/shawon-crosen/teleport-challenge/pull/1#discussion_r2129154241)

> What is this `user` field?

**@codingllama** on `tjper/teleport` [→](https://github.com/tjper/teleport/pull/2#discussion_r827449254)

> Why do we need this functionality? Let's remove it if we don't.

**@greedy52** on `xSudoNymx/job-manager` [→](https://github.com/xSudoNymx/job-manager/pull/1#discussion_r2709762719)

> why not use the official Timestamp type?

**@espadolini** on `Zephan92/teleport-handson` [→](https://github.com/Zephan92/teleport-handson/pull/3#discussion_r2736605466)

> Why can't this be done as part of the execution of the server?

---

## reproducibility-and-tooling

_Build reproducibility; buf/protoc/golangci-lint versions; committed .pb.go files; Docker vs pure Go implementations._

**19 quotes** from `9` distinct reviewers across `12` candidate submissions.

**@codingllama** on `adalton/teleport-exercise` [→](https://github.com/adalton/teleport-exercise/pull/2#discussion_r771447049)

> If you want to distinguish unset/default from STDOUT, better start tags at 1.
> 
> https://developers.google.com/protocol-buffers/docs/style#enums

**@nklaassen** on `benmoss/job-worker-service` [→](https://github.com/benmoss/job-worker-service/pull/3#discussion_r2684012549)

> It's really not obvious to me why you set r.lastVersion to `version + 1` and not just `version`

**@rosstimothy** on `clydotron/job_worker_service` [→](https://github.com/clydotron/job_worker_service/pull/1#discussion_r1599969091)

> While a docker container works, we want to see you write as much of your own code as possible. Please drop docker in favor of interacting with cgroups yourself directly in the library.

**@rosstimothy** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/1#discussion_r1786191858)

> Alternatively there is also [buf](https://buf.build/)

**@codingllama** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/1#discussion_r1786298692)

> My 2c: linters are great but you can just eyeball the style for the challenge, if that saves you time.

**@espadolini** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/2#discussion_r1787701881)

> (then again, my mind has also been irreparably ruined by the elegance of `buf generate` so maybe my personal opinions on the matter shouldn't really be considered)

**@codingllama** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/2#discussion_r1787749339)

> nit: Don't silence commands (other than echo and the likes), this hides the doings of the Makefile. `make -s` can be used for a true silent mode.
> 
> (No need to follow up on this everywhere, the time spent is not worth it.)

**@codingllama** on `creack/telepilot` [→](https://github.com/creack/telepilot/pull/2#discussion_r1787753277)

> This is an interesting way to get this setup consistently. 👀
> 
> I would likely go native with golangci-lint, as in my experience it's super slow on docker, but everything else seems fine. (Arguably Buf does a bunch to help too.)
> 
> Another idea I toyed with is having a bunch of shims under a "bin/" folder (bin/buf, bin/golangci-lint, etc) and have those hold the docker incantations (or whatever is used). I find it a bit simpler than all the Makefile "build" / "run" rules. The root Makefile can add $PWD/bin to the $PATH or use those directly. (I'm also vaguely remember some external tooling that helps with the shims.)

**@rosstimothy** on `GevorgGal/jobworker` [→](https://github.com/GevorgGal/jobworker/pull/3#discussion_r3002859755)

> Should we pin to a specific version of protoc? I'm seeing a diff locally because I have a different version of protoc than the one that was used to generated the files committed in this PR.

**@rosstimothy** on `GevorgGal/jobworker` [→](https://github.com/GevorgGal/jobworker/pull/3#discussion_r3003175600)

> I asked because I had 25+ failures when I ran this locally since I have a config file in one of the places golangci-lint looks for one. Even if you only add a config file with default values, that would be better than the caller being subject to random failures.

**@rosstimothy** on `GevorgGal/jobworker` [→](https://github.com/GevorgGal/jobworker/pull/3#discussion_r3003320394)

> Could we switch to buf now for reproducible generation?

**@strideynet** on `gstelang/job-worker-service` [→](https://github.com/gstelang/job-worker-service/pull/4#discussion_r1740683343)

> Using `bufio.Scanner` means you'll be handling output on a line by line basis. I'm not sure this is a good assumption - some jobs could output a large amount of data within a single line and there's no requirement that what they are outputting might be clear line-by-line logs.

**@andreiko** on `ilyazz/jobs` [→](https://github.com/ilyazz/jobs/pull/4#discussion_r1007435643)

> Are the files supposed to be generated on every build or committed to git? Would be nice to have something that clarifies that.

**@AntonAM** on `jaylane/job-scheduler` [→](https://github.com/jaylane/job-scheduler/pull/1#discussion_r1714301890)

> Maybe a bit more details, like alogs and reasons for choosing this version?

**@codingllama** on `joshuarubin/teleport-job-worker` [→](https://github.com/joshuarubin/teleport-job-worker/pull/4#discussion_r1747606909)

> nit: Write a proper nolint comment. There's not much point in using nolintlint if we ignore it.

**@codingllama** on `joshuarubin/teleport-job-worker` [→](https://github.com/joshuarubin/teleport-job-worker/pull/4#discussion_r1747638720)

> Suggestion: This file is all file-manipulation, so while the logic is linux-specific you could drop the build constraint, save yourself a couple of `//nonlints` and even test it regardless of platform.

**@greedy52** on `kiakeshmiri/process-runner` [→](https://github.com/kiakeshmiri/process-runner/pull/1#discussion_r1766866138)

> not very concerned on which tool will be used. but your sample below do answer most things I was looking for like a CA, and some CSR values.

**@zmb3** on `kovyrin/teleport-exec` [→](https://github.com/kovyrin/teleport-exec/pull/2#discussion_r795835640)

> I would use the tools in the `os` package for making temporary dirs/files. This will make sure you leverage the appropriate temp directory for the system.

**@codingllama** on `tjper/teleport` [→](https://github.com/tjper/teleport/pull/1#discussion_r822072824)

> That's OK, we don't need all of that in design. I was mostly looking for algorithm/key length.

---

## command-invocation

_Questions about command and argument parsing, shell interpretation, and data transport._

**16 quotes** from `8` distinct reviewers across `7` candidate submissions.

**@greedy52** on `m3talsmith/jobberthehut` [→](https://github.com/m3talsmith/jobberthehut/pull/1#discussion_r2499188566)

> does `command` include args?

**@rosstimothy** on `m3talsmith/jobberthehut` [→](https://github.com/m3talsmith/jobberthehut/pull/1#discussion_r2500026851)

> How are arguments to the command provided? Are they included in this string? What problems might that present?

**@rosstimothy** on `m3talsmith/jobberthehut` [→](https://github.com/m3talsmith/jobberthehut/pull/1#discussion_r2500043611)

> Is a string capable of transporting binary data?

**@tigrato** on `mikewurtz/taskman` [→](https://github.com/mikewurtz/taskman/pull/1#discussion_r2027167697)

> Is there a way of using the command without quotes?
> having to quote the command + args seems a bit annoying 😆

**@timothyb89** on `mikewurtz/taskman` [→](https://github.com/mikewurtz/taskman/pull/1#discussion_r2028829559)

> I think this is still a bit ambiguous, since some rule is needed to decide whether or not to use shell interpretation (and which shell, etc).
> 
> I'd argue the simpler and better solution might be to let the user's own shell do the parsing, and only support the 2nd case. If the user wants to run a shell command, they can just use e.g. `sh -c 'command with spaces'`

**@awly** on `renatoaguimaraes/golang-job-scheduler` [→](https://github.com/renatoaguimaraes/golang-job-scheduler/pull/2#discussion_r623381461)

> If your command is `sh -c "foo | bar"`, pipes should work
> ```suggestion
> // Start runs a Linux single command with arguments.
> ```

**@rosstimothy** on `RichyHBM/teleport-challenge` [→](https://github.com/RichyHBM/teleport-challenge/pull/1#discussion_r2048947465)

> How are any arguments provided?

**@rosstimothy** on `RichyHBM/teleport-challenge` [→](https://github.com/RichyHBM/teleport-challenge/pull/1#discussion_r2048947623)

> How are any arguments provided?

**@rosstimothy** on `RichyHBM/teleport-challenge` [→](https://github.com/RichyHBM/teleport-challenge/pull/1#discussion_r2048954046)

> Is a string capable of transporting non-UTF 8 data?

**@eriktate** on `RichyHBM/teleport-challenge` [→](https://github.com/RichyHBM/teleport-challenge/pull/1#discussion_r2051223889)

> Would splitting the command from the arguments reduce this risk at all?

**@eriktate** on `RichyHBM/teleport-challenge` [→](https://github.com/RichyHBM/teleport-challenge/pull/1#discussion_r2052834593)

> I agree we would probably want to do more if command injection were a larger concern, but for this challenge your approach sounds reasonable to me

**@rosstimothy** on `shawon-crosen/teleport-challenge` [→](https://github.com/shawon-crosen/teleport-challenge/pull/1#discussion_r2130920501)

> Does this include the binary and any arguments that may be passed to it?

**@smallinsky** on `supby/job-worker` [→](https://github.com/supby/job-worker/pull/1#discussion_r716420256)

> it would be nice to include also command argument here.

**@rosstimothy** on `tjper/teleport` [→](https://github.com/tjper/teleport/pull/1#discussion_r821057464)

> How do I provide arguments to a command?

**@rosstimothy** on `tjper/teleport` [→](https://github.com/tjper/teleport/pull/1#discussion_r821060978)

> How do I pass arguments to a command?

**@codingllama** on `tjper/teleport` [→](https://github.com/tjper/teleport/pull/3#discussion_r829327706)

> My 2c: Make FetchJob (and all other funcs) take the ID as a string. The underlying format doesn't matter for API users, plus we get to ditch the uuid.Parse above.

---

## code-organization

_Code structure, package layout, architectural patterns, and Go idioms (goroutines, context, flags, locks, clock abstraction)._

**16 quotes** from `4` distinct reviewers across `3` candidate submissions.

**@russjones** on `atburke/teleport_interview` [→](https://github.com/atburke/teleport_interview/pull/1#discussion_r645899241)

> This logic should be internal to your database and not exposed to the caller. So start the goroutine in `NewMySqlDatabase`.

**@r0mant** on `atburke/teleport_interview` [→](https://github.com/atburke/teleport_interview/pull/1#discussion_r645900821)

> Nit: Move this closer to where it's used.

**@r0mant** on `atburke/teleport_interview` [→](https://github.com/atburke/teleport_interview/pull/1#discussion_r645900863)

> Nit: Move this closer to where it's used.

**@r0mant** on `atburke/teleport_interview` [→](https://github.com/atburke/teleport_interview/pull/1#discussion_r645901853)

> Pass context to Run method and do a select on this and context.Done() channel - this way this goroutine will be cancelable.

**@awly** on `ibeckermayer/teleport-interview` [→](https://github.com/ibeckermayer/teleport-interview/pull/2#discussion_r552920723)

> nit: extract this into a function, so that the fake authentication logic is contained and separate from normal code

**@awly** on `ibeckermayer/teleport-interview` [→](https://github.com/ibeckermayer/teleport-interview/pull/2#discussion_r552924196)

> nit: you don't need to use `malformedRequest` as a pointer everywhere

**@russjones** on `ibeckermayer/teleport-interview` [→](https://github.com/ibeckermayer/teleport-interview/pull/2#discussion_r553699892)

> I don't think this works, if I am not mistaken you have to call `flag.Parse()` before using the variables.

**@awly** on `ibeckermayer/teleport-interview` [→](https://github.com/ibeckermayer/teleport-interview/pull/3#discussion_r555978145)

> you don't need to lock since `sm.timeout` is never modified

**@awly** on `ibeckermayer/teleport-interview` [→](https://github.com/ibeckermayer/teleport-interview/pull/3#discussion_r555980670)

> nit: might be simpler to just export the `env` field in `Config` and get rid of the constructor

**@awly** on `ibeckermayer/teleport-interview` [→](https://github.com/ibeckermayer/teleport-interview/pull/3#discussion_r555985011)

> nit: create constants for these plan names, to prevent typos

**@awly** on `ibeckermayer/teleport-interview` [→](https://github.com/ibeckermayer/teleport-interview/pull/3#discussion_r555985534)

> nit: this constructor func seems redundant, if you export the fields of `Config`

**@awly** on `ibeckermayer/teleport-interview` [→](https://github.com/ibeckermayer/teleport-interview/pull/3#discussion_r556140087)

> ```suggestion
> 	mtx     sync.Mutex
> ```

**@awly** on `ibeckermayer/teleport-interview` [→](https://github.com/ibeckermayer/teleport-interview/pull/3#discussion_r556140327)

> nit: put the comment above the field, not on the same line
> also, wrap the long comment over multiple lines

**@awly** on `ibeckermayer/teleport-interview` [→](https://github.com/ibeckermayer/teleport-interview/pull/5#discussion_r557577199)

> nit: all of this should happen under a single write lock
> if you have 2 concurrent `DeleteSession` calls, they might intertwine their read/write locks and return the wrong bool.

**@awly** on `ibeckermayer/teleport-interview` [→](https://github.com/ibeckermayer/teleport-interview/pull/5#discussion_r557579849)

> nit: this check should be `splitToken := strings.Fields(reqToken); if len(splitToken) != 2 || splitToken[0] != "Bearer" { ...`, to avoid other invalid values like `Bearer foo bar` for example

**@alex-kovoy** on `tobocop/go-teleport-directory-browser` [→](https://github.com/tobocop/go-teleport-directory-browser/pull/2#discussion_r698897152)

> will there be a separate context provider for app level methods? if so, why not just call it `AppContextProvider` ?

---

## auth-storage

_Session token and cookie storage mechanisms, including security considerations around localStorage vs httpOnly cookies._

**15 quotes** from `3` distinct reviewers across `3` candidate submissions.

**@alex-kovoy** on `ibeckermayer/teleport-interview` [→](https://github.com/ibeckermayer/teleport-interview/pull/1#discussion_r551412543)

> >  its default `CookieStore` implementation apparently uses a [non-OWASP-compliant](https://curtisvermeeren.github.io/2018/05/13/Golang-Gorilla-Sessions.html) storage mechanism.
> 
> Just wonder what makes it non OWASP compliant ?

**@alex-kovoy** on `ibeckermayer/teleport-interview` [→](https://github.com/ibeckermayer/teleport-interview/pull/1#discussion_r551430598)

> I think you can implement the Use double submit cookies mechanism which would be simpler to implement. Also, could you specify how you are going to pass a CSRF token to the web-client?

**@alex-kovoy** on `ibeckermayer/teleport-interview` [→](https://github.com/ibeckermayer/teleport-interview/pull/1#discussion_r551584588)

> What to store in the cookie is optional, it does not look like that package enforces it. The blog post, it seems, says that if you do not want to store session information in the DB, you can store it in the cookie if your requirements allow it but at the end of the day it's up to you.

**@alex-kovoy** on `ibeckermayer/teleport-interview` [→](https://github.com/ibeckermayer/teleport-interview/pull/1#discussion_r551611467)

> > Using `localStorage` does increase the attack surface for XSS attacks in general
> 
> I would not say that it increases the attack surface because XSS targets the rendering layer, un-sanitized inputs, etc. If an attacker managed to run malicious JS in the victim browser then it's already too late.

**@alex-kovoy** on `ibeckermayer/teleport-interview` [→](https://github.com/ibeckermayer/teleport-interview/pull/3#discussion_r555538027)

> I like your solution with reducer but I think it could be simplified a bit (at least for this app). For example, we can do something like this:
> 
> ```tsx
> function AppContext (props) {
>   const [ state, setState ] = React.useState(() => {
>     const session = localStorage.getItem('session');
>     return {
>       session
>     }    
>   })
> 
>   const setSession = session => {
>     localStorage.setItem('session', JSON.stringify(session);
>     setState({
>       ...state,
>       session
>     })
>   }
> 
>   const clearSession = () => {
>     setSession(null)    
>   }
>   
>   const ctx = {
>     ...state,
>     setSession,
>     clearSession
>   }
> 
>   return (
>     <StoreContext.Provider value={ctx} children={props.children}/>      
>   );
> }
> 
> ```
> Also, you may want to create a file `session.ts` with 2 helper methods, so you can get the token outside from virtual DOM (like your api.ts).
> 
> ```ts
> export function setSession(session){
>   localStorage.setItem('session', JSON.stringify(session);
> }
> 
> export function getSession(){
>   return JSON.parse(localStorage.getItem('session'))
> }
> ```
> 
> But again, your reducer looks good enough to me.

**@alex-kovoy** on `ibeckermayer/teleport-interview` [→](https://github.com/ibeckermayer/teleport-interview/pull/3#discussion_r556046344)

> > Unless I'm missing something, I don't think I'll need setSession
> 
> This is to keep localStorage getters/setters methods in one place for tracking purposes (like item key names, etc). Then you can call it from within AppContext.setSession (instead of referencing the localStorage directly).

**@alex-kovoy** on `ibeckermayer/teleport-interview` [→](https://github.com/ibeckermayer/teleport-interview/pull/3#discussion_r556092626)

> > Also, you may want to create a file session.ts 
> 
> (correction) instead of creating `sesion.ts`, you can create `localStorage.ts` and put these methods in there. So all window.localStorage setter/getters are in 1 place.

**@alex-kovoy** on `ibeckermayer/teleport-interview` [→](https://github.com/ibeckermayer/teleport-interview/pull/4#discussion_r556205464)

> Our previous discussion was around Application level storage. If you are not planning to use it for anything else but the session, in this case you might not need any of it  (React Context/Provider/Consumer). Below is an example of 3 files (localStorage.js, Login.js, and Authenticated.js)
> 
> ```js
> // localStorage.ts
> const localStorage = {
>   setSession(value){
>     window.localStorage.setItem('session', JSON.stringify(value));
>   }
> 
>   getSession(){
>     return JSON.parse(window.localStorage.getItem('session'));
>   }
> }
> 
> export default localStorage;
> 
> 
> 
> // login.ts
> import localStorage from './localStorage';
> 
> // ...
> 
> const tryLogin = async e => {    
>   e.preventDefault();
>   try {
>     const session = await api.post('/login', {
>       email: e.target.email.value,
>       password: e.target.password.value,
>     });
> 
>     // will add a session to local storage
>     localStorage.setSession(session);    
>     
>     // will trigger a re-render
>     history.push('/dashboard');      
>   } catch (error) {
>     //
>   }
> };
> 
> // Authenticated.js
> import localStorage from './localStorage';
> 
> export default function Authenticated(props) {
>   const [authenticated] = React.useState(() => !!localStorage.getSession()));
> 
>   React.useEffect(() => {
>     if(!authenticated){
>       // return to the login screen. This will trigger a global re-render
>       history.push('/login')
>     }
>   }, []);
> 
> 
>   if (!authenticated) {
>      return null;
>   }
> 
>   return props.children;
> }
> ```

**@alex-kovoy** on `ibeckermayer/teleport-interview` [→](https://github.com/ibeckermayer/teleport-interview/pull/4#discussion_r556242177)

> > maybe just go back to calling it the more general "store"
> 
> Yeah, constraining it to a session was a bit confusing. AppStore|AppContext sounds better if you are planing to use it for app level storage.
> 
> > a mutation to the global state, and then each component "subscribes" to pieces of the global state relevant to them, and automatically reacts to changes 
> 
> Got it. I am familiar with this approach and used it in the past as well although now I prefer hooks with a minimum global state (when it is possible of course)
> 
> > Or do you still want me to change it to something closer to your example?
> 
> Up to you but if you want to set up a global state I would just name it accordingly.

**@alex-kovoy** on `tobocop/go-teleport-directory-browser` [→](https://github.com/tobocop/go-teleport-directory-browser/pull/3#discussion_r699745227)

> This state is not transferable between multiple browser tabs. Right now it works only in 1 browser tab.

**@alex-kovoy** on `tobocop/go-teleport-directory-browser` [→](https://github.com/tobocop/go-teleport-directory-browser/pull/3#discussion_r699748694)

> i am not sure that this is CSRF double cookie implementation, from the first glace it does not look it does anything that mitigates CSRF.

**@russjones** on `tobocop/go-teleport-directory-browser` [→](https://github.com/tobocop/go-teleport-directory-browser/pull/3#discussion_r699783424)

> Typically the server will generate and send the cookie to the client, then the client submits that cookie. I'll look at your submission in more detail tomorrow morning.

**@alex-kovoy** on `tobocop/go-teleport-directory-browser` [→](https://github.com/tobocop/go-teleport-directory-browser/pull/3#discussion_r699785260)

> The cookie needs to be set by the server and client needs to send it back with each http request. Current solution wont work with sub-domains because it's possible to set a cookie from javascript for parent domains. For example, a user can visit a subdomain website which will send API requests to your endpoints on user behalf.

**@r0mant** on `zship/teleport-challenge` [→](https://github.com/zship/teleport-challenge/pull/1#discussion_r718769353)

> Are there any downsides to using local storage for the access token vs cookies?

**@alex-kovoy** on `zship/teleport-challenge` [→](https://github.com/zship/teleport-challenge/pull/5#discussion_r726714332)

> it could be more readable if usage of localStorage could be located in one place. Right now localStorage is referenced in different files (actions and store) which makes it harder to see where the source of truth is. It feels that you are trying to follow Redux patterns which looks unnecessary here (both server and client side). I suggest to put actions and store handling in 1 file and rename it to `useLocalStorage`.

---

## auth0-least-privilege

_Auth0 scopes, OAuth permissions, MFA factors, and adherence to principle of least privilege._

**10 quotes** from `4` distinct reviewers across `5` candidate submissions.

**@russjones** on `ibeckermayer/teleport-interview` [→](https://github.com/ibeckermayer/teleport-interview/pull/2#discussion_r552310832)

> I wouldn't include this logic in `main`. Instead create a package that has your server in it and start that instead. Ideally your `main` would just be flag parsing and then calling that package to start the server. Maybe doing some signal handling but I think that's out of scope for your task.

**@alex-kovoy** on `zship/teleport-challenge` [→](https://github.com/zship/teleport-challenge/pull/4#discussion_r725679587)

> I find this logic a bit hard to follow where `getActions` returns a set of actions to execute (?). The caller in this case goes through all actions and decides which one to run. This makes it unclear who is responsible for managing the session state `getActions` method or a `caller`?  What if `getActions` returns 2 actions for a session: `delete` and `update`?

**@oeric** on `nick-hinds/teleport_challenge` [→](https://github.com/nick-hinds/teleport_challenge/pull/1#discussion_r2362754219)

> Can you expand on which scopes in particular you'd need?

**@oeric** on `nick-hinds/teleport_challenge` [→](https://github.com/nick-hinds/teleport_challenge/pull/1#discussion_r2362760354)

> Can you provide more info on which MFA factors you'd allow/configure, and your reasoning for them?

**@oeric** on `nick-hinds/teleport_challenge` [→](https://github.com/nick-hinds/teleport_challenge/pull/1#discussion_r2362768147)

> What Github Actions would you be building? Please provide more info (what the triggers are, what they'd be doing, etc.)

**@russjones** on `nick-hinds/teleport_challenge` [→](https://github.com/nick-hinds/teleport_challenge/pull/1#discussion_r2366293683)

> This is a really old version of the Auth0 Terraform Provider. Any reason you're using such an old version?

**@russjones** on `nick-hinds/teleport_challenge` [→](https://github.com/nick-hinds/teleport_challenge/pull/1#discussion_r2366293900)

> Similar to the Auth0 Terraform Provider, this is a really old version. Any reason you're using an old version of this package?

**@russjones** on `nick-hinds/teleport_challenge` [→](https://github.com/nick-hinds/teleport_challenge/pull/1#discussion_r2366294882)

> Why do you need access to `*:actions`? Maybe I am misunderstanding something but Auth0 actions are not needed for this challenge? https://auth0.com/docs/customize/actions

**@russjones** on `tedmist1/challenge-auth0` [→](https://github.com/tedmist1/challenge-auth0/pull/1#discussion_r2492511341)

> > * Register a sample user in the tenant's database via HTTP POST request to Auth0 Management API endpoint. (`POST https://{auth0_domain}/api/v2/users`).
> 
> How do you intend to use this endpoint. It looks like in some configuration an initial password is required.
> 
> https://auth0.com/docs/api/management/v2/users/post-users

**@rcanderson23** on `Chili-Man/teleport-sre-challenge` [→](https://github.com/Chili-Man/teleport-sre-challenge/pull/1#discussion_r1416166814)

> What RBAC permissions are required for your application to work correctly in the kubernetes cluster?

---

## tls-and-hashing

_Password hashing (bcrypt, argon2), HMAC handling, and timing attack prevention._

**9 quotes** from `3` distinct reviewers across `3` candidate submissions.

**@alex-kovoy** on `ibeckermayer/teleport-interview` [→](https://github.com/ibeckermayer/teleport-interview/pull/1#discussion_r551437099)

> How are you planning to store and verify passwords?

**@awly** on `ibeckermayer/teleport-interview` [→](https://github.com/ibeckermayer/teleport-interview/pull/3#discussion_r555982667)

> ```suggestion
> 	if !auth.CheckPasswordHash(body.Password, account.PasswordHash) {
> ```

**@alex-kovoy** on `tobocop/go-teleport-directory-browser` [→](https://github.com/tobocop/go-teleport-directory-browser/pull/2#discussion_r698906068)

> what does this method do? Should it use stored password hash to compare?

**@russjones** on `tobocop/go-teleport-directory-browser` [→](https://github.com/tobocop/go-teleport-directory-browser/pull/2#discussion_r698907880)

> Why are you passing a `hmac(password)` to bcrypt?

**@russjones** on `tobocop/go-teleport-directory-browser` [→](https://github.com/tobocop/go-teleport-directory-browser/pull/2#discussion_r699401010)

> I'm still not sure I understand the need for the hmac?
> 
> The output of `bcrypt.GenerateFromPassword` is a string that looks like the following which embeds the salt + hash (and other information) in it. You can pre-generate this and just check it into source.
> 
> ```
> $2a$10$J2dMyKwK4188qqjTnk/vkurcdEUMRt/Bq0zoa6isxdsxHWKIPFoyu
> ```
> 
> Then all you have to do is call `bcrypt.CompareHashAndPassword` with the above and the passed in password.
> 
> Using bcrypt alone this way, I don't see how a plaintext password would be stored? Maybe I am missing something?

**@russjones** on `tobocop/go-teleport-directory-browser` [→](https://github.com/tobocop/go-teleport-directory-browser/pull/2#discussion_r699418418)

> The bcrypt hash is not just the hash, it's the salt + hash. When you do `bcrypt.CompareHashAndPassword` it just extracts the salt from the passed in bcrypt hash (this thing: `$2a$10$J2dMyKwK4188qqjTnk/vkurcdEUMRt/Bq0zoa6isxdsxHWKIPFoyu`), then uses it to recompute the hash with the extract salt + plaintext password, and then do a constant time compare.

**@russjones** on `tobocop/go-teleport-directory-browser` [→](https://github.com/tobocop/go-teleport-directory-browser/pull/3#discussion_r699778279)

> I would use https://pkg.go.dev/crypto/subtle#ConstantTimeCompare here, otherwise you might open yourself up to timing attacks.

**@alex-kovoy** on `tobocop/go-teleport-directory-browser` [→](https://github.com/tobocop/go-teleport-directory-browser/pull/3#discussion_r699789890)

> `found, passwordHash` is more precise

**@russjones** on `zship/teleport-challenge` [→](https://github.com/zship/teleport-challenge/pull/1#discussion_r719879201)

> You don't need to do all of this for this challenge, just using a KDF like bcrypt is fine.

---

## input-validation

_Input sanitization, validation, prevention of open redirects, and escape handling._

**8 quotes** from `4` distinct reviewers across `3` candidate submissions.

**@awly** on `ibeckermayer/teleport-interview` [→](https://github.com/ibeckermayer/teleport-interview/pull/5#discussion_r557577645)

> should logout redirect to the login page?

**@r0mant** on `tobocop/go-teleport-directory-browser` [→](https://github.com/tobocop/go-teleport-directory-browser/pull/2#discussion_r698879387)

> Why do username/password need to be escaped if they're sent in the request body?

**@alex-kovoy** on `zship/teleport-challenge` [→](https://github.com/zship/teleport-challenge/pull/1#discussion_r718776929)

> how redirects will be handled back to original URL (after a user has been successfully authenticated) ?

**@alex-kovoy** on `zship/teleport-challenge` [→](https://github.com/zship/teleport-challenge/pull/1#discussion_r718777350)

> Are there going to be any input validation regarding these requests?

**@russjones** on `zship/teleport-challenge` [→](https://github.com/zship/teleport-challenge/pull/1#discussion_r719878767)

> What about open redirect?

**@alex-kovoy** on `zship/teleport-challenge` [→](https://github.com/zship/teleport-challenge/pull/1#discussion_r719884952)

> hm...why `/` need to be encoded in the `redirect` parameter?

**@alex-kovoy** on `zship/teleport-challenge` [→](https://github.com/zship/teleport-challenge/pull/4#discussion_r725680512)

> This message is a bit misleading (as input validation message) as it's used for `empty` passwords as well

**@alex-kovoy** on `zship/teleport-challenge` [→](https://github.com/zship/teleport-challenge/pull/5#discussion_r726701650)

> by following you comment it's unclear why `isValidUrl` is used here.  Would `if(redirect.startsWith('/') {  ... } ` be enough?

---

## setup-configuration

_Details on certificate generation, configuration, addressing, and hardcoded vs dynamic values._

**5 quotes** from `3` distinct reviewers across `4` candidate submissions.

**@greedy52** on `m3talsmith/jobberthehut` [→](https://github.com/m3talsmith/jobberthehut/pull/1#discussion_r2499163907)

> how is the client configured?

**@rosstimothy** on `m3talsmith/jobberthehut` [→](https://github.com/m3talsmith/jobberthehut/pull/1#discussion_r2505259609)

> I would suggest to reduce complexity by making this a unary look up.

**@rosstimothy** on `neildo/tjob` [→](https://github.com/neildo/tjob/pull/1#discussion_r1759124477)

> Should there be an option to allow overriding the default server address too?

**@rosstimothy** on `rexposadas/teleport` [→](https://github.com/rexposadas/teleport/pull/1#discussion_r1790653267)

> What will the hardcoded values be?

**@tigrato** on `rsteinkeXJ/teleport_interview` [→](https://github.com/rsteinkeXJ/teleport_interview/pull/1#discussion_r3123737322)

> feel free to hardcode these values server side. no need to accept them from the client

---

## csrf-protection

_CSRF protection mechanisms, double submit cookies, and token-based approaches._

**2 quotes** from `2` distinct reviewers across `2` candidate submissions.

**@r0mant** on `atburke/teleport_interview` [→](https://github.com/atburke/teleport_interview/pull/1#discussion_r645902542)

> Yeah, I think if you move this logic for checking csrf token, retrieving and checking the sesion etc, that would make it a bit cleaner.

**@alex-kovoy** on `ibeckermayer/teleport-interview` [→](https://github.com/ibeckermayer/teleport-interview/pull/1#discussion_r551425607)

> > Mitigates CSRF attacks
> 
> It would prevent only _some_ CSRF attacks, so it should not be considered as a safe protection mechanism alone. For example, if some GET endpoint changes the user state then an attacker can use one-click URL to carry an attack.

---

