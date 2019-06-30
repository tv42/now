// Now tells what time it is in other time zones. The first argument
// identifies a time zone either by shorthand (EST, NYC) or by time zone
// file base name, such as Yellowknife or Paris.
//
//     $ now Paris
//     Sun 2019-06-30 23:03:39 CEST+02:00 Paris
//     $ now Adelaide
//     Mon 2019-07-01 06:33:50 ACST+09:30 Adelaide
//     $
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

var prog = filepath.Base(os.Args[0])

func usage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", prog)
	fmt.Fprintf(os.Stderr, "  %s [TIMEZONE]\n", prog)
	fmt.Fprintf(os.Stderr, "\n")
	flag.PrintDefaults()
}

func main() {
	flag.Usage = usage
	flag.Parse()
	if flag.NArg() > 1 {
		usage()
		os.Exit(2)
	}
	zone := ""
	t := time.Now()
	if flag.NArg() > 0 {
		zone = flag.Arg(0)
		if tz, ok := timeZone[zone]; ok {
			zone = tz
		} else if tz, ok = timeZone[toUpper(zone)]; ok {
			zone = tz
		}
		t = t.In(loadZone(zone))
	}
	const format = "Mon 2006-01-02 15:04:05 MST-07:00"
	fmt.Printf("%s %s\n", t.Format(format), zone)
}

func loadZone(zone string) *time.Location {
	loc, err := time.LoadLocation(zone)
	if err == nil {
		return loc
	}
	// Pure ASCII, but OK. Allow us to say "paris" as well as "Paris".
	if len(zone) > 0 && 'a' <= zone[0] && zone[0] <= 'z' {
		zone = string(zone[0]+'A'-'a') + string(zone[1:])
	}
	// See if there's a file with that name in /usr/share/zoneinfo
	files, _ := filepath.Glob("/usr/share/zoneinfo/*/" + zone)
	if len(files) >= 1 {
		if len(files) > 1 {
			fmt.Fprintf(os.Stderr, "now: multiple time zones; using first of %v\n", files)
		}
		loc, err = time.LoadLocation(files[0][len("/usr/share/zoneinfo/"):])
		if err == nil {
			return loc
		}
	}
	fmt.Fprintf(os.Stderr, "now: %s\n", err)
	os.Exit(1)
	return nil

}

// Pure ASCII
func toUpper(s string) string {
	var b = make([]byte, len(s))
	for i := range b {
		c := s[i]
		if 'a' <= c && c <= 'z' {
			c -= ' '
		}
		b[i] = c
	}
	return string(b)
}

// from /usr/share/zoneinfo
var timeZone = map[string]string{
	"GMT":     "Europe/London",
	"BST":     "Europe/London",
	"BSDT":    "Europe/London",
	"CET":     "Europe/Paris",
	"UTC":     "",
	"PST":     "America/Los_Angeles",
	"PDT":     "America/Los_Angeles",
	"LA":      "America/Los_Angeles",
	"LAX":     "America/Los_Angeles",
	"MST":     "America/Denver",
	"MDT":     "America/Denver",
	"CST":     "America/Chicago",
	"CDT":     "America/Chicago",
	"Chicago": "America/Chicago",
	"EST":     "America/New_York",
	"EDT":     "America/New_York",
	"NYC":     "America/New_York",
	"NY":      "America/New_York",
	"AEST":    "Australia/Sydney",
	"AEDT":    "Australia/Sydney",
	"AWST":    "Australia/Perth",
	"AWDT":    "Australia/Perth",
	"ACST":    "Australia/Adelaide",
	"ACDT":    "Australia/Adelaide",
}
