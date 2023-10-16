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

Features we plan to add:
- Option to turn off case sensitivity
- Option to end search after 1st match
- Option to load search terms from file, and write results to file and/or copy matching files to new folder
