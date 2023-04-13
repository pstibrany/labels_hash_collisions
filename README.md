# Hash collisions finder

This collection of programs can be used to find Prometheus series with colliding `labels.Hash()` values.
With a little change, it can be adapted to find collisions for any 64bit hash function.

This is not smart in any way... it simply generates a lot of random series (10B) and relies on probability to find collision.
With 10B randomly generated series, we have ~93% probability of finding a collision.

There are following programs included:

- `cmd/gen` generates files with random series and their hashes. 
  By default, it produces `10_000_000_000` series and hashes, which takes about 270 GB of space and about 70 minutes on my machine, when using `-c 4` option, writing files to external SSD.
- `cmd/collisions` finds hash collisions in generated files. 
  This takes about 2h on my machine, with 1000 generated files and 10_000_000 hashes per file.
  Be sure to redirect output to a file, to avoid losing logged collision due to limited buffer size of your terminal.
- `cmd/list` can be used to view content of generated files. Only useful for testing.
- `cmd/check` shows some colliding series (see below). 

Note that files only store random string, but how exactly that string is used to generate a series identifier depends on `cmd/gen/main.go`.
By default, it produces series like `metric{lbl="<random string>"}`, but it can be modified to generate any series from the random string.
See `generateRandomEntries` function in `cmd/gen/main.go` for details.

Some found collisions:

```
{__name__="metric", lbl="qeYKm3"} 15994195474147469050
{__name__="metric", lbl="2fUczT"} 15994195474147469050

{__name__="metric", lbl1="value", lbl2="l6CQ5y"} 12938137947073075402
{__name__="metric", lbl1="value", lbl2="v7uDlF"} 12938137947073075402

{__name__="metric", lbl1="W7qx", lbl2="zqqr"} 3743944359508544375
{__name__="metric", lbl1="Z00w", lbl2="wuwb"} 3743944359508544375

{__name__="metric", lbl1="59zo", lbl2="ucIY"} 11564299526604017765
{__name__="metric", lbl1="ThBT", lbl2="XYrv"} 11564299526604017765

{__name__="pqrw", lbl="Aanhoh"} 2946529260388296395
{__name__="sBmm", lbl="pwdthe"} 2946529260388296395
```

To find collisions with new `labels.Hash()` when using `stringlabels` tag, pass `-tags=stringlabels` when building `cmd/gen`. Other programs don't need this tag.

Here are some collisions found with `-tags=stringlabels`:

```
{__name__="metric", lbl="HFnEaGl"} 1999701020165299582
{__name__="metric", lbl="RqcXatm"} 1999701020165299582

{__name__="metric", lbl="gfIS7Ce"} 2592529552659112591
{__name__="metric", lbl="x5tSfjf"} 2592529552659112591
```
