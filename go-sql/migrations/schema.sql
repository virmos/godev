CREATE TABLE departments (
   id INT PRIMARY 							KEY     			NOT NULL,
   "name"           						CHAR(255)   		NOT NULL,
   parent_id								INT,
   constraint sr_fk_depart_depart 	foreign key (parent_id)		references departments(id),
   created_at	  							TIMESTAMP
);
ALTER TABLE departments ALTER COLUMN created_at SET DEFAULT null;
ALTER TABLE departments DROP CONSTRAINT sr_fk_depart_depart;
ALTER TABLE departments ADD CONSTRAINT sr_fk_depart_depart FOREIGN KEY (parent_id) REFERENCES departments(id) ON DELETE RESTRICT ON UPDATE RESTRICT;

alter table departments alter column created_at set default now();
CREATE TRIGGER set_timestamp
    BEFORE UPDATE ON departments
    FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp()


CREATE TABLE employees (
   id 			  INT PRIMARY KEY    		NOT NULL,
   "name"         CHAR(255)    				NOT NULL,
   first_name     CHAR(255)    				NOT NULL,
   last_name      CHAR(255)    				NOT NULL,
   email       	  CHAR(255),
   department_id  INT,
   created_at	  TIMESTAMP
);
ALTER TABLE employees ALTER COLUMN created_at SET DEFAULT null;
ALTER TABLE employees ADD CONSTRAINT sr_fk_depart_emp FOREIGN KEY (department_id) REFERENCES departments(id) ON DELETE RESTRICT ON UPDATE RESTRICT;

alter table employees alter column created_at set default now();
CREATE TRIGGER set_timestamp
    BEFORE UPDATE ON employees
    FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp()



-------------------------INSERT--------------------------------
INSERT INTO departments (id, "name", parent_id, created_at) VALUES (1, 'Depart 1', , NULL);
INSERT INTO departments (id, "name", parent_id, created_at) VALUES (1, 'Depart 2', , NULL);
INSERT INTO departments (id, "name", parent_id, created_at) VALUES (1, 'Depart 1', , NULL);
INSERT INTO departments (id, "name", parent_id, created_at) VALUES (1, 'Depart 1', , NULL);
INSERT INTO departments (id, "name", parent_id, created_at) VALUES (1, 'Depart 1', , NULL);
INSERT INTO departments (id, "name", parent_id, created_at) VALUES (1, 'Depart 1', , NULL);
INSERT INTO departments (id, "name", parent_id, created_at) VALUES (1, 'Depart 1', , NULL);

INSERT INTO employees (id, "name", first_name, last_name, email, created_at) VALUES (1, 'Klee', 'Klee', 'Genshin', 'Klee@gmail.com', NULL);
INSERT INTO employees (id, "name", first_name, last_name, email, created_at) VALUES (2, 'Collei', 'Collei', 'Genshin', 'Collei@gmail.com', NULL);
INSERT INTO employees (id, "name", first_name, last_name, email, created_at) VALUES (3, 'Xinyan', 'Xinyan', 'Genshin', 'Xinyan@gmail.com', NULL);
INSERT INTO employees (id, "name", first_name, last_name, email, created_at) VALUES (4, 'Xiao', 'Xiao', 'Genshin', 'Xiao@gmail.com', NULL);
INSERT INTO employees (id, "name", first_name, last_name, email, created_at) VALUES (5, 'Xinyan', 'Xinyan', 'Genshin', 'Xinyan@gmail.com', NULL);
INSERT INTO employees (id, "name", first_name, last_name, email, created_at) VALUES (6, 'Qiqi', 'Qiqi', 'Genshin', 'Qiqi@gmail.com', NULL);
INSERT INTO employees (id, "name", first_name, last_name, email, created_at) VALUES (7, 'Diluc', 'Diluc', 'Genshin', 'Diluc@gmail.com', NULL);
INSERT INTO employees (id, "name", first_name, last_name, email, created_at) VALUES (8, 'Yoimiya', 'Yoimiya', 'Genshin', 'Yoimiya@gmail.com', NULL);



-------------------------DELETE--------------------------------
ALTER TABLE departments DROP CONSTRAINT sr_fk_depart_depart;
ALTER TABLE employees DROP CONSTRAINT sr_fk_depart_emp;
drop table departments;
drop table employees;






























