CREATE TABLE public."SourceRegions" (
	r_id serial NOT NULL,
	r_name varchar NULL,
	r_date_create timestamp NOT NULL,
	r_date_update timestamp NOT NULL
);
ALTER TABLE public."SourceRegions" ADD CONSTRAINT sourceregions_pk UNIQUE (r_id);
ALTER TABLE public."Files" ADD CONSTRAINT files_fr FOREIGN KEY (f_area) REFERENCES public."SourceRegions"(r_id) ON DELETE CASCADE ON UPDATE CASCADE;

INSERT INTO public."SourceResources" (sr_id, sr_name, sr_fullname) VALUES(1, 'notifications44', 'Извещения по 44');
INSERT INTO public."SourceResources" (sr_id, sr_name, sr_fullname) VALUES(2, 'protocols44', 'Протоколы по 44');
INSERT INTO public."SourceResources" (sr_id, sr_name, sr_fullname) VALUES(3, 'notifications223', 'Извещения по 223');
INSERT INTO public."SourceResources" (sr_id, sr_name, sr_fullname) VALUES(4, 'protocols223', 'Протоколы по 223');



