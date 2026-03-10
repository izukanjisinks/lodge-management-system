package email

import "fmt"

const loginURL = "http://localhost:5173/login"

// loginButton returns the reusable HTML login button block
func loginButton() string {
	return fmt.Sprintf(`
              <table width="100%%%%" cellpadding="0" cellspacing="0" style="margin:30px 0;">
                <tr>
                  <td align="center">
                    <a href="%s"
                       style="display:inline-block; padding:14px 28px; background-color:#1e3c72; color:#ffffff; text-decoration:none; font-size:15px; font-weight:bold; border-radius:6px;">
                       Login Now
                    </a>
                  </td>
                </tr>
              </table>`, loginURL)
}

// LeaveRequestAssignedTemplate generates HTML for leave request assignment notification
func LeaveRequestAssignedTemplate(requesterName, leaveType string, days int, startDate, endDate string) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8" />
<meta name="viewport" content="width=device-width, initial-scale=1.0" />
<title>Leave Request Assigned</title>
</head>
<body style="margin:0; padding:0; background-color:#f4f6f9; font-family: Arial, Helvetica, sans-serif;">

  <table width="100%%" cellpadding="0" cellspacing="0" style="background-color:#f4f6f9; padding:40px 0;">
    <tr>
      <td align="center">

        <table width="600" cellpadding="0" cellspacing="0" style="background:#ffffff; border-radius:8px; box-shadow:0 4px 12px rgba(0,0,0,0.08); overflow:hidden;">

          <tr>
            <td style="background:linear-gradient(135deg, #1e3c72, #2a5298); padding:30px; text-align:center;">
              <h1 style="color:#ffffff; margin:0; font-size:22px; letter-spacing:0.5px;">
                Leave Request Assigned for Review
              </h1>
            </td>
          </tr>

          <tr>
            <td style="padding:40px 30px; color:#333333; font-size:15px; line-height:1.6;">

              <p style="margin-top:0;">Hello,</p>

              <p>A leave request has been assigned to you for review.</p>

              <table width="100%%" cellpadding="0" cellspacing="0" style="margin:25px 0;">
                <tr>
                  <td style="background-color:#f0f4ff; border:1px solid #d6e0f0; padding:20px; border-radius:6px;">
                    <table width="100%%" cellpadding="0" cellspacing="0">
                      <tr>
                        <td style="padding:6px 0; font-weight:bold; color:#555; width:120px;">Employee:</td>
                        <td style="padding:6px 0; color:#333;">%s</td>
                      </tr>
                      <tr>
                        <td style="padding:6px 0; font-weight:bold; color:#555;">Leave Type:</td>
                        <td style="padding:6px 0; color:#333;">%s</td>
                      </tr>
                      <tr>
                        <td style="padding:6px 0; font-weight:bold; color:#555;">Duration:</td>
                        <td style="padding:6px 0; color:#333;">%d days</td>
                      </tr>
                      <tr>
                        <td style="padding:6px 0; font-weight:bold; color:#555;">Period:</td>
                        <td style="padding:6px 0; color:#333;">%s to %s</td>
                      </tr>
                    </table>
                  </td>
                </tr>
              </table>

              <p>Please log in to the HR System to review and take action on this request.</p>

              %s

              <p style="margin-bottom:0;">
                Best regards,<br/>
                <strong>HR System</strong>
              </p>

            </td>
          </tr>

          <tr>
            <td style="background-color:#f8f9fc; padding:20px; text-align:center; font-size:12px; color:#888;">
              This is an automated message from the HR System. Please do not reply to this email.
            </td>
          </tr>

        </table>

      </td>
    </tr>
  </table>

