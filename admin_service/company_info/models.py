from django.db import models


class FooterAddresses(models.Model):
    address = models.TextField(verbose_name='Адрес футера', blank=False, null=False)

    def __str__(self):
        return f'Адрес футера - {self.address}'

    class Meta:
        managed = False
        db_table = 'footer_addresses'
        verbose_name = 'Адреса в футере'
        verbose_name_plural = 'Адреса в футере'


class FooterObjects(models.Model):
    icon = models.ImageField(verbose_name='Иконка',blank=True , null=True, upload_to='photos/footer_objects/')
    link = models.URLField(verbose_name='Ссылка',blank=True , null=True)
    name = models.CharField(max_length=30, verbose_name='Название', blank=False, null=False)

    class Meta:
        managed = False
        db_table = 'footer_objects'
        verbose_name = 'Иконки в футере'
        verbose_name_plural = 'Иконки в футере'

    def __str__(self):
        return f'Иконка футера - {self.name}'


class HeaderPhones(models.Model):
    phone = models.CharField(verbose_name='Номер телефона', max_length=20, blank=False, null=False)

    def __str__(self):
        return f'Телефон хэдера - {self.phone}'

    class Meta:
        managed = False
        db_table = 'header_phones'
        verbose_name = 'Телефоны в хэдера'
        verbose_name_plural = 'Телефоны в хэдера'


class FooterPhones(models.Model):
    phone = models.CharField(verbose_name='Номер телефона', max_length=20, blank=False, null=False)

    def __str__(self):
        return f'Телефон футера - {self.phone}'

    class Meta:
        managed = False
        db_table = 'footer_phones'
        verbose_name = 'Телефоны в футере'
        verbose_name_plural = 'Телефоны в футере'


class PickUpPointTime(models.Model):
    mon = models.CharField(verbose_name="Понедельник",blank=True , null=True, max_length=20)
    tue = models.CharField(verbose_name="Вторник",blank=True , null=True, max_length=20)
    wen = models.CharField(verbose_name="Среда",blank=True , null=True, max_length=20)
    thu = models.CharField(verbose_name="Четверг",blank=True , null=True, max_length=20)
    fri = models.CharField(verbose_name="Пятница",blank=True , null=True, max_length=20)
    sat = models.CharField(verbose_name="Суббота",blank=True , null=True, max_length=20)
    sun = models.CharField(verbose_name="Воскресенье",blank=True , null=True, max_length=20)

    class Meta:
        managed = False
        db_table = 'pick_up_point_time'
        verbose_name = 'Режим работы филиала'
        verbose_name_plural = 'Режим работы филиала'

    def __str__(self):
        return f'Режим работы филиала - {self.id}'


class PickUpPoint(models.Model):
    phone1 = models.CharField(verbose_name='Номер телефона 1', max_length=20,blank=True , null=True)
    phone2 = models.CharField(verbose_name='Номер телефона 2', max_length=20,blank=True , null=True)
    phone3 = models.CharField(verbose_name='Номер телефона 3', max_length=20, blank=True, null=True)
    email1 = models.EmailField(verbose_name='Электронная почта 1', blank=True, null=True)
    email2 = models.EmailField(verbose_name='Электронная почта 2', blank=True, null=True)
    address = models.TextField(verbose_name='Адрес', blank=False, null=False)
    pick_up_point_stock_title = models.ForeignKey(
        'PickUpPointStockTitle', models.DO_NOTHING, blank=True, null=True, verbose_name='Акционный заголовок')
    pick_up_point_time = models.ForeignKey('PickUpPointTime', models.DO_NOTHING, blank=True, null=True, verbose_name='Время работы')
    coordinate_x = models.CharField(verbose_name='Координаты Х', max_length=255, blank=True, null=True)
    coordinate_y = models.CharField(verbose_name='Координаты Y', max_length=255, blank=True, null=True)

    class Meta:
        managed = False
        db_table = 'pick_up_point'
        verbose_name = 'Филиалы'
        verbose_name_plural = 'Филиалы'

    def __str__(self):
        return f'Филиал - {self.address}'

class PickUpPointStockTitle(models.Model):
    title = models.TextField(blank=False, null=False, verbose_name='Описание')

    class Meta:
        managed = False
        db_table = 'pick_up_point_stock_title'
        verbose_name = 'Уведомление филиала'
        verbose_name_plural = 'Уведомление филиалов'

    def __str__(self):
        return f'Заголовок - {self.title}'

class PickUpPointStockDescription(models.Model):
    id = models.BigAutoField(primary_key=True, verbose_name='id')
    description = models.CharField(max_length=300, blank=False, null=False, verbose_name='Описание')
    pick_up_point_stock_title = models.ForeignKey('PickUpPointStockTitle', models.DO_NOTHING, blank=True, null=True, verbose_name='Уведомление филиала')

    class Meta:
        managed = False
        db_table = 'pick_up_point_stock_description'
        verbose_name = 'Описание уведомления филиала'
        verbose_name_plural = 'Описание уведомления филиала'

    def __str__(self):
        return f'Описание - {self.description}'

class Requisites(models.Model):
    text = models.TextField(blank=False, null=False, verbose_name='Реквизиты')

    class Meta:
        managed = False
        db_table = 'requisites'
        verbose_name = 'Реквизиты'
        verbose_name_plural = 'Реквизиты'

    def __str__(self):
        return f'Реквизиты'


class PrivacyPolicy(models.Model):
    text = models.TextField(blank=False, null=False, verbose_name='Реквизиты')

    class Meta:
        managed = False
        db_table = 'privacy_policy'
        verbose_name = 'Политика конфиденциальности'
        verbose_name_plural = 'Политика конфиденциальности'

    def __str__(self):
        return f'Политика конфиденциальности'