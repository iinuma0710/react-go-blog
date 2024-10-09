CREATE TABLE `user`
(
    `id`         BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ユーザの識別子',
    `name`       VARCHAR(20) NOT NULL COMMENT 'ユーザ名',
    `password`   VARCHAR(80) NOT NULL COMMENT 'パスワードハッシュ',
    `role`       VARCHAR(80) NOT NULL COMMENT 'ユーザのロール',
    `created_at` DATETIME(6) NOT NULL COMMENT 'レコードの作成日時',
    `updated_at` DATETIME(6) NOT NULL COMMENT 'レコードの更新日時',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uix_name` (`name`) USING BTREE
) ENGINE=INNODB DEFAULT CHARSET=utf8mb4 COMMENT='ユーザ';

CREATE TABLE `article`
(
    `id`         BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '記事の識別子',
    `title`      VARCHAR(128)    NOT NULL COMMENT '記事のタイトル',
    -- `content`    VARCHAR(20)     NOT NULL COMMENT '記事の本文',
    `status`     VARCHAR(20)     NOT NULL COMMENT '記事のステータス',
    -- `author_id`  BIGINT UNSIGNED NOT NULL COMMENT '記事作成者のユーザID',
    `created_at` DATETIME(6)     NOT NULL COMMENT 'レコードの作成日時',
    -- `updated_at` DATETIME(6)     NOT NULL COMMENT 'レコードの更新日時',
    PRIMARY KEY (`id`)
) ENGINE=INNODB DEFAULT CHARSET=utf8mb4 COMMENT='ブログ記事';