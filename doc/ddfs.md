Here is an example of using goworker with ddfs.

Step 1: grab a large text file:
```
$ wget http://discoproject.org/media/text/chekhov.txt
```

Step 2: Split the text file into multiple parts:
```
$ split -l 2000 chekhov.txt
```

Step 3: Push the data to ddfs
```
$ ddfs push data:chekhov ./xa?
```

Step 4: Now run the job over the ddfs tag used in step 3:
```
$GOPATH/bin/jobpack -W $GOPATH/src/github.com/discoproject/goworker/examples/count_words.go -I "tag://data:chekhov"
```
