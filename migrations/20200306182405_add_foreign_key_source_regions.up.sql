ALTER TABLE public."Files" ADD CONSTRAINT files_fk1 FOREIGN KEY (f_source_regions_id) REFERENCES public."SourceResources"(sr_id) ON DELETE CASCADE ON UPDATE CASCADE;
