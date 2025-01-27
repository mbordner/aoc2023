``` 
              ┌───────────────┐
              │               │
              │               │
              └──────────┐3==7│
                         │    │
     ┌──────────┐        │    │
     │          │        │    │
     │          │        │    │
     │          │        │    │
     └─────┐2==8└────────┘    │
           │             6    │
           │             |    │
┌──────────┘             |    │
│          5             |    │
└──┐1===┌────────────────┐4===└──────┐
   │    │                │           │
   │    │                │           │
   │===9└───────────┐    │           │
   │                │    │           │
   │                │    └───────────┘
   └────────────────┘
```

create all lines []Line
store map of points to lines.  map[Pos][]Line

for each dir N, E, S, W
    get all lines
        if !has connected next line (clockwise)
            add lines 
                # need to define how far to go..  should find next point to connect to
                # for N, we add 1, 2, 3, 4 to connect to E
                # for E, we add 5, 6
                # for S, we add 7, 8, 9