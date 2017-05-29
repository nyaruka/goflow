import json
import time
import fabric
import subprocess
import tempfile
import requests

from fabconfig import *
from fabric.api import *
from fabric.contrib.console import confirm as _confirm
from fabric.contrib.files import *
from getpass import getpass
from os import path

try:
    import colorama
    from colorama import Fore, Back, Style
except Exception as e:
    print e
    print "Missing colorama library.  Run: % pip install colorama"
    exit(1)

try:
    import route53
except:
    print Fore.RED + "Missing Route53 library.  Run: % pip install route53"
    exit(1)

try:
    from slacker import Slacker
except:
    print Fore.RED + "Missing slack library.  Run: % pip install slacker"

if not 'AWS_ACCESS_KEY_ID' in os.environ or not 'AWS_SECRET_ACCESS_KEY' in os.environ:
    print Fore.RED + "Missing AWS_ACCESS_KEY_ID or AWS_SECRET_ACCESS_KEY in environment."
    exit(1)

colorama.init(autoreset=True)

env.HOME = os.getenv('HOME')
env.user = 'ubuntu'
env.hosts = ['localhost']

# Some examples of what these variables represent
#
# fab_config.py settings
# name: "textit"
# repo: "textit"
# settings_dir: "temba"
#
# Translate to..
#
# env.app_dir:                 /home/textit/live
# env.manage_dir:              /home/textit/live                               (holds manage.py)
# env.settings_dir:            /home/textit/live/temba                         (holds apps: settings.py, channels, msgs, etc)
# env.settings_template_dir:   /home/textit/live/env/src/textitapp/settings    (where to find settings file templates)

env.forward_agent = True
env.config = config

fabric.state.output['running'] = False
fabric.state.output['everything'] = False

# local key files
env.key_filename = [
    '%s/.ssh/staging.pem' % env.HOME,
]

from StringIO import StringIO
import sys

def status(text):
    print Fore.YELLOW + " ** " + text

def status_pending(text):
    print Fore.CYAN + "    " + text

def print_row(row, col_size):
    colors = [Fore.YELLOW, Fore.GREEN, Fore.CYAN, Fore.MAGENTA]

    line = ""
    for (index, col) in enumerate(row):
        line += colors[index] + col
        line += " " * (col_size - len(col))

    print line

@task
def choose_version():
    tags = local("git for-each-ref --sort='*authordate' --format='%(refname:short) %(contents:subject)' refs/tags", capture=True).splitlines()

    choices = []
    max_ver_len = 1
    for tag in tags[-10:]:
        version, _comment = tag.split(' ', 1)
        if _comment.startswith(" * "):
            _comment = "\n".join(_comment.split(" * ")).strip()

        choices.append((version, _comment))
        if len(version) > max_ver_len:
            max_ver_len = len(version)

    if not choices:
        abort(Fore.RED + "No versions to ship, do a % fab rev before deploying")

    print
    print_row(["Version", "Comment"], max_ver_len + 2)
    print Style.DIM + "-" * 100

    max_ver_len += 2

    colors = [Fore.YELLOW, Fore.GREEN, Fore.CYAN, Fore.MAGENTA]
    for choice in choices:
        version = choice[0]
        comments = choice[1].split("\n")

        print Fore.YELLOW + version + " " * (max_ver_len-len(version)) + Fore.GREEN + comments[0]
        if len(comments) > 1:
            for comment in comments[1:]:
                print Fore.GREEN + " " * max_ver_len + comment

    print
    version = prompt(Fore.CYAN + "Select " + Fore.YELLOW + "VERSION" + Fore.CYAN + "", default=choices[-1][0])

    if '-' in version and not _confirm(Fore.CYAN + "Are you sure want to deploy a non-master version?", default=False):
        print Fore.RED + "Cancelled"
        return

    # get the asset id for this token
    headers = {'Authorization': 'token %s' % os.getenv('GITHUB_TOKEN')}
    response = requests.get('https://api.github.com/repos/%s/%s/releases' % (env.config['gh_account'], env.config['gh_repo']), headers=headers)
    for release in response.json():
        if release['name'] == version:
            for asset in release['assets']:
                if asset['name'].find("linux_amd64") > 0:
                    env.config['asset_url'] = asset['url']
                    break

        if env.config['asset_url'] is not None:
            break

    if env.config.get('asset_url') is None:
        print Fore.RED + "No release found for %s" % version
        return

    env.config['version'] = version
    status("Found assset url: %s\n" % env.config['asset_url'])

@task
def debug():
    fabric.state.output['running'] = True
    fabric.state.output['everything'] = True
    fabric.state.output['debug'] = True

