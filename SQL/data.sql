USE hrms;

-- 插入用户数据
INSERT INTO users (username, password_hash, user_type, department, position, hire_date, salary_base, is_active, created_at)
VALUES
    ('admin01', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsM8lE9lBOsl7iKTVKIUi', 'admin', '管理部', '系统管理员', '2022-01-01', 25000.00, TRUE, '2022-01-01 08:00:00'),
    ('emp001', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsM8lE9lBOsl7iKTVKIUi', 'employee', '技术部', 'Go开发工程师', '2023-01-15', 15000.00, TRUE, '2023-01-15 09:00:00'),
    ('emp002', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsM8lE9lBOsl7iKTVKIUi', 'employee', '人力资源部', 'HR主管', '2022-05-20', 18000.00, TRUE, '2022-05-20 09:00:00'),
    ('emp003', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsM8lE9lBOsl7iKTVKIUi', 'employee', '市场部', '市场专员', '2023-03-10', 12000.00, TRUE, '2023-03-10 09:00:00'),
    ('emp004', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsM8lE9lBOsl7iKTVKIUi', 'employee', '财务部', '会计', '2022-08-15', 14000.00, TRUE, '2022-08-15 09:00:00'),
    ('emp005', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsM8lE9lBOsl7iKTVKIUi', 'employee', '技术部', '前端开发工程师', '2023-02-20', 14000.00, TRUE, '2023-02-20 09:00:00'),
    ('app2024', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsM8lE9lBOsl7iKTVKIUi', 'applicant', NULL, NULL, NULL, NULL, TRUE, '2024-04-15 14:30:00'),
    ('app2025', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsM8lE9lBOsl7iKTVKIUi', 'applicant', NULL, NULL, NULL, NULL, TRUE, '2024-04-20 10:15:00');

-- 更新应聘者信息
UPDATE users
SET education = '南京大学 计算机科学 本科',
    work_experience = '腾讯 后端开发 2年',
    skills = 'Go, MySQL, Redis, Docker',
    resume_path = '/resumes/app2024.pdf'
WHERE username = 'app2024';

UPDATE users
SET education = '浙江大学 市场营销 硕士',
    work_experience = '阿里巴巴 市场专员 1年',
    skills = '数据分析, 市场调研, 英语流利',
    resume_path = '/resumes/app2025.pdf'
WHERE username = 'app2025';

-- 插入职位数据
INSERT INTO jobs (title, description, requirements, salary_range, post_date, expiration_date, status)
VALUES
    ('高级Go工程师', '负责公司核心业务系统的开发与维护，参与分布式系统架构设计，优化系统性能，保证系统稳定性。', '1. 3年以上Go语言开发经验\n2. 熟悉微服务架构，有大型分布式系统开发经验\n3. 精通MySQL、Redis等数据库\n4. 有良好的编码习惯和问题解决能力', '25k-40k', '2024-04-01 09:00:00', '2024-06-30', 'open'),
    ('HR实习生', '协助招聘工作，整理简历，安排面试，维护人事档案，参与员工培训活动组织。', '1. 人力资源相关专业在校生\n2. 熟练使用Office办公软件\n3. 良好的沟通能力和团队协作精神\n4. 责任心强，工作认真细致', '4k-6k', '2024-04-10 14:00:00', '2024-05-31', 'open'),
    ('前端开发工程师', '负责公司产品的Web前端开发，与后端团队协作实现产品功能，优化用户体验。', '1. 2年以上前端开发经验\n2. 精通HTML、CSS、JavaScript\n3. 熟练使用React或Vue框架\n4. 有良好的UI/UX感觉', '15k-25k', '2024-04-15 10:00:00', '2024-06-15', 'open');

-- 插入申请数据
INSERT INTO applications (user_id, job_id, apply_date, status)
VALUES
    (7, 1, '2024-04-16 10:30:00', 'pending'),
    (7, 2, '2024-04-17 14:20:00', 'interviewed'),
    (8, 3, '2024-04-21 09:45:00', 'pending');

-- 插入考勤数据
INSERT INTO attendances (user_id, clock_in, clock_out, status)
VALUES
    (2, '2024-05-01 08:58:00', '2024-05-01 18:02:00', 'normal'),
    (2, '2024-05-02 09:05:00', '2024-05-02 17:55:00', 'late'),
    (3, '2024-05-01 08:45:00', '2024-05-01 18:10:00', 'normal'),
    (3, '2024-05-02 08:50:00', '2024-05-02 17:30:00', 'early_leave');

-- 插入薪资数据
INSERT INTO salaries (user_id, month, base, bonus, deductions, payment_date)
VALUES
    (1, '202404', 25000.00, 5000.00, 4500.00, '2024-05-05'),
    (2, '202404', 15000.00, 3000.00, 2450.00, '2024-05-05'),
    (3, '202404', 18000.00, 2500.00, 3150.00, '2024-05-05');

-- 插入培训数据
INSERT INTO trainings (title, description, start_time, end_time, location, capacity)
VALUES
    ('Go语言性能优化', '深入讲解Go性能调优技巧，包括内存管理、并发控制、GC优化等内容，帮助开发人员提升系统性能。', '2024-06-01 14:00:00', '2024-06-01 17:00:00', '线上会议室', 50),
    ('新员工入职培训', '公司制度与文化培训，包括公司历史、组织架构、规章制度、企业文化等内容，帮助新员工快速融入团队。', '2024-05-15 09:00:00', '2024-05-16 17:00:00', '公司会议室', 30);

-- 插入培训记录数据
INSERT INTO training_records (user_id, training_id, enroll_time, status, score)
VALUES
    (2, 1, '2024-05-10 10:30:00', 'registered', NULL),
    (3, 2, '2024-05-05 14:20:00', 'completed', 92);

-- 插入公告数据
INSERT INTO notices (title, content, publish_time, expire_time, target_type, department)
VALUES
    ('五一放假通知', '根据国家法定节假日安排，公司定于5月1日-5月5日放假，5月6日正常上班。请各部门做好工作交接和安排，确保假期业务正常运转。', '2024-04-25 10:00:00', '2024-05-06 00:00:00', 'all', NULL),
    ('技术部例会通知', '本周五下午2点在技术部会议室召开技术部全体会议，请所有技术部成员准时参加，会议将讨论近期项目进展和技术难题。', '2024-05-08 09:30:00', '2024-05-11 00:00:00', 'department', '技术部');

-- 插入权限数据
INSERT INTO permissions (perm_code, description)
VALUES
    ('USER_MANAGE', '用户管理权限'),
    ('SALARY_MANAGE', '薪资管理权限'),
    ('ATTENDANCE_MANAGE', '考勤管理权限'),
    ('JOB_MANAGE', '岗位管理权限'),
    ('TRAINING_MANAGE', '培训管理权限'),
    ('NOTICE_MANAGE', '公告管理权限');

-- 插入角色数据
INSERT INTO roles (role_name, description)
VALUES
    ('系统管理员', '拥有全部系统权限'),
    ('HR经理', '人力资源管理权限'),
    ('部门主管', '部门管理权限'),
    ('普通员工', '基本权限');

-- 插入角色权限数据
INSERT INTO role_permissions (role_id, perm_id)
VALUES
    (1, 1), (1, 2), (1, 3), (1, 4), (1, 5), (1, 6),
    (2, 1), (2, 2), (2, 3), (2, 4), (2, 5),
    (3, 3), (3, 5), (3, 6);

-- 插入用户角色数据
INSERT INTO user_roles (user_id, role_id)
VALUES
    (1, 1),  -- admin01 系统管理员
    (2, 3),  -- emp001 部门主管
    (3, 2),  -- emp002 HR经理
    (4, 4),  -- emp003 普通员工
    (5, 4),  -- emp004 普通员工
    (6, 4);  -- emp005 普通员工
