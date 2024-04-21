# Count Min Sketch

## Problem Statement
Counting the cardinality of a stream with billions of events poses a challenge in terms of memory usage and processing speed. Traditional solutions like hash maps become impractical at large scales due to their space requirements.

## Solution Overview
The Count Min Sketch offers a probabilistic data structure solution. It sacrifices some accuracy for fast approximation and reduced space usage. By employing mathematical formulas to determine memory usage based on acceptable error rates, it efficiently processes stream data.

## Key Features
Operations: Supports two primary operations: add(x) and count(x).
Parameters: Defined by the number of buckets (b) and the number of hash functions (l). The compression achieved by having l significantly less than n (the size of the dataset) introduces errors, mitigated by implementing independent trials. These parameters remain independent of dataset size, crucial for processing large datasets.


## Additional Functionality
### Top K
Combining Count Min Sketches with a priority queue enables the extraction of top K features. For more robust functionality, consider the Heavy Hitter algorithm, which offers improved performance.

## Resources
1. [Heavy Hitter](https://www.usenix.org/system/files/conference/atc18/atc18-gong.pdf)
2. [topk](https://github.com/segmentio/topk/blob/main/topk.go)