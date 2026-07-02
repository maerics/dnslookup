// Command dnslookup resolves DNS records for a host and prints them as JSON or YAML.
package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/alecthomas/kong"
	"gopkg.in/yaml.v3"
)

type CLI struct {
	Host  string   `arg:"" help:"Hostname to look up."`
	Types []string `short:"t" default:"A,AAAA,CNAME,MX,TXT,NS" help:"Record types to query (comma-separated)."`
	YAML  bool     `help:"Output YAML instead of JSON."`
}

type MXRecord struct {
	Host string `json:"host" yaml:"host"`
	Pref uint16 `json:"pref" yaml:"pref"`
}

type Result struct {
	Host  string     `json:"host" yaml:"host"`
	A     []string   `json:"a,omitempty" yaml:"a,omitempty"`
	AAAA  []string   `json:"aaaa,omitempty" yaml:"aaaa,omitempty"`
	CNAME string     `json:"cname,omitempty" yaml:"cname,omitempty"`
	MX    []MXRecord `json:"mx,omitempty" yaml:"mx,omitempty"`
	TXT   []string   `json:"txt,omitempty" yaml:"txt,omitempty"`
	NS    []string   `json:"ns,omitempty" yaml:"ns,omitempty"`
}

func main() {
	var cli CLI
	kong.Parse(&cli)

	result := Result{Host: cli.Host}
	for _, t := range cli.Types {
		switch strings.ToUpper(strings.TrimSpace(t)) {
		case "A":
			result.A = lookupIP(cli.Host, false)
		case "AAAA":
			result.AAAA = lookupIP(cli.Host, true)
		case "CNAME":
			if cname, err := net.LookupCNAME(cli.Host); err == nil {
				result.CNAME = cname
			}
		case "MX":
			if mxs, err := net.LookupMX(cli.Host); err == nil {
				for _, mx := range mxs {
					result.MX = append(result.MX, MXRecord{Host: mx.Host, Pref: mx.Pref})
				}
			}
		case "TXT":
			if txt, err := net.LookupTXT(cli.Host); err == nil {
				result.TXT = txt
			}
		case "NS":
			if nss, err := net.LookupNS(cli.Host); err == nil {
				for _, ns := range nss {
					result.NS = append(result.NS, ns.Host)
				}
			}
		default:
			fmt.Fprintf(os.Stderr, "dnslookup: unknown record type %q\n", t)
			os.Exit(1)
		}
	}

	var out []byte
	var err error
	if cli.YAML {
		out, err = yaml.Marshal(result)
	} else {
		out, err = json.MarshalIndent(result, "", "  ")
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "dnslookup:", err)
		os.Exit(1)
	}
	fmt.Println(string(out))
}

func lookupIP(host string, v6 bool) []string {
	ips, err := net.LookupIP(host)
	if err != nil {
		return nil
	}
	var out []string
	for _, ip := range ips {
		isV4 := ip.To4() != nil
		if isV4 == !v6 {
			out = append(out, ip.String())
		}
	}
	return out
}
