from django.db import models
from django.core.exceptions import ValidationError


class Cart(models.Model):
    id = models.BigAutoField(primary_key=True, verbose_name='id')
    user = models.ForeignKey('auth_app.User', models.DO_NOTHING, blank=True, null=True, verbose_name='Пользователь')
    in_order = models.BooleanField(blank=True, null=True, verbose_name='Корзина в заказе')
    cart_product = models.ManyToManyField("CartProduct", through="Cartm2Ms", verbose_name='Продукты корзины')

    class Meta:
        managed = True
        db_table = 'cart'
        verbose_name = 'Корзина'
        verbose_name_plural = 'Корзины'

    def __str__(self):
        return f'Корзина - {self.id}'


class CartProduct(models.Model):
    id = models.BigAutoField(primary_key=True, verbose_name='id')
    product_uuid = models.ForeignKey('Product', models.DO_NOTHING, db_column='product_uuid', blank=True, null=True, verbose_name='id товара')
    count = models.BigIntegerField(blank=True, null=True, verbose_name='Кол-во')
    total_price = models.DecimalField(max_digits=65535, decimal_places=2, blank=True, null=True, verbose_name='Сумма')

    class Meta:
        managed = False
        db_table = 'cart_product'
        verbose_name = 'Продукты корзины'
        verbose_name_plural = 'Продукты корзины'

    def __str__(self):
        return f'Продукт - {self.product_uuid.title} | Кол-во - {self.count} | Итоговая цена - {self.total_price}'


class Cartm2Ms(models.Model):
    cart = models.ForeignKey(Cart, models.DO_NOTHING, primary_key=True)
    cart_product = models.ForeignKey(CartProduct, models.DO_NOTHING)

    class Meta:
        managed = False
        db_table = 'cartm2ms'
        unique_together = (('cart', 'cart_product'),)
        verbose_name = 'Продукты корзины'
        verbose_name_plural = 'Продукты корзины'

    def __str__(self):
        return f'Продукты'


class Category(models.Model):
    uuid = models.CharField(max_length=250, primary_key=True, verbose_name='UUID')
    title = models.CharField(max_length=250, blank=True, null=True, verbose_name='Название')
    image = models.ImageField(blank=True, null=True, upload_to="photos/categories", verbose_name='Фото')
    can_to_view = models.BooleanField(blank=True, null=True, verbose_name='Разрешено отображать на сайте')
    level = models.BigIntegerField(blank=True, null=True, verbose_name='Уровень')
    parent_uuid = models.CharField(max_length=250, blank=True, null=True, verbose_name='UUID отца')

    class Meta:
        db_table = 'category'
        verbose_name = 'Категория'
        verbose_name_plural = 'Категории'
        ordering = ['uuid']

    def __str__(self):
        return f'Категория - {self.title}'


class CategoryProduct(models.Model):
    category = models.ForeignKey('Category', db_column='category_uuid', on_delete=models.CASCADE)
    product = models.ForeignKey('Product', db_column='product_uuid', on_delete=models.CASCADE)

    class Meta:
        managed = False
        db_table = 'category_product'
        unique_together = (('category', 'product'),)


class Favourite(models.Model):
    user = models.ForeignKey('auth.User', models.DO_NOTHING, blank=True, null=True)

    class Meta:
        managed = False
        db_table = 'favourite'


class FavouriteProduct(models.Model):
    product_uuid = models.ForeignKey('Product', models.DO_NOTHING, db_column='product_uuid', blank=True, null=True)

    class Meta:
        managed = False
        db_table = 'favourite_product'


class Favouritesm2Ms(models.Model):
    favourite = models.OneToOneField(Favourite, models.DO_NOTHING, primary_key=True)
    favourite_product = models.ForeignKey(FavouriteProduct, models.DO_NOTHING)

    class Meta:
        managed = False
        db_table = 'favouritesm2ms'
        unique_together = (('favourite', 'favourite_product'),)


class Order(models.Model):
    id = models.BigAutoField(primary_key=True, verbose_name='id')
    user = models.ForeignKey('auth_app.User', models.DO_NOTHING, blank=True, null=True, verbose_name='Пользователь')
    order_status = models.ForeignKey('cms.OrderStatus', models.DO_NOTHING, blank=True, null=True, verbose_name='Статус')
    cart = models.ForeignKey(Cart, models.DO_NOTHING, blank=True, null=True, verbose_name='Корзина')
    created_at = models.DateTimeField(blank=True, null=True, verbose_name='Дата создания')
    payment_method = models.ForeignKey('cms.PaymentMethod', models.DO_NOTHING, blank=True, null=True, verbose_name='Способ оплаты')
    total_price = models.DecimalField(max_digits=65535, decimal_places=2, blank=True, null=True, verbose_name='Итоговая цена')
    cancel = models.BooleanField(blank=True, null=True, verbose_name='Отменен')
    promo = models.BooleanField(blank=True, null=True, verbose_name='Использование промокода')

    class Meta:
        managed = False
        db_table = 'order'
        verbose_name = 'Заказ'
        verbose_name_plural = 'Заказы'

    def __str__(self):
        return f'{self.id}'


