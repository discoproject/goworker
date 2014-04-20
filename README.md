After cloning this repo, and setting GOPATH you can issue:

$ go build disco/jobpack

to compile a basic utility for creation and submission of jobpacks. Then use the
following command to submit the jobpack to the disco master as a job.

$ ./jobpack localhost:8989  src/disco/examples/ http://discoproject.org/media/text/chekhov.txt
