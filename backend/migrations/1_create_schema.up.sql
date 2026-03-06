CREATE TYPE question_level AS ENUM ('junior', 'middle', 'senior');

CREATE TABLE topics (
    id          SERIAL PRIMARY KEY,
    slug        VARCHAR(100) NOT NULL UNIQUE,
    name        VARCHAR(200) NOT NULL,
    description TEXT,
    icon        VARCHAR(10),
    sort_order  INT          NOT NULL DEFAULT 0
);

CREATE TABLE tags (
    id         SERIAL PRIMARY KEY,
    slug       VARCHAR(100) NOT NULL UNIQUE,
    name       VARCHAR(200) NOT NULL
);

CREATE TABLE questions (
    id         SERIAL PRIMARY KEY,
    topic_id   INT REFERENCES topics (id) ON DELETE SET NULL,
    title      VARCHAR(500)   NOT NULL,
    slug       VARCHAR(200)   NOT NULL UNIQUE,
    answer     TEXT           NOT NULL DEFAULT '',
    level      question_level NOT NULL
);

CREATE INDEX idx_questions_topic_id ON questions (topic_id);
CREATE INDEX idx_questions_level    ON questions (level);
CREATE INDEX idx_questions_slug     ON questions (slug);

CREATE TABLE question_tags (
    question_id INT NOT NULL REFERENCES questions (id) ON DELETE CASCADE,
    tag_id      INT NOT NULL REFERENCES tags     (id) ON DELETE CASCADE,
    PRIMARY KEY (question_id, tag_id)
);

CREATE INDEX idx_question_tags_tag_id ON question_tags (tag_id);