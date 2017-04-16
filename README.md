# opti-transport
## Summary
The package contains functions to solve [transportation
problems](https://en.wikipedia.org/wiki/Transportation_theory_(mathematics)).
It has minimal tax method to create starting
solution and optimize function that finds negative
cells by potential method and redistributes weights
of transportation by cycle of the cell shift.

## Notes(or Problems)
- the package doesn't provide concurrency. You __CAN'T__ work
by few goroutines under one `Solving` struct(also it's ridiculous
because it can do worse or/and get data race).
- methods of the package require to save data in `Condition`
 struct after getting `Solving` struct.
- you can change the precision of floats in `NewCondition`
- now you can't get result matrix