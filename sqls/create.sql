create database camp;

# 成员登录表
create table users
(
    name varchar(20) primary key,
    password varchar(20)
)engine=innoDB;

# 成员信息表
create table userinfo
(
    userid   int(64) primary key auto_imcrement,
    username varchar(20),
    nickname varchar(20),
    usertype int
)engine=innoDB;

# 系统内置管理员