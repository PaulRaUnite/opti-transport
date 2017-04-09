# opti-transport
## Summary
The package contains functions to solve [transportation
problems](https://en.wikipedia.org/wiki/Transportation_theory_(mathematics)).
It has minimal tax method to create starting
solution and optimize function that finds cells by potential
method and redistributes weights of transportation by
cycle shift.

## Notes(or Problems)
- the package doesn't provide concurrency. You CAN'T work
by few goroutines under one Solving struct(also it's ridiculous
because it can do worse or/and get data race).
- methods of the package require to save data in Condition
 struct after getting Solving struct.
- taxes matrix MUST BE save because of reusing inside 
Condition struct(it's about NewCondition constructor),
so don't reuse it while algorithm works.