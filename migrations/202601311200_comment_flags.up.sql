-- Comment flags: any user can flag a comment for inappropriateness; flags visible only to moderators
CREATE TABLE IF NOT EXISTS comment_flags (
    id          SERIAL PRIMARY KEY,
    tenant_id   INT NOT NULL,
    comment_id  INT NOT NULL,
    user_id     INT NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    reason      TEXT,
    CONSTRAINT comment_flags_tenant_id_fkey FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    CONSTRAINT comment_flags_comment_id_fkey FOREIGN KEY (comment_id) REFERENCES comments(id),
    CONSTRAINT comment_flags_user_fkey FOREIGN KEY (tenant_id, user_id) REFERENCES users(tenant_id, id),
    CONSTRAINT comment_flags_unique UNIQUE (tenant_id, comment_id, user_id)
);

CREATE INDEX IF NOT EXISTS comment_flags_tenant_comment ON comment_flags (tenant_id, comment_id);
