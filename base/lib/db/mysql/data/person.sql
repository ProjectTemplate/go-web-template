create table person
(
    id          int auto_increment primary key,
    name        varchar(50)  default ''  not null,
    age         int          default 0   not null,
    create_time int unsigned default '0' not null comment '创建时间'
);

INSERT INTO test.person (id, name, age, create_time) VALUES (1, '张三', 1, 0);
INSERT INTO test.person (id, name, age, create_time) VALUES (2, '李四', 2, 0);
INSERT INTO test.person (id, name, age, create_time) VALUES (3, '王五', 3, 0);
