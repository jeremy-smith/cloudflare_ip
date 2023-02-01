# cloudflare_ip

Small util for updating a DNS record on Cloudflare.

"DynDNS" example

    You have a small server running on your home network, but you do not have
    a static IP, so the IP might change without warning and you can no
    longer access your server from outside the LAN.

    cloudflare_ip will check your external IP using some generic
    web service which returns a json response and will update the DNS
    record on Cloudflare to reflect it.

    Run cloudflare_ip with cron to make sure your IP is updated on a regular
    schedule.

It can be used for setting any record type, not just A records, but since it
can only set the value to the IP (and only accepts IPv4) it's not very useful
for much else.

You will need a Cloudflare access token. Go into your profile
on Cloudflare and create an access token with edit rights for Zone.DNS.

Copy `config_template.yml` into `config.yml` and fill in the fields.

