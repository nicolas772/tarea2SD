# Laboratorio 2:
##Información grupo
Grupo: 14
Fecha de entrega: 5 Noviembre 2021
## Integrantes:
    - Nicolás Araya 201773550-3
    - Matías Carvajal 201873559-0
    - Fabián Jiménez 201673503-8






# Resumen
En esta entrega se ubicaran 4 componentes del sistema en 4 máquinas virtuales distintas. 


## Máquina virtual dist53: Lider + Data Node
Contiene el código que permite ejecutar la instancia de líder que se relaciona en forma de sevidor con jugador y en forma de cliente con pozo.  Contiene el código de Datanode que permite conseguir la información respecto a las jugadas de los jugadores.


## Máquina virtual dist54: Pozo + Data node
Contiene el código que permite ejecutar la instancia de pozo que se relaciona en forma de servidor con líder. Contiene el código de Datanode que actúa como base de datos para conseguir la información asociada a las jugadas.


## Máquina virtual dist55: Jugadores + Data node
Contiene el código que permite ejecutar la instancia de jugador que se relaciona en forma de cliente con líder. Actúa como base de datos para registrar información de las jugadas.


## Máquina virtual dist56: Name node
Contiene el código que permite ejecutar la instancia de Namenode que de forma síncrona distribuye los registros de jugadas entre los 3 datanode.


# Instrucciones
Para esta entrega se han realizado 4 formas de realizar ejecuciones entre máquinas virtuales que describen 4 relaciones distintas entre distintos componentes del sistema.




## 1) Relación Jugador - Líder
 En está la funcionalidad Jugador-Líder, el jugador interactúa con líder mediante una conexión gRPC en el juego Luz Roja, Luz Verde (Etapa 1), en donde el jugador solo podrá ser aceptado dentro del juego, si no ha superado la cantidad máxima de conexiones realizadas jugador-líder (16). En cada ronda el jugador ingresará un número entre el 1 y el 10, si el número es mayor al ingresado aleatoriamente por el líder (entre 6 y 10) será eliminado. Los valores ingresados se acumularán en un puntaje el cual debe tener un valor mayor de 21 antes de terminar la ronda 4, si no, también será eliminado. En caso de no superar el número del líder dado en cada ronda y de acumular un puntaje mayor o igual a 21 al terminar la ronda 4, habrá ganado la etapa.


###: Ejecución makefile:
En la máquina virtual dist53:
        cd lider/tarea2SDJugadorLider/pruebaJugadorLider/
        make run


En la máquina virtual dist:55:
        cd tarea2SDJugadorLider/pruebaJugadorLider/
        make run


###: En la interfaz:
* En la interfaz del jugador se mostrará si se desea ingresar a una partida o salir.
* Si se ingresa a la partida, se le pedirá al jugador ingresar un número entre 1 y 10,
* El servidor líder responderá automáticamente con un número entre 6 y 10. Si este es menor al ingresado por el jugador, se eliminará al jugador.
* Se acumulará el puntaje ingresado por el jugador.
* Se repite este proceso de elegir números 3 veces más. Si el jugador no supera una cantidad acumulada mayor a 21 entonces se eliminará, en caso contrario superará la etapa.




##2) Relación Líder - Pozo


Para esta prueba, el pozo será el servidor y el líder será cliente. Se establecen dos conexiones, una para la conexión sincrona a través de gRPC y protobuf y otra asincrona con rabbitmq. En la interacción se le permitirá al líder ver el monto actual y se le dará posibilidad de eliminar jugadores. También se crea el archivo eliminados.txt con el formato que se pide por enunciado.


### Ejecución makefile:
En la máquina virtual dist 54:
        cd pozo/tarea2SD
        make run


En la máquina virtual dist 53:
        cd lider/tarea2sd
        make run [se intentó hacer makefile pero no se logró hacer desde la carpeta más externa porque la forma en que fueron exportador los módulos generaban error]


* En la terminal de líder se muestra un menú para interacción, el cual dará dos opciones: ver el monto actual y eliminar a un jugador. También se crea el archivo eliminados.txt con el formato que se pide por enunciado.


##3) Relación Date node - Líder / Jugadores / Pozo 


### : Ejecución makefile:


make run




##4) Relación Name node - Data node


### : Ejecución makefile:


        make run