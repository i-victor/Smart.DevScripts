----
-- phpLiteAdmin database dump (http://www.phpliteadmin.org/)
-- phpLiteAdmin version: 1.9.7.1
-- Exported: 4:38pm on October 27, 2017 (UTC)
-- database file: ./db/smtpd.sqlite3

--- HowTo Setup
-- cd /root
-- openssl req -new -x509 -nodes -newkey rsa:4096 -keyout server.key -out server.crt -days 1095
-- # move certs to /etc/ssl{/private}
-- groupadd -g 2100 vmail
-- useradd -g vmail -u 2100 -s /sbin/nologin -d /var/vmail -m vmail
-- test smtpd config: smtpd -n
---

----
BEGIN TRANSACTION;

----
-- Table structure for valias
----
CREATE TABLE valias (
alias_id INT(8) PRIMARY KEY NOT NULL,
addr varchar(42) NOT NULL,
alias varchar(42) NOT NULL
);

----
-- Data dump for valias, a total of 5 rows
----
INSERT INTO "valias" ("alias_id","addr","alias") VALUES ('1','root@unix.uxm','admin__at__unix__dot__uxm');
INSERT INTO "valias" ("alias_id","addr","alias") VALUES ('2','postmaster@unix.uxm','admin__at__unix__dot__uxm');
INSERT INTO "valias" ("alias_id","addr","alias") VALUES ('3','admin@unix.uxm','admin__at__unix__dot__uxm');
INSERT INTO "valias" ("alias_id","addr","alias") VALUES ('4','unixman@unix.uxm','admin__at__unix__dot__uxm');
INSERT INTO "valias" ("alias_id","addr","alias") VALUES ('5','webmaster@unix.uxm','admin__at__unix__dot__uxm');

----
-- Table structure for vdomains
----
CREATE TABLE vdomains (
domain_id INT(11) PRIMARY KEY NOT NULL,
domain VARCHAR(42) NOT NULL
);

----
-- Data dump for vdomains, a total of 1 rows
----
INSERT INTO "vdomains" ("domain_id","domain") VALUES ('1','unix.uxm');

----
-- Table structure for users
----
CREATE TABLE users (
user_id INT(8) PRIMARY KEY NOT NULL,
username TEXT NOT NULL,
domain TEXT NOT NULL,
mailbox TEXT NOT NULL,
password TEXT NULL,
home TEXT NOT NULL,
uid INTEGER NOT NULL,
gid INTEGER NOT NULL
);

----
-- Data dump for users, a total of 0 rows
----

----
-- Table structure for userinfo
----
CREATE TABLE 'userinfo' (
user_id INT(11) PRIMARY KEY NOT NULL,
user VARCHAR(42) UNIQUE NOT NULL,
password VARCHAR(255) NOT NULL,
uid INT(42) DEFAULT 2100,
gid INT(42) DEFAULT 2100,
maildir VARCHAR(255) DEFAULT '/var/empty/vmail'
);

----
-- Data dump for userinfo, a total of 1 rows ; passw: admin11 ; smtpctl encrypt my-personnal-password
----
INSERT INTO "userinfo" ("user_id","user","password","uid","gid","maildir") VALUES ('1','admin__at__unix__dot__uxm','$2b$10$PhiClsq3ct9./udaU8Yyx.pXr6cjBdr4kSsQM0JlANcKwZLROFH5y','2100','2100','/var/empty/vmail');

----
-- structure for index sqlite_autoindex_valias_1 on table valias
----
;

----
-- structure for index sqlite_autoindex_vdomains_1 on table vdomains
----
;

----
-- structure for index sqlite_autoindex_users_1 on table users
----
;

----
-- structure for index sqlite_autoindex_userinfo_1 on table userinfo
----
;

----
-- structure for index sqlite_autoindex_userinfo_2 on table userinfo
----
;
COMMIT;
