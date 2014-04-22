Go Worker for Disco
===================

This is an implementation of [Disco worker protocol](http://disco.readthedocs.org/en/latest/howto/worker.html) in golang.
See [discoproject.org] (http://discoproject.org) for more information.

To use this worker install the three libraries via `go get`

```
$ go get github.com/discoproject/goworker/jobutil
$ go get github.com/discoproject/goworker/jobpack
$ go get github.com/discoproject/goworker/worker
```

Then, use the `jobpack` command to run workers. An example:
```
$ $GOPATH/jobpack  -W examples/ -I http://discoproject.org/media/text/chekhov.txt
```

To submit a job to Disco master.

Warning: This is a work in progress and it is not ready for production use.

This implementation requires golang v1.1 or later.

Build Status: [Travis-CI](http://travis-ci.org/discoproject/goworker) :: ![Travis-CI](https://secure.travis-ci.org/discoproject/goworker.png)
