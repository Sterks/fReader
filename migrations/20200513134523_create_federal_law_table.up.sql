CREATE TABLE public."FederalLaw" (
	fl_id serial NOT NULL,
	fl_name_law varchar NULL,
	fl_comment varchar NULL
);
ALTER TABLE public."FederalLaw" ADD CONSTRAINT federallaw_pk PRIMARY KEY (fl_id);
ALTER TABLE public."FederalLaw" ADD CONSTRAINT federallaw_un UNIQUE (fl_id);


