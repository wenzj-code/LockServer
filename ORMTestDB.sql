/*
Navicat MySQL Data Transfer

Source Server         : 192.168.2.142
Source Server Version : 50535
Source Host           : 192.168.2.142:3306
Source Database       : ORMTestDB

Target Server Type    : MYSQL
Target Server Version : 50535
File Encoding         : 65001

Date: 2018-09-20 11:12:11
*/

SET FOREIGN_KEY_CHECKS=0;

-- ----------------------------
-- Table structure for hotel_base_info
-- ----------------------------
DROP TABLE IF EXISTS `hotel_base_info`;
CREATE TABLE `hotel_base_info` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `appid` varchar(255) DEFAULT NULL COMMENT '取token appid',
  `secret` varchar(255) DEFAULT NULL COMMENT '取token secret',
  `app_domain` varchar(100) DEFAULT NULL COMMENT '应用服务器域名',
  `device_domain` varchar(100) DEFAULT NULL COMMENT '设备服务器域名',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=52 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for hotel_public_room
-- ----------------------------
DROP TABLE IF EXISTS `hotel_public_room`;
CREATE TABLE `hotel_public_room` (
  `id` int(11) NOT NULL,
  `gw_id` int(11) NOT NULL,
  `device_id` varchar(32) NOT NULL,
  `roomdevno` varchar(255) NOT NULL,
  `status` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- ----------------------------
-- Table structure for hotel_room_info
-- ----------------------------
DROP TABLE IF EXISTS `hotel_room_info`;
CREATE TABLE `hotel_room_info` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `gw_id` bigint(32) DEFAULT NULL COMMENT '网关id',
  `device_id` varchar(32) DEFAULT '' COMMENT '设备id，唯一',
  `roomdevno` varchar(255) DEFAULT NULL COMMENT '设备房间号 设备服务器使用该字段 ',
  `status` int(11) DEFAULT NULL,
  `barry` float DEFAULT NULL,
  `hotel_id` bigint(32) DEFAULT NULL COMMENT '公寓表主键ID',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=875 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for t_gateway_info
-- ----------------------------
DROP TABLE IF EXISTS `t_gateway_info`;
CREATE TABLE `t_gateway_info` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `deleted` bit(1) DEFAULT b'0',
  `hotel_id` bigint(32) DEFAULT NULL COMMENT '公寓表主键ID',
  `gateway_id` varchar(32) COLLATE utf8_unicode_ci DEFAULT '' COMMENT '网关设备id',
  `status` int(2) NOT NULL DEFAULT '0' COMMENT '设备是否在线  1 在线  0 不在线',
  PRIMARY KEY (`id`),
  KEY `FK_61iomm8p53wtvnlnljxq3wems` (`hotel_id`)
) ENGINE=MyISAM AUTO_INCREMENT=12 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='网关信息表';

-- ----------------------------
-- Table structure for t_user_info
-- ----------------------------
DROP TABLE IF EXISTS `t_user_info`;
CREATE TABLE `t_user_info` (
  `uid` int(11) NOT NULL AUTO_INCREMENT,
  `user_name` varchar(255) NOT NULL DEFAULT '',
  `email` varchar(60) NOT NULL DEFAULT '',
  PRIMARY KEY (`uid`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=latin1;
