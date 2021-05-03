## step1:

run `python3 server.py` on server

## step2

run `curl -s -w %{http_code} -o /dev/null -d "$(cat data/v1.json|base64)" http://127.0.0.1:5000/panel` on clients

* expect HTTP 200OK responses

## step3

visit http://127.0.0.1:5000/panel/editor to edit config, submit to save config

## step4

visit dedicated url for sub info