@task
def deploy():
    if os.getenv('GITHUB_TOKEN') is None:
        print Fore.RED + "Must have GITHUB_TOKEN environment variable set"
        return

    execute(choose_version)

    confirm = prompt((Fore.CYAN + "Deploy " + Fore.YELLOW + "%(version)s" + Fore.CYAN + " to " + Fore.YELLOW + "%(user)s@%(host)s" + Fore.CYAN + "?") % env.config, default="n")
    if confirm.lower() != 'y':
        print Fore.RED + "Cancelled"
        return

    env.hosts = [env.config['host']]
    env.HOME = "/home/%(user)s" % env.config

    execute(do_deploy)

@task
def chat(message, color=None):
    if not color:
        color = 'gray'

    try:
        client = Slacker('xoxp-8715334054-8715616337-12805581127-14740b8822')
        client.chat.post_message('#code', message, username='deploy', icon_url='https://feedback-assets.s3.amazonaws.com/eric/deploy')
        # print "[SLACK] %s" % message
    except:
        status(Fore.RED + "Trouble contacting slack, going on without notifications")

def do_deploy(db_file=None, do_buildout=True, quick=False):
    version = env.config['version']

    # steal the SSH_AUTH_SOCK so our ssh-agent keys are forwarded
    status("configuring ssh")
    hijack_sock(env.config['user'])

    # install the requested version into the releases directory
    live_dir = path.join(env.HOME, 'live')
    install_dir = path.join(env.HOME, 'releases', version)

    import getpass
    env.config['local_user'] = getpass.getuser();
    execute(chat,"%(local_user)s is deploying %(version)s to %(user)s@%(host)s" % env.config)
    execute(install_version, version=version, install_dir=install_dir)
    execute(stop_server)

    # link up our new version
    live_dir = path.join(env.HOME, 'live')
    run_user('rm -Rf %s' % live_dir)
    run_user('ln -s %s %s' % (install_dir, live_dir))

    # finally start the server
    execute(start_server)

    execute(chat, "Deployment successful for %(user)s@%(host)s" % env.config)
    status("server ready")
    env.config['success'] = True

@task
def install_version(version, install_dir='.', do_buildout=True):
    status("removing existing dir")
    run_user('rm -Rf %s' % install_dir)
    run_user('mkdir -p %s' % install_dir)

    with cd(install_dir):
        status("fetching release %s" % env.config['version'])
        run_user('curl -L -H "Accept: application/octet-stream" "%s?access_token=%s" -o flowserver.tar.gz' % (env.config['asset_url'], os.getenv('GITHUB_TOKEN')))
        run_user('tar zxpf flowserver.tar.gz')        

@task
def stop_server():
    """
    Start app via supervisor
    """
    status("stopping services")
    execute(chat, "Bringing server down for %(user)s@%(host)s" % env.config)

    # stop the process if it is running
    sudo('supervisorctl stop %(user)s' % env.config)

    if 'processes' in config and env.config['processes']:
        for process in env.config['processes']:
            sudo('supervisorctl stop %s_%s' % (env.config['user'], process))

    elif 'processes' in env.config:
        for process in env.config['processes']:
            sudo('supervisorctl stop %s_%s' % (env.config['user'], process))

    elif 'celery' in config and env.config['celery']:
        sudo('supervisorctl stop %(user)s_celery' % env.config)

    # at this point supervisor will have stopped all python processes, but it is possible that some
    # may still exist and keep working, so kill all python processes for this user
    with settings(warn_only=True):
        sudo('killall -9 python', user=env.config['user'])
    with settings(warn_only=True):
        sudo('killall -9 celery', user=env.config['user'])

@task
def start_server():
    """
    Start app via supervisor
    """
    execute(chat, "Bringing server up for %(user)s@%(host)s" % env.config)
    status("starting services")

    # start the process
    sudo('supervisorctl start %(user)s' % env.config)

    if 'processes' in config and env.config['processes']:
        for process in env.config['processes']:
            sudo('supervisorctl start %s_%s' % (env.config['user'], process))

    elif 'processes' in env.config:
        for process in env.config['processes']:
            sudo('supervisorctl start %s_%s' % (env.config['user'], process))

    elif 'celery' in config and env.config['celery']:
        sudo('supervisorctl start %(user)s_celery' % env.config)
 
def hijack_sock(user):
    """
    Hijacks the ssh key forwarding sock so we can fetch git
    source as a specific user instead of root. This means the
    server machines will not require our git ssh key, only
    machines which initiate the deployment do.
    """
    sock = run('echo $SSH_AUTH_SOCK')
    sock = sock.split('/')
    sudo('sudo chown -R %s:%s /tmp/%s' % (user, user, sock[2]))


def run_user(cmd):
    with prefix("source %s/env.sh" % (env.HOME)):
        # if we aren't in debug mode, hide everything
        if not fabric.state.output['debug']:
            with settings(hide('everything')):
                return sudo(cmd, user=env.config['user'])

        # otherwise, show it all
        else:
            return sudo(cmd, user=env.config['user'])