</body>
</html>
`, requesterName, leaveType, days, startDate, endDate, loginButton())
}

// LeaveRequestApprovedTemplate generates HTML for leave request approval notification
func LeaveRequestApprovedTemplate(employeeName, leaveType string, days int, startDate, endDate, approverName string) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8" />
<meta name="viewport" content="width=device-width, initial-scale=1.0" />
<title>Leave Request Approved</title>
</head>
<body style="margin:0; padding:0; background-color:#f4f6f9; font-family: Arial, Helvetica, sans-serif;">

  <table width="100%%" cellpadding="0" cellspacing="0" style="background-color:#f4f6f9; padding:40px 0;">
    <tr>
      <td align="center">

        <table width="600" cellpadding="0" cellspacing="0" style="background:#ffffff; border-radius:8px; box-shadow:0 4px 12px rgba(0,0,0,0.08); overflow:hidden;">

          <tr>
            <td style="background:linear-gradient(135deg, #1e6d3a, #28a745); padding:30px; text-align:center;">
              <h1 style="color:#ffffff; margin:0; font-size:22px; letter-spacing:0.5px;">
                Leave Request Approved
              </h1>
            </td>
          </tr>

          <tr>
            <td style="padding:40px 30px; color:#333333; font-size:15px; line-height:1.6;">

              <p style="margin-top:0;">Hello %s,</p>

              <p>Great news! Your leave request has been approved.</p>

              <table width="100%%" cellpadding="0" cellspacing="0" style="margin:25px 0;">
                <tr>
                  <td style="background-color:#eafbef; border:1px solid #c3e6cb; padding:20px; border-radius:6px;">
                    <table width="100%%" cellpadding="0" cellspacing="0">
                      <tr>
                        <td style="padding:6px 0; font-weight:bold; color:#555; width:120px;">Leave Type:</td>
                        <td style="padding:6px 0; color:#333;">%s</td>
                      </tr>
                      <tr>
                        <td style="padding:6px 0; font-weight:bold; color:#555;">Duration:</td>
                        <td style="padding:6px 0; color:#333;">%d days</td>
                      </tr>
                      <tr>
                        <td style="padding:6px 0; font-weight:bold; color:#555;">Period:</td>
                        <td style="padding:6px 0; color:#333;">%s to %s</td>
                      </tr>
                      <tr>
                        <td style="padding:6px 0; font-weight:bold; color:#555;">Approved by:</td>
                        <td style="padding:6px 0; color:#333;">%s</td>
                      </tr>
                    </table>
                  </td>
                </tr>
              </table>

              <p>Your leave has been officially approved. Enjoy your time off!</p>

              %s

              <p style="margin-bottom:0;">
                Best regards,<br/>
                <strong>HR System</strong>
              </p>

            </td>
          </tr>

          <tr>
            <td style="background-color:#f8f9fc; padding:20px; text-align:center; font-size:12px; color:#888;">
              This is an automated message from the HR System. Please do not reply to this email.
            </td>
          </tr>

        </table>

      </td>
    </tr>
  </table>

</body>
</html>
`, employeeName, leaveType, days, startDate, endDate, approverName, loginButton())
}

