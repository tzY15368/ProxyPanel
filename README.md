# Raysub

Built with flask, this simple server aims to reduce the pain of those who run more than one v2ray instances for whatever reasons.  

Currently, the v2ray instances are set-up using theses [scripts](https://v2raytech.com/v2ray-all-in-one-script-vless-tcp-xtls-support/). Using ws+tls, it utilizes Nginx as a reverse proxy and sets up https for you. This simple server utilizes the https part.

- For each v2ray instance, the script provides a vmess link, with base64-encoded configuration json. 

- To generate a subscription file, one needs to concat all the vmess links and base64-encode them once more.

Things start getting painful when you want to modify the generated subscription files, like adding aliases to files or manually modifying addresses. Raysub does the encoding part and saving to subscription file part for you.

*Note: This is probably not the right way to do it, but it did solve my problems.*

## step 0

Set up nginx+v2ray with ws+tls with the script above.  

A configuration file should be generated for nginx in `/etc/nginx/conf.d/` with filename `[YOUR-DOMAIN].conf`

## step1

setup reverse proxy on nginx conf  

edit `/etc/nginx/conf.d/[YOUR-DOMAIN].conf`, add a proxy pass for the raysub server  

``` conf
location /panel {
    proxy_pass http://127.0.0.1:5000;
}
```

reload nginx `nginx -s reload`

modify `server.py`, edit the `WRITE_TO` const to indicate the absolute path of the subscription file, be it web-accessible or not.

run `python3 server.py` on server

## step2

run `curl -s -w %{http_code} -o /dev/null -d "$(cat /etc/v2ray/v2ray.json|base64)" http://127.0.0.1:5000/panel` on the v2ray host servers

*this bash script base64-encodes your v2ray config, and sends it to the raysub server*

- **expect HTTP 200OK responses**

## step3

visit [http://127.0.0.1:5000/panel/editor](http://127.0.0.1:5000/panel/editor) to edit config, submit to save config

## step4

visit dedicated url for sub info
