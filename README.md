conf
====

Really, really simple key=val configuration file parser for Go

Example
=======

Each parameter is in the form key=val. Keys are strings (can have
embedded whitespace) and values are either strings or unsigned
integers.

Comments must be on lines by themselves and start with #

   # Moon Landing Configuration File

   program.name=Apollo
   sequence=11
   commander=N. A. Armstrong
   lm pilot=E. E. Aldrin
   cm pilot=M. Collins
   year=1969
   capcom=C. M. Duke

API
===

Read the config file with ReadConfigFile, extract the values with
GetString, GetUint and check for any extra values (often typos) with
CheckUnread.

    // ReadConfigFile reads an entire config file and return a Config
    // structure that can be used to extract single values.
    func ReadConfigFile(file string) (c *Config, err error)

    // GetString reads a string value from the parsed config file and
    // return it (or if it is missing then return the default value passed
    // in d)
	func (c *Config) GetString(k string, d string) (v string)
	
    // GetUint reads an unsigned integer from the parsed config file and
    // return it (or if it is missing then return the default value passed
    // in d)
    func (c *Config) GetUint(k string, d uint) (v uint)

    // CheckUnread checks to see if any of the parameters in the file
    // have not been read and returns a string containing those that have
    // not been read.
    func (c *Config) CheckUnread() (s string)



