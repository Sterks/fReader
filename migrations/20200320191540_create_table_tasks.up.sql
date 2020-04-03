CREATE TABLE public."Tasks" (
	ts_id serial NOT NULL,
	ts_name varchar NULL,
	ts_data_start timestamp NULL,
	ts_run_times int4 NULL,
	ts_comment text NULL,
	CONSTRAINT tasks_pk PRIMARY KEY (ts_id)
);