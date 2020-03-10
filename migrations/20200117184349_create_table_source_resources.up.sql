CREATE TABLE public."SourceResources" (
	sr_id serial NOT NULL,
	sr_name varchar NULL,
	sr_fullname varchar NULL
);

ALTER TABLE public."SourceResources" ADD CONSTRAINT sourceresources_pk PRIMARY KEY (sr_id);
