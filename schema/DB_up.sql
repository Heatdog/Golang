SET NAMES utf8;

DROP TABLE IF EXISTS userDB;
CREATE TABLE userDB(
    user_id INT AUTO_INCREMENT NOT NULL,
    login VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    token VARCHAR(255),
    PRIMARY KEY (user_id)
);