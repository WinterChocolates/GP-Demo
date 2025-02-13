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

-- 考勤表
CREATE TABLE attendances (
                             attendance_id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
                             user_id INT UNSIGNED NOT NULL,
                             clock_in DATETIME NOT NULL,
                             clock_out DATETIME,
                             status ENUM('normal', 'late', 'early_leave') DEFAULT 'normal',
                             date_clock_in DATE AS (DATE(clock_in)) STORED,
                             INDEX idx_user_date (user_id, date_clock_in),
                             FOREIGN KEY (user_id) REFERENCES users(user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

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

-- 权限表
CREATE TABLE permissions (
                             perm_id SMALLINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
                             perm_code VARCHAR(50) UNIQUE NOT NULL COMMENT '权限标识符',
                             description VARCHAR(200)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 角色表
CREATE TABLE roles (
                       role_id SMALLINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
                       role_name VARCHAR(50) UNIQUE NOT NULL,
                       description VARCHAR(200)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 角色权限表
CREATE TABLE role_permissions (
                                  role_id SMALLINT UNSIGNED,
                                  perm_id SMALLINT UNSIGNED,
                                  PRIMARY KEY (role_id, perm_id),
                                  FOREIGN KEY (role_id) REFERENCES roles(role_id),
                                  FOREIGN KEY (perm_id) REFERENCES permissions(perm_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 用户角色表
CREATE TABLE user_roles (
                            user_id INT UNSIGNED,
                            role_id SMALLINT UNSIGNED,
                            PRIMARY KEY (user_id, role_id),
                            FOREIGN KEY (user_id) REFERENCES users(user_id),
                            FOREIGN KEY (role_id) REFERENCES roles(role_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
