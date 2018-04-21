/*
Navicat MySQL Data Transfer

Source Server         : 192.168.191.131
Source Server Version : 50535
Source Host           : 192.168.191.131:3306
Source Database       : HOTEL

Target Server Type    : MYSQL
Target Server Version : 50535
File Encoding         : 65001

Date: 2018-04-21 15:18:53
*/

SET FOREIGN_KEY_CHECKS=0;

-- ----------------------------
-- Table structure for t_device_info
-- ----------------------------
DROP TABLE IF EXISTS `t_device_info`;
CREATE TABLE `t_device_info` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `device_id` varchar(32) NOT NULL,
  `status` tinyint(4) NOT NULL DEFAULT '0' COMMENT '0 不在线，1在线',
  `user_id` int(11) NOT NULL COMMENT '关联表t_user_info',
  `gw_id` int(11) NOT NULL COMMENT '关联表t_gateway_info',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- ----------------------------
-- Table structure for t_gateway_info
-- ----------------------------
DROP TABLE IF EXISTS `t_gateway_info`;
CREATE TABLE `t_gateway_info` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `device_id` varchar(32) NOT NULL,
  `status` tinyint(4) NOT NULL DEFAULT '0' COMMENT '0不在线，1在线',
  `user_id` int(11) NOT NULL COMMENT '关联表t_user_info',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- ----------------------------
-- Table structure for t_user_info
-- ----------------------------
DROP TABLE IF EXISTS `t_user_info`;
CREATE TABLE `t_user_info` (
  `id` int(11) NOT NULL,
  `user_account` varchar(32) CHARACTER SET utf8 COLLATE utf8_unicode_ci NOT NULL COMMENT '登录名(唯一)',
  `user_pwd` varchar(64) NOT NULL COMMENT '登录密码(md5加密)',
  `hotel_name` varchar(125) NOT NULL COMMENT '酒店名',
  `hotel_addr` varchar(255) DEFAULT '' COMMENT '酒店地址',
  `hotel_phone` varchar(64) DEFAULT '' COMMENT '酒店电话',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
