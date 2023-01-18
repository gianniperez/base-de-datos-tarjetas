package main

import (
	"bytes"
	"encoding/json"
	bolt "go.etcd.io/bbolt"
	"fmt"
	"log"
	"os"
	"strconv"
)

type Cliente struct {
	NroCliente int
	Nombre     string
	Apellido   string
	Domicilio  string
	Telefono   string
}

type Tarjeta struct {
	NroTarjeta   string
	NroCliente   int
	ValidaDesde  string
	ValidaHasta  string
	CodSeguridad string
	LimiteCompra float64
	Estado       string
}

type Comercio struct {
	NroComercio  int
	Nombre       string
	Domicilio    string
	CodigoPostal string
	Telefono     string
}

type Compra struct {
	NroOperacion int
	NroTarjeta   string
	NroComercio  int
	Fecha        string
	Monto        float64
	Pagado       bool
}

var db *bolt.DB
var err error
var opcion int = 0
var Menu string
var clientes []Cliente
var tarjetas []Tarjeta
var comercios []Comercio
var compras []Compra

func main() {
	

	Menu := 
`

% Bienvenido %
~ Opciones de acciones para realizar:
		
	[1]. Crear base NoSQL de datos: [Tarjetas]. 
	[2]. Rellenar arreglos de clientes, tarjetas, comercios, y compras. 
	[3]. Insertar datos de los clientes en DB.
	[4]. Insertar datos de las tarjetas en DB.
	[5]. Insertar datos de los comercios en DB.
	[6]. Insertar datos de las compras en DB.
	[7]. SALIR.
	
`

	for opcion != 7 {

		fmt.Print(Menu,"\n")
		fmt.Print("Ingrese el número de opción a realizar: ")
		fmt.Scanf("%d", &opcion)
		fmt.Print("Operación solicitada: ", opcion)
		fmt.Print("\n")
		if(opcion <= 0 || opcion >= 8){ 
			fmt.Print("Ingrese una opción válida. \n")
			os.Exit(0) 
		}

	switch opcion{
		case 1:
			fmt.Print("Creando base... \n")
			db, err = bolt.Open("tarjetas.db", 0600, nil)
			if err != nil {
				log.Fatal(err)
			}
			defer db.Close()
			fmt.Print("Base de datos creada. \n")
		case 2:
			fmt.Print("Insertando info en arrays... \n\n")
			rellenarArraysConDatos()
			fmt.Print("Datos ingresados en arrays! \n\n")
		case 3:
			fmt.Print("Insertando clientes... \n")
			insertarDatosClientes()
			fmt.Print("Clientes insertados! \n")
		case 4:
			fmt.Print("Insertando tarjetas..\n\n")
			insertarDatosTarjetas()
			fmt.Print("Tarjetas ingresadas! \n\n")
		case 5:
			fmt.Print("Insertando Comercios... \n")
			insertarDatosComercios()
			fmt.Print("Comercios ingresados! \n\n")
		case 6:	
			fmt.Print("Insertando compras... \n\n")
			insertarDatosCompras()
			fmt.Print("Compras ingresadas! \n\n")
		case 7:
			os.Exit(0)
		}
	}
}

func rellenarArraysConDatos() {
	
	// clientes
	clientes = append(clientes, Cliente{4, "Guadalupe", "Torres", "Av. Victorica 1128", "541134521789"})
	clientes = append(clientes, Cliente{5, "Tomas", "Gutierrez", "Formosa y Ruta 52", "541167789012"})
	clientes = append(clientes, Cliente{6, "Juan", "Ugarte", "Quevedo 3365", "541146578974"})
	
	// tarjetas
	tarjetas = append(tarjetas, Tarjeta{"1294382965767449", 4, "201912", "202212", "2364", 10500, "suspendida"})
	tarjetas = append(tarjetas, Tarjeta{"9293437257327169", 5, "202109", "202409", "937", 67000, "anulada"})
	tarjetas = append(tarjetas, Tarjeta{"5827624652643290", 6, "202211", "202511", "8246", 8000, "vigente"})
	
	// comercios
	comercios = append(comercios, Comercio{1, "Ferreteria El Cosito", "Cervantes 565","C1425EID","541165849035"})
	comercios = append(comercios, Comercio{2, "Carniceria La Vaca Feliz", "Gral. Guemes 897","C1425EID","541199999999"})
	comercios = append(comercios, Comercio{3, "Bar El Obrero", "Donado 65","E6325EID","541174392057"})
	
	// compras
	compras = append(compras, Compra{1, "1294382965767449", 1, "2022-09-11", 2500, false})
	compras = append(compras, Compra{2, "9293437257327169", 2, "2022-06-20", 50250, false})
	compras = append(compras, Compra{3, "5827624652643290", 3, "2022-08-15", 7800, false})
}

