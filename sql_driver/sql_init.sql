CREATE DATABASE IF NOT EXISTS stockhome;
USE stockhome;

DROP TABLES IF EXISTS `stock`, `history`;

CREATE TABLE `stock` (
  `itemID` int DEFAULT NULL,
  `amount` int DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `history` (
  `historyID` int NOT NULL AUTO_INCREMENT,
  `action` int DEFAULT NULL,
  `userID` int DEFAULT NULL,
  `itemID` int DEFAULT NULL,
  `amount` int DEFAULT NULL,
  `date` text,
  `time` text,
  PRIMARY KEY (`historyID`)
) ENGINE=InnoDB AUTO_INCREMENT=539 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
