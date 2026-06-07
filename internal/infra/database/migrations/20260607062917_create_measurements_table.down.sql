DROP INDEX IF EXISTS measurements_device_room_index;
DROP INDEX IF EXISTS measurements_deleted_date_index;
DROP INDEX IF EXISTS measurements_created_date_index;
DROP INDEX IF EXISTS measurements_room_id_index;
DROP INDEX IF EXISTS measurements_device_id_index;

DROP TABLE IF EXISTS public.measurements;