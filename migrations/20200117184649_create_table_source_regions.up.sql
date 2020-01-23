CREATE TABLE public."SourceRegions" (
	r_id serial NOT NULL,
	r_name varchar NULL,
	r_date_create timestamp NOT NULL,
	r_date_update timestamp NOT NULL
);
ALTER TABLE public."SourceRegions" ADD CONSTRAINT sourceregions_pk PRIMARY KEY (r_id);