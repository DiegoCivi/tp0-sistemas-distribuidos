import os

MSG = 'TEST'

# Generate a docker-compose with no clients
os.system('python3 docker-compose-generator.py 0')

# Build and run the server container
os.system('make docker-compose-up')

# The image from ./Dockerfile is built.
os.system(f'docker build --build-arg=MSG="{MSG}" -t test-conn ./scripts/srv-connection')

# The container is started from the previous image. It is connected to the same docker network as the server container.
# Inside the container a simple script that uses netcat to send a message to the server is runned.
srv_ans = os.popen('docker run --rm --network tp0_testing_net test-conn').read().strip('\n')

# Check if the servere is working properly by answering the same message it received.
print(f'\n[TEST] #### The echo-server works: {srv_ans == MSG} ####\n')

# Stop the servers container
os.system('make docker-compose-down')