// LeaveRequestRejectedTemplate generates HTML for leave request rejection notification
func LeaveRequestRejectedTemplate(employeeName, leaveType string, days int, startDate, endDate, reviewerName, reason string) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8" />
<meta name="viewport" content="width=device-width, initial-scale=1.0" />
<title>Leave Request Not Approved</title>
</head>
<body style="margin:0; padding:0; background-color:#f4f6f9; font-family: Arial, Helvetica, sans-serif;">

  <table width="100%%" cellpadding="0" cellspacing="0" style="background-color:#f4f6f9; padding:40px 0;">
    <tr>
      <td align="center">

        <table width="600" cellpadding="0" cellspacing="0" style="background:#ffffff; border-radius:8px; box-shadow:0 4px 12px rgba(0,0,0,0.08); overflow:hidden;">

          <tr>
            <td style="background:linear-gradient(135deg, #c0392b, #e74c3c); padding:30px; text-align:center;">
              <h1 style="color:#ffffff; margin:0; font-size:22px; letter-spacing:0.5px;">
                Leave Request Not Approved
              </h1>
            </td>
          </tr>

          <tr>
            <td style="padding:40px 30px; color:#333333; font-size:15px; line-height:1.6;">

              <p style="margin-top:0;">Hello %s,</p>

              <p>We regret to inform you that your leave request has not been approved.</p>

              <table width="100%%" cellpadding="0" cellspacing="0" style="margin:25px 0;">
                <tr>
                  <td style="background-color:#fff8e1; border:1px solid #ffe0b2; padding:20px; border-radius:6px;">
                    <table width="100%%" cellpadding="0" cellspacing="0">
                      <tr>
                        <td style="padding:6px 0; font-weight:bold; color:#555; width:120px;">Leave Type:</td>
                        <td style="padding:6px 0; color:#333;">%s</td>
                      </tr>
                      <tr>
                        <td style="padding:6px 0; font-weight:bold; color:#555;">Duration:</td>
                        <td style="padding:6px 0; color:#333;">%d days</td>
                      </tr>
                      <tr>
                        <td style="padding:6px 0; font-weight:bold; color:#555;">Period:</td>
                        <td style="padding:6px 0; color:#333;">%s to %s</td>
                      </tr>
                      <tr>
                        <td style="padding:6px 0; font-weight:bold; color:#555;">Reviewed by:</td>
                        <td style="padding:6px 0; color:#333;">%s</td>
                      </tr>
                      %s
                    </table>
                  </td>
                </tr>
              </table>

              <p>If you have any questions, please contact your manager or HR department.</p>

              %s

              <p style="margin-bottom:0;">
                Best regards,<br/>
                <strong>HR System</strong>
              </p>

            </td>
          </tr>

          <tr>
            <td style="background-color:#f8f9fc; padding:20px; text-align:center; font-size:12px; color:#888;">
              This is an automated message from the HR System. Please do not reply to this email.
            </td>
          </tr>

        </table>

      </td>
    </tr>
  </table>

</body>
</html>
`, employeeName, leaveType, days, startDate, endDate, reviewerName, formatReason(reason), loginButton())
}

func formatReason(reason string) string {
	if reason != "" {
		return fmt.Sprintf(`
                      <tr>
                        <td style="padding:6px 0; font-weight:bold; color:#555; width:120px;">Reason:</td>
                        <td style="padding:6px 0; color:#333;">%s</td>
                      </tr>`, reason)
	}
	return ""
}

// GenericTaskAssignedTemplate generates HTML for generic task assignment notification
func GenericTaskAssignedTemplate(recipientName, taskName, taskDescription string) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8" />
<meta name="viewport" content="width=device-width, initial-scale=1.0" />
<title>New Task Assigned</title>
</head>
<body style="margin:0; padding:0; background-color:#f4f6f9; font-family: Arial, Helvetica, sans-serif;">

  <table width="100%%" cellpadding="0" cellspacing="0" style="background-color:#f4f6f9; padding:40px 0;">
    <tr>
      <td align="center">

        <table width="600" cellpadding="0" cellspacing="0" style="background:#ffffff; border-radius:8px; box-shadow:0 4px 12px rgba(0,0,0,0.08); overflow:hidden;">

          <tr>
            <td style="background:linear-gradient(135deg, #1e3c72, #2a5298); padding:30px; text-align:center;">
              <h1 style="color:#ffffff; margin:0; font-size:22px; letter-spacing:0.5px;">
                New Task Assigned
              </h1>
            </td>
          </tr>

          <tr>
            <td style="padding:40px 30px; color:#333333; font-size:15px; line-height:1.6;">

              <p style="margin-top:0;">Hello %s,</p>

              <p>A new task has been assigned to you for your action.</p>

              <table width="100%%" cellpadding="0" cellspacing="0" style="margin:25px 0;">
                <tr>
                  <td style="background-color:#f0f4ff; border:1px solid #d6e0f0; padding:20px; border-radius:6px;">
                    <table width="100%%" cellpadding="0" cellspacing="0">
                      <tr>
                        <td style="padding:6px 0; font-weight:bold; color:#555; width:120px;">Task:</td>
                        <td style="padding:6px 0; color:#333;">%s</td>
                      </tr>
                      <tr>
                        <td style="padding:6px 0; font-weight:bold; color:#555;">Description:</td>
                        <td style="padding:6px 0; color:#333;">%s</td>
                      </tr>
                    </table>
                  </td>
                </tr>
              </table>

              <p>Please log in to the HR System to review and take action on this task.</p>

              %s

              <p style="margin-bottom:0;">
                Best regards,<br/>
                <strong>HR System</strong>
              </p>

            </td>
          </tr>

          <tr>
            <td style="background-color:#f8f9fc; padding:20px; text-align:center; font-size:12px; color:#888;">
              This is an automated message from the HR System. Please do not reply to this email.
            </td>
          </tr>

        </table>

      </td>
    </tr>
  </table>

</body>
</html>
`, recipientName, taskName, taskDescription, loginButton())
}

