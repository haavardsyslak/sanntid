Exercise 1: Shared variable
===================================

Sharing a variable
---------------------
When running the code, we get a random number each time since the two threads tries to access the same resource at the same time, thus creating a race condition.

runtime.GOMAXPROCS() defines how many processors we should use.
When using 1 processor instead of 2, we get a result of 0 since the two threads alternates in incrementing and decrementing the variable.

Sharing a variable, but properly it workds
---------------------
