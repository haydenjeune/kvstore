# Storage Engines

The `store` package provides a number of implementations of the `KvStore` interface ([store.go](store.go)) they are detailed in rough order of complexity below.

## `InMemHashMapStorage`

Uses a hash map to link keys to their stored values in memory.

### Advantages

- Fast reads and writes (independent of number of total records)

### Disadvantages

- Key and Value data must all fit into memory
- Data will be lost if the process exits

### Conclusion

A fast and simple implementation if your data doesn't need to be persisted, and is small enough to fit into memory.

## `FsAppendOnlyStorage`

This storage engine revolves around a single file to store key-value pairs. The file is append only, so updates are just added to the end of the file. This means that on read, the whole file must be scanned, and the value associated with the last occurrence of a given key is the one to return.

This file looks something like:

```
a key, the value
another key, hahaha
third, 3rd
a key, updated!
```

In this naive implementation, be sure not to set any keys of values that contain the string `", "` or a newline

### Advantages

- No memory restrictions
- Fast writes (independent of the number of total records)
- Append only write model allows multiple concurrent reading threads easily (but still only one write)
- Data is persisted even if the process exits

### Disadvantages

- Slow reads (O(n) in the number of total records)
- Unnecessary storage (superseded values never get deleted/overwritten)

### Conclusion

Solves some of the problems of `InMemHashMapStorage` in that the data doesn't need to fit into memory, and won't be lost if the process exits. However, pays for this with a big reduction in read speeds.

One way to speed up reads would be to use some kind of in memory index to map between key in memory, and value on disk.

## `InMemSortedKVStorage`

Another way of storing key-value pairs in an easily accessible way in-memory is to use some kind of search tree to store the data. The simplest implementation of this (as I have done) uses a binary search tree, where the node position is determined by the key. We build the tree as new keys are added. With the BST implementation, if we add keys in sequential order, we will encounter worst case insert and search times of O(n). If a more advanced tree structure is used, like Red/Black or AVL Trees, insert and search times can be reduced to log(n).

### Advantages

- Fastish writes and reads, possibly O(log(n)) with the number of records

### Disadvantages

- Key and Value data must all fit into memory
- Data will be lost if the process exits
