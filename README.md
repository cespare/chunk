# chunk

Chunk is a tiny utility for grabbing an arbitrary chunk from a file, given an offset and length.

You'd think that `dd` would be the right tool, but unfortunately it makes this task quite difficult.

## Installation

    go get github.com/cespare/chunk

## Usage

```
$ chunk -h
Usage: chunk [OPTIONS] FILENAME OFFSET
where OPTIONS are:
  -end=0: Ending offset
  -len=0: Length of chunk
(Exactly one of -end, -len must be given.)
```

## Example

```
$ cat f.txt
foo bar
baz
$ chunk -len 7 f.txt 3
 bar
ba
```
