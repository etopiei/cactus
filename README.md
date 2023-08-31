# Cactus

![Cactus Logo](./small_cactus.png)

A friendly (not very strong) chess engine written in Go.

## Running

1. Install `xboard`
2. Build the app
```
$ go build
```

3. Play against the engine:
```
xboard -fUCI -fcp ./cactus
```

(For some reason after you start the app you have to go to the 'Mode' menu and select 'Machine White')