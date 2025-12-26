## Doubt

why was this not working and giving the output in a random and broken order -

```go
func filterOldEmails(emails []email) {
isOldChan := make(chan bool)

    go func() {
    	for _, e := range emails {
    		if e.date.Before(time.Date(2020, 0, 0, 0, 0, 0, 0, time.UTC)) {
    			isOldChan <- true
    			continue
    		}
    		isOldChan <- false
    	}
    }()

    go func() {
    	isOld := <-isOldChan
    	fmt.Println("email 1 is old:", isOld)
    	isOld = <-isOldChan
    	fmt.Println("email 2 is old:", isOld)
    	isOld = <-isOldChan
    	fmt.Println("email 3 is old:", isOld)
    }()

}
```

to prevent the deadlock, as channel write will be blocking.

and this is working-

```go
func filterOldEmails(emails []email) {
isOldChan := make(chan bool)

    go func() {
    	for _, e := range emails {
    		if e.date.Before(time.Date(2020, 0, 0, 0, 0, 0, 0, time.UTC)) {
    			isOldChan <- true
    			continue
    		}
    		isOldChan <- false
    	}
    }()

    isOld := <-isOldChan
    fmt.Println("email 1 is old:", isOld)
    isOld = <-isOldChan
    fmt.Println("email 2 is old:", isOld)
    isOld = <-isOldChan
    fmt.Println("email 3 is old:", isOld)

}
```

## Answer

So the simple ques is why is it working properly when i use a single goroutine for the first one only, and why it doesnt work when i put both the things/blocks in diff goroutines. i mean just why? be detailed.

This is an excellent concurrency question, and the confusion you’re feeling is 100% valid.
What’s happening here is not random, it’s unsynchronized goroutines + program exit + uncoordinated channel usage.

I’ll explain this slowly, precisely, and mechanically, because this is a core Go concept.

First: the MOST IMPORTANT RULE you’re violating

If main returns, the program exits — even if goroutines are still running.

Go does NOT wait for goroutines automatically.

This single rule explains 80% of what you’re seeing.

Let’s analyze the BROKEN version (both blocks in goroutines)

```go
func filterOldEmails(emails []email) {
isOldChan := make(chan bool)

    go func() {
    	for _, e := range emails {
    		isOldChan <- e.date.Before(...)
    	}
    }()

    go func() {
    	isOld := <-isOldChan
    	fmt.Println("email 1 is old:", isOld)
    	isOld = <-isOldChan
    	fmt.Println("email 2 is old:", isOld)
    	isOld = <-isOldChan
    	fmt.Println("email 3 is old:", isOld)
    }()

}
```

Step-by-step execution (CRITICAL)
1️⃣ filterOldEmails is called

Channel created

Two goroutines started

filterOldEmails immediately returns

⚠️ There is NO synchronization keeping the function alive.

2️⃣ What happens when filterOldEmails returns?

If this function is called from main:

```go
func main() {
filterOldEmails(emails)
}
```

Then:

```sql
main returns
↓
PROGRAM EXITS
↓
ALL goroutines are KILLED immediately
```

No cleanup. No waiting. No mercy.

3️⃣ Why output is “random / broken”

Because scheduling is non-deterministic:

Sometimes:

producer goroutine runs a bit

consumer goroutine prints one line

program exits mid-way

Sometimes:

nothing prints

sometimes partial output

sometimes reordered output

👉 This is undefined scheduling, not randomness.

Why there is NO deadlock here (important correction)

You said:

“to prevent deadlock, as channel write will be blocking”

Actually:

You did not deadlock

You had premature program termination

A deadlock would look like:

```go
fatal error: all goroutines are asleep - deadlock!
```

You didn’t get that.

Now the WORKING version (only one goroutine)

```go
func filterOldEmails(emails []email) {
isOldChan := make(chan bool)

    go func() {
    	for _, e := range emails {
    		isOldChan <- e.date.Before(...)
    	}
    }()

    isOld := <-isOldChan
    fmt.Println("email 1 is old:", isOld)
    isOld = <-isOldChan
    fmt.Println("email 2 is old:", isOld)
    isOld = <-isOldChan
    fmt.Println("email 3 is old:", isOld)

}
```

Why THIS works perfectly
Key difference:

👉 The receiver is NOT in a goroutine

That means:

filterOldEmails blocks on <-isOldChan

Function cannot return early

Program stays alive

Producer and consumer synchronize correctly

Timeline

```css
producer goroutine ──▶ sends value ──▶
main goroutine ◀── receives value
```

Channel acts as a synchronization point.

The CORE difference (this is the real answer)
Version What keeps function alive?
Both in goroutines ❌ Nothing
One in main goroutine ✅ Blocking channel receive
Channels do TWO things (important)

Most beginners think channels are just for data.

❌ WRONG.

Channels provide:

Data transfer

Synchronization

Your working version uses both.
Your broken version uses only data, no synchronization.

How to FIX the two-goroutine version properly

You must wait for the goroutines.

✅ Option 1: sync.WaitGroup (idiomatic)

```go
func filterOldEmails(emails []email) {
isOldChan := make(chan bool)
var wg sync.WaitGroup
wg.Add(2)

    go func() {
    	defer wg.Done()
    	for _, e := range emails {
    		isOldChan <- e.date.Before(...)
    	}
    	close(isOldChan)
    }()

    go func() {
    	defer wg.Done()
    	i := 1
    	for isOld := range isOldChan {
    		fmt.Printf("email %d is old: %v\n", i, isOld)
    		i++
    	}
    }()

    wg.Wait()

}
```

✅ Option 2: Do NOT use goroutine for consumer (simplest)

Which is what your working code already does 👍

One more VERY important bug (you missed this)

```go
time.Date(2020, 0, 0, ...)
```

⚠️ Months and days are 1-based in Go.

This is invalid:

```go
month = 0
day = 0
```

Correct:

```go
time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
```

Final mental model (lock this in)

Goroutines don’t keep programs alive.
Only blocking operations or explicit synchronization do.

TL;DR (tattoo-level summary)

Your first version fails because:

function returns

program exits

goroutines die

Your second version works because:

main goroutine blocks on channel

Channels synchronize ONLY when you actually block on them

To use multiple goroutines:

use WaitGroup, or

block in main goroutine
