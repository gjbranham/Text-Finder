# Text Finder

Text Finder is a search tool that works similarly to `grep`. It utilizes Go's concurrency features to dramatically speed up the search process.

Clone:

`git clone https://github.com/gjbranham/Text-Finder.git`

Build:

`cd Text-Finder && make build`

This will compile and install the binary in `[starting dir]/Text-Finder/bin/`. 

We did not include a `run` Makefile target because the application is based around command-line arguments. Doing so would complicate running the program.

Invoke:

`./Text-Finder/bin/text-finder -r -d ~/Documents foo bar baz`

This will search for the strings `["foo", "bar", "baz"]` in all files recursively starting at `~/Documents`. It will search all lines of each non-binary file.

Some notes:

- Currently does not follow symlinks
- Unit tests are in development

Features we plan to add:
- Option to turn off case sensitivity
- Option to terminate search after 1st match
- Option to load search terms from file, and write results to file and/or copy matching files to folder
 