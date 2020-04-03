ALTER TABLE public."Files" ADD CONSTRAINT files_fk1 FOREIGN KEY (f_source_resources_id) REFERENCES public."SourceResources"(sr_id) ON DELETE CASCADE ON UPDATE CASCADE;
