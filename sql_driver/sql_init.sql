CREATE DATABASE IF NOT EXISTS stockhome;
USE stockhome;

DROP TABLES IF EXISTS `stock`, `history`;

CREATE TABLE IF NOT EXISTS `stock` (
  `itemID` int unsigned NOT NULL,
  `amount` int NOT NULL,
  PRIMARY KEY (`ItemID`)
);

CREATE TABLE IF NOT EXISTS `history` (
  `historyID` int unsigned NOT NULL AUTO_INCREMENT,
  `action` bool NOT NULL,
  `userID` int unsigned NOT NULL,
  `itemID` int unsigned NOT NULL,
  `amount` int unsigned NOT NULL,
  `date` date NOT NULL,
  `time` time NOT NULL,
  PRIMARY KEY (`historyID`)
) auto_increment 1;
