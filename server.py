from flask import Flask, request, redirect, url_for, render_template_string
import json
import base64
import logging
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger("raysub")
logger.setLevel(logging.INFO)

WRITE_TO = "/var/www/html/robots.txt"  # "/usr/share/nginx/html"
CLIENTS_FILE = "./clients.json"

app = Flask(__name__)


class ClientCFG(object):
    def __init__(self):
        self.v = "2"
        self.ps = None  # 别名
        self.add = None  # server ip/host
        self.port = 443
        self.id = None
        self.aid = None  # alterid
        self.net = None  # network
        self.type = None
        self.host = None
        self.path = None
        self.tls = None

    def __str__(self) -> str:
        return json.dumps(self.__dict__)

    def loads(self, input: str):
        cfg = json.loads(input)
        inbounds = cfg['inbounds']
        settings = inbounds[0]['settings']
        streamSettings = inbounds[0]['streamSettings']
        clients = settings['clients']
        self.id = clients[0]['id']
        self.aid = clients[0]['alterId']
        self.type = "none"
        self.tls = "tls"
        self.add = streamSettings['wsSettings']['headers']['Host']
        self.host = streamSettings['wsSettings']['headers']['Host']
        self.path = streamSettings['wsSettings']['path']
        self.net = streamSettings['network']

    def dumps(self) -> str:
        return json.dumps(self.__dict__)


def gen_sub():
    logger.info("--- started generating sub info ---")
    with open(CLIENTS_FILE, 'a+') as f:
        f.seek(0)
        raw = f.read()
        if len(raw) == 0:
            logger.warning("clients file empty, skipping sub gen")
            return
        data = json.loads(raw)
        if(type(data) is not list):
            raise TypeError
    intermediate = ["vmess://{}\n".format(base64.b64encode(
        json.dumps(c).encode('utf-8')).decode('utf-8')) for c in data]
    sub_str = base64.b64encode("".join(intermediate).encode('utf-8'))
    with open(WRITE_TO, 'w') as f:
        f.write(sub_str.decode('utf-8'))
    logger.info("sub info of {} servers written to {}".format(
        len(intermediate), WRITE_TO))


@app.route('/panel/editor')
def panel_editor():
    with open(CLIENTS_FILE, 'a+') as f:
        f.seek(0)
        data = f.read()
    return render_template_string("""
    <form action="{}" method="POST">
    <input type="submit"><br/>
    <textarea name="cfg" rows="50" cols="100">{}</textarea>
    </form>
    """.format(url_for('panel_editor'), data))


@app.route('/panel/editor', methods=['POST'])
def panel_editor_post():
    data = request.values.get('cfg')
    with open(CLIENTS_FILE, 'w+') as f:
        f.write(data)
    gen_sub()
    return redirect(url_for('panel_editor'))


@app.route('/panel')
def panel_get():
    with open(CLIENTS_FILE, 'r+') as f:
        f.seek(0)
        data = f.read()
    return data


@app.route('/panel', methods=['POST'])
def panel_post():
    data = request.get_data()
    raw_data = base64.b64decode(data).decode('utf-8')
    cfg = ClientCFG()
    cfg.loads(raw_data)
    with open(CLIENTS_FILE, 'r+') as cf:
        raw = cf.read()
        if len(raw) != 0:
            old_cfgs = json.loads(raw)
            if type(old_cfgs) is not list:
                raise TypeError
            new_cfgs = list(old_cfgs)
        else:
            new_cfgs = []
        new_cfgs.append(cfg.__dict__)
        cf.seek(0)
        cf.truncate()
        cf.write(json.dumps(new_cfgs))
        logging.info("updated cfg.")

    gen_sub()
    return ''


if __name__ == "__main__":
    gen_sub()
    app.run(debug=False)
