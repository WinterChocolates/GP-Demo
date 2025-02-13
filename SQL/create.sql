-- 创建数据库
CREATE DATABASE IF NOT EXISTS hrms
    CHARACTER SET utf8mb4
    COLLATE utf8mb4_unicode_ci;
USE hrms;

-- 用户表
CREATE TABLE users (
                       user_id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
                       username VARCHAR(50) UNIQUE NOT NULL,
                       password_hash CHAR(60) NOT NULL,
                       user_type ENUM('admin', 'employee', 'applicant') NOT NULL,
                       department VARCHAR(50),
                       position VARCHAR(50),
                       hire_date DATE,
                       salary_base DECIMAL(10,2),
                       education TEXT,
                       work_experience TEXT,
                       skills TEXT,
                       resume_path VARCHAR(255),
                       is_active BOOLEAN DEFAULT TRUE,
                       created_at DATETIME DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 插入用户数据（密码应为BCrypt哈希，示例使用明文仅用于演示）
INSERT INTO users (username, password_hash, user_type, department, position, hire_date, salary_base)
VALUES
    ('admin01', '123456', 'admin', NULL, NULL, NULL, NULL),
    ('emp001', '123456', 'employee', '技术部', 'Go开发工程师', '2023-01-15', 15000.00),
    ('emp002', '123456', 'employee', '人力资源部', 'HR主管', '2022-05-20', 18000.00),
    ('app2024', '123456', 'applicant', NULL, NULL, NULL, NULL);

-- 更新应聘者信息
UPDATE users
SET education = '南京大学 计算机科学 本科',
    work_experience = '腾讯 后端开发 2年',
    skills = 'Go, MySQL, Redis',
    resume_path = '/resumes/app2024.pdf'
WHERE username = 'app2024';

-- 职位表
CREATE TABLE jobs (
                      job_id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
                      title VARCHAR(100) NOT NULL,
                      description TEXT NOT NULL,
                      requirements TEXT NOT NULL,
                      salary_range VARCHAR(50),
                      post_date DATETIME DEFAULT CURRENT_TIMESTAMP,
                      expiration_date DATE,
                      status ENUM('open', 'closed') DEFAULT 'open'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

INSERT INTO jobs (title, description, requirements, salary_range, expiration_date)
VALUES
    ('高级Go工程师', '负责分布式系统开发', '3年以上Go开发经验', '25k-40k', '2024-06-30'),
    ('HR实习生', '协助招聘工作', '人力资源专业优先', '4k-6k', '2024-05-31');

-- 申请表
CREATE TABLE applications (
                              application_id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
                              user_id INT UNSIGNED NOT NULL,
                              job_id INT UNSIGNED NOT NULL,
                              apply_date DATETIME DEFAULT CURRENT_TIMESTAMP,
                              status ENUM('pending', 'interviewed', 'hired', 'rejected') DEFAULT 'pending',
                              FOREIGN KEY (user_id) REFERENCES users(user_id),
                              FOREIGN KEY (job_id) REFERENCES jobs(job_id),
                              INDEX idx_user_id (user_id),
                              INDEX idx_job_id (job_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

INSERT INTO applications (user_id, job_id, status)
VALUES
    (4, 1, 'pending'),
    (4, 2, 'interviewed');

-- 考勤表（使用生成列实现日期索引）
CREATE TABLE attendances (
                             attendance_id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
                             user_id INT UNSIGNED NOT NULL,
                             clock_in DATETIME NOT NULL,
                             clock_out DATETIME,
                             status ENUM('normal', 'late', 'early_leave') DEFAULT 'normal',
    -- 添加虚拟列用于日期索引
                             date_clock_in DATE AS (DATE(clock_in)) STORED,
                             INDEX idx_user_date (user_id, date_clock_in),
                             FOREIGN KEY (user_id) REFERENCES users(user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

INSERT INTO attendances (user_id, clock_in, clock_out, status)
VALUES
    (2, '2024-05-01 08:58:00', '2024-05-01 18:02:00', 'normal'),
    (2, '2024-05-02 09:05:00', '2024-05-02 17:55:00', 'late');

-- 薪资表
CREATE TABLE salaries (
                          salary_id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
                          user_id INT UNSIGNED NOT NULL,
                          month CHAR(6) NOT NULL COMMENT '格式: YYYYMM',
                          base DECIMAL(10,2) NOT NULL,
                          bonus DECIMAL(10,2) DEFAULT 0.00,
                          deductions DECIMAL(10,2) DEFAULT 0.00,
                          payment_date DATE,
                          UNIQUE KEY uniq_user_month (user_id, month),
                          FOREIGN KEY (user_id) REFERENCES users(user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

INSERT INTO salaries (user_id, month, base, bonus, deductions, payment_date)
VALUES
    (2, '202404', 15000.00, 3000.00, 2450.00, '2024-05-05'),
    (3, '202404', 18000.00, 2500.00, 3150.00, '2024-05-05');

-- 培训表
CREATE TABLE trainings (
                           training_id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
                           title VARCHAR(100) NOT NULL,
                           description TEXT,
                           start_time DATETIME NOT NULL,
                           end_time DATETIME NOT NULL,
                           location VARCHAR(100),
                           capacity INT UNSIGNED
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

INSERT INTO trainings (title, description, start_time, end_time, location, capacity)
VALUES
    ('Go语言性能优化', '深入讲解Go性能调优', '2024-06-01 14:00:00', '2024-06-01 17:00:00', '线上会议室', 50),
    ('新员工入职培训', '公司制度与文化培训', '2024-05-15 09:00:00', '2024-05-16 17:00:00', '公司会议室', 30);

-- 培训记录表
CREATE TABLE training_records (
                                  record_id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
                                  user_id INT UNSIGNED NOT NULL,
                                  training_id INT UNSIGNED NOT NULL,
                                  enroll_time DATETIME DEFAULT CURRENT_TIMESTAMP,
                                  status ENUM('registered', 'completed', 'canceled') DEFAULT 'registered',
                                  score TINYINT UNSIGNED,
                                  FOREIGN KEY (user_id) REFERENCES users(user_id),
                                  FOREIGN KEY (training_id) REFERENCES trainings(training_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

INSERT INTO training_records (user_id, training_id, status)
VALUES
    (2, 1, 'registered'),
    (3, 2, 'completed');

-- 公告表
CREATE TABLE notices (
                         notice_id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
                         title VARCHAR(200) NOT NULL,
                         content TEXT NOT NULL,
                         publish_time DATETIME DEFAULT CURRENT_TIMESTAMP,
                         expire_time DATETIME,
                         target_type ENUM('all', 'department') DEFAULT 'all',
                         department VARCHAR(50)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

INSERT INTO notices (title, content, expire_time, target_type)
VALUES
    ('五一放假通知', '5月1日-5月5日放假，请各部门做好工作安排', '2024-05-06 00:00:00', 'all'),
    ('技术部例会通知', '本周五下午2点技术部全体会议', '2024-05-11 00:00:00', 'department');

-- 权限表
CREATE TABLE permissions (
                             perm_id SMALLINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
                             perm_code VARCHAR(50) UNIQUE NOT NULL COMMENT '权限标识符',
                             description VARCHAR(200)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

INSERT INTO permissions (perm_code, description)
VALUES
    ('USER_MANAGE', '用户管理权限'),
    ('SALARY_MANAGE', '薪资管理权限'),
    ('ATTENDANCE_VIEW', '考勤查看权限');

-- 角色表
CREATE TABLE roles (
                       role_id SMALLINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
                       role_name VARCHAR(50) UNIQUE NOT NULL,
                       description VARCHAR(200)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

INSERT INTO roles (role_name, description)
VALUES
    ('系统管理员', '拥有全部系统权限'),
    ('HR经理', '人力资源管理权限'),
    ('部门主管', '部门管理权限');

-- 角色权限表
CREATE TABLE role_permissions (
                                  role_id SMALLINT UNSIGNED,
                                  perm_id SMALLINT UNSIGNED,
                                  PRIMARY KEY (role_id, perm_id),
                                  FOREIGN KEY (role_id) REFERENCES roles(role_id),
                                  FOREIGN KEY (perm_id) REFERENCES permissions(perm_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

INSERT INTO role_permissions
VALUES
    (1, 1),
    (1, 2),
    (1, 3),
    (2, 1),
    (2, 3),
    (3, 3);

-- 用户角色表
CREATE TABLE user_roles (
                            user_id INT UNSIGNED,
                            role_id SMALLINT UNSIGNED,
                            PRIMARY KEY (user_id, role_id),
                            FOREIGN KEY (user_id) REFERENCES users(user_id),
                            FOREIGN KEY (role_id) REFERENCES roles(role_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

INSERT INTO user_roles
VALUES
    (1, 1),  -- admin01 系统管理员
    (3, 2),  -- emp002 HR经理
    (2, 3);  -- emp001 部门主管
