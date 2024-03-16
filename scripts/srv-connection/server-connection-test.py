import os

# TODO: Add the case where the network is not up

os.system('docker build -t test-conn ./scripts/srv-connection')
srv_ans = os.popen('docker run --rm --network tp0_testing_net test-conn').read()
print(f"The server ans was: {srv_ans}")
