drop database if exists camp;
create database camp;
use camp;

# member
create table member
(
    member_id int(64) auto_increment primary key,
    member_name varchar(32) not null,
    member_nickname varchar(32) not null,
    member_password varchar(32) not null,
    member_type int not null,
    is_deleted boolean default false not null
);

#member索引
create index idx_member_name
    on member (member_name);

create index idx_is_deleted_id
    on member (is_deleted,member_id);

# course
create table course
(
    course_id int auto_increment,
    course_name varchar(32) not null,
    course_available int not null,
    constraint course_pk
        primary key (course_id)
);


create table teacher_schedule
(
    schedule_id int auto_increment,
    teacher_id int not null,
    course_id int not null,
    constraint teacher_schedule_pk
        primary key (schedule_id)
);

create index idx_course_id
    on teacher_schedule (course_id);

create index idx_teacher_id
    on teacher_schedule (teacher_id);

create table student_schedule
(
    schedule_id int auto_increment,
    student_id int not null,
    course_id int not null,
    constraint student_schedule_pk
        primary key (schedule_id)
);

create index idx_course_id
    on student_schedule (course_id);

create index idx_student_id
    on student_schedule (student_id);


# 系统内置管理员
INSERT INTO camp.member (member_name, member_nickname, member_password, member_type) VALUES('JudgeAdmin', 'JudgeAdmin', 'JudgePassword2022',1);
