drop procedure if exists init_member;
delimiter //

# ---------- 创建存储过程初始化成员表 ------------ #
create procedure init_member(num int)
begin
    declare musername varchar(32);
    declare mnickname varchar(32);
    declare mpassword varchar(32);
    declare musertype int;
    declare m_deleted boolean;
    declare i int;
    declare lower char(26) default 'abcdefghijklmnopqrstuvwxyz';
    declare j int;

    set i = 1;

    while(i <= num)
    do
        set j = 1;
        set musername = '';
        while(j <= 10)
            do
                set musername = concat(musername, substring(lower, floor(1+rand()*26), 1));
                set j = j + 1;
            end while;

        set mnickname = musername;
        set mpassword = concat('P', concat(musername, '1'));
        set musertype = floor(1+rand()*3);
        set m_deleted = floor(rand()*2);
        if musername not in (select member_name from member) then
            insert into member (member_name, member_nickname, member_password, member_type, is_deleted) values (musername, mnickname, mpassword, musertype, m_deleted);
            set i = i + 1;
        end if;
    end while;
end //

delimiter ;