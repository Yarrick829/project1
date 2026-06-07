ALTER TABLE public.rooms
    DROP CONSTRAINT IF EXISTS rooms_organization_id_fkey;

ALTER TABLE public.rooms
    ADD CONSTRAINT rooms_organization_id_fkey
    FOREIGN KEY (organization_id)
    REFERENCES public.organizations(id);