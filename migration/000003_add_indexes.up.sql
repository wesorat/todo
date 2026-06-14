CREATE INDEX idx_lists_user_id ON lists(user_id);
CREATE INDEX idx_items_list_id ON items(list_id);
CREATE INDEX idx_refresh_user_id ON refresh(user_id);