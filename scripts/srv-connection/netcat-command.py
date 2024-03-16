import os

# Send message to server with netcat
os.system('echo "THIS IS A TEST" | nc server 12345') # TODO: params like port and ip should be read from config file. The msg should be received through the dockerfile
os.system('exit')