// PasswordResetTemplate generates HTML for password reset notification
func PasswordResetTemplate(temporaryPassword string) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8" />
<meta name="viewport" content="width=device-width, initial-scale=1.0" />
<title>Password Reset</title>
</head>
<body style="margin:0; padding:0; background-color:#f4f6f9; font-family: Arial, Helvetica, sans-serif;">

  <table width="100%%" cellpadding="0" cellspacing="0" style="background-color:#f4f6f9; padding:40px 0;">
    <tr>
      <td align="center">

        <table width="600" cellpadding="0" cellspacing="0" style="background:#ffffff; border-radius:8px; box-shadow:0 4px 12px rgba(0,0,0,0.08); overflow:hidden;">

          <tr>
            <td style="background:linear-gradient(135deg, #1e3c72, #2a5298); padding:30px; text-align:center;">
              <h1 style="color:#ffffff; margin:0; font-size:22px; letter-spacing:0.5px;">
                Password Reset Notification
              </h1>
            </td>
          </tr>

          <tr>
            <td style="padding:40px 30px; color:#333333; font-size:15px; line-height:1.6;">

              <p style="margin-top:0;">Hello,</p>

              <p>Your password has been <strong>reset by an administrator</strong>.</p>

              <table width="100%%" cellpadding="0" cellspacing="0" style="margin:25px 0;">
                <tr>
                  <td style="background-color:#f0f4ff; border:1px dashed #2a5298; padding:20px; text-align:center; border-radius:6px;">
                    <p style="margin:0; font-size:14px; color:#555;">Your Temporary Password</p>
                    <p style="margin:8px 0 0; font-size:20px; font-weight:bold; letter-spacing:2px; color:#1e3c72;">
                      %s
                    </p>
                  </td>
                </tr>
              </table>

              <p>You will be required to change this password when you next log in.</p>

              %s

              <p style="color:#a94442; background:#fdecea; padding:12px; border-radius:4px; font-size:14px;">
                If you did not expect this change, please contact your administrator immediately.
              </p>

              <p style="margin-bottom:0;">
                Best regards,<br/>
                <strong>HR System</strong>
              </p>

            </td>
          </tr>

          <tr>
            <td style="background-color:#f8f9fc; padding:20px; text-align:center; font-size:12px; color:#888;">
              This is an automated message from the HR System. Please do not reply to this email.
            </td>
          </tr>

        </table>

      </td>
    </tr>
  </table>

