CREATE TABLE IF NOT EXISTS fleet
(
    fleet_id varchar(10) primary key
);
CREATE TABLE IF NOT EXISTS driver
(
    driver_id varchar(10) primary key,
    name      varchar(30) NOT NULL,
    year      int         NOT NULL,
    sex       int         NOT NULL,
    fleet_id  varchar(10) NOT NULL,
    position  int         NOT NULL,
    passwd    varchar(20) NOT NULL,
    FOREIGN KEY (fleet_id) REFERENCES fleet (fleet_id)
    CHECK (sex in (0,1))
);


CREATE TABLE IF NOT EXISTS line
(
    line_id  varchar(10) PRIMARY KEY,
    fleet_id varchar(10) NOT NULL,
    FOREIGN KEY (fleet_id) REFERENCES fleet (fleet_id)
);

CREATE TABLE IF NOT EXISTS bus
(
    bus_id  varchar(20) PRIMARY KEY,
    line_id varchar(10) NOT NULL,
    FOREIGN KEY (line_id) REFERENCES line (line_id)
);

CREATE TABLE IF NOT EXISTS stop
(
    stop_id varchar(15) primary key
);

CREATE TABLE IF NOT EXISTS line_stop
(
    stop_id    varchar(15),
    line_id    varchar(10),
    stop_order int NOT NULL,
    FOREIGN KEY (stop_id) REFERENCES stop (stop_id),
    FOREIGN KEY (line_id) REFERENCES line (line_id),
    primary key (stop_id, line_id)
);

CREATE TABLE IF NOT EXISTS driver_line
(
    driver_id varchar(10) primary key,
    line_id   varchar(10) NOT NULL,
    position  int         NOT NULL,
    FOREIGN KEY (line_id) REFERENCES line (line_id)
);


CREATE TABLE IF NOT EXISTS violation_type
(
    violation_type_id varchar(20) primary key
);

CREATE TABLE IF NOT EXISTS violation_record
(
    violation_id      int primary key AUTO_INCREMENT,
    violation_type_id varchar(20) NOT NULL,
    time              int         NOT NULL,
    driver_id         varchar(10) NOT NULL,
    bus_id            varchar(20) NOT NULL,
    fleet_id          varchar(20) NOT NULL,
    stop_id           varchar(15) NOT NULL,
    line_id           varchar(10) NOT NULL,
    FOREIGN KEY (violation_type_id) references violation_type (violation_type_id),
    FOREIGN KEY (driver_id) references driver (driver_id),
    FOREIGN KEY (line_id) references line (line_id),
    FOREIGN KEY (bus_id) references bus (bus_id),
    FOREIGN KEY (fleet_id) references fleet (fleet_id),
    FOREIGN KEY (stop_id) references stop (stop_id)
);

INSERT INTO fleet (fleet_id) values ('0');
INSERT INTO driver (driver_id,name,year,sex,fleet_id,position,passwd) values ('root','Ruitian Zhong',2003,1,'0',0,'Your password');
Insert into violation_type values ('闯红灯'),('超速'),('不礼让行人'),('不按规定按喇叭'),('不按车道行驶'),('在车厢没有关好时行车'),('接打电话');
CREATE INDEX violation_record_index on violation_record(time);
CREATE INDEX violation_record_index_driver_id on violation_record(driver_id);