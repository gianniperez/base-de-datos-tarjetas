
CREATE OR REPLACE FUNCTION alerta_rechazo() RETURNS TRIGGER AS $$
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
EXECUTE PROCEDURE alerta_compras();
