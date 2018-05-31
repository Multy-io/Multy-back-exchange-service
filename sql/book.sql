create or replace function fill_book() returns void
	language plpgsql
as $$
BEGIN

INSERT INTO exchanges(title, create_date)
SELECT DISTINCT exchange_title, CURRENT_TIMESTAMP
FROM sa_book_orders
ON CONFLICT DO NOTHING;

INSERT INTO currencies(title, code, create_date, native_id)
SELECT DISTINCT target_title, target_code,  CURRENT_TIMESTAMP , target_native_id
  from sa_book_orders
UNION
  SELECT DISTINCT reference_title, reference_code,  CURRENT_TIMESTAMP, reference_native_id
  from sa_book_orders
ON CONFLICT DO NOTHING;

  INSERT INTO exchanges_pairs
SELECT DISTINCT on (ex.id, tcr.id, rcr.id) ex.id exchange_id, tcr.id target_id, rcr.id reference_id, time_stamp, false as is_calculated
FROM sa_book_orders sr, exchanges ex, currencies tcr,  currencies rcr
WHERE sr.exchange_title = ex.title and sr.target_code = tcr.code and sr.reference_code = rcr.code
 ON CONFLICT DO NOTHING;

  TRUNCATE book_orders;

  INSERT into book_orders
SELECT ex.id exchange_id, tcr.id target_id, rcr.id reference_id, sr.time_stamp, sr.price, sr.amount, sr.is_ask
from sa_book_orders sr, exchanges ex, currencies tcr,  currencies rcr
WHERE sr.exchange_title = ex.title and sr.target_code = tcr.code and sr.reference_code = rcr.code;


  TRUNCATE sa_book_orders;
END;
$$
;


create or REPLACE view book_view as
SELECT ex.title, tcr.code target, rcr.code reference, bo.price price, bo.amount, bo.is_ask
FROM book_orders bo, exchanges ex,  currencies tcr,  currencies rcr
WHERE bo.exchange_id = ex.id and bo.target_id = tcr.id and bo.reference_id = rcr.id
ORDER BY ex.title, is_ask, price;


SELECT title, sum(amount), is_ask
FROM book_view
GROUP BY title, is_ask;

SELECT title, round(price,0) price, sum(amount), is_ask
FROM book_view
WHERE title = 'BITFINEX'
GROUP BY title, is_ask, round(price,0)
ORDER BY is_ask  , price;



SELECT ex.title, tcr.code, rcr.code, round(bo.price,0) price, sum(bo.amount), bo.is_ask
FROM book_orders bo, exchanges ex,  currencies tcr,  currencies rcr
WHERE bo.exchange_id = ex.id and bo.target_id = tcr.id and bo.reference_id = rcr.id
  GROUP BY ex.title, tcr.code, rcr.code, round(bo.price,0), bo.is_ask
ORDER BY ex.title, is_ask, price DESC ;