func insertarDatosClientes(){
	
	for _, cliente := range clientes{

		//Conversion a json data de cliente
		data, err := json.Marshal(cliente)
		if err != nil { // Control de error
			log.Fatal(err)
		}
		
		//Insertar en DB cada cliente.
		createUpdateBucket(db, "cliente", []byte(strconv.Itoa(cliente.NroCliente)), data) //Insertar en bucket creado o existente.
		
		//Leer el resultado insertado para mostrar.
		resultado, err := leerFilaDeBucket(db, "cliente", []byte(strconv.Itoa(cliente.NroCliente))) //Extrayendolo del bucket con modo lectura
		
		prettyJSON, err := formatJSON(resultado)
		fmt.Println(string(prettyJSON))
	}

	fmt.Print("Insercción de datos de clientes terminada. \n") // Mensaje final de operacion terminada.

}

func insertarDatosTarjetas(){
	
	for _, tarjeta := range tarjetas{

		//Conversion a json data de tarjeta
		data, err := json.Marshal(tarjeta)
		if err != nil { // Control de error
			log.Fatal(err)
		}
		
		//Insertar en DB cada tarjeta.
		createUpdateBucket(db, "tarjeta", []byte(tarjeta.NroTarjeta), data) //Insertar en bucket creado o existente.
		
		//Leer el resultado insertado para mostrar.
		resultado, err := leerFilaDeBucket(db, "tarjeta", []byte(tarjeta.NroTarjeta)) //Extrayendolo del bucket con modo lectura
		
		prettyJSON, err := formatJSON(resultado)
		fmt.Println(string(prettyJSON))
	}

	fmt.Print("Insercción de datos de tarjetas terminada. \n") // Mensaje final de operacion terminada.
}

func insertarDatosComercios(){
	
	for _, comercio := range comercios{

		//Conversion a json data de comercio
		data, err := json.Marshal(comercio)
		if err != nil { // Control de error
			log.Fatal(err)
		}
		
		//Insertar en DB cada comercio.
		createUpdateBucket(db, "comercio", []byte(strconv.Itoa(comercio.NroComercio)), data) //Insertar en bucket creado o existente.
		
		//Leer el resultado insertado para mostrar.
		resultado, err := leerFilaDeBucket(db, "comercio", []byte(strconv.Itoa(comercio.NroComercio))) //Extrayendolo del bucket con modo lectura
		
		prettyJSON, err := formatJSON(resultado)
		fmt.Println(string(prettyJSON))
	}

	fmt.Print("Insercción de datos de comercios terminada. \n") // Mensaje final de operacion terminada.
	
}

func insertarDatosCompras(){
	
	for _, compra := range compras{

		//Conversion a json data de compra
		data, err := json.Marshal(compra)
		if err != nil { // Control de error
			log.Fatal(err)
		}
		
		//Insertar en DB cada compra.
		createUpdateBucket(db, "comercio", []byte(strconv.Itoa(compra.NroComercio)), data) //Insertar en bucket creado o existente.
		
		//Leer el resultado insertado para mostrar.
		resultado, err := leerFilaDeBucket(db, "comercio", []byte(strconv.Itoa(compra.NroComercio))) //Extrayendolo del bucket con modo lectura
		
		prettyJSON, err := formatJSON(resultado)
		fmt.Println(string(prettyJSON))
	}

	fmt.Print("Insercción de datos de compras terminada. \n") // Mensaje final de operacion terminada.
}

/* Funciones aux para insertar y mostrar datos insertados */

func createUpdateBucket(db *bolt.DB, bucketName string, clave []byte, valor []byte) error {
	//Abre transaccion de escritura en la DB
	tx, err := db.Begin(true)
	if err != nil{ // Control de error de apertura de transaccion de escritura
		return err
	}
	defer tx.Rollback() // Rollback en caso de salir antes (no realiza el update)

	// Creo el bucket si no existe, y si existe lo llamo o cosumo
	b, _ := tx.CreateBucketIfNotExists([]byte(bucketName))

	// Inserto el par clave valor pasado por parametro y lo ingreso en el bucket
	err = b.Put(clave, valor)
	if err != nil { // Control de error de ingreso de datos
		return err
	}

	//Cierra transaccion. Commit los cambios.
	if err := tx.Commit(); err != nil {
		return err // Control de error en commit
	}

	return nil
}

func leerFilaDeBucket(db *bolt.DB, bucketName string, key []byte) ([]byte, error) {
	//Donde guardamos lo obtenido
	var buffer []byte

	//Abrimos una transaccion de lectura
	err := db.View(func(tx *bolt.Tx) error { // obtenemos y en tal caso devolvemos el error
		b := tx.Bucket([]byte(bucketName)) // Buscamos el bucket
		buffer = b.Get(key) // Obtenemos datos atravez del id
		return nil
	})

	return buffer, err //Retornamos buffer
}

/* Funcion sacada de: https://www.yellowduck.be/posts/pretty-print-json-with-go/
   solo para estilizar como se ve el json por consola */

func formatJSON(data []byte) ([]byte, error) {
    var out bytes.Buffer
    
    err := json.Indent(&out, data, "", "   ")
    if err == nil {
        return out.Bytes(), err
    }
    return data, nil
}
