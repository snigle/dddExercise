# dddExercise

The goal of the program is to check validity and availability of a username on multiple platforms using a DDD implementation.\
A fake twitter implementation has been created. Don't need to fix it as there is no public API to do that.

## TODO

* Add github implementation with simple http request on `https://github.com/<username>` and check status code (200/404).
* Add validation specification of username in the domain (max length for example).
* Parallelize the calls to different backends (FakeTwitter, github, ...)
* Use CLI interface of API handler to execute the usecase.

## Bonus

Add logs for debug (errors, http call, ...)