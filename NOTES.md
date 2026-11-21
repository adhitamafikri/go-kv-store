# Learning Notes

> This note contains the lesson learned from each stage of the project

## Stage 1

> In this stage, I mostly learned about writing, updating/overwriting, deleting values from the memory and preventing values from being accidentally mutated by the outside statements

- Directly assigning `slice` copies its `header`, but not the underlying bytes
  - The `header` includes: **Memory address**, **length**, **capacity**
  - Memory address: _pointer to the slice's starting element in the underlying array_

- Enforce **independent ownership** on the value for preventing the accidental mutation on the original data
  - Leverage **copy** data before storing the data to the storage AND before returning the retrieved data from DB
    - We can simply use `bytes.Clone(value)`
    - Or we can use these set of operations:
      - Initialize a **zero-filled slice** with `cloned := make([]byte, len(value))`
      - Fill the new zero-filled slice with `copy(cloned, value)`. This copies the `value` bytes into the `cloned` slice. Then you can store that `cloned` slice into the database
      - Example:

        ```go
        cloned := make([]byte, len(value)) // allocate independent new slice for the backing storage
        copy(cloned, value) // Copy the bytes from that value into that backing storage
        db.data[string(key)] = cloned // Store the new slice into the map
        ```

  - Copying does not make bytes immutable. Cloning on both `Put` and `Get` prevents the callers from mutating the stored bytes through the input or returned slices. Callers must use `Put` to replace the stored value

...

## Stage 10

TBD

## Closing

TBD
