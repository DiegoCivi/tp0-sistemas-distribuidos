import os

MSG = 'THIS IS A TEST'

# The image from ./Dockerfile is built.
os.system('docker build -t test-conn ./scripts/srv-connection')

# The container is started from the previous image. It is connected to the same docker network as the server container.
# Inside the container a simple script that uses netcat to send a message to the server is runned.
srv_ans = os.popen('docker run --rm --network tp0_testing_net test-conn').read().strip('\n')

# Check if the servere is working properly by answering the same message it received.
print(f"The echo-server works: {srv_ans == MSG}")
