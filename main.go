// Command dnslookup resolves DNS records for a host and prints them as JSON or YAML.
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"

	"github.com/alecthomas/kong"
	"gopkg.in/yaml.v3"
)

type CLI struct {
	Version kong.VersionFlag `help:"Print version and exit." short:"v"`
	Hosts   []string         `arg:"" name:"host" help:"Hostname(s) to look up."`
	Types   []string         `short:"t" default:"A,AAAA,CNAME,MX,NS,TXT" help:"Record types to query (comma-separated)."`
	YAML    bool             `help:"Output YAML instead of JSON." short:"y"`
}

// Linked at build time.
var version, commit, date string

func getVersionString() string {
	if version == "" {
		return "(unknown)"
	}
	return fmt.Sprintf("v%v, commit=%v, timestamp=%v", version, commit, date)
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
	kong.Parse(&cli, kong.UsageOnError(), kong.Vars{"version": getVersionString()})

	var buf bytes.Buffer
	for _, host := range cli.Hosts {
		result, err := lookupHost(host, cli.Types)
		if err != nil {
			fmt.Fprintln(os.Stderr, "dnslookup:", err)
			os.Exit(1)
		}

		var out []byte
		if cli.YAML {
			out, err = yaml.Marshal(result)
			buf.WriteString("---\n")
		} else {
			out, err = json.MarshalIndent(result, "", "  ")
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, "dnslookup:", err)
			os.Exit(1)
		}
		buf.Write(out)
		buf.WriteString("\n")
	}
	printResult(buf.Bytes(), cli.YAML)
}

// lookupHost resolves the requested record types for host, returning an
// error if an unknown record type is requested.
func lookupHost(host string, types []string) (Result, error) {
	result := Result{Host: host}
	for _, t := range types {
		switch strings.ToUpper(strings.TrimSpace(t)) {
		case "A":
			result.A = lookupIP(host, false)
		case "AAAA":
			result.AAAA = lookupIP(host, true)
		case "CNAME":
			if cname, err := net.LookupCNAME(host); err == nil {
				result.CNAME = cname
			}
		case "MX":
			if mxs, err := net.LookupMX(host); err == nil {
				for _, mx := range mxs {
					result.MX = append(result.MX, MXRecord{Host: mx.Host, Pref: mx.Pref})
				}
			}
		case "TXT":
			if txt, err := net.LookupTXT(host); err == nil {
				result.TXT = txt
			}
		case "NS":
			if nss, err := net.LookupNS(host); err == nil {
				for _, ns := range nss {
					result.NS = append(result.NS, ns.Host)
				}
			}
		default:
			return Result{}, fmt.Errorf("unknown record type %q", t)
		}
	}
	return result, nil
}

func lookupIP(host string, v6 bool) []string {
	ips, err := net.LookupIP(host)
	if err != nil {
		return nil
	}
	return filterIPs(ips, v6)
}

// filterIPs returns the string form of the IPv6 addresses in ips if v6 is
// true, or the IPv4 addresses otherwise.
func filterIPs(ips []net.IP, v6 bool) []string {
	var out []string
	for _, ip := range ips {
		isV4 := ip.To4() != nil
		if isV4 == !v6 {
			out = append(out, ip.String())
		}
	}
	return out
}

// printResult prints data (a stream of JSON documents, or YAML documents
// separated by "---") by piping it through jq (JSON) or yq (YAML), falling
// back to a plain print if the tool is unavailable or fails.
func printResult(data []byte, useYAML bool) {
	name, args := "jq", []string{"."}
	if useYAML {
		name, args = "yq", []string{"e", "."}
	}
	if path, err := exec.LookPath(name); err == nil {
		cmd := exec.Command(path, args...)
		cmd.Stdin = bytes.NewReader(data)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if cmd.Run() == nil {
			return
		}
	}
	fmt.Println(string(data))
}
