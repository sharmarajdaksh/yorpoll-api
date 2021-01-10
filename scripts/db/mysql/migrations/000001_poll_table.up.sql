CREATE TABLE IF NOT EXISTS `poll` (
    `id` CHAR(36) NOT NULL,
    `title` TEXT NOT NULL,
    `poll_description` TEXT NOT NULL,
    `created` BIGINT NOT NULL,
    `modified` BIGINT NOT NULL,
    `expiry` BIGINT NOT NULL,
    PRIMARY KEY(`id`),
    INDEX (`id`)
);
