# Text Finder

Text Finder is an experimental search tool that works similarly to `grep`. It utilizes Go's concurrency features to speed up the search process.

Clone:

`git clone https://github.com/gjbranham/Text-Finder.git`

Build:

`cd Text-Finder && go build`

Invoke:

`./text-finder -r -d ~/Documents foo bar baz`

This will search for the strings `["foo", "bar", "baz"]` in all files recursively starting at `~/Documents`. It will search all lines of each non-binary file.

Some notes:

- Currently does not follow symlinks
- Unit tests are in development

Features we plan to add:
- Option to turn off case sensitivity
- Option to terminate search after 1st match
- Option to load search terms from file, and write results to file and/or copy matching files to folder
