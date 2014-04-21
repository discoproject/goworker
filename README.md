Go Worker for Disco
===================

This is an implementation of [Disco worker protocol](http://disco.readthedocs.org/en/latest/howto/worker.html) in golang.
See [discoproject.org] (http://discoproject.org) for more information.

To use this worker set the GOPATH to point to the root directory of this repo and then:

```
$ go install disco/jobpack
```

creates a basic utility for creation and submission of jobpacks.  It can be used as

```
$ ./bin/jobpack localhost:8989  src/disco/examples/ http://discoproject.org/media/text/chekhov.txt
```

To submit a job to Disco master.

Warning: This is a work in progress and it is not ready for production use.

This implementation requires golang v1.1 or later.

Build Status: [Travis-CI](http://travis-ci.org/pooya/goworker) :: ![Travis-CI](https://secure.travis-ci.org/pooya/goworker.png)
