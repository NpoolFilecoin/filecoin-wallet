/*
 Navicat MySQL Data Transfer

 Source Server         : 165
 Source Server Type    : MySQL
 Source Server Version : 50732
 Source Host           : 127.0.0.1:3306
 Source Schema         : wallet

 Target Server Type    : MySQL
 Target Server Version : 50732
 File Encoding         : 65001

 Date: 30/01/2021 11:40:56
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;
create database if not exists `wallet`; SET character_set_client = utf8; use wallet;
-- ----------------------------
-- Table structure for balance_transfer_request
-- ----------------------------
CREATE TABLE if not exists `balance_transfer_request` (
  `id` varchar(64) CHARACTER SET latin1 COLLATE latin1_swedish_ci NOT NULL,
  `creator` varchar(255) CHARACTER SET latin1 COLLATE latin1_swedish_ci NULL DEFAULT NULL COMMENT '',
  `reviewer` varchar(255) CHARACTER SET latin1 COLLATE latin1_swedish_ci NULL DEFAULT NULL COMMENT '',
  `from` varchar(512) NULL DEFAULT 0 COMMENT '',
  `to` varchar(512) NULL DEFAULT 0 COMMENT '',
  `amount` double(6, 2) NULL DEFAULT 0 COMMENT '',
  `status` varchar(32) NULL DEFAULT 0 COMMENT '',
  `time` datetime NOT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = latin1 COLLATE = latin1_swedish_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for balance_withdraw_request
-- ----------------------------
CREATE TABLE if not exists `balance_withdraw_request` (
  `id` varchar(64) CHARACTER SET latin1 COLLATE latin1_swedish_ci NOT NULL,
  `creator` varchar(255) CHARACTER SET latin1 COLLATE latin1_swedish_ci NULL DEFAULT NULL COMMENT '',
  `reviewer` varchar(255) CHARACTER SET latin1 COLLATE latin1_swedish_ci NULL DEFAULT NULL COMMENT '',
  `owner` varchar(512) NULL DEFAULT 0 COMMENT '',
  `miner` varchar(512) NULL DEFAULT 0 COMMENT '',
  `amount` double(6, 2) NULL DEFAULT 0 COMMENT '',
  `status` varchar(32) NULL DEFAULT 0 COMMENT '',
  `time` datetime NOT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = latin1 COLLATE = latin1_swedish_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for filecoin_customer
-- ----------------------------
CREATE TABLE if not exists `filecoin_customer` (
  `id` varchar(64) CHARACTER SET latin1 COLLATE latin1_swedish_ci NOT NULL,
  `customer_name` varchar(255) CHARACTER SET latin1 COLLATE latin1_swedish_ci NULL DEFAULT NULL COMMENT '',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = latin1 COLLATE = latin1_swedish_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for filecoin_miner
-- ----------------------------
CREATE TABLE if not exists `filecoin_miner` (
  `id` varchar(64) CHARACTER SET latin1 COLLATE latin1_swedish_ci NOT NULL,
  `customer_id` varchar(64) CHARACTER SET latin1 COLLATE latin1_swedish_ci NOT NULL,
  `miner_id` varchar(255) CHARACTER SET latin1 COLLATE latin1_swedish_ci NULL DEFAULT NULL COMMENT '',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = latin1 COLLATE = latin1_swedish_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for filecoin_account
-- ----------------------------
CREATE TABLE if not exists `filecoin_account` (
  `id` varchar(64) CHARACTER SET latin1 COLLATE latin1_swedish_ci NOT NULL,
  `address` varchar(64) CHARACTER SET latin1 COLLATE latin1_swedish_ci NOT NULL,
  `wallet_type` varchar(64) CHARACTER SET latin1 COLLATE latin1_swedish_ci NOT NULL,
  `customer_id` varchar(64) CHARACTER SET latin1 COLLATE latin1_swedish_ci NOT NULL,
  `miner_id` varchar(255) CHARACTER SET latin1 COLLATE latin1_swedish_ci NULL DEFAULT NULL COMMENT '',
  `miner_wallet_type` varchar(255) CHARACTER SET latin1 COLLATE latin1_swedish_ci NULL DEFAULT NULL COMMENT '',
  `have_private_key` boolean,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = latin1 COLLATE = latin1_swedish_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for filecoin_transfer_target
-- ----------------------------
CREATE TABLE if not exists `filecoin_transfer_target` (
  `id` varchar(64) CHARACTER SET latin1 COLLATE latin1_swedish_ci NOT NULL,
  `address` varchar(255) CHARACTER SET latin1 COLLATE latin1_swedish_ci NOT NULL,
  `target_addresses` varchar(32768) CHARACTER SET latin1 COLLATE latin1_swedish_ci NOT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = latin1 COLLATE = latin1_swedish_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for filecoin_transfer_history
-- ----------------------------
CREATE TABLE if not exists `filecoin_transfer_history` (
  `id` varchar(64) CHARACTER SET latin1 COLLATE latin1_swedish_ci NOT NULL,
  `from` varchar(255) CHARACTER SET latin1 COLLATE latin1_swedish_ci NOT NULL,
  `to` varchar(255) CHARACTER SET latin1 COLLATE latin1_swedish_ci NOT NULL,
  `review_status` varchar(64) CHARACTER SET latin1 COLLATE latin1_swedish_ci NOT NULL,
  `status` varchar(64) CHARACTER SET latin1 COLLATE latin1_swedish_ci NOT NULL,
  `cid` varchar(255) CHARACTER SET latin1 COLLATE latin1_swedish_ci NOT NULL,
  `msg` varchar(4096) CHARACTER SET latin1 COLLATE latin1_swedish_ci NOT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = latin1 COLLATE = latin1_swedish_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table for review_history
-- ----------------------------
CREATE TABLE if not exists `review_history` (
  `request_id` varchar(64) CHARACTER SET latin1 COLLATE latin1_swedish_ci NOT NULL,
  `from` varchar(255) CHARACTER SET latin1 COLLATE latin1_swedish_ci NOT NULL,
  `from_owner` varchar(255) CHARACTER SET latin1 COLLATE latin1_swedish_ci NOT NULL,
  `creator` varchar(255) CHARACTER SET latin1 COLLATE latin1_swedish_ci NOT NULL,
  `to` varchar(255) CHARACTER SET latin1 COLLATE latin1_swedish_ci NOT NULL,
  `to_owner` varchar(255) CHARACTER SET latin1 COLLATE latin1_swedish_ci NOT NULL,
  `amount` double(6, 2) CHARACTER SET latin1 COLLATE latin1_swedish_ci NOT NULL,
  `cid` varchar(255) CHARACTER SET latin1 COLLATE latin1_swedish_ci NOT NULL,
  `gas_limit` varchar(255) CHARACTER SET latin1 COLLATE latin1_swedish_ci NOT NULL,
  `gas_feecap` varchar(255) CHARACTER SET latin1 COLLATE latin1_swedish_ci NOT NULL,
  `gas_premium` varchar(255) CHARACTER SET latin1 COLLATE latin1_swedish_ci NOT NULL,
  `reviewer` varchar(255) CHARACTER SET latin1 COLLATE latin1_swedish_ci NOT NULL,
  `status` varchar(255) CHARACTER SET latin1 COLLATE latin1_swedish_ci NOT NULL,
  `time` datetime NOT NULL,
  `type` varchar(255) CHARACTER SET latin1 COLLATE latin1_swedish_ci NOT NULL,
  PRIMARY KEY (`cid`) USING BTREE
)ENGINE = InnoDB CHARACTER SET = latin1 COLLATE = latin1_swedish_ci ROW_FORMAT = Dynamic;
SET FOREIGN_KEY_CHECKS = 1;
