# hasher
A simple application to calculate many types of hashes concurrently, mostly from Go's standard library.

Normally you'll want to hash files, but you can also feed data from the standard input and calculate one single hash from there.

Hashes are calculated concurrently using a simple pipeline concurrency pattern. "Workers" as used here are not exactly the typical multiprocessing worker, but more like in how many pipelines will the workload be distributed. For more info, check [here](https://blog.golang.org/pipelines). Or don't. If you're a user you shouldn't care about implementation :) .

To ask for help, simply run

```
hasher
```

or

```
hasher -help
```
