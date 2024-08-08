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

-- Dumping structure for table asset.user
CREATE TABLE IF NOT EXISTS `user` (
  `user_id` int(11) NOT NULL AUTO_INCREMENT,
  `username` varchar(50) NOT NULL,
  `password` varchar(50) NOT NULL,
  `nama_lengkap` varchar(50) DEFAULT NULL,
  `alamat` varchar(50) DEFAULT NULL,
  `jenis_kelamin` enum('L','P') DEFAULT 'L',
  `tanggal_lahir` date DEFAULT NULL,
  `email` varchar(50) NOT NULL,
  `nomor_telepon` varchar(13) NOT NULL,
  `foto_profil` varchar(50) DEFAULT '',
  `ktp` varchar(225) DEFAULT '',
  `created_at` datetime DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `login_timestamp` datetime DEFAULT NULL,
  PRIMARY KEY (`user_id`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- Dumping data for table asset.user: ~3 rows (approximately)
REPLACE INTO `user` (`user_id`, `username`, `password`, `nama_lengkap`, `alamat`, `jenis_kelamin`, `tanggal_lahir`, `email`, `nomor_telepon`, `foto_profil`, `ktp`, `created_at`, `deleted_at`, `updated_at`, `login_timestamp`) VALUES
	(1, 'superadmin1', 'am12345', 'superadmin1', NULL, 'L', '2024-08-08', 'superadmin1@gmail.com', '08987654321', '', '', '2024-08-08 14:02:36', NULL, NULL, '2024-08-08 14:02:36'),
	(2, 'admin_surveyor', 'am12345', 'admin_surveyor', NULL, 'L', '2024-08-08', 'admin_surveyor@gmail.com', '08987654321', '', '', '2024-08-08 14:03:05', NULL, NULL, '2024-08-08 14:03:05'),
	(3, 'admin_verifikator', 'am12345', 'admin_verifikator', NULL, 'L', '2024-08-08', 'admin_verifikator@gmail.com', '08987654321', '', '', '2024-08-08 14:03:17', NULL, NULL, '2024-08-08 14:03:17');

-- Dumping structure for table asset.user_detail
CREATE TABLE IF NOT EXISTS `user_detail` (
  `user_detail_id` int(11) DEFAULT NULL,
  `kelas` int(11) NOT NULL,
  `status` int(11) DEFAULT NULL,
  `tipe` int(11) DEFAULT NULL,
  `first_login` enum('Y','N') NOT NULL DEFAULT 'Y',
  `denied_by_admin` enum('Y','N') NOT NULL DEFAULT 'N',
  KEY `user_detail_id` (`user_detail_id`),
  KEY `kelas` (`kelas`),
  CONSTRAINT `user_detail_id` FOREIGN KEY (`user_detail_id`) REFERENCES `user` (`user_id`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `user_kelas_id` FOREIGN KEY (`kelas`) REFERENCES `kelas` (`kelas_id`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- Dumping data for table asset.user_detail: ~2 rows (approximately)
REPLACE INTO `user_detail` (`user_detail_id`, `kelas`, `status`, `tipe`, `first_login`, `denied_by_admin`) VALUES
	(1, 1, 1, 8, 'Y', 'N'),
	(2, 1, 1, 8, 'Y', 'N'),
	(3, 1, 1, 8, 'Y', 'N');

-- Dumping structure for table asset.user_privilege
CREATE TABLE IF NOT EXISTS `user_privilege` (
  `user_privilege_id` int(11) NOT NULL AUTO_INCREMENT,
  `privilege_id` int(11) NOT NULL,
  `user_id` int(11) NOT NULL,
  PRIMARY KEY (`user_privilege_id`),
  KEY `Column 1` (`privilege_id`) USING BTREE,
  KEY `Column 2` (`user_id`) USING BTREE,
  CONSTRAINT `FK_userprivilege_privilegeid` FOREIGN KEY (`privilege_id`) REFERENCES `privilege` (`privilege_id`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `FK_userprivilege_userid` FOREIGN KEY (`user_id`) REFERENCES `user` (`user_id`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- Dumping data for table asset.user_privilege: ~2 rows (approximately)
REPLACE INTO `user_privilege` (`user_privilege_id`, `privilege_id`, `user_id`) VALUES
	(1, 17, 1),
	(2, 17, 2),
	(3, 17, 3);

-- Dumping structure for table asset.user_role
CREATE TABLE IF NOT EXISTS `user_role` (
  `user_role_id` int(11) NOT NULL AUTO_INCREMENT,
  `user_id` int(11) NOT NULL,
  `role_id` int(11) NOT NULL,
  PRIMARY KEY (`user_role_id`),
  KEY `user_id` (`user_id`),
  KEY `role_id` (`role_id`),
  CONSTRAINT `FK_userrole_roleid` FOREIGN KEY (`role_id`) REFERENCES `role` (`role_id`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `FK_userrole_userid` FOREIGN KEY (`user_id`) REFERENCES `user` (`user_id`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- Dumping data for table asset.user_role: ~2 rows (approximately)
REPLACE INTO `user_role` (`user_role_id`, `user_id`, `role_id`) VALUES
	(1, 1, 8),
	(2, 2, 8),
	(3, 3, 8);

/*!40103 SET TIME_ZONE=IFNULL(@OLD_TIME_ZONE, 'system') */;
/*!40101 SET SQL_MODE=IFNULL(@OLD_SQL_MODE, '') */;
/*!40014 SET FOREIGN_KEY_CHECKS=IFNULL(@OLD_FOREIGN_KEY_CHECKS, 1) */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40111 SET SQL_NOTES=IFNULL(@OLD_SQL_NOTES, 1) */;
