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


-- Dumping database structure for asetmanajemen
CREATE DATABASE IF NOT EXISTS `asetmanajemen` /*!40100 DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci */;
USE `asetmanajemen`;

-- Dumping structure for table asetmanajemen.asset
CREATE TABLE IF NOT EXISTS `asset` (
  `id_asset_parent` int(11) NOT NULL AUTO_INCREMENT,
  `nama` varchar(50) DEFAULT NULL,
  `nama_legalitas` varchar(255) DEFAULT NULL,
  `nomor_legalitas` varchar(50) DEFAULT NULL,
  `tipe` varchar(50) DEFAULT NULL,
  `nilai` int(11) DEFAULT 0,
  `luas` float DEFAULT 0,
  `titik_koordinat` varchar(50) DEFAULT NULL,
  `batas_koordinat` varchar(50) DEFAULT NULL,
  `kondisi` varchar(50) DEFAULT NULL,
  `id_asset_child` varchar(50) DEFAULT NULL,
  `alamat` varchar(50) DEFAULT NULL,
  `status_pengecekan` varchar(50) DEFAULT NULL,
  `status_verifikasi` varchar(50) DEFAULT NULL,
  `hak_akses` varchar(50) DEFAULT NULL,
  `status_asset` varchar(50) DEFAULT NULL,
  `masa_sewa` date DEFAULT NULL,
  PRIMARY KEY (`id_asset_parent`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- Dumping data for table asetmanajemen.asset: ~0 rows (approximately)

-- Dumping structure for table asetmanajemen.perusahaan
CREATE TABLE IF NOT EXISTS `perusahaan` (
  `perusahaan_id` int(11) NOT NULL AUTO_INCREMENT,
  `user_id` int(11) DEFAULT NULL,
  `sertifikat_perusahaan` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`perusahaan_id`),
  KEY `user_id` (`user_id`),
  CONSTRAINT `FK_perusahaan_userid` FOREIGN KEY (`user_id`) REFERENCES `user` (`user_id`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- Dumping data for table asetmanajemen.perusahaan: ~0 rows (approximately)

-- Dumping structure for table asetmanajemen.privilege
CREATE TABLE IF NOT EXISTS `privilege` (
  `privilege_id` int(11) NOT NULL AUTO_INCREMENT,
  `nama_privilege` varchar(50) DEFAULT NULL,
  PRIMARY KEY (`privilege_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- Dumping data for table asetmanajemen.privilege: ~0 rows (approximately)

-- Dumping structure for table asetmanajemen.role
CREATE TABLE IF NOT EXISTS `role` (
  `role_id` int(11) NOT NULL AUTO_INCREMENT,
  `nama_role` varchar(50) DEFAULT NULL,
  PRIMARY KEY (`role_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- Dumping data for table asetmanajemen.role: ~0 rows (approximately)

-- Dumping structure for table asetmanajemen.surveyor
CREATE TABLE IF NOT EXISTS `surveyor` (
  `suveyor_id` int(11) NOT NULL AUTO_INCREMENT,
  `user_id` int(11) DEFAULT NULL,
  `lokasi` varchar(50) DEFAULT NULL,
  `availability_surveyor` int(1) DEFAULT 0,
  PRIMARY KEY (`suveyor_id`),
  KEY `id_user` (`user_id`),
  CONSTRAINT `FK_surveyor_userid` FOREIGN KEY (`user_id`) REFERENCES `user` (`user_id`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- Dumping data for table asetmanajemen.surveyor: ~0 rows (approximately)

-- Dumping structure for table asetmanajemen.survey_request
CREATE TABLE IF NOT EXISTS `survey_request` (
  `id_transaksi_jual_sewa` int(11) NOT NULL AUTO_INCREMENT,
  `user_id` int(11) NOT NULL DEFAULT 0,
  `id_asset` int(11) NOT NULL DEFAULT 0,
  `dateline` date DEFAULT NULL,
  `status_request` varchar(50) DEFAULT NULL,
  PRIMARY KEY (`id_transaksi_jual_sewa`),
  KEY `user_id` (`user_id`),
  KEY `id_asset` (`id_asset`),
  CONSTRAINT `FK_surveyreq_idasset` FOREIGN KEY (`id_asset`) REFERENCES `asset` (`id_asset_parent`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `FK_surveyreq_userid` FOREIGN KEY (`user_id`) REFERENCES `user` (`user_id`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- Dumping data for table asetmanajemen.survey_request: ~0 rows (approximately)

-- Dumping structure for table asetmanajemen.transaction_request
CREATE TABLE IF NOT EXISTS `transaction_request` (
  `id_transaksi_jual_sewa` int(11) NOT NULL AUTO_INCREMENT,
  `user_id` int(11) DEFAULT NULL,
  `id_asset` int(11) DEFAULT NULL,
  `tipe` varchar(50) DEFAULT NULL,
  `masa_sewa` date DEFAULT NULL,
  `meeting_log` text DEFAULT NULL,
  PRIMARY KEY (`id_transaksi_jual_sewa`),
  KEY `user_id` (`user_id`),
  KEY `id_asset` (`id_asset`),
  CONSTRAINT `FK_transactionreq_assetid` FOREIGN KEY (`id_asset`) REFERENCES `asset` (`id_asset_parent`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `FK_transactionreq_userid` FOREIGN KEY (`user_id`) REFERENCES `user` (`user_id`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- Dumping data for table asetmanajemen.transaction_request: ~0 rows (approximately)

-- Dumping structure for table asetmanajemen.user
CREATE TABLE IF NOT EXISTS `user` (
  `user_id` int(11) NOT NULL AUTO_INCREMENT,
  `username` varchar(50) DEFAULT NULL,
  `password` varchar(50) DEFAULT NULL,
  `nama_lengkap` varchar(50) DEFAULT NULL,
  `alamat` varchar(50) DEFAULT NULL,
  `jenis_kelamin` varchar(1) DEFAULT NULL,
  `tanggal_lahir` date DEFAULT NULL,
  `email` varchar(50) DEFAULT NULL,
  `nomor_telepon` varchar(13) DEFAULT NULL,
  `foto_profil` varchar(50) DEFAULT NULL,
  `created_at` datetime DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `login_timestamp` datetime DEFAULT NULL,
  `ktp` varchar(225) DEFAULT NULL,
  PRIMARY KEY (`user_id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- Dumping data for table asetmanajemen.user: ~1 rows (approximately)
REPLACE INTO `user` (`user_id`, `username`, `password`, `nama_lengkap`, `alamat`, `jenis_kelamin`, `tanggal_lahir`, `email`, `nomor_telepon`, `foto_profil`, `created_at`, `deleted_at`, `updated_at`, `login_timestamp`, `ktp`) VALUES
	(1, 'tes1', 'tes1', 'tes1', 'tes1', 'L', '2024-07-16', 'tes1@gmail.com', '08123456789', 'tes1', '2024-07-16 16:45:12', '2024-07-16 16:45:13', '2024-07-16 16:45:14', '2024-07-16 18:17:46', NULL);

-- Dumping structure for table asetmanajemen.user_privilege
CREATE TABLE IF NOT EXISTS `user_privilege` (
  `user_privilege_id` int(11) NOT NULL AUTO_INCREMENT,
  `privilege_id` int(11) DEFAULT NULL,
  `user_id` int(11) DEFAULT NULL,
  PRIMARY KEY (`user_privilege_id`),
  KEY `Column 1` (`privilege_id`) USING BTREE,
  KEY `Column 2` (`user_id`) USING BTREE,
  CONSTRAINT `FK_userprivilege_privilegeid` FOREIGN KEY (`privilege_id`) REFERENCES `privilege` (`privilege_id`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `FK_userprivilege_userid` FOREIGN KEY (`user_id`) REFERENCES `user` (`user_id`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- Dumping data for table asetmanajemen.user_privilege: ~0 rows (approximately)

-- Dumping structure for table asetmanajemen.user_role
CREATE TABLE IF NOT EXISTS `user_role` (
  `user_role_id` int(11) NOT NULL AUTO_INCREMENT,
  `user_id` int(11) DEFAULT NULL,
  `role_id` int(11) DEFAULT NULL,
  PRIMARY KEY (`user_role_id`),
  KEY `user_id` (`user_id`),
  KEY `role_id` (`role_id`),
  CONSTRAINT `FK_userrole_roleid` FOREIGN KEY (`role_id`) REFERENCES `role` (`role_id`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `FK_userrole_userid` FOREIGN KEY (`user_id`) REFERENCES `user` (`user_id`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- Dumping data for table asetmanajemen.user_role: ~0 rows (approximately)

/*!40103 SET TIME_ZONE=IFNULL(@OLD_TIME_ZONE, 'system') */;
/*!40101 SET SQL_MODE=IFNULL(@OLD_SQL_MODE, '') */;
/*!40014 SET FOREIGN_KEY_CHECKS=IFNULL(@OLD_FOREIGN_KEY_CHECKS, 1) */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40111 SET SQL_NOTES=IFNULL(@OLD_SQL_NOTES, 1) */;
