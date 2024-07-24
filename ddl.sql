CREATE SCHEMA careers AUTHORIZATION postgres;

CREATE TABLE careers.applications (
	id bigserial NOT NULL,
	"name" text NULL,
	email text NULL,
	phone_no int8 NULL,
	created_at timestamp NULL,
	resume text NULL,
	linkedin_url text DEFAULT ''::text NULL,
	portfolio_url text DEFAULT ''::text NULL
);

CREATE TABLE public.contact_vulnerability_enquiry_form (
	id int8 DEFAULT nextval('form_data_id_seq'::regclass) NOT NULL,
	"name" text DEFAULT ''::text NULL,
	email text DEFAULT ''::text NULL,
	mobile text DEFAULT ''::text NULL,
	created_at text NULL,
	company text DEFAULT ''::text NULL,
	message text DEFAULT ''::text NULL,
	url text DEFAULT ''::text NULL,
	form_type text DEFAULT ''::text NULL
);