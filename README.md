# Parte 1

## Ejercicio N° 1

Se agrego en el docker compose un nuevo servicio llamado **client2**. Es identico al
servicio **client1**, pero se le cambio el _container name_, y su _CLI_ID_.

## Ejercicio N° 1.1
Para este ejercicio se creo el script _scripts/compose-gen/docker-compose-generator.py_ en python al que se le indica la cantidad de clientes y reescribe el docker compose con esa cantidad. Para correrlo hay que ejecutar `python3 scripts/compose-gen/docker-compose-generator.py [cantidad_de_clientes]`.


## Ejercicio N° 2

Para poder montar los **docker volumes** se crearon 2 nuevas carpetas (una en el lado del servidor y otra del lado del cliente) llamadas _config_. En ellas se encuentran ambos archivos de configuracion, _config.ini_ y _config.yaml_. Sobre estas carpetas se montan los **volume** y esto se hace desde el **docker

## Ejercicio N° 3

Con el _docker-compose-generator.py_ del ejercicio 1.1 se modifica el compose para que solo exista el servidor (y no levantar los clientes sin motivo ya que no los necesitamos para este test) y se ejecuta `python3 /scripts/srv-connection/server-connection-test.py`. Este script de python modifica el docker compose y buildea y corre una nueva imagen a partir de _scripts/srv-conn/Dockerfile_, que tambien la conecta a la network _tp0_testing_net_. Este nuevo container corre el script _scripts/srv-conn/netcat-command.py_ el cual usa **netcat** para mandarle un mensaje al servidor y que este le conteste.

## Ejercicio N° 4

Del lado del server, a la clase **Server** se le agrego *self._stop_server* para tener un valor booleano el cual determinara si se sigue el ciclo de aceptar conexiones o no. Este valor booleano empieza en **False**, indicando que no se tiene que detener el server. Pero ademas se agrego `__exit_gracefully(self, *args)` la cual se ejecutara apenas llege una señal **SIGTERM**. Desde esta funcion se cambia el valoor booleano a **True** para que no se continue con un nuevo ciclo en el server, despues de terminar el que se estaba ejecutando en ese momento. Ademas de eso, se le hace un _shutdown()_ al socket. Esto puede producir una **exception** la cual es catcheada. Si salta esa exception, se deja de seguir aceptando conexiones y se procede a cerrar los sockets.

Del lado del cliente, se creo un nuevo channel que recibe las señales **SIGTERM**. Y se aprovecho el _select statement_, entonces despues de cada ciclo del _for_, se fija si ese nuevo channel tiene contenido. Si lo tiene es porque en algun momento ya le mandaron **SIGTERM** y no hay que continuar. 

## Ejercicio N° 5

Se crearon los archivos _communication.py_ y _communication.go_ en donde encontramos las funciones que permiten la comunicacion entre servidor y cliente, usando nuestro propio protocolo de comunicacion. Comentaremos como es el protocolo, pero cabe aclarar que algunas cosas se modificaron en los ejercicios siguientes para satisfacer el incremento de requisitos.

El protocolo es muy simple. Se usa un header que siempre tendra un largo de 4 bytes. En este header se informa la longitud del mensaje. Entonces si tenemos el mensaje "Hola!", el header sera "0005" y el mensaje completo que se envia por el socket sera "0005Hola!". Para evitar un short-write, se envia el mensaje y se va contando cuantos bytes se escribieron. Si no se escribieron todos los bytes, se sigue enviado desde el byte que no se pudo escribir en el socket. Del lado del lector, este sabe que siempre primero tiene que leer 4 bytes, asi consigue el header y sabe cuantos bytes mas tiene que leer para conseguir el mensaje completo. No para de hacer intentos de leer el socket hasta que no se haya leido la cantidad de bytes indicada por el header.

Por el lado de la serializacion y des-serializacion, se juntan todos los parametros en una sola string separados por un "/" y una vez que se lee todo ese mensaje se hace un split para conseguir los parametros separados.

## Ejercicio N° 6

En este ejercicio, se necesito agregar un byte al header. Este byte que lo llamamos end_flag, se usa para avisarle al lector que ya se termino con el envio de datos y que puede dejar de leer en el loop y cerrar la conexion.

Ademas se cambiaron las funciones de escritura del protocolo. Esto se hizo porque en el ejercicio 5 no se tomo en cuenta el caso de por ejemplo usar letras con tilde, las cuales ocupan 2 bytes en vez de uno y generaban problemas en el protocolo, ahora se trabaja directamente con bytes y no con las strings.

Se usa ReadLine para leer linea por line eel archivo y se lo va a agregando a un batch el cual tiene limite de bets y de bytes. Una vez que se envio todo el archivo se termina la conexion. Del lado del server, este por cada batch que recibe, lo separa en bets y las guarda usando _store_bets()_.

**IMPORTANTE:** Si lo corren, es necesario que en la carpeta _./client/data_ se haga unzip a la carpeta _dataset.zip_. Ahora el docker volume es en esa carpeta y no en la carpeta _./client/config_.

## Ejercicio N° 7

Dentro de lo que es el protocolo no hay cambios. En esa parte se agregaron los requisitos que se pedian. El cliente ahora no termina la conexion una vez que envio todas las bets, si no que avisa que termino y se pone a escuchar para recibir los ganadores del sorteo. El server una vez que terminaron los 5 clientes, empieza a leer el archivo buscando los ganadores y por cada uno que encuentra le envia su documento al cliente de la agencia correspondiente.

## Ejercicio N° 8

Como los thread de python no permiten el paralelismoo, se recurrio al multiprocessing. Por cada cliente que se conecta al server, este le crea un proceso el cual se encargara de manejar a ese cliente. Para la comunicacion entre procesos se hizo uso de los Pipes y para la sincronizacion de esos se recurrio a un Semaforo binario. Cada proceso se comunica con el cliente respectivo de la misma manera que lo hacia en el ejercicio y el unico proceso que existia. A medida que llegan las bets, los procesos intentan tomar el semaforo para poder escribir todos en el mismo archivo sin generar ningun problema que el paralelismo pueda traer. Cuando el cliente termina de mandar todas las bets, el proceso que lo trataba le comunica al proceso "padre" que su cliente ya termino por un pipe. El proceso "padre", que seria el que spawneo a todos los demas procesos, estara esperando a traves del pipe a que le llegue el aviso de todos sus procesos "hijos". Una vez que suceda eso, inicia el sorteo de la loteria. Por cada ganador que encuentra en el archivo, ahora no le manda un mensaje directo al cliente, si no que le manda el documento al procesos que se encargue de ese cliente y este es el que se lo manda por el socket al cliente. De esta manera el proceso "padre" no se estaria procupando por short-writes ya que le deja esta preocupacion a los demas procesos.
