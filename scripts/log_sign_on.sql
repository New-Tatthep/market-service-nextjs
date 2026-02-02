CREATE OR REPLACE PROCEDURE initial_log_sign_on()
LANGUAGE plpgsql
AS $$
--     DECLARE
--         a integer := 10;
--         b integer := 20;
    BEGIN
        IF EXISTS (SELECT FROM pg_catalog.pg_tables WHERE schemaname = 'public' AND tablename  = 'log_sign_on') THEN
            RAISE NOTICE 'Table public.log_sign_on already exists.';
        ELSE

            CREATE TABLE log_sign_on (
                last_update int8 NOT NULL,
                username varchar(50) NOT NULL,
                "action" varchar(10) NOT NULL,
                message text NULL,
                ip_address varchar(50) NULL,
                machine_id varchar(255) NULL,
                source_action varchar(20) NOT NULL,
                log_code varchar(100) NOT NULL,
                status varchar(10) NOT NULL,
                login_type varchar(10) NULL,
                CONSTRAINT pk_employee_sign_on_log PRIMARY KEY (log_code)
            );


        END IF;

       
    END
$$;

CALL initial_log_sign_on();

DROP PROCEDURE IF EXISTS initial_log_sign_on();

ALTER TABLE log_sign_on ADD COLUMN IF NOT EXISTS user_agent varchar(255) NOT NULL DEFAULT '';
