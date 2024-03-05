from django.db import models


class Banner(models.Model):
    title = models.CharField(verbose_name='Заголовок', max_length=255, blank=True, null=True)
    link = models.URLField(verbose_name='Ссылка', blank=True, null=True)
    description = models.TextField(verbose_name='Описание', blank=True, null=True)
    image_right = models.ImageField(
        verbose_name='Изображение справа', max_length=255, blank=True, null=True, upload_to="photos/banner")

    class Meta:
        managed = False
        db_table = 'banner'
        verbose_name = 'Баннер на главной странице'
        verbose_name_plural = 'Баннер'

    def __str__(self):
        return f'Баннер - {self.title}'


class Blog(models.Model):
    title = models.CharField(verbose_name='Заголовок', max_length=255, blank=True, null=True)
    description = models.TextField(verbose_name='Описание', blank=True, null=True)
    short_description = models.TextField(verbose_name='Краткое описание', blank=True, null=True)
    image = models.ImageField(verbose_name='Фотография', max_length=255, blank=True, null=True, upload_to="photos/blog")
    date = models.DateTimeField(verbose_name='Дата', blank=False, null=False)

    class Meta:
        managed = False
        db_table = 'blog'
        verbose_name = 'Блог'
        verbose_name_plural = 'Блог'

    def __str__(self):
        return f'Блог - {self.title}'


class Vacancies(models.Model):
    title = models.CharField(verbose_name='Заголовок', max_length=255, blank=True, null=True)
    first_phone = models.CharField(verbose_name='Номер телефона 1', max_length=20, blank=True, null=True)
    second_phone = models.CharField(verbose_name='Номер телефона 2', max_length=20, blank=True, null=True)
    email = models.EmailField(verbose_name='Электронная почта', max_length=255, blank=True, null=True)

    class Meta:
        managed = False
        db_table = 'vacancy'
        verbose_name = 'Вакансии'
        verbose_name_plural = 'Вакансии'

    def __str__(self):
        return f'Вакансия - {self.title}'


class CurrentPromotions(models.Model):
    title = models.CharField(verbose_name='Заголовок', max_length=255, blank=True, null=True)
    description = models.TextField(verbose_name='Описание', blank=True, null=True)
    image = models.ImageField(
        verbose_name='Фотография', max_length=255, blank=True, null=True, upload_to="photos/curren_promotions")
    date = models.DateTimeField(blank=False, null=False)

    class Meta:
        managed = False
        db_table = 'current_promotions'
        verbose_name = 'Акции'
        verbose_name_plural = 'Акции'

    def __str__(self):
        return f'Акция - {self.title}'


class EmailStatic(models.Model):
    car_image = models.FileField(blank=True, null=True, verbose_name='Фото "Автомобиля"', upload_to="photos/email")
    cart_image = models.FileField(blank=True, null=True, verbose_name='Фото "Корзины"', upload_to="photos/email")
    like_image = models.FileField(blank=True, null=True, verbose_name='Фото "Лайка"', upload_to="photos/email")
    logo_image = models.FileField(blank=True, null=True, verbose_name='Фото "Лого"', upload_to="photos/email")
    courier_email = models.CharField(max_length=250, blank=False, null=False, verbose_name='Почта курьера')

    class Meta:
        managed = False
        db_table = 'email_static'
        verbose_name = 'Файл Email'
        verbose_name_plural = 'Файлы Email'
        ordering = ['id']

    def __str__(self):
        return f'Электронная почта - {self.courier_email}'


class PaymentMethod(models.Model):
    title = models.CharField(max_length=250, blank=True, null=True, verbose_name='Название')
    description = models.TextField(blank=True, null=True, verbose_name='Описание')
    icon = models.FileField(blank=True, null=True, verbose_name='Иконка', upload_to="photos/payment_methods/icons")
    image = models.FileField(blank=True, null=True, upload_to="photos/payment_methods/images", verbose_name='Фотография')

    class Meta:
        managed = False
        db_table = 'payment_method'
        verbose_name = 'Способ оплаты'
        verbose_name_plural = 'Способы оплаты'
        ordering = ['id']

    def __str__(self):
        return f'{self.title}'


class OrderStatus(models.Model):
    name = models.CharField(max_length=250, blank=True, null=True, verbose_name='Статус')

    class Meta:
        managed = False
        db_table = 'order_status'
        verbose_name = 'Статус заказа'
        verbose_name_plural = 'Статусы заказов'
        ordering = ['id']

    def __str__(self):
        return f'{self.name}'


class MetaTags(models.Model):
    title = models.TextField(verbose_name='Название')
    tag = models.TextField(verbose_name='Код')

    class Meta:
        managed = False
        db_table = 'meta_tags'
        verbose_name = 'Мета тэги'
        verbose_name_plural = 'Мета тэги'
        ordering = ['id']

    def __str__(self):
        return f'{self.title}'