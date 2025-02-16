ALTER TABLE inventory 
ADD CONSTRAINT inventory_unique UNIQUE (user_id, merch_id);
