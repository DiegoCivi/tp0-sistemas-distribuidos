import os
import sys

msg = os.environ.get('ENV_MSG', 'ERR')

# Send message to server with netcat
os.system(f'echo "{msg}" | nc server 12345') # TODO: params like port and ip should be read from config file. The msg should be received through the dockerfile
os.system('exit')