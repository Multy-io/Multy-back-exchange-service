create table if not exists exchanges
(
	id serial not null
		constraint exchanges_pkey
			primary key,
	title varchar not null,
	create_date timestamp not null
)
;

create unique index if not exists exchanges_title_uindex
	on exchanges (title)
;

create table if not exists currencies
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

create unique index if not exists currencies_code_uindex
	on currencies (code)
;


create table if not exists rates
(
	exchange_id integer not null,
	target_id integer not null,
	reference_id integer not null,
	time_stamp timestamp not null,
	rate real not null,
	is_calculated boolean not null
);


create table if not exists exchanges_pairs
(
	exchange_id integer not null,
	target_id integer not null,
	reference_id integer not null,
	time_stamp timestamp not null,
	is_calculated boolean not null,
	constraint exchanges_pairs_exchange_id_reference_id_target_id_is_calculate
		unique (exchange_id, reference_id, target_id, is_calculated)
);

-- SA Source Abstraction - Raw data
create table if not exists sa_rates
(
	exchange_title varchar not null,
	target_title varchar not null,
	target_code varchar not null,
	target_native_id integer not null,
	reference_title varchar not null,
	reference_code varchar not null,
	reference_native_id integer not null,
	time_stamp timestamp not null,
	rate real not null,
	is_calculated boolean not null
)
;

create view rates_view as
SELECT ex.title AS exchange_title,
    tcr.code AS target_code,
    rcr.code AS reference_code,
    r.time_stamp,
    r.rate,
    r.is_calculated
   FROM rates r,
    currencies tcr,
    currencies rcr,
    exchanges ex
  WHERE ((r.exchange_id = ex.id) AND (r.target_id = tcr.id) AND (r.reference_id = rcr.id));


create or replace function fill_rates() returns void
	language plpgsql
as $$
BEGIN

INSERT INTO exchanges(title, create_date)
SELECT DISTINCT exchange_title, CURRENT_TIMESTAMP
FROM sa_rates
ON CONFLICT DO NOTHING;

INSERT INTO currencies(title, code, create_date, native_id)
SELECT DISTINCT target_title, target_code,  CURRENT_TIMESTAMP , target_native_id
  from sa_rates
UNION
  SELECT DISTINCT reference_title, reference_code,  CURRENT_TIMESTAMP, reference_native_id
  from sa_rates
ON CONFLICT DO NOTHING;

INSERT INTO exchanges_pairs
SELECT DISTINCT on (ex.id, tcr.id, rcr.id, sr.is_calculated) ex.id exchange_id, tcr.id target_id, rcr.id reference_id, time_stamp, sr.is_calculated as is_calculated
FROM sa_rates sr, exchanges ex, currencies tcr,  currencies rcr
WHERE sr.exchange_title = ex.title and sr.target_code = tcr.code and sr.reference_code = rcr.code
 ON CONFLICT DO NOTHING;


INSERT INTO rates
SELECT ex.id exchange_id, tcr.id target_id, rcr.id reference_id, sr.time_stamp, sr.rate, sr.is_calculated as is_calculated
FROM sa_rates sr, exchanges ex, currencies tcr,  currencies rcr
WHERE sr.exchange_title = ex.title and sr.target_code = tcr.code and sr.reference_code = rcr.code;

     TRUNCATE sa_rates;
    END;
$$
;

-- Get exchange info at given point in time for given exchange from one currency to set of others
create or replace function getrates(p_time_stamp timestamp with time zone, p_exchange_title character varying, p_target_code character varying, p_referencies_codes character varying[])
	returns TABLE(o_exchnage_title character varying, o_target_code character varying, o_reference_code character varying, o_time_stamp timestamp without time zone, o_rate real)
	language plpgsql
as $$
BEGIN

      RETURN QUERY SELECT  DISTINCT on (reference_code) exchange_title as o_exchnage_title, target_code as o_target_code, reference_code as o_reference_code, time_stamp as o_time_stamp, rate as o_rate
FROM rates_view
WHERE time_stamp <= p_time_stamp -- add some kind of a window (a day?) limiting number of rows matching by time.
	and exchange_title = p_exchange_title
	and target_code = p_target_code
	and reference_code = ANY(p_referencies_codes)
ORDER by reference_code, time_stamp DESC;
    END;
$$
;

