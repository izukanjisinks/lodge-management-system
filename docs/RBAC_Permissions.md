# Role-Based Access Control (RBAC)

Lodge Management System — Permission Reference
Last updated: 2026-03-17

---

## Roles

| Role | Description |
|---|---|
| **Admin** | Full system access. Manages users, configuration, and all data. |
| **Manager** | Oversees operations. Can approve bookings, view reports, and manage rooms. |
| **Receptionist** | Front-desk staff. Handles bookings, clients, and invoices. Read-only on rooms. |
| **Cleaner** | Housekeeping staff. Has a dedicated simplified dashboard showing only their assigned rooms and cleaning schedule. |

---

## Page Access

| Page | Admin | Manager | Receptionist | Cleaner |
|---|:---:|:---:|:---:|:---:|
| Dashboard (staff) | ✅ | ✅ | ✅ | ❌ |
| Dashboard (cleaner) | ❌ | ❌ | ❌ | ✅ |
| Rooms | ✅ | ✅ | ✅ | ✅ (read-only) |
| Bookings | ✅ | ✅ | ✅ | ❌ |
| Meals | ✅ | ✅ | ✅ | ❌ |
| Invoices | ✅ | ✅ | ✅ | ❌ |
| Individual Clients | ✅ | ✅ | ✅ | ❌ |
| Corporate Clients | ✅ | ✅ | ✅ | ❌ |
| Reports | ✅ | ✅ | ❌ | ❌ |
| Workflow Editor | ✅ | ❌ | ❌ | ❌ |
| Task Inbox | ✅ | ✅ | ✅ | ❌ |
| System Users | ✅ | ❌ | ❌ | ❌ |

---

## Action-Level Permissions

### Rooms
| Action | Admin | Manager | Receptionist | Cleaner |
|---|:---:|:---:|:---:|:---:|
| View rooms | ✅ | ✅ | ✅ | ✅ |
| Create / Edit room | ✅ | ✅ | ❌ | ❌ |
| Delete room | ✅ | ✅ | ❌ | ❌ |
| View cleaning assignments | ✅ | ✅ | ❌ | ✅ |
| Manage cleaning assignments | ✅ | ✅ | ❌ | ❌ |

### Bookings
| Action | Admin | Manager | Receptionist | Cleaner |
|---|:---:|:---:|:---:|:---:|
| View bookings | ✅ | ✅ | ✅ | ❌ |
| Create / Edit booking | ✅ | ✅ | ✅ | ❌ |
| Delete booking | ✅ | ✅ | ❌ | ❌ |
| Approve / Reject booking | ✅ | ✅ | ❌ | ❌ |

### Invoices
| Action | Admin | Manager | Receptionist | Cleaner |
|---|:---:|:---:|:---:|:---:|
| View invoices | ✅ | ✅ | ✅ | ❌ |
| Download invoice PDF | ✅ | ✅ | ✅ | ❌ |
| Update invoice status | ✅ | ✅ | ❌ | ❌ |
| Delete invoice | ✅ | ✅ | ❌ | ❌ |

### Clients
| Action | Admin | Manager | Receptionist | Cleaner |
|---|:---:|:---:|:---:|:---:|
| View clients | ✅ | ✅ | ✅ | ❌ |
| Create / Edit client | ✅ | ✅ | ✅ | ❌ |
| Delete client | ✅ | ✅ | ❌ | ❌ |

### System Users
| Action | Admin | Manager | Receptionist | Cleaner |
|---|:---:|:---:|:---:|:---:|
| View users | ✅ | ❌ | ❌ | ❌ |
| Create / Edit user | ✅ | ❌ | ❌ | ❌ |
| Delete user | ✅ | ❌ | ❌ | ❌ |

### Workflow
| Action | Admin | Manager | Receptionist | Cleaner |
|---|:---:|:---:|:---:|:---:|
| View workflow editor | ✅ | ❌ | ❌ | ❌ |
| Edit workflow steps / transitions | ✅ | ❌ | ❌ | ❌ |
| View task inbox | ✅ | ✅ | ✅ | ❌ |
| Approve / Reject tasks | ✅ | ✅ | ❌ | ❌ |

### Reports
| Action | Admin | Manager | Receptionist | Cleaner |
|---|:---:|:---:|:---:|:---:|
| View reports & analytics | ✅ | ✅ | ❌ | ❌ |

---

## Dev Login Credentials

> For development/testing only. Remove before production.

| Email | Password | Role |
|---|---|---|
| admin@lodge.dev | admin123 | Admin |
| manager@lodge.dev | manager123 | Manager |
| receptionist@lodge.dev | reception123 | Receptionist |
| cleaner@lodge.dev | cleaner123 | Cleaner |

---

## Implementation

Permission checks are enforced in two places:

1. **Navigation guard** — `src/composables/usePermissions.ts` maps each route to its allowed roles. Clicking a restricted sidebar item shows an "Access Denied" dialog instead of navigating.

2. **UI visibility** — Action buttons (delete, edit, approve, etc.) are conditionally rendered based on role using the same composable.

> Unauthenticated users are redirected to `/login` by the global router guard in `src/router/index.ts`.
