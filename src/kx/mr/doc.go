/*
Key of MapReduce:
    Is seperating the 'what' of distributed processing from the 'how'!
    MapReduce = select key, aggr(value) from ... group by key order by key

MapReduce in hadoop:
    map: (k1, v1) => [(k2, v2)] 
    reduce: (k2, [v2]) => [(k3, v3)]

    Extra roles:
        partitioner, which controls the assignment of words to reducers.
        combiner





WordCount as example:
                    lines
       -------------------------------
      |               |               |
      V               V               V
    mapper          mapper          mapper
      |               |               |
      V               V               V
    a[1] b[2]    a[2] c[5]       b[2] a[6] c[9]
    -------------------------------------------
    shuffle and sort: aggregate values by keys
    -------------------------------------------
    a[1][2][6]   b[2][2]    c[5][9]
      |             |           |
       -------------------------
                    |
                 reducer
                    |
            a[9] b[4] c[14]

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
