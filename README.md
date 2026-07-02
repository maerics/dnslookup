# dnslookup

Resolve DNS records for a host and print them as JSON or YAML.

## Usage
```
Usage: dnslookup <host> [flags]

Arguments:
  <host>    Hostname to look up.

Flags:
  -h, --help       Show context-sensitive help.
  -v, --version    Print version and exit.
  -t, --types=A,AAAA,CNAME,MX,NS,TXT,...
                   Record types to query (comma-separated).
      --yaml       Output YAML instead of JSON.
```

## Examples
```sh
$ dnslookup example.com
{
  "host": "example.com",
  "a": [
    "104.20.23.154",
    "172.66.147.243"
  ],
  "aaaa": [
    "2606:4700:10::ac42:93f3",
    "2606:4700:10::6814:179a"
  ],
  "cname": "example.com.",
  "ns": [
    "a.iana-servers.net.",
    "b.iana-servers.net."
  ]
}
$ dnslookup example.com -t ns --yaml
host: example.com
ns:
    - a.iana-servers.net.
    - b.iana-servers.net.
```
