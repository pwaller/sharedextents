sharedextents - show proportion of physical extents shared between two files
============================================================================

Discover when two files on a CoW filesystem share identical physical data.

#### Warning: Proof-of-concept. Alpha quality software, use at your own risk.

**See Also**:

* [`fienode`](https://github.com/pwaller/fienode), for finding identical matches via a hash.
* [`filefrag -v`](https://en.wikipedia.org/wiki/E2fsprogs), for listing physical extents.

## Installation

#### From source

```
go get github.com/pwaller/sharedextents
```

#### Binary Download

See [releases](https://github.com/pwaller/sharedextents/releases/).

* [sharedextents-linux-amd64]( https://github.com/pwaller/sharedextents/releases/download/v1.0/sharedextents-linux-amd64) v1.0

## Usage: `sharedextents <a> <b>`

For example:

```
$ sharedextents a b
1142734848 / 1679605760 bytes (68.04%)

# (note: files a and b share some physical extents)
```

Exit is status `0` if extents are shared between `a` and `b`, and `1` otherwise.

## Caveats

There may be bugs. This will delete all your data and eat your cat.
When it does, that is your problem. Keep backups, folks.

#### License

MIT.
