package goopt

// An almost-drop-in replacement for flag.  It is intended to work
// basically the same way, but to parse flags like getopt does.

import (
	"os"
	"fmt"
	"time"
	"bytes"
	"path"
	"tabwriter"
	"strings"
)

var opts = make([]opt, 0, 100)

var Usage = func() string {
	if Summary != "" {
		return fmt.Sprintf("Usage of %s:\n\t",os.Args[0]) +
			Summary + "\n" + Help()
	}
	return fmt.Sprintf("Usage of %s:\n%s",os.Args[0], Help())
}
var Summary = ""
var Author = ""
var Version = ""
var Suite = ""
var Vars = make(map[string]string)

func Expand(x string) string {
	for k,v := range Vars {
		x = strings.Join(strings.Split(x, k, 0), v)
	}
	return x
}

var Help = func() string {
	h0 := new(bytes.Buffer)
	h := tabwriter.NewWriter(h0,0,8,2,' ',0)
	if (len(opts) > 1) { fmt.Fprintln(h, "Options:") }
	for _, o := range opts {
		fmt.Fprint(h,"  ")
		if len(o.shortnames) > 0 {
			for _,sn:= range o.shortnames[0:len(o.shortnames)-1] {
				fmt.Fprintf(h, "-%c, ", sn)
			}
			fmt.Fprintf(h, "-%c", o.shortnames[len(o.shortnames)-1])
			if o.allowsArg != "" { fmt.Fprintf(h, " %s", o.allowsArg) }
		}
		fmt.Fprintf(h,"\t")
		if len(o.names) > 0 {
			for _,n:= range o.names[0:len(o.names)-1] {
				fmt.Fprintf(h, "%s, ", n)
			}
			fmt.Fprint(h, o.names[len(o.names)-1])
			if o.allowsArg != "" { fmt.Fprintf(h, "=%s", o.allowsArg) }
		}
		fmt.Fprintf(h, "\t%v\n", Expand(o.help))
	}
	h.Flush()
	return h0.String()
}

var Synopsis = func() string {
	h := new(bytes.Buffer)
	for _, o := range opts {
		fmt.Fprint(h," [")
		switch {
		case len(o.shortnames) == 0:
			for _,n:= range o.names[0:len(o.names)-1] {
				fmt.Fprintf(h, "\\-\\-%s|", n[2:])
			}
			fmt.Fprintf(h, "\\-\\-%s", o.names[len(o.names)-1][2:])
			if o.allowsArg != "" { fmt.Fprintf(h, " %s", o.allowsArg) }
		case len(o.names) == 0:
			for _,c:= range o.shortnames[0:len(o.shortnames)-1] {
				fmt.Fprintf(h, "\\-%c|", c)
			}
			fmt.Fprintf(h, "\\-%c", o.shortnames[len(o.shortnames)-1])
			if o.allowsArg != "" { fmt.Fprintf(h, " %s", o.allowsArg) }
		default:
			for _,c:= range o.shortnames {
				fmt.Fprintf(h, "\\-%c|", c)
			}
			for _,n:= range o.names[0:len(o.names)-1] {
				fmt.Fprintf(h, "\\-\\-%s|", n[2:])
			}
			fmt.Fprintf(h, "\\-\\-%s", o.names[len(o.names)-1][2:])
			if o.allowsArg != "" { fmt.Fprintf(h, " %s", o.allowsArg) }
		}
		fmt.Fprint(h, "]")
	}
	return h.String()
}

var Description = func() string {
	return `To add a description to your program, define goopt.Description.

If you want paragraphs, just use two newlines in a row, like latex.`
}

type opt struct {
	names               []string
	shortnames, help    string
	needsArg bool
	allowsArg string
	process             func(string) os.Error // returns error when it's illegal
}

func addOpt(o opt) {
	newnames := make([]string,0,100)
	for _, n := range o.names {
		switch {
		case len(n) < 2:
			panic("Invalid very short flag: " + n)
		case n[0] != '-':
			panic("Invalid flag, doesn't start with '-':" + n)
		case len(n) == 2:
			o.shortnames = o.shortnames + string(n[1])
		case n[1] != '-':
			panic("Invalid long flag, doesn't start with '--':" + n)
		default:
			newnames = Append(newnames, n)
		}
	}
	o.names = newnames
	if len(opts) == cap(opts) { // reallocate
		// Allocate double what's needed, for future growth.
		newOpts := make([]opt, len(opts), len(opts)*2)
		for i, oo := range opts {
			newOpts[i] = oo
		}
		opts = newOpts
	}
	opts = opts[0 : 1+len(opts)]
	opts[len(opts)-1] = o
}

