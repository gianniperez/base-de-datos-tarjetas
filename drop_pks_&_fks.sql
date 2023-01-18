-- FK’s Drops

ALTER TABLE tarjeta DROP CONSTRAINT tarjeta_nrocliente_fk;
ALTER TABLE compra DROP CONSTRAINT compra_nrotarjeta_fk;
ALTER TABLE compra DROP CONSTRAINT compra_nrocomercio_fk;
ALTER TABLE rechazo DROP CONSTRAINT rechazo_nrocomercio_fk;
ALTER TABLE cabecera DROP CONSTRAINT cabecera_nrotarjeta_fk;
ALTER TABLE alerta DROP CONSTRAINT alerta_nrorechazo_fk;

-- PK’s Drops

ALTER TABLE cliente DROP CONSTRAINT cliente_pk;
ALTER TABLE tarjeta DROP CONSTRAINT tarjeta_pk;
ALTER TABLE comercio DROP CONSTRAINT comercio_pk;
ALTER TABLE compra DROP CONSTRAINT compra_pk;
ALTER TABLE rechazo DROP CONSTRAINT rechazo_pk;
ALTER TABLE cierre DROP CONSTRAINT cierre_pk;
ALTER TABLE cabecera DROP CONSTRAINT cabecera_pk;
ALTER TABLE detalle DROP CONSTRAINT detalle_pk;
ALTER TABLE alerta DROP CONSTRAINT alerta_pk;

