# dnslookup

Resolve DNS records for one or more hosts and print them as a stream of JSON
or YAML documents.

## Usage
```
Usage: dnslookup <host> ... [flags]

Arguments:
  <host> ...    Hostname(s) to look up.

Flags:
  -h, --help       Show context-sensitive help.
  -v, --version    Print version and exit.
  -t, --types=A,AAAA,CNAME,MX,NS,TXT,...
                   Record types to query (comma-separated).
  -y, --yaml       Output YAML instead of JSON.
```

Output is one JSON document per host, or one YAML document per host
separated by `---`, in the order hosts were given.

## Examples
```sh
$ dnslookup example.com google.com -t ns
{
  "host": "example.com",
  "ns": [
    "hera.ns.cloudflare.com.",
    "elliott.ns.cloudflare.com."
  ]
}
{
  "host": "google.com",
  "ns": [
    "ns3.google.com.",
    "ns1.google.com.",
    "ns2.google.com.",
    "ns4.google.com."
  ]
}

$ dnslookup example.com google.com -t ns --yaml
---
host: example.com
ns:
  - hera.ns.cloudflare.com.
  - elliott.ns.cloudflare.com.
---
host: google.com
ns:
  - ns3.google.com.
  - ns1.google.com.
  - ns2.google.com.
  - ns4.google.com.
```
