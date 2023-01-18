package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"fmt"
	"os"
	"log"
	"time"
)

type consumo struct {
	nroconsumo   int
	nrotarjeta   string
	codseguridad string
	nrocomercio  int
	monto        float64
}

const dbInfo string = "user=postgres password=123 host=localhost dbname=tarjetas sslmode=disable"

var db *sql.DB
var err error
var opcion int = 0
var Menu string
var nroConsumo int = 1
var TotalDeConsumos int = 0
var cantidadAlertas int = 0

func main() {

Menu := 
`% Bienvenido % 
~ Opciones de acciones para realizar:
		
	[1]. Crear Base de datos de Tarjetas. 
	[2]. Crear Tablas. 
	[3]. Crear PK's y FK's. 
	[4]. Borrar PK's y FK's. 
	[5]. Insertar datos en tablas. 
	[6]. Crear y cargar funciones. 
	[7]. Autorizar compras. 
	[8]. Generar Resumenes. 
	[9]. SALIR.`
	
	for opcion != 9 {
		
		fmt.Print(Menu,"\n")
		fmt.Print("Ingrese el número de opción a realizar: ")
		fmt.Scanf("%d", &opcion)
		fmt.Print("Operación solicitada: ", opcion)
		fmt.Print("\n")
		if(opcion <= 0 || opcion >= 10){ 
			fmt.Print("Ingrese una opción válida. \n")
			os.Exit(0) 
		}
		
	switch opcion{
		case 1:
			fmt.Print("Creando Base de datos... \n")
			createDatabase()
			fmt.Print("Base de datos creada! \n")
		case 2:
			fmt.Print("Creando Tablas... \n")
			createTables()
			fmt.Print("Tablas de base de datos creadas! \n")
		case 3:
			fmt.Print("Creando PK's y FK's... \n")
			createPksAndFks()
			fmt.Print("PK's y FK's de tablas creadas! \n")
		case 4:
			fmt.Print("Borrando PK's y FK's... \n")
			dropPksAndFks()
			fmt.Print("PK's y FK's eliminadas! \n")
		case 5:
			fmt.Print("Insertando datos en tablas... \n")
			insertValues()
			fmt.Print("Datos insertados en las tablas! \n")
		case 6:
			fmt.Print("Creando y cargando funciones... \n")
			createFunctionAutorizaciones()
			createFunctionResumenes()
			createFunctionAlertas()
			go informarAlertaNueva()
			fmt.Print("Funciones Creadas! \n")
		case 7:
			fmt.Print("Autorizando compras... \n")
			autorizarCompras()
			fmt.Print("Compras autorizadas! \n")
		case 8:
			fmt.Print("Generando resumenes... \n")
			generarResumenes()
			fmt.Print("Resumenes generados! \n")
		case 9:
			os.Exit(0)
		}
	}
}

