# Vul Warnings

## News

- [x] Aliyun
- [x] Cert360
- [x] TencentTI
- [x] GitHub CVE Search

## Pusher

- [x] Feishu 飞书
- [x] Qiwei 企业微信
- [x] Dingding 钉钉

## Usage

### Step 1 : Create MySQL/Mariadb Database

```sql
create database `you_database_name` default character set utf8mb4 collate utf8mb4_unicode_ci;
create user 'username'@'%' identified by 'you_password';
-- create user 'username'@'127.0.0.1' identified by 'you_password';
-- create user 'username'@'localhost' identified by 'you_password';
grant all on `you_database_name`.* to username;
flush privileges;
```

### Step 2 : config.yaml

Modify the `config.yaml` in the path of binary file.

1. Print template config.yaml by run command: `./vulwarning config`
2. And you can save it. `./vulwarning config > config.yaml`
3. Write anything you want.

**pusher config:** Just push message to which you set key

**Example:** `example.config.yaml`

[CRON Expression Format](#cron-expression-format)

### Step 3 : Init Database

Init Database and Run First Crawl without pushing message.

`./vulwarning initdb`

### Step 4 : Service

Install the vulwarning service. `./vulwarning install`

You can found other usage about service `./vulwarning help`

Also, you can run `./vulwarning` without any argv

Have a good time~ `./vulwarning start`

### Debug Mode

You can open Debug Mode by `config.yaml` or setenv DEBUG=1

### CRON Expression Format

https://pkg.go.dev/github.com/robfig/cron/v3#hdr-CRON_Expression_Format

| Field name   | Mandatory? | Allowed values  | Allowed special characters |
| ------------ | ---------- | --------------- | -------------------------- |
| Minutes      | Yes        | 0-59            | * / , -                    |
| Hours        | Yes        | 0-23            | * / , -                    |
| Day of month | Yes        | 1-31            | * / , - ?                  |
| Month        | Yes        | 1-12 or JAN-DEC | * / , -                    |
| Day of week  | Yes        | 0-6 or SUN-SAT  | * / , - ?                  |

## LICENSE

[WTFPL](LICENSE)