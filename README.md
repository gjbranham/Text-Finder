# Text Finder

Text Finder is a search tool that works similarly to `grep`. It utilizes Go's concurrency features to dramatically speed up the search process.

### Clone:

`$ git clone https://github.com/gjbranham/Text-Finder.git`

### Build:

`$ make build`

This will compile and install the binary in `[starting dir]/Text-Finder/bin/`. 

We did not include a `run` Makefile target because the application is based around command-line arguments. Doing so would complicate invocation.

### Run unit tests:

`$ make test`

### Invoke:

`$ ./Text-Finder/bin/text-finder -r -d ~/Documents foo bar baz`

This will search for the strings `["foo", "bar", "baz"]` in all files recursively starting at `~/Documents`. It will search all lines of each non-binary file.

### Notes:

Planned features:
- Option to terminate search after 1st match
- Option to load search terms from file, and write results to file and/or copy matching files to folder
 
