CREATE TABLE IF NOT EXISTS `poll_option` (
    `id` CHAR(36) NOT NULL,
    `title` TEXT NOT NULL,
    `votes` BIGINT NOT NULL,
    `poll_id` CHAR(36) NOT NULL, 
    PRIMARY KEY(`id`),
    INDEX(`id`),
    FOREIGN KEY (`poll_id`) REFERENCES poll(`id`) ON DELETE CASCADE
);
