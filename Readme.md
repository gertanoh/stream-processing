# Working with Count min Sketch ?

What problem does it try to solve ?
How do you count the cardinality of a stream with bilions of events ? 
A classic computer science problem that can be solve using a hash map.
However, that solution has impratical space requirements at large scale.
Probabilistic data structures sacrifie some accuracy, but produce fast approximate and use less space.
There are some mathematical formula to determine the memory usage  with respect to the acceptable error rate.
Let's build some of these algorithms for stream data processing 

##  count min sketch
count min sketch supports two operations : add(x) and count(x).
The count min sketch has two parameters: the number of buckets b and the number of hash functions l.
l being really less than n, this compression leads to erros. l is used to implement independent trials, which reduce the errors. There are mathematical formulas for it. The magic is also that these factors are independent of the size of the dataset, which let's not forget is in the billions.

## Top K
Combining count min sketches and a priority queue will result in top K feature.
There is a better algorithm called [Heavy Hitter](https://www.usenix.org/system/files/conference/atc18/atc18-gong.pdf). Here is a go [package](https://github.com/segmentio/topk/tree/main) that implements it.


## references
1. [Heavy Hitter](https://www.usenix.org/system/files/conference/atc18/atc18-gong.pdf)
2. [topk](https://github.com/segmentio/topk/blob/main/topk.go)