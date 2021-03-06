Locks vs Channels
======================

Which is the faster or better way to share state in Go programming?
In other words, is it a good idea to: "Do not communicate by sharing 
memory; instead, share memory by communicating." This conversation 
came up at work, so here is a simple answer to the question.

These two simple programs add by one to a shared variable a million
times. The results below were run on a four core Intel Core i7-4750HQ 
CPU @ 2.00GHz. Each function call to process performs a lock, add, 
and unlock, with locking using a `sync.Mutex`.

Results : chanadder.go
----------------------

    $ chanadder --threads 8
    Total iterations: 1000000; total time: 221.342108ms, 221ns/op
        Thread 07 ran: 125000 times
        Thread 00 ran: 125000 times
        Thread 01 ran: 125000 times
        Thread 02 ran: 125000 times
        Thread 03 ran: 125000 times
        Thread 04 ran: 125000 times
        Thread 05 ran: 125000 times
        Thread 06 ran: 125000 times

Results : lockadder.go
----------------------

    $ lockadder --threads 8
    Total iterations: 1000000; total time: 109.827864ms, 109ns/op
        Thread 06 ran: 1 times
        Thread 07 ran: 1 times
        Thread 02 ran: 1 times
        Thread 03 ran: 1 times
        Thread 04 ran: 1 times
        Thread 00 ran: 243524 times
        Thread 01 ran: 747439 times
        Thread 05 ran: 9039 times

Wait, The Locking Code isn't Fair!
----------------------------------

Looking at the results you might complain that the lock version
does not distribute the work evenly! You're right, it doesn't!
There were actually two version of the code using the lock. In
the original code each goroutine received a fair share of the
work. Intrestingly the run time was about the same ranging from
110ms to about 125ms. I changed the code to actually match
chanadder.go a bit more, which does fair work sharing "automatically"
because of the use of channels.