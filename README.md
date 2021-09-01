# chunk

Chunk is a tiny utility for grabbing an arbitrary chunk from a file, given an
offset and length.

You'd think that `dd` would be the right tool, but unfortunately it makes this
task quite difficult.

## Installation

    go install github.com/cespare/chunk@latest

## Usage

```
chunk -h
Usage: ./chunk [OPTIONS] FILENAME OFFSET
where OPTIONS are:
  -end value
        Ending offset
  -len value
        Length of chunk
Exactly one of -end, -len must be given.
Numbers may be written as 1000, 1e3, or 1kB.
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