func VisitAllNames(f func (string)) {
	for _,o := range opts {
		for _,n := range o.names {
			f(n)
		}
	}
}

func NoArg(names []string, help string, process func() os.Error) {
	addOpt(opt{names, "", help, false, "", func(s string) os.Error {
			if s != "" {
				return os.NewError("unexpected flag: " + s)
			}
			return process() }})
}

func ReqArg(names []string, argname, help string, process func(string) os.Error) {
	addOpt(opt{names, "", help, true, argname, process})
}

func OptArg(names []string, def, help string, process func(string) os.Error) {
	addOpt(opt{names, "", help, false, def, func(s string) os.Error {
			if s == "" {
				return process(def)
			}
			return process(s)
	}})
}

func Alternatives(names, vs []string, help string) *string {
	out := new(string)
	*out = vs[0]
	f := func(s string) os.Error {
		for _,v := range vs {
			if s == v {
				*out = v
				return nil
			}
		}
		return os.NewError("invalid flag: "+s)
	}
	possibilities := "["+vs[0]
	for _,v := range vs[1:] {
		possibilities += "|" + v
	}
	possibilities += "]"
	ReqArg(names, possibilities, help, f)
	return out
}

func Bool(name string, d bool, help string) *bool {
	b := new(bool)
	*b = d
	f := func(s string) os.Error {
		//fmt.Println("Got", name, "of", s)
		switch s {
		case "true", "True", "yes", "":
			*b = true
		case "false", "False", "no":
			*b = false
		default:
			return os.NewError("bad boolean flag: " + s)
		}
		return nil
	}
	addOpt(opt{[]string{name}, "", help, false, fmt.Sprintf("%v",d), f})
	return b
}

func String(names []string, d string, help string) *string {
	s := new(string)
	*s = d
	f := func(ss string) os.Error {
		*s = ss
		return nil
	}
	ReqArg(names, d, help, f)
	return s
}

func Strings(names []string, d string, help string) []string {
	s := make([]string,0,100)
	f := func(ss string) os.Error {
		s = Append(s, ss)
		return nil
	}
	ReqArg(names, d, help, f)
	return s
}

func Flag(yes []string, no []string, helpyes, helpno string) *bool {
	b := new(bool)
	y := func() os.Error {
		*b = true
		return nil
	}
	n := func() os.Error {
		*b = false
		return nil
	}
	if len(yes) > 0 {
		NoArg(yes, helpyes, y)
	}
	if len(no) > 0 {
		NoArg(no, helpno, n)
	}
	return b
}

func failnoting(s string, e os.Error) {
	if e != nil {
		fmt.Println(Usage())
		fmt.Println("\n"+s, e.String())
		os.Exit(1)
	}
}

var Args []string

