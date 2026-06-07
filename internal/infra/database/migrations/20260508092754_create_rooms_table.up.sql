CREATE TABLE IF NOT EXISTS public.rooms
(
    id              serial PRIMARY KEY,
    organization_id         integer NOT NULL REFERENCES public.users(id),
    name            varchar(250) NOT NULL,
    description     text,
    city            varchar(100) NOT NULL,
    address         varchar(250) NOT NULL,
    lat             float4 NOT NULL,
    lon             float4 NOT NULL,
    created_date timestamptz NOT NULL,
    updated_date timestamptz NOT NULL,
    deleted_date timestamptz 
);
