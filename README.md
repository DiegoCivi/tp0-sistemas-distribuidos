# Parte 1

## Ejercicio N° 1

Se agrego en el docker compose un nuevo servicio llamado **client2**. Es identico al
servicio **client1**, pero se le cambio el _container name_, y su _CLI_ID_.

## Ejercicio N° 1.1
Para este ejercicio se creo el script _scripts/compose-gen/docker-compose-generator.py_ en python al que se le indica la cantidad de clientes y reescribe el docker compose con esa cantidad. Para correrlo hay que ejecutar `python3 scripts/compose-gen/docker-compose-generator.py [cantidad_de_clientes]`.


## Ejercicio N° 2

Para poder montar los **docker volumes** se crearon 2 nuevas carpetas (una en el lado del servidor y otra del lado del cliente) llamadas _config_. En ellas se encuentran ambos archivos de configuracion, _config.ini_ y _config.yaml_. Sobre estas carpetas se montan los **volume** y esto se hace desde el **docker