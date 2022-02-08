create database camp;
use camp;

# member
create table member
(
    member_id int auto_increment,
    member_name varchar(32) not null,
    member_nickname varchar(32) not null,
    member_password varchar(32) not null,
    member_type int not null,
    is_deleted int default 0 not null,
    constraint member_pk
        primary key (member_id)
);

# course
create table course
(
    course_id int auto_increment,
    course_name varchar(32) not null,
    course_available int not null,
    constraint course_pk
        primary key (course_id)
);

# course_schedule
create table course_schedule
(
    schedule_id int auto_increment,
    course_id int not null,
    member_id int not null,
    member_type int not null,
    constraint course_schedule_pk
        primary key (schedule_id)
);

# 系统内置管理员
INSERT INTO camp.member (member_name, member_nickname, member_password, member_type) VALUES('JudgeAdmin', 'JudgeAdmin', 'JudgePassword2022',1);

#member索引
create index idx_member_name
    on member (member_name);

create index idx_is_deleted_id
    on member (is_deleted,member_id);