class Product(models.Model):
    uuid = models.CharField(max_length=250, primary_key=True, verbose_name='UUID')
    title = models.CharField(max_length=250, blank=True, null=True, verbose_name='Название')
    description = models.TextField(blank=True, null=True, verbose_name='Описание')
    vendor_code = models.CharField(max_length=250, blank=True, null=True, verbose_name='Вендор-код')
    base_unit = models.CharField(max_length=250, blank=True, null=True, verbose_name='Единица измерения')
    image = models.ImageField(blank=True, null=True, upload_to="photos/products", verbose_name='Фото')
    price = models.DecimalField(max_digits=10, decimal_places=2, blank=True, null=True, verbose_name='Цена')
    can_to_view = models.BooleanField(blank=True, null=True, verbose_name='Разрешено отображать на сайте')
    category = models.ManyToManyField('Category', through='CategoryProduct', through_fields=('product', 'category'), verbose_name='Категории')

    class Meta:
        db_table = 'product'
        verbose_name = 'Продукт'
        verbose_name_plural = 'Продукты'
        ordering = ['uuid']

    def __str__(self):
        return f'Продукт - {self.title}'


class ProductStatistic(models.Model):
    rating = models.DecimalField(max_digits=2, decimal_places=2, blank=True, null=True, verbose_name='Рейтинг')
    sales_count = models.BigIntegerField(blank=True, null=True, verbose_name='Кол-во продаж')
    request_detail_count = models.BigIntegerField(blank=True, null=True, verbose_name='Кол-во запросов')
    product_uuid = models.ForeignKey(Product, models.DO_NOTHING, db_column='product_uuid', blank=True, null=True, verbose_name='Продукт')

    class Meta:
        managed = False
        db_table = 'product_statistic'
        verbose_name = 'Статистика продукта'
        verbose_name_plural = 'Статистика продуктов'
        ordering = ['id']

    def __str__(self):
        return f'{self.product_uuid}'


class PromoCode(models.Model):
    promo_code = models.CharField(max_length=250, blank=False, null=False, verbose_name='Промокод')
    discount_percent = models.BigIntegerField(blank=True, null=True, verbose_name='Скидка в %')
    discount_sum = models.BigIntegerField(blank=True, null=True, verbose_name='Скидка в р.')
    number_of_uses = models.BigIntegerField(blank=False, null=False, verbose_name='Кол-во использований')
    expires_at = models.DateTimeField(blank=False, null=False, verbose_name='Дата истичения действия')
    created_at = models.DateTimeField(blank=False, null=False, verbose_name='Дата создания')

    class Meta:
        managed = False
        db_table = 'promo_code'
        verbose_name = 'Промокод'
        verbose_name_plural = 'Промокоды'
        ordering = ['id']

    def clean(self):
        if not self.discount_sum and not self.discount_percent:
            raise ValidationError('Необходимо заполнить одно из полей "Скидка в %" или "Скидка в р."')
        if self.discount_sum and self.discount_percent:
            raise ValidationError('Можно заполнить только одно из полей "Скидка в %" или "Скидка в р."')

    def __str__(self):
        return f'Промокод №{self.id}'
    

class ProductFile(models.Model):
    product_files = models.ForeignKey('ProductFiles', models.DO_NOTHING, primary_key=True)
    product_uuid = models.ForeignKey(Product, models.DO_NOTHING, db_column='product_uuid')

    class Meta:
        managed = False
        db_table = 'product_file'
        unique_together = (('product_files', 'product_uuid'),)
        verbose_name = 'Файл продукта'
        verbose_name_plural = 'Файлы продуктов'

    def __str__(self):
        return f'Продукт'


class ProductFiles(models.Model):
    id = models.BigAutoField(primary_key=True, verbose_name='Id')
    title = models.CharField(max_length=255, blank=True, null=True, verbose_name='Название')
    document = models.FileField(blank=True, null=True, upload_to="photos/products", verbose_name='Файл')
    products = models.ManyToManyField("Product", through="ProductFile", verbose_name='Продукты')

    class Meta:
        managed = False
        db_table = 'product_files'
        verbose_name = 'Файл продукта'
        verbose_name_plural = 'Файлы продуктов'
        ordering = ['id']

    def __str__(self):
        return f'Файл - {self.title}'