USE hrms;

-- 插入用户数据
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

-- 插入职位数据
INSERT INTO jobs (title, description, requirements, salary_range, expiration_date)
VALUES
    ('高级Go工程师', '负责分布式系统开发', '3年以上Go开发经验', '25k-40k', '2024-06-30'),
    ('HR实习生', '协助招聘工作', '人力资源专业优先', '4k-6k', '2024-05-31');

-- 插入申请数据
INSERT INTO applications (user_id, job_id, status)
VALUES
    (4, 1, 'pending'),
    (4, 2, 'interviewed');

-- 插入考勤数据
INSERT INTO attendances (user_id, clock_in, clock_out, status)
VALUES
    (2, '2024-05-01 08:58:00', '2024-05-01 18:02:00', 'normal'),
    (2, '2024-05-02 09:05:00', '2024-05-02 17:55:00', 'late');

-- 插入薪资数据
INSERT INTO salaries (user_id, month, base, bonus, deductions, payment_date)
VALUES
    (2, '202404', 15000.00, 3000.00, 2450.00, '2024-05-05'),
    (3, '202404', 18000.00, 2500.00, 3150.00, '2024-05-05');

-- 插入培训数据
INSERT INTO trainings (title, description, start_time, end_time, location, capacity)
VALUES
    ('Go语言性能优化', '深入讲解Go性能调优', '2024-06-01 14:00:00', '2024-06-01 17:00:00', '线上会议室', 50),
    ('新员工入职培训', '公司制度与文化培训', '2024-05-15 09:00:00', '2024-05-16 17:00:00', '公司会议室', 30);

-- 插入培训记录数据
INSERT INTO training_records (user_id, training_id, status)
VALUES
    (2, 1, 'registered'),
    (3, 2, 'completed');

-- 插入公告数据
INSERT INTO notices (title, content, expire_time, target_type)
VALUES
    ('五一放假通知', '5月1日-5月5日放假，请各部门做好工作安排', '2024-05-06 00:00:00', 'all'),
    ('技术部例会通知', '本周五下午2点技术部全体会议', '2024-05-11 00:00:00', 'department');

-- 插入权限数据
INSERT INTO permissions (perm_code, description)
VALUES
    ('USER_MANAGE', '用户管理权限'),
    ('SALARY_MANAGE', '薪资管理权限'),
    ('ATTENDANCE_VIEW', '考勤查看权限');

-- 插入角色数据
INSERT INTO roles (role_name, description)
VALUES
    ('系统管理员', '拥有全部系统权限'),
    ('HR经理', '人力资源管理权限'),
    ('部门主管', '部门管理权限');

-- 插入角色权限数据
INSERT INTO role_permissions
VALUES
    (1, 1),
    (1, 2),
    (1, 3),
    (2, 1),
    (2, 3),
    (3, 3);

-- 插入用户角色数据
INSERT INTO user_roles
VALUES
    (1, 1),  -- admin01 系统管理员
    (3, 2),  -- emp002 HR经理
    (2, 3);  -- emp001 部门主管
