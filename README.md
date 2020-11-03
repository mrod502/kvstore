# kvstore
## a key-value store with concurrent-acces safety and optional element expiration

# testing
- some data from stocktwits was used as test data to be stored in the structure.
- 10,000,000 keys generated using xid were assigned data in the map.
- all data were accessed, and the time taken was recorded.
- the performance stats were calculated and put into `results.log`.
- a call graph was created to get an idea of the calls which contribute the most to the overall runtime.


# Under the hood

- A read-write mutex is used at the top-level to allow for concurrent reads from the map by multiple processes, optimizing for read performance.

- If enabled, the janitor periodically checks to see if an element has expired - and if this is the case, deletes it.

- Another approach to this would be to have an element-wise mutex lock and pointer structures stored in the map. This would allow concurrently editing the
    data in the struct at the cost of some extra memory overhead (a mutex per-element rather than per-map).
    This however allows the user to circumvent the thread-safety aspect like so:
    ```go
    var store = NewStore(false)
    
    type obj struct{
        MeaningOfLife uint32
    }

    //add some data
    store.Set("someKey",&obj{MeaningOfLife:34},time.Now().Unix()+120)

    var val = store.Get("someKey").(*obj)
    
    go func(){
    time.Sleep(time.Second)
    val.MeaningOfLife = 42 //unsafe editing of value not explicitly prohibited.
    }()

    store.Delete("someKey") // val will be deleted before the Field is set
    ```

# Optimization

- If there is some underlying pattern in the generated keys, a custom hash function could be used to translate the
    key to an int in some efficient manner (maps are often optimized for int keys)