/*
Option is the CLI options object.

Worker scan and parse a single dlog file.
Each Worker will run in it's own goroutine.

Dlog analyzer has many kinds(such as amf), which is only interested in 
some specific kind of info. 
So Worker has many sub structs, which should implement
'IsLineValid' and [map/reduct | ExtractLineInfo].

Attention:
    For performance issue, IsLineValid must be implemented in main go program, while
    map/reduce can be any runnable script file, e.g python/php/ruby/nodejs, etc.

Manager is the manager of all dlog goroutines.
There will be a single manager in runtime.

amf is a kind of Worker, which just parse 'AMF_SLOW' related log lines.

Supervisor simulates Erlang/OTP supervisors.

combiner实现本地key的聚合，对map输出的key排序，value进行迭代，类似本地的reducer

Data flow:
    manager               worker             executor        
      |                     |                   |
      | ch                  | ch                |
      \------------------->run---------------->map
                                                |
                                                V
              TransformData          MapData    |
 collectWorkers<---------transform<------------/ 
      |                                          
      V                                          
      |     ReduceData                           
    merge------->---------------------------->reduce
                                                |
                                 ReduceResult   |
   exportDB<--------------------------<--------/ 
                                               
*/
package dlog
