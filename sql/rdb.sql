




-- CREATE UNIQUE INDEX exchanges_id_uindex ON public.exchanges (id);

CREATE OR REPLACE FUNCTION fill_rates()
    RETURNS void AS $$
    BEGIN

INSERT INTO exchanges(title, create_date)
SELECT DISTINCT exchange_title,  current_timestamp
FROM sa_rates
ON CONFLICT DO NOTHING;

      INSERT INTO currencies(title, code, create_date, native_id)
SELECT DISTINCT target_title, target_code,  current_timestamp, target_native_id
  from sa_rates
UNION
  SELECT DISTINCT reference_title, reference_code,  current_timestamp, reference_native_id
  from sa_rates
ON CONFLICT DO NOTHING;

      INSERT INTO rates
      SELECT ex.id exchange_id, tcr.id target_id, rcr.id reference_id, sr.time_stamp, sr.rate
FROM sa_rates sr, exchanges ex, currencies tcr,  currencies rcr
WHERE sr.exchange_title = ex.title and sr.target_code = tcr.code and sr.reference_code = rcr.code;
      TRUNCATE sa_rates;
    END;
    $$ LANGUAGE plpgsql;



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

create function fill_rates() returns void
	language plpgsql
as $$
BEGIN

INSERT INTO exchanges(title, create_date)
SELECT DISTINCT exchange_title,  current_timestamp
FROM sa_rates
ON CONFLICT DO NOTHING;

      INSERT INTO currencies(title, code, create_date, native_id)
SELECT DISTINCT target_title, target_code,  current_timestamp, target_native_id
  from sa_rates
UNION
  SELECT DISTINCT reference_title, reference_code,  current_timestamp, reference_native_id
  from sa_rates
ON CONFLICT DO NOTHING;

      INSERT INTO rates
      SELECT ex.id exchange_id, tcr.id target_id, rcr.id reference_id, sr.time_stamp, sr.rate
FROM sa_rates sr, exchanges ex, currencies tcr,  currencies rcr
WHERE sr.exchange_title = ex.title and sr.target_code = tcr.code and sr.reference_code = rcr.code;
      TRUNCATE sa_rates;
    END;
$$
;

