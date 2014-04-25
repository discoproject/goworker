Go Worker for Disco
===================

This is an implementation of [Disco worker protocol](http://disco.readthedocs.org/en/latest/howto/worker.html) in golang.
See [discoproject.org] (http://discoproject.org) for more information.

There is a sample worker in the examples directory.  In order to run this worker, you need the jobpack utility:

```
$ go get github.com/discoproject/goworker/jobpack
$ $GOPATH/bin/jobpack -W $GOPATH/src/github.com/discoproject/goworker/examples/count_words.go -I http://discoproject.org/media/text/chekhov.txt
```

Warning: This is a work in progress and it is not ready for production use.

This implementation requires golang v1.1 or later.

Build Status: [Travis-CI](http://travis-ci.org/discoproject/goworker) :: ![Travis-CI](https://secure.travis-ci.org/discoproject/goworker.png)
