ALTER TABLE photos MODIFY id BIGINT unsigned;
ALTER TABLE photos_keywords MODIFY photo_id BIGINT unsigned;
ALTER TABLE details MODIFY photo_id BIGINT unsigned;
ALTER TABLE photos_labels MODIFY photo_id BIGINT unsigned;
ALTER TABLE files MODIFY id BIGINT unsigned;
ALTER TABLE files_share MODIFY file_id BIGINT unsigned;