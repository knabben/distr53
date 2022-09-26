# Route53 Distributed Database

Accordingly, to [RFC1464](https://www.rfc-editor.org/rfc/rfc1464) the TXT
record of DNS format consists of the attribute name followed by the value
of the attribute. The name and value are
separated by an equals sign (=).

Following this approach this project uses Route53 as the zone file to 
store a global distributed database with high availability, your keys
are partitioned across different records with consistent hash, permitting
a more consistent access of your registers.

## Compiling

This software uses the AWS SDK, and `buraksezer/consistent` for the consistent hash
values for each member is stored on a double linked list. To compile run:

```
$ go install ./...
$ export PATH=${PATH}:${HOME}/go/bin
$ distr53
```

## Saving a new record

Make sure you have your credentials accessible with your temporary token. 
Set your HostedZone ID in the *HOSTED_ZONE* variable, you need to set the domain used
as well.

Suppose you want to save multiple key=value pairs *a=b,c=d,e=f,g=h,i=j*

```shell
$ export HOSTED_ZONE=`aws route53 list-hosted-zones | jq ".HostedZones[0].Id" -r`
$ export DOMAIN="opssec.in"

$ distr53 --keyvalue a=b,c=d,e=f,g=h,i=j --hostedzone=$HOSTED_ZONE --domain=$DOMAIN
2022/09/25 21:22:28 Adding the following values: ["a=b" "g=h"] on mbe
2022/09/25 21:22:30 Finished with status: PENDING

2022/09/25 21:22:30 Adding the following values: ["e=f"] on mba
2022/09/25 21:22:30 Finished with status: PENDING

2022/09/25 21:22:30 Adding the following values: ["c=d"] on mbb
2022/09/25 21:22:31 Finished with status: PENDING

2022/09/25 21:22:31 Adding the following values: ["i=j"] on mbc
2022/09/25 21:22:31 Finished with status: PENDING
```

## Accessing the records

To access the records is pretty simple, use dig and point to AWS nameserver.
Use the TXT record type and the member subdomain.

```shell
$ dig @ns-703.awsdns-23.net. -t txt mbe.opssec.in. +short
"a=b"
"g=h"

$ dig @ns-703.awsdns-23.net. -t txt mbb.opssec.in. +short
"c=d"
```

Note: Use this on PUBLIC hosted zones at your own risk.