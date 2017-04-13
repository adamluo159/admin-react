-- MySQL dump 10.13  Distrib 5.7.16, for Win64 (x86_64)
--
-- Host: localhost    Database: zonelog0
-- ------------------------------------------------------
-- Server version	5.7.16

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `log_charge`
--

DROP TABLE IF EXISTS `log_charge`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `log_charge` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '序号，自增',
  `zid` int(11) unsigned DEFAULT NULL COMMENT '区服ID',
  `cid` int(11) DEFAULT NULL COMMENT '渠道ID',
  `uid` bigint(20) unsigned DEFAULT NULL COMMENT '用户UID',
  `rid` bigint(20) unsigned DEFAULT NULL COMMENT '角色ID',
  `rname` varchar(100) DEFAULT '' COMMENT '角色名称',
  `orderno` varchar(50) DEFAULT '' COMMENT '订单号',
  `charge` float(10,2) unsigned DEFAULT NULL COMMENT '充值金额',
  `status` tinyint(3) unsigned DEFAULT '0' COMMENT '充值是否成功：0失败 1成功',
  `time` int(11) unsigned DEFAULT NULL COMMENT '充值时间',
  PRIMARY KEY (`id`),
  KEY `I_uid_rid` (`uid`,`rid`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `log_charge`
--

LOCK TABLES `log_charge` WRITE;
/*!40000 ALTER TABLE `log_charge` DISABLE KEYS */;
/*!40000 ALTER TABLE `log_charge` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `log_coins`
--

DROP TABLE IF EXISTS `log_coins`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `log_coins` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '序号，自增',
  `zid` int(11) unsigned DEFAULT NULL COMMENT '区服ID',
  `cid` int(11) DEFAULT NULL COMMENT '渠道ID',
  `uid` bigint(20) unsigned DEFAULT NULL COMMENT '用户UID',
  `rid` bigint(20) unsigned DEFAULT NULL COMMENT '角色ID',
  `rname` varchar(100) DEFAULT '' COMMENT '角色名称',
  `act` tinyint(3) unsigned DEFAULT '0' COMMENT '操作类型  0增加 1减少',
  `oldnum` bigint(20) unsigned DEFAULT '0' COMMENT '操作前数量',
  `actnum` bigint(20) unsigned DEFAULT '0' COMMENT '操作数量',
  `newnum` bigint(20) unsigned DEFAULT '0' COMMENT '操作后的数量',
  `type` smallint(8) unsigned DEFAULT '0' COMMENT '操作类型',
  `mark` varchar(200) DEFAULT '' COMMENT '备注信息',
  `time` int(11) unsigned DEFAULT NULL COMMENT '操作时间',
  PRIMARY KEY (`id`),
  KEY `I_uid_rid` (`uid`,`rid`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `log_coins`
--

LOCK TABLES `log_coins` WRITE;
/*!40000 ALTER TABLE `log_coins` DISABLE KEYS */;
/*!40000 ALTER TABLE `log_coins` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `log_email`
--

DROP TABLE IF EXISTS `log_email`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `log_email` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '序号，自增',
  `zid` int(11) unsigned DEFAULT NULL COMMENT '区服ID',
  `cid` int(11) DEFAULT NULL COMMENT '渠道ID',
  `uid` bigint(20) unsigned DEFAULT NULL COMMENT '用户UID',
  `rid` bigint(20) unsigned DEFAULT NULL COMMENT '角色ID',
  `rname` varchar(100) DEFAULT '' COMMENT '角色名称',
  `email_id` bigint(20) unsigned DEFAULT NULL COMMENT '邮件ID',
  `email_type` tinyint(3) unsigned DEFAULT '0' COMMENT '邮件发送类型：0发件  1收件',
  `email_title` varchar(200) DEFAULT '' COMMENT '邮件名称',
  `email_man` varchar(100) DEFAULT '' COMMENT '发/接件人',
  `email_files` varchar(200) DEFAULT '' COMMENT '附件名称',
  `email_num` int(11) unsigned DEFAULT '0' COMMENT '邮件数量',
  `time` int(11) unsigned DEFAULT NULL COMMENT '操作时间',
  PRIMARY KEY (`id`),
  KEY `I_uid_rid` (`uid`,`rid`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `log_email`
--

LOCK TABLES `log_email` WRITE;
/*!40000 ALTER TABLE `log_email` DISABLE KEYS */;
/*!40000 ALTER TABLE `log_email` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `log_jewel`
--

DROP TABLE IF EXISTS `log_jewel`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `log_jewel` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '序号，自增',
  `zid` int(11) unsigned DEFAULT NULL COMMENT '区服ID',
  `cid` int(11) DEFAULT NULL COMMENT '渠道ID',
  `uid` bigint(20) unsigned DEFAULT NULL COMMENT '用户UID',
  `rid` bigint(20) unsigned DEFAULT NULL COMMENT '角色ID',
  `rname` varchar(100) DEFAULT '' COMMENT '角色名称',
  `act` tinyint(3) unsigned DEFAULT '0' COMMENT '操作类型  0增加 1减少',
  `oldnum` bigint(20) unsigned DEFAULT '0' COMMENT '操作前数量',
  `actnum` bigint(20) unsigned DEFAULT '0' COMMENT '操作数量',
  `newnum` bigint(20) unsigned DEFAULT '0' COMMENT '操作后的数量',
  `type` smallint(8) unsigned DEFAULT '0' COMMENT '操作类型',
  `mark` varchar(200) DEFAULT '' COMMENT '备注信息',
  `time` int(11) unsigned DEFAULT NULL COMMENT '操作时间',
  PRIMARY KEY (`id`),
  KEY `I_uid_rid` (`uid`,`rid`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `log_jewel`
--

LOCK TABLES `log_jewel` WRITE;
/*!40000 ALTER TABLE `log_jewel` DISABLE KEYS */;
/*!40000 ALTER TABLE `log_jewel` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `log_prop`
--

DROP TABLE IF EXISTS `log_prop`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `log_prop` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '序号，自增',
  `zid` int(11) unsigned DEFAULT NULL COMMENT '区服ID',
  `cid` int(11) DEFAULT NULL COMMENT '渠道ID',
  `uid` bigint(20) unsigned DEFAULT NULL COMMENT '用户UID',
  `rid` bigint(20) unsigned DEFAULT NULL COMMENT '角色ID',
  `rname` varchar(100) DEFAULT '' COMMENT '角色名称',
  `prop_id` int(11) unsigned DEFAULT NULL COMMENT '道具ID',
  `prop_name` varchar(100) DEFAULT '' COMMENT '道具名称',
  `act` tinyint(3) unsigned DEFAULT '0' COMMENT '操作类型  0增加 1减少',
  `oldnum` bigint(20) unsigned DEFAULT '0' COMMENT '操作前数量',
  `actnum` bigint(20) unsigned DEFAULT '0' COMMENT '操作数量',
  `newnum` bigint(20) unsigned DEFAULT '0' COMMENT '操作后的数量',
  `type` smallint(8) unsigned DEFAULT '0' COMMENT '操作类型',
  `mark` varchar(200) DEFAULT '' COMMENT '备注信息',
  `time` int(11) unsigned DEFAULT NULL COMMENT '操作时间',
  PRIMARY KEY (`id`),
  KEY `I_uid_rid` (`uid`,`rid`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `log_prop`
--

LOCK TABLES `log_prop` WRITE;
/*!40000 ALTER TABLE `log_prop` DISABLE KEYS */;
/*!40000 ALTER TABLE `log_prop` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `log_sidejob`
--

DROP TABLE IF EXISTS `log_sidejob`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `log_sidejob` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '序号，自增',
  `zid` int(11) unsigned DEFAULT NULL COMMENT '区服ID',
  `cid` int(11) DEFAULT NULL COMMENT '渠道ID',
  `uid` bigint(20) unsigned DEFAULT NULL COMMENT '用户UID',
  `rid` bigint(20) unsigned DEFAULT NULL COMMENT '角色ID',
  `rname` varchar(100) DEFAULT '' COMMENT '角色名称',
  `bpropid` int(11) unsigned DEFAULT NULL COMMENT '变化前物品ID',
  `apropid` int(11) unsigned DEFAULT NULL COMMENT '变化后物品ID',
  `plevel` int(11) unsigned DEFAULT NULL COMMENT '玩家等级',
  `status` tinyint(3) unsigned DEFAULT '0' COMMENT '是否操作成功:0失败 1成功',
  `type` smallint(8) unsigned DEFAULT '0' COMMENT '操作类型',
  `mark` varchar(200) DEFAULT '' COMMENT '备注信息',
  `time` int(11) unsigned DEFAULT NULL COMMENT '操作时间',
  PRIMARY KEY (`id`),
  KEY `I_uid_rid` (`uid`,`rid`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `log_sidejob`
--

LOCK TABLES `log_sidejob` WRITE;
/*!40000 ALTER TABLE `log_sidejob` DISABLE KEYS */;
/*!40000 ALTER TABLE `log_sidejob` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2017-04-05 20:08:13
