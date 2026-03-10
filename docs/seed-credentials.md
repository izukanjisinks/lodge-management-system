# Seed User Credentials

This document lists all user accounts created by the seed script (`make seed`). These credentials are for **development and testing only**.

## Admin Users

| Email | Password | Role | Description |
|-------|----------|------|-------------|
| admin@hr-system.com | `Admin@123` | Super Admin | Full system access |
| hr.manager@hr-system.com | `HrManager@123` | HR Manager | Manage employees, leave, payroll, recruitment |

## Manager Users

| Email | Password | Role | Department |
|-------|----------|------|------------|
| eng.manager@hr-system.com | `Manager@123` | Manager | Engineering |
| fin.manager@hr-system.com | `Manager@123` | Manager | Finance |
| sales.manager@hr-system.com | `Manager@123` | Manager | Sales & Marketing |
| ops.manager@hr-system.com | `Manager@123` | Manager | Operations |

## Employee Users

All employee accounts use the same password: **`Employee@123`**

| Email | Employee Name | Department | Position |
|-------|---------------|------------|----------|
| alice.smith@hr-system.com | Alice Smith | Backend Engineering | Software Engineer II |
| bob.jones@hr-system.com | Bob Jones | Backend Engineering | Senior Software Engineer |
| carol.white@hr-system.com | Carol White | Backend Engineering | Software Engineer I |
| david.brown@hr-system.com | David Brown | Frontend Engineering | Senior Frontend Engineer |
| eve.davis@hr-system.com | Eve Davis | Frontend Engineering | Frontend Engineer |
| frank.miller@hr-system.com | Frank Miller | Talent Acquisition (HR) | Talent Acquisition Specialist |
| grace.wilson@hr-system.com | Grace Wilson | Payroll (HR) | Payroll Specialist |
| henry.moore@hr-system.com | Henry Moore | Accounting (Finance) | Senior Accountant |
| iris.taylor@hr-system.com | Iris Taylor | Sales | Account Executive |
| jack.anderson@hr-system.com | Jack Anderson | Marketing (Sales) | Marketing Specialist |

## Quick Reference

```bash
# Super Admin
Email:    admin@hr-system.com
Password: Admin@123

# HR Manager
Email:    hr.manager@hr-system.com
Password: HrManager@123

# Department Managers
Email:    {eng|fin|sales|ops}.manager@hr-system.com
Password: Manager@123

# Regular Employees
Email:    {firstname}.{lastname}@hr-system.com
Password: Employee@123
```

## Security Notes

⚠️ **WARNING:** These are development credentials only. In production:
- Use strong, unique passwords for each user
- Enforce password complexity requirements
- Implement password rotation policies
- Enable multi-factor authentication
- Never commit real credentials to version control
