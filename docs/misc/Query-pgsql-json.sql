
-- #START :: v.161006

SELECT id, data#>>'{extra,os,0}' FROM public.test_tbl ORDER BY data->>'pad_urls_id' ASC

SELECT * FROM distrib.apps WHERE json_data#>>'{pad_data,downloads}' ILIKE '%play.google%';

SELECT * FROM distrib.apps WHERE json_data#>'{extra,os}' ? 'mac'
SELECT * FROM distrib.apps WHERE json_data#>>'{pad_data,downloads}' ILIKE '%play.google%'

CREATE INDEX apps__idx__json_data_progr_name ON distrib.apps USING btree (((json_data #>> '{pad_data,program,name}')));
CREATE INDEX apps__idx__json_data_progr_categ ON distrib.apps USING btree (((json_data #>> '{pad_data,program,category}')));
CREATE INDEX apps__idx__json_data_progr_scateg ON distrib.apps USING btree (((json_data #>> '{pad_data,program,sub_category}')));
CREATE INDEX apps__idx__json_data_progr_os ON website.apps USING btree (((json_data #> '{extra,os}')));

CREATE INDEX apps__idx__json_data_progr_os ON website.apps USING btree (COALESCE((json_data #> '{extra,os}'::jsonb),'[]'::jsonb));

CREATE INDEX apps__idx__json_data_progr_cost ON website.apps USING btree (CAST(('0' || COALESCE((json_data #>> '{pad_data,cost,dollars}'),'0')) AS NUMERIC));
CREATE INDEX apps__idx__json_data_progr_fsize ON website.apps USING btree (CAST(('0' || COALESCE((json_data #>> '{pad_data,program,file_size}'),'0')) AS BIGINT));

CREATE INDEX apps__idx__json_data_progr_review_rating ON website.apps USING btree (CAST(('0' || COALESCE(("review_json"#>>'{rating}'),'0')) AS NUMERIC));
CREATE INDEX apps__idx__json_data_progr_review_pdate ON website.apps USING btree (((review_json ->> 'publish_date')));

SELECT * FROM distrib.apps WHERE json_data#>'{extra,os}' ? 'linux'

UPDATE "apps" SET "scateg_id" = (json_data #>> '{pad_data,program,sub_category}')
UPDATE "apps" SET "scateg_id" = ( SELECT "categs"."id" FROM "categs" WHERE ("categs"."name" = "apps"."scateg_id" AND "categs"."parent_id" = "apps"."categ_id") LIMIT 1 OFFSET 0)


SELECT slug, CAST(('0' || COALESCE((json_data #>> '{pad_data,program,file_size}'),'0')) AS BIGINT) as fsize FROM apps WHERE (json_data #>> '{pad_data,program,file_size}') LIKE '%.%'
ORDER BY CAST(('0' || COALESCE((json_data #>> '{pad_data,program,file_size}'),'0')) AS BIGINT) ASC LIMIT 10 OFFSET 0

UPDATE "apps" SET os_win = TRUE WHERE ("json_data"#>'{extra,os}' ? 'windows')
UPDATE "apps" SET os_mac = TRUE WHERE ("json_data"#>'{extra,os}' ? 'mac')

UPDATE "apps" SET os_lin = TRUE WHERE ("json_data"#>'{extra,os}' ? 'linux');
UPDATE "apps" SET os_and = TRUE WHERE ("json_data"#>'{extra,os}' ? 'android');
UPDATE "apps" SET os_ios = TRUE WHERE ("json_data"#>'{extra,os}' ? 'ios');

--

CREATE TABLE test_url (url varchar NOT NULL, url_tsvector tsvector NOT NULL);
CREATE OR REPLACE FUNCTION generate_url_tsvector(varchar)
RETURNS tsvector
LANGUAGE sql
AS $_$
    SELECT to_tsvector(regexp_replace($1, '[^\w]+', ' ', 'gi'));
$_$;

CREATE OR REPLACE FUNCTION before_insert_test_url()
RETURNS TRIGGER
LANGUAGE plpgsql AS $_$
BEGIN
  NEW.url_tsvector := generate_url_tsvector(NEW.url);
  RETURN NEW;
END
$_$;

CREATE TRIGGER before_insert_test_url_trig
BEFORE INSERT ON test_url
FOR EACH ROW EXECUTE PROCEDURE before_insert_test_url();

--

SELECT DISTINCT jsonb_array_elements_text(topics)::text FROM quotes;

SELECT (jsonb_array_elements(authors)::jsonb #> '{id}')::jsonb FROM quotes;

SELECT (jsonb_array_elements(authors)::jsonb #>> '{id}')::text FROM quotes;
SELECT * FROM ( SELECT *, (jsonb_array_elements(authors)::jsonb #>> '{id}')::text AS author_id FROM quotes ) a WHERE a.author_id = '0000000002';

SELECT * FROM quotes WHERE topics ? 'T000000001';

SELECT * FROM quotes WHERE topics ?| ARRAY(SELECT id FROM topics WHERE slug = 'topic-one'); -- inneficient !!!
SELECT a.id, b.* FROM topics a INNER JOIN quotes b ON (b.topics ? a.id) LIMIT 25 OFFSET 0 -- a better solution

SELECT a.dom, b.pad_urls_id, b.source FROM blacklist_domains a INNER JOIN radix_urls b ON (b.dom_data ? a.dom) LIMIT 25 OFFSET 0


SELECT to_tsvector('english', 'a fat  cat sat on a mat - it ate a fat rats');
SELECT to_tsvector('a fat  cat sat on a mat - it ate a fat rats');
UPDATE "quotes" SET "fts_quote" = to_tsvector("quote")
UPDATE "quotes" SET "fts_quote" = to_tsvector(smart_deaccent_string("quote"))
SELECT * FROM quotes WHERE fts_quote @@ to_tsquery('days')
SELECT * FROM quotes WHERE fts_quote @@ to_tsquery(smart_deaccent_string('days & quote'))

INSERT INTO quotes VALUES ('ABCDEF1230', 'some-quote', '0000000000', '["0000000001", "0000000002"]', '["topic1", "topic2"]', 'Some Quote of the day
by Räksmörgås Jösefsson', '["some", "quote", "day"]', '''day'':5 ''josefsson'':8 ''quot'':2 ''raksmorga'':7', '[]');

--

SELECT COALESCE(token, '') as host FROM ts_debug('http://example.com/stuff/index.html') WHERE alias = 'host';
SELECT ('["' || COALESCE(token, '') || '"]')::jsonb AS host FROM ts_debug('http://www.example.com/stuff/index.html'

UPDATE radix_urls a SET dom_data = ( SELECT ('["' || COALESCE(token, '') || '"]')::jsonb FROM ts_debug(a.value) WHERE alias = 'host' LIMIT 1 OFFSET 0 );


-- select substring( 'http://www.arandomsite.co.uk' from '^[^:]*://(?:[^/:]*:[^/@]*@)?(?:[^/:.]*\.)+([^:/]+)' ) as tld;

-- select substring( 'http://www.arandomsite.com' from '^[^:]*://(?:[^/:]*:[^/@]*@)?(?:[^/:.]*\.)+((?:[^/:.]*\.)+)' ) as tld;

SELECT * FROM tmp.test_url WHERE url_data ? 'www.google.com'
-- SELECT COALESCE(token, '') as host FROM ts_debug('http://example.com/stuff/index.html') WHERE alias = 'host';

UPDATE radix_urls a SET dom_data = COALESCE(( SELECT ('["' || COALESCE(token, '') || '"]')::jsonb FROM ts_debug(a.value) WHERE alias = 'host' LIMIT 1 OFFSET 0 ), '[]');


UPDATE blacklist_domains a SET dom_info = COALESCE(( SELECT radix_info FROM radix_subdomains b WHERE a.dom = b.radix LIMIT 1 OFFSET 0 ), '')

-- View Sizes
SELECT relname, relpages FROM pg_class ORDER BY relpages DESC;

-- #END
