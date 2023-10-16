# Text Finder

Text Finder is an experimental search tool utilizing some of Go's concurrency features that works similarly to `grep`.

Build:

`go build`

Invoke the tool:

`./text-finder -r -d ~/Documents some text`

In this example the tool will search for the strings `["some", "text"]` in all files recursively starting at `~/Documents`

Some notes:

- Ignores most binary files by default
- Currently does not support searching symlinks
- Unit tests are in development
- More features are planned to be added
