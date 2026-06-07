CREATE TABLE IF NOT EXISTS public.devices
(
    id                  BIGSERIAL PRIMARY KEY,

    organization_id      BIGINT NOT NULL REFERENCES public.organizations(id),
    room_id              BIGINT REFERENCES public.rooms(id),

    guid                VARCHAR(255) NOT NULL,
    inventory_number    VARCHAR(255) NOT NULL,
    serial_number       VARCHAR(255) NOT NULL,
    characteristics    TEXT,

    category            VARCHAR(50) NOT NULL,
    units               VARCHAR(100),

    power_consumption   DOUBLE PRECISION NOT NULL DEFAULT 0,

    created_date        TIMESTAMPTZ NOT NULL,
    updated_date        TIMESTAMPTZ NOT NULL,
    deleted_date        TIMESTAMPTZ,

    CONSTRAINT devices_category_check
        CHECK (category IN ('SENSOR', 'ACTUATOR'))
);

CREATE UNIQUE INDEX IF NOT EXISTS devices_guid_uindex
    ON public.devices (guid)
    WHERE deleted_date IS NULL;

CREATE INDEX IF NOT EXISTS devices_organization_id_index
    ON public.devices (organization_id);

CREATE INDEX IF NOT EXISTS devices_room_id_index
    ON public.devices (room_id);

CREATE INDEX IF NOT EXISTS devices_category_index
    ON public.devices (category);

CREATE INDEX IF NOT EXISTS devices_deleted_date_index
    ON public.devices (deleted_date);