func createDatabase() {
	db, err = sql.Open("postgres", "user=postgres password=123 host=localhost dbname=postgres sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec("DROP DATABASE IF EXISTS tarjetas")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("CREATE DATABASE tarjetas")
	if err != nil {
		log.Fatal(err)
	}
	
	db, err = sql.Open("postgres", "user=postgres password=123 host=localhost dbname=tarjetas sslmode=disable")
}

func createTables() {
	_, err = db.Exec(`CREATE TABLE cliente (nrocliente SERIAL, nombre TEXT, apellido TEXT, domicilio TEXT, telefono CHAR(12));
	CREATE TABLE tarjeta (nrotarjeta CHAR(16), nrocliente INT, validadesde CHAR(6), validahasta CHAR(6), codseguridad CHAR(4), limitecompra DECIMAL(8,2), estado CHAR(10));
	CREATE TABLE comercio (nrocomercio SERIAL, nombre TEXT, domicilio TEXT, codigopostal CHAR(8), telefono CHAR(12));
	CREATE TABLE compra (nrooperacion SERIAL, nrotarjeta CHAR(16), nrocomercio INT, fecha TIMESTAMP, monto DECIMAL(7,2), pagado BOOLEAN);
	CREATE TABLE rechazo (nrorechazo SERIAL, nrotarjeta CHAR(16), nrocomercio INT, fecha TIMESTAMP, monto DECIMAL(7,2), motivo TEXT);
	CREATE TABLE cierre (año INT, mes INT, terminacion INT, fechainicio DATE, fechacierre DATE, fechavto DATE);
	CREATE TABLE cabecera (nroresumen SERIAL, nombre TEXT, apellido TEXT, domicilio TEXT, nrotarjeta CHAR(16), desde DATE, hasta DATE, vence DATE, total DECIMAL(8,2));
	CREATE TABLE detalle (nroresumen INT, nrolinea INT, fecha DATE, nombrecomercio TEXT, monto DECIMAL(7,2));
	CREATE TABLE alerta (nroalerta SERIAL, nrotarjeta CHAR(16), fecha TIMESTAMP, nrorechazo INT, codalerta INT, descripcion TEXT);
	CREATE TABLE consumo (nroconsumo SERIAL,nrotarjeta CHAR(16), codseguridad CHAR(4), nrocomercio INT, monto DECIMAL(7,2));`)
	if err != nil {
		log.Fatal(err)
	}
}

func createPksAndFks() {
	_, err = db.Exec(`ALTER TABLE cliente ADD CONSTRAINT cliente_pk PRIMARY KEY (nrocliente);
	ALTER TABLE tarjeta ADD CONSTRAINT tarjeta_pk PRIMARY KEY (nrotarjeta);
	ALTER TABLE comercio ADD CONSTRAINT comercio_pk PRIMARY KEY (nrocomercio);
	ALTER TABLE compra ADD CONSTRAINT compra_pk PRIMARY KEY (nrooperacion);
	ALTER TABLE rechazo ADD CONSTRAINT rechazo_pk PRIMARY KEY (nrorechazo);
	ALTER TABLE cierre ADD CONSTRAINT cierre_pk PRIMARY KEY (año, mes, terminacion);
	ALTER TABLE cabecera ADD CONSTRAINT cabecera_pk PRIMARY KEY (nroresumen);
	ALTER TABLE detalle ADD CONSTRAINT detalle_pk PRIMARY KEY (nroresumen, nrolinea);
	ALTER TABLE alerta ADD CONSTRAINT alerta_pk PRIMARY KEY (nroalerta);
	ALTER TABLE tarjeta ADD CONSTRAINT tarjeta_nrocliente_fk FOREIGN KEY (nrocliente) REFERENCES cliente (nrocliente);
	ALTER TABLE compra ADD CONSTRAINT compra_nrotarjeta_fk FOREIGN KEY (nrotarjeta) REFERENCES tarjeta (nrotarjeta);
	ALTER TABLE compra ADD CONSTRAINT compra_nrocomercio_fk FOREIGN KEY (nrocomercio) REFERENCES comercio (nrocomercio);
	ALTER TABLE rechazo ADD CONSTRAINT rechazo_nrocomercio_fk FOREIGN KEY (nrocomercio) REFERENCES comercio (nrocomercio);
	ALTER TABLE cabecera ADD CONSTRAINT cabecera_nrotarjeta_fk FOREIGN KEY (nrotarjeta) REFERENCES tarjeta (nrotarjeta);
	ALTER TABLE alerta ADD CONSTRAINT alerta_nrorechazo_fk FOREIGN KEY (nrorechazo) REFERENCES rechazo (nrorechazo);`)
	if err != nil {
		log.Fatal(err)
	}
}

func dropPksAndFks() {
	_, err = db.Exec(`
	ALTER TABLE tarjeta DROP CONSTRAINT tarjeta_nrocliente_fk;
	ALTER TABLE compra DROP CONSTRAINT compra_nrotarjeta_fk;
	ALTER TABLE compra DROP CONSTRAINT compra_nrocomercio_fk;
	ALTER TABLE rechazo DROP CONSTRAINT rechazo_nrocomercio_fk;
	ALTER TABLE cabecera DROP CONSTRAINT cabecera_nrotarjeta_fk;
	ALTER TABLE alerta DROP CONSTRAINT alerta_nrorechazo_fk;
	ALTER TABLE cliente DROP CONSTRAINT cliente_pk;
	ALTER TABLE tarjeta DROP CONSTRAINT tarjeta_pk;
	ALTER TABLE comercio DROP CONSTRAINT comercio_pk;
	ALTER TABLE compra DROP CONSTRAINT compra_pk;
	ALTER TABLE rechazo DROP CONSTRAINT rechazo_pk;
	ALTER TABLE cierre DROP CONSTRAINT cierre_pk;
	ALTER TABLE cabecera DROP CONSTRAINT cabecera_pk;
	ALTER TABLE detalle DROP CONSTRAINT detalle_pk;
	ALTER TABLE alerta DROP CONSTRAINT alerta_pk;
	`)
	if err != nil {
		log.Fatal(err)
	}
}

func insertValues() {
	_, err = db.Exec(`INSERT INTO cliente (nombre, apellido, domicilio, telefono) VALUES ('Matias', 'Avila', '9 de Julio 2302', '541112321232');
	INSERT INTO cliente (nombre, apellido, domicilio, telefono) VALUES ('Giannina', 'Perez', 'Panamericana Km 36.5', '541145678970');
	INSERT INTO cliente (nombre, apellido, domicilio, telefono) VALUES ('Lucas', 'Prieto', 'Av. Pres. Arturo Umberto Illia 3770', '541142335678');
	INSERT INTO cliente (nombre, apellido, domicilio, telefono) VALUES ('Guadalupe', 'Torres', 'Av. Victorica 1128', '541134521789');
	INSERT INTO cliente (nombre, apellido, domicilio, telefono) VALUES ('Tomas', 'Gutierrez', 'Formosa y Ruta 52', '541167789012');
	INSERT INTO cliente (nombre, apellido, domicilio, telefono) VALUES ('Juan', 'Ugarte', 'Quevedo 3365', '541146578974');
	INSERT INTO cliente (nombre, apellido, domicilio, telefono) VALUES ('Sergio', 'Messi', 'Av. Del Libertador 6820', '541156920932');
	INSERT INTO cliente (nombre, apellido, domicilio, telefono) VALUES ('Santiago', 'Pereyra', 'Paraná 3745', '541154648972');
	INSERT INTO cliente (nombre, apellido, domicilio, telefono) VALUES ('Karina', 'Castan', 'San Martín 546', '541176853412');
	INSERT INTO cliente (nombre, apellido, domicilio, telefono) VALUES ('Emiliano', 'Ayala', 'Sarmiento 2157', '541156748921');
	INSERT INTO cliente (nombre, apellido, domicilio, telefono) VALUES ('Fabian', 'Moreno', 'Juan Julian Lastra 2400', '541145769276');
	INSERT INTO cliente (nombre, apellido, domicilio, telefono) VALUES ('Jorge', 'Peron', 'Av Nestor C. Kirchner 1142', '541125678970');
	INSERT INTO cliente (nombre, apellido, domicilio, telefono) VALUES ('Federico', 'Santillan', 'Josiah Williams 209', '541178929283');
	INSERT INTO cliente (nombre, apellido, domicilio, telefono) VALUES ('Sebastian', 'Rodriguez', 'Perito Moreno 1460', '541176526341');
	INSERT INTO cliente (nombre, apellido, domicilio, telefono) VALUES ('Camila', 'Martin', 'San Martín 800', '541154327801');
	INSERT INTO cliente (nombre, apellido, domicilio, telefono) VALUES ('Tania', 'Rojo', 'Av. España 309', '541134679086');
	INSERT INTO cliente (nombre, apellido, domicilio, telefono) VALUES ('Belen', 'Aristimuño', 'Av. Libertad 254', '541197367361');
	INSERT INTO cliente (nombre, apellido, domicilio, telefono) VALUES ('Morena', 'Dominguez', 'Av. Sáenz Peña 318', '541121340968');
	INSERT INTO cliente (nombre, apellido, domicilio, telefono) VALUES ('Paola', 'Tugas', 'Enrique Bodereau 7571', '541197656291');
	INSERT INTO cliente (nombre, apellido, domicilio, telefono) VALUES ('Paulo', 'Gomez', 'Recta Martinoli 8357', '541187561264');
	
	INSERT INTO comercio (nombre, domicilio, codigopostal, telefono) VALUES ('Ferreteria El Cosito', 'Cervantes 565','C1425EID','541165849035');
	INSERT INTO comercio (nombre, domicilio, codigopostal, telefono) VALUES ('Carniceria La Vaca Feliz', 'Gral. Guemes 897','C1425EID','541199999999');
	INSERT INTO comercio (nombre, domicilio, codigopostal, telefono) VALUES ('Bar El Obrero', 'Donado 65','E6325EID','541174392057');
	INSERT INTO comercio (nombre, domicilio, codigopostal, telefono) VALUES ('Merceria Doña Tita', 'Pres. Juan Domingo Perón 4739','F3325EID','541129495473');
	INSERT INTO comercio (nombre, domicilio, codigopostal, telefono) VALUES ('Farmacia Balcarce', 'Paraná 3822','F6721EID','541139572659');
	INSERT INTO comercio (nombre, domicilio, codigopostal, telefono) VALUES ('Supermecado Los 3 Hermanos', 'Antonio Saenz 2041','B6721EID','541148573927');
	INSERT INTO comercio (nombre, domicilio, codigopostal, telefono) VALUES ('Hotel San Carlos', 'Av. Ing Agustín Rocca 249','G8902EID','541160382046');
	INSERT INTO comercio (nombre, domicilio, codigopostal, telefono) VALUES ('Dietetica El Obeso', 'Conquista Del Desierto 230','H9012EID','541140682540');
	INSERT INTO comercio (nombre, domicilio, codigopostal, telefono) VALUES ('Kiosco La Esquina', 'Dr. Salvador Sallares 63','F5462EID','541168392028');
	INSERT INTO comercio (nombre, domicilio, codigopostal, telefono) VALUES ('Verduleria Don Pedro', 'Bernardo De Monteagudo 3351','A7810EID','541147592038');
	INSERT INTO comercio (nombre, domicilio, codigopostal, telefono) VALUES ('Panaderia La Ponderosa', 'Brig. Gral. Juan M. De Rosas 14446','G2091EID','541100374993');
	INSERT INTO comercio (nombre, domicilio, codigopostal, telefono) VALUES ('Veterinaria Cura Bicho', 'Av. Eva Duarte de Peron 1474','H5674EID','541100000384');
	INSERT INTO comercio (nombre, domicilio, codigopostal, telefono) VALUES ('Tienda de Ropa La Gastadera', 'Hipólito Yrigoyen 1807','C4512EID','541103058239');
	INSERT INTO comercio (nombre, domicilio, codigopostal, telefono) VALUES ('Vidrieria El Reflejo', 'Av Luro 5975','F3241EID','541100000037');
	INSERT INTO comercio (nombre, domicilio, codigopostal, telefono) VALUES ('Gomeria La Rueda Feliz', '9 De Julio 1185','E4312EID','541100000070');
	INSERT INTO comercio (nombre, domicilio, codigopostal, telefono) VALUES ('Mecanica Austral', 'Calle 8 788','G3212EID','541100000800');
	INSERT INTO comercio (nombre, domicilio, codigopostal, telefono) VALUES ('Lavadero Blanca', 'Calle 12 esq. 58 Nº 1249','C6786EID','541100000500');
	INSERT INTO comercio (nombre, domicilio, codigopostal, telefono) VALUES ('Colchones Dulce Sueños', 'Sarmiento 2685','C2341EID','541100000900');
	INSERT INTO comercio (nombre, domicilio, codigopostal, telefono) VALUES ('Salon de Fiestas Wonka', 'Catamarca 1865','F0443EID','541100001000');
	INSERT INTO comercio (nombre, domicilio, codigopostal, telefono) VALUES ('Restaurant Ratatouille', 'Bartolomé Mitre 2642','G3245EID','541100009843');

	INSERT INTO tarjeta (nrotarjeta, nrocliente, validadesde, validahasta, codseguridad, limitecompra, estado) VALUES ('9320584378172604', 1, '201106', '201406', '906', 20000, 'anulada');
	INSERT INTO tarjeta (nrotarjeta, nrocliente, validadesde, validahasta, codseguridad, limitecompra, estado) VALUES ('4982394782623736', 2, '202008', '202308', '2958', 25000, 'vigente');
	INSERT INTO tarjeta (nrotarjeta, nrocliente, validadesde, validahasta, codseguridad, limitecompra, estado) VALUES ('5129043812284623', 3, '202203', '202503', '753', 30800, 'vigente');
	INSERT INTO tarjeta (nrotarjeta, nrocliente, validadesde, validahasta, codseguridad, limitecompra, estado) VALUES ('1294382965767449', 4, '201912', '202212', '2364', 10500, 'suspendida');
	INSERT INTO tarjeta (nrotarjeta, nrocliente, validadesde, validahasta, codseguridad, limitecompra, estado) VALUES ('9293437257327169', 5, '202109', '202409', '937', 67000, 'anulada');
	INSERT INTO tarjeta (nrotarjeta, nrocliente, validadesde, validahasta, codseguridad, limitecompra, estado) VALUES ('5827624652643290', 6, '202211', '202511', '8246', 8000, 'vigente');
	INSERT INTO tarjeta (nrotarjeta, nrocliente, validadesde, validahasta, codseguridad, limitecompra, estado) VALUES ('3928217404085943', 7, '202212', '202512', '1485', 90000, 'vigente');
	INSERT INTO tarjeta (nrotarjeta, nrocliente, validadesde, validahasta, codseguridad, limitecompra, estado) VALUES ('1982747364536562', 8, '201301', '201601', '842', 70500, 'anulada');
	INSERT INTO tarjeta (nrotarjeta, nrocliente, validadesde, validahasta, codseguridad, limitecompra, estado) VALUES ('9012348282748326', 9, '202010', '202310', '520', 21000, 'vigente');
	INSERT INTO tarjeta (nrotarjeta, nrocliente, validadesde, validahasta, codseguridad, limitecompra, estado) VALUES ('8274236578326572', 10, '201911', '202211', '7856', 14000, 'suspendida');
	INSERT INTO tarjeta (nrotarjeta, nrocliente, validadesde, validahasta, codseguridad, limitecompra, estado) VALUES ('8324732653627823', 11, '202204', '202504', '5848', 34500, 'vigente');
	INSERT INTO tarjeta (nrotarjeta, nrocliente, validadesde, validahasta, codseguridad, limitecompra, estado) VALUES ('7632462483439834', 12, '201212', '201512', '530', 48600, 'anulada');
	INSERT INTO tarjeta (nrotarjeta, nrocliente, validadesde, validahasta, codseguridad, limitecompra, estado) VALUES ('7216435365643723', 13, '200112', '200412', '740', 98000, 'anulada');
	INSERT INTO tarjeta (nrotarjeta, nrocliente, validadesde, validahasta, codseguridad, limitecompra, estado) VALUES ('3624376458982394', 14, '202007', '202307', '2485', 17900, 'suspendida');
	INSERT INTO tarjeta (nrotarjeta, nrocliente, validadesde, validahasta, codseguridad, limitecompra, estado) VALUES ('6347467826428439', 15, '202006', '202306', '773', 15000, 'vigente');
	INSERT INTO tarjeta (nrotarjeta, nrocliente, validadesde, validahasta, codseguridad, limitecompra, estado) VALUES ('5923892848377829', 16, '201003', '201303', '4575', 40000, 'suspendida');
	INSERT INTO tarjeta (nrotarjeta, nrocliente, validadesde, validahasta, codseguridad, limitecompra, estado) VALUES ('3784736427463790', 17, '202003', '202303', '2397', 38000, 'vigente');
	INSERT INTO tarjeta (nrotarjeta, nrocliente, validadesde, validahasta, codseguridad, limitecompra, estado) VALUES ('3209340294208483', 18, '202108', '202408', '608', 16500, 'suspendida');
	INSERT INTO tarjeta (nrotarjeta, nrocliente, validadesde, validahasta, codseguridad, limitecompra, estado) VALUES ('2743472943783934', 19, '202109', '202409', '453', 55500, 'vigente');
	INSERT INTO tarjeta (nrotarjeta, nrocliente, validadesde, validahasta, codseguridad, limitecompra, estado) VALUES ('8437284736756736', 19, '200202', '200502', '192', 66500, 'anulada');
	INSERT INTO tarjeta (nrotarjeta, nrocliente, validadesde, validahasta, codseguridad, limitecompra, estado) VALUES ('3464536272389793', 20, '202201', '202501', '1256', 13800, 'vigente');
	INSERT INTO tarjeta (nrotarjeta, nrocliente, validadesde, validahasta, codseguridad, limitecompra, estado) VALUES ('1273672467364703', 20, '200506', '200806', '9243', 98000, 'anulada');
	
	INSERT INTO cierre (año, mes, terminacion, fechainicio, fechacierre, fechavto) VALUES (2022, 1, 0, '2022-01-01', '2022-01-31', '2022-02-10');
	INSERT INTO cierre (año, mes, terminacion, fechainicio, fechacierre, fechavto) VALUES (2022, 2, 0, '2022-02-01', '2022-02-28', '2022-03-10');
	INSERT INTO cierre (año, mes, terminacion, fechainicio, fechacierre, fechavto) VALUES (2022, 3, 0, '2022-03-01', '2022-03-31', '2022-04-10');
	INSERT INTO cierre (año, mes, terminacion, fechainicio, fechacierre, fechavto) VALUES (2022, 4, 0, '2022-04-01', '2022-04-30', '2022-05-10');
	INSERT INTO cierre (año, mes, terminacion, fechainicio, fechacierre, fechavto) VALUES (2022, 5, 0, '2022-05-01', '2022-05-31', '2022-06-10');
	INSERT INTO cierre (año, mes, terminacion, fechainicio, fechacierre, fechavto) VALUES (2022, 6, 0, '2022-06-01', '2022-06-30', '2022-07-10');
	INSERT INTO cierre (año, mes, terminacion, fechainicio, fechacierre, fechavto) VALUES (2022, 7, 0, '2022-07-01', '2022-07-31', '2022-08-10');
	INSERT INTO cierre (año, mes, terminacion, fechainicio, fechacierre, fechavto) VALUES (2022, 8, 0, '2022-08-01', '2022-08-31', '2022-09-10');
	INSERT INTO cierre (año, mes, terminacion, fechainicio, fechacierre, fechavto) VALUES (2022, 9, 0, '2022-09-01', '2022-09-30', '2022-10-10');
	INSERT INTO cierre (año, mes, terminacion, fechainicio, fechacierre, fechavto) VALUES (2022, 10, 0, '2022-10-01', '2022-10-31', '2022-11-10');
	INSERT INTO cierre (año, mes, terminacion, fechainicio, fechacierre, fechavto) VALUES (2022, 11, 0, '2022-11-01', '2022-11-30', '2022-12-10');
	INSERT INTO cierre (año, mes, terminacion, fechainicio, fechacierre, fechavto) VALUES (2022, 12, 0, '2022-12-01', '2022-12-31', '2023-01-10');

	INSERT INTO cierre (año, mes, terminacion, fechainicio, fechacierre, fechavto) VALUES (2022, 1, 1, '2022-01-01', '2022-01-31', '2022-02-11');
	INSERT INTO cierre (año, mes, terminacion, fechainicio, fechacierre, fechavto) VALUES (2022, 2, 1, '2022-02-01', '2022-02-28', '2022-03-11');
	INSERT INTO cierre (año, mes, terminacion, fechainicio, fechacierre, fechavto) VALUES (2022, 3, 1, '2022-03-01', '2022-03-31', '2022-04-11');
	INSERT INTO cierre (año, mes, terminacion, fechainicio, fechacierre, fechavto) VALUES (2022, 4, 1, '2022-04-01', '2022-04-30', '2022-05-11');
	INSERT INTO cierre (año, mes, terminacion, fechainicio, fechacierre, fechavto) VALUES (2022, 5, 1, '2022-05-01', '2022-05-31', '2022-06-11');
	INSERT INTO cierre (año, mes, terminacion, fechainicio, fechacierre, fechavto) VALUES (2022, 6, 1, '2022-06-01', '2022-06-30', '2022-07-11');
	INSERT INTO cierre (año, mes, terminacion, fechainicio, fechacierre, fechavto) VALUES (2022, 7, 1, '2022-07-01', '2022-07-31', '2022-08-11');
	INSERT INTO cierre (año, mes, terminacion, fechainicio, fechacierre, fechavto) VALUES (2022, 8, 1, '2022-08-01', '2022-08-31', '2022-09-11');
	INSERT INTO cierre (año, mes, terminacion, fechainicio, fechacierre, fechavto) VALUES (2022, 9, 1, '2022-09-01', '2022-09-30', '2022-10-11');
	INSERT INTO cierre (año, mes, terminacion, fechainicio, fechacierre, fechavto) VALUES (2022, 10, 1, '2022-10-01', '2022-10-31', '2022-11-11');
	INSERT INTO cierre (año, mes, terminacion, fechainicio, fechacierre, fechavto) VALUES (2022, 11, 1, '2022-11-01', '2022-11-30', '2022-12-11');
	INSERT INTO cierre (año, mes, terminacion, fechainicio, fechacierre, fechavto) VALUES (2022, 12, 1, '2022-12-01', '2022-12-31', '2023-01-11');

	INSERT INTO cierre (año, mes, terminacion, fechainicio, fechacierre, fechavto) VALUES (2022, 1, 2, '2022-01-01', '2022-01-31', '2022-02-12');
	INSERT INTO cierre (año, mes, terminacion, fechainicio, fechacierre, fechavto) VALUES (2022, 2, 2, '2022-02-01', '2022-02-28', '2022-03-12');
	INSERT INTO cierre (año, mes, terminacion, fechainicio, fechacierre, fechavto) VALUES (2022, 3, 2, '2022-03-01', '2022-03-31', '2022-04-12');
	INSERT INTO cierre (año, mes, terminacion, fechainicio, fechacierre, fechavto) VALUES (2022, 4, 2, '2022-04-01', '2022-04-30', '2022-05-12');
	INSERT INTO cierre (año, mes, terminacion, fechainicio, fechacierre, fechavto) VALUES (2022, 5, 2, '2022-05-01', '2022-05-31', '2022-06-12');
	INSERT INTO cierre (año, mes, terminacion, fechainicio, fechacierre, fechavto) VALUES (2022, 6, 2, '2022-06-01', '2022-06-30', '2022-07-12');
	INSERT INTO cierre (año, mes, terminacion, fechainicio, fechacierre, fechavto) VALUES (2022, 7, 2, '2022-07-01', '2022-07-31', '2022-08-12');
	INSERT INTO cierre (año, mes, terminacion, fechainicio, fechacierre, fechavto) VALUES (2022, 8, 2, '2022-08-01', '2022-08-31', '2022-09-12');
	INSERT INTO cierre (año, mes, terminacion, fechainicio, fechacierre, fechavto) VALUES (2022, 9, 2, '2022-09-01', '2022-09-30', '2022-10-12');
	INSERT INTO cierre (año, mes, terminacion, fechainicio, fechacierre, fechavto) VALUES (2022, 10, 2, '2022-10-01', '2022-10-31', '2022-11-12');
	INSERT INTO cierre (año, mes, terminacion, fechainicio, fechacierre, fechavto) VALUES (2022, 11, 2, '2022-11-01', '2022-11-30', '2022-12-12');
	INSERT INTO cierre (año, mes, terminacion, fechainicio, fechacierre, fechavto) VALUES (2022, 12, 2, '2022-12-01', '2022-12-31', '2023-01-12');

	INSERT INTO cierre (año, mes, terminacion, fechainicio, fechacierre, fechavto) VALUES (2022, 1, 3, '2022-01-01', '2022-01-31', '2022-02-13');
	INSERT INTO cierre (año, mes, terminacion, fechainicio, fechacierre, fechavto) VALUES (2022, 2, 3, '2022-02-01', '2022-02-28', '2022-03-13');
	INSERT INTO cierre (año, mes, terminacion, fechainicio, fechacierre, fechavto) VALUES (2022, 3, 3, '2022-03-01', '2022-03-31', '2022-04-13');
	INSERT INTO cierre (año, mes, terminacion, fechainicio, fechacierre, fechavto) VALUES (2022, 4, 3, '2022-04-01', '2022-04-30', '2022-05-13');
	INSERT INTO cierre (año, mes, terminacion, fechainicio, fechacierre, fechavto) VALUES (2022, 5, 3, '2022-05-01', '2022-05-31', '2022-06-13');
	INSERT INTO cierre (año, mes, terminacion, fechainicio, fechacierre, fechavto) VALUES (2022, 6, 3, '2022-06-01', '2022-06-30', '2022-07-13');
	INSERT INTO cierre (año, mes, terminacion, fechainicio, fechacierre, fechavto) VALUES (2022, 7, 3, '2022-07-01', '2022-07-31', '2022-08-13');
	INSERT INTO cierre (año, mes, terminacion, fechainicio, fechacierre, fechavto) VALUES (2022, 8, 3, '2022-08-01', '2022-08-31', '2022-09-13');
	INSERT INTO cierre (año, mes, terminacion, fechainicio, fechacierre, fechavto) VALUES (2022, 9, 3, '2022-09-01', '2022-09-30', '2022-10-13');
	INSERT INTO cierre (año, mes, terminacion, fechainicio, fechacierre, fechavto) VALUES (2022, 10, 3, '2022-10-01', '2022-10-31', '2022-11-13');
	INSERT INTO cierre (año, mes, terminacion, fechainicio, fechacierre, fechavto) VALUES (2022, 11, 3, '2022-11-01', '2022-11-30', '2022-12-13');
	INSERT INTO cierre (año, mes, terminacion, fechainicio, fechacierre, fechavto) VALUES (2022, 12, 3, '2022-12-01', '2022-12-31', '2023-01-13');

	INSERT INTO cierre (año, mes, terminacion, fechainicio, fechacierre, fechavto) VALUES (2022, 1, 4, '2022-01-01', '2022-01-31', '2022-02-14');
	INSERT INTO cierre (año, mes, terminacion, fechainicio, fechacierre, fechavto) VALUES (2022, 2, 4, '2022-02-01', '2022-02-28', '2022-03-14');
	INSERT INTO cierre (año, mes, terminacion, fechainicio, fechacierre, fechavto) VALUES (2022, 3, 4, '2022-03-01', '2022-03-31', '2022-04-14');
	INSERT INTO cierre (año, mes, terminacion, fechainicio, fechacierre, fechavto) VALUES (2022, 4, 4, '2022-04-01', '2022-04-30', '2022-05-14');
	INSERT INTO cierre (año, mes, terminacion, fechainicio, fechacierre, fechavto) VALUES (2022, 5, 4, '2022-05-01', '2022-05-31', '2022-06-14');
	INSERT INTO cierre (año, mes, terminacion, fechainicio, fechacierre, fechavto) VALUES (2022, 6, 4, '2022-06-01', '2022-06-30', '2022-07-14');
	INSERT INTO cierre (año, mes, terminacion, fechainicio, fechacierre, fechavto) VALUES (2022, 7, 4, '2022-07-01', '2022-07-31', '2022-08-14');
	INSERT INTO cierre (año, mes, terminacion, fechainicio, fechacierre, fechavto) VALUES (2022, 8, 4, '2022-08-01', '2022-08-31', '2022-09-14');
	INSERT INTO cierre (año, mes, terminacion, fechainicio, fechacierre, fechavto) VALUES (2022, 9, 4, '2022-09-01', '2022-09-30', '2022-10-14');
	INSERT INTO cierre (año, mes, terminacion, fechainicio, fechacierre, fechavto) VALUES (2022, 10, 4, '2022-10-01', '2022-10-31', '2022-11-14');
	INSERT INTO cierre (año, mes, terminacion, fechainicio, fechacierre, fechavto) VALUES (2022, 11, 4, '2022-11-01', '2022-11-30', '2022-12-14');
	INSERT INTO cierre (año, mes, terminacion, fechainicio, fechacierre, fechavto) VALUES (2022, 12, 4, '2022-12-01', '2022-12-31', '2023-01-14');

	INSERT INTO cierre (año, mes, terminacion, fechainicio, fechacierre, fechavto) VALUES (2022, 1, 5, '2022-01-01', '2022-01-31', '2022-02-15');
	INSERT INTO cierre (año, mes, terminacion, fechainicio, fechacierre, fechavto) VALUES (2022, 2, 5, '2022-02-01', '2022-02-28', '2022-03-15');
	INSERT INTO cierre (año, mes, terminacion, fechainicio, fechacierre, fechavto) VALUES (2022, 3, 5, '2022-03-01', '2022-03-31', '2022-04-15');
	INSERT INTO cierre (año, mes, terminacion, fechainicio, fechacierre, fechavto) VALUES (2022, 4, 5, '2022-04-01', '2022-04-30', '2022-05-15');
	INSERT INTO cierre (año, mes, terminacion, fechainicio, fechacierre, fechavto) VALUES (2022, 5, 5, '2022-05-01', '2022-05-31', '2022-06-15');
	INSERT INTO cierre (año, mes, terminacion, fechainicio, fechacierre, fechavto) VALUES (2022, 6, 5, '2022-06-01', '2022-06-30', '2022-07-15');
	INSERT INTO cierre (año, mes, terminacion, fechainicio, fechacierre, fechavto) VALUES (2022, 7, 5, '2022-07-01', '2022-07-31', '2022-08-15');
	INSERT INTO cierre (año, mes, terminacion, fechainicio, fechacierre, fechavto) VALUES (2022, 8, 5, '2022-08-01', '2022-08-31', '2022-09-15');
	INSERT INTO cierre (año, mes, terminacion, fechainicio, fechacierre, fechavto) VALUES (2022, 9, 5, '2022-09-01', '2022-09-30', '2022-10-15');
	INSERT INTO cierre (año, mes, terminacion, fechainicio, fechacierre, fechavto) VALUES (2022, 10, 5, '2022-10-01', '2022-10-31', '2022-11-15');
	INSERT INTO cierre (año, mes, terminacion, fechainicio, fechacierre, fechavto) VALUES (2022, 11, 5, '2022-11-01', '2022-11-30', '2022-12-15');
	INSERT INTO cierre (año, mes, terminacion, fechainicio, fechacierre, fechavto) VALUES (2022, 12, 5, '2022-12-01', '2022-12-31', '2023-01-15');
	
	INSERT INTO cierre (año, mes, terminacion, fechainicio, fechacierre, fechavto) VALUES (2022, 1, 6, '2022-01-01', '2022-01-31', '2022-02-16');
	INSERT INTO cierre (año, mes, terminacion, fechainicio, fechacierre, fechavto) VALUES (2022, 2, 6, '2022-02-01', '2022-02-28', '2022-03-16');
	INSERT INTO cierre (año, mes, terminacion, fechainicio, fechacierre, fechavto) VALUES (2022, 3, 6, '2022-03-01', '2022-03-31', '2022-04-16');
	INSERT INTO cierre (año, mes, terminacion, fechainicio, fechacierre, fechavto) VALUES (2022, 4, 6, '2022-04-01', '2022-04-30', '2022-05-16');
	INSERT INTO cierre (año, mes, terminacion, fechainicio, fechacierre, fechavto) VALUES (2022, 5, 6, '2022-05-01', '2022-05-31', '2022-06-16');
	INSERT INTO cierre (año, mes, terminacion, fechainicio, fechacierre, fechavto) VALUES (2022, 6, 6, '2022-06-01', '2022-06-30', '2022-07-16');
	INSERT INTO cierre (año, mes, terminacion, fechainicio, fechacierre, fechavto) VALUES (2022, 7, 6, '2022-07-01', '2022-07-31', '2022-08-16');
	INSERT INTO cierre (año, mes, terminacion, fechainicio, fechacierre, fechavto) VALUES (2022, 8, 6, '2022-08-01', '2022-08-31', '2022-09-16');
	INSERT INTO cierre (año, mes, terminacion, fechainicio, fechacierre, fechavto) VALUES (2022, 9, 6, '2022-09-01', '2022-09-30', '2022-10-16');
	INSERT INTO cierre (año, mes, terminacion, fechainicio, fechacierre, fechavto) VALUES (2022, 10, 6, '2022-10-01', '2022-10-31', '2022-11-16');
	INSERT INTO cierre (año, mes, terminacion, fechainicio, fechacierre, fechavto) VALUES (2022, 11, 6, '2022-11-01', '2022-11-30', '2022-12-16');
	INSERT INTO cierre (año, mes, terminacion, fechainicio, fechacierre, fechavto) VALUES (2022, 12, 6, '2022-12-01', '2022-12-31', '2023-01-16');
	
	INSERT INTO cierre (año, mes, terminacion, fechainicio, fechacierre, fechavto) VALUES (2022, 1, 7, '2022-01-01', '2022-01-31', '2022-02-17');
	INSERT INTO cierre (año, mes, terminacion, fechainicio, fechacierre, fechavto) VALUES (2022, 2, 7, '2022-02-01', '2022-02-28', '2022-03-17');
	INSERT INTO cierre (año, mes, terminacion, fechainicio, fechacierre, fechavto) VALUES (2022, 3, 7, '2022-03-01', '2022-03-31', '2022-04-17');
	INSERT INTO cierre (año, mes, terminacion, fechainicio, fechacierre, fechavto) VALUES (2022, 4, 7, '2022-04-01', '2022-04-30', '2022-05-17');
	INSERT INTO cierre (año, mes, terminacion, fechainicio, fechacierre, fechavto) VALUES (2022, 5, 7, '2022-05-01', '2022-05-31', '2022-06-17');
	INSERT INTO cierre (año, mes, terminacion, fechainicio, fechacierre, fechavto) VALUES (2022, 6, 7, '2022-06-01', '2022-06-30', '2022-07-17');
	INSERT INTO cierre (año, mes, terminacion, fechainicio, fechacierre, fechavto) VALUES (2022, 7, 7, '2022-07-01', '2022-07-31', '2022-08-17');
	INSERT INTO cierre (año, mes, terminacion, fechainicio, fechacierre, fechavto) VALUES (2022, 8, 7, '2022-08-01', '2022-08-31', '2022-09-17');
	INSERT INTO cierre (año, mes, terminacion, fechainicio, fechacierre, fechavto) VALUES (2022, 9, 7, '2022-09-01', '2022-09-30', '2022-10-17');
	INSERT INTO cierre (año, mes, terminacion, fechainicio, fechacierre, fechavto) VALUES (2022, 10, 7, '2022-10-01', '2022-10-31', '2022-11-17');
	INSERT INTO cierre (año, mes, terminacion, fechainicio, fechacierre, fechavto) VALUES (2022, 11, 7, '2022-11-01', '2022-11-30', '2022-12-17');
	INSERT INTO cierre (año, mes, terminacion, fechainicio, fechacierre, fechavto) VALUES (2022, 12, 7, '2022-12-01', '2022-12-31', '2023-01-17');
	
	INSERT INTO cierre (año, mes, terminacion, fechainicio, fechacierre, fechavto) VALUES (2022, 1, 8, '2022-01-01', '2022-01-31', '2022-02-18');
	INSERT INTO cierre (año, mes, terminacion, fechainicio, fechacierre, fechavto) VALUES (2022, 2, 8, '2022-02-01', '2022-02-28', '2022-03-18');
	INSERT INTO cierre (año, mes, terminacion, fechainicio, fechacierre, fechavto) VALUES (2022, 3, 8, '2022-03-01', '2022-03-31', '2022-04-18');
	INSERT INTO cierre (año, mes, terminacion, fechainicio, fechacierre, fechavto) VALUES (2022, 4, 8, '2022-04-01', '2022-04-30', '2022-05-18');
	INSERT INTO cierre (año, mes, terminacion, fechainicio, fechacierre, fechavto) VALUES (2022, 5, 8, '2022-05-01', '2022-05-31', '2022-06-18');
	INSERT INTO cierre (año, mes, terminacion, fechainicio, fechacierre, fechavto) VALUES (2022, 6, 8, '2022-06-01', '2022-06-30', '2022-07-18');
	INSERT INTO cierre (año, mes, terminacion, fechainicio, fechacierre, fechavto) VALUES (2022, 7, 8, '2022-07-01', '2022-07-31', '2022-08-18');
	INSERT INTO cierre (año, mes, terminacion, fechainicio, fechacierre, fechavto) VALUES (2022, 8, 8, '2022-08-01', '2022-08-31', '2022-09-18');
	INSERT INTO cierre (año, mes, terminacion, fechainicio, fechacierre, fechavto) VALUES (2022, 9, 8, '2022-09-01', '2022-09-30', '2022-10-18');
	INSERT INTO cierre (año, mes, terminacion, fechainicio, fechacierre, fechavto) VALUES (2022, 10, 8, '2022-10-01', '2022-10-31', '2022-11-18');
	INSERT INTO cierre (año, mes, terminacion, fechainicio, fechacierre, fechavto) VALUES (2022, 11, 8, '2022-11-01', '2022-11-30', '2022-12-18');
	INSERT INTO cierre (año, mes, terminacion, fechainicio, fechacierre, fechavto) VALUES (2022, 12, 8, '2022-12-01', '2022-12-31', '2023-01-18');
	
	INSERT INTO cierre (año, mes, terminacion, fechainicio, fechacierre, fechavto) VALUES (2022, 1, 9, '2022-01-01', '2022-01-31', '2022-02-19');
	INSERT INTO cierre (año, mes, terminacion, fechainicio, fechacierre, fechavto) VALUES (2022, 2, 9, '2022-02-01', '2022-02-28', '2022-03-19');
	INSERT INTO cierre (año, mes, terminacion, fechainicio, fechacierre, fechavto) VALUES (2022, 3, 9, '2022-03-01', '2022-03-31', '2022-04-19');
	INSERT INTO cierre (año, mes, terminacion, fechainicio, fechacierre, fechavto) VALUES (2022, 4, 9, '2022-04-01', '2022-04-30', '2022-05-19');
	INSERT INTO cierre (año, mes, terminacion, fechainicio, fechacierre, fechavto) VALUES (2022, 5, 9, '2022-05-01', '2022-05-31', '2022-06-19');
	INSERT INTO cierre (año, mes, terminacion, fechainicio, fechacierre, fechavto) VALUES (2022, 6, 9, '2022-06-01', '2022-06-30', '2022-07-19');
	INSERT INTO cierre (año, mes, terminacion, fechainicio, fechacierre, fechavto) VALUES (2022, 7, 9, '2022-07-01', '2022-07-31', '2022-08-19');
	INSERT INTO cierre (año, mes, terminacion, fechainicio, fechacierre, fechavto) VALUES (2022, 8, 9, '2022-08-01', '2022-08-31', '2022-09-19');
	INSERT INTO cierre (año, mes, terminacion, fechainicio, fechacierre, fechavto) VALUES (2022, 9, 9, '2022-09-01', '2022-09-30', '2022-10-19');
	INSERT INTO cierre (año, mes, terminacion, fechainicio, fechacierre, fechavto) VALUES (2022, 10, 9, '2022-10-01', '2022-10-31', '2022-11-19');
	INSERT INTO cierre (año, mes, terminacion, fechainicio, fechacierre, fechavto) VALUES (2022, 11, 9, '2022-11-01', '2022-11-30', '2022-12-19');
	INSERT INTO cierre (año, mes, terminacion, fechainicio, fechacierre, fechavto) VALUES (2022, 12, 9, '2022-12-01', '2022-12-31', '2023-01-19');
	
	INSERT INTO consumo (nrotarjeta, codseguridad, nrocomercio, monto) VALUES ('5827624652643290','8246', 3, 3500);
	INSERT INTO consumo (nrotarjeta, codseguridad, nrocomercio, monto) VALUES ('2930489589375756','235', 7, 10000);
	INSERT INTO consumo (nrotarjeta, codseguridad, nrocomercio, monto) VALUES ('9320584378172604','906', 1, 7000);
	INSERT INTO consumo (nrotarjeta, codseguridad, nrocomercio, monto) VALUES ('1294382965767449','2364', 6, 3500);
	INSERT INTO consumo (nrotarjeta, codseguridad, nrocomercio, monto) VALUES ('6347467826428439','546', 2, 960);
	INSERT INTO consumo (nrotarjeta, codseguridad, nrocomercio, monto) VALUES ('3464536272389793','1256', 1, 25000);
	INSERT INTO consumo (nrotarjeta, codseguridad, nrocomercio, monto) VALUES ('3928217404085943','1485', 1, 5000);
	INSERT INTO consumo (nrotarjeta, codseguridad, nrocomercio, monto) VALUES ('3928217404085943','1485', 2, 7000);
	INSERT INTO consumo (nrotarjeta, codseguridad, nrocomercio, monto) VALUES ('9012348282748326','520', 3, 10000);
	INSERT INTO consumo (nrotarjeta, codseguridad, nrocomercio, monto) VALUES ('9012348282748326','520', 4, 9000);
	INSERT INTO consumo (nrotarjeta, codseguridad, nrocomercio, monto) VALUES ('3784736427463790','2397', 9, 50000);
	INSERT INTO consumo (nrotarjeta, codseguridad, nrocomercio, monto) VALUES ('3784736427463790','2397', 9, 80000);`)
	if err != nil {
		log.Fatal(err)
	}
}

func createFunctionAutorizaciones() {
	_, err = db.Exec(`CREATE OR REPLACE FUNCTION autorizar_compra(n_tarjeta tarjeta.nrotarjeta%type,
											cod_seg tarjeta.codseguridad%type,
												n_comercio compra.nrocomercio%type,
													monto_compra compra.monto%type) RETURNS boolean as $$
					DECLARE
						tarjeta_fila record;
						fecha_actual DATE;
						fecha_vencimiento DATE;
						comercio_encontrado INT;
						fecha_de_vencimiento_text TEXT;
						monto_total_compras_tarjeta_actual compra.monto%type;
					BEGIN

						SELECT * INTO tarjeta_fila FROM tarjeta t WHERE n_tarjeta = t.nrotarjeta;
						
						IF NOT found then
							INSERT INTO rechazo (nrotarjeta, nrocomercio, fecha, monto, motivo) 
								VALUES (n_tarjeta, n_comercio, current_timestamp, monto_compra, 'Tarjeta no valida o no vigente.');
						
							return false;
						ELSE
							
							SELECT CAST(validahasta AS TEXT) INTO fecha_de_vencimiento_text FROM tarjeta t WHERE n_tarjeta = t.nrotarjeta;

							SELECT SUM(monto) INTO monto_total_compras_tarjeta_actual FROM compra c WHERE n_tarjeta = c.nrotarjeta and c.pagado = false;
							IF monto_total_compras_tarjeta_actual IS NULL then
								monto_total_compras_tarjeta_actual := 0;
							END IF;

							fecha_actual := CURRENT_DATE;
							fecha_vencimiento := TO_DATE(fecha_de_vencimiento_text, 'YYYYMM');

							IF tarjeta_fila.codseguridad != cod_seg then
								INSERT INTO rechazo (nrotarjeta, nrocomercio, fecha, monto, motivo)  
									VALUES (n_tarjeta, n_comercio, current_timestamp, monto_compra, 'Codigo de seguridad invalido.');
								
								return false;
							
							ELSIF monto_compra + monto_total_compras_tarjeta_actual  >= tarjeta_fila.limitecompra then
								INSERT INTO rechazo (nrotarjeta, nrocomercio, fecha, monto, motivo)
									VALUES (n_tarjeta, n_comercio, current_timestamp, monto_compra, 'Supera límite de tarjeta');

								return false;
							
							ELSIF  fecha_actual > fecha_vencimiento then
								INSERT INTO rechazo (nrotarjeta, nrocomercio, fecha, monto, motivo)
									VALUES (n_tarjeta, n_comercio, current_timestamp, monto_compra, 'Plazo de vigencia expirado.');
							
								return false;

							ELSIF tarjeta_fila.estado = 'suspendida' then
								INSERT INTO rechazo (nrotarjeta, nrocomercio, fecha, monto, motivo)
									VALUES (n_tarjeta, n_comercio, current_timestamp, monto_compra, 'La Tarjeta se encuentra suspendida.');
							
								return false;
							
							ELSE
								INSERT INTO compra (nrotarjeta, nrocomercio, fecha, monto, pagado)
									VALUES (n_tarjeta, n_comercio, current_timestamp, monto_compra, false);

								return true;
							END IF;
						END IF;
					END;
					$$ LANGUAGE plpgsql;`)
	if err != nil {
		log.Fatal(err)
	}
}

func createFunctionResumenes() {
	_, err = db.Exec(`CREATE OR REPLACE FUNCTION generar_resumen (n_cliente cliente.nrocliente%TYPE, aux_año INT, aux_mes INT) RETURNS void AS $$

						DECLARE

							n_linea INT := 1;
							aux_cliente RECORD;
							aux_compra RECORD;
							aux_tarjeta RECORD;
							aux_cierre RECORD;
							aux_comercio RECORD;
							n_resumen cabecera.nroresumen%type;
							monto_total cabecera.total%type;

						BEGIN

							SELECT * INTO aux_cliente FROM cliente WHERE nrocliente = n_cliente;
								IF NOT FOUND THEN
									RAISE 'El número de cliente % no existe.', n_cliente;
								END IF;

							FOR aux_tarjeta IN SELECT * FROM tarjeta WHERE nrocliente = aux_cliente.nrocliente LOOP
						
								monto_total := 0;

								SELECT * INTO aux_cierre FROM cierre WHERE año = aux_año AND mes = aux_mes
									AND terminacion = substring(aux_tarjeta.nrotarjeta, 16, 1)::INT;

								INSERT INTO cabecera (nombre, apellido, domicilio, nrotarjeta, desde, hasta, vence, total)
									VALUES (aux_cliente.nombre, aux_cliente.apellido, aux_cliente.domicilio, aux_tarjeta.nrotarjeta, aux_cierre.fechainicio, aux_cierre.fechacierre, aux_cierre.fechavto, monto_total);

								SELECT nroresumen INTO n_resumen FROM cabecera WHERE nrotarjeta = aux_tarjeta.nrotarjeta
									AND desde = aux_cierre.fechainicio AND hasta = aux_cierre.fechacierre;

								FOR aux_compra IN SELECT * FROM compra WHERE nrotarjeta = aux_tarjeta.nrotarjeta AND fecha >= aux_cierre.fechainicio AND fecha <= aux_cierre.fechacierre AND pagado = false LOOP

									SELECT * INTO aux_comercio FROM comercio WHERE nrocomercio = aux_compra.nrocomercio;

									INSERT INTO detalle (nroresumen, nrolinea, fecha, nombrecomercio, monto)
										VALUES (n_resumen, n_linea, aux_compra.fecha, aux_comercio.nombre, aux_compra.monto);
									n_linea := n_linea + 1;
									monto_total := monto_total + aux_compra.monto;
									UPDATE compra SET pagado = true WHERE nrooperacion = aux_compra.nrooperacion;
								END LOOP;

								UPDATE cabecera SET total = monto_total WHERE nrotarjeta = aux_tarjeta.nrotarjeta
										AND desde = aux_cierre.fechainicio AND hasta = aux_cierre.fechacierre;
										
								END LOOP;
							END;
					$$ LANGUAGE plpgsql;`)
	if err != nil {
		log.Fatal(err)
	}
}

func createFunctionAlertas() {
	_, err = db.Exec(`CREATE OR REPLACE FUNCTION alerta_rechazo() RETURNS TRIGGER AS $$
						BEGIN
							
							PERFORM * FROM rechazo r WHERE r.nrotarjeta = NEW.nrotarjeta
															AND r.nrorechazo != NEW.nrorechazo
															AND EXTRACT(DAY FROM r.fecha) = EXTRACT(DAY FROM new.fecha)
															AND EXTRACT(MONTH FROM r.fecha) = EXTRACT(MONTH FROM new.fecha)
															AND EXTRACT(YEAR FROM r.fecha) = EXTRACT(YEAR FROM new.fecha)
															AND r.motivo = 'Supera límite de tarjeta'
															AND new.motivo = 'Supera límite de tarjeta';
							IF FOUND THEN
								UPDATE tarjeta SET estado = 'suspendida' WHERE nrotarjeta = new.nrotarjeta;
								INSERT INTO alerta (nrotarjeta, fecha, nrorechazo, codalerta, descripcion) VALUES (new.nrotarjeta, CURRENT_TIMESTAMP, new.nrorechazo, 32, 'Tarjeta suspendida por compras excedidas del límite');
							END IF;
							INSERT INTO alerta (nrotarjeta, fecha, nrorechazo, codalerta, descripcion) VALUES (NEW.nrotarjeta, CURRENT_TIMESTAMP, NEW.nrorechazo, 0, NEW.motivo);
							RETURN NEW;
						END;
						$$ LANGUAGE plpgsql;
						
						CREATE OR REPLACE FUNCTION es_mismo_dia(fecha1 TIMESTAMP, fecha2 TIMESTAMP) RETURNS BOOLEAN AS $$
						DECLARE
							anio_fecha1 INT;
							anio_fecha2 INT;
							mes_fecha1 INT;
							mes_fecha2 INT;
							dia_fecha1 INT;
							dia_fecha2 INT;
						BEGIN
							SELECT EXTRACT INTO anio_fecha1 (YEAR FROM fecha1);
							SELECT EXTRACT INTO anio_fecha2 (YEAR FROM fecha2);
							SELECT EXTRACT INTO mes_fecha1 (MONTH FROM fecha1);
							SELECT EXTRACT INTO mes_fecha2 (MONTH FROM fecha2);
							SELECT EXTRACT INTO dia_fecha1 (DAY FROM fecha1);
							SELECT EXTRACT INTO dia_fecha2 (DAY FROM fecha2);
							
							IF anio_fecha1 = anio_fecha2 AND mes_fecha1 = mes_fecha2 AND dia_fecha1 = dia_fecha2 THEN
								RETURN TRUE;
							ELSE
								RETURN  FALSE;
							END IF;
						END;
						$$ LANGUAGE plpgsql;
						
						
						CREATE OR REPLACE FUNCTION alerta_compras() RETURNS TRIGGER AS $$
						DECLARE
							ultima_compra record;
							codigo_postal_ultima_compra comercio.codigopostal%TYPE;
							codigo_postal_compra_actual comercio.codigopostal%TYPE;
							diferencia_minutos INT;
							mismo_dia BOOLEAN;
							
						BEGIN
							SELECT INTO ultima_compra * FROM compra WHERE nrotarjeta = new.nrotarjeta ORDER BY fecha DESC LIMIT 1;
							SELECT INTO codigo_postal_ultima_compra codigopostal FROM comercio WHERE nrocomercio = ultima_compra.nrocomercio;
							SELECT INTO codigo_postal_compra_actual codigopostal FROM comercio WHERE nrocomercio = new.nrocomercio;
							mismo_dia := es_mismo_dia(NEW.fecha, ultima_compra.fecha);
							SELECT EXTRACT INTO diferencia_minutos (MINUTES FROM (NEW.fecha - ultima_compra.fecha));
							
							IF NEW.nrocomercio != ultima_compra.nrocomercio AND codigo_postal_compra_actual = codigo_postal_ultima_compra AND mismo_dia = true AND diferencia_minutos < 1 THEN
								INSERT INTO alerta (nrotarjeta, fecha, codalerta, descripcion) VALUES (NEW.nrotarjeta, CURRENT_TIMESTAMP, 1, 'Se realizaron dos compras en el mismo minuto en tiendas distintas');
							END IF;
							
							IF codigo_postal_compra_actual != codigo_postal_ultima_compra AND mismo_dia = true AND diferencia_minutos < 5 THEN
								INSERT INTO alerta (nrotarjeta, fecha, codalerta, descripcion) VALUES (NEW.nrotarjeta, CURRENT_TIMESTAMP, 5, 'Se realizaron dos compras en 5 minutos en localidades distintas');
							END IF;
							RETURN NEW;
						END;
						$$ LANGUAGE plpgsql;
						
						CREATE OR REPLACE TRIGGER alerta_rechazo_trigger
						AFTER INSERT ON rechazo
						FOR EACH ROW
						EXECUTE PROCEDURE alerta_rechazo();
						
						CREATE OR REPLACE TRIGGER alerta_compras_trigger
						BEFORE INSERT ON compra
						FOR EACH ROW
						EXECUTE PROCEDURE alerta_compras();`)
	if err != nil {
		log.Fatal(err)
	}
}

func autorizarCompras(){
	
	for TotalDeConsumos != 12 {
	/* Busqueda de datos por nroConsumo */
		row, err := db.Query(`SELECT * FROM consumo WHERE nroconsumo = $1`, nroConsumo)
		
	//Control de error de query anterior
		if err != nil{ 
			log.Fatal(err) 
		} 
		
	defer row.Close() //Close de consumo de row al terminar.
	var c consumo // var de type struc consumo.
	var retorno bool // var para mostar si se autorizo o no una compra.
	
	// Asignacion de datos de query a struc c.
		for row.Next(){
			if err = row.Scan(&c.nroconsumo, &c.nrotarjeta, &c.codseguridad, &c.nrocomercio, &c.monto); 
			//Control de error
			err != nil{
				log.Fatal(err)
			}
		}
	
	    //Consumo de funcion autorizar_compra()
		row, err = db.Query(`SELECT autorizar_compra($1::CHAR(16), $2::char(4), $3::INT, $4::DECIMAL(7,2));`, c.nrotarjeta, c.codseguridad, c.nrocomercio, c.monto)
		//Control de error
		if err != nil{ 
			log.Fatal(err)
		}
		
		// Obtencion de return de query
		if row.Next() {
			if err = row.Scan(&retorno); err != nil { //Control de error
				log.Fatal(err)
			}
		}
		//Control de error
		if err != nil{ 
			log.Fatal(err)
		}
		
		if !retorno {
			fmt.Print("\nCompra Rechazada: ", nroConsumo, " (numero de consumo)")
		} 
		
		//Aumento de nroConsumo mostrado
		nroConsumo = nroConsumo + 1
		
		//Suma hasta llegar al tope total de consumos
		TotalDeConsumos = TotalDeConsumos + 1
	}
	
	fmt.Print("\n-- No hay más compras para autorizar\n")
	
}

func informarAlertaNueva() {
	for opcion != 9 {
		var countAlertas int = 0
		
		row, err := db.Query(`SELECT count(*) FROM alerta;`)
		if err != nil {
			log.Fatal(err)
		}
		if row.Next() {
			if err = row.Scan(&countAlertas); err != nil {
				log.Fatal(err)
			}
		}
		if err != nil {
			log.Fatal(err)
		}
		if cantidadAlertas < countAlertas {
			cantidadAlertas = countAlertas
			fmt.Print("\n Nueva alerta detectada \n")
		}
		time.Sleep(60 * time.Second)
	}
}

func generarResumenes() {
	_, err = db.Query(`SELECT generar_resumen(1, 2022, 2);
				SELECT generar_resumen(4, 2022, 9);
				SELECT generar_resumen(5, 2022, 6);
				SELECT generar_resumen(6, 2022, 8);
				SELECT generar_resumen(11, 2022, 12);`)
	if err != nil {
		log.Fatal(err)
	}
}