CREATE TABLE IF NOT EXISTS public.events
(
id BIGSERIAL PRIMARY KEY,

device_id       BIGINT NOT NULL
                REFERENCES public.devices(id),

room_id         BIGINT NOT NULL
                REFERENCES public.rooms(id),

action          VARCHAR(10) NOT NULL,

created_date    TIMESTAMPTZ NOT NULL,
updated_date    TIMESTAMPTZ NOT NULL,
deleted_date    TIMESTAMPTZ

);

CREATE INDEX IF NOT EXISTS events_device_id_index
ON public.events (device_id);

CREATE INDEX IF NOT EXISTS events_room_id_index
ON public.events (room_id);

CREATE INDEX IF NOT EXISTS events_created_date_index
ON public.events (created_date);

CREATE INDEX IF NOT EXISTS events_deleted_date_index
ON public.events (deleted_date);

CREATE INDEX IF NOT EXISTS events_device_room_index
ON public.events (
device_id,
room_id,
created_date
);

CREATE INDEX IF NOT EXISTS events_device_action_index
ON public.events (
device_id,
action,
created_date
);