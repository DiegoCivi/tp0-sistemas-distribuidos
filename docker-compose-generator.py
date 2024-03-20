import yaml
import sys

ARGS_QUANTITY = 2
CLIENT_QUANTITY_INDEX = 1
INDEX_DIFFERENCE = 1

docker_compose = {
    'version': '3.9',
    'name': 'tp0',
    'services': {
        'server': {
            'container_name': 'server',
            'image': 'server:latest',
            'entrypoint': 'python3 /main.py',
            'environment': [
                'PYTHONUNBUFFERED=1', 'LOGGING_LEVEL=DEBUG'],
            'networks': ['testing_net'],
            'volumes': ['./server/config:/config']
        },
    },
    'networks': {
        'testing_net': {
            'ipam': {
                'driver': 'default',
                'config': [{'subnet': '172.25.125.0/24'}]
            }
        }
    }
}

def main():
    if len(sys.argv) > ARGS_QUANTITY:
        print(f"Only one argument needed. Received {len(sys.argv) - INDEX_DIFFERENCE} arguments.")
        return 1
    
    try:
        client_quantity = int(sys.argv[CLIENT_QUANTITY_INDEX])
    except:
        print("Argument received could not be parsed.")
        return 1
    
    for id in range(client_quantity):
        client = f"client{id + INDEX_DIFFERENCE}"
        docker_compose['services'][client] = {
            'container_name': client,
            'image': 'client:latest',
            'entrypoint': '/client',
            'environment': [f'CLI_ID={id+1}', 'CLI_LOG_LEVEL=DEBUG'],
            'networks': ['testing_net'],
            'depends_on': ['server'],
            'volumes': ['./client/config:/config']
        }

    with open('docker-compose-dev.yaml', 'w') as file:
        yaml.dump(docker_compose, file, default_flow_style=False)

main()