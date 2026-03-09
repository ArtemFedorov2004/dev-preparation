CREATE TYPE progress_status AS ENUM ('learned', 'need_review', 'dont_know');

CREATE TABLE IF NOT EXISTS user_progress (
                                             user_id     TEXT        NOT NULL,
                                             question_id BIGINT      NOT NULL REFERENCES questions(id) ON DELETE CASCADE,
    status      progress_status NOT NULL,
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, question_id)
    );

CREATE INDEX IF NOT EXISTS idx_user_progress_user_id ON user_progress(user_id);

CREATE TABLE IF NOT EXISTS bookmarks (
                                         user_id     TEXT   NOT NULL,
                                         question_id BIGINT NOT NULL REFERENCES questions(id) ON DELETE CASCADE,
    bookmarked_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, question_id)
    );

CREATE INDEX IF NOT EXISTS idx_bookmarks_user_id ON bookmarks(user_id);

CREATE TABLE IF NOT EXISTS view_history (
                                            user_id     TEXT   NOT NULL,
                                            question_id BIGINT NOT NULL REFERENCES questions(id) ON DELETE CASCADE,
    viewed_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, question_id)
    );

CREATE INDEX IF NOT EXISTS idx_view_history_user_id_viewed_at ON view_history(user_id, viewed_at);