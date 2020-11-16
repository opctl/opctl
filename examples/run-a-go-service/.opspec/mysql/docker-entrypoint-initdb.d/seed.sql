CREATE TABLE `users` (
    `uid` INT(10) NOT NULL AUTO_INCREMENT,
    `username` VARCHAR(64) NULL DEFAULT NULL,
    PRIMARY KEY (`uid`)
);
INSERT INTO users (`uid`, `username`) VALUES(1337,"tester");