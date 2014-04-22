Go Worker for Disco
===================

This is an implementation of [Disco worker protocol](http://disco.readthedocs.org/en/latest/howto/worker.html) in golang.
See [discoproject.org] (http://discoproject.org) for more information.

To use this worker install the three libraries via `go get`

```
$ go get github.com/dinedal/goworker/jobutil
$ go get github.com/dinedal/goworker/jobpack
$ go get github.com/dinedal/goworker/worker
```

Then, use the `jobpack` command to run workers. An example:
```
$ jobpack localhost:8989  $GOPATH/github.com/dinedal/goworker/examples/ http://discoproject.org/media/text/chekhov.txt
```

To submit a job to Disco master.

Warning: This is a work in progress and it is not ready for production use.

This implementation requires golang v1.1 or later.

Build Status: [Travis-CI](http://travis-ci.org/pooya/goworker) :: ![Travis-CI](https://secure.travis-ci.org/pooya/goworker.png)
