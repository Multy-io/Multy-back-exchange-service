




-- CREATE UNIQUE INDEX exchanges_id_uindex ON public.exchanges (id);

CREATE OR REPLACE FUNCTION fill_rates()
    RETURNS void AS $$
    BEGIN

INSERT INTO exchanges(title, create_date)
SELECT DISTINCT exchange_title,  time_stamp
FROM sa_rates
ON CONFLICT DO NOTHING;

INSERT INTO currencies(title, code, create_date, native_id)
SELECT DISTINCT target_title, target_code,  time_stamp, target_native_id
  from sa_rates
UNION
  SELECT DISTINCT reference_title, reference_code,  time_stamp, reference_native_id
  from sa_rates
ON CONFLICT DO NOTHING;

INSERT INTO exchanges_pairs
SELECT DISTINCT on (ex.id, tcr.id, rcr.id ) ex.id exchange_id, tcr.id target_id, rcr.id reference_id, time_stamp
FROM sa_rates sr, exchanges ex, currencies tcr,  currencies rcr
WHERE sr.exchange_title = ex.title and sr.target_code = tcr.code and sr.reference_code = rcr.code
 ON CONFLICT DO NOTHING;

INSERT INTO rates
SELECT ex.id exchange_id, tcr.id target_id, rcr.id reference_id, sr.time_stamp, sr.rate
FROM sa_rates sr, exchanges ex, currencies tcr,  currencies rcr
WHERE sr.exchange_title = ex.title and sr.target_code = tcr.code and sr.reference_code = rcr.code;
      TRUNCATE sa_rates;
    END;
    $$ LANGUAGE plpgsql;


CREATE OR REPLACE FUNCTION getRates(p_time_stamp TIMESTAMP WITH TIME ZONE, p_exchange_title VARCHAR, p_target_code VARCHAR, p_referencies_codes VARCHAR[] )
  RETURNS TABLE (
  o_exchnage_title VARCHAR,
  o_target_code varchar,
  o_reference_code varchar,
  o_time_stamp TIMESTAMP ,
  o_rate REAL
  )
AS $$
    BEGIN

      RETURN QUERY SELECT  DISTINCT on (reference_code) exchange_title as o_exchnage_title, target_code as o_target_code, reference_code as o_reference_code, time_stamp as o_time_stamp, rate as o_rate
FROM rates_view
WHERE time_stamp <= p_time_stamp
and exchange_title = p_exchange_title
and target_code = p_target_code
and reference_code = ANY(p_referencies_codes)
ORDER by reference_code, time_stamp DESC;
    END;
    $$
LANGUAGE plpgsql;


CREATE or REPLACE VIEW rates_view as
SELECT ex.title as exchange_title, tcr.code as target_code, rcr.code as reference_code, r.time_stamp, r.rate
FROM rates r, currencies tcr, currencies rcr, exchanges ex
WHERE  r.exchange_id = ex.id AND r.target_id = tcr.id and r.reference_id = rcr.id;



create table userinfo
(
	uid serial not null
		constraint userinfo_pkey
			primary key,
	username varchar(100) not null,
	departname varchar(500) not null,
	created date
)
;

create table exchanges
(
	id serial not null
		constraint exchanges_pkey
			primary key,
	title varchar not null,
	create_date timestamp not null
)
;

create unique index exchanges_title_uindex
	on exchanges (title)
;

create table currencies
(
	id serial not null
		constraint currencies_pkey
			primary key,
	code varchar not null,
	title varchar not null,
	create_date timestamp not null,
	native_id integer not null
)
;

create unique index currencies_code_uindex
	on currencies (code)
;

create table rates
(
	exchange_id integer not null,
	target_id integer not null,
	reference_id integer not null,
	time_stamp timestamp not null,
	rate real not null
)
;

create table exchanges_pairs
(
	exchange_id integer not null,
	target_id integer not null,
	reference_id integer not null,
	time_stamp timestamp not null
)
;

create table sa_rates
(
	exchange_title varchar not null,
	target_title varchar not null,
	target_code varchar not null,
	target_native_id integer not null,
	reference_title varchar not null,
	reference_code varchar not null,
	reference_native_id integer not null,
	time_stamp timestamp not null,
	rate real not null
)
;


SELECT CURRENT_TIMESTAMP - INTERVAL '5 minutes'


SELECT  DISTINCT on (reference_code) exchange_title, target_code, reference_code, time_stamp, rate
FROM (SELECT ex.title as exchange_title, tcr.code as target_code, rcr.code as reference_code, r.time_stamp, r.rate
FROM rates r, currencies tcr, currencies rcr, exchanges ex
WHERE  r.exchange_id = ex.id AND r.target_id = tcr.id and r.reference_id = rcr.id) r
WHERE time_stamp <= (CURRENT_TIMESTAMP  - INTERVAL '5 minutes' + INTERVAL '3 hours')
and exchange_title = 'BINANCE'
and target_code = 'ETH'
and reference_code in ('USDT', 'BTC')
ORDER by reference_code, time_stamp DESC;


SELECT COUNT(*)
from (
SELECT ex.title as exchange_title, tcr.code as target_code, rcr.code as reference_code, r.time_stamp, r.rate
FROM rates r, currencies tcr, currencies rcr, exchanges ex
WHERE  r.exchange_id = ex.id AND r.target_id = tcr.id and r.reference_id = rcr.id) a;


SELECT COUNT (*)
FROM rates;




SELECT * from getRates(CURRENT_TIMESTAMP, 'BINANCE', 'ETH', '{USDT, BTC}');


SELECT DISTINCT on (ex.id, tcr.id, rcr.id ) ex.id exchange_id, tcr.id target_id, rcr.id reference_id, time_stamp
FROM sa_rates sr, exchanges ex, currencies tcr,  currencies rcr
WHERE sr.exchange_title = ex.title and sr.target_code = tcr.code and sr.reference_code = rcr.code;