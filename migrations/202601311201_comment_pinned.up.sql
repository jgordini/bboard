-- Pinned comments: moderators can pin/unpin; pinned comments appear first
ALTER TABLE comments ADD COLUMN IF NOT EXISTS pinned_at TIMESTAMPTZ NULL;
ALTER TABLE comments ADD COLUMN IF NOT EXISTS pinned_by_id INT NULL;

ALTER TABLE comments
   ADD CONSTRAINT comments_pinned_by_id_fkey
   FOREIGN KEY (pinned_by_id, tenant_id) REFERENCES users(id, tenant_id);
