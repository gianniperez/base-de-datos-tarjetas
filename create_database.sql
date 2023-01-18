DROP DATABASE IF EXISTS tarjetas;
CREATE DATABASE tarjetas;
\c tarjetas
CREATE TABLE cliente (nrocliente SERIAL, nombre TEXT, apellido TEXT, domicilio TEXT, telefono CHAR(12));
CREATE TABLE tarjeta (nrotarjeta CHAR(16), nrocliente INT, validadesde CHAR(6), validahasta CHAR(6), codseguridad CHAR(4), limitecompra DECIMAL(8,2), estado CHAR(10));
CREATE TABLE comercio (nrocomercio SERIAL, nombre TEXT, domicilio TEXT, codigopostal CHAR(8), telefono CHAR(12));
CREATE TABLE compra (nrooperacion SERIAL, nrotarjeta CHAR(16), nrocomercio INT, fecha TIMESTAMP, monto DECIMAL(7,2), pagado BOOLEAN);
CREATE TABLE rechazo (nrorechazo SERIAL, nrotarjeta CHAR(16), nrocomercio INT, fecha TIMESTAMP, monto DECIMAL(7,2), motivo TEXT);
CREATE TABLE cierre (año INT, mes INT, terminacion INT, fechainicio DATE, fechacierre DATE, fechavto DATE);
CREATE TABLE cabecera (nroresumen SERIAL, nombre TEXT, apellido TEXT, domicilio TEXT, nrotarjeta CHAR(16), desde DATE, hasta DATE, vence DATE, total DECIMAL(8,2));
CREATE TABLE detalle (nroresumen INT, nrolinea INT, fecha DATE, nombrecomercio TEXT, monto DECIMAL(7,2));
CREATE TABLE alerta (nroalerta SERIAL, nrotarjeta CHAR(16), fecha TIMESTAMP, nrorechazo INT, codalerta INT, descripcion TEXT);
CREATE TABLE consumo (nrotarjeta CHAR(16), codseguridad CHAR(4), nrocomercio INT, monto DECIMAL(7,2));