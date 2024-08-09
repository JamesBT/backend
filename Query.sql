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

-- Dumping structure for table asset.asset
CREATE TABLE IF NOT EXISTS `asset` (
  `id_asset_parent` int(11) NOT NULL AUTO_INCREMENT,
  `nama` varchar(50) NOT NULL,
  `nama_legalitas` varchar(255) NOT NULL,
  `nomor_legalitas` varchar(50) NOT NULL,
  `tipe` enum('L','B','A') NOT NULL DEFAULT 'L',
  `nilai` float NOT NULL,
  `luas` float NOT NULL DEFAULT 0,
  `titik_koordinat` varchar(225) NOT NULL,
  `batas_koordinat` varchar(225) NOT NULL,
  `kondisi` varchar(50) NOT NULL,
  `id_asset_child` varchar(50) DEFAULT '',
  `alamat` text NOT NULL,
  `status_pengecekan` enum('Y','N') NOT NULL DEFAULT 'N',
  `status_verifikasi` enum('Y','N') NOT NULL DEFAULT 'N',
  `hak_akses` varchar(50) DEFAULT '',
  `status_asset` enum('S','T') NOT NULL DEFAULT 'T',
  `masa_sewa` date DEFAULT NULL,
  PRIMARY KEY (`id_asset_parent`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- Dumping data for table asset.asset: ~2 rows (approximately)
REPLACE INTO `asset` (`id_asset_parent`, `nama`, `nama_legalitas`, `nomor_legalitas`, `tipe`, `nilai`, `luas`, `titik_koordinat`, `batas_koordinat`, `kondisi`, `id_asset_child`, `alamat`, `status_pengecekan`, `status_verifikasi`, `hak_akses`, `status_asset`, `masa_sewa`) VALUES
	(1, 'tes1', 'tes1', 'tes1', 'L', 1000000, 100, '-7.27290454460171, 112.74271229250712', '-7.272567138080667, 112.74277687674497', 'bagus', '', 'Jl. Panglima Sudirman No.101-103, Embong Kaliasin, Kec. Genteng, Surabaya, Jawa Timur 60271', 'N', 'N', '', 'T', NULL),
	(2, 'tes2', 'tes2', 'tes2', 'L', 1000000, 100, '-7.27290454460171, 112.74271229250712', '-7.272567138080667, 112.74277687674497', 'bagus', '', 'ini testing ke 2', 'N', 'N', '', 'T', NULL);

-- Dumping structure for table asset.kelas
CREATE TABLE IF NOT EXISTS `kelas` (
  `kelas_id` int(11) NOT NULL AUTO_INCREMENT,
  `kelas_nama` varchar(225) NOT NULL,
  `kelas_modal_minimal` float NOT NULL DEFAULT 0,
  `kelas_modal_maksimal` float NOT NULL DEFAULT 0,
  PRIMARY KEY (`kelas_id`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- Dumping structure for table asset.notification
CREATE TABLE IF NOT EXISTS `notification` (
  `notification_id` int(11) NOT NULL AUTO_INCREMENT,
  `user_id_sender` int(11) NOT NULL DEFAULT 0,
  `user_id_receiver` int(11) NOT NULL DEFAULT 0,
  `notification_detail` text DEFAULT NULL,
  PRIMARY KEY (`notification_id`),
  KEY `user_id_sender` (`user_id_sender`),
  KEY `user_id_receiver` (`user_id_receiver`),
  CONSTRAINT `notification_user_id_receiver` FOREIGN KEY (`user_id_receiver`) REFERENCES `user` (`user_id`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `notification_user_id_sender` FOREIGN KEY (`user_id_sender`) REFERENCES `user` (`user_id`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- Dumping data for table asset.notification: ~0 rows (approximately)

-- Dumping structure for table asset.perusahaan
CREATE TABLE IF NOT EXISTS `perusahaan` (
  `perusahaan_id` int(11) NOT NULL AUTO_INCREMENT,
  `user_id` int(11) DEFAULT NULL,
  `sertifikat_perusahaan` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`perusahaan_id`),
  KEY `user_id` (`user_id`),
  CONSTRAINT `FK_perusahaan_userid` FOREIGN KEY (`user_id`) REFERENCES `user` (`user_id`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- Dumping data for table asset.perusahaan: ~0 rows (approximately)

-- Dumping structure for table asset.privilege
CREATE TABLE IF NOT EXISTS `privilege` (
  `privilege_id` int(11) NOT NULL AUTO_INCREMENT,
  `nama_privilege` varchar(50) DEFAULT NULL,
  PRIMARY KEY (`privilege_id`)
) ENGINE=InnoDB AUTO_INCREMENT=22 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- Dumping data for table asset.privilege: ~21 rows (approximately)
REPLACE INTO `privilege` (`privilege_id`, `nama_privilege`) VALUES
	(1, 'buat role baru'),
	(2, 'menyetujui buat role baru'),
	(3, 'ubah role admin/user'),
	(4, 'lihat aset'),
	(5, 'tambah wishlist'),
	(6, 'ajukan meeting'),
	(7, 'transaksi'),
	(8, 'tambah surveyor'),
	(9, 'lihat data surveyor'),
	(10, 'assign surveyor'),
	(11, 'verifikasi hasil survei'),
	(12, 'tambah aset'),
	(13, 'ubah public/private'),
	(14, 'gabungkan/pecahkan aset'),
	(15, 'tambah kelas'),
	(16, 'set meeting user dan pemilik'),
	(17, 'read only '),
	(18, 'terima assignment'),
	(19, 'submit hasil survei'),
	(20, 'tambah ke wishlist'),
	(21, 'ajukan meeting dengan pemilik');

-- Dumping structure for table asset.role
CREATE TABLE IF NOT EXISTS `role` (
  `role_id` int(11) NOT NULL AUTO_INCREMENT,
  `nama_role` varchar(50) DEFAULT NULL,
  PRIMARY KEY (`role_id`)
) ENGINE=InnoDB AUTO_INCREMENT=12 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- Dumping data for table asset.role: ~11 rows (approximately)
REPLACE INTO `role` (`role_id`, `nama_role`) VALUES
	(1, 'super_admin'),
	(2, 'admin_surveyor'),
	(3, 'admin_verifikator'),
	(4, 'admin_aset'),
	(5, 'admin_mitra'),
	(6, 'admin_direktur'),
	(7, 'surveyor'),
	(8, 'user'),
	(9, 'user_manajer'),
	(10, 'user_staff'),
	(11, 'user_direktur');

-- Dumping structure for table asset.surveyor
CREATE TABLE IF NOT EXISTS `surveyor` (
  `suveyor_id` int(11) NOT NULL AUTO_INCREMENT,
  `user_id` int(11) DEFAULT NULL,
  `lokasi` varchar(50) DEFAULT NULL,
  `availability_surveyor` enum('Y','N') DEFAULT 'Y',
  PRIMARY KEY (`suveyor_id`),
  KEY `id_user` (`user_id`),
  CONSTRAINT `FK_surveyor_userid` FOREIGN KEY (`user_id`) REFERENCES `user` (`user_id`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- Dumping data for table asset.surveyor: ~1 rows (approximately)
REPLACE INTO `surveyor` (`suveyor_id`, `user_id`, `lokasi`, `availability_surveyor`) VALUES
	(1, 5, '', 'Y');

-- Dumping structure for table asset.survey_request
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

-- Dumping data for table asset.survey_request: ~0 rows (approximately)

-- Dumping structure for table asset.transaction_request
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

-- Dumping data for table asset.transaction_request: ~0 rows (approximately)

-- Dumping structure for table asset.user
CREATE TABLE IF NOT EXISTS `user` (
  `user_id` int(11) NOT NULL AUTO_INCREMENT,
  `username` varchar(50) NOT NULL,
  `password` varchar(50) NOT NULL,
  `nama_lengkap` varchar(50) DEFAULT NULL,
  `alamat` varchar(50) DEFAULT '',
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
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- Dumping data for table asset.user: ~5 rows (approximately)
REPLACE INTO `user` (`user_id`, `username`, `password`, `nama_lengkap`, `alamat`, `jenis_kelamin`, `tanggal_lahir`, `email`, `nomor_telepon`, `foto_profil`, `ktp`, `created_at`, `deleted_at`, `updated_at`, `login_timestamp`) VALUES
	(1, 'superadmin1', 'am12345', 'superadmin1', '', 'L', '2024-08-08', 'superadmin1@gmail.com', '08987654321', '', '', '2024-08-08 14:02:36', NULL, NULL, '2024-08-09 11:53:04'),
	(2, 'admin_surveyor', 'am12345', 'admin_surveyor', '', 'L', '2024-08-08', 'admin_surveyor@gmail.com', '08987654321', '', '', '2024-08-08 14:03:05', NULL, NULL, '2024-08-08 14:03:05'),
	(3, 'admin_verifikator', 'am12345', 'admin_verifikator', '', 'L', '2024-08-08', 'admin_verifikator@gmail.com', '08987654321', '', '', '2024-08-08 14:03:17', NULL, NULL, '2024-08-08 14:03:17'),
	(4, 'test_db_server', 'am12345', 'test_db_server', '', 'L', '2024-08-08', 'test_db_server@gmail.com', '08987654321', '', '', '2024-08-08 18:15:29', NULL, NULL, '2024-08-08 18:15:29'),
	(5, 'test_surveyor', 'am12345', 'test_surveyor', '', 'L', '2024-08-09', 'test_surveyor@gmail.com', '08987654321', '', '', '2024-08-09 11:43:54', NULL, NULL, '2024-08-09 11:53:15');

-- Dumping structure for table asset.user_detail
CREATE TABLE IF NOT EXISTS `user_detail` (
  `user_detail_id` int(11) DEFAULT NULL,
  `user_kelas_id` int(11) NOT NULL,
  `status` int(11) DEFAULT NULL,
  `tipe` int(11) DEFAULT NULL,
  `first_login` enum('Y','N') NOT NULL DEFAULT 'Y',
  `denied_by_admin` enum('Y','N') NOT NULL DEFAULT 'N',
  KEY `user_detail_id` (`user_detail_id`),
  KEY `kelas` (`user_kelas_id`) USING BTREE,
  CONSTRAINT `user_detail_id` FOREIGN KEY (`user_detail_id`) REFERENCES `user` (`user_id`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `user_kelas_id` FOREIGN KEY (`user_kelas_id`) REFERENCES `kelas` (`kelas_id`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- Dumping data for table asset.user_detail: ~5 rows (approximately)
REPLACE INTO `user_detail` (`user_detail_id`, `user_kelas_id`, `status`, `tipe`, `first_login`, `denied_by_admin`) VALUES
	(1, 1, 1, 8, 'Y', 'N'),
	(2, 1, 1, 8, 'Y', 'N'),
	(3, 1, 1, 8, 'Y', 'N'),
	(4, 1, 1, 8, 'Y', 'N'),
	(5, 1, 1, 8, 'Y', 'N');

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
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- Dumping data for table asset.user_privilege: ~3 rows (approximately)
REPLACE INTO `user_privilege` (`user_privilege_id`, `privilege_id`, `user_id`) VALUES
	(1, 17, 1),
	(2, 17, 2),
	(3, 17, 3),
	(4, 17, 4),
	(5, 17, 5);

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
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- Dumping data for table asset.user_role: ~5 rows (approximately)
REPLACE INTO `user_role` (`user_role_id`, `user_id`, `role_id`) VALUES
	(1, 1, 8),
	(2, 2, 8),
	(3, 3, 8),
	(4, 4, 8),
	(5, 5, 7);

/*!40103 SET TIME_ZONE=IFNULL(@OLD_TIME_ZONE, 'system') */;
/*!40101 SET SQL_MODE=IFNULL(@OLD_SQL_MODE, '') */;
/*!40014 SET FOREIGN_KEY_CHECKS=IFNULL(@OLD_FOREIGN_KEY_CHECKS, 1) */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40111 SET SQL_NOTES=IFNULL(@OLD_SQL_NOTES, 1) */;
