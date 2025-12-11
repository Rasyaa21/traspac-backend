package utils

import (
	"fmt"
	"gin-backend-app/internal/dto/request"
)

func BuildVerificationEmailHTML(data *request.EmailData) (string, error) {
    html := fmt.Sprintf(`
    <!DOCTYPE html>
    <html lang="en">
    <head>
        <meta charset="UTF-8">
        <meta name="viewport" content="width=device-width, initial-scale=1.0">
        <title>Email Verification</title>
        <style>
            body { 
                margin: 0; 
                padding: 0; 
                font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif; 
                background-color: #f5f5f5; 
            }
            .container { 
                max-width: 600px; 
                margin: 0 auto; 
                background-color: white; 
                border-radius: 10px; 
                overflow: hidden; 
                box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1); 
            }
            .header { 
                background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%); 
                padding: 30px; 
                text-align: center; 
            }
            .header h1 { 
                color: white; 
                margin: 0; 
                font-size: 28px; 
                font-weight: 600; 
            }
            .content { 
                padding: 40px 30px; 
                text-align: center; 
            }
            .content p { 
                color: #333; 
                font-size: 16px; 
                line-height: 1.6; 
                margin: 20px 0; 
            }
            .otp-code { 
                display: inline-block; 
                padding: 20px 30px; 
                background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%); 
                color: white; 
                font-size: 32px; 
                font-weight: bold; 
                letter-spacing: 8px; 
                border-radius: 15px; 
                margin: 30px 0; 
                border: 3px solid #f0f0f0;
                box-shadow: 0 4px 15px rgba(102, 126, 234, 0.3);
            }
            .otp-info { 
                background-color: #e3f2fd; 
                border: 1px solid #2196f3; 
                border-radius: 8px; 
                padding: 15px; 
                margin: 20px 0; 
                color: #1565c0; 
            }
            .footer { 
                background-color: #f8f9fa; 
                padding: 20px; 
                text-align: center; 
                color: #666; 
                font-size: 14px; 
            }
            .icon { 
                width: 60px; 
                height: 60px; 
                background: rgba(255, 255, 255, 0.2); 
                border-radius: 50%%; 
                display: inline-flex; 
                align-items: center; 
                justify-content: center; 
                margin-bottom: 20px; 
                font-size: 24px;
            }
        </style>
    </head>
    <body>
        <div class="container">
            <div class="header">
                <div class="icon">✉️</div>
                <h1>%s</h1>
            </div>
            <div class="content">
                <p>%s</p>
                <div class="otp-info">
                    <strong>Your verification code is:</strong>
                </div>
                <div class="otp-code">%s</div>
                <div class="otp-info">
                    <strong>Important:</strong> This code will expire in 10 minutes for your security.
                    <br>Please do not share this code with anyone.
                </div>
                <p style="margin-top: 30px; font-size: 14px; color: #666;">
                    Enter this code in the verification form to complete your email verification.
                </p>
            </div>
            <div class="footer">
                <p>If you didn't request this verification, please ignore this email.</p>
                <p>&copy; 2025 Your App Name. All rights reserved.</p>
            </div>
        </div>
    </body>
    </html>
    `, data.Title, data.Message, data.OTPCode)

    return html, nil
}

func BuildResetPasswordEmailHTML(data *request.EmailData) (string, error) {
    html := fmt.Sprintf(`
    <!DOCTYPE html>
    <html lang="en">
    <head>
        <meta charset="UTF-8">
        <meta name="viewport" content="width=device-width, initial-scale=1.0">
        <title>Password Reset</title>
        <style>
            body { 
                margin: 0; 
                padding: 0; 
                font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif; 
                background-color: #f5f5f5; 
            }
            .container { 
                max-width: 600px; 
                margin: 0 auto; 
                background-color: white; 
                border-radius: 10px; 
                overflow: hidden; 
                box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1); 
            }
            .header { 
                background: linear-gradient(135deg, #ff6b6b 0%%, #ee5a52 100%%); 
                padding: 30px; 
                text-align: center; 
            }
            .header h1 { 
                color: white; 
                margin: 0; 
                font-size: 28px; 
                font-weight: 600; 
            }
            .content { 
                padding: 40px 30px; 
                text-align: center; 
            }
            .content p { 
                color: #333; 
                font-size: 16px; 
                line-height: 1.6; 
                margin: 20px 0; 
            }
            .otp-code { 
                display: inline-block; 
                padding: 20px 30px; 
                background: linear-gradient(135deg, #ff6b6b 0%%, #ee5a52 100%%); 
                color: white; 
                font-size: 32px; 
                font-weight: bold; 
                letter-spacing: 8px; 
                border-radius: 15px; 
                margin: 30px 0; 
                border: 3px solid #f0f0f0;
                box-shadow: 0 4px 15px rgba(255, 107, 107, 0.3);
            }
            .warning { 
                background-color: #fff3cd; 
                border: 1px solid #ffeaa7; 
                border-radius: 8px; 
                padding: 15px; 
                margin: 20px 0; 
                color: #856404; 
            }
            .footer { 
                background-color: #f8f9fa; 
                padding: 20px; 
                text-align: center; 
                color: #666; 
                font-size: 14px; 
            }
            .icon { 
                width: 60px; 
                height: 60px; 
                background: rgba(255, 255, 255, 0.2); 
                border-radius: 50%%; 
                display: inline-flex; 
                align-items: center; 
                justify-content: center; 
                margin-bottom: 20px; 
                font-size: 24px;
            }
        </style>
    </head>
    <body>
        <div class="container">
            <div class="header">
                <div class="icon">🔐</div>
                <h1>%s</h1>
            </div>
            <div class="content">
                <p>%s</p>
                <div class="warning">
                    <strong>Your password reset code is:</strong>
                </div>
                <div class="otp-code">%s</div>
                <div class="warning">
                    <strong>Security Notice:</strong> This code will expire in 10 minutes for your security.
                    <br>Please do not share this code with anyone.
                </div>
                <p style="margin-top: 30px; font-size: 14px; color: #666;">
                    Enter this code in the password reset form to continue.
                </p>
            </div>
            <div class="footer">
                <p>If you didn't request this password reset, please ignore this email or contact support.</p>
                <p>&copy; 2025 Your App Name. All rights reserved.</p>
            </div>
        </div>
    </body>
    </html>
    `, data.Title, data.Message, data.OTPCode)

    return html, nil
}
