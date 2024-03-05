from email.mime.multipart import MIMEMultipart
from email.mime.application import MIMEApplication
from os.path import basename
import smtplib

import os

async def send_email(receiver_email: str, subject: str, file_path: str):
    '''Отправка html письма'''
    sender_email = os.environ.get('SENDER_EMAIL')
    password = os.environ.get('SENDER_EMAIL_PASSWORD')

    message = MIMEMultipart("alternative")
    message["Subject"] = subject
    message["From"] = sender_email
    message["To"] = receiver_email

    with open(file_path, "rb") as fil:
        part = MIMEApplication(
            fil.read(),
            Name=basename(subject + '.pdf')
        )
    # After the file is closed
    part['Content-Disposition'] = 'attachment; filename="%s"' % basename(subject + '.pdf')
    message.attach(part)

    # Создание безопасного подключения с сервером и отправка сообщения
    server = smtplib.SMTP_SSL('smtp.mail.ru')
    server.set_debuglevel(1)
    server.ehlo()
    server.login(sender_email, password)
    server.auth_plain()

    server.sendmail(sender_email, receiver_email, message.as_string())
    server.close()