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
  `id_asset` int(11) NOT NULL AUTO_INCREMENT,
  `id_asset_parent` int(11) DEFAULT NULL,
  `id_asset_child` int(11) DEFAULT NULL,
  `id_join` varchar(50) DEFAULT NULL,
  `perusahaan_id` int(11) DEFAULT NULL,
  `nama` varchar(50) NOT NULL,
  `tipe` enum('L','B','A') NOT NULL DEFAULT 'L',
  `nomor_legalitas` varchar(50) NOT NULL,
  `file_legalitas` varchar(225) NOT NULL DEFAULT '',
  `status_asset` enum('S','T') NOT NULL DEFAULT 'T',
  `surat_kuasa` varchar(225) NOT NULL DEFAULT '',
  `alamat` text NOT NULL DEFAULT '',
  `kondisi` text NOT NULL DEFAULT '',
  `titik_koordinat` text NOT NULL DEFAULT '',
  `batas_koordinat` text NOT NULL DEFAULT '',
  `luas` float NOT NULL DEFAULT 0,
  `nilai` float NOT NULL DEFAULT 0,
  `provinsi` varchar(255) NOT NULL DEFAULT '',
  `usage` text NOT NULL DEFAULT '\'\'',
  `status_pengecekan` enum('Y','N') NOT NULL DEFAULT 'N',
  `status_verifikasi` enum('Y','N') NOT NULL DEFAULT 'N',
  `status_publik` enum('Y','N') NOT NULL DEFAULT 'N',
  `hak_akses` varchar(50) NOT NULL DEFAULT '',
  `masa_sewa` date DEFAULT NULL,
  `created_at` datetime NOT NULL DEFAULT current_timestamp(),
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id_asset`) USING BTREE,
  KEY `perusahaan_id` (`perusahaan_id`),
  KEY `id_asset_parent` (`id_asset_parent`),
  KEY `id_asset_child` (`id_asset_child`),
  CONSTRAINT `id_asset_child` FOREIGN KEY (`id_asset_child`) REFERENCES `asset` (`id_asset`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `id_asset_parent` FOREIGN KEY (`id_asset_parent`) REFERENCES `asset` (`id_asset`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `id_asset_perusahaan` FOREIGN KEY (`perusahaan_id`) REFERENCES `perusahaan` (`perusahaan_id`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- Dumping data for table asset.asset: ~0 rows (approximately)

-- Dumping structure for table asset.asset_gambar
CREATE TABLE IF NOT EXISTS `asset_gambar` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `id_asset_gambar` int(11) NOT NULL,
  `link_gambar` varchar(225) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`),
  KEY `fk_asset_gambar` (`id_asset_gambar`),
  CONSTRAINT `fk_asset_gambar` FOREIGN KEY (`id_asset_gambar`) REFERENCES `asset` (`id_asset`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- Dumping data for table asset.asset_gambar: ~0 rows (approximately)

-- Dumping structure for table asset.asset_tags
CREATE TABLE IF NOT EXISTS `asset_tags` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `id_asset` int(11) NOT NULL,
  `id_tags` int(11) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `fk_asset` (`id_asset`),
  KEY `fk_tags` (`id_tags`),
  CONSTRAINT `fk_id_aset` FOREIGN KEY (`id_asset`) REFERENCES `asset` (`id_asset`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `fk_tags` FOREIGN KEY (`id_tags`) REFERENCES `tags` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- Dumping data for table asset.asset_tags: ~0 rows (approximately)

-- Dumping structure for table asset.kelas
CREATE TABLE IF NOT EXISTS `kelas` (
  `kelas_id` int(11) NOT NULL AUTO_INCREMENT,
  `kelas_nama` varchar(225) NOT NULL,
  `kelas_modal_minimal` float NOT NULL DEFAULT 0,
  `kelas_modal_maksimal` float NOT NULL DEFAULT 0,
  PRIMARY KEY (`kelas_id`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- Dumping data for table asset.kelas: ~4 rows (approximately)
REPLACE INTO `kelas` (`kelas_id`, `kelas_nama`, `kelas_modal_minimal`, `kelas_modal_maksimal`) VALUES
	(1, 'mikro', 0, 500000000),
	(2, 'kecil', 0, 2500000000),
	(3, 'menengah', 0, 10000000000),
	(4, 'makro', 0, 100000000000);

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
  `status` enum('V','N') NOT NULL DEFAULT 'N',
  `denied_by_admin` enum('Y','N') NOT NULL DEFAULT 'N',
  `name` varchar(255) NOT NULL,
  `username` varchar(50) NOT NULL,
  `lokasi` varchar(255) NOT NULL,
  `tipe` enum('L','B','A') NOT NULL,
  `kelas` int(11) NOT NULL DEFAULT 1,
  `dokumen_kepemilikan` varchar(255) NOT NULL DEFAULT '',
  `dokumen_perusahaan` varchar(255) NOT NULL DEFAULT '',
  `modal_awal` float NOT NULL,
  `deskripsi` text NOT NULL DEFAULT '',
  `created_at` datetime DEFAULT current_timestamp(),
  `updated_at` datetime DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  `login_timestamp` datetime DEFAULT current_timestamp(),
  PRIMARY KEY (`perusahaan_id`),
  KEY `kelas` (`kelas`),
  CONSTRAINT `fk_perusahaan_kelas` FOREIGN KEY (`kelas`) REFERENCES `kelas` (`kelas_id`) ON DELETE NO ACTION ON UPDATE NO ACTION
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
  `lokasi` varchar(255) DEFAULT NULL,
  `availability_surveyor` enum('Y','N') DEFAULT 'Y',
  PRIMARY KEY (`suveyor_id`),
  KEY `id_user` (`user_id`),
  CONSTRAINT `FK_surveyor_userid` FOREIGN KEY (`user_id`) REFERENCES `user` (`user_id`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- Dumping data for table asset.surveyor: ~0 rows (approximately)

-- Dumping structure for table asset.survey_request
CREATE TABLE IF NOT EXISTS `survey_request` (
  `id_transaksi_jual_sewa` int(11) NOT NULL AUTO_INCREMENT,
  `user_id` int(11) NOT NULL DEFAULT 0,
  `id_asset` int(11) NOT NULL DEFAULT 0,
  `created_at` datetime NOT NULL DEFAULT current_timestamp(),
  `dateline` date NOT NULL,
  `status_request` enum('O','F','R') NOT NULL DEFAULT 'O',
  `data_lengkap` enum('Y','N') NOT NULL DEFAULT 'N',
  `usage_old` varchar(255) NOT NULL DEFAULT '',
  `usage_new` varchar(255) NOT NULL DEFAULT '',
  `luas_old` float NOT NULL DEFAULT 0,
  `luas_new` float NOT NULL DEFAULT 0,
  `nilai_old` float NOT NULL DEFAULT 0,
  `nilai_new` float NOT NULL DEFAULT 0,
  `kondisi_old` text NOT NULL DEFAULT '',
  `kondisi_new` text NOT NULL DEFAULT '',
  `batas_koordinat_old` text NOT NULL DEFAULT '',
  `batas_koordinat_new` text NOT NULL DEFAULT '',
  `tags_old` int(11) DEFAULT NULL,
  `tags_new` int(11) DEFAULT NULL,
  PRIMARY KEY (`id_transaksi_jual_sewa`),
  KEY `user_id` (`user_id`),
  KEY `id_asset` (`id_asset`),
  KEY `tags_old` (`tags_old`),
  KEY `tags_new` (`tags_new`),
  CONSTRAINT `id_surveryreq_tag_old` FOREIGN KEY (`tags_old`) REFERENCES `tags` (`id`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `id_surveyreq_asset` FOREIGN KEY (`id_asset`) REFERENCES `asset` (`id_asset`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `id_surveyreq_tag_new` FOREIGN KEY (`tags_new`) REFERENCES `tags` (`id`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `id_surveyreq_userid` FOREIGN KEY (`user_id`) REFERENCES `user` (`user_id`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- Dumping data for table asset.survey_request: ~1 rows (approximately)

-- Dumping structure for table asset.tags
CREATE TABLE IF NOT EXISTS `tags` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `nama` varchar(255) NOT NULL DEFAULT '',
  `detail` text NOT NULL DEFAULT '',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- Dumping data for table asset.tags: ~2 rows (approximately)
REPLACE INTO `tags` (`id`, `nama`, `detail`) VALUES
	(1, 'cafe', ''),
	(2, 'kantor', '');

-- Dumping structure for table asset.transaction_request
CREATE TABLE IF NOT EXISTS `transaction_request` (
  `id_transaksi_jual_sewa` int(11) NOT NULL AUTO_INCREMENT,
  `user_id` int(11) DEFAULT NULL,
  `id_asset` int(11) DEFAULT NULL,
  `perusahaan_id` int(11) NOT NULL,
  `tipe` varchar(50) DEFAULT NULL,
  `masa_sewa` date DEFAULT NULL,
  `meeting_log` text DEFAULT NULL,
  PRIMARY KEY (`id_transaksi_jual_sewa`),
  KEY `user_id` (`user_id`),
  KEY `id_asset` (`id_asset`),
  KEY `perusahaan_id` (`perusahaan_id`),
  CONSTRAINT `FK_transactionreq_userid` FOREIGN KEY (`user_id`) REFERENCES `user` (`user_id`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `id_transactionreq_asset` FOREIGN KEY (`id_asset`) REFERENCES `asset` (`id_asset`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `id_transactionreq_perusahaan` FOREIGN KEY (`perusahaan_id`) REFERENCES `perusahaan` (`perusahaan_id`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- Dumping data for table asset.transaction_request: ~0 rows (approximately)

-- Dumping structure for table asset.user
CREATE TABLE IF NOT EXISTS `user` (
  `user_id` int(11) NOT NULL AUTO_INCREMENT,
  `username` varchar(50) NOT NULL,
  `password` varchar(50) NOT NULL,
  `nama_lengkap` varchar(50) DEFAULT NULL,
  `alamat` varchar(50) NOT NULL DEFAULT '',
  `jenis_kelamin` enum('L','P') DEFAULT 'L',
  `tanggal_lahir` date DEFAULT NULL,
  `email` varchar(50) NOT NULL,
  `nomor_telepon` varchar(13) NOT NULL,
  `foto_profil` varchar(50) DEFAULT '',
  `ktp` varchar(225) DEFAULT '',
  `created_at` datetime DEFAULT current_timestamp(),
  `deleted_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `login_timestamp` datetime DEFAULT current_timestamp(),
  PRIMARY KEY (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- Dumping data for table asset.user: ~0 rows (approximately)

-- Dumping structure for table asset.user_detail
CREATE TABLE IF NOT EXISTS `user_detail` (
  `id_user_detail` int(11) NOT NULL AUTO_INCREMENT,
  `user_detail_id` int(11) DEFAULT NULL,
  `user_kelas_id` int(11) NOT NULL,
  `status` enum('V','N') NOT NULL DEFAULT 'N',
  `tipe` int(11) DEFAULT NULL,
  `first_login` enum('Y','N') NOT NULL DEFAULT 'Y',
  `denied_by_admin` enum('Y','N') NOT NULL DEFAULT 'N',
  `kode_otp` int(4) NOT NULL DEFAULT 0,
  `status_verifikasi_otp` enum('V','N') NOT NULL DEFAULT 'N',
  PRIMARY KEY (`id_user_detail`),
  KEY `user_detail_id` (`user_detail_id`),
  KEY `kelas` (`user_kelas_id`),
  CONSTRAINT `user_detail_id` FOREIGN KEY (`user_detail_id`) REFERENCES `user` (`user_id`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `user_kelas_id` FOREIGN KEY (`user_kelas_id`) REFERENCES `kelas` (`kelas_id`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- Dumping data for table asset.user_detail: ~0 rows (approximately)

-- Dumping structure for table asset.user_perusahaan
CREATE TABLE IF NOT EXISTS `user_perusahaan` (
  `id_user_perusahaan` int(11) NOT NULL AUTO_INCREMENT,
  `id_user` int(11) DEFAULT NULL,
  `id_perusahaan` int(11) DEFAULT NULL,
  PRIMARY KEY (`id_user_perusahaan`),
  KEY `id_user` (`id_user`),
  KEY `id_perusahaan` (`id_perusahaan`),
  CONSTRAINT `id_userperusahaan_perusahaan` FOREIGN KEY (`id_perusahaan`) REFERENCES `perusahaan` (`perusahaan_id`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `id_userperusahaan_user` FOREIGN KEY (`id_user`) REFERENCES `user` (`user_id`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- Dumping data for table asset.user_perusahaan: ~0 rows (approximately)

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
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- Dumping data for table asset.user_privilege: ~0 rows (approximately)

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
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- Dumping data for table asset.user_role: ~0 rows (approximately)

/*!40103 SET TIME_ZONE=IFNULL(@OLD_TIME_ZONE, 'system') */;
/*!40101 SET SQL_MODE=IFNULL(@OLD_SQL_MODE, '') */;
/*!40014 SET FOREIGN_KEY_CHECKS=IFNULL(@OLD_FOREIGN_KEY_CHECKS, 1) */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40111 SET SQL_NOTES=IFNULL(@OLD_SQL_NOTES, 1) */;
