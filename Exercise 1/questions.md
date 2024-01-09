Exercise 1 - Theory questions
-----------------------------

### Concepts

What is the difference between *concurrency* and *parallelism*?
> Concurrency means that multiple processes are making progress "at the same time", but not necessarily simultaniously. 
> Parallellism means that multiple processes are making progress "at the same time", but also simultaniously. 

What is the difference between a *race condition* and a *data race*? 
> Data race is a race condition where two concurrent processes attempt to access the same resource simultaniously.
 
*Very* roughly - what does a *scheduler* do, and how does it do it?
> The scheduler determines which thread is to be excecuted next and for how long.


### Engineering

Why would we use multiple threads? What kinds of problems do threads solve?
> By using multiple threads, we can detangle code and make certain parts of the program independent from each other, even though they run on the same CPU.

Some languages support "fibers" (sometimes called "green threads") or "coroutines"? What are they, and why would we rather use them over threads?
> Fibers divide threads into several tasks, and independent tasks do their own work as opposed to work being passed around as function. Fibers are lightweight and managed by the application, whereas threads are managed by the operating system. 

Does creating concurrent programs make the programmer's life easier? Harder? Maybe both?
> Depends on the application. In some cases we can benefit greatly from concurrency, in other cases, the increased complexity of implementing concurrent code is not neccessary.

What do you think is best - *shared variables* or *message passing*?
> Shared variable.


