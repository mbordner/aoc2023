* create graph of all corner nodes in outer path
  * connect them with directions
* create row and column heaps

``` 
              ┌───────────────┐
              │               │
              │               │
              └──────────┐    │
                         │    │
     ┌──────────┐        │    │
     │4  5      │        │    │
     │          │        │    │
     │          │        │    │
     └─────┐    └────────┘    │
           │                  │
           │                  │
┌──────────┘                  │
│1  2   6                     │
└──┐3   ┌────────────────┐    └──────┐
   │    │                │           │
   │    │                │           │
   │    └───────────┐    │           │
   │                │    │           │
   │                │    └───────────┘
   └────────────────┘
```

- loop through columns
- loop through rows of cur col
- if left facing column end, else:
  - if it is NW corner:
    - if corners do not cross border...  (next col over should be found)
      - find first right corner (may not exist expect (NE+NW or NE))
      - find first lower corner,  (expect SW, or NE)
        - find first right corner (may not exist, but if does, expect NE or SE)
          - from 2 right corners (or one), pick most left
          - from this col, choose or create other (above or below)

- create line segments (top -> bottom, left -> right)
- create shape from line segments
  - map lines to shape (line -> array of shapes)
  - map points to shape (point -> array of shapes)

- continue on next row past previous shape

- loop through shapes add up areas
  - loop through lines, remove any lines that point to more than one shape

