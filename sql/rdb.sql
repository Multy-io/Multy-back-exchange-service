




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
SELECT DISTINCT on (ex.id, tcr.id, rcr.id ) ex.id exchange_id, tcr.id target_id, rcr.id reference_id, time_stamp, FALSE as is_calculated
FROM sa_rates sr, exchanges ex, currencies tcr,  currencies rcr
WHERE sr.exchange_title = ex.title and sr.target_code = tcr.code and sr.reference_code = rcr.code
 ON CONFLICT DO NOTHING;


INSERT INTO rates
SELECT ex.id exchange_id, tcr.id target_id, rcr.id reference_id, sr.time_stamp, sr.rate, false as is_calculated
FROM sa_rates sr, exchanges ex, currencies tcr,  currencies rcr
WHERE sr.exchange_title = ex.title and sr.target_code = tcr.code and sr.reference_code = rcr.code
UNION
SELECT ex.id exchange_id, tcr.id target_id, rcr.id reference_id, sr.time_stamp, sr.rates, true as is_calculated
FROM sa_cross_rates sr, exchanges ex, currencies tcr,  currencies rcr
WHERE sr.exchange_title = ex.title and sr.tc_code = tcr.code and sr.rc_cross_code = rcr.code;

INSERT INTO exchanges_pairs
SELECT DISTINCT on (ex.id, tcr.id, rcr.id ) ex.id exchange_id, tcr.id target_id, rcr.id reference_id, time_stamp,is_calculated
FROM sa_cross_rates sr, exchanges ex, currencies tcr,  currencies rcr
WHERE sr.exchange_title = ex.title and sr.tc_code = tcr.code and sr.rc_cross_code = rcr.code
 ON CONFLICT DO NOTHING;

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


SELECT exchange_id, target_id, *
FROM exchanges_pairs ex, currencies rcr
WHERE ex.reference_id = rcr.id and code in ('USDT','USD');




SELECT e.title, tc.code, rc.code
FROM exchanges_pairs ep, exchanges e, currencies tc, currencies rc
WHERE ep.exchange_id = e.id and ep.target_id = tc.id and ep.reference_id = rc.id
and e.title = 'HITBTC'
ORDER  by  1




SELECT * from (
SELECT DISTINCT target_id
FROM exchanges_pairs) t,
(SELECT DISTINCT reference_id
FROM exchanges_pairs) r,
(SELECT DISTINCT exchanges.id
FROM exchanges) ex;


TRUNCATE TABLE public.currencies
    CONTINUE IDENTITY
    RESTRICT;

TRUNCATE TABLE exchanges
    CONTINUE IDENTITY
    RESTRICT;

TRUNCATE TABLE exchanges_pairs
    CONTINUE IDENTITY
    RESTRICT;


TRUNCATE TABLE rates
    CONTINUE IDENTITY
    RESTRICT;


TRUNCATE TABLE sa_rates
    CONTINUE IDENTITY
    RESTRICT;

CREATE or replace view sa_cross_pairs as
SELECT ex.id ex_id, tc.id tc_id, rc.id rc_id, ex.title, tc.code as tc_code, rc.code as rc_code from (
SELECT a.id as exchange_id, a.target_id, a.reference_id, exp.exchange_id as exist
FROM exchanges_pairs exp
RIGHT OUTER JOIN
(
SELECT * from (
SELECT DISTINCT target_id
FROM exchanges_pairs
WHERE is_calculated = FALSE) t,
(SELECT DISTINCT reference_id
FROM exchanges_pairs
WHERE is_calculated = FALSE) r,
(SELECT DISTINCT exchanges.id
FROM exchanges) ex
WHERE t.target_id != r.reference_id
) a
on exp.reference_id = a.reference_id and exp.target_id = a.target_id and exp.exchange_id = a.id and exp.is_calculated = FALSE ) result, exchanges ex, currencies tc, currencies rc, exchanges_pairs ep
WHERE result.exist is NULL and result.exchange_id = ex.id and result.target_id = tc.id and result.reference_id = rc.id and result.exchange_id = ep.exchange_id and result.target_id = ep.target_id and ep.is_calculated = FALSE
ORDER  by ex.title, 2;


SELECT *
from sa_cross_pairs;


CREATE or replace view sa_cross_curriences as
SELECT c.title exchange_title, c.tc_code, cr.cross_code rf_code,  c.rc_code as rc_cross_code
FROM
( SELECT *
from sa_cross_pairs) c ,
( SELECT DISTINCT exchange_id cross_exchange_id, reference_id cross_reference_id, is_calculated, ex.title cross_title, c.code cross_code
from exchanges_pairs ep, exchanges ex, currencies c
WHERE ep.exchange_id = ex.id and reference_id = c.id)  cr
WHERE c.ex_id = cr.cross_exchange_id and c.rc_id != cr.cross_reference_id;


SELECT *
from sa_cross_curriences;

SELECT cr.exchange_title, cr.tc_code, cr.rf_code, sar.rate, cr.rc_cross_code--, sar_cross.rate
from sa_cross_curr cr, sa_rates sar--, sa_rates sar_cross
WHERE cr.exchange_title = sar.exchange_title and cr.tc_code = sar.target_code and cr.rf_code = sar.reference_code


CREATE or replace view sa_cross_rates as;


SELECT cr.exchange_title, cr.tc_code, cr.rc_cross_code,  sar.rate/sar_cross.rate as rates, true as is_calculated, sar.time_stamp
from sa_cross_curriences cr, sa_rates sar, sa_rates sar_cross
WHERE cr.exchange_title = sar.exchange_title and cr.tc_code = sar.target_code and cr.rf_code = sar.reference_code
and cr.exchange_title = sar_cross.exchange_title and cr.rc_cross_code = sar_cross.target_code and cr.rf_code = sar_cross.reference_code
UNION all
SELECT cr.exchange_title, cr.tc_code, cr.rc_cross_code, sar.rate *  sar_cross.rate rates, true as is_calculated, sar.time_stamp
from sa_cross_curriences cr, sa_rates sar, sa_rates sar_cross
WHERE cr.exchange_title = sar.exchange_title and cr.tc_code = sar.target_code and cr.rf_code = sar.reference_code
and cr.exchange_title = sar_cross.exchange_title and  cr.rc_cross_code = sar_cross.reference_code and cr.rf_code = sar_cross.target_code;

SELECT *
FROM  sa_cross_rates;



 SELECT  DISTINCT on (reference_code) exchange_title as o_exchnage_title, target_code as o_target_code, reference_code as o_reference_code, time_stamp as o_time_stamp, rate as o_rate
FROM rates_view
WHERE exchange_title = p_exchange_title
and target_code = p_target_code
and reference_code = 'STEEM'
ORDER by reference_code, time_stamp DESC;