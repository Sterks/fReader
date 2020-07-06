ALTER TABLE public."SourceRegions" ADD CONSTRAINT sourceregions_fk FOREIGN KEY (r_fz_law) REFERENCES "FederalLaw"(fl_id);

