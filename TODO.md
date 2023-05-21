#TODO LIST

## TODO
 - Test `World.GC()` method to optimize/clean storage after removing a good amount of components or entities
 - Consider calling GC() automatically (if configured) after more than 50% of entities inside a component storage are "dead"
 - Consider manual memory management for component storages https://dgraph.io/blog/post/manual-memory-management-golang-jemalloc/


### MMM
https://dgraph.io/blog/post/manual-memory-management-golang-jemalloc/
https://github.com/dgraph-io/ristretto/blob/750f5be31aadcf02486cb9d40471cabf4a72cbd4/z/allocator.go#L38
https://github.com/dgraph-io/ristretto/blob/master/z/calloc_jemalloc.go