/*
Option is the CLI options object.

Dlog stands for a single dlog file executor.
Each Dlog will run in it's own goroutine.

Dlog analyzer has many kinds(such as amf), which is only interested in 
some specific kind of info. 
So Dlog has many sub structs, which should implement
'IsLineValid' and [map/reduct | ExtractLineInfo].

Attention:
    For performance issue, IsLineValid must be implemented in main go program, while
    map/reduce can be any runnable script file, e.g python/php/ruby/nodejs, etc.

Manager is the manager of all dlog goroutines.
There will be a single manager in runtime.

amf is a kind of Dlog, which just parse 'AMF_SLOW' related log lines.
*/
package dlog
