/*
MapReduce related package.

MapReduce = select key, aggr(value) from ... group by key order by key

MapReduce in hadoop:
    map: (k1, v1) => [(k2, v2)] 
    reduce: (k2, [v2]) => [(k3, v3)]
*/
package mr
