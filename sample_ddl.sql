-- Table: gpstrace.gpstrace

-- DROP TABLE gpstrace.gpstrace;

CREATE TABLE gpstrace.gpstrace
(
    event_id character varying(16) COLLATE pg_catalog."default",
    asset_id character varying(16) COLLATE pg_catalog."default",
    gps_data jsonb
) PARTITION BY LIST (event_id) 
WITH (
    OIDS = FALSE
)
TABLESPACE pg_default;

ALTER TABLE gpstrace.gpstrace
    OWNER to gpswriter;

-- Partitions SQL

CREATE TABLE gpstrace."DEMO1" PARTITION OF gpstrace.gpstrace
    FOR VALUES IN ('DEMO1');

CREATE TABLE gpstrace."DEMO2" PARTITION OF gpstrace.gpstrace
    FOR VALUES IN ('DEMO2');
