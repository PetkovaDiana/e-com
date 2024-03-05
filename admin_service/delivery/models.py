from django.db import models


class CdekDelivery(models.Model):
    id = models.BigAutoField(primary_key=True, verbose_name='id')
    delivery_type = models.ForeignKey('DeliveryType', models.DO_NOTHING, blank=True, null=True, verbose_name='Тип доставки')
    pick_up_point_address = models.CharField(max_length=100, blank=True, null=True, verbose_name='Адрес пункта выдачи')

    class Meta:
        managed = False
        db_table = 'cdek_delivery'
        verbose_name = 'CDEK'
        verbose_name_plural = 'CDEK'

    def __str__(self):
        return f'CDEK доставка - {self.id}'


class CourierDelivery(models.Model):
    address = models.TextField(verbose_name='Адрес', blank=True, null=True)
    apartment_office = models.CharField(verbose_name='Офис / квартира', blank=True, null=True, max_length=30)
    index = models.CharField(verbose_name='Индекс', blank=True, null=True, max_length=30)
    entrance = models.CharField(verbose_name='Подъезд', blank=True, null=True, max_length=30)
    intercom = models.CharField(verbose_name='Домофон', blank=True, null=True, max_length=30)
    floor = models.CharField(verbose_name='Этаж', blank=True, null=True, max_length=10)
    note = models.TextField(verbose_name='Комментарий', blank=True, null=True)
    delivery_type = models.ForeignKey(
        'DeliveryType', models.DO_NOTHING, blank=True, null=True, verbose_name='Номер заказа')

    class Meta:
        managed = False
        db_table = 'courier_delivery'
        verbose_name = 'Курьерская доставка'
        verbose_name_plural = 'Курьерская доствка'

    def __str__(self):
        return f'Курьерская доставка - {self.id}'


class DeliveryTypeInfo(models.Model):
    title = models.CharField(verbose_name='Название', max_length=30, blank=False, null=False)
    description = models.CharField(verbose_name='Описание', max_length=255, null=True, blank=True)
    icon = models.ImageField(verbose_name='Иконка', upload_to='photos/delivery_types/')
    can_delivery = models.BooleanField(verbose_name='Возможность доставки')
    delivery_price = models.IntegerField(verbose_name='Стоимость доставки ₽', max_length=10)


    class Meta:
        managed = False
        db_table = 'delivery_type_info'
        verbose_name = 'Виды доставки'
        verbose_name_plural = 'Виды доставки'

    def __str__(self):
        return f'Доставка - {self.title}'


class DeliveryType(models.Model):
    id = models.BigAutoField(primary_key=True)
    order = models.ForeignKey('catalog.Order', models.DO_NOTHING, blank=True, null=True, verbose_name='Заказ')

    class Meta:
        managed = False
        db_table = 'delivery_type'
        verbose_name = 'Способ доставки'
        verbose_name_plural = 'Способ доставки'

    def __str__(self):
        return f'Доставка - {self.order.id}'


class SelfDelivery(models.Model):
    delivery_type = models.ForeignKey(
        DeliveryType, models.DO_NOTHING, blank=True, null=True, verbose_name='Номер заказа')
    pick_up_point = models.ForeignKey(
        'company_info.PickUpPoint', models.DO_NOTHING, blank=True, null=True, verbose_name='Адрес филиала')

    class Meta:
        managed = False
        db_table = 'self_delivery'
        verbose_name = 'Самовывоз'
        verbose_name_plural = 'Самовывоз'


    def __str__(self):
        return f'Самовывоз - {self.delivery_type.order.id}'
    
class CdekDeliveryInfo(models.Model):
    description = models.TextField(blank=False, null=False, verbose_name='Описание')

    class Meta:
        managed = False
        db_table = 'cdek_delivery_info'
        verbose_name = 'Информация o доставке CDEK'
        verbose_name_plural = 'Информация o доставке CDEK'

    def __str__(self):
        return f'Информация о доставке CDEK - {self.description}'


class CourierDeliveryInfo(models.Model):
    description = models.TextField(blank=True, null=True, verbose_name='Описание')
    courier_delivery_time_info = models.ForeignKey('CourierDeliveryTimeInfo', models.DO_NOTHING, blank=True, null=True, verbose_name='Рассписание курьера')

    class Meta:
        managed = False
        db_table = 'courier_delivery_info'
        verbose_name = 'Информация o доставке курьером'
        verbose_name_plural = 'Информация o доставке курьером'

    def __str__(self):
        return f'Информация о доставке курьером'


class CourierDeliveryTimeInfo(models.Model):
    mon = models.CharField(max_length=30, blank=True, null=True, verbose_name='Понедельник')
    tue = models.CharField(max_length=30, blank=True, null=True, verbose_name='Вторник')
    wen = models.CharField(max_length=30, blank=True, null=True, verbose_name='Среда')
    thu = models.CharField(max_length=30, blank=True, null=True, verbose_name='Четверг')
    fri = models.CharField(max_length=30, blank=True, null=True, verbose_name='Пятница')
    sat = models.CharField(max_length=30, blank=True, null=True, verbose_name='Суббота')
    sun = models.CharField(max_length=30, blank=True, null=True, verbose_name='Воскресенье')

    class Meta:
        managed = False
        db_table = 'courier_delivery_time_info'
        verbose_name = 'Рассписание курьера'
        verbose_name_plural = 'Рассписание курьера'

    def __str__(self):
        return 'Информация о рассписании курьера'