func Parse(extraopts func() []string) {
	// First we'll add the "--help" option.
	addOpt(opt{[]string{"--help"}, "", "show usage message", false, "",
		func(string) os.Error {
			fmt.Println(Usage())
			os.Exit(0)
			return nil
	}})
	// Let's now tally all the long option names, so we can use this to
	// find "unique" options.
	longnames := []string{"--list-options", "--create-manpage"}
	for _, o := range opts {
		longnames = Cat(longnames, o.names)
	}
	// Now let's check if --list-options was given, and if so, list all
	// possible options.
	if Any(func(a string) bool {return match(a, longnames)=="--list-options"},
		os.Args[1:]) {
		if extraopts != nil {
			for _, o := range extraopts() {
				fmt.Println(o)
			}
		}
		VisitAllNames(func (n string) { fmt.Println(n) })
		os.Exit(0)
	}
	// Now let's check if --create-manpage was given, and if so, create a
	// man page.
	if Any(func(a string) bool {return match(a, longnames)=="--create-manpage"},
		os.Args[0:]) {
		makeManpage()
		os.Exit(0)
	}
	for i:=0; i<len(os.Args);i++ {
		a := os.Args[i]
		if a == "--" {
			for _,aa := range os.Args[i:len(Args)] {
				Args = Append(Args, aa)
			}
			break
		}
		if len(a) > 1 && a[0] == '-' && a[1] != '-' {
			//fmt.Println("looking at short option",a)
			for _, o := range opts {
				//fmt.Println("checking in shortnames ", o.shortnames)
				for j, s := range a[1:len(a)] {
					for _, c := range o.shortnames {
						//fmt.Println("comparing ",string(c)," with ",string(s))
						if c == s {
							switch {
							case o.allowsArg != "" && j+1 == len(a)-1 && len(os.Args) > i+1 &&
									len(os.Args[i+1]) > 1 && os.Args[i+1][0] != '-':
								// next arg looks like a flag!
								failnoting("Error in flag -"+string(c)+":",
									         o.process(os.Args[i+1]))
								i++ // skip next arg in looking for flags...
							case o.needsArg:
								fmt.Printf("Flag -%c requires argument!\n", c)
								os.Exit(1)
							default:
								failnoting("Error in flag -"+string(c)+":",
									o.process(""))
							}
							break
						}
					}
				}
			}
		} else {
			foundone := false
		optloop: for _, o := range opts {
				for _, n := range o.names {
					if a == n {
						if o.allowsArg != "" && len(os.Args) > i+1 && len(os.Args[i+1]) > 1 && os.Args[i+1][0] != '-' {
							// next arg looks like a flag!
							failnoting("Error in flag "+n+":",
								o.process(os.Args[i+1]))
							i++ // skip next arg in looking for flags...
						} else if o.needsArg {
							fmt.Println("Flag",a,"requires argument!")
							os.Exit(1)
						} else { // no (optional) argument was provided...
							failnoting("Error in flag "+n+":", o.process(""))
						}
						foundone = true
						break optloop
					} else if o.allowsArg != "" && len(a) > len(n)+1 &&
						a[0:len(n)] == n && a[len(n)] == '=' {
						failnoting("Error in flag "+a+":", o.process(a[len(n)+1 : len(a)]))
						foundone = true
						break optloop
					}
				}
			}
			if !foundone && len(a) > 2 && a[0] == '-' && a[1] == '-' {
				failnoting("Bad flag:", os.NewError(a))
			}
			if !foundone {
				Args = Append(Args, a)
			}
		}
	}
}

func match(x string, allflags []string) string {
	for _,f := range allflags {
		if f == x {
			return x
		}
	}
	out := ""
	for _,f := range allflags {
		if len(f) >= len(x) && f[0:len(x)] == x {
			if out == "" {
				out = f
			} else {
				return ""
			}
		}
	}
	return out
}

func makeManpage() {
	_,progname :=  path.Split(os.Args[0])
	version := Version
	if Suite != "" { version = Suite + " " + version }
	fmt.Printf(".TH \"%s\" 1 \"%s\" \"%s\" \"%s\"\n", progname,
		time.LocalTime().Format("January 2, 2006"), version, Suite)
	fmt.Println(".SH NAME")
	if Summary != "" {
		fmt.Println(progname,"\\-",Summary)
	} else {
		fmt.Println(progname)
	}
	fmt.Println(".SH SYNOPSIS")
	fmt.Println(progname, Synopsis())
	fmt.Println(".SH DESCRIPTION")
	fmt.Println(formatParagraphs(Description()))
	fmt.Println(".SH OPTIONS")
	for _, o := range opts {
		fmt.Println(".TP")
		switch {
		case len(o.shortnames) == 0:
			for _,n:= range o.names[0:len(o.names)-1] {
				fmt.Printf("\\-\\-%s,", n[2:])
			}
			fmt.Printf("\\-\\-%s", o.names[len(o.names)-1][2:])
			if o.allowsArg != "" { fmt.Printf( " %s", o.allowsArg) }
		case len(o.names) == 0:
			for _,c:= range o.shortnames[0:len(o.shortnames)-1] {
				fmt.Printf( "\\-%c,", c)
			}
			fmt.Printf("\\-%c", o.shortnames[len(o.shortnames)-1])
			if o.allowsArg != "" { fmt.Printf( " %s", o.allowsArg) }
		default:
			for _,c:= range o.shortnames {
				fmt.Printf("\\-%c,", c)
			}
			for _,n:= range o.names[0:len(o.names)-1] {
				fmt.Printf("\\-\\-%s,", n[2:])
			}
			fmt.Printf( "\\-\\-%s", o.names[len(o.names)-1][2:])
			if o.allowsArg != "" { fmt.Printf( " %s", o.allowsArg) }
		}
		fmt.Printf("\n%s\n", Expand(o.help))
	}
	if Author != "" {
		fmt.Printf(".SH AUTHOR\n%s\n", Author)
	}
}

func formatParagraphs(x string) string {
	h := new(bytes.Buffer)
	lines := strings.Split(x, "\n", 0)
	for _,l := range lines {
		if l == "" {
			fmt.Fprintln(h, ".PP")
		} else {
			fmt.Fprintln(h, l)
		}
	}
	return h.String()
}
