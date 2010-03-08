package goopt

// An almost-drop-in replacement for flag.  It is intended to work
// basically the same way, but to parse flags like getopt does.

import (
	"os"
	"fmt"
	"bytes"
	"tabwriter"
)

var opts = make([]opt, 0, 100)

var Usage = func() string {
  return fmt.Sprintf("Usage of %s:\n%s", os.Args[0], Help())
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
		fmt.Fprintf(h, "\t%v\n", o.help)
	}
	h.Flush()
	return h0.String()
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
			addString(n, &newnames)
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

func addString(x string, xs *[]string) {
	if len(*xs) == cap(*xs) { // reallocate
		// Allocate double what's there, for future growth.
		newxs := make([]string, len(*xs), 1+len(*xs)*2)
		for i, oo := range *xs {
			newxs[i] = oo
		}
		*xs = newxs
	}
	*xs = (*xs)[0 : 1+len(*xs)]
	(*xs)[len(*xs)-1] = x
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

func String(name string, d string, help string) *string {
	s := new(string)
	*s = d
	f := func(ss string) os.Error {
		*s = ss
		return nil
	}
	ReqArg([]string{name}, d, help, f)
	return s
}

func Strings(names []string, d string, help string) []string {
	s := make([]string,0,100)
	f := func(ss string) os.Error {
		addString(ss, &s)
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
	NoArg(yes, helpyes, y)
	NoArg(no, helpno, n)
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

func Parse() {
	addOpt(opt{[]string{"--help"}, "", "show usage message", false, "",
		func(string) (e os.Error) { _,e = fmt.Println(Usage()); os.Exit(0); return }})
	for i:=0; i<len(os.Args);i++ {
		a := os.Args[i]
		if a == "--" {
			for _,aa := range os.Args[i:len(Args)] {
				addString(aa,&Args)
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
				addString(a,&Args)
			}
		}
	}
}
