SET NAMES utf8mb4;

-- ----------------------------
-- Table structure for p_article
-- ----------------------------
DROP TABLE IF EXISTS `p_article`;
CREATE TABLE `p_article`  (
                              `id` int NOT NULL AUTO_INCREMENT,
                              `title` varchar(256) CHARACTER SET utf8mb4  NULL DEFAULT NULL COMMENT '文章标题',
                              `content` longtext CHARACTER SET utf8mb4  NOT NULL COMMENT '文章内容',
                              `parseContent` longtext CHARACTER SET utf8mb4  NOT NULL COMMENT '解析后的文章内容',
                              `contentModel` varchar(32) CHARACTER SET utf8mb4  NULL DEFAULT NULL COMMENT '文章内容类型:html/markdown',
                              `type` varchar(32) CHARACTER SET utf8mb4  NULL DEFAULT NULL COMMENT '文章类型:article文章,page页面',
                              `summary` varchar(1024) CHARACTER SET utf8mb4  NULL DEFAULT NULL COMMENT '文章摘要',
                              `metaKeywords` varchar(512) CHARACTER SET utf8mb4  NULL DEFAULT NULL COMMENT 'SEO关键字',
                              `metaDescription` varchar(512) CHARACTER SET utf8mb4  NULL DEFAULT NULL COMMENT 'SEO描述',
                              `thumbnail` varchar(256) CHARACTER SET utf8mb4  NULL DEFAULT NULL COMMENT '缩略图',
                              `slug` varchar(128) CHARACTER SET utf8mb4  NULL DEFAULT NULL COMMENT 'slug',
                              `isTop` int NULL DEFAULT 0 COMMENT '是否置顶0:否,1:是',
                              `status` int NULL DEFAULT 0 COMMENT '状态0:已发布,1:草稿',
                              `viewCount` int NULL DEFAULT 0 COMMENT '访问量',
                              `greatCount` int NULL DEFAULT 0 COMMENT '访问量',
                              `isComment` int NULL DEFAULT 1 COMMENT '是否允许评论0:否,1是',
                              `flag` varchar(256) CHARACTER SET utf8mb4  NULL DEFAULT NULL COMMENT '标识',
                              `template` varchar(256) CHARACTER SET utf8mb4  NULL DEFAULT NULL COMMENT '模板',
                              `createTime` datetime NOT NULL COMMENT '创建时间',
                              `updateTime` datetime NULL DEFAULT NULL COMMENT '更新时间',
                              `createUserId` int NOT NULL COMMENT '添加人',
                              `updateUserId` int NULL DEFAULT NULL COMMENT '更新人',
                              `visibility` int NOT NULL DEFAULT 0 COMMENT '是否可见, 0是, 1否',
                              PRIMARY KEY (`id`) USING BTREE,
                              INDEX `slug`(`slug`) USING BTREE,
                              INDEX `isTop`(`isTop`) USING BTREE,
                              INDEX `type`(`type`) USING BTREE,
                              INDEX `status`(`status`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 187 CHARACTER SET = utf8mb4  ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of p_article
-- ----------------------------

-- ----------------------------
-- Table structure for p_article_category
-- ----------------------------
DROP TABLE IF EXISTS `p_article_category`;
CREATE TABLE `p_article_category`  (
                                       `id` int NOT NULL AUTO_INCREMENT COMMENT '主键',
                                       `articleId` int NOT NULL COMMENT '文章id',
                                       `categoryId` int NOT NULL COMMENT '分类id',
                                       `createTime` datetime NOT NULL COMMENT '创建时间',
                                       `updateTime` datetime NULL DEFAULT NULL COMMENT '更新时间',
                                       `createUserId` int NULL DEFAULT NULL COMMENT '添加人',
                                       `updateUserId` int NULL DEFAULT NULL COMMENT '更新人',
                                       PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 56 CHARACTER SET = utf8mb4  ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of p_article_category
-- ----------------------------

-- ----------------------------
-- Table structure for p_article_tag
-- ----------------------------
DROP TABLE IF EXISTS `p_article_tag`;
CREATE TABLE `p_article_tag`  (
                                  `id` int NOT NULL AUTO_INCREMENT,
                                  `articleId` int NOT NULL COMMENT '文章id',
                                  `tagId` int NOT NULL COMMENT '标签id',
                                  `createTime` datetime NULL DEFAULT NULL COMMENT '创建时间',
                                  `updateTime` datetime NULL DEFAULT NULL COMMENT '更新时间',
                                  `createUserId` int NULL DEFAULT NULL COMMENT '添加人',
                                  `updateUserId` int NULL DEFAULT NULL COMMENT '更新人',
                                  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 73 CHARACTER SET = utf8mb4  ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of p_article_tag
-- ----------------------------

-- ----------------------------
-- Table structure for p_attach
-- ----------------------------
DROP TABLE IF EXISTS `p_attach`;
CREATE TABLE `p_attach`  (
                             `id` int NOT NULL AUTO_INCREMENT COMMENT '主键',
                             `name` varchar(256) CHARACTER SET utf8mb4  NOT NULL COMMENT '附件名',
                             `remark` varchar(512) CHARACTER SET utf8mb4  NULL DEFAULT NULL COMMENT '附件描述',
                             `path` varchar(512) CHARACTER SET utf8mb4  NOT NULL COMMENT '附件路径',
                             `flag` varchar(256) CHARACTER SET utf8mb4  NULL DEFAULT NULL COMMENT '标识',
                             `mineType` varchar(128) CHARACTER SET utf8mb4  NULL DEFAULT NULL COMMENT '文件类型mineType',
                             `type` varchar(32) CHARACTER SET utf8mb4  NULL DEFAULT NULL COMMENT '文件类型',
                             `createTime` datetime NOT NULL COMMENT '创建时间',
                             `updateTime` datetime NULL DEFAULT NULL COMMENT '更新时间',
                             `configId` int NOT NULL COMMENT '存储策略id',
                             `url` varchar(512) CHARACTER SET utf8mb4  NOT NULL COMMENT '访问路径',
                             `attachGroup` varchar(256) CHARACTER SET utf8mb4  NOT NULL DEFAULT 'default' COMMENT '附件分组',
                             `storage` int NOT NULL COMMENT '存储器类型',
                             `createUserId` int NULL DEFAULT NULL COMMENT '添加人',
                             `updateUserId` int NULL DEFAULT NULL COMMENT '更新人',
                             PRIMARY KEY (`id`) USING BTREE,
                             INDEX `type`(`type`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 250 CHARACTER SET = utf8mb4  ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of p_attach
-- ----------------------------

-- ----------------------------
-- Table structure for p_category
-- ----------------------------
DROP TABLE IF EXISTS `p_category`;
CREATE TABLE `p_category`  (
                               `id` int NOT NULL AUTO_INCREMENT COMMENT '主键',
                               `name` varchar(256) CHARACTER SET utf8mb4  NOT NULL COMMENT '分类名',
                               `pid` int NOT NULL DEFAULT -1 COMMENT '父级id',
                               `desc` varchar(512) CHARACTER SET utf8mb4  NULL DEFAULT NULL COMMENT '描述',
                               `metaKeywords` varchar(256) CHARACTER SET utf8mb4  NULL DEFAULT NULL COMMENT 'SEO关键字',
                               `thumbnail` varchar(256) CHARACTER SET utf8mb4  NULL DEFAULT NULL COMMENT '封面图',
                               `slug` varchar(128) CHARACTER SET utf8mb4  NULL DEFAULT NULL COMMENT 'slug',
                               `metaDescription` varchar(256) CHARACTER SET utf8mb4  NULL DEFAULT NULL COMMENT 'SEO描述内容',
                               `status` int NOT NULL DEFAULT 0 COMMENT '状态0:正常,1禁用',
                               `createTime` datetime NOT NULL COMMENT '创建时间',
                               `updateTime` datetime NULL DEFAULT NULL COMMENT '更新时间',
                               `createUserId` int NULL DEFAULT NULL COMMENT '添加人',
                               `updateUserId` int NULL DEFAULT NULL COMMENT '更新人',
                               PRIMARY KEY (`id`) USING BTREE,
                               INDEX `status`(`status`) USING BTREE,
                               INDEX `slug`(`slug`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 11 CHARACTER SET = utf8mb4  COMMENT = '分类表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of p_category
-- ----------------------------
-- ----------------------------
-- Table structure for p_comment
-- ----------------------------
DROP TABLE IF EXISTS `p_comment`;
CREATE TABLE `p_comment`  (
                              `id` int NOT NULL AUTO_INCREMENT COMMENT '主键',
                              `articleId` int NOT NULL COMMENT '文章id',
                              `pid` int NULL DEFAULT -1 COMMENT '父级id',
                              `topPid` int NULL DEFAULT -1 COMMENT '顶层父级id',
                              `userId` int NULL DEFAULT NULL COMMENT '用户iD',
                              `content` varchar(2048) CHARACTER SET utf8mb4  NULL DEFAULT NULL COMMENT '评论内容',
                              `status` int NULL DEFAULT 0 COMMENT '状态:0正常,1:待审核',
                              `avatar` varchar(256) CHARACTER SET utf8mb4  NULL DEFAULT NULL COMMENT '头像',
                              `website` varchar(256) CHARACTER SET utf8mb4  NULL DEFAULT NULL COMMENT '网站地址',
                              `email` varchar(256) CHARACTER SET utf8mb4  NULL DEFAULT NULL COMMENT '邮箱',
                              `userName` varchar(256) CHARACTER SET utf8mb4  NULL DEFAULT NULL COMMENT '评论人',
                              `ip` varchar(256) CHARACTER SET utf8mb4  NULL DEFAULT NULL COMMENT 'ip',
                              `device` varchar(256) CHARACTER SET utf8mb4  NULL DEFAULT NULL COMMENT '设备类型',
                              `createTime` datetime NOT NULL COMMENT '创建时间',
                              `updateTime` datetime NULL DEFAULT NULL COMMENT '更新时间',
                              `createUserId` int NULL DEFAULT NULL COMMENT '添加人',
                              `updateUserId` int NULL DEFAULT NULL COMMENT '更新人',
                              PRIMARY KEY (`id`) USING BTREE,
                              INDEX `articleId`(`articleId`) USING BTREE,
                              INDEX `status`(`status`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 89 CHARACTER SET = utf8mb4  ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of p_comment
-- ----------------------------

-- ----------------------------
-- Table structure for p_link
-- ----------------------------
DROP TABLE IF EXISTS `p_link`;
CREATE TABLE `p_link`  (
                           `id` int NOT NULL AUTO_INCREMENT COMMENT '主键',
                           `name` varchar(256) CHARACTER SET utf8mb4  NOT NULL COMMENT '网站名',
                           `logo` varchar(256) CHARACTER SET utf8mb4  NULL DEFAULT NULL COMMENT '网站logo',
                           `desc` varchar(512) CHARACTER SET utf8mb4  NULL DEFAULT NULL COMMENT '网站描述',
                           `address` varchar(256) CHARACTER SET utf8mb4  NOT NULL COMMENT '网站地址',
                           `createTime` datetime NOT NULL COMMENT '创建时间',
                           `updateTime` datetime NULL DEFAULT NULL COMMENT '更新时间',
                           `createUserId` int NULL DEFAULT NULL COMMENT '添加人',
                           `updateUserId` int NULL DEFAULT NULL COMMENT '更新人',
                           PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 11 CHARACTER SET = utf8mb4  ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of p_link
-- ----------------------------

-- ----------------------------
-- Table structure for p_tag
-- ----------------------------
DROP TABLE IF EXISTS `p_tag`;
CREATE TABLE `p_tag`  (
                          `id` int NOT NULL AUTO_INCREMENT COMMENT '主键',
                          `name` varchar(256) CHARACTER SET utf8mb4  NOT NULL COMMENT '标签名',
                          `color` varchar(128) CHARACTER SET utf8mb4  NULL DEFAULT NULL COMMENT '颜色',
                          `thumbnail` varchar(256) CHARACTER SET utf8mb4  NULL DEFAULT NULL COMMENT '缩略图',
                          `slug` varchar(128) CHARACTER SET utf8mb4  NULL DEFAULT NULL COMMENT 'slug',
                          `createUserId` int NULL DEFAULT NULL COMMENT '添加人',
                          `updateUserId` int NULL DEFAULT NULL COMMENT '更新人',
                          `createTime` datetime NOT NULL COMMENT '创建时间',
                          `updateTime` datetime NULL DEFAULT NULL COMMENT '修改时间',
                          PRIMARY KEY (`id`) USING BTREE,
                          INDEX `slug`(`slug`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 27 CHARACTER SET = utf8mb4  ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of p_tag
-- ----------------------------

-- ----------------------------
-- Table structure for p_user
-- ----------------------------
DROP TABLE IF EXISTS `p_user`;
CREATE TABLE `p_user`  (
                           `id` int NOT NULL AUTO_INCREMENT COMMENT '主键',
                           `account` varchar(32) CHARACTER SET utf8mb4  NOT NULL COMMENT '账户',
                           `userName` varchar(32) CHARACTER SET utf8mb4  NOT NULL COMMENT '账户名',
                           `password` varchar(32) CHARACTER SET utf8mb4  NOT NULL COMMENT '密码',
                           `salt` varchar(32) CHARACTER SET utf8mb4  NOT NULL COMMENT '盐值',
                           `status` int NOT NULL DEFAULT 0 COMMENT '状态:0正常,1禁用',
                           `avatar` varchar(512) CHARACTER SET utf8mb4  NULL DEFAULT NULL COMMENT '头像',
                           `email` varchar(128) CHARACTER SET utf8mb4  NULL DEFAULT NULL COMMENT '邮箱',
                           `website` varchar(256) CHARACTER SET utf8mb4  NULL DEFAULT NULL COMMENT '网站地址',
                           `createTime` datetime NOT NULL COMMENT '创建时间',
                           `updateTime` datetime NULL DEFAULT NULL COMMENT '更新时间',
                           `createUserId` int NULL DEFAULT NULL COMMENT '添加人',
                           `updateUserId` int NULL DEFAULT NULL COMMENT '更新人',
                           `remark` varchar(500) CHARACTER SET utf8mb4  NULL DEFAULT NULL COMMENT '备注',
                           `mobile` varchar(11) CHARACTER SET utf8mb4  NULL DEFAULT '' COMMENT '手机号码',
                           `sex` tinyint NULL DEFAULT NULL COMMENT '用户性别',
                           `loginIp` varchar(50) CHARACTER SET utf8mb4  NULL DEFAULT '' COMMENT '最后登录IP',
                           `loginDate` datetime NULL DEFAULT NULL COMMENT '最后登录时间',
                           PRIMARY KEY (`id`) USING BTREE,
                           INDEX `account`(`account`) USING BTREE,
                           INDEX `status`(`status`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 56 CHARACTER SET = utf8mb4  ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of p_user
-- ----------------------------
