from django.db import models

class User(models.Model):
    id = models.BigAutoField(primary_key=True, verbose_name='Id')
    name = models.CharField(blank=True, null=True, verbose_name='Имя', max_length=250)
    surname = models.CharField(blank=True, null=True, verbose_name='Фамилия', max_length=250)
    company_name = models.CharField(blank=True, null=True, verbose_name='Имя компании', max_length=250)
    company_address = models.CharField(blank=True, null=True, verbose_name='Адресс компании', max_length=250)
    inn = models.CharField(blank=True, null=True, verbose_name='ИНН', max_length=20)
    kpp = models.CharField(blank=True, null=True, verbose_name='КПП', max_length=20)
    phone = models.CharField(blank=False, null=False, verbose_name='Телефон', max_length=250)
    manager_name = models.CharField(blank=True, null=True, verbose_name='Имя менеджера', max_length=250)

    class Meta:
        managed = False
        db_table = 'user'
        verbose_name = 'Пользователи'
        verbose_name_plural = 'Пользователи'
        ordering = ['id']

    def __str__(self):
        if self.inn == "" and self.company_name == "" and self.phone != "":
            return f'Физ. лицо: Имя - {self.name} | Номер телефона - {self.phone}'
        elif self.inn != "" and self.company_name != "" and self.phone != "":
            return f'Юр. лицо: Компания - {self.company_name} | Номер телефона - {self.phone}'
        else:
            return f'Незарегестрированный пользователь: Id - {self.id}'

class Email(models.Model):
    id = models.BigAutoField(primary_key=True, verbose_name='Id')
    email = models.CharField(blank=True, null=True, verbose_name='email', max_length=250)
    can_to_send_news = models.BooleanField(blank=True, null=True, verbose_name='Разрешение на отправку акций')
    can_to_send_personal_offers = models.BooleanField(blank=True, null=True, verbose_name='Разрешение на отправку персональных предложений')
    user = models.ForeignKey('User', models.DO_NOTHING, blank=True, null=True, verbose_name='Id пользователя')

    class Meta:
        managed = False
        db_table = 'email'
        verbose_name = 'Email'
        verbose_name_plural = verbose_name
        ordering = ['id']

    def __str__(self):
        return f'Email - {self.email}'


class SiteReview(models.Model):
    rating = models.IntegerField(verbose_name='Рейтинг')
    comment = models.TextField(blank=True, null=True, verbose_name='Комментарий')
    created_at = models.DateTimeField(blank=True, null=True, verbose_name='Дата создания')

    class Meta:
        managed = False
        db_table = 'site_review'
        verbose_name = 'Отзыв сайта'
        verbose_name_plural = 'Отзывы сайта'
        ordering = ['rating']

    def __str__(self):
        return f'Рейтинг - {self.rating} | Комментарий - {self.comment}'


