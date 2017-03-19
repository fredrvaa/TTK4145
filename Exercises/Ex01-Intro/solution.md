Reasons for concurrency & parallelism
=====================================

Concurrency vs parallelism
--------------------------
From the Go blog: http://blog.golang.org/concurrency-is-not-parallelism  
"Concurrency is the *composition* of independently executing processes, while parallelism is the simultaneous *execution* of (possibly related) computations. Concurrency is about *dealing with* lots of things at once. Parallelism is about *doing* lots of things at once."

Why multicore machines
----------------------
 - CPU frequency can't get faster
   - Power consumption is too high
     - (power consumption = a + b*freq^5, apparently: http://forums.anandtech.com/showthread.php?t=2195927)
   - Heat generated is too high, which requires even more additional power for cooling
   - The time it takes for a signal to reach one physical end of the cpu to another is multiple clock cycles
     - Only way to go faster is to go smaller, and we are fast approaching those limits for silicon
 - These CPUs are being used in "general-purpose" machines
   - OS is running multiple concurrent processes/threads anyway, why not run them on separate cores?

    
Why concurrency
---------------

 - The need to model things that happen "simultaneously"
   - See chapter below (conc. being easier/harder)
   - Dividing into threads is "natural" when we have a 1-to-1 correspondence between threads and "real-world simultaneous-ness"
 - The need to exploit multicore hardware
   - The only way to get all the performance out of a multicore CPU is by using all the cores
 - The need to utilize hardware efficiently (even a single core)
   - Example: Delays/sleeps:
     - We should do something useful in the "downtime", instead of spinning cycles
   - Example: CPUs are (generally) faster than IO (memory, disk, network)
     - (This is true also on microcontrollers. UART, SPI have clock dividers.)
     - Instead of busy-waiting for these operations to complete, we can do something useful by switching to another thread


Concurrency being "Easier"/"Harder"
-----------------------------------
This is a code quality issue. It is technically possible to write everything as a "single thread" (as long as it can run on a single core), but the resulting program will have bits of "independent" programs interleaved as a big mess.
             
Creating concurrent programs is easier if the system consists of multiple (somewhat) independent parts, like
 - Two periodic functions with different periods (how would you even do this with a single thread?)
 - Handling multiple network requests
 - Polling different independent hardware, and generating separate events for each
 - Things that are "embarrassingly parallel": http://en.wikipedia.org/wiki/Embarrassingly_parallel (but this is mostly to exploit the computing power available)
        
Creating concurrent programs is harder if the concurrent parts have to interact, as we need some way to share data safely: we need communication and/or synchronization
 - The `i=i+1` problem demonstrates that strange things may happen if we don't share data safely
 - All threads share the same processor and memory (for common CPU architectures, anyway), so this also creates a form of interaction: Running one thread means we're not running another. This is a problem for real-time (ie. timing-sensitive) applications
    
    
Different kinds of concurrent execution
---------------------------------------
See also http://stackoverflow.com/questions/3324643/processes-threads-green-threads-protothreads-fibers-coroutines-whats-the
    
The exact definitions seem to vary a little, depending on who you are asking. The important part is knowing something about:
 - Who decides what gets to run: The OS or the language runtime?
 - When we change what gets to run: When "our time is up" (preemptive scheduling), or only when we explicitly say "run something else" (non-preemptive)?
    
<!-- -->
 - Process:
   - An executing program (as opposed to just the code itself)
   - Provides the resources needed to execute a program:
     - Executable code
     - One or more threads
     - An address space
       - The OS and hardware enforces that we cannot access the memory of one process from another

<!-- -->
 - Thread:
   - A "unit of execution"
     - Has its own stack
     - Scheduled by the OS
     - Contains data structures for saving the current state of execution (the "context"), so its execution can be paused and resumed
   - Share memory within the same process
     - Some languages/extensions offer "Thread-local storage" of the heap, such that also the heap is not shared. eg Go, D
   
<!-- -->
 - Green Thread:
   - Managed & (preemptively) scheduled by the runtime (ie. not OS-managed)
     - A scheduler is part of the core library of that language
     - Often mapped N:M over several threads. See eg. GOMAXPROCS
   - May or may not share memory
    
<!-- -->
 - Co-routine:
   - Cooperatively (ie. manually, non-preemptively) scheduled by the user/programmer. No runtime, not OS-managed.
     - All co-routines must call "yield"/"reschedule" (or its equivalent), else other co-routines will not run (hence the need to be "cooperative")
     - Sometimes mapped N:M over several threads, but usually just multiplexed over a single thread.
   - Usually share memory
   - May not need to save/restore the execution context, saving significant management overhead
    
What thread spawning functions create
-------------------------------------

 - `pthread_create`:
   - Creates an OS-managed thread
    
<!-- -->
 - `threading.Thread`:
   - Creates an OS-managed thread
   - BUT: see GIL explanation below
    
<!-- -->
 - `go`:
   - Creates a co-routine
   - The calls to "yield" (called "preemption points") are inserted when compiling
     - Preemption points are added on:
       - System calls (like reading/writing terminal/files/network/etc, and all other calls to the OS)
       - Sending/receiving on channels
       - Calls to sleep-functions
       - Occasionally when calling functions, see https://golang.org/doc/go1.2#preemption
   - Effectively this looks like "green threads" from the programmers point of view, so that is arguably also a somewhat correct answer. However, the following code wouldn't hang the entire program if goroutines were truly green threads:
     - `go func(){ for{} }()`
     - The behaviour of the Go scheduler is always improving, so this answer is subject to change
     - See https://github.com/golang/go/issues/543, this issue has been "open" since 2010
        
    
    

Global Interpreter Lock
-----------------------
The python interpreter can only interpret one piece of code at a time, so no two threads can run the interpreter at once. Spawning more OS threads does not increase performance! The OS will just run whatever thread the python interpreter has made runnable.

GIL Workaround
--------------
 - Spawn more interpreters (Use the "Multiprocessing" module)
 - Share memory between the interpreters (which run in their own processes) using memory mapping
        
    
GOMAXPROCS
----------
 - Distributes/maps goroutines over more OS threads
 - Switching between OS threads is (much) slower than switching between goroutines, so increasing GOMAXPROCS does not necessarily increase performance
 - Problems with unresponsive Go programs (like the `go func(){ for{} }()` example) should not be "solved" by  increasing GOMAXPROCS
    
    
    
    
    
    
    
    
    
    
    
    
    
    
    
    
    
