# File Search

File search search tool utilizing some of Go's concurrency features. Currently very primitive, but works similarly to `grep`.

Build:

`go build`

Invoke the tool:

`./text-finder -r -d ~/Documents some text to search for`

In this example the tool will search for the strings `["some", "text", "to", "search", "for"]` in all files recursively starting at `~/Documents`

- Ignores most binary files by default
- Currently does not support searching symlinks
