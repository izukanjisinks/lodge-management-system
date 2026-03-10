-- Remove seeded demo data
DELETE FROM leave_balances WHERE employee_id IN (SELECT id FROM employees WHERE employee_number LIKE 'EMP-%');
DELETE FROM employee_documents WHERE employee_id IN (SELECT id FROM employees WHERE employee_number LIKE 'EMP-%');
DELETE FROM emergency_contacts WHERE employee_id IN (SELECT id FROM employees WHERE employee_number LIKE 'EMP-%');
DELETE FROM employees WHERE employee_number IN ('EMP-001','EMP-002','EMP-003','EMP-004','EMP-005','EMP-006','EMP-007','EMP-008','EMP-009','EMP-010','EMP-011','EMP-012','EMP-013','EMP-014','EMP-015');
DELETE FROM users WHERE email IN ('hr.manager@hr-system.com','eng.manager@hr-system.com','fin.manager@hr-system.com','sales.manager@hr-system.com','ops.manager@hr-system.com','alice.smith@hr-system.com','bob.jones@hr-system.com','carol.white@hr-system.com','david.brown@hr-system.com','eve.davis@hr-system.com','frank.miller@hr-system.com','grace.wilson@hr-system.com','henry.moore@hr-system.com','iris.taylor@hr-system.com','jack.anderson@hr-system.com');
DELETE FROM positions WHERE code IN ('SWE-1','SWE-2','SWE-SR','SWE-ST','FE-1','FE-SR','DO-1','ENG-MGR','HR-GEN','HR-TAS','HR-MGR','PAY-SP','FIN-ACC1','FIN-ACC2','FIN-MGR','SAL-REP','SAL-AE','MKT-SP','SAL-MGR','OPS-AN','OPS-MGR');
DELETE FROM departments WHERE code IN ('ENG-BE','ENG-FE','ENG-DO','HR-TA','HR-PAY','FIN-ACC','FIN-TRE','SAL-SLS','SAL-MKT','ENG','HR','FIN','SAL','OPS');
DELETE FROM leave_types WHERE code IN ('AL','SL','PL','UL','CL','ML');