</body>
</html>
`, temporaryPassword, loginButton())
}

// WelcomeEmployeeTemplate generates HTML for new employee welcome notification
func WelcomeEmployeeTemplate(firstName, lastName, userEmail, password string) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8" />
<meta name="viewport" content="width=device-width, initial-scale=1.0" />
<title>Welcome to HR System</title>
</head>
<body style="margin:0; padding:0; background-color:#f4f6f9; font-family: Arial, Helvetica, sans-serif;">

  <table width="100%%" cellpadding="0" cellspacing="0" style="background-color:#f4f6f9; padding:40px 0;">
    <tr>
      <td align="center">

        <table width="600" cellpadding="0" cellspacing="0" style="background:#ffffff; border-radius:8px; box-shadow:0 4px 12px rgba(0,0,0,0.08); overflow:hidden;">

          <tr>
            <td style="background:linear-gradient(135deg, #1e3c72, #2a5298); padding:30px; text-align:center;">
              <h1 style="color:#ffffff; margin:0; font-size:22px; letter-spacing:0.5px;">
                Welcome to HR System!
              </h1>
            </td>
          </tr>

          <tr>
            <td style="padding:40px 30px; color:#333333; font-size:15px; line-height:1.6;">

              <p style="margin-top:0;">Hello %s %s,</p>

              <p>Welcome aboard! Your employee account has been successfully created in our HR System.</p>

              <table width="100%%" cellpadding="0" cellspacing="0" style="margin:25px 0;">
                <tr>
                  <td style="background-color:#eafbef; border:1px solid #c3e6cb; padding:15px; border-radius:6px; text-align:center; font-weight:bold; color:#1e6d3a;">
                    Your account is now active!
                  </td>
                </tr>
              </table>

              <p>You can now access the system using the following credentials:</p>

              <table width="100%%" cellpadding="0" cellspacing="0" style="margin:25px 0;">
                <tr>
                  <td style="background-color:#f0f4ff; border:1px solid #d6e0f0; padding:20px; border-radius:6px;">
                    <table width="100%%" cellpadding="0" cellspacing="0">
                      <tr>
                        <td style="padding:6px 0; font-weight:bold; color:#555; width:160px;">Email:</td>
                        <td style="padding:6px 0; color:#333; font-family:monospace; background-color:#f4f4f4; padding:2px 6px; border-radius:3px;">%s</td>
                      </tr>
                      <tr>
                        <td style="padding:6px 0; font-weight:bold; color:#555;">Temporary Password:</td>
                        <td style="padding:6px 0; color:#333; font-family:monospace; background-color:#f4f4f4; padding:2px 6px; border-radius:3px;">%s</td>
                      </tr>
                    </table>
                  </td>
                </tr>
              </table>

              <p style="color:#856404; background:#fff3cd; padding:12px; border-radius:4px; font-size:14px;">
                <strong>Important Security Notice:</strong> For your security, you will be required to change this temporary password upon your first login. Please choose a strong password that meets our security requirements.
              </p>

              %s

              <p>If you have any questions or need assistance, please contact the HR department.</p>

              <p>We look forward to working with you!</p>

              <p style="margin-bottom:0;">
                Best regards,<br/>
                <strong>HR System</strong>
              </p>

            </td>
          </tr>

          <tr>
            <td style="background-color:#f8f9fc; padding:20px; text-align:center; font-size:12px; color:#888;">
              This is an automated message from the HR System. Please do not reply to this email.
            </td>
          </tr>

        </table>

      </td>
    </tr>
  </table>

</body>
</html>
`, firstName, lastName, userEmail, password, loginButton())
}

// PayslipReadyTemplate generates HTML for payslip ready notification
func PayslipReadyTemplate(firstName, period string) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8" />
<meta name="viewport" content="width=device-width, initial-scale=1.0" />
<title>Payslip Ready</title>
</head>
<body style="margin:0; padding:0; background-color:#f4f6f9; font-family: Arial, Helvetica, sans-serif;">

  <table width="100%%%%" cellpadding="0" cellspacing="0" style="background-color:#f4f6f9; padding:40px 0;">
    <tr>
      <td align="center">

        <table width="600" cellpadding="0" cellspacing="0" style="background:#ffffff; border-radius:8px; box-shadow:0 4px 12px rgba(0,0,0,0.08); overflow:hidden;">

          <tr>
            <td style="background:linear-gradient(135deg, #1e6d3a, #28a745); padding:30px; text-align:center;">
              <h1 style="color:#ffffff; margin:0; font-size:22px; letter-spacing:0.5px;">
                Your Payslip is Ready
              </h1>
            </td>
          </tr>

          <tr>
            <td style="padding:40px 30px; color:#333333; font-size:15px; line-height:1.6;">

              <p style="margin-top:0;">Hi %s,</p>

              <p>Your payslip for <strong>%s</strong> has been processed and is now available.</p>

              <p>Log in to the HR portal to view your full payslip breakdown.</p>

              <p style="margin-bottom:0;">
                Best regards,<br/>
                <strong>HR System</strong>
              </p>

            </td>
          </tr>

          <tr>
            <td style="background-color:#f8f9fc; padding:20px; text-align:center; font-size:12px; color:#888;">
              This is an automated message from the HR System. Please do not reply to this email.
            </td>
          </tr>

        </table>

      </td>
    </tr>
  </table>

</body>
</html>
`, firstName, period)
}
