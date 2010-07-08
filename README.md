goopt
-----

A getopt-like processor of command-line flags.  It works much like the
"flag" package, only it processes arguments in a more standard
getopt-like way.

Documentation
-------------
Once the package is installed via goinstall, use the following to view the documentation:

  # godoc --http=:6060

If you installed it from github, you will want to do this from the source directory:

  # godoc --http=:6060 --path=.

This will run in the foreground, so do it in a terminal without anything important in it.
Then you can go to http://localhost:6060/ and navigate via the package directory to the
documentation or the left-hand navigation, depending on if it was goinstalled or run from
a git clone.
