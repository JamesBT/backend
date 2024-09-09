-- --------------------------------------------------------
-- Host:                         127.0.0.1
-- Server version:               11.4.2-MariaDB - mariadb.org binary distribution
-- Server OS:                    Win64
-- HeidiSQL Version:             12.6.0.6765
-- --------------------------------------------------------

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET NAMES utf8 */;
/*!50503 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;


-- Dumping database structure for asset
CREATE DATABASE IF NOT EXISTS `asset` /*!40100 DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci */;
USE `asset`;

-- Dumping structure for table asset.progress
CREATE TABLE IF NOT EXISTS `progress` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_id` int(11) DEFAULT NULL,
  `id_asset` int(11) DEFAULT NULL,
  `nama` varchar(255) NOT NULL DEFAULT '',
  `proposal` varchar(255) NOT NULL DEFAULT '',
  `tanggal_meeting` date DEFAULT NULL,
  `waktu_meeting` time DEFAULT NULL,
  `tempat_meeting` varchar(255) DEFAULT NULL,
  `waktu_mulai_meeting` time DEFAULT NULL,
  `waktu_selesai_meeting` time DEFAULT NULL,
  `notes` text DEFAULT NULL,
  `file` varchar(255) DEFAULT NULL,
  `tipe_file` enum('L','C','A') DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `user_id` (`user_id`),
  KEY `id_asset` (`id_asset`),
  KEY `proposal` (`proposal`),
  KEY `nama_progress` (`nama`) USING BTREE,
  CONSTRAINT `fk_progress_asset` FOREIGN KEY (`id_asset`) REFERENCES `asset` (`id_asset`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `fk_progress_nama` FOREIGN KEY (`nama`) REFERENCES `transaction_request` (`nama_progress`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `fk_progress_proposal` FOREIGN KEY (`proposal`) REFERENCES `transaction_request` (`proposal`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `fk_progress_user` FOREIGN KEY (`user_id`) REFERENCES `user` (`user_id`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- Dumping data for table asset.progress: ~0 rows (approximately)

/*!40103 SET TIME_ZONE=IFNULL(@OLD_TIME_ZONE, 'system') */;
/*!40101 SET SQL_MODE=IFNULL(@OLD_SQL_MODE, '') */;
/*!40014 SET FOREIGN_KEY_CHECKS=IFNULL(@OLD_FOREIGN_KEY_CHECKS, 1) */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40111 SET SQL_NOTES=IFNULL(@OLD_SQL_NOTES, 1) */;
