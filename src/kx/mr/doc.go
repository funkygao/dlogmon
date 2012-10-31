/*
Key of MapReduce:
    Is seperating the 'what' of distributed processing from the 'how'!
    MapReduce = select key, aggr(value) from ... group by key order by key

Conceptually in MapReduce one can think of the mapper being applied to all input 
key-value pairs and the reducer being applied to all values associated with the same key.

For programmer, he needs only think about mapper and reducer, and optionally the combiner
and the partitioner.
All other aspects of execution are handled transparently by the execution framework.

MapReduce in hadoop:
    map: (k1, v1) => [(k2, v2)] 
    reduce: (k2, [v2]) => [(k3, v3)]

    Extra roles:
        partitioner, which controls the assignment of words to reducers.
        combiner, which is a mini reducer whose input hte output of mapper





WordCount as example:
                    lines
                      |
                      | split
                      |
       -------------------------------
      |               |               |
      V               V               V
    mapper          mapper          mapper
      |               |               |
      V               V               V
    a[1] b[2]    a[2] a[5]       b[2] a[6] c[9]
      |               |               |
      V               V               V
    combiner       combiner        combiner
      |               |               |
      V               V               V
    a[1] b[2]       a[7]         b[2] a[6] c[9]

    --------------------------------------------
   | shuffle and sort: aggregate values by keys |
    --------------------------------------------

    a[1][2][6][5]   b[2][2]    c[5][9]
      |                |           |
       ----------------------------
                    |
                 reducer
                    |
            a[11] b[4] c[14]

MR in python:
    class Mapper:
        def map(docid, doc):
            for term in doc:
                emit(term, 1)

    class Reducer:
        def reduce(term, counts):
            sum = 0
            for count in counts:
                sum += count
            emit(term, sum)

*/